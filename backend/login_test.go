package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Dass33/administratum/backend/internal/auth"
	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func TestReturnLoginData(t *testing.T) {
	// Setup test database
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "file:test.db"
	}

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Initialize database schema
	err = initializeTestDatabase(db)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	// Create database queries
	dbQueries := database.New(db)

	// Setup test configuration
	cfg := &apiConfig{
		db:       dbQueries,
		platform: PlatformDev,
		jwt_key:  "test-jwt-key-for-testing-purposes-only",
	}

	// Create test context
	ctx := context.Background()

	// Create test data in the database
	testUser, testTable, _ := createTestDataInDB(t, ctx, dbQueries)

	// Test the ReturnLoginData function
	t.Run("ReturnLoginData with complete data", func(t *testing.T) {
		// Create a test response writer
		w := httptest.NewRecorder()

		// Call the function
		cfg.ReturnLoginData(w, testUser, ctx, 200)

		// Check response status
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		// Parse the response
		var response LoginData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Log the response for debugging
		log.Printf("=== ReturnLoginData Test Response ===")
		log.Printf("Email: %s", response.Email)
		log.Printf("Token: %s", response.Token)
		log.Printf("OpenedTable ID: %s", response.OpenedTable.ID)
		log.Printf("OpenedTable Name: %s", response.OpenedTable.Name)
		log.Printf("OpenedSheet ID: %s", response.OpenedSheet.ID)
		log.Printf("OpenedSheet Name: %s", response.OpenedSheet.Name)
		log.Printf("Number of table names: %d", len(response.TableIdNames))
		log.Printf("Number of columns in sheet: %d", len(response.OpenedSheet.Columns))
		log.Printf("Number of branches in table: %d", len(response.OpenedTable.BranchesNames))

		// Verify response fields
		if response.Email != testUser.Email {
			t.Errorf("Expected email %s, got %s", testUser.Email, response.Email)
		}

		if response.Token == "" {
			t.Errorf("Expected non-empty token")
		}

		// Check cookies
		cookies := w.Result().Cookies()
		foundRefreshCookie := false
		for _, cookie := range cookies {
			if cookie.Name == auth.RefreshTokenName {
				foundRefreshCookie = true
				if cookie.Value == "" {
					t.Errorf("Expected non-empty refresh token cookie value")
				}
				break
			}
		}
		if !foundRefreshCookie {
			t.Errorf("Expected refresh token cookie to be set")
		}

		log.Printf("=== Test completed successfully ===")
	})

	// Test with user that has no opened sheet
	t.Run("ReturnLoginData with no opened sheet", func(t *testing.T) {
		// Create a user without opened sheet
		userNoSheet := testUser
		userNoSheet.OpenedSheet = nil

		w := httptest.NewRecorder()
		cfg.ReturnLoginData(w, userNoSheet, ctx, 200)

		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var response LoginData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		log.Printf("=== ReturnLoginData Test (No Opened Sheet) ===")
		log.Printf("Email: %s", response.Email)
		log.Printf("Token: %s", response.Token)
		log.Printf("OpenedTable ID: %s", response.OpenedTable.ID)
		log.Printf("OpenedSheet ID: %s", response.OpenedSheet.ID)
		log.Printf("Number of table names: %d", len(response.TableIdNames))

		// Should still have table names but empty opened table/sheet
		if len(response.TableIdNames) != 1 {
			t.Errorf("Expected 1 table name, got %d", len(response.TableIdNames))
		}

		log.Printf("=== No opened sheet test completed successfully ===")
	})

	// Clean up test data
	cleanupTestData(t, ctx, dbQueries, testUser.ID, testTable.ID)
}

func createTestDataInDB(t *testing.T, ctx context.Context, db *database.Queries) (database.User, database.Table, database.Sheet) {
	// Create test user
	hashedPassword, err := auth.HashPassword("testpassword123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	userParams := database.CreateUserParams{
		Email:          "test@example.com",
		HashedPassword: hashedPassword,
	}
	user, err := db.CreateUser(ctx, userParams)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test table
	table, err := db.CreateTable(ctx, sql.NullString{String: "https://example.com/game", Valid: true})
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Create test branch
	branchParams := database.CreateBranchParams{
		Name:    "Test Branch",
		TableID: table.ID,
	}
	branch, err := db.CreateBranch(ctx, branchParams)
	if err != nil {
		t.Fatalf("Failed to create test branch: %v", err)
	}

	// Create test sheet
	sheetParams := database.CreateSheetParams{
		Name:     "Test Sheet",
		RowCount: 10,
		BranchID: branch.ID,
	}
	sheet, err := db.CreateSheet(ctx, sheetParams)
	if err != nil {
		t.Fatalf("Failed to create test sheet: %v", err)
	}

	// Update user with opened sheet
	updateUserOpenedSheetParams := database.UpdateUserOpenedSheetParams{
		ID:          user.ID,
		OpenedSheet: &sheet.ID,
	}
	user, err = db.UpdateUserOpenedSheet(ctx, updateUserOpenedSheetParams)
	if err != nil {
		t.Fatalf("Failed to update user with opened sheet: %v", err)
	}

	// Create user table relationship
	userTableParams := database.CreateUserTableParams{
		UserID:     user.ID,
		TableID:    table.ID,
		Permission: "read_write",
	}
	_, err = db.CreateUserTable(ctx, userTableParams)
	if err != nil {
		t.Fatalf("Failed to create user table relationship: %v", err)
	}

	log.Printf("=== Test Data Setup ===")
	log.Printf("User ID: %s", user.ID)
	log.Printf("User Email: %s", user.Email)
	log.Printf("Table ID: %s", table.ID)
	log.Printf("Table Name: %s", table.Name)
	log.Printf("Branch ID: %s", branch.ID)
	log.Printf("Branch Name: %s", branch.Name)
	log.Printf("Sheet ID: %s", sheet.ID)
	log.Printf("Sheet Name: %s", sheet.Name)
	log.Printf("Opened Sheet ID: %s", *user.OpenedSheet)

	return user, table, sheet
}

func initializeTestDatabase(db *sql.DB) error {
	// Read and execute schema files
	schemaFiles := []string{
		"sql/schema/001_users.sql",
		"sql/schema/002_refresh.sql",
		"sql/schema/003_tables.sql",
		"sql/schema/004_opened_table.sql",
		"sql/schema/005_rename_rows.sql",
		"sql/schema/006_tables_name.sql",
		"sql/schema/007_opened_sheet.sql",
		"sql/schema/008_column_data_index.sql",
		"sql/schema/009_column_data_column_id.sql",
		"sql/schema/010_remove_sheet_id.sql",
	}

	for _, file := range schemaFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		_, err = db.Exec(string(content))
		if err != nil {
			return err
		}
	}

	return nil
}

func cleanupTestData(t *testing.T, ctx context.Context, db *database.Queries, userID, tableID uuid.UUID) {
	// Clean up test data
	// In a real implementation, you would delete the test data from the database
	log.Printf("Cleaning up test data for user %s and table %s", userID, tableID)
}

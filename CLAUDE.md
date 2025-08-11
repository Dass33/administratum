# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Administratum is a full-stack web application that makes configuration management for web apps more accessible. It provides a spreadsheet-like interface for managing structured data with version control through branches.

**Architecture**: 
- **Backend**: Go REST API using Chi router, SQLC for type-safe SQL queries, Goose for migrations, and Turso (LibSQL) database
- **Frontend**: React/TypeScript SPA with Vite build system, TailwindCSS, and Material-UI components
- **Deployment**: Backend on Google Cloud Run, frontend on GitHub Pages

## Key Concepts

The application manages hierarchical data structures:
- **Tables/Projects**: Top-level containers representing different configurations
- **Branches**: Version control for tables, allowing different variants (similar to git branches)  
- **Sheets**: Collections of columns within a branch (different data categories)
- **Columns**: Individual data fields with types (text, number, bool, etc.)
- **Column Data**: The actual values stored in columns

## Development Commands

### Backend (Go)
```bash
# From backend/ directory
go run .                    # Run development server
go test ./...              # Run all tests
go build -o administratum  # Build binary

# Database migrations (requires .env with DATABASE_URL)
../scripts/migrateup.sh    # Apply migrations
../scripts/migratedown.sh  # Rollback migrations

# Production build
../scripts/buildprod.sh    # Cross-compile for Linux
```

### Frontend (React/TypeScript)
```bash
# From frontend/ directory  
npm run dev        # Development server (Vite)
npm run build      # Production build (TypeScript + Vite)
npm run lint       # ESLint
npm run preview    # Preview production build
npm run deploy     # Deploy to GitHub Pages
```

## Architecture Details

### Backend Structure
- **main.go**: Entry point, configures Chi router and CORS for allowed origins
- **internal/database/**: SQLC-generated database layer with type-safe queries
- **sql/queries/**: Raw SQL queries used by SQLC
- **sql/schema/**: Database migrations managed by Goose
- **Handler files** (*.go): HTTP endpoint implementations for CRUD operations
- **Authentication**: JWT-based auth with refresh tokens, bcrypt password hashing

### Frontend Structure  
- **App.tsx**: Main component with resizable panels and modal management
- **AppContext.tsx**: Global state management using React Context
- **Component files**: Modular React components for UI sections
- **Authentication flow**: JWT tokens with refresh token rotation

### Database Schema
Key entities: Users, Tables, Branches, Sheets, Columns, ColumnData, UserTable (permissions), RefreshTokens

### Authentication & Permissions
- JWT access tokens with refresh token rotation
- User-table permissions (read/write) stored in UserTable junction table  
- Middleware authentication on protected endpoints
- CORS configured for localhost:5173 (dev) and GitHub Pages (prod)

## Environment Setup

**Backend** requires `.env` file:
- `PORT`: Server port
- `DATABASE_URL`: Turso database connection string  
- `PLATFORM`: "dev" or "production"
- `JWT_KEY`: Secret for JWT signing

## Code Generation

- **SQLC**: `sqlc generate` in backend/ regenerates database layer from SQL queries
- SQL queries in `sql/queries/` generate type-safe Go functions in `internal/database/`
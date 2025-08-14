# 📜 Administratum

A full-stack web application that makes configuration management for web apps more accessible through a spreadsheet-like interface with built-in version control.

## 🎯 Overview

Administratum provides a user-friendly way to manage structured configuration data:

- **Spreadsheet Interface**: Intuitive table-based data management
- **Branching**: Branch-based workflow similar to Git for configuration variants
- **Type Safety**: Strongly typed columns (text, number, boolean, etc.)
- **Game View**: Users can see changes to the config in real time

## 🏗️ Architecture

**Frontend**: React/TypeScript SPA and TailwindCSS styling
**Backend**: Go REST API with Chi router and type-safe database queries
**Database**: Turso with automated migrations
**Deployment**: Backend on Google Cloud Run, frontend on GitHub Pages

## 🛠️ Technology Stack

### Frontend
- **TypeScript**
- **React**
- **Vite**
- **TailwindCSS**

### Backend
- **Go**
- **Chi**
- **SQLC**
- **Goose**
- **Turso**

### Infrastructure
- **Google Cloud Run**
- **GitHub Pages**

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+

### Backend Setup
```bash
cd backend/
cp .env.example .env  # Configure your environment variables
go mod download
../scripts/migrateup.sh    # Apply database migrations
go run .                   # Start development server
```

### Frontend Setup
```bash
cd frontend/
npm install
npm run dev    # Start development server at http://localhost:5173
```

## 📁 Project Structure

```
administratum/
├── backend/           # Go API
│   ├── internal/      # Internal Go packages
│   ├── sql/           # Database queries and migrations
│   └── *.go           # HTTP handlers and main application
├── frontend/
│   ├── src/           # React TypeScript SPA
│   └── public/
└── scripts/           # Build and deployment scripts
```

## 🔧 Development Commands

### Backend
```bash
go run .                    # Development server
go build -o administratum  # Build binary
../scripts/buildprod.sh    # Production build
```

### Frontend
```bash
npm run dev        # Development server
npm run build      # Production build
npm run deploy     # Deploy to GitHub Pages
```

### Database
```bash
../scripts/migrateup.sh    # Apply migrations
../scripts/migratedown.sh  # Rollback migrations
sqlc generate              # Regenerate database layer
```

## 🔐 Authentication

- JWT-based authentication with refresh token rotation
- User-table permissions system (read/write access)

## 🎒 Roadmap

- [ ] Fix visual bugs
- [ ] Add JSON importing functionality

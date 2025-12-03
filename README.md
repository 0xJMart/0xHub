# 0xHub - Project Hub with CRD Operator

A modern project hub that displays projects from Kubernetes Custom Resource Definitions (CRDs) through a React frontend, Go backend API, and Kubernetes operator.

## Architecture

- **React Frontend** - Displays projects in a dynamic, grid-based layout
- **Go Backend API** - Serves project data to the frontend
- **Kubernetes Operator** - Watches for Project CRDs and syncs them to the backend

## Project Structure

```
0xHub/
├── frontend/          # React application
├── backend/           # Go API server
├── operator/          # Kubernetes operator
├── crd/              # Custom Resource Definitions
└── README.md
```

## Quick Start

### Prerequisites

- **Go 1.21+** - [Install Go](https://golang.org/doc/install)
- **Node.js 18+** - [Install Node.js](https://nodejs.org/)
- **npm** or **yarn** - Comes with Node.js

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080` and automatically seed 6 sample projects.

**API Endpoints:**
- `GET /api/health` - Health check
- `GET /api/projects` - Get all projects
- `GET /api/projects/:id` - Get a specific project
- `POST /api/projects` - Create a new project
- `PUT /api/projects/:id` - Update a project
- `DELETE /api/projects/:id` - Delete a project

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

The frontend will start on `http://localhost:5173` (or another port if 5173 is taken).

Open your browser and navigate to the URL shown in the terminal to see the project hub.

### Running Both Services

To run both the backend and frontend simultaneously, open two terminal windows:

**Terminal 1 (Backend):**
```bash
cd backend
go run cmd/server/main.go
```

**Terminal 2 (Frontend):**
```bash
cd frontend
npm run dev
```

## Development

### Backend Development

The backend uses:
- **Gin** - HTTP web framework
- **In-memory store** - Projects are stored in memory (data is lost on restart)
- **CORS** - Configured to allow requests from the frontend

The backend automatically seeds sample projects on startup. You can modify the `seedProjects` function in `backend/cmd/server/main.go` to add or change sample data.

### Frontend Development

The frontend uses:
- **React 18+** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework

The frontend is configured to proxy API requests to `http://localhost:8080` during development.

### Project Model

Projects have the following structure:
```typescript
{
  id: string          // Unique identifier
  name: string        // Project name
  description: string // Project description
  url: string         // Project URL
  icon?: string       // Optional icon URL
  category?: string   // Optional category
  status?: string     // Optional status (e.g., "active")
}
```

## Phase Status

### Phase 1: Core Hub ✅ **Completed**
- Project structure and directory setup
- Go backend API with in-memory store
- REST endpoints for projects
- CORS middleware configuration
- React frontend with TypeScript
- Vite setup with React 18+
- Responsive project grid/card layout
- API integration
- Tailwind CSS styling
- Manual project seeding (6 sample projects)
- Documentation

### Phase 2: CRD Definition ✅ **Completed**
- Project CRD YAML with full schema validation
- OpenAPI validation rules
- Example project resources
- Testing scripts for kind

### Phase 3: Kubernetes Operator ✅ **Completed**
- Operator project structure with controller-runtime
- Project controller with reconcile logic
- Backend sync client (HTTP client)
- CRD to backend project mapping
- Status updates on Project CRDs
- RBAC manifests
- Dockerfile and deployment manifests
- Kind testing script
- Documentation

**Next Steps:** Phase 4 - Integration & Polish

## License

MIT


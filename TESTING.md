# Testing Strategy for 0xHub

This document outlines the comprehensive testing strategy for the 0xHub project.

## Testing Philosophy

We follow a multi-layered testing approach:
1. **Unit Tests** - Test individual components in isolation
2. **Integration Tests** - Test component interactions
3. **End-to-End (E2E) Tests** - Test the full flow from CRD to frontend

## Test Structure

```
0xHub/
├── backend/
│   ├── internal/
│   │   ├── handlers/
│   │   │   └── projects_test.go
│   │   ├── store/
│   │   │   └── store_test.go
│   │   └── models/
│   │       └── project_test.go
│   └── cmd/
│       └── server/
│           └── main_test.go
├── operator/
│   ├── controllers/
│   │   └── project_controller_test.go
│   └── internal/
│       └── backend/
│           └── client_test.go
├── frontend/
│   ├── src/
│   │   ├── __tests__/
│   │   │   ├── api.test.ts
│   │   │   └── App.test.tsx
│   │   └── components/
│   │       └── __tests__/
│   │           └── ProjectCard.test.tsx
│   └── e2e/
│       └── basic.spec.ts
└── tests/
    └── e2e/
        └── full-integration.sh
```

## Backend Testing

### Unit Tests

**Store Tests** (`backend/internal/store/store_test.go`)
- Test CRUD operations
- Test concurrent access (race conditions)
- Test edge cases (empty store, non-existent IDs)

**Handler Tests** (`backend/internal/handlers/projects_test.go`)
- Test all HTTP endpoints
- Test request validation
- Test error handling
- Test CORS headers

**Model Tests** (`backend/internal/models/project_test.go`)
- Test JSON serialization/deserialization
- Test validation logic (if any)

### Running Backend Tests

```bash
cd backend
go test ./...
go test -v ./...  # Verbose output
go test -cover ./...  # With coverage
```

## Operator Testing

### Unit Tests

**Controller Tests** (`operator/controllers/project_controller_test.go`)
- Test reconcile logic for create/update/delete
- Test error handling and retry logic
- Test status updates
- Use fake client from controller-runtime

**Backend Client Tests** (`operator/internal/backend/client_test.go`)
- Test HTTP client methods
- Test error handling
- Use HTTP test server for mocking

### Running Operator Tests

```bash
cd operator
go test ./...
go test -v ./...
go test -cover ./...
```

## Frontend Testing

### Unit Tests

**Component Tests** (`frontend/src/components/__tests__/`)
- Test ProjectCard rendering
- Test search/filter functionality
- Test category tabs
- Use React Testing Library

**API Tests** (`frontend/src/__tests__/api.test.ts`)
- Test API client functions
- Mock fetch calls
- Test error handling

### Integration Tests

**App Tests** (`frontend/src/__tests__/App.test.tsx`)
- Test full component integration
- Test user interactions
- Test loading and error states

### Running Frontend Tests

```bash
cd frontend
npm test  # Run tests in watch mode
npm run test:coverage  # With coverage
```

## End-to-End Testing

### Full Integration Test

**Script** (`tests/e2e/full-integration.sh`)
- Sets up kind cluster
- Deploys operator
- Creates Project CRD
- Verifies sync to backend
- Verifies frontend display (optional, requires browser automation)

### Manual E2E Testing

1. Start backend: `cd backend && go run cmd/server/main.go`
2. Start frontend: `cd frontend && npm run dev`
3. Deploy operator to cluster (or run locally)
4. Create Project CRD: `kubectl apply -f crd/examples/example-project.yaml`
5. Verify:
   - Operator logs show sync
   - Backend API returns project: `curl http://localhost:8080/api/projects`
   - Frontend displays project

## Test Coverage Goals

- **Backend**: >80% coverage
- **Operator**: >70% coverage (controller logic)
- **Frontend**: >70% coverage (components and utilities)

## Continuous Integration

Tests should run:
- On every pull request
- Before merging to main
- On every commit (optional, can be done locally)

## Test Data

- Use fixtures for consistent test data
- Use factories for generating test projects
- Clean up test data after tests

## Mocking Strategy

- **Backend**: Use in-memory store (already implemented)
- **Operator**: Use fake Kubernetes client from controller-runtime
- **Frontend**: Mock fetch API calls
- **E2E**: Use real services in kind cluster

## Performance Testing

- Load testing for backend API (future)
- Frontend performance metrics (future)
- Operator reconciliation performance (future)


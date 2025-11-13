# Sarana AI Take-Home Test

Full-stack Notes Application with Authentication, CRUD Operations, Image Upload, and Logging.

## Tech Stack

### Backend (sarana-ai-take-home-test-be)
- **Language**: Golang
- **Framework**: Fiber (v2)
- **Database**: PostgreSQL
- **Authentication**: JWT with bcrypt password hashing
- **Features**:
  - User registration and login
  - CRUD operations for notes
  - Image upload for notes
  - Request/Response logging middleware

### Frontend (sarana-ai-take-home-test-fe)
- **Framework**: Next.js 16 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: React Context API
- **Features**:
  - User authentication (login/register)
  - Protected dashboard
  - Notes management (create, view, delete)
  - Image upload for notes

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Database**: PostgreSQL 16

## Project Structure

```
.
├── sarana-ai-take-home-test-be/        # Backend (Golang/Fiber)
│   ├── database/                        # Database connection & schema
│   ├── handlers/                        # HTTP handlers (auth, notes)
│   ├── middleware/                      # JWT auth & logging middleware
│   ├── models/                          # Data models
│   ├── utils/                           # JWT utilities
│   ├── uploads/                         # Uploaded images directory
│   ├── main.go                          # Application entry point
│   ├── go.mod                           # Go dependencies
│   └── .env.example                     # Environment variables template
│
├── sarana-ai-take-home-test-fe/        # Frontend (Next.js)
│   ├── src/
│   │   ├── app/                         # Next.js pages
│   │   │   ├── login/                   # Login page
│   │   │   ├── register/                # Register page
│   │   │   ├── dashboard/               # Protected dashboard
│   │   │   ├── layout.tsx               # Root layout with AuthProvider
│   │   │   └── page.tsx                 # Home page
│   │   ├── contexts/                    # React contexts
│   │   │   └── AuthContext.tsx          # Authentication context
│   │   └── lib/                         # Utilities
│   │       └── api.ts                   # API client
│   ├── package.json                     # NPM dependencies
│   ├── next.config.ts                   # Next.js configuration
│   └── .env.local.example               # Environment variables template
│
├── sarana-ai-take-home-test-misc/      # Docker configurations
│   ├── Dockerfile.be                    # Backend Dockerfile
│   └── Dockerfile.fe                    # Frontend Dockerfile
│
└── docker-compose.yml                   # Docker Compose orchestration
```

## API Endpoints

### Authentication
- `POST /register` - Register a new user
- `POST /login` - Login and get JWT token

### Notes (Protected - requires JWT)
- `POST /notes` - Create a new note
- `GET /notes` - Get all notes for the authenticated user
- `GET /notes/:id` - Get a specific note
- `DELETE /notes/:id` - Delete a note (with ownership check)
- `POST /notes/:id/image` - Upload an image for a note

### Health
- `GET /health` - Health check endpoint

## Features

### Backend Features
1. **Authentication**
   - Password hashing with bcrypt
   - JWT token generation and validation
   - Protected routes with JWT middleware

2. **Notes Management**
   - Create, read, and delete notes
   - User-specific notes (isolation)
   - Ownership verification for delete operations

3. **Image Upload**
   - Multipart form data handling
   - File type validation
   - Image storage in uploads directory
   - Database path reference

4. **Logging Middleware**
   - Logs all requests to database
   - Captures: datetime, method, endpoint, headers, request body, response body, status code
   - Masks sensitive headers (Authorization)

### Frontend Features
1. **Authentication**
   - Register and login pages
   - JWT token storage in localStorage
   - Protected routes with authentication check

2. **Dashboard**
   - Display all user notes
   - Create new notes with title and content
   - Delete notes with confirmation
   - Upload images for notes
   - Responsive grid layout

3. **State Management**
   - React Context for authentication state
   - Persistent login across page reloads

## Getting Started

### Prerequisites
- Docker and Docker Compose installed
- Git

### Running with Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   cd d:\workspace\sarana-ai\sarana-ai-take-home-test
   ```

2. **Start all services**
   ```bash
   docker-compose up --build
   ```

   This will start:
   - PostgreSQL database on port 5432
   - Backend API on port 8080
   - Frontend application on port 3000

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Health Check: http://localhost:8080/health

### Running Locally (Development)

#### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd sarana-ai-take-home-test-be
   ```

2. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Set up PostgreSQL database**
   ```bash
   # Using Docker
   docker run --name postgres-notes -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=notesapp -p 5432:5432 -d postgres:16-alpine
   ```

4. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. **Run the backend**
   ```bash
   go run main.go
   ```

#### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd sarana-ai-take-home-test-fe
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Configure environment variables**
   ```bash
   cp .env.local.example .env.local
   # Edit .env.local with your API URL
   ```

4. **Run the frontend**
   ```bash
   npm run dev
   ```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Notes Table
```sql
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    image_path VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Logs Table
```sql
CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    datetime TIMESTAMP NOT NULL,
    method VARCHAR(10) NOT NULL,
    endpoint VARCHAR(500) NOT NULL,
    headers TEXT,
    request_body TEXT,
    response_body TEXT,
    status_code INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Environment Variables

### Backend (.env)
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=notesapp
JWT_SECRET=your-super-secret-jwt-key-change-in-production
PORT=8080
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Testing the Application

1. **Register a new user**
   - Go to http://localhost:3000
   - Click "Register"
   - Enter email and password
   - You'll be redirected to the dashboard

2. **Create a note**
   - Fill in the title and content
   - Click "Create Note"

3. **Upload an image**
   - Click "Choose File" on any note
   - Select an image file
   - The image will be uploaded and displayed

4. **Delete a note**
   - Click "Delete" on any note
   - Confirm the deletion

## Docker Commands

```bash
# Start all services
docker-compose up

# Start in detached mode
docker-compose up -d

# Rebuild and start
docker-compose up --build

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# View logs
docker-compose logs -f

# View logs for specific service
docker-compose logs -f backend
```

## Security Features

- Password hashing with bcrypt
- JWT token authentication
- Protected API endpoints
- Authorization header masking in logs
- Note ownership verification
- SQL injection prevention with parameterized queries

## Future Enhancements

- Update note functionality
- Note search and filtering
- Tags and categories
- User profile management
- Refresh token implementation
- Rate limiting
- API documentation with Swagger
- Unit and integration tests
- CI/CD pipeline

## License

This project is for take-home test purposes.

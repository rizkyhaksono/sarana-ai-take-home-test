# Backend 

This backend is used for Sarana AI's take home test. It is built using Go and the Fiber web framework. The backend provide a RESTful API for user authentication and note management, with logs.
Also documentation is provided using Swagger, and Scalar.

## 📐 New Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         main.go                             │
│  • Database connection                                      │
│  • Loki initialization                                      │
│  • Fiber app setup                                          │
│  • Middleware configuration                                 │
│  • Swagger metadata (@title, @version, etc.)                │
│  • Documentation routes                                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ calls SetupRoutes()
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    routes/routes.go                         │
│  • Static file serving (/uploads)                           │
│  • Public routes (register, login, health)                  │
│  • Protected routes group (JWT middleware)                  │
│  • Notes endpoints                                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ routes to handlers
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   handlers/ (with Swagger)                  │
│                                                             │
│  ┌────────────────┐         ┌──────────────────┐           │
│  │  auth.go       │         │   notes.go       │           │
│  │ • Register()   │         │ • CreateNote()   │           │
│  │ • Login()      │         │ • GetNotes()     │           │
│  │                │         │ • GetNote()      │           │
│  │ @Summary       │         │ • DeleteNote()   │           │
│  │ @Description   │         │ • UploadImage()  │           │
│  │ @Tags          │         │                  │           │
│  │ @Router        │         │ @Security        │           │
│  └────────────────┘         │ @Summary         │           │
│                             │ @Router          │           │
│                             └──────────────────┘           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ delegates to
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      services/                              │
│  • auth_service.go                                          │
│  • note_service.go                                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ uses
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                     database/                               │
│  • Connection management                                    │
│  • Schema initialization                                    │
└─────────────────────────────────────────────────────────────┘
```

## 📚 Documentation Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    Developer writes code                    │
│                                                             │
│  handlers/auth.go                handlers/notes.go          │
│  // @Summary Register           // @Summary Create note    │
│  // @Router /register [post]    // @Router /notes [post]   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ $ swag init
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   Auto-generated docs/                      │
│                                                             │
│  ├── docs.go           (Go code)                            │
│  ├── swagger.json      (OpenAPI 2.0)                        │
│  └── swagger.yaml      (YAML format)                        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ served by
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Documentation UIs                        │
│                                                             │
│  ┌──────────────────┐         ┌──────────────────┐         │
│  │  Scalar UI       │         │   Swagger UI     │         │
│  │  /docs           │         │   /swagger/*     │         │
│  │                  │         │                  │         │
│  │  ✓ Modern        │         │  ✓ Traditional  │         │
│  │  ✓ Dark mode     │         │  ✓ Try it out   │         │
│  │  ✓ Interactive   │         │  ✓ Schemas      │         │
│  │  ✓ Code samples  │         │                  │         │
│  └──────────────────┘         └──────────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

## 🔄 Request Flow Example

```
┌──────────────┐
│   Client     │
└──────┬───────┘
       │
       │ GET /docs
       ▼
┌─────────────────────┐
│    main.go          │
│  /docs handler      │
│  Serves Scalar HTML │
└──────┬──────────────┘
       │
       │ Scalar loads /swagger/doc.json
       ▼
┌─────────────────────┐
│  fiber-swagger      │
│  Serves swagger.json│
└──────┬──────────────┘
       │
       │ JSON with all endpoints
       ▼
┌─────────────────────┐
│  Scalar UI          │
│  Renders beautiful  │
│  API documentation  │
└─────────────────────┘
```

## 🎯 Protected Route Flow

```
┌──────────────┐
│   Client     │
└──────┬───────┘
       │
       │ POST /notes
       │ Authorization: Bearer TOKEN
       ▼
┌─────────────────────────┐
│   routes/routes.go      │
│   JWT Middleware        │
│   Validates token       │
└──────┬──────────────────┘
       │
       │ Token valid, userID extracted
       ▼
┌─────────────────────────┐
│   handlers/notes.go     │
│   CreateNote()          │
│   Gets userID from ctx  │
└──────┬──────────────────┘
       │
       │ Delegates business logic
       ▼
┌─────────────────────────┐
│   services/note_service │
│   CreateNote()          │
│   Validates & saves     │
└──────┬──────────────────┘
       │
       │ Stores in database
       ▼
┌─────────────────────────┐
│   database/             │
│   PostgreSQL with UUID  │
└─────────────────────────┘
```

## 📦 File Dependencies

```
main.go
  ├─ routes/routes.go
  │    └─ handlers/*.go
  │         └─ services/*.go
  │              └─ database/
  │
  ├─ docs/docs.go (generated)
  │    ├─ swagger.json
  │    └─ swagger.yaml
  │
  └─ middleware/
       ├─ auth.go
       └─ logger.go
```

## 🏗️ Code Organization Benefits

### Before
```
main.go (120+ lines)
  ├─ All route definitions
  ├─ Handler logic mixed
  └─ No documentation
```

### After
```
main.go (85 lines)
  ├─ Configuration only
  └─ Swagger metadata

routes/routes.go (35 lines)
  └─ Clean route definitions

handlers/*.go (with Swagger)
  ├─ Swagger annotations
  └─ Handler logic

docs/ (auto-generated)
  └─ OpenAPI specification
```

## 🎨 Documentation UI Comparison

| Feature | Scalar | Swagger UI |
|---------|--------|------------|
| Modern Design | ✅ | ❌ |
| Dark Mode | ✅ | ❌ |
| Code Examples | ✅ Multiple languages | ⚠️ Limited |
| Mobile Friendly | ✅ | ⚠️ Partial |
| Search | ✅ Fast | ✅ Basic |
| Try It Out | ✅ | ✅ |
| Setup | Zero config | Zero config |

## 🔧 Development Workflow

```
1. Write handler code
   └─ Add Swagger annotations
      └─ Run: swag init
         └─ Test in Scalar UI
            └─ Commit changes
```

## 📝 Maintenance

### Adding New Endpoint

```go
// 1. Add handler with Swagger annotations
// @Summary Your summary
// @Router /path [method]
func YourHandler(c *fiber.Ctx) error {
    // Implementation
}

// 2. Add route in routes/routes.go
app.Post("/path", handlers.YourHandler)

// 3. Regenerate docs
$ swag init

// 4. Test at http://localhost:8080/docs
```

### Updating Existing Endpoint

```go
// 1. Update handler and annotations
// 2. Run: swag init
// 3. Refresh browser
```

## 🎓 Best Practices Applied

1. ✅ **Separation of Concerns**: Routes separated from main
2. ✅ **Documentation as Code**: Swagger annotations in handlers
3. ✅ **Auto-generation**: No manual JSON editing
4. ✅ **Multiple UIs**: Scalar (modern) + Swagger (traditional)
5. ✅ **Type Safety**: Go structs define API contracts
6. ✅ **Maintainability**: Easy to find and update routes
7. ✅ **Developer Experience**: Interactive documentation

---

**Architecture Version**: 2.0  
**Last Updated**: November 13, 2025

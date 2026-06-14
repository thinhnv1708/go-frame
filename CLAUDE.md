# Project Development Guidelines (AGENT.md)

This document defines the architectural patterns, design principles, and coding standards for the project. Any AI Assistant contributing to this codebase must read, understand, and strictly adhere to these guidelines to ensure consistency and clean structure.

---

## 1. Architectural Overview (Clean & Layered Architecture)

The project is structured around a **Clean/Layered Architecture** combined with **Manual Dependency Injection (DI)**. The layers are strictly decoupled:

```
[Client] 
   │
   ▼ (HTTP Request)
[Handler] ───► Bind & Validate (dto/request)
   │
   ▼ (Entity / DTO)
[Service] ───► Core Business Logic
   │
   ▼ (Entity)
[Repository] ───► Data Persistence
   │
   ▼ (Database Model)
[Database Layer (ORM / Driver)]
```

### Core Rules for Data Flow:
* **Handler Layer**: Parses incoming HTTP requests, binds payloads to Request DTOs, performs validation, calls the appropriate Service, and registers returned errors in the framework's HTTP context.
* **Service Layer**: Orchestrates business logic. It works exclusively with Domain Entities and DTOs. It does not access Database Models or database connection pools directly.
* **Repository Layer**: Encapsulates persistence logic. It queries the database engine and maps the resulting persistent models/documents to pure Domain Entities before returning them to the Service.
* **DTO (Data Transfer Object) Layer**: Separated into `request` and `response` packages, serving as the communication schema between Client and Handler.

---

## 2. Directory Structure and Responsibilities

* **`cmd/`**: Entry point of the application. Handles container building, server initialization, and graceful shutdown setups.
* **`internal/`**: The private application code.
  * **`app/`**: Application bootstrap and lifecycle management.
  * **`config/`**: Configuration structures and loaders.
  * **`database/`**: Database connections and persistence models/schemas.
  * **`di/`**: Centralized dependency injection configurations.
  * **`dto/`**: Data Transfer Objects (request/response schemas).
  * **`entity/`**: Domain entities representing business models.
  * **`exception/`**: Predefined business exceptions and custom error structures.
  * **`handler/`**: Transport handlers (HTTP/RPC) responsible for request parsing and output formatting.
  * **`logger/`**: Logging infrastructure.
  * **`middleware/`**: Shared HTTP middlewares (e.g., centralized error handling).
  * **`provider/`**: Protocol implementations (such as HTTP servers and router configurations).
  * **`repository/`**: Persistence layer interfaces and concrete implementations.
  * **`security/`**: Hashing, cryptography, and token generation utilities.
  * **`service/`**: Business logic layer interfaces and concrete implementations.
  * **`validation/`**: Framework-level validation error parsing.

---

## 3. Naming & Coding Conventions

### 3.1. Database Models vs. Domain Entities
* **Database Models**: Placed in the database layer. Represent the physical schema. They contain database tags, indices, and timestamps.
* **Domain Entities**: Placed in the entity layer. Represent the business representation. They contain only raw Go types with no persistence or framework tags.
* **Mapping**: Repository implementations must implement helper functions to map between persistence models and domain entities.

### 3.2. Interfaces & Concrete Implementations
* **Interface Definition**: All interfaces for a layer must be defined together in a centralized interface file (e.g., `interface.go`) of that package.
* **Concrete Structs**: Implementation structs must be named with the suffix `Imp` (e.g., `ServiceImp`, `RepositoryImp`).
* **Constructors**: Instantiation functions must follow the `New<Name>` convention and return the interface type instead of the concrete struct.

### 3.3. Error Handling
* Business errors must be predefined as global variables of a custom error type containing an internal application code, message, and corresponding HTTP status code.
* Services return these custom error instances upon failure.
* Handlers pass errors to the framework's HTTP context instead of formatting JSON responses. A centralized middleware intercepts errors from the context, identifies custom exceptions, and outputs a standardized error JSON structure.

### 3.4. Transaction Management
* Transactions must be propagated via context boundaries.
* Repositories should resolve the active database/session transaction from the context to participate in transactions.
* Services orchestrate transactional boundaries using a transaction manager abstraction.

### 3.5. Manual Dependency Injection (DI)
* Do not use automated DI libraries or code generators.
* All components must be wired manually inside the DI package.
* When adding a new Repository, Service, or Handler, manually instantiate and register them in the dependency provider builder.

---

## 4. Step-by-Step Guide to Implementing a New Feature

When adding a new capability, follow this sequence:

1. **Define the Domain Entity**: Create the core business representation in the entity layer.
2. **Define the Database Model**: Create the persistence schema in the database layer (if database persistence is required).
3. **Define and Implement the Repository**: 
   * Declare the new repository interface in the centralized repository interface file.
   * Create the concrete repository implementation.
   * Write mappers to convert between the Database Model and Domain Entity.
4. **Create Request & Response DTOs**: Define schemas for API payloads in the DTO layer.
5. **Define and Implement the Service**:
   * Declare the new service interface in the centralized service interface file.
   * Create the concrete service implementation carrying the business logic.
6. **Create the Handler**: Implement transport handling, request binding, validation, and error propagation.
7. **Register Routes**: Implement sub-routers and register them in the main routing configuration.
8. **Configure Dependency Injection**: Instantiate and wire up the new repository, service, and handler in the manual DI provider builder.

---

## 5. Development Checklist
* [ ] Are Domain Entities strictly separated from Database/Persistence Models?
* [ ] Are conversion mappers implemented in the Repository layer?
* [ ] Does the constructor function return the Interface type instead of a concrete struct?
* [ ] Are business exceptions declared as global variables of the custom error type?
* [ ] Are handler-level errors propagated via the framework context to delegate formatting to the middleware?
* [ ] Are the new components registered manually in the DI container?
* [ ] Are unit tests written for the new Service and Repository layers?

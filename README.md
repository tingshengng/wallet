# Starting the service locally
1. Start postgres
```bash
make pg-up
```
2. Start the service
```bash
make up
```

# Techstack used
- Language: Go
- Web framework: Gin
- Database: Postgres
- ORM: GORM
- Cache: InMemoryCache
- Unit test mock generation: GoMock

# Design details

### Code Structure
- cmd/api/main.go: Entry point where the service is initialized.
- internal/models: Defines the data models used by the service.
- internal/handlers: Contains the API handlers.
- internal/middleware: Handles authentication logic before requests reach the handlers.
- internal/services: The business logic layer
- internal/repositories: Provides the database access layer.
- internal/cache: Contains the cache implementation

### Login and Authentication

Since this is a wallet service, all APIs must be authenticated to a user. However, as authentication is not the main focus of this project, we simulate it using a mock login API. The mock API accepts an email address and returns a short-lived authentication token valid for 4 hours.

This approach mimics a magic link login flow, where a user initiates a login via email, and the server email the user with a login link containing the token. In real-world applications, authentication would typically include more complex steps. Besides, there would be other steps like identity verification during account creation, which is out of scope for this project.

### Deposit and Withdrawal
For simplicity, deposit and withdrawal are implemented as straightforward API calls with immediate responses. In real-world systems, these operations are better suited to an event-driven architecture, where the API sends a message to a queue and a queue worker processes it asynchronously. This is because such operations usually involve external payment service integrations, which can take time to complete.

### Caching
Caching is primarily used for the transaction history API. We cache the 10 most recent transactions per user, as these are typically the most frequently accessed. Limiting the cache to 10 transactions helps manage memory usage efficiently.

The cache is implemented as a lightweight in-memory store. When a user retrieves their transaction history, we populate the cache. When a deposit or withdrawal occurs, we evict the cache for that user. The main limitation of this approach is that in-memory caches don’t scale well across multiple service instances, especially with our action-based eviction strategy. However, since the cache interface is already defined, it can be easily replaced with a distributed solution like Redis in the future.



# Unit test
Unit tests focus on the core logic in the service layer. We use golang/mock to generate mock implementations of the repository layer for testing.
```bash
go test ./... -v
```

# API Endpoints for testing

**Login**
```bash
curl --location 'http://localhost:8888/api/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "satoshi@gmail.com",
    "name": "Satoshi Nakamoto"
}'
```

**Deposit**
```bash
curl --location 'http://localhost:8888/api/deposit' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {token-from-login-resp}' \
--data '{
    "amount": 800
}'
```

**Withdraw**
```bash
curl --location 'http://localhost:8888/api/withdraw' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {token-from-login-resp}' \
--data '{
    "amount": 200
}'
```

**Transfer**
```bash
curl --location 'http://localhost:8888/api/transfer' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {token-from-login-resp}' \
--data '{
    "to_user_id": "9d49cb09-609e-4200-8684-53c6983c9a40",
    "amount": 10
}'
```

**Get Balance**
```bash
curl --location 'http://localhost:8888/api/balance' \
--header 'Authorization: Bearer {token-from-login-resp}'
```

**Get Transaction History**
```bash
curl --location 'http://localhost:8888/api/transactions?type=deposit' \
--header 'Authorization: Bearer {token-from-login-resp}'
```




# Starting the service locally
1. Add an .env file, sample .env for local development
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=wallet_db
DB_SSLMODE=disable
```
2. Start postgres
```bash
make pg-up
```
3. Start the wallet service
```bash
make up
```

# Techstack used
- Language: Go
- Web framework: Gin
- Database: Postgres
- ORM: GORM
- Cache: In-memory cache
- Unit test mock generation: GoMock

# Design details

### Code Structure
- `cmd/api/main.go`: Entry point where the service is initialized, and endpoints are defined.
- `internal/models`: Defines the data models used by the service.
- `internal/handlers`: Contains the API handlers.
- `internal/middleware`: Handles authentication logic before requests reach the handlers.
- `internal/services`: The business logic layer, this is where the main logic of the wallet service is implemented.
- `internal/repositories`: Provides the database access layer.
- `internal/cache`: Contains the cache implementation

### Simplified Login and Authentication
Since this is a wallet service, all APIs must be authenticated to a user. However, as authentication is not the main focus of this project, we intentionally simplify it with a mock login API. The mock API accepts an email address and returns a short-lived authentication token valid for 4 hours.

This approach mimics a magic link login flow, where a user initiates a login via email, and the server email the user with a login link containing the token. In real-world applications, authentication would typically include more complex steps. Besides, there would be things like identity verification etc. during account creation, which is out of scope for this project.

### Simplified Deposit and Withdrawal
To keep things simple, we have skipped integration with external payment services. Deposit and withdrawal are also implemented as straightforward API calls with immediate responses. In real-world systems, these operations are better suited to an event-driven architecture, where the API sends a message to a queue and a queue worker processes it asynchronously. This approach is better suited for handling external payment integrations, which may involve retries, failures, or delays.

### Caching
Caching is primarily used for the transaction history API. We store the 10 most recent transactions per user, as these are typically the most frequently accessed. At the same time, caching too many transactions can lead to high memory usage. Limiting it to 10 strikes a balance between performance and resource efficiency.

The cache is implemented as a lightweight in-memory store. When a user retrieves their transaction history, we populate the cache. When a deposit or withdrawal occurs, we evict the cache for that user. The main limitation of this approach is that in-memory caches don’t scale well across multiple service instances, especially with our action-based eviction strategy. However, since the cache interface is already defined, it can be easily replaced with a distributed solution like Redis in the future if scalability becomes a concern.

### UUID as primary identifier
We use UUIDs as primary identifiers for all models instead of auto-incrementing IDs. This is a security-conscious choice that makes it significantly harder to guess or enumerate resources based on predictable ID patterns.

# Unit test
Unit tests focus on the core functions in the service layer: Deposit, Withdrawal and Transfer. 
We use golang/mock to generate mock implementations of the repository layer for testing.
To run unit tests:
```bash
go test ./... -v
```

# API Endpoints for testing

**Login**
```bash
curl --location '{baseUrl}/api/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "satoshi@gmail.com",
    "name": "Satoshi"
}'
```

**Deposit**
```bash
curl --location '{baseUrl}/api/deposit' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {token-from-login-response}' \
--data '{
    "amount": 800
}'
```

**Withdraw**
```bash
curl --location '{baseUrl}/api/withdraw' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {token-from-login-response}' \
--data '{
    "amount": 200
}'
```

**Transfer**
```bash
curl --location '{baseUrl}/api/transfer' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {token-from-login-response}' \
--data '{
    "to_user_id": "9d49cb09-609e-4200-8684-53c6983c9a40",
    "amount": 10
}'
```

**Get Balance**
```bash
curl --location '{baseUrl}/api/balance' \
--header 'Authorization: Bearer {token-from-login-response}'
```

**Get Transaction History**
```bash
curl --location '{baseUrl}/api/transactions?type=deposit' \
--header 'Authorization: Bearer {token-from-login-response}'
```

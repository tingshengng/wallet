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

# Design details
**Code structure**

- cmd/api/main.go: The entry point where we initialize our service.
- internal/models: The data models of the service
- internal/handlers: The API handlers
- internal/middleware: Handle the authentication logics before each requests hit the handlers
- internal/repositories: The database access layer
- internal/cache: The cache layer

**Login and authentication**

We implemented a mock login API as it is not the focus of this project.
**Deposit and Withdraw**


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




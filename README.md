# Starting the service locally
1. Start postgres
```bash
make pg-up
```
2. Run migrations
```bash
make migration-up
```
3. Start the service
```bash
go run cmd/api/main.go
```

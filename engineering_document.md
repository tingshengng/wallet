# User Story
- User can deposit money into his/her wallet
- User can withdraw money from his/her wallet
- User can send money to another user
- User can check his/her wallet balance
- User can view his/her transaction history

# RESTful API
- Deposit to specify user wallet
- Withdraw from specify user wallet
- Transfer from one user to another user
- Get specify user balance
- Get specify user transaction history

# Tech stack
- Language: Go
- Web framework: Fiber
- Database: Postgres
- ORM: GORM
- 


# Prompt
Please create a simple Golang boilerplate with the tech stack given, where the web framework is Fiber. The database should be using Postgres, and the ORM should be GORM.

For postgres, we will need a simple docker-compose file to run the postgres container, with a docker volume for data persistence. Please create a Makefile to run the postgres container with make pg-up.
The postgres credentials should be stored in a .env file that the Go application can read.

For migration, we should not use auto migration, but instead, we should use raw SQL for migration. The SQL files should be stored in a migrations folder. It can be triggered to run migration manually when we want, which you can create a Makefile target for it. 
We also need to revert the migration if needed. So we will need make migration-up and make migration-down.
I would like to utilize GORM to generate the SQL files, but it will just generate the SQL which I can review and manually run it, the migration command would be make migration-generate.
We also need a make migration-create migrationName for creating a new migration file.

Above are the tech stack requirements.
Now, for the applicaiton logic. 
Please create a new table called as Test, with a few columns: id, uuid, name, created_at, updated_at, deleted_at. 
All timestamps should be timestamp with timezone.
Do not run the migration, just create a raw SQL file in the migrations folder. 
Next, create two simple API to POST to the Test table and GET from the Test Table.

With this, I should be able to run make pg-up to start the postgres container, make migration-up to run the migration and then I can call my API to insert and get data from the Test table.
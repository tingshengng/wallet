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


I need to build a simple Wallet App.
First, we need to create the models. You can recommend me if there are additional columns that you think is needed.
1. User: For storing user. 
- columns: id(uuid), name, email, created_at, updated_at, deleted_at. 
2. UserToken: For storing short lived user token for authentication
- columns: id, user_id(uuid), token, created_at, updated_at, expires_at. 
3. Wallet: For storing user wallet details. For simplicity, each user will only have one wallet.
- columns: id(uuid), user_id(uuid), balance, created_at, updated_at, deleted_at. 
4. Transaction: For storing user transaction history. (type: deposit, withdraw, transfer; status: pending, success, failed)
- columns: id(uuid), from_user_id(uuid), to_user_id(uuid), amount, type, status, created_at, updated_at, deleted_at. 
Please create the models and update in the migrations function.

Next, we will create the APIs. Below are the APIs and features needed:
1. User Login Mock API: Just a mock API for login. We will call the API with an email in the request body. It will create a user if not exist, along with the wallet. Next, this API will return a short lived user token that last for 24 hours. If the Login API is called before expiry, the token will expiry time will be refreshed 24 hours. The other APIs below will need to use this token for authentication and throw unauthorized error if the token is expired.
2. Deposit money to user wallet: Just take whatever amount of money from the request body and add it to the user wallet.
3. Withdraw money from user wallet: Need to have a minimum balance check.
4. Transfer money to another user: Need to have a minimum balance check.
5. Get user balance: User can only get their own balance.
6. Get user transaction history: User can only get their own transaction history. Besides, this API need to support pagination and filter by type, status.

Other API requirements:
- Since all APIs need to use the user token for authentication, we need to create a middleware to check the user token. The middleware should be used for all APIs except the login API. Then, the user object should be passed in to the API handler.


Repository layer:
We will create a repository layer for database operations. Please create an interface for all the DB operations. Please separate the different models' DB operations into different files. The repository layer should be used by the API handler.
Please only create the methods that are needed for the APIs above.

Please make the following changes:
1. Move Login handler to a separate file called user.go
2. For all the request body struct in every handler, please move the struct definition below the Handler struct and give a relevant name.
3. For the "id" columns, please auto generate a uuid for it.

The Login handler have multiple issues.
It needs to first find by email
- if not exist, create the user and wallet. Next, create a token and return it.
- if exist, search for the user token by userId and update the expiry time to another 24 hours.
 Take note that your current implementaion uses h.UserTokenRepo.FindByToken which is wrong. Please also update the repository layers accordingly


For the models, I want to add indexes. Please help me update the models accordingly.
For the User model, I want to add index for email.
For the UserToken model, I want to add index for user_id and token.
For the Wallet model, I want to add index for user_id.
For the Transaction model, I want to add index for user_id and to_user_id.

I want to implement an in-memory cache for get transaction history API.
The cache key should be the userID, and the value should be the 10 most recent transaction history.
Therefore, we only need to set the cache when the request is page=1 and page_size=10.
For cache eviction, we have to monitor the deposit, withdrawal and transfer APIs, when these APIs are called, we need to invalidate the cache for the user. Take note that for transfer API, we need to invalidate the cache for both the sender and the recipient.

Please use this library for the in memory cache https://github.com/patrickmn/go-cache. You will need to declare an interface for the cache so that in the future we can move to Redis easily.
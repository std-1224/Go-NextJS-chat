create migration

migrate create -ext sql -dir db/migrations add_users_table



# FORMAL

1. postgres image from docker hub
    • create the go-chat database
2. database setup
    • make a connection to the db
    • add a db migration file to create the 'users' table
3. /signup endpoint to create a new user
    • repository + service + handler (dependencies)



# Architecture
- DB
- Repository
- Service
- Handler
- Routes

# api
`go run cmd/api`
The interface for the e-shop system, does all the http request handling.

# admin
`go run cmd/admin`
Admin cli tool, used to manage the system, add products, etc.

# user
`go run cmd/user`
User cli tool, used to interact with the system as a user, browse products, make orders, etc.

# scheduler
`go run cmd/scheduler`
The scheduler polls events table for upcoming, not preheated events, and preheats them 
by setting up values in redis.

# worker
`go run cmd/worker`
The worker listens to redis streams for seckill events, and processes them by creating 
actual orders in the database.

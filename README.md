# Shorty 

## Setup

### Dependencies

`dep` is used to manage third party libraries:

`$ dep ensure`


### Docker Compose
Used for Mongo, Postgres, and Redis.

Install docker on your machine and run the following command from app root.

`docker-compose up`

If for some reason the `init.sql`  fails to run then the best solution is to retry via:

`docker-compose stop`

`docker ps -a` get the container ID for `shorty_postgres_1`

`docker rm -v {containerID}`

### Run
`go run cmd/web/main.go`

## Test
`go test ../.`



## API

### `POST /link`
Creates a short url from a long url 

### `GET /link/{id}`
Retrieves short link data

### `GET /{short-link}`
Redirects to full URL
``

# Explanation of Work

#### design philosophy
My intention was to go beyond the simplest version of this application which may consist of http handlers doing "all of the work". <br /> 
The application structure uses *some* patterns borrowed from Domain Driven Design in order to decouple various responsibilities. <br />  


Document storage for links

Redis cache probably

Postgres for stats, 2 tables: raw and snapshot

^daily job to process stats

^ short link at seven chars is ~5 billion possibilities, 6 ==  ~2 billion
# Disclaimer

* No work has been done to properly setup MongoDB, Redis, Postgres. This includes basic security.
* Per assignment, no auth or user management is provided.
* the document store (in this case mongo) would have to be rethought for production, perhaps sharding
* the `.env` file is committed for convenience - this would never happen in a real world scenario
* Main should ensure that DB connections are valid; at this point that check is not in place

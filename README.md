# Smurl

## Setup

### Dependencies

`dep` is used to manage third party libraries:

`$ dep ensure`


### Docker Compose
Used for Mongo and Redis.

Install docker on your machine and run the following command from app root.

`docker-compose up`

### Run
`go run cmd/smurl/web/main.go`

## Test
`go test ../.`



## API

### `POST /short`
Creates a short url from a long url 

### `GET /short/{id}`
Retrieves short link data

### `GET /{short-link}`
Redirects to full URL

``

# Explanation of Work
decoupled services: link, keygen, stats

document storage for links

redis cache probably

postgre for stats, 2 tables: raw and snapshot

^daily job to process stats

^ short link at seven chars is ~5 billion possibilities

# Disclaimer

* No work has been done to properly setup MongoDB or Redis. This includes basic security.
* Per assignment, no auth or user management is provided.
* the document store (in this case mongo) would have to be rethought for production, perhaps sharding

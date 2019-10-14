# Shorty 

## Setup

### Dependencies

`dep` is used to manage third party libraries:

`$ dep ensure`


### Docker Compose
Install docker on your machine and run the following command from app root.

`docker-compose up`

If for some reason the `init.sql`  fails to run then the best solution is to retry via:

`docker-compose stop`

`docker ps -a` get the container ID for `shorty_postgres_1`

`docker rm -v {containerID}`

### Run
`go run cmd/web/main.go`

## Test
`go test -v ./...`

### curltest.sh

A small script is included that will sequence 100000 parallel post requests to the create link endpoint.
This script depends on `parallel`.
IF running MacOS ensure this dependency is available, if not:

`brew install parallel`

`parallel --citation`

run the script with

`./curltest.sh`


## API

### `POST /links`
Creates a short url from a long url 
```$xslt
{
  "destination": "valid url"
}
```
Destination must be the complete URL. For example `http://wwww.somewebsite.com` is valid whereas `somewebsite.com` is not.

### `GET /{slug}`
Redirects to destination URL

### `GET /links/{slug}/stats`
Retrieves short link data

``

# Explanation of Work

My intention was to go beyond the simplest version of this application which may consist of http handlers doing "all of the work". <br /> 
The application structure uses *some* patterns borrowed from Domain Driven Design in order to decouple various responsibilities. <br />  



The application logic is provided by two services and one small package.
The `link` service is responsible for the management of shortlinks. This service depends on the `keygen` package; a simple random key generator.
The key generation for the `link` service followed a simple idea that a string of six characters comprised of 63 possible characters could produce a total of 62,523,502,209 possible shortlinks (63^6).
The `stats` service is responsible for tracking the utilization of a given short link.

The shortlinks are persisted to MonogoDB. Document storage seems ideal  as the design anticipates millions (perhaps more)  shortlinks which themselves with little need for relations.
The shortlink statistics are persisted to Postgres. My assumption is that an analytics type service would probably benefit from an RDMS.





#### short comings 
The processing of shortlink usage statistics would be better serviced by a periodic job to process the raw view counts and create a kind of snapshot.
In keeping with the "this assignment should be accomplished in a day, perhaps 4 hours" I did not feel it was appropriate to spend the additional time to implement this.

Similarly, I suspect that the call to save a `view` should be done processed via some sort of queue in order to ease load on Postgres.

Testing is provided only in the simplest form. A basic integration test is valuable for each handler. Given sufficient time, the application should fulfull the testing pyramid as described by Mike Cohn.
Test setup as found in `TestMain` is a mess and not acceptable for production code.
Lastly, there is no environment distinction between test and anything else.

No caching has been implemented. It makes sense to provide a cash of shortlinks with an LRU strategy.

The `stats` package, while decoupled, is highly specific to `link`. Ideally, `stats` would work in a more general.
The snapshot concept is implemented as a simple query. This is problematic in two ways. First, the `stats_raw` table could grow to a problematic size. Second, the query's performance will degrade over time.

Logging is practically absent. Each request *should* be given a UID which can be used a logger.

Last, as I am at the suggested time limit I am forgoing the typical refactoring that often happens when one reads one's own code after a good sleep.




#### Disclaimer

* No work has been done to properly setup MongoDB or Postgres. This includes basic security.
* Per assignment, no auth or user management is provided.
* the document store (in this case mongo) would have to be rethought for production, perhaps partitioning. 
* the `.env` file is committed for convenience - this would never happen in a real world scenario
* Main should ensure that DB connections are valid; at this point that check is not in place
* curltest.sh could be expanded to target more endpoints

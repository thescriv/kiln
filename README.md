Kiln mid exercice

## Presentation

This Repo tried to handle the exercice case of Kiln ->.
Expected time by Kiln was around 3hours.

## Setup

Download and install docker :

    ```bash
    brew install docker
    ```

Then run the docker-compose file :

```bash
    docker compose up -d
```

The docker-compose.yml provide two databases : mysql and mysql_test.
-   mysql is used by the principal service to run the core project.
-   mysql_test is used by the test to insert and read from a db.

Download and install magefile:

```bash
    git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go && export PATH=$(go env GOPATH)/bin:$PATH && cd ..
```

copy content from `.env.example` to `.env` and adjust as needed.
Value provided in `.env.example` match connection specified in `docker-compose.yml`

## Usage

`go run cmd/main.go`

The project expose a single endpoint : `xtz/delegations` which return the last delegations found in DB.
You can provide a bunch of query params :

-   `year`
    -   Fetch delegations found in DB for a given year.
-   `page`
    -   Fetch delegations for a given page.
-   `limit`
    -   Limit the delegations fetching to a limit, as a lot of delegations can be found, by default the value is set to `100`.

### Worker Specification

-   If no delegations exist in db, the worker will start to fetch all delegations from previous day (`time.Now().AddDate(0, 0, -1)`)
-   If delegations exist in db, the worker will pull the more recent one based on the field `timestamp` and will fetch all new delegations based on the `timestamp` found.

## MageFile

`mage tezos:fetchDelegationsFromYear`

This repo provide a magefile `tezos:fetchDelegationsFromYear` and wait a for a parameter `year`.
The magefile `fetchDelegationsFromYear` will poll all delegations from the year provided and it will store it in the mysqlDB provided by docker.

## Test case

you can run test

``` bash
    go test ./...
```
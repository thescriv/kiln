Kiln mid exercice

## Presentation

This Repo tried to handle the exercice case of Kiln -> https://kilnfi.notion.site/Mid-Backend-Exercice-ceb13c8dd6114dcc98d4ba5053370bd8.
Expected time by Kiln was around 3hours.
it took me around 5-6hours to handle all things i wanted to handle.
As the code can be a bit "fat", i tried to not overload reviewer and store logic in small repository as their is a lot of use case to handle.

I loose a lot of time with a bunch of error with database configuration and error debugging.

Overall the test was interesting, challenging and i learned new things, thank you.

## Setup

Download and install docker :

    ```bash
    brew install docker
    ```

Then run the docker-compose file :

```bash
    docker compose up -d
```

Download and install magefile:

```bash
    git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go && export PATH=$(go env GOPATH)/bin:$PATH && cd ..
```

copy content from `.env.example` to `.env` and adjust as needed.
Value provided in `.env.example` match connection specified in `docker-compose.yml`

## Usage

`go run cmd/main.go`

The repo provide two databases : mysql and mysql_test.

-   mysql is used by the principal service to run the core project.
-   mysql_test is used by the test to insert and read from a db.
    -   I choosed to use a testDB instead of a mock db to win time and not declare mocking function.

The project expose a single endpoint : `xtz/delegations` which return the last delegations found in DB.
You can provide a bunch of query params :

-   `year`
    -   Fetch delegations found in DB for a given year.
-   `page`
    -   Fetch delegations for a given page.
-   `limit`
    -   Limit the delegations fetching to a limit, as a lot of delegations can be found, by default the value is set to `100`.

### Worker Specification

-   If no delegations exist in db, the worker will start to fetch all delegations from `time.Now().AddDate(0, 0, -1)`
-   If delegations exist in db, the worker will pull the more recent one based on the field `timestamp` and will fetch all new delegations based on the `timestamp` found.

## MageFile

`mage tezos:fetchDelegationsFromYear`

This repo provide a magefile `tezos:fetchDelegationsFromYear` and wait a for a parameter `year`.
The magefile `fetchDelegationsFromYear` will poll all delegations from the year provided and it will store it in the mysqlDB provided by docker.

## Test case

A bunch of test can be found in this repo :

-   pkg/utilworker
-   pkg/tezos
-   pkg/miscellaneous
-   cmd/worker
-   cmd/xtz

With a time constraint i aimed to only test the core of the project.

## What to do next

-   Add more test
-   Improve sql request to accelerate queries
-   Improve error handling
    -   Right now the worker and magefile only panic or return if something happend, with more error handling the code can be more resilient.
-   Add information stored in delegations if needed.
-   `cmd/worker/delegations` and `magefiles/fetchDelegationsFromYear` use almost the same code, we can think of a better solution to handle this.
-   Add quality of comment, i tried to keep everything documented but sometimes comment can be a bit flaky/poor.
-   fix inconsistency in the codebase
    -   flaky importation of models, sometimes I use the system of . import and sometimes i call `models`

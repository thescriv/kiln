Kiln mid exercice

## Presentation

This Repo tried to handle the exercice case of Kiln -> https://kilnfi.notion.site/Mid-Backend-Exercice-ceb13c8dd6114dcc98d4ba5053370bd8.
Expected time by Kiln was around 3hours.
it took me around 5-6hours to handle all things i wanted to handle.
As the code can be a bit "fat", i tried to not overload reviewer and store logic in small repository as their is a lot of use case to handle.

I loose a lot of time with a bunch of error with database configuration and error debugging.

Overall the test was interesting, challenging and i could've learn new things, thank you.

## Setup

Download and install docker & docker compose -> `brew install docker`
Download and install magefile -> `git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go && export PATH=$(go env GOPATH)/bin:$PATH && cd ..`

copy content from `.env.example` to `.env`

## Usage

When launching the principal process a worker is launched in a goroutine to fetch recent delegations from tezos blockchain.

-   If no delegations exist in db, the worker will start to fetch all delegations from `time.Now().AddDate(0, 0, -1)`
-   If delegations exist in db, the worker will pull the more recent one based on the field `timestamp` and will fetch all new delegations based on the `timestamp` found.

The project expose a single endpoint : `xtz/delegations` which return the last delegations found in DB.
You can provide a bunch of query params :

-   `year`
    -   Fetch delegations found in DB for a given year.
-   `page`
    -   Fetch delegations for a given page.
-   `limit`
    -   Limit the delegations fetching to a limit, as a lot of delegations can be found, by default the value is set to `100`.

## MageFile

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
-   Add more comment, right now their is not a lot of comment as i tried not to over-exxagerate the time spent on the exercice.

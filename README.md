## full-text-search-api

A RESTful full text search service using Postgres.

[![Build Status](https://github.com/anothrNick/full-text-search-api/workflows/Bump%20version/badge.svg)](https://github.com/anothrNick/full-text-search-api/workflows/Bump%20version/badge.svg)
[![Stable Version](https://img.shields.io/github/v/tag/anothrNick/full-text-search-api)](https://img.shields.io/github/v/tag/anothrNick/full-text-search-api)

Create records keyed by `project` that are indexed and searchable using [Postgres full text search](https://www.postgresql.org/docs/12/textsearch.html).

### Usage

##### Fields

|Field|Description|
|-----|-----------|
|`data`|Records are searched on the `data` field. Can be any data type (`number`, `string`, `array`, `object`)|
|`meta`|Provided for any relevant information needed to associate with a discovered record. Can be any data type (`number`, `string`, `array`, `object`)|
|`project`|Project `string`, used to organize searchable data into buckets. This is also the partition key for Postgres.|

##### cURL Commands

```sh
# create some data for a project called 'test'
curl -X POST -d '{"data": {"name": "Bart", "some_other_count": 30}, "meta":{"source_id": 10222}}' localhost:5001/test
curl -X POST -d '{"data": {"name": "Jimbo", "age": 30}, "meta":{"source_id": 10222}}' localhost:5001/test
curl -X POST -d '{"data": {"one": "another query for debug"}, "meta":{}}' localhost:5001/test
curl -X POST -d '{"data": {"one": "test query"}, "meta":{}}' localhost:5001/test

# full text search with `_search` query parameter
# multiple results containing `query`
$ curl -s "localhost:5001/test?_search=query" | jq "."
{
  "items": [
    {
      "id": "83947653-454d-4e29-85a4-0c4614fb3ed3",
      "project": "test",
      "data": {
        "one": "test query"
      },
      "meta": {}
    },
    {
      "id": "0b8961ea-7b76-44eb-ab49-3f2a23bf7f3f",
      "project": "test",
      "data": {
        "one": "another query for debug"
      },
      "meta": {}
    },
    {
      "id": "32c3675f-b922-4733-9041-9c7939b3894d",
      "project": "test",
      "data": {
        "one": "another query for debug"
      },
      "meta": {
        "source_id": 10001
      }
    }
  ]
}

# multiple results containing 30
$ curl -s "localhost:5001/test?_search=30" | jq "."
{
  "items": [
    {
      "id": "45de46fd-af9e-4f54-84ed-241a9f803b2d",
      "project": "test",
      "data": {
        "age": 30,
        "name": "Jimbo"
      },
      "meta": {
        "source_id": 10222
      }
    },
    {
      "id": "082bd15b-0e52-489b-9834-b31a07903b60",
      "project": "test",
      "data": {
        "name": "Bart",
        "some_other_count": 30
      },
      "meta": {
        "source_id": 10222
      }
    }
  ]
}

# single result for '"age": 30'
$ curl -s "localhost:5001/test?_search=%22age%22%3A%2030" | jq "."
{
  "items": [
    {
      "id": "45de46fd-af9e-4f54-84ed-241a9f803b2d",
      "project": "test",
      "data": {
        "age": 30,
        "name": "Jimbo"
      },
      "meta": {
        "source_id": 10222
      }
    }
  ]
}
```

### Run Locally

This project runs two docker containers; api and postgres.

```sh
# build without cache
rebuild:
	docker-compose build --no-cache

# build image
build:
	docker-compose build

# spin up containers (postgres and api)
up:
	docker-compose up -d

# bring down running containers
down:
	docker-compose down

# stop running containers
stop:
	docker-compose stop

# remove containers and images
remove:
	docker-compose rm -f

# cleanup images and volumes
clean:
	docker-compose down --rmi all -v --remove-orphans
```

### License

MIT Copyright &copy; Nick Sjostrom

Testing 
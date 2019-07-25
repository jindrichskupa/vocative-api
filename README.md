# Czech Vocative API Server

## Description

Public REST API for Czech Vocatives based on Minister of the Interior database

### TODO

* [ ] Automate load from CSV
* [ ] Generate exports - full (CSV, JSON, SQLite)
* [ ] Improvements...

## Build

### Download dependencies

```bash
dep ensure
```

### Make

```bash
make

# or cross-compile for Linux amd64
make linux64

# or docker image
make docker
```

## Run

### Configuration ENV variables

Variables and default values

```bash
VOCATIVE_LISTEN_IP=0.0.0.0
VOCATIVE_LISTEN_PORT=8080
VOCATIVE_DB_HOSTNAME=localhost
VOCATIVE_DB_PORT=5432
VOCATIVE_DB_USER=postgres
VOCATIVE_DB_PASSWORD=password
VOCATIVE_DB_NAME=vocativedb
VOCATIVE_DB_RETRIES=10
```

### Docker compose

1. update `docker-compose.yml` environment variables values
2. run `docker-compose up -d`

### Standalone

```bash
VOCATIVE_DB_NAME=vocative_db_test ./vocative-api
```

## Database setup

### Requirements

1. PostgreSQL
2. Basic PostgreSQL extenstions

### Setup

1. download ziped CSV files
2. extract .zip files
3. update paths in `import.sql` script
4. run `import.sql` script

### Data URLs

* Firstnames:
* Surnames:

## Usage

### Rest API

* [openapi.yml](./openapi.yml)

### API Call Examples

* search for firstnames based on part of name and gender

```bash
curl 'localhost:8080/vocative/firstnames/search?gender=male&name=Jan' | jq
```

* get all firstnames

```bash
curl 'localhost:8090/vocative/firstnames' | jq
```

* get vocative for a name

```bash
curl 'localhost:8080/vocative/firstnames/Jan' | jq
```

* get vocatives for name (firstname, surname and gender)

```bash
curl 'localhost:8090/vocative?firstname=Jindrich&surname=Skupa&gender=male&limit=2' | jq
```

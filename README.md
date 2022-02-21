# Golang Base API

[![Quality Assurance](https://github.com/brice-74/golang-base-api/actions/workflows/qa.yml/badge.svg)](https://github.com/brice-74/golang-base-api/actions/workflows/qa.yml)
[![codecov](https://codecov.io/gh/brice-74/golang-base-api/branch/master/graph/badge.svg?token=M5MV59TD3S)](https://codecov.io/gh/brice-74/golang-base-api)

Golang Base API is a GraphQL REST API implemented with JWT authentication and user session backup.

## Installation

You will need Docker installed on your machine. Please read the
[documentation here](https://docs.docker.com/get-docker/).

## Usage

Copy the `.env.example` file in a new `.env.dev` file and replace
the values and secrets if necessary.

:running::whale2: You can now start the project in two different ways, in both cases, the command will start the Docker containers with the API HTTP server and the Postgres database:

```bash
make run/dev/air # Start using cosmtrek/air
make run/dev/reflex # Start using cespare/reflex
```

:point_right::elephant: Make sure to perform migrations will containers are running:

```bash
make db/migrations/up # Create postgres tables
```

Enjoy ! :sunglasses:

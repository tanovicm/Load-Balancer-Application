# Distributed Expense Tracker

A distributed system for tracking user expenses, built with Go. It uses a microservices architecture with REST and gRPC communication.

## Features

- **JWT Authentication**: Register, login, (bonus: session invalidation)
- **Bank Accounts**: Create, fetch (with balance), delete
- **Expenses**: Create, list, delete

## Architecture

- **Load Balancer**: Exposes a REST API, handles JWT auth, converts REST to gRPC, routes requests to Workers using Round-Robin.
- **Worker Service**: Handles bank account and expense logic. Registers with Load Balancer over gRPC.

## Tech Stack

- Go, gRPC, REST, JWT

## Requirements

- Redis (Latest)
- Golang (1.14)

## Boot up

```
# start redis
redis-server

# start load balancer and two workers
./boot.sh
```

## Additional worker

```
go run worker/*.go <port>
```

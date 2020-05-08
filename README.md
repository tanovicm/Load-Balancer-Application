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
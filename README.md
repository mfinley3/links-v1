# links-v1 

A service for creating and managing short links.  
Upon creation of a short link, a client can immediately begin using it and get metrics about it.

Complete documentation/readme in progress...

## Setup

Most values are hardcoded. So the API should run out of the box with the required dependencies started.

### Start Dependencies 

From the project root run:

```bash
docker-compose -f docker/docker-compose.yml up
```

This will start both Postgres and Influxdb

### Start API 

From the project root run:

```bash
go run cmd/api/main.go
```



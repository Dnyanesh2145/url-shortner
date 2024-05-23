# Quickstart: compose-dev-env for khulke URL Shortener installation and setup
You can open this compose application with Docker Dev Environments feature of Docker Desktop version 4.12 or later.

[Click here to Open in Docker Dev Environments](https://open.docker.com/dashboard/dev-envs?url=https://bitbucket.org/loktantram-admin/khulke-url-shortener/src/master/) From Docker Desktop

Isolated Development can be done using Visual Studio code IDE, which you have to choose in docker environment in above step and install all recommended VS Code extensions shown in IDE.

This example is based on the [`nginx-golang-postgres`](https://github.com/docker/awesome-compose/tree/master/nginx-golang-postgres) sample of [`awesome-compose` repository](https://github.com/docker/awesome-compose/).

## Khulke Golang prototype application structure
### GoFiber Micro service API server with an Nginx proxy webserver, PostgreSQL database and Redis as queue/cache

This prototype is the base of all golang projects used for khulke product backend
Project structure:
```
.
├── Dockerfile
├── README.md
├── compose-dev.yaml
├── config
│   ├── constants.go
│   └── env.go
├── database
│   ├── database.go
│   └── redis.go
├── go.mod
├── go.sum
├── helpers
│   ├── helpers.go
│   └── models.go
├── main.go
├── migrations.sql
├── routes
│   ├── resolve.go
│   └── shorten.go
└── utils
    └── loggers.go
```

[_docker-compose.yaml_](compose-dev.yaml) used for development and running the services using docker-compose
```
services:
  backend:
    build: 
      context: .
      target: development
    ...
  db:
    image: postgres
    ...
  redis:
    image: redis
    ...
```
The compose file defines an application with 3 services `redis`, `db` and `backend`.
When deploying the application, docker-compose maps port 8080 of the proxy service container to port 8080(or 8081 some times in VS Code) of the host as specified in the file.
Make sure port 8080 on the host is not already being in use.

After the application starts, navigate to `http://localhost:8080` in your web browser or run:
```
$ curl http://localhost:8080/healthcheck

```

## Features
- low footprint on disk and memory
- extreme performance(with go-channels) and scalable in infra and development
- shortens URLs from whitelisted IPs via POST, for e.g:
    ```
    curl --location 'http://localhost:8081/v1/shorten' --header 'Content-Type: application/json' --data '{"url": "https://www.khulke.com/roundtable/all1"}'

    {"shortURL":"http://localhost:8081/u/Lmyefb"}           
    ```
- above shortURL will be 301 redirected to original URL on browsing
- health check is there at: http://localhost:8081/health

## Contributers
[dnyaneshrode07@gmail.com](mailto:dnyaneshrode07@gmail.com)

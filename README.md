# gin-gonic REST API using MongoDB
A simple Gin Gonic REST API using MongoDB & Redis

This a gin-gonic application, to provide an example on how to create REST API, integrated with 
MongoDB in Dropwizard with an OpenAPI specification. 

This example is a simple RESTful API to easily manage books I have read.

This is meant to be a playground project for all things gin-gonic, so the design or project structure is not something I 
would advocate.

This is my take on [Building Distributed Applications in Gin](https://github.com/PacktPublishing/building-distributed-applications-in-gin)
repository and (a great) book by [Mohamed Labouardy](https://www.labouardy.com/).

**Dependencies used:**
 * using [gin-gonic](https://github.com/gin-gonic/gin#gin-web-framework) v1.7.7 web framework
 * using [viper](https://github.com/spf13/viper) as a configuration solution
 * using [mongo-db](https://www.mongodb.com/) as NoSQL DB
 * using [redis](https://redis.io/) to cache `GET /books` and `GET /books/:id` resources
 * using [gin/sessions](github.com/gin-contrib/sessions) to handle session cookies
 * using [jwt-go](github.com/dgrijalva/jwt-go) to provide an implementation of JWT
 * using [x/crypto](golang.org/x/crypto), Go Cryptography package 
 * using [nancy](https://github.com/sonatype-nexus-community/nancy), tool to check for vulnerabilities in your Golang dependencies
 
## How to start the Gin-gonic application
```shell
go run main.go

# or via docker-compose
docker-compose up
```

## How to run audit and tests 
```shell
make audit test
```

## Swagger OpenAPI specification

* To generate the swagger spec file and have the API spec served from the Swagger UI
```shell
swagger generate spec --scan-models -o ./swagger.json && swagger serve  --port 8081 --path docs -F swagger ./swagger.json
```


## CURL commands to interact with the REST API
### Get all the books
```shell
curl -X GET 'localhost:8080/books'
```

### Search book by tag
```shell
curl -X GET 'localhost:8080/books/search?tag=nasa'
```

### Create a new book entry
```shell
curl -X POST 'localhost:8080/books' \
--data '{"name": "Moondust", "author": "Andrew Smith", "publisher": "Bloomsbury Publishing PLC", "published_at": {"month":"July", "year":"2009"}, "tags":["space exploration", "astronauts", "nasa"], "review":4.6}'
```

### Update a book entry
```shell
curl -X PUT 'localhost:8080/books/c8n5pb2kq9ndfcl9os7g' \
--data '{"name": "Moondust", "author": "Andrew Smith", "publisher": "Bloomsbury Publishing PLC", "published_at": {"month":"July", "year":"2009"}, "tags":["space exploration", "astronauts", "nasa", "JPL"], "review":4.7}'
```

### Delete a book entry
```shell
curl -X DELETE 'localhost:8080/books/c8n5pb2kq9ndfcl9os7g'
```

## Required Docker images

Create a MongoDB container
```shell
docker pull mongo
docker create --name mongodb -it -p 27017:27017 mongo
```

Create a Redis container
```shell
docker pull redis
docker create --name redis -it -p 6379:6379 redis
```

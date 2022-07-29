# gin-gonic REST API using MongoDB
A simple Gin Gonic REST API using MongoDB

This a gin-gonic application, to provide an example on how to create REST API, integrated with 
MongoDB in Dropwizard with an OpenAPI specification. 
This example is a simple RESTful API to easily manage books read.

 * using [gin-gonic](https://github.com/gin-gonic/gin#gin-web-framework) v1.7.7 web framework
 * using [viper](https://github.com/spf13/viper) as a configuration solution
 * using [mongo-db](https://www.mongodb.com/) as NoSQL DB
 * using [redis](https://redis.io/) to cache `GET /books` and `GET /books/:id` resources
 * using [gin sessions](github.com/gin-contrib/sessions v0.0.5) to handle session cookies
 * using [jwt-go](github.com/dgrijalva/jwt-go) to provide an implementation of JWT

## How to start the Gin-gonic application
```shell
go run pkg/main.go
```

## Swagger OpenAPI specification

* To generate the swagger spec file
```shell
swagger generate spec -o ./swagger.json
```

* To have the API specification served from the Swagger UI
```shell
swagger serve --port 8081 --path docs -F swagger ./swagger.json
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

## TODO

1. Improve error handling
2. Improve logging
3. Adds users resources and store them on the MongoDB
4. Return user profile
5. Create refresh endpoint
6. Testing
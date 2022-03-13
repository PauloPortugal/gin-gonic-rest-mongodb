package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	err := router.Run()
	if err != nil {
		return
	}
}

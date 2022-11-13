package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type test struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

var tests = []test{
	{ID: "1", Name: "hello"},
	{ID: "2", Name: "world"},
}

func getTest(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tests)
}

func postTest(c *gin.Context) {
	var newTest test
	if err := c.BindJSON(&newTest); err != nil {
		fmt.Println("Error")
	}
	tests = append(tests, newTest)
	c.IndentedJSON(http.StatusCreated, newTest)
}

func main() {

	router := gin.Default()
	router.GET("/test", getTest)
	router.POST("/test", postTest)

	router.Run("localhost:8000")
}

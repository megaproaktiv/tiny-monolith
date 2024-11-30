package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {

	var err error

	if err != nil {
		//
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", Index)
	r.Static("/css", "./css")
	r.GET("/process", Process)
	r.Run() // listen and serve on 0.0.0.0:8080

}

func handleSigTerms() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received SIGTERM, exiting")

		os.Exit(1)
	}()
}

type Status struct {
	Valid bool
}

func Process(c *gin.Context) {
	validation := Status{Valid: true}
	c.HTML(http.StatusOK, "process.html", validation)
	c.Next()
}

func Index(c *gin.Context) {
	validation := Status{Valid: true}
	c.HTML(http.StatusOK, "index.html", validation)
	c.Next()
}

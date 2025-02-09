package main

import (
	_ "embed"
	"fmt"
	"strings"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"htmx/log"

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
	r.Static("/include", "./include")
	r.POST("/process", Process)
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

func Index(c *gin.Context) {
	validation := Status{Valid: true}
	c.HTML(http.StatusOK, "index.html", validation)
	c.Next()
}

func Process(c *gin.Context) {
	validation := Status{Valid: true}
	goodBoy := c.PostForm("goodboy")
	log.Logger.Info("Process", "goodboy", goodBoy)
	if goodBoy == "me" {
		validation.Valid = true
	}
	if strings.ToLower(goodBoy) == "donald" {
		validation.Valid = false
	}
	c.HTML(http.StatusOK, "process.html", validation)
	c.Next()
}

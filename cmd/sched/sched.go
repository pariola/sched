package main

import (
	"log"
	"sched/handlers"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	h := handlers.New()

	e.POST("/create", h.CreateJob)
	e.POST("/update/:job", h.UpdateJob)
	e.POST("/cancel/:job", h.CancelJob)

	if err := e.Start(":5000"); err != nil {
		log.Fatal(err)
	}
}

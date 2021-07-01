package main

import (
	"gmgalvan/edChallenge2021/internal/m2m"
	"gmgalvan/edChallenge2021/internal/transport"
	"gmgalvan/edChallenge2021/internal/usecases"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	// setup m2m client
	client := &http.Client{}
	m2m := m2m.NewM2M(client)

	// usecases
	uc := usecases.NewNomicsReport(m2m)

	// transport
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	transport := transport.NewTrasport(uc)
	transport.InitRoutes(e)

	// Start server
	s := http.Server{
		Addr:    ":8080",
		Handler: e,
	}
	log.Println("Server Start")
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

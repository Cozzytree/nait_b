package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cozzytree/nait/internal/server"
	"github.com/joho/godotenv"
)

func gracefullShutdown(s *http.Server, doneChan chan bool) {
	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	ctx, stop = context.WithTimeout(
		context.Background(),
		time.Second*2)

	defer stop()

	log.Println("server shutting")
	<-ctx.Done()

	err := s.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	doneChan <- true
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	s := server.InitServer()

	doneChan := make(chan bool, 1)

	go gracefullShutdown(s, doneChan)

	log.Println("server started on port", s.Addr)

	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	<-doneChan
	log.Println("server shutting down")
}

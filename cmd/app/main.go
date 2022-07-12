package main

import (
	"fmt"
	"github.com/Vallghall/wb-test-l0/internal/handlers"
	"github.com/Vallghall/wb-test-l0/internal/services"
	"github.com/Vallghall/wb-test-l0/internal/storage"
	"github.com/Vallghall/wb-test-l0/internal/storage/postgres"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %v\n", err)
	}

	db := postgres.NewConnection(&postgres.Configs{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PW"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})

	s := storage.New(db)
	srv := services.New(s)
	srv.SynchronizeCash()
	h := handlers.New(srv)

	urlOption := stan.NatsURL(fmt.Sprintf(
		"nats://%s:%s",
		os.Getenv("NATS_HOST"),
		os.Getenv("NATS_PORT")))

	sc, err := stan.Connect("nats-test", "test-client", urlOption)
	if err != nil {
		log.Fatalf("nats streaming connection failed: %v\n", err)
	}

	ch := make(chan struct{})
	go subscribe(ch, sc, srv)
	go systemShutdown(ch, sc, db)
	log.Println("haha")
	log.Fatalf("server shutting down: %v\n", h.Routes().Run())
}

func subscribe(ch chan struct{}, sc stan.Conn, s *services.Service) {
	sub, err := sc.Subscribe(
		"test",
		func(msg *stan.Msg) {
			log.Printf("Received message: %s\n", string(msg.Data))
			s.DataService.CashMessage(msg.Data)
			s.DataService.PersistMessage(msg.Data)
		},
		stan.DeliverAllAvailable(),
		stan.DurableName("dur-test"),
	)
	if err != nil {
		log.Fatalf("subscription failed: %v", err)
	}
	log.Println("Successfully subscribed!!!")

	defer func() {
		sub.Unsubscribe()
		sub.Close()
	}()

	<-ch
	log.Println("blah")
}

func systemShutdown(ch chan<- struct{}, leftovers ...io.Closer) {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	sig := <-sign
	log.Printf("received signal %v. SHUTTING DOWN\n", sig)
	ch <- struct{}{}
	for _, closer := range leftovers {
		if err := closer.Close(); err != nil {
			log.Printf("failed to close resource: %v", err)
		}
	}
}

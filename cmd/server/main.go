package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/whiterthanwhite/businessinsight/internal/db"
)

type ServerLogger struct {
	Handler http.Handler
}

func (s *ServerLogger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request URI: %s\n", req.RequestURI)
	log.Printf("Method: %s\n", req.Method)

	s.Handler.ServeHTTP(w, req)
}

func main() {
	dbConnectionStr := os.Getenv("DBCONNECTIONSTR")
	if dbConnectionStr == "" {
		log.Fatalln("Database connections string is not specified!")
	}

	ctx, cancel := context.WithCancel(context.Background())

	interruptSig := make(chan os.Signal, 1)
	signal.Notify(interruptSig, os.Interrupt)
	go func() {
		<-interruptSig
		cancel()
	}()

	conn, err := db.Conntect(ctx, dbConnectionStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	err = conn.InitTables(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/helloworld", func(w http.ResponseWriter, req *http.Request) {
		sqlMessage, err := conn.HelloWorld(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(sqlMessage))
	})
	mux.HandleFunc("/currtime", func(w http.ResponseWriter, req *http.Request) {
		sqlMessage, err := conn.CurrentTime(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(sqlMessage))
	})

	sl := &ServerLogger{
		Handler: mux,
	}

	go func() {
		log.Println("Server started")
		if err := http.ListenAndServe(":8080", sl); err != nil {
			fmt.Println(err)
		}
	}()

	<-ctx.Done()
	fmt.Print("\r")
	log.Println("Server stopped")
}

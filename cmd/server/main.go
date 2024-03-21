package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/whiterthanwhite/businessinsight/internal/db"
	"github.com/whiterthanwhite/businessinsight/internal/entities/currency"
	"github.com/whiterthanwhite/businessinsight/internal/middleware"
)

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

	conn, err := db.Connect(ctx, dbConnectionStr)
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
	mux.HandleFunc("/currency/add", func(w http.ResponseWriter, req *http.Request) {
		currenciesJSON, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currencies, err := currency.ParseJSON(currenciesJSON)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = conn.InsertCurrency(ctx, currencies)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("currencies added")
	})
	mux.HandleFunc("/currencies", func(w http.ResponseWriter, req *http.Request) {
		currencies, err := conn.GetCurrencies(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currenciesJSON, err := json.Marshal(currencies)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(currenciesJSON)
	})

	sl := &middleware.ServerLogger{
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

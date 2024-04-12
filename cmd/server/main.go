package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/whiterthanwhite/businessinsight/internal/db"
	"github.com/whiterthanwhite/businessinsight/internal/handlerfunctions"
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

	mux, err := createCustomMux(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	rh := &middleware.ReactHelper{
		Handler: mux,
	}
	sl := &middleware.ServerLogger{
		Handler: rh,
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

func createCustomMux(parentCtx context.Context) (*http.ServeMux, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	conn, err := db.GetInstance()
	if err != nil {
		return nil, err
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

	mux.HandleFunc("/currencies/add", handlerfunctions.AddCurrenciesHandlerFunc())
	mux.HandleFunc("/currencies", handlerfunctions.GetCurrenciesHandlerFunc())
	mux.HandleFunc("/currencies/delete", handlerfunctions.DeleteCurrenciesHandlerFunc())

	mux.HandleFunc("/accounts", handlerfunctions.GetAccountsHandlerFunction())
	mux.HandleFunc("/accounts/add", handlerfunctions.AddAccountsHandlerFunction())
	mux.HandleFunc("/accounts/delete", handlerfunctions.DeleteAccountsHandlerFunction())

	mux.HandleFunc("/categories", handlerfunctions.GetCategoriesHandlerFunction())
	mux.HandleFunc("/categories/add", handlerfunctions.AddCategoryHandlerFunction())
	mux.HandleFunc("/categories/delete", handlerfunctions.DeleteCategoriesHandlerFunctions())

	mux.HandleFunc("/operations", handlerfunctions.GetOperationsHandlerFunction())
	mux.HandleFunc("/operations/add", handlerfunctions.AddOperationsHandlerFunction())
	mux.HandleFunc("/operations/delete", handlerfunctions.DeleteOperationsHandlerFunction())

	mux.HandleFunc("/accountStatistics", handlerfunctions.GetAccountStatisticsHandlerFunction())

	return mux, nil
}

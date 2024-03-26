package handlerfunctions

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/whiterthanwhite/businessinsight/internal/db"
	"github.com/whiterthanwhite/businessinsight/internal/entities/currency"
)

func AddCurrenciesHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn := db.GetInstance()
		if conn == nil {
			err := errors.New("Database variable was not initialized!")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currenciesJSON, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currencies, err := currency.ParseJSON(currenciesJSON)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = conn.InsertCurrency(ctx, currencies)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("currencies added")
	}
}

func GetCurrenciesHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn := db.GetInstance()
		if conn == nil {
			err := errors.New("Database variable was not initialized!")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currencies, err := conn.GetCurrencies(ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(currencies)

		currenciesJSON, err := json.Marshal(currencies)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(currenciesJSON)
	}
}

func DeleteCurrenciesHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn := db.GetInstance()
		if conn == nil {
			err := errors.New("Database variable was not initialized!")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currenciesJSON, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currencies, err := currency.ParseJSON(currenciesJSON)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = conn.DeleteCurrencies(ctx, currencies)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Currencies deleted!")
	}
}

package handlerfunctions

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/whiterthanwhite/businessinsight/internal/db"
	"github.com/whiterthanwhite/businessinsight/internal/entities/account"
	"github.com/whiterthanwhite/businessinsight/internal/entities/category"
	"github.com/whiterthanwhite/businessinsight/internal/entities/currency"
)

// Currency handler functions
func GetCurrenciesHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn, err := db.GetInstance()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currencies, err := conn.GetCurrencies(ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currenciesJSON, err := json.Marshal(currencies)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(string(currenciesJSON))
		w.Write(currenciesJSON)
	}
}

func AddCurrenciesHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn, err := db.GetInstance()
		if err != nil {
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

		for _, newCurrency := range currencies {
			xCurrency, err := conn.GetCurrency(ctx, &newCurrency)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if xCurrency != nil {
				if xCurrency.Description != newCurrency.Description {
					err = conn.UpdateCurrency(ctx, &newCurrency)
					if err != nil {
						log.Println(err.Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			} else {
				err = conn.InsertCurrency(ctx, &newCurrency)
				if err != nil {
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		log.Println("currencies added")
	}
}

func DeleteCurrenciesHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn, err := db.GetInstance()
		if err != nil {
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

// Account handler functions
func GetAccountsHandlerFunction() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn, err := db.GetInstance()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accounts, err := conn.GetAccounts(ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accountsJson, err := json.Marshal(&accounts)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(string(accountsJson))
		w.Write(accountsJson)
	}
}

func AddAccountsHandlerFunction() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn, err := db.GetInstance()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newAccounts, err := account.ParseJSON(requestBody)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, newAccount := range newAccounts {
			xAccount, err := conn.GetAccount(ctx, &newAccount)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if xAccount != nil {
				if xAccount.Name != newAccount.Name || xAccount.CurrencyCode != newAccount.CurrencyCode {
					err = conn.UpdateAccount(ctx, &newAccount)
					if err != nil {
						log.Println(err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			} else {
				err = conn.InsertAccount(ctx, &newAccount)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}
}

func DeleteAccountsHandlerFunction() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn, err := db.GetInstance()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accounts, err := account.ParseJSON(requestBody)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = conn.DeleteAccount(ctx, &accounts[0])
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// category handler fucntions
func AddCategoryHandlerFunction() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categories, err := category.ParseJSON(requestBody)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conn, err := db.GetInstance()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, category := range categories {
			xCategory, err := conn.GetCategory(ctx, &category)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if xCategory != nil {
				if xCategory.Type != category.Type || xCategory.Name != category.Name ||
					xCategory.Description != category.Description {

					err = conn.UpdateCategory(ctx, &category)
					if err != nil {
						log.Println(err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			} else {
				err = conn.InsertCategory(ctx, &category)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		log.Println(`categories were added`)
	}
}

func GetCategoriesHandlerFunction() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn, err := db.GetInstance()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categories, err := conn.GetCategories(ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBody, err := json.Marshal(&categories)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(responseBody)
	}
}

func DeleteCategoriesHandlerFunctions() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categories, err := category.ParseJSON(requestBody)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conn, err := db.GetInstance()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, category := range categories {
			err = conn.DeleteCategory(ctx, &category)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

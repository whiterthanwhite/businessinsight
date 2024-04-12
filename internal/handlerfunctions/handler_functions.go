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
	"github.com/whiterthanwhite/businessinsight/internal/entities/operation"
	"github.com/whiterthanwhite/businessinsight/internal/entities/operation_type"
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

// Operation handler fucntions
func AddOperationsHandlerFunction() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		var operations []operation.Operation
		err = json.Unmarshal(requestBody, &operations)
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		conn, err := db.GetInstance()
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		var fromOperationSet, toOperationSet bool
		var lastTransactionNo int
		for _, operation := range operations {
			xOperation, err := conn.GetOperation(ctx, &operation)
			if err != nil {
				log.Println(err)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			if xOperation == nil {
				if operation.Type == operation_type.Transfer {
					if fromOperationSet && toOperationSet {
						fromOperationSet, toOperationSet = false, false
					}
					if !fromOperationSet && !toOperationSet {
						lastTransactionNo, err = conn.GetMaxTransactionNo(ctx)
						if err != nil {
							log.Println(err)
							http.Error(rw, err.Error(), http.StatusInternalServerError)
							return
						}
						lastTransactionNo++
					}
					operation.TransactionNo = lastTransactionNo
					if operation.Amount < 0 {
						fromOperationSet = true
					} else {
						toOperationSet = true
					}
				}
				if err = conn.InsertOperation(ctx, &operation); err != nil {
					log.Println(err)
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				if !xOperation.Compare(&operation) {
					if err = conn.UpdateOperation(ctx, &operation); err != nil {
						log.Println(err)
						http.Error(rw, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}
		}
	}
}

func GetOperationsHandlerFunction() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn, err := db.GetInstance()
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		operations, err := conn.GetOperations(ctx)
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBody, err := json.Marshal(&operations)
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Write(responseBody)
	}
}

func DeleteOperationsHandlerFunction() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn, err := db.GetInstance()
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		var operations []operation.Operation
		if err = json.Unmarshal(requestBody, &operations); err != nil {
			log.Println(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, operation := range operations {
			if err = conn.DeleteOperation(ctx, &operation); err != nil {
				log.Println(err)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

// Statics handler functions
func GetAccountStatisticsHandlerFunction() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		conn, err := db.GetInstance()
		if err != nil {
			log.Println(err.Error())
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		accountsStatistics, err := conn.GetAccountStatistics(ctx)
		if err != nil {
			log.Println(err.Error())
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		responseBodyJson, err := json.Marshal(&accountsStatistics)
		if err != nil {
			log.Println(err.Error())
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Write(responseBodyJson)
	}
}

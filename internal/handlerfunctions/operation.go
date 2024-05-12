package handlerfunctions

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/whiterthanwhite/businessinsight/internal/db"
	"github.com/whiterthanwhite/businessinsight/internal/entities/operation"
	"github.com/whiterthanwhite/businessinsight/internal/entities/operation_type"
)

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
			operation.CreationDate = operation.DateTime
			operation.CreationTime = operation.DateTime
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

package boltdb

import (
	"encoding/json"
	"errors"

	"github.com/pieterlouw/go-jsonapigateway_tmpl/gateway"

	"github.com/boltdb/bolt"
)

var (
	// ErrTranNotFound is returned when transaction is not found in database
	ErrTranNotFound = errors.New("transaction not found")
)

//AuthTransactionResponseService represents bolt version of gateway.AuthTransactionResponseService
type AuthTransactionResponseService struct {
	Db         *bolt.DB
	BucketName string
}

//Set sets the tranID key to the value of newItem
func (s AuthTransactionResponseService) Set(tranID string, newItem gateway.AuthTransactionResponse) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(s.BucketName))
		if err != nil {
			return err
		}
		encoded, err := json.Marshal(newItem)
		if err != nil {
			return err
		}
		return b.Put([]byte(tranID), encoded)
	})
}

//Get value for key tranID
func (s AuthTransactionResponseService) Get(tranID string) (gateway.AuthTransactionResponse, error) {
	var rsp gateway.AuthTransactionResponse
	var err error
	var dbValue []byte

	//read value from db bucket
	err = s.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.BucketName))

		if b == nil {
			return ErrTranNotFound
		}
		dbValue = b.Get([]byte(tranID))
		if dbValue == nil {
			return ErrTranNotFound
		}
		return nil
	})

	if err != nil {
		return rsp, err
	}

	if err = json.Unmarshal(dbValue, &rsp); err != nil {
		return rsp, err
	}
	return rsp, err
}

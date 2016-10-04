package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/pieterlouw/go-jsonapigateway_tmpl/boltdb"
	"github.com/pieterlouw/go-jsonapigateway_tmpl/gateway"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

const appVersion = "0.1.0"

var (
	lookupTable map[string]string
	appConfig   gateway.AppConfig
	err         error
)

func init() {
	appConfig, err = gateway.ReadConfig("config.json")
	if err != nil {
		fmt.Printf("There was a problem reading the config file. Error=%v", err)
		os.Exit(1)
	}

}

func main() {
	log.Printf("msg=Application started	version=%s\n", appVersion)
	lookupTable = make(map[string]string)

	//Opens BoltDB file. It will be created if it doesn't exist.
	db, err := bolt.Open(appConfig.BoltDBName, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(appConfig.BoltDBTranRspBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	authResponseStore := &boltdb.AuthTransactionResponseService{
		Db:         db,
		BucketName: appConfig.BoltDBTranRspBucket,
	}

	//load lookup table
	contents, err := ioutil.ReadFile(appConfig.LookupFileName)
	if err != nil {
		log.Fatalln(err)
	}

	for _, line := range strings.Split(string(contents), "\n") {
		data := strings.Split(line, "=")
		if len(data) == 2 {
			lookupTable[data[0]] = strings.Replace(data[1], "\r", "", -1)
		}
	}

	if len(lookupTable) > 0 {
		log.Printf("Numbers loaded in lookup table: %d", len(lookupTable))
	} else {
		log.Fatalln("No items in lookup table")
	}

	gateway.NewRemoteAPI(appConfig.RemoteURL, appConfig.APIUsername, appConfig.APIPassword, appConfig.APIClientID, lookupTable, authResponseStore)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/report/payments/{transactionID}", gateway.PaymentReport).Methods("GET")
	router.Handle("/api/v1/ping", handler(gateway.Ping)).Methods("GET")
	router.Handle("/api/v1/payments", handler(gateway.Payment)).Methods("POST")
	router.Handle("/api/v1/payments/{transactionID}", handler(gateway.QueryPayment)).Methods("GET")

	log.Fatal(http.ListenAndServe(appConfig.ListeningPort, router))

}

//handle defintion to add some panic and error logging middleware to http handlers
type handler func(resp http.ResponseWriter, req *http.Request) (gateway.Result, error)

func (h handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	runHandler(resp, req, h)
}

func runHandler(w http.ResponseWriter, r *http.Request, fn func(http.ResponseWriter, *http.Request) (gateway.Result, error)) {
	var err error

	defer func() {
		if rv := recover(); rv != nil {
			err = errors.New("handler panic")
			logError(r, err, rv)
			handleError(w, r, http.StatusInternalServerError, err)
		}
	}()

	result, err := fn(w, r)
	if err != nil {
		logError(r, err, nil)
		handleError(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(result); err != nil {
		logError(r, err, r)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.Header().Set("cache-control", "no-cache")

	w.WriteHeader(status)
}

func logError(req *http.Request, err error, rv interface{}) {
	if err != nil {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Error serving %s (route %s): %s\n", req.URL, mux.CurrentRoute(req).GetName(), err)
		if rv != nil {
			fmt.Fprintln(&buf, rv)
			buf.Write(debug.Stack())
		}
		log.Print(buf.String())
	}
}

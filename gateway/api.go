package gateway

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (

	//RemoteAPI holds info to be used to make calls to remote verification API
	api *remoteAPI
)

//NewRemoteAPI initializes API
func NewRemoteAPI(url string, authUser string, authPass string, clientID string, lookupTable map[string]string, service AuthTransactionResponseService) {
	api = &remoteAPI{
		URL:                            url,
		BasicAuthUsername:              authUser,
		BasicAuthPassword:              authPass,
		ClientID:                       clientID,
		MSISDNLookupTable:              lookupTable,
		AuthTransactionResponseService: service,
	}
}

//RemoteAPI holds info to be used to make calls to remote verification API
type remoteAPI struct {
	URL                            string
	BasicAuthUsername              string
	BasicAuthPassword              string
	ClientID                       string
	MSISDNLookupTable              map[string]string
	AuthTransactionResponseService AuthTransactionResponseService
}

//AuthTransactionRequest data to send to API
type AuthTransactionRequest struct {
	Amount        int    `json:"amount_in_cents"`
	APIClientID   string `json:"clientid"`
	TransactionID string `json:"tranid"`
	UserID        string `json:"userid"`
}

//PingResponse data response from /ping API method
type PingResponse struct {
	Success    bool   `json:"success"`
	ResultText string `json:"resultText"`
}

//AuthTransactionResponse data from /Auth API method
type AuthTransactionResponse struct {
	Success       bool   `json:"success"`
	ResultText    string `json:"resultText"`
	Amount        int    `json:"amount_in_cents"`
	TransactionID string `json:"tranid"`
	UserID        string `json:"userid"`
	Authorized    bool   `json:"authorized"`
}

//AuthTransactionResponseService interface type
type AuthTransactionResponseService interface {
	Set(tranID string, newItem AuthTransactionResponse) error
	Get(tranID string) (AuthTransactionResponse, error)
}

func (r remoteAPI) doRequest(method string, function string, payload []byte) ([]byte, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s%s", r.URL, function)
	log.Printf("URL: %s\n", url)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return []byte(""), err
	}

	req.SetBasicAuth(r.BasicAuthUsername, r.BasicAuthPassword)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	// Slurp up the response's body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	// Close the response's body.
	if err := resp.Body.Close(); err != nil {
		return []byte(""), err
	}

	return body, nil
}

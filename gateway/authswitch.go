package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Result general structure to be used as return value in handler calls
type Result struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

//PaymentRequest define the parameters that will be received by Payment handler
type PaymentRequest struct {
	Amount   int    `json:"amount"`
	DeviceID string `json:"deviceID"`
	MSISDN   string `json:"msisdn"`
	TSN      string `json:"tsn"`
}

//Ping function handler that will translate local api call to remote api
func Ping(w http.ResponseWriter, r *http.Request) (Result, error) {
	var pingRsp PingResponse
	var result Result

	rsp, err := api.doRequest("GET", "isAlive", nil)
	if err != nil {
		return result, err
	}

	if err = json.Unmarshal(rsp, &pingRsp); err != nil {
		return result, err
	}

	if pingRsp.Success == false {
		result.Code = 99
		result.Text = fmt.Sprintf("Service not available Status=[%s]", pingRsp.ResultText)
	} else {
		result.Code = 0
		result.Text = "Service is up"
	}

	return result, err
}

//Payment function handler that will translate local api call to remote call
func Payment(w http.ResponseWriter, r *http.Request) (Result, error) {
	var result Result
	var AuthReq AuthTransactionRequest
	var AuthRsp AuthTransactionResponse
	var output bytes.Buffer
	var paymentRequest PaymentRequest

	//get and parse request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return result, err
	}

	if err := r.Body.Close(); err != nil {
		return result, err
	}

	if err = json.Unmarshal(body, &paymentRequest); err != nil {
		result.Code = 400
		result.Text = "Invalid request"
		return result, err
	}

	log.Printf("msg=Rx Auth Request DeviceID=%s TSN=%s Amount=%d MSISDN=%s\n", paymentRequest.DeviceID, paymentRequest.TSN, paymentRequest.Amount, paymentRequest.MSISDN)

	//lookup relevant userID for MSISDN
	userID, ok := api.MSISDNLookupTable[paymentRequest.MSISDN]
	if ok == false {
		result.Code = 403
		result.Text = "No UserID found for MSISDN"
		return result, err
	}

	tranID := fmt.Sprintf("%s-%s", paymentRequest.DeviceID, paymentRequest.TSN)
	AuthReq.APIClientID = api.ClientID
	AuthReq.TransactionID = tranID
	AuthReq.UserID = userID
	AuthReq.Amount = paymentRequest.Amount

	payload, err := json.Marshal(AuthReq)
	if err != nil {
		return result, err
	}

	log.Printf("Tx Auth Request: %s", payload)

	rsp, err := api.doRequest("POST", "auth", payload)
	if err != nil {
		return result, err
	}
	log.Printf("Auth Response: %s", rsp)

	if err = json.Unmarshal([]byte(rsp), &AuthRsp); err != nil {
		return result, err
	}

	err = api.AuthTransactionResponseService.Set(tranID, AuthRsp)
	if err != nil {
		log.Printf("Could not save transaction response to database Error=%s\n", err)
	}

	if AuthRsp.Success == false {
		result.Code = 1
		output.WriteString(fmt.Sprintf("FAILED (%d)", AuthRsp.ResultText))
		result.Text = output.String()
	} else {
		if AuthRsp.Authorized {
			result.Code = 0
			output.WriteString("APPROVED")
			result.Text = output.String()
		} else {
			result.Code = 2
			output.WriteString(fmt.Sprintf("DECLINED (%d)", AuthRsp.ResultText))
			result.Text = output.String()
		}
	}

	return result, err
}

//QueryPayment function handler that will translate local api call to remote call
func QueryPayment(w http.ResponseWriter, r *http.Request) (Result, error) {
	var result Result
	var AuthRsp AuthTransactionResponse
	var output bytes.Buffer
	tsn := mux.Vars(r)["transactionID"]

	rsp, err := api.doRequest("GET", fmt.Sprintf("auth/%s", tsn), nil)
	if err != nil {
		return result, err
	}
	if err = json.Unmarshal(rsp, &AuthRsp); err != nil {
		return result, err
	}

	if AuthRsp.Success == false {
		result.Code = 1
		output.WriteString(fmt.Sprintf("FAILED (%d)", AuthRsp.ResultText))
		result.Text = output.String()
	} else {
		if AuthRsp.Authorized {
			result.Code = 0
			output.WriteString("APPROVED")
			result.Text = output.String()
		} else {
			result.Code = 2
			output.WriteString(fmt.Sprintf("DECLINED (%d)", AuthRsp.ResultText))
			result.Text = output.String()
		}
	}
	return result, err
}

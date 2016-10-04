package gateway

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var templ = template.Must(template.New("main").Parse(templateStr))
var templErr = template.Must(template.New("err").Parse(templateErrStr))

//PaymentReport handler function to generate HTML report for given transactionID
func PaymentReport(w http.ResponseWriter, r *http.Request) {
	tsn := mux.Vars(r)["transactionID"]

	w.Header().Set("cache-control", "no-cache")
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	rsp, err := api.AuthTransactionResponseService.Get(tsn)
	//log.Printf("api.VerifyTransactionResponseService.Get: [%v] [%v]", rsp, err)
	if err != nil {
		err = templErr.Execute(w, &struct {
			Error error
		}{
			Error: err,
		})
	} else {
		templ.Execute(w, &struct {
			Data AuthTransactionResponse
		}{
			Data: rsp,
		})
	}

}

const templateStr = `
<html>
<head>
<title>View Transaction</title>
</head>
<body>

{{if .Data}}
    <h1>View Transaction - {{.Data.TransactionID }}</h1>
    <p><b>Success:</b>{{.Data.Success}}</p>
    <p><b>ResultText:</b>{{.Data.ResultText}}</p>
    <p><b>Amount (cents):</b>{{.Data.Amount}}</p>
    <p><b>UserID:</b>{{.Data.UserID}}</p>
    <p><b>Authorized:</b>{{.Data.Authorized}}</p>
{{end}}
</body>
</html>`

const templateErrStr = `
<html>
<head>
<title>View Transaction</title>
</head>
<body>

{{if .Error}}
    <p>Error:{{.Error}}</p>
{{end}}

</body>
</html>`

package gateway

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//AppConfig properties used in app
type AppConfig struct {
	ListeningPort       string `json:"listeningPort"`
	RemoteURL           string `json:"remoteURL"`
	APIUsername         string `json:"apiUsername"`
	APIPassword         string `json:"apiPassword"`
	APIClientID         string `json:"apiClientID"`
	EnableTrace         bool   `json:"enableTrace"`
	EnableDebug         bool   `json:"enableDebug"`
	LogToStderr         bool   `json:"logToStderr"`
	LogToFile           bool   `json:"logToFile"`
	LogFilePath         string `json:"logFilePath"`
	LookupFileName      string `json:"lookupFileName"`
	BoltDBName          string `json:"boltDBName"`
	BoltDBTranRspBucket string `json:"boltDBTranRspBucket"`
}

// ReadConfig reads info from config file
func ReadConfig(configfile string) (AppConfig, error) {
	_, err := os.Stat(configfile)
	if err != nil {
		return AppConfig{}, err
	}

	data, err := ioutil.ReadFile(configfile) // For read access.
	if err != nil {
		return AppConfig{}, err
	}

	var config AppConfig

	err = json.Unmarshal(data, &config)
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}

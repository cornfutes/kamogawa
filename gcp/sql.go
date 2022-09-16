package gcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"kamogawa/types"
	"log"
	"net/http"
)

type ErrorGloudSQLListError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorGloudSQLList struct {
	Error ErrorGCFListError `json:"error"`
}

type CloudSQLDBInstance struct {
	Name string `json:"name"`
}

type ResponeCloudSQLList struct {
	Items []CloudSQLDBInstance `json:"items"`
}

var apiHostCloudSQL = "https://sqladmin.googleapis.com"

var ScopeCloudSQL = "https://www.googleapis.com/auth/sqlservice.admin"

func CloudSQLListInstances(user types.User, projectId string) (*ResponeCloudSQLList, *ErrorGloudSQLList) {
	url := apiHostCloudSQL + "/v1/projects/" + projectId + "/instances"
	var bearer = "Bearer " + user.AccessToken.String
	fmt.Printf("Request "+url+" with token %v\n", bearer)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Error while making request to API")
	}
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error on response. %v\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}

	tee := io.TeeReader(resp.Body, buf)

	reader1, _ := ioutil.ReadAll(tee)
	reader2, _ := ioutil.ReadAll(buf)

	var responseSuccess ResponeCloudSQLList
	err = json.Unmarshal(reader1, &responseSuccess)
	if err != nil {
		log.Printf("Original %v\n", reader1)
		panic(reader1)
	}

	var responseError ErrorGloudSQLList
	err = json.Unmarshal(reader2, &responseError)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Error %v\n", string(reader2))

	return &responseSuccess, &responseError
}

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
	// https://www.golangprograms.com/dynamic-json-parser-without-struct-in-golang.html
)

type ErrorGCFListError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorGCFList struct {
	Error ErrorGCFListError `json:"error"`
}

type GCFFunction struct {
	Name  string `json:"name"`
	State string `json:"status"`
}

type ResponseGFFList struct {
	Functions []GCFFunction `json:"functions"`
}

// TODO: dynamicall;y retrieve this.
// var GCFLocations [25]string = [25]string{
// 	"asia-east1",
// 	"asia-east2",
// 	"asia-northeast1",
// 	"asia-northeast2",
// 	"asia-northeast3",
// 	"asia-south1",
// 	"asia-southeast1",
// 	"asia-southeast2",
// 	"australia-southeast1",
// 	"europe-central2",
// 	"europe-north1",
// 	"europe-west1",
// 	"europe-west2",
// 	"europe-west3",
// 	"europe-west4",
// 	"europe-west6",
// 	"northamerica-northeast1",
// 	"southamerica-east1",
// 	"us-central1",
// 	"us-east1",
// 	"us-east4",
// 	"us-west1",
// 	"us-west2",
// 	"us-west3",
// 	"us-west4",
// }

var apiHost = "https://cloudfunctions.googleapis.com"

var GCFScopes = "https://www.googleapis.com/auth/cloud-platform"

func GCFListFunctions(user types.User, projectId string, locationId *string) (*ResponseGFFList, *ErrorGCFList) {
	// Query globally
	var parent string
	if locationId == nil {
		parent = "projects/" + projectId + "/locations/-"
	} else {
		parent = "projects/" + projectId + "/locations/" + *locationId
	}

	url := apiHost + "/v1/" + parent + "/functions"
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

	var responseSuccess ResponseGFFList
	err = json.Unmarshal(reader1, &responseSuccess)
	if err != nil {
		log.Printf("Original %v\n", reader1)
		panic(reader1)
	}

	var responseError ErrorGCFList
	err = json.Unmarshal(reader2, &responseError)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Error %v\n", string(reader2))

	return &responseSuccess, &responseError
}

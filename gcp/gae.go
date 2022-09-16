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

// {
//   "services": [
//     {
//       "name": "apps/linear-cinema-360910/services/default",
//       "id": "default",
//       "split": {
//         "allocations": {
//           "20220830t021415": 1
//         }
//       }
//     }
//   ]
// }

type GAEServiceTrafficeAllocation struct {
}

type GAEServiceTrafficAllocations struct {
	Allocations GAEServiceTrafficeAllocation `json:"allocations"`
}

// https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps.services#Service
type GAEService struct {
	Name  string                       `json:"name"`
	Id    string                       `json:"id"`
	Split GAEServiceTrafficAllocations `json:"split"`
}

type GAEListServicesResponse struct {
	Services []GAEService `json:"services"`
}

//
// Response {
// 	  "error": {
// 	    "code": 404,
// 	    "message": "App does not exist.",
// 	    "status": "NOT_FOUND"
// 	   }
// 	}

// https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps.services.versions#Version
type GAEVersion struct {
	Name          string `json:"name"`
	Id            string `json:"id"`
	ServingStatus string `json:"servingStatus"`
}

type GAEListVersionsResponse struct {
	Versions []GAEVersion `json:"versions"`
}

// https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps.services.versions.instances#Instance
type GAEInstance struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	VMNamm string `json:"vmName"`
}

type GAEListInstancesResponse struct {
	Instances     []GAEVersion `json:"instances"`
	NextPageToken string       `json:"nextPageToken"`
}

type ErrorAdminAPI struct {
	Error ErrorAdminAPIError `json:"error"`
}
type ErrorAdminAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps.services/list
func GAEListServices(user types.User, projectId string) (*GAEListServicesResponse, *ErrorAdminAPI) {
	apiAdminApiUrl := "https://appengine.googleapis.com/v1/apps/" + projectId + "/services"
	var bearer = "Bearer " + user.AccessToken.String
	fmt.Printf("Token %v\n", bearer)

	req, err := http.NewRequest("GET", apiAdminApiUrl, nil)
	if err != nil {
		panic("Error while making request to project")
	}
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	apiProjectsResp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer apiProjectsResp.Body.Close()

	buf := &bytes.Buffer{}
	tee := io.TeeReader(apiProjectsResp.Body, buf)
	reader1, _ := ioutil.ReadAll(tee)
	reader2, _ := ioutil.ReadAll(buf)

	fmt.Printf("Response %v \n", buf.String())

	var responseSuccess GAEListServicesResponse
	err = json.Unmarshal(reader1, &responseSuccess)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Services %v\n", responseSuccess)
	var responseError ErrorAdminAPI
	err = json.Unmarshal(reader2, &responseError)
	if err != nil {
		panic(err)
	}
	return &responseSuccess, &responseError
}

// https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps.services.versions/list
func GAEListVersions(user types.User, projectId string, serviceId string) (*GAEListVersionsResponse, *ErrorAdminAPI) {
	apiAdminApiUrl := "https://appengine.googleapis.com/v1/apps/" + projectId + "/services/" + serviceId + "/versions/"
	var bearer = "Bearer " + user.AccessToken.String
	fmt.Printf("Token %v\n", bearer)

	req, err := http.NewRequest("GET", apiAdminApiUrl, nil)
	if err != nil {
		panic("Error while making request to project")
	}
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	apiProjectsResp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer apiProjectsResp.Body.Close()

	buf := &bytes.Buffer{}
	tee := io.TeeReader(apiProjectsResp.Body, buf)
	reader1, _ := ioutil.ReadAll(tee)
	reader2, _ := ioutil.ReadAll(buf)

	fmt.Printf("Response %v \n", buf.String())

	var responseSuccess GAEListVersionsResponse
	err = json.Unmarshal(reader1, &responseSuccess)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Services %v\n", responseSuccess)
	var responseError ErrorAdminAPI
	err = json.Unmarshal(reader2, &responseError)
	if err != nil {
		panic(err)
	}
	return &responseSuccess, &responseError
}

// https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps.services.versions.instances/list
func GAEListInstances(user types.User, projectId string, serviceId string, versionId string) (*GAEListInstancesResponse, *ErrorAdminAPI) {
	apiAdminApiUrl := "https://appengine.googleapis.com/v1/apps/" + projectId + "/services/" + serviceId + "/versions/" + versionId + "/instances"
	var bearer = "Bearer " + user.AccessToken.String
	fmt.Printf("Token %v\n", bearer)

	req, err := http.NewRequest("GET", apiAdminApiUrl, nil)
	if err != nil {
		panic("Error while making request to project")
	}
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	apiProjectsResp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer apiProjectsResp.Body.Close()

	buf := &bytes.Buffer{}
	tee := io.TeeReader(apiProjectsResp.Body, buf)
	reader1, _ := ioutil.ReadAll(tee)
	reader2, _ := ioutil.ReadAll(buf)

	fmt.Printf("Response %v \n", buf.String())

	var responseSuccess GAEListInstancesResponse
	err = json.Unmarshal(reader1, &responseSuccess)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Services %v\n", responseSuccess)
	var responseError ErrorAdminAPI
	err = json.Unmarshal(reader2, &responseError)
	if err != nil {
		panic(err)
	}
	return &responseSuccess, &responseError
}

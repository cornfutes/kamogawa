package gcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"kamogawa/cache/gaecache"
	"kamogawa/config"
	"kamogawa/types"
	"kamogawa/types/gcp/gaetypes"
	"log"
	"net/http"

	"gorm.io/gorm"
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

//
// Response {
// 	  "error": {
// 	    "code": 404,
// 	    "message": "App does not exist.",
// 	    "status": "NOT_FOUND"
// 	   }
// 	}

// https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps.services/list
func GAEListServices(db *gorm.DB, user types.User, projectId string, useCache bool) (*gaetypes.GAEListServicesResponse, *gaetypes.ErrorAdminAPI) {
	if config.CacheEnabled && useCache {
		responseSuccess, err := gaecache.ReadServicesCache(db, projectId)
		if err == nil {
			return responseSuccess, &gaetypes.ErrorAdminAPI{}
		}
	}

	responseSuccess, responseError := GAEListServicesMain(user, projectId)
	if responseSuccess == nil {
		return nil, responseError
	}

	if config.CacheEnabled {
		gaecache.WriteServicesCache(db, projectId, responseSuccess)
	}

	return responseSuccess, responseError
}

func GAEListServicesMain(user types.User, projectId string) (*gaetypes.GAEListServicesResponse, *gaetypes.ErrorAdminAPI) {
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

	var responseSuccess gaetypes.GAEListServicesResponse
	err = json.Unmarshal(reader1, &responseSuccess)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Services %v\n", responseSuccess)
	var responseError gaetypes.ErrorAdminAPI
	err = json.Unmarshal(reader2, &responseError)
	if err != nil {
		panic(err)
	}
	return &responseSuccess, &responseError
}

// https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps.services.versions/list
func GAEListVersions(db *gorm.DB, user types.User, projectId string, serviceName string, useCache bool) (*gaetypes.GAEListVersionsResponse, *gaetypes.ErrorAdminAPI) {
	if config.CacheEnabled && useCache {
		responseSuccess, err := gaecache.ReadVersionsCache(db, projectId, serviceName)
		if err == nil {
			return responseSuccess, &gaetypes.ErrorAdminAPI{}
		}
	}

	responseSuccess, responseError := GAEListVersionsMain(user, projectId, serviceName)
	if responseSuccess == nil {
		return nil, responseError
	}

	if config.CacheEnabled {
		gaecache.WriteVersionsCache(db, projectId, serviceName, responseSuccess)
	}

	return responseSuccess, responseError
}

func GAEListVersionsMain(user types.User, projectId string, serviceName string) (*gaetypes.GAEListVersionsResponse, *gaetypes.ErrorAdminAPI) {
	apiAdminApiUrl := "https://appengine.googleapis.com/v1/apps/" + projectId + "/services/" + serviceName + "/versions/"
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

	var responseSuccess gaetypes.GAEListVersionsResponse
	err = json.Unmarshal(reader1, &responseSuccess)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Services %v\n", responseSuccess)
	var responseError gaetypes.ErrorAdminAPI
	err = json.Unmarshal(reader2, &responseError)
	if err != nil {
		panic(err)
	}
	return &responseSuccess, &responseError
}

// https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps.services.versions.instances/list
func GAEListInstances(db *gorm.DB, user types.User, projectId string, serviceName string, versionName string, useCache bool) (*gaetypes.GAEListInstancesResponse, *gaetypes.ErrorAdminAPI) {
	if config.CacheEnabled && useCache {
		responseSuccess, err := gaecache.ReadInstancesCache(db, projectId, serviceName, versionName)
		if err == nil {
			return responseSuccess, &gaetypes.ErrorAdminAPI{}
		}
	}

	responseSuccess, responseError := GAEListInstancesMain(user, projectId, serviceName, versionName)
	if responseSuccess == nil {
		return nil, responseError
	}

	if config.CacheEnabled {
		gaecache.WriteInstancesCache(db, projectId, serviceName, versionName, responseSuccess)
	}

	return responseSuccess, responseError
}

func GAEListInstancesMain(user types.User, projectId string, serviceName string, versionName string) (*gaetypes.GAEListInstancesResponse, *gaetypes.ErrorAdminAPI) {
	apiAdminApiUrl := "https://appengine.googleapis.com/v1/apps/" + projectId + "/services/" + serviceName + "/versions/" + versionName + "/instances"
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

	var responseSuccess gaetypes.GAEListInstancesResponse
	err = json.Unmarshal(reader1, &responseSuccess)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Services %v\n", responseSuccess)
	var responseError gaetypes.ErrorAdminAPI
	err = json.Unmarshal(reader2, &responseError)
	if err != nil {
		panic(err)
	}
	return &responseSuccess, &responseError
}

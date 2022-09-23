package gcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"kamogawa/cache/gcpcache"
	"kamogawa/cache/gcpcache/gcecache"
	"kamogawa/config"
	"kamogawa/types"
	"kamogawa/types/gcp/coretypes"
	"kamogawa/types/gcp/gcetypes"
	"log"
	"net/http"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//   "error": {
//     "code": 401,
//     "message": "Request had invalid authentication credentials. Expected OAuth 2 access token, login cookie or other valid authentication credential. See https://developers.google.com/identity/sign-in/web/devconsole-project.",
//     "status": "UNAUTHENTICATED",
//     "details": [
//       {
//         "@type": "type.googleapis.com/google.rpc.ErrorInfo",
//         "reason": "ACCESS_TOKEN_TYPE_UNSUPPORTED",
//         "metadata": {
//           "service": "cloudresourcemanager.googleapis.com",
//           "method": "google.cloudresourcemanager.v1.Projects.ListProjects"
//         }
//       }
//     ]
//   }
// }

func GCPListProjects(db *gorm.DB, user types.User, useCache bool) ([]coretypes.ProjectDB, *gcetypes.ErrorResponse) {
	if config.CacheEnabled && useCache {
		projectDBs := gcecache.ReadProjectsCache2(db, user)
		return projectDBs, nil
	}

	responseSuccess, responseError := GCPListProjectsMain(db, user)
	if responseSuccess == nil {
		return []coretypes.ProjectDB{}, responseError
	}

	projectDBs := make([]coretypes.ProjectDB, 0, len(responseSuccess.Projects))
	for _, v := range responseSuccess.Projects {
		projectDBs = append(projectDBs, coretypes.ProjectToProjectDB(&v))
	}

	sort.Slice(projectDBs, func(i, j int) bool {
		if projectDBs[i].Name != projectDBs[j].Name {
			return projectDBs[i].Name < projectDBs[j].Name
		}
		return projectDBs[i].ProjectId < projectDBs[j].ProjectId
	})

	if config.CacheEnabled {
		gcecache.WriteProjectsCache2(db, user, projectDBs)
	}

	return projectDBs, responseError
}

func GCPListProjectsMain(db *gorm.DB, user types.User) (*gcetypes.ListProjectResponse, *gcetypes.ErrorResponse) {
	apiProjectsUrl := "https://cloudresourcemanager.googleapis.com/v1/projects?filter=lifecycleState:ACTIVE"
	log.Printf("User %v\n", user.AccessToken)
	if !user.AccessToken.Valid {
		panic("Access Token expected but not found %v\n")
	}
	var bearer = "Bearer " + user.AccessToken.String

	fmt.Printf("Token %v\n", bearer)

	req, err := http.NewRequest("GET", apiProjectsUrl, nil)
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

	var responseSuccess gcetypes.ListProjectResponse
	err = json.Unmarshal(reader1, &responseSuccess)
	if err != nil {
		panic(err)
	}
	var responseError gcetypes.ErrorResponse
	err = json.Unmarshal(reader2, &responseError)
	if err != nil {
		panic(err)
	}
	if responseError.Error.Code == 401 {
		if strings.HasPrefix(responseError.Error.Message, "Request had invalid authentication credentials.") {
			return nil, &responseError
		}
	}
	fmt.Printf("response raw %v \n", string(reader1))
	fmt.Printf("Fetched %v projects\n", len(responseSuccess.Projects))
	return &responseSuccess, &responseError
}

func GCPListProjectAPIs(db *gorm.DB, user types.User, projectDB coretypes.ProjectDB, useCache bool) ([]coretypes.GCPProjectAPI, *gcetypes.ErrorResponse) {
	if config.CacheEnabled && useCache {
		return gcpcache.ReadGCPProjectAPIsCache(db, user, projectDB), nil
	}

	responseSuccess, responseError := gcpListProjectAPIsMain(user, projectDB)
	if responseError != nil {
		return []coretypes.GCPProjectAPI{}, responseError
	}

	gcpProjectAPI := coretypes.GCPProjectAPI{ProjectId: projectDB.ProjectId}
	gcpProjectAPI.API.Scan(responseSuccess)
	gcpProjectAPI.IsGAEEnabled = gcpIsGAEEnabled(user, projectDB)
	gcpProjectAPIs := []coretypes.GCPProjectAPI{gcpProjectAPI}

	if config.CacheEnabled {
		gcpcache.WriteGCPProjectAPIsCache(db, gcpProjectAPIs)
	}

	return gcpProjectAPIs, responseError
}

func gcpListProjectAPIsMain(user types.User, projectDB coretypes.ProjectDB) (string, *gcetypes.ErrorResponse) {
	url := "https://serviceusage.googleapis.com/v1/projects/" + projectDB.ProjectId + "/services?filter=state:ENABLED"
	log.Printf("User %v\n", user.AccessToken)
	if !user.AccessToken.Valid {
		panic("Access Token expected but not found %v\n")
	}
	var bearer = "Bearer " + user.AccessToken.String

	fmt.Printf("Token %v\n", bearer)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Error while making request to project")
	}
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bytes, _ := io.ReadAll(resp.Body)
		return string(bytes), nil
	} else {
		var responseError gcetypes.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&responseError); err != nil {
			panic(err)
		}
		return "", &responseError
	}
}

// gcpIsGAEEnabled GAE is special, need to check if it's initialized not just if API enabled
func gcpIsGAEEnabled(user types.User, projectDB coretypes.ProjectDB) bool {
	url := "https://appengine.googleapis.com/v1/apps/" + projectDB.ProjectId
	log.Printf("User %v\n", user.AccessToken)
	if !user.AccessToken.Valid {
		panic("Access Token expected but not found %v\n")
	}
	var bearer = "Bearer " + user.AccessToken.String

	fmt.Printf("Token %v\n", bearer)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Error while making request to project")
	}
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		var responseError gcetypes.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&responseError); err != nil {
			panic(err)
		}
		return !strings.HasPrefix(responseError.Error.Message, "App does not exist.")
	}
	return true
}

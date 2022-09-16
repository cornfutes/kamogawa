package gcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"kamogawa/cache/gcecache"
	"kamogawa/config"
	"kamogawa/types"
	"kamogawa/types/gcp/gcetypes"
	"log"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs"
	"gorm.io/gorm"
)

var ScopeGCE = "https://www.googleapis.com/auth/compute.readonly"

func GCEListInstances(db *gorm.DB, user types.User, projectId string, useCache bool) (*gcetypes.GCEAggregatedInstances, *gcetypes.ErrorGCEListInstance) {
	if config.CacheEnabled && useCache {
		responseSuccess, err := gcecache.ReadInstancesCache(db, projectId)
		if err == nil {
			return responseSuccess, &gcetypes.ErrorGCEListInstance{}
		}
	}

	responseSuccess, responseError := GCEListInstancesMain(user, projectId)
	if responseSuccess == nil {
		return nil, responseError
	}

	if config.CacheEnabled {
		gcecache.WriteInstancesCache(db, user, projectId, responseSuccess)
	}

	return responseSuccess, responseError
}

func GCEListInstancesMain(user types.User, projectId string) (*gcetypes.GCEAggregatedInstances, *gcetypes.ErrorGCEListInstance) {
	url := "https://compute.googleapis.com/compute/v1/projects/" + projectId + "/aggregated/instances"
	var bearer = "Bearer " + user.AccessToken.String
	fmt.Printf("Token %v\n", bearer)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Error while making request to API")
	}
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}

	tee := io.TeeReader(resp.Body, buf)

	reader1, _ := ioutil.ReadAll(tee)
	reader2, _ := ioutil.ReadAll(buf)

	jsonParsed, err := gabs.ParseJSON(reader1)
	if err != nil {
		panic(err)
	}

	result := gcetypes.GCEAggregatedInstances{}
	// Iterating address objects
	a, _ := jsonParsed.S("items").ChildrenMap()
	for key := range a {
		if jsonParsed.ExistsP("items." + key + ".warning") {
		} else {
			zone := gcetypes.ZoneMetadata{}
			zone.Zone = strings.Split(key, "/")[1]
			json.Unmarshal(jsonParsed.Search("items", key, "instances").Bytes(), &zone.Instances)
			result.Zones = append(result.Zones, zone)
		}
	}

	var responseError gcetypes.ErrorGCEListInstance
	err = json.Unmarshal(reader2, &responseError)
	if err != nil {
		panic(err)
	}

	return &result, &responseError
}

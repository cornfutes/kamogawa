package main

import (
	"encoding/csv"
	"fmt"
	"github.com/tjarratt/babble"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	babbler := babble.NewBabbler()
	babbler.Count = 1

	iCsvFile, err := os.Create("test/instances.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	iCsvWriter := csv.NewWriter(iCsvFile)

	pCsvFile, err := os.Create("test/projects.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	pCsvWriter := csv.NewWriter(pCsvFile)

	sCsvFile, err := os.Create("test/gaeservices.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	sCsvWriter := csv.NewWriter(sCsvFile)

	vCsvFile, err := os.Create("test/gaeversions.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	vCsvWriter := csv.NewWriter(vCsvFile)

	gaeICsvFile, err := os.Create("test/gaeinstances.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	gaeICsvWriter := csv.NewWriter(gaeICsvFile)

	// int, str, str, str, str
	iCsvWriter.Write([]string{"id", "created_at", "updated_at", "deleted_at", "name", "status", "project_id", "zone"})

	// int, int, str, str, str, 1337gamer@gmail.com
	pCsvWriter.Write([]string{"created_at", "updated_at", "deleted_at", "project_number", "project_id", "life_cycle_state", "name", "create_time", "email"})
	//
	//// str...
	//sCsvWriter.Write([]string{"name", "id", "project_id"})
	//
	//// str...
	//vCsvWriter.Write([]string{"name", "id", "serving_status", "service_id", "parent_id"})
	//
	//// str
	//iCsvWriter.Write([]string{"name", "id", "vm_name", "version_id"})

	instanceId := 0
	for numProjects := 0; numProjects < 10; numProjects++ {
		project_number := fmt.Sprintf("%v-%v-%v", babbler.Babble(), babbler.Babble(), strconv.Itoa(rand.Intn(100)))
		project_id := babbler.Babble()
		name := babbler.Babble()
		life_cycle_state := "ACTIVE"
		email := "1337gamer@gmail.com"
		pCsvWriter.Write([]string{"", "", "", project_number, project_id, life_cycle_state, name, "", email})

		for numZone := 0; numZone < 3; numZone++ {
			zones := []string{"zones/us-west1-a", "zones/us-west1-b", "zones/us-west1-c"}
			zone := zones[numZone]
			for numInstances := 0; numInstances < 10; numInstances++ {
				instanceId++
				id := strconv.Itoa(instanceId)
				name := babbler.Babble()
				status := "ACTIVE"
				iCsvWriter.Write([]string{id, "", "", "", name, status, project_id, zone})
			}
		}

		for numServices := 0; numServices < 10; numServices++ {
			serviceName := babbler.Babble()
			serviceId := fmt.Sprintf("apps/%v/services/%v", project_id, serviceName)
			sCsvWriter.Write([]string{serviceName, serviceId, project_id})

			for numVersions := 0; numVersions < 2; numVersions++ {
				versionId := babbler.Babble()
				serving_status := "SERVING"
				versionName := fmt.Sprintf("%v/versions/%v", serviceId, versionId)
				vCsvWriter.Write([]string{versionName, versionId, serving_status, serviceId})

				for numInstances := 0; numInstances < 10; numInstances++ {
					instanceId++
					name := babbler.Babble()
					id := strconv.Itoa(instanceId)
					vm_name := babbler.Babble()
					gaeICsvWriter.Write([]string{name, id, vm_name, versionName})
				}
			}
		}
	}

	iCsvWriter.Flush()
	iCsvFile.Close()
	pCsvWriter.Flush()
	pCsvFile.Close()
	sCsvWriter.Flush()
	sCsvFile.Close()
	vCsvWriter.Flush()
	vCsvFile.Close()
	gaeICsvWriter.Flush()
	gaeICsvFile.Close()
}

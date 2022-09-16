package gaetypes

type GAEServiceTrafficAllocation struct {
}

type GAEServiceTrafficAllocations struct {
	Allocations GAEServiceTrafficAllocation `json:"allocations"`
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
	Id     string `json:"id`
	VMName string `json:"vmName"`
}

type GAEListInstancesResponse struct {
	Instances     []GAEInstance `json:"instances"`
	NextPageToken string        `json:"nextPageToken"`
}

type ErrorAdminAPI struct {
	Error ErrorAdminAPIError `json:"error"`
}
type ErrorAdminAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

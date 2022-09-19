package coretypes

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type GCPServiceUsageServicesListResponse struct {
	Services []GCPServiceUsageServicesListResponseService `json:"services"`
}

type GCPServiceUsageServicesListResponseService struct {
	Config GCPServiceUsageServicesListResponseConfig `json:"config"`
}

type GCPServiceUsageServicesListResponseConfig struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func JSONToGCPServiceUsageServicesListResponse(j datatypes.JSON) GCPServiceUsageServicesListResponse {
	var resp GCPServiceUsageServicesListResponse
	if err := json.Unmarshal(j, &resp); err != nil {
		panic(err)
	}
	return resp
}

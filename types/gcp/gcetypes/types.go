package gcetypes

//type ProjectParent struct {
//	Type string `json:"type"`
//	Id   string `json:"id"`
//}

type Project struct {
	ProjectNumber  string `json:"projectNumber"`
	ProjectId      string `json:"projectId"`
	LifeCycleState string `json:"lifecycleState"`
	Name           string `json:"name"`
	CreateTime     string `json:"createTime"`
	//Parent         ProjectParent
}

type ListProjectResponse struct {
	Projects []Project `json:"projects"`
}

type ErrorResponse struct {
	Error ErrorResponseError `json:"error"`
}

type errorResponseErrorDetailMetadata struct {
	Service string `json:"service"`
	Method  string `json:"method"`
}
type errorResponseErrorDetail struct {
	Type     string                           `json:"@type"`
	Reason   string                           `json:"reason"`
	Metadata errorResponseErrorDetailMetadata `json:"metadata"`
}
type ErrorResponseError struct {
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Status  string                     `json:"status"`
	Details []errorResponseErrorDetail `json:"details"`
}

type ErrorGCEListInstanceErrorError struct {
	Message string `json:"message"`
	Domain  string `json:"domain"`
	Reason  string `json:"reason"`
}

type ErrorGCEListInstanceError struct {
	Code    int                              `json:"code"`
	Message string                           `json:"message"`
	Errors  []ErrorGCEListInstanceErrorError `json:"errors"`
}

type ErrorGCEListInstance struct {
	Error ErrorGCEListInstanceError `json:"error"`
}

type Warning struct {
	Code string `json:"code"`
}
type ZoneWarning struct {
	Warning Warning `json:"warning"`
}

type GCEAggregatedInstances struct {
	Zones []ZoneMetadata
}

type ZoneMetadata struct {
	Zone      string
	Instances []GCEInstance
}

type GCEInstance struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

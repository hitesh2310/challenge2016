package models

type Distributor struct {
	Id                string   `json:"id,omitempty"`
	IncludeCode       []string `json:"includeCode,omitempty"`
	ExcludeCode       []string `json:"excludeCode,omitempty"`
	HeadDistributorId string   `json:"headDistributor,omitempty"`
	SubDistributorId  string   `json:"subDistributor,omitempty"`
}

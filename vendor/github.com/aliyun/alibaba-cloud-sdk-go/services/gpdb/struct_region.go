package gpdb

// Region is a nested struct in gpdb response
type Region struct {
	RegionId string `json:"RegionId" xml:"RegionId"`
	Zones    Zones  `json:"Zones" xml:"Zones"`
}

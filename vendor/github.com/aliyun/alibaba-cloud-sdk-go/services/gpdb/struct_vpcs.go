package gpdb

// Vpcs is a nested struct in gpdb response
type Vpcs struct {
	Vpc []VpcItem `json:"Vpc" xml:"Vpc"`
}

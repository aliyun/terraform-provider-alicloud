package gpdb

// Vpc is a nested struct in gpdb response
type Vpc struct {
	VpcItem []VpcItem `json:"VpcItem" xml:"VpcItem"`
}

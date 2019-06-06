package gpdb

// Tag is a nested struct in gpdb response
type Tag struct {
	Key   string `json:"Key" xml:"Key"`
	Value string `json:"Value" xml:"Value"`
}

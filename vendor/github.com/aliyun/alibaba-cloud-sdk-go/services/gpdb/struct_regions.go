package gpdb

// Regions is a nested struct in gpdb response
type Regions struct {
	Region []Region `json:"Region" xml:"Region"`
}

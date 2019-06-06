package gpdb

// Zones is a nested struct in gpdb response
type Zones struct {
	Zone []Zone `json:"Zone" xml:"Zone"`
}

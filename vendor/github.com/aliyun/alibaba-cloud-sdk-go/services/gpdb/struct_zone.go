package gpdb

// Zone is a nested struct in gpdb response
type Zone struct {
	ZoneId     string `json:"ZoneId" xml:"ZoneId"`
	VpcEnabled bool   `json:"VpcEnabled" xml:"VpcEnabled"`
}

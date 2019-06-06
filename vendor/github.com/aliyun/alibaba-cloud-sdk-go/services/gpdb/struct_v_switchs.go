package gpdb

// VSwitchs is a nested struct in gpdb response
type VSwitchs struct {
	VSwitch []VSwitch `json:"VSwitch" xml:"VSwitch"`
}

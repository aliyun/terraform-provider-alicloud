package gpdb

// VSwitches is a nested struct in gpdb response
type VSwitches struct {
	VSwitch []VSwitchItem `json:"VSwitch" xml:"VSwitch"`
}

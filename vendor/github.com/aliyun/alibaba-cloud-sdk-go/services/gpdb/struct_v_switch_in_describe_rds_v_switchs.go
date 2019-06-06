package gpdb

// VSwitchInDescribeRdsVSwitchs is a nested struct in gpdb response
type VSwitchInDescribeRdsVSwitchs struct {
	VSwitchItem []VSwitchItem `json:"VSwitchItem" xml:"VSwitchItem"`
}

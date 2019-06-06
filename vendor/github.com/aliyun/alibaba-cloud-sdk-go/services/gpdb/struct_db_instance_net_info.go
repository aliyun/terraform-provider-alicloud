package gpdb

// DBInstanceNetInfo is a nested struct in gpdb response
type DBInstanceNetInfo struct {
	ConnectionString string `json:"ConnectionString" xml:"ConnectionString"`
	IPAddress        string `json:"IPAddress" xml:"IPAddress"`
	IPType           string `json:"IPType" xml:"IPType"`
	Port             string `json:"Port" xml:"Port"`
	VPCId            string `json:"VPCId" xml:"VPCId"`
	VSwitchId        string `json:"VSwitchId" xml:"VSwitchId"`
}

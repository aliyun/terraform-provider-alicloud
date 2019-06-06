package gpdb

// LogFile is a nested struct in gpdb response
type LogFile struct {
	FileID         string `json:"FileID" xml:"FileID"`
	LogStatus      string `json:"LogStatus" xml:"LogStatus"`
	LogDownloadURL string `json:"LogDownloadURL" xml:"LogDownloadURL"`
	LogSize        string `json:"LogSize" xml:"LogSize"`
	LogStartTime   string `json:"LogStartTime" xml:"LogStartTime"`
	LogEndTime     string `json:"LogEndTime" xml:"LogEndTime"`
}

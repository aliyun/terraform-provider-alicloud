package gpdb

// SQLRecord is a nested struct in gpdb response
type SQLRecord struct {
	DBName              string `json:"DBName" xml:"DBName"`
	AccountName         string `json:"AccountName" xml:"AccountName"`
	HostAddress         string `json:"HostAddress" xml:"HostAddress"`
	SQLText             string `json:"SQLText" xml:"SQLText"`
	TotalExecutionTimes int    `json:"TotalExecutionTimes" xml:"TotalExecutionTimes"`
	ReturnRowCounts     int    `json:"ReturnRowCounts" xml:"ReturnRowCounts"`
	ExecuteTime         string `json:"ExecuteTime" xml:"ExecuteTime"`
	ThreadID            string `json:"ThreadID" xml:"ThreadID"`
}

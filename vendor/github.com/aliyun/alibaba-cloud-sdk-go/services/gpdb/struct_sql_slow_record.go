package gpdb

// SQLSlowRecord is a nested struct in gpdb response
type SQLSlowRecord struct {
	HostAddress        string `json:"HostAddress" xml:"HostAddress"`
	DBName             string `json:"DBName" xml:"DBName"`
	SQLText            string `json:"SQLText" xml:"SQLText"`
	QueryTimes         int    `json:"QueryTimes" xml:"QueryTimes"`
	LockTimes          int    `json:"LockTimes" xml:"LockTimes"`
	ParseRowCounts     int    `json:"ParseRowCounts" xml:"ParseRowCounts"`
	ReturnRowCounts    int    `json:"ReturnRowCounts" xml:"ReturnRowCounts"`
	ExecutionStartTime string `json:"ExecutionStartTime" xml:"ExecutionStartTime"`
}

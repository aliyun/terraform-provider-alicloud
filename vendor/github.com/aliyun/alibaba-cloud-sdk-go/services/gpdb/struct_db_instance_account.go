package gpdb

// DBInstanceAccount is a nested struct in gpdb response
type DBInstanceAccount struct {
	DBInstanceId       string `json:"DBInstanceId" xml:"DBInstanceId"`
	AccountName        string `json:"AccountName" xml:"AccountName"`
	AccountStatus      string `json:"AccountStatus" xml:"AccountStatus"`
	AccountDescription string `json:"AccountDescription" xml:"AccountDescription"`
}

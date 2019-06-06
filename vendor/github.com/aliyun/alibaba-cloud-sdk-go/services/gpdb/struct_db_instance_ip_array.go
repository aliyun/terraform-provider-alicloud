package gpdb

// DBInstanceIPArray is a nested struct in gpdb response
type DBInstanceIPArray struct {
	DBInstanceIPArrayName      string `json:"DBInstanceIPArrayName" xml:"DBInstanceIPArrayName"`
	DBInstanceIPArrayAttribute string `json:"DBInstanceIPArrayAttribute" xml:"DBInstanceIPArrayAttribute"`
	SecurityIPList             string `json:"SecurityIPList" xml:"SecurityIPList"`
}

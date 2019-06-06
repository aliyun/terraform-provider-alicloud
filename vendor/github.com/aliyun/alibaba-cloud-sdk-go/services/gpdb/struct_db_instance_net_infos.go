package gpdb

// DBInstanceNetInfos is a nested struct in gpdb response
type DBInstanceNetInfos struct {
	DBInstanceNetInfo []DBInstanceNetInfo `json:"DBInstanceNetInfo" xml:"DBInstanceNetInfo"`
}

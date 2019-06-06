package gpdb

// Accounts is a nested struct in gpdb response
type Accounts struct {
	DBInstanceAccount []DBInstanceAccount `json:"DBInstanceAccount" xml:"DBInstanceAccount"`
}

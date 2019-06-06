package gpdb

// ItemsInDescribeDBInstanceAttribute is a nested struct in gpdb response
type ItemsInDescribeDBInstanceAttribute struct {
	DBInstanceAttribute []DBInstanceAttribute `json:"DBInstanceAttribute" xml:"DBInstanceAttribute"`
}

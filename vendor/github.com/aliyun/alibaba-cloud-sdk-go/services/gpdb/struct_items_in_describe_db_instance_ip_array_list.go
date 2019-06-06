package gpdb

// ItemsInDescribeDBInstanceIPArrayList is a nested struct in gpdb response
type ItemsInDescribeDBInstanceIPArrayList struct {
	DBInstanceIPArray []DBInstanceIPArray `json:"DBInstanceIPArray" xml:"DBInstanceIPArray"`
}

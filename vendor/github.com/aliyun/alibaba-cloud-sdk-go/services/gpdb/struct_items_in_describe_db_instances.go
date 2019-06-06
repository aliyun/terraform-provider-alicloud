package gpdb

// ItemsInDescribeDBInstances is a nested struct in gpdb response
type ItemsInDescribeDBInstances struct {
	DBInstance []DBInstance `json:"DBInstance" xml:"DBInstance"`
}

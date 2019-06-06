package gpdb

// ItemsInDescribeSlowLogRecords is a nested struct in gpdb response
type ItemsInDescribeSlowLogRecords struct {
	SQLSlowRecord []SQLSlowRecord `json:"SQLSlowRecord" xml:"SQLSlowRecord"`
}

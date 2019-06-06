package gpdb

// ItemsInDescribeSQLLogRecords is a nested struct in gpdb response
type ItemsInDescribeSQLLogRecords struct {
	SQLRecord []SQLRecord `json:"SQLRecord" xml:"SQLRecord"`
}

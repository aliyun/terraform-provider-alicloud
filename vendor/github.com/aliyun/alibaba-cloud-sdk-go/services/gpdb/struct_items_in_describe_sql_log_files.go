package gpdb

// ItemsInDescribeSQLLogFiles is a nested struct in gpdb response
type ItemsInDescribeSQLLogFiles struct {
	LogFile []LogFile `json:"LogFile" xml:"LogFile"`
}

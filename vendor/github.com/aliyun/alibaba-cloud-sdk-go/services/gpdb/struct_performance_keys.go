package gpdb

// PerformanceKeys is a nested struct in gpdb response
type PerformanceKeys struct {
	PerformanceKey []string `json:"PerformanceKey" xml:"PerformanceKey"`
}

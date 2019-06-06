package gpdb

// TagsInDescribeDBInstanceAttribute is a nested struct in gpdb response
type TagsInDescribeDBInstanceAttribute struct {
	Tag []Tag `json:"Tag" xml:"Tag"`
}

package gpdb

// TagsInDescribeDBInstances is a nested struct in gpdb response
type TagsInDescribeDBInstances struct {
	Tag []Tag `json:"Tag" xml:"Tag"`
}

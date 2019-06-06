package gpdb

// TagResources is a nested struct in gpdb response
type TagResources struct {
	TagResource []TagResource `json:"TagResource" xml:"TagResource"`
}

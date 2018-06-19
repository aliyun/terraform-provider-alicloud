package alicloud

type IndexFiledType string

const (
	TextType   = IndexFiledType("text")
	LongType   = IndexFiledType("long")
	DoubleType = IndexFiledType("double")
	JsonType   = IndexFiledType("json")
)

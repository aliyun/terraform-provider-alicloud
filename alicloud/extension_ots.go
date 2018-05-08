package alicloud

type PrimaryKeyTypeString string

const (
	IntegerType = PrimaryKeyTypeString("Integer")
	StringType  = PrimaryKeyTypeString("String")
	BinaryType  = PrimaryKeyTypeString("Binary")
)

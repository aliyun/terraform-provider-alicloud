package alicloud

type InstanceType string
type DRDSInstancePayType string

const (
	PrivateType  = InstanceType("PRIVATE")
	PublicType   = InstanceType("PUBLIC")
	PrivateType_ = InstanceType("1")
	PublicType_  = InstanceType("0")
	DRDSInstancePostPayType = DRDSInstancePayType("drdsPost")

)
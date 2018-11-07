package alicloud

type InstanceType string
type DRDSInstancePayType string

const (
	PrivateType             = InstanceType("1")
	PublicType              = InstanceType("0")
	DRDSInstancePostPayType = DRDSInstancePayType("PostPaid")
	DRDSInstancePrePayType  = DRDSInstancePayType("PrePaid")
)

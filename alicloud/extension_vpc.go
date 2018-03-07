package alicloud

type NatGatewaySpec string

const (
	NatGatewaySmallSpec  = NatGatewaySpec("Small")
	NatGatewayMiddleSpec = NatGatewaySpec("Middle")
	NatGatewayLargeSpec  = NatGatewaySpec("Large")
)

const (
	EcsInstance = "EcsInstance"
	SlbInstance = "SlbInstance"
	Nat         = "Nat"
	HaVip       = "HaVip"
)

type RouterType string
type Role string
type Spec string

const (
	VRouter = RouterType("VRouter")
	VBR     = RouterType("VBR")

	InitiatingSide = Role("InitiatingSide")
	AcceptingSide  = Role("AcceptingSide")

	Small1  = Spec("Small.1")
	Small2  = Spec("Small.2")
	Small5  = Spec("Small.5")
	Middle1 = Spec("Middle.1")
	Middle2 = Spec("Middle.2")
	Middle5 = Spec("Middle.5")
	Large1  = Spec("Large.1")
	Large2  = Spec("Large.2")
)

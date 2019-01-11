package alicloud

type RecordType string

const (
	RecordA     = RecordType("A")
	RecordCNAME = RecordType("CNAME")
	RecordTXT   = RecordType("TXT")
	RecordMX    = RecordType("MX")
	RecordPTR   = RecordType("PTR")
)

var PvtzThrottlingUserCatcher = Catcher{PvtzThrottlingUser, 30, 2}

func PvtzInvoker() Invoker {
	i := Invoker{}
	i.AddCatcher(PvtzThrottlingUserCatcher)
	i.AddCatcher(ServiceBusyCatcher)
	return i
}

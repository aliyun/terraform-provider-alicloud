package alicloud

type RecordType string

const (
	RecordA     = RecordType("A")
	RecordCNAME = RecordType("CNAME")
	RecordTXT   = RecordType("TXT")
	RecordMX    = RecordType("MX")
	RecordPTR   = RecordType("PTR")
)

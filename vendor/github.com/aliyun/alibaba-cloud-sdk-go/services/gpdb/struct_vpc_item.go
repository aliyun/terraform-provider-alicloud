package gpdb

// VpcItem is a nested struct in gpdb response
type VpcItem struct {
	VpcId       string    `json:"VpcId" xml:"VpcId"`
	VpcName     string    `json:"VpcName" xml:"VpcName"`
	Bid         string    `json:"Bid" xml:"Bid"`
	AliUid      string    `json:"AliUid" xml:"AliUid"`
	RegionNo    string    `json:"RegionNo" xml:"RegionNo"`
	CidrBlock   string    `json:"CidrBlock" xml:"CidrBlock"`
	IsDefault   bool      `json:"IsDefault" xml:"IsDefault"`
	Status      string    `json:"Status" xml:"Status"`
	GmtCreate   string    `json:"GmtCreate" xml:"GmtCreate"`
	GmtModified string    `json:"GmtModified" xml:"GmtModified"`
	VSwitchs    []VSwitch `json:"VSwitchs" xml:"VSwitchs"`
}

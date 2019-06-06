package gpdb

// VSwitchItem is a nested struct in gpdb response
type VSwitchItem struct {
	VSwitchId   string `json:"VSwitchId" xml:"VSwitchId"`
	VSwitchName string `json:"VSwitchName" xml:"VSwitchName"`
	IzNo        string `json:"IzNo" xml:"IzNo"`
	Bid         string `json:"Bid" xml:"Bid"`
	AliUid      string `json:"AliUid" xml:"AliUid"`
	RegionNo    string `json:"RegionNo" xml:"RegionNo"`
	CidrBlock   string `json:"CidrBlock" xml:"CidrBlock"`
	IsDefault   bool   `json:"IsDefault" xml:"IsDefault"`
	Status      string `json:"Status" xml:"Status"`
	GmtCreate   string `json:"GmtCreate" xml:"GmtCreate"`
	GmtModified string `json:"GmtModified" xml:"GmtModified"`
}

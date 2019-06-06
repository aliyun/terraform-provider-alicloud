package gpdb

// VSwitch is a nested struct in gpdb response
type VSwitch struct {
	GmtCreate   string `json:"GmtCreate" xml:"GmtCreate"`
	IsDefault   bool   `json:"IsDefault" xml:"IsDefault"`
	VSwitchId   string `json:"VSwitchId" xml:"VSwitchId"`
	VSwitchName string `json:"VSwitchName" xml:"VSwitchName"`
	GmtModified string `json:"GmtModified" xml:"GmtModified"`
	Status      string `json:"Status" xml:"Status"`
	CidrBlock   string `json:"CidrBlock" xml:"CidrBlock"`
	IzNo        string `json:"IzNo" xml:"IzNo"`
}

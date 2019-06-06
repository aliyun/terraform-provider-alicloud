package gpdb

// DBInstanceAttribute is a nested struct in gpdb response
type DBInstanceAttribute struct {
	DBInstanceId          string                            `json:"DBInstanceId" xml:"DBInstanceId"`
	PayType               string                            `json:"PayType" xml:"PayType"`
	DBInstanceClassType   string                            `json:"DBInstanceClassType" xml:"DBInstanceClassType"`
	RegionId              string                            `json:"RegionId" xml:"RegionId"`
	ConnectionString      string                            `json:"ConnectionString" xml:"ConnectionString"`
	Port                  string                            `json:"Port" xml:"Port"`
	Engine                string                            `json:"Engine" xml:"Engine"`
	EngineVersion         string                            `json:"EngineVersion" xml:"EngineVersion"`
	DBInstanceClass       string                            `json:"DBInstanceClass" xml:"DBInstanceClass"`
	DBInstanceCpuCores    int                               `json:"DBInstanceCpuCores" xml:"DBInstanceCpuCores"`
	DBInstanceMemory      int                               `json:"DBInstanceMemory" xml:"DBInstanceMemory"`
	DBInstanceStorage     int                               `json:"DBInstanceStorage" xml:"DBInstanceStorage"`
	DBInstanceDiskMBPS    int                               `json:"DBInstanceDiskMBPS" xml:"DBInstanceDiskMBPS"`
	HostType              string                            `json:"HostType" xml:"HostType"`
	DBInstanceGroupCount  string                            `json:"DBInstanceGroupCount" xml:"DBInstanceGroupCount"`
	DBInstanceNetType     string                            `json:"DBInstanceNetType" xml:"DBInstanceNetType"`
	DBInstanceStatus      string                            `json:"DBInstanceStatus" xml:"DBInstanceStatus"`
	DBInstanceDescription string                            `json:"DBInstanceDescription" xml:"DBInstanceDescription"`
	LockMode              string                            `json:"LockMode" xml:"LockMode"`
	LockReason            string                            `json:"LockReason" xml:"LockReason"`
	ReadDelayTime         string                            `json:"ReadDelayTime" xml:"ReadDelayTime"`
	CreationTime          string                            `json:"CreationTime" xml:"CreationTime"`
	ExpireTime            string                            `json:"ExpireTime" xml:"ExpireTime"`
	MaintainStartTime     string                            `json:"MaintainStartTime" xml:"MaintainStartTime"`
	MaintainEndTime       string                            `json:"MaintainEndTime" xml:"MaintainEndTime"`
	AvailabilityValue     string                            `json:"AvailabilityValue" xml:"AvailabilityValue"`
	MaxConnections        int                               `json:"MaxConnections" xml:"MaxConnections"`
	SecurityIPList        string                            `json:"SecurityIPList" xml:"SecurityIPList"`
	ZoneId                string                            `json:"ZoneId" xml:"ZoneId"`
	InstanceNetworkType   string                            `json:"InstanceNetworkType" xml:"InstanceNetworkType"`
	VpcId                 string                            `json:"VpcId" xml:"VpcId"`
	ConnectionMode        string                            `json:"ConnectionMode" xml:"ConnectionMode"`
	Tags                  TagsInDescribeDBInstanceAttribute `json:"Tags" xml:"Tags"`
}

package gpdb

// DBInstance is a nested struct in gpdb response
type DBInstance struct {
	DBInstanceId          string                    `json:"DBInstanceId" xml:"DBInstanceId"`
	DBInstanceDescription string                    `json:"DBInstanceDescription" xml:"DBInstanceDescription"`
	PayType               string                    `json:"PayType" xml:"PayType"`
	InstanceNetworkType   string                    `json:"InstanceNetworkType" xml:"InstanceNetworkType"`
	ConnectionMode        string                    `json:"ConnectionMode" xml:"ConnectionMode"`
	RegionId              string                    `json:"RegionId" xml:"RegionId"`
	ZoneId                string                    `json:"ZoneId" xml:"ZoneId"`
	ExpireTime            string                    `json:"ExpireTime" xml:"ExpireTime"`
	DBInstanceStatus      string                    `json:"DBInstanceStatus" xml:"DBInstanceStatus"`
	Engine                string                    `json:"Engine" xml:"Engine"`
	EngineVersion         string                    `json:"EngineVersion" xml:"EngineVersion"`
	DBInstanceNetType     string                    `json:"DBInstanceNetType" xml:"DBInstanceNetType"`
	LockMode              string                    `json:"LockMode" xml:"LockMode"`
	LockReason            string                    `json:"LockReason" xml:"LockReason"`
	CreateTime            string                    `json:"CreateTime" xml:"CreateTime"`
	VpcId                 string                    `json:"VpcId" xml:"VpcId"`
	VSwitchId             string                    `json:"VSwitchId" xml:"VSwitchId"`
	Tags                  TagsInDescribeDBInstances `json:"Tags" xml:"Tags"`
}

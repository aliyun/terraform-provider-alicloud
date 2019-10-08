package emr

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// Group is a nested struct in emr response
type Group struct {
	Name               string                                `json:"Name" xml:"Name"`
	DataDiskSize       int                                   `json:"DataDiskSize" xml:"DataDiskSize"`
	BizId              string                                `json:"BizId" xml:"BizId"`
	GmtCreate          string                                `json:"GmtCreate" xml:"GmtCreate"`
	MinSize            int                                   `json:"MinSize" xml:"MinSize"`
	Description        string                                `json:"Description" xml:"Description"`
	CpuCount           int                                   `json:"CpuCount" xml:"CpuCount"`
	SysDiskCategory    string                                `json:"SysDiskCategory" xml:"SysDiskCategory"`
	PayType            string                                `json:"PayType" xml:"PayType"`
	GmtModified        string                                `json:"GmtModified" xml:"GmtModified"`
	Id                 int64                                 `json:"Id" xml:"Id"`
	DefaultCooldown    int                                   `json:"DefaultCooldown" xml:"DefaultCooldown"`
	SysDiskSize        int                                   `json:"SysDiskSize" xml:"SysDiskSize"`
	ActiveRuleCategory string                                `json:"ActiveRuleCategory" xml:"ActiveRuleCategory"`
	MemSize            int                                   `json:"MemSize" xml:"MemSize"`
	Status             string                                `json:"Status" xml:"Status"`
	SpotStrategy       string                                `json:"SpotStrategy" xml:"SpotStrategy"`
	ScalingGroupId     string                                `json:"ScalingGroupId" xml:"ScalingGroupId"`
	DataDiskCategory   string                                `json:"DataDiskCategory" xml:"DataDiskCategory"`
	HostGroupId        string                                `json:"HostGroupId" xml:"HostGroupId"`
	DataDiskCount      int                                   `json:"DataDiskCount" xml:"DataDiskCount"`
	MaxSize            int                                   `json:"MaxSize" xml:"MaxSize"`
	InstanceTypeList   []string                              `json:"InstanceTypeList" xml:"InstanceTypeList"`
	ScalingConfig      ScalingConfig                         `json:"ScalingConfig" xml:"ScalingConfig"`
	ScalingRuleList    ScalingRuleList                       `json:"ScalingRuleList" xml:"ScalingRuleList"`
	UserList           []UserInfo                            `json:"UserList" xml:"UserList"`
	SpotPriceLimits    SpotPriceLimitsInListScalingTaskGroup `json:"SpotPriceLimits" xml:"SpotPriceLimits"`
}

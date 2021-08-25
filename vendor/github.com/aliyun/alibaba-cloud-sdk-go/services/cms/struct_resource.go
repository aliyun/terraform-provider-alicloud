package cms

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

// Resource is a nested struct in cms response
type Resource struct {
	TemplateId                string                                      `json:"TemplateId" xml:"TemplateId"`
	Name                      string                                      `json:"Name" xml:"Name"`
	Category                  string                                      `json:"Category" xml:"Category"`
	Unit                      string                                      `json:"Unit" xml:"Unit"`
	SystemEventTemplates      string                                      `json:"SystemEventTemplates" xml:"SystemEventTemplates"`
	Desc                      string                                      `json:"Desc" xml:"Desc"`
	BindUrl                   string                                      `json:"BindUrl" xml:"BindUrl"`
	Dimensions                string                                      `json:"Dimensions" xml:"Dimensions"`
	GroupName                 string                                      `json:"GroupName" xml:"GroupName"`
	ComparisonOperator        string                                      `json:"ComparisonOperator" xml:"ComparisonOperator"`
	ServiceId                 string                                      `json:"ServiceId" xml:"ServiceId"`
	RestVersion               string                                      `json:"RestVersion" xml:"RestVersion"`
	RegionId                  string                                      `json:"RegionId" xml:"RegionId"`
	Tag                       string                                      `json:"Tag" xml:"Tag"`
	InstanceId                string                                      `json:"InstanceId" xml:"InstanceId"`
	DynamicTagRuleId          string                                      `json:"DynamicTagRuleId" xml:"DynamicTagRuleId"`
	Expression                string                                      `json:"Expression" xml:"Expression"`
	NetworkType               string                                      `json:"NetworkType" xml:"NetworkType"`
	Description               string                                      `json:"Description" xml:"Description"`
	Periods                   string                                      `json:"Periods" xml:"Periods"`
	Type                      string                                      `json:"Type" xml:"Type"`
	Times                     int                                         `json:"Times" xml:"Times"`
	Level                     int                                         `json:"Level" xml:"Level"`
	InstanceName              string                                      `json:"InstanceName" xml:"InstanceName"`
	Dimension                 string                                      `json:"Dimension" xml:"Dimension"`
	PreCondition              string                                      `json:"PreCondition" xml:"PreCondition"`
	HostAvailabilityTemplates string                                      `json:"HostAvailabilityTemplates" xml:"HostAvailabilityTemplates"`
	Id                        int64                                       `json:"Id" xml:"Id"`
	Threshold                 string                                      `json:"Threshold" xml:"Threshold"`
	GroupFounderTagValue      string                                      `json:"GroupFounderTagValue" xml:"GroupFounderTagValue"`
	GmtCreate                 int64                                       `json:"GmtCreate" xml:"GmtCreate"`
	Namespace                 string                                      `json:"Namespace" xml:"Namespace"`
	GroupId                   int64                                       `json:"GroupId" xml:"GroupId"`
	GmtModified               int64                                       `json:"GmtModified" xml:"GmtModified"`
	GroupFounderTagKey        string                                      `json:"GroupFounderTagKey" xml:"GroupFounderTagKey"`
	MetricName                string                                      `json:"MetricName" xml:"MetricName"`
	Labels                    string                                      `json:"Labels" xml:"Labels"`
	ProcessMonitorTemplates   string                                      `json:"ProcessMonitorTemplates" xml:"ProcessMonitorTemplates"`
	Statistics                string                                      `json:"Statistics" xml:"Statistics"`
	TemplateIds               TemplateIds                                 `json:"TemplateIds" xml:"TemplateIds"`
	Vpc                       Vpc                                         `json:"Vpc" xml:"Vpc"`
	Region                    Region                                      `json:"Region" xml:"Region"`
	AlertResults              []Result                                    `json:"AlertResults" xml:"AlertResults"`
	Tags                      TagsInDescribeMonitorGroupInstanceAttribute `json:"Tags" xml:"Tags"`
	ContactGroups             ContactGroupsInDescribeMonitorGroups        `json:"ContactGroups" xml:"ContactGroups"`
	AlertTemplates            AlertTemplates                              `json:"AlertTemplates" xml:"AlertTemplates"`
}

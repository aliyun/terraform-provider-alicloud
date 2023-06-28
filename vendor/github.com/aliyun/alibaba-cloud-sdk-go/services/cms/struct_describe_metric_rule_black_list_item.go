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

// DescribeMetricRuleBlackListItem is a nested struct in cms response
type DescribeMetricRuleBlackListItem struct {
	EffectiveTime   string        `json:"EffectiveTime" xml:"EffectiveTime"`
	UpdateTime      string        `json:"UpdateTime" xml:"UpdateTime"`
	CreateTime      string        `json:"CreateTime" xml:"CreateTime"`
	ScopeType       string        `json:"ScopeType" xml:"ScopeType"`
	IsEnable        bool          `json:"IsEnable" xml:"IsEnable"`
	Namespace       string        `json:"Namespace" xml:"Namespace"`
	Category        string        `json:"Category" xml:"Category"`
	EnableEndTime   int64         `json:"EnableEndTime" xml:"EnableEndTime"`
	Name            string        `json:"Name" xml:"Name"`
	EnableStartTime int64         `json:"EnableStartTime" xml:"EnableStartTime"`
	Id              string        `json:"Id" xml:"Id"`
	Instances       []string      `json:"Instances" xml:"Instances"`
	ScopeValue      []string      `json:"ScopeValue" xml:"ScopeValue"`
	Metrics         []MetricsItem `json:"Metrics" xml:"Metrics"`
}
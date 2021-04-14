package r_kvstore

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

// Children is a nested struct in r_kvstore response
type Children struct {
	Id                  int64  `json:"Id" xml:"Id"`
	Name                string `json:"Name" xml:"Name"`
	BizType             string `json:"BizType" xml:"BizType"`
	ReplicaSize         int    `json:"ReplicaSize" xml:"ReplicaSize"`
	Modifier            int    `json:"Modifier" xml:"Modifier"`
	ServiceVersion      string `json:"ServiceVersion" xml:"ServiceVersion"`
	DiskSizeMB          int    `json:"DiskSizeMB" xml:"DiskSizeMB"`
	Nickname            string `json:"Nickname" xml:"Nickname"`
	PrimaryInsName      string `json:"PrimaryInsName" xml:"PrimaryInsName"`
	ClassCode           string `json:"ClassCode" xml:"ClassCode"`
	Creator             int    `json:"Creator" xml:"Creator"`
	ResourceGroupName   string `json:"ResourceGroupName" xml:"ResourceGroupName"`
	Health              string `json:"Health" xml:"Health"`
	BinlogRetentionDays int    `json:"BinlogRetentionDays" xml:"BinlogRetentionDays"`
	UserId              string `json:"UserId" xml:"UserId"`
	LockReason          string `json:"LockReason" xml:"LockReason"`
	Service             string `json:"Service" xml:"Service"`
	Capacity            int64  `json:"Capacity" xml:"Capacity"`
	BandWidth           int64  `json:"BandWidth" xml:"BandWidth"`
	Connections         int64  `json:"Connections" xml:"Connections"`
	Items               []Item `json:"Items" xml:"Items"`
}

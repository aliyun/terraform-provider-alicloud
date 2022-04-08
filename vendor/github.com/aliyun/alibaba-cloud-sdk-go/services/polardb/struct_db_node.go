package polardb

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

// DBNode is a nested struct in polardb response
type DBNode struct {
	MaxIOPS          int    `json:"MaxIOPS" xml:"MaxIOPS"`
	DBNodeClass      string `json:"DBNodeClass" xml:"DBNodeClass"`
	FailoverPriority int    `json:"FailoverPriority" xml:"FailoverPriority"`
	DBNodeRole       string `json:"DBNodeRole" xml:"DBNodeRole"`
	MaxConnections   int    `json:"MaxConnections" xml:"MaxConnections"`
	RegionId         string `json:"RegionId" xml:"RegionId"`
	ZoneId           string `json:"ZoneId" xml:"ZoneId"`
	DBNodeStatus     string `json:"DBNodeStatus" xml:"DBNodeStatus"`
	CreationTime     string `json:"CreationTime" xml:"CreationTime"`
	DBNodeId         string `json:"DBNodeId" xml:"DBNodeId"`
	ImciSwitch       string `json:"ImciSwitch" xml:"ImciSwitch"`
	HotReplicaMode   string `json:"HotReplicaMode" xml:"HotReplicaMode"`
}

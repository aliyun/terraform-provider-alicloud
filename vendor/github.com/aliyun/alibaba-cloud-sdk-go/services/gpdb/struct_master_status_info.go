package gpdb

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

// MasterStatusInfo is a nested struct in gpdb response
type MasterStatusInfo struct {
	NormalNodeNum       int `json:"NormalNodeNum" xml:"NormalNodeNum"`
	ExceptionNodeNum    int `json:"ExceptionNodeNum" xml:"ExceptionNodeNum"`
	NotSyncingNodeNum   int `json:"NotSyncingNodeNum" xml:"NotSyncingNodeNum"`
	SyncedNodeNum       int `json:"SyncedNodeNum" xml:"SyncedNodeNum"`
	PreferredNodeNum    int `json:"PreferredNodeNum" xml:"PreferredNodeNum"`
	NotPreferredNodeNum int `json:"NotPreferredNodeNum" xml:"NotPreferredNodeNum"`
}

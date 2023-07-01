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

// SQLSlowLog is a nested struct in polardb response
type SQLSlowLog struct {
	SQLText              string `json:"SQLText" xml:"SQLText"`
	ReturnMaxRowCount    int64  `json:"ReturnMaxRowCount" xml:"ReturnMaxRowCount"`
	CreateTime           string `json:"CreateTime" xml:"CreateTime"`
	MaxExecutionTime     int64  `json:"MaxExecutionTime" xml:"MaxExecutionTime"`
	ParseTotalRowCounts  int64  `json:"ParseTotalRowCounts" xml:"ParseTotalRowCounts"`
	TotalLockTimes       int64  `json:"TotalLockTimes" xml:"TotalLockTimes"`
	TotalExecutionTimes  int64  `json:"TotalExecutionTimes" xml:"TotalExecutionTimes"`
	DBNodeId             string `json:"DBNodeId" xml:"DBNodeId"`
	SQLHASH              string `json:"SQLHASH" xml:"SQLHASH"`
	ParseMaxRowCount     int64  `json:"ParseMaxRowCount" xml:"ParseMaxRowCount"`
	MaxLockTime          int64  `json:"MaxLockTime" xml:"MaxLockTime"`
	ReturnTotalRowCounts int64  `json:"ReturnTotalRowCounts" xml:"ReturnTotalRowCounts"`
	DBName               string `json:"DBName" xml:"DBName"`
	TotalExecutionCounts int64  `json:"TotalExecutionCounts" xml:"TotalExecutionCounts"`
}

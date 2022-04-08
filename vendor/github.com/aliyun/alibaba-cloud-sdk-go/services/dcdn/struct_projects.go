package dcdn

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

// Projects is a nested struct in dcdn response
type Projects struct {
	ProjectName  string  `json:"ProjectName" xml:"ProjectName"`
	Type         string  `json:"Type" xml:"Type"`
	DomainName   string  `json:"DomainName" xml:"DomainName"`
	FieldName    string  `json:"FieldName" xml:"FieldName"`
	SamplingRate float64 `json:"SamplingRate" xml:"SamplingRate"`
	DataCenter   string  `json:"DataCenter" xml:"DataCenter"`
	SLSRegion    string  `json:"SLSRegion" xml:"SLSRegion"`
	SLSProject   string  `json:"SLSProject" xml:"SLSProject"`
	SLSLogStore  string  `json:"SLSLogStore" xml:"SLSLogStore"`
	BusinessType string  `json:"BusinessType" xml:"BusinessType"`
}

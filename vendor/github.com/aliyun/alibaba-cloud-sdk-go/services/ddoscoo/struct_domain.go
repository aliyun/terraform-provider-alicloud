package ddoscoo

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

// Domain is a nested struct in ddoscoo response
type Domain struct {
	Domain        string        `json:"Domain" xml:"Domain"`
	CcEnabled     bool          `json:"CcEnabled" xml:"CcEnabled"`
	CcRuleEnabled bool          `json:"CcRuleEnabled" xml:"CcRuleEnabled"`
	CcTemplate    string        `json:"CcTemplate" xml:"CcTemplate"`
	SslProtocols  string        `json:"SslProtocols" xml:"SslProtocols"`
	SslCiphers    string        `json:"SslCiphers" xml:"SslCiphers"`
	Http2Enable   bool          `json:"Http2Enable" xml:"Http2Enable"`
	Cname         string        `json:"Cname" xml:"Cname"`
	CertName      string        `json:"CertName" xml:"CertName"`
	WhiteList     []string      `json:"WhiteList" xml:"WhiteList"`
	BlackList     []string      `json:"BlackList" xml:"BlackList"`
	ProxyTypeList []ProxyConfig `json:"ProxyTypeList" xml:"ProxyTypeList"`
	RealServers   []RealServer  `json:"RealServers" xml:"RealServers"`
}

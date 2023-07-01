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

// IspCity is a nested struct in cms response
type IspCity struct {
	Region                string `json:"Region" xml:"Region"`
	Country               string `json:"Country" xml:"Country"`
	AreaEn                string `json:"Area.en" xml:"Area.en"`
	IPV4ProbeCount        string `json:"IPV4ProbeCount" xml:"IPV4ProbeCount"`
	CityNameZhCN          string `json:"CityName.zh_CN" xml:"CityName.zh_CN"`
	IspNameZhCN           string `json:"IspName.zh_CN" xml:"IspName.zh_CN"`
	IspName               string `json:"IspName" xml:"IspName"`
	CountryEn             string `json:"Country.en" xml:"Country.en"`
	City                  string `json:"City" xml:"City"`
	BrowserProbeCount     int    `json:"BrowserProbeCount" xml:"BrowserProbeCount"`
	WinProbeCount         int    `json:"WinProbeCount" xml:"WinProbeCount"`
	RegionZhCN            string `json:"Region.zh_CN" xml:"Region.zh_CN"`
	IspNameEn             string `json:"IspName.en" xml:"IspName.en"`
	CountryZhCN           string `json:"Country.zh_CN" xml:"Country.zh_CN"`
	CityNameEn            string `json:"CityName.en" xml:"CityName.en"`
	Isp                   string `json:"Isp" xml:"Isp"`
	CityName              string `json:"CityName" xml:"CityName"`
	RegionEn              string `json:"Region.en" xml:"Region.en"`
	IPV6ProbeCount        string `json:"IPV6ProbeCount" xml:"IPV6ProbeCount"`
	AreaZhCN              string `json:"Area.zh_CN" xml:"Area.zh_CN"`
	Ipv4BrowserProbeCount int    `json:"Ipv4BrowserProbeCount" xml:"Ipv4BrowserProbeCount"`
	APIProbeCount         string `json:"APIProbeCount" xml:"APIProbeCount"`
	IPPool                IPPool `json:"IPPool" xml:"IPPool"`
}

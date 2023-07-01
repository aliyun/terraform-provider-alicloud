package cbn

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

// TransitRouterMulticastGroup is a nested struct in cbn response
type TransitRouterMulticastGroup struct {
	GroupIpAddress                     string `json:"GroupIpAddress" xml:"GroupIpAddress"`
	TransitRouterAttachmentId          string `json:"TransitRouterAttachmentId" xml:"TransitRouterAttachmentId"`
	VSwitchId                          string `json:"VSwitchId" xml:"VSwitchId"`
	NetworkInterfaceId                 string `json:"NetworkInterfaceId" xml:"NetworkInterfaceId"`
	PeerTransitRouterMulticastDomainId string `json:"PeerTransitRouterMulticastDomainId" xml:"PeerTransitRouterMulticastDomainId"`
	Status                             string `json:"Status" xml:"Status"`
	GroupSource                        bool   `json:"GroupSource" xml:"GroupSource"`
	GroupMember                        bool   `json:"GroupMember" xml:"GroupMember"`
	MemberType                         string `json:"MemberType" xml:"MemberType"`
	SourceType                         string `json:"SourceType" xml:"SourceType"`
	ResourceType                       string `json:"ResourceType" xml:"ResourceType"`
	ResourceOwnerId                    int64  `json:"ResourceOwnerId" xml:"ResourceOwnerId"`
	ResourceId                         string `json:"ResourceId" xml:"ResourceId"`
	ConnectPeerId                      string `json:"ConnectPeerId" xml:"ConnectPeerId"`
	TransitRouterMulticastDomainId     string `json:"TransitRouterMulticastDomainId" xml:"TransitRouterMulticastDomainId"`
}

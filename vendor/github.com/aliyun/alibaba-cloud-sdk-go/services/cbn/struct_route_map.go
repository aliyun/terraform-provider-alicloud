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

// RouteMap is a nested struct in cbn response
type RouteMap struct {
	Status                             string                        `json:"Status" xml:"Status"`
	RouteMapId                         string                        `json:"RouteMapId" xml:"RouteMapId"`
	CenId                              string                        `json:"CenId" xml:"CenId"`
	CenRegionId                        string                        `json:"CenRegionId" xml:"CenRegionId"`
	Description                        string                        `json:"Description" xml:"Description"`
	MapResult                          string                        `json:"MapResult" xml:"MapResult"`
	Priority                           int                           `json:"Priority" xml:"Priority"`
	NextPriority                       int                           `json:"NextPriority" xml:"NextPriority"`
	CidrMatchMode                      string                        `json:"CidrMatchMode" xml:"CidrMatchMode"`
	AsPathMatchMode                    string                        `json:"AsPathMatchMode" xml:"AsPathMatchMode"`
	CommunityMatchMode                 string                        `json:"CommunityMatchMode" xml:"CommunityMatchMode"`
	CommunityOperateMode               string                        `json:"CommunityOperateMode" xml:"CommunityOperateMode"`
	Preference                         int                           `json:"Preference" xml:"Preference"`
	TransmitDirection                  string                        `json:"TransmitDirection" xml:"TransmitDirection"`
	SourceInstanceIdsReverseMatch      bool                          `json:"SourceInstanceIdsReverseMatch" xml:"SourceInstanceIdsReverseMatch"`
	DestinationInstanceIdsReverseMatch bool                          `json:"DestinationInstanceIdsReverseMatch" xml:"DestinationInstanceIdsReverseMatch"`
	GatewayZoneId                      string                        `json:"GatewayZoneId" xml:"GatewayZoneId"`
	MatchAddressType                   string                        `json:"MatchAddressType" xml:"MatchAddressType"`
	TransitRouterRouteTableId          string                        `json:"TransitRouterRouteTableId" xml:"TransitRouterRouteTableId"`
	SourceInstanceIds                  SourceInstanceIds             `json:"SourceInstanceIds" xml:"SourceInstanceIds"`
	DestinationInstanceIds             DestinationInstanceIds        `json:"DestinationInstanceIds" xml:"DestinationInstanceIds"`
	SourceRouteTableIds                SourceRouteTableIds           `json:"SourceRouteTableIds" xml:"SourceRouteTableIds"`
	DestinationRouteTableIds           DestinationRouteTableIds      `json:"DestinationRouteTableIds" xml:"DestinationRouteTableIds"`
	SourceRegionIds                    SourceRegionIds               `json:"SourceRegionIds" xml:"SourceRegionIds"`
	SourceChildInstanceTypes           SourceChildInstanceTypes      `json:"SourceChildInstanceTypes" xml:"SourceChildInstanceTypes"`
	DestinationChildInstanceTypes      DestinationChildInstanceTypes `json:"DestinationChildInstanceTypes" xml:"DestinationChildInstanceTypes"`
	DestinationCidrBlocks              DestinationCidrBlocks         `json:"DestinationCidrBlocks" xml:"DestinationCidrBlocks"`
	RouteTypes                         RouteTypes                    `json:"RouteTypes" xml:"RouteTypes"`
	MatchAsns                          MatchAsns                     `json:"MatchAsns" xml:"MatchAsns"`
	MatchCommunitySet                  MatchCommunitySet             `json:"MatchCommunitySet" xml:"MatchCommunitySet"`
	OperateCommunitySet                OperateCommunitySet           `json:"OperateCommunitySet" xml:"OperateCommunitySet"`
	PrependAsPath                      PrependAsPath                 `json:"PrependAsPath" xml:"PrependAsPath"`
	DestinationRegionIds               DestinationRegionIds          `json:"DestinationRegionIds" xml:"DestinationRegionIds"`
	OriginalRouteTableIds              OriginalRouteTableIds         `json:"OriginalRouteTableIds" xml:"OriginalRouteTableIds"`
	SrcZoneIds                         SrcZoneIds                    `json:"SrcZoneIds" xml:"SrcZoneIds"`
}

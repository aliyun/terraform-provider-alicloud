package alicloud

import (
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type VpcService struct {
	client *connectivity.AliyunClient
}

func (s *VpcService) DescribeEip(id string) (eip vpc.EipAddress, err error) {

	request := vpc.CreateDescribeEipAddressesRequest()
	request.RegionId = string(s.client.Region)
	request.AllocationId = id
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeEipAddresses(request)
	})
	if err != nil {
		return eip, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.DescribeEipAddressesResponse)
	if len(response.EipAddresses.EipAddress) <= 0 || response.EipAddresses.EipAddress[0].AllocationId != id {
		return eip, WrapErrorf(Error(GetNotFoundMessage("Eip", id)), NotFoundMsg, ProviderERROR)
	}
	eip = response.EipAddresses.EipAddress[0]
	return
}

func (s *VpcService) DescribeEipAssociation(id string) (object vpc.EipAddress, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	object, err = s.DescribeEip(parts[0])
	if err != nil {
		err = WrapError(err)
		return
	}
	if object.InstanceId != parts[1] {
		err = WrapErrorf(Error(GetNotFoundMessage("Eip Association", id)), NotFoundMsg, ProviderERROR)
	}

	return
}

func (s *VpcService) DescribeNatGateway(id string) (nat vpc.NatGateway, err error) {
	request := vpc.CreateDescribeNatGatewaysRequest()
	request.RegionId = string(s.client.Region)
	request.NatGatewayId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidNatGatewayId.NotFound"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeNatGatewaysResponse)
		if len(response.NatGateways.NatGateway) <= 0 || response.NatGateways.NatGateway[0].NatGatewayId != id {
			return WrapErrorf(Error(GetNotFoundMessage("NatGateway", id)), NotFoundMsg, ProviderERROR)
		}
		nat = response.NatGateways.NatGateway[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeVpc(id string) (v vpc.Vpc, err error) {
	request := vpc.CreateDescribeVpcsRequest()
	request.RegionId = s.client.RegionId
	request.VpcId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpcs(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVpcID.NotFound", "Forbidden.VpcNotFound"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVpcsResponse)
		if len(response.Vpcs.Vpc) < 1 || response.Vpcs.Vpc[0].VpcId != id {
			return WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundMsg, ProviderERROR)
		}
		v = response.Vpcs.Vpc[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeVpcWithTeadsl(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}

	action := "DescribeVpcs"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"VpcId":    id,
	}

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVpcID.NotFound", "Forbidden.VpcNotFound"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		v, err := jsonpath.Get("$.Vpcs.Vpc", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, id, "$.Vpcs.Vpc", response)
		}
		if len(v.([]interface{})) < 1 || v.([]interface{})[0].(map[string]interface{})["VpcId"].(string) != id {
			return WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
		object = v.([]interface{})[0].(map[string]interface{})
		return nil
	})
	return
}

func (s *VpcService) ListTagResources(resourceId, resourceType string) (object interface{}, err error) {
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceId.1": resourceId,
		"ResourceType": resourceType,
		"MaxResults":   PageSizeLarge,
	}
	tags := make([]interface{}, 0)

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, resourceId, "$.TagResources.TagResource", response))
			}
			tags = append(tags, v.([]interface{})...)
			nextToken, _ := jsonpath.Get("$.NextToken", response)
			request["NextToken"] = nextToken
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, resourceId, action, AlibabaCloudSdkGoERROR)
			return
		}
		if request["NextToken"].(string) == "" {
			break
		}
	}

	return tags, nil
}

func (s *VpcService) VpcStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcWithTeadsl(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		status := object["Status"].(string)
		for _, failState := range failStates {
			if status == failState {
				return object, status, WrapError(Error(FailedToReachTargetStatus, status))
			}
		}

		return object, status, nil
	}
}

func (s *VpcService) DescribeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.RegionId = s.client.RegionId
	request.VSwitchId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if response.VSwitchId != id {
			return WrapErrorf(Error(GetNotFoundMessage("vswitch", id)), NotFoundMsg, ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}

func (s *VpcService) VSwitchStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVSwitch(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}

		return object, object.Status, nil
	}
}

func (s *VpcService) DescribeSnatEntry(id string) (snat vpc.SnatTableEntry, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return snat, WrapError(err)
	}
	request := vpc.CreateDescribeSnatTableEntriesRequest()
	request.RegionId = string(s.client.Region)
	request.SnatTableId = parts[0]
	request.PageSize = requests.NewInteger(PageSizeLarge)

	for {
		invoker := NewInvoker()
		var response *vpc.DescribeSnatTableEntriesResponse
		var raw interface{}
		err = invoker.Run(func() error {
			raw, err = s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeSnatTableEntries(request)
			})
			response, _ = raw.(*vpc.DescribeSnatTableEntriesResponse)
			return err
		})

		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		//so judge the snatEntries length priority
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidSnatTableId.NotFound", "InvalidSnatEntryId.NotFound"}) {
				return snat, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return snat, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if len(response.SnatTableEntries.SnatTableEntry) < 1 {
			break
		}

		for _, snat := range response.SnatTableEntries.SnatTableEntry {
			if snat.SnatEntryId == parts[1] {
				return snat, nil
			}
		}

		if len(response.SnatTableEntries.SnatTableEntry) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return snat, WrapError(err)
		}
		request.PageNumber = page
	}

	return snat, WrapErrorf(Error(GetNotFoundMessage("SnatEntry", id)), NotFoundMsg, ProviderERROR)
}

func (s *VpcService) DescribeForwardEntry(id string) (entry vpc.ForwardTableEntry, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return entry, WrapError(err)
	}
	forwardTableId, forwardEntryId := parts[0], parts[1]
	request := vpc.CreateDescribeForwardTableEntriesRequest()
	request.RegionId = string(s.client.Region)
	request.ForwardTableId = forwardTableId
	request.ForwardEntryId = forwardEntryId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeForwardTableEntries(request)
		})
		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		//so judge the snatEntries length priority
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidForwardEntryId.NotFound", "InvalidForwardTableId.NotFound"}) {
				return WrapErrorf(Error(GetNotFoundMessage("ForwardEntry", id)), NotFoundMsg, ProviderERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, "ForwardEntry", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeForwardTableEntriesResponse)

		if len(response.ForwardTableEntries.ForwardTableEntry) > 0 {
			entry = response.ForwardTableEntries.ForwardTableEntry[0]
			return nil
		}

		return WrapErrorf(Error(GetNotFoundMessage("ForwardEntry", id)), NotFoundMsg, ProviderERROR)
	})
	return
}

func (s *VpcService) QueryRouteTableById(routeTableId string) (rt vpc.RouteTable, err error) {
	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = s.client.RegionId
	request.RouteTableId = routeTableId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTables(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, routeTableId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouteTablesResponse)
		if len(response.RouteTables.RouteTable) == 0 ||
			response.RouteTables.RouteTable[0].RouteTableId != routeTableId {
			return WrapErrorf(Error(GetNotFoundMessage("RouteTable", routeTableId)), NotFoundMsg, ProviderERROR)
		}
		rt = response.RouteTables.RouteTable[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeRouteEntry(id string) (*vpc.RouteEntry, error) {
	v := &vpc.RouteEntry{}
	parts, err := ParseResourceId(id, 5)
	if err != nil {
		return v, WrapError(err)
	}
	rtId, cidr, nexthop_type, nexthop_id := parts[0], parts[2], parts[3], parts[4]

	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = s.client.RegionId
	request.RouteTableId = rtId

	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			response, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouteTables(request)
			})
			raw = response
			return err
		}); err != nil {
			return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouteTablesResponse)
		if len(response.RouteTables.RouteTable) < 1 {
			break
		}
		for _, table := range response.RouteTables.RouteTable {
			for _, entry := range table.RouteEntrys.RouteEntry {
				if entry.DestinationCidrBlock == cidr && entry.NextHopType == nexthop_type && entry.InstanceId == nexthop_id {
					return &entry, nil
				}
			}
		}
		if len(response.RouteTables.RouteTable) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return v, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return v, WrapErrorf(Error(GetNotFoundMessage("RouteEntry", id)), NotFoundMsg, ProviderERROR)
}

func (s *VpcService) DescribeRouterInterface(id, regionId string) (ri vpc.RouterInterfaceType, err error) {
	request := vpc.CreateDescribeRouterInterfacesRequest()
	request.RegionId = regionId
	values := []string{id}
	filter := []vpc.DescribeRouterInterfacesFilter{
		{
			Key:   "RouterInterfaceId",
			Value: &values,
		},
	}
	request.Filter = &filter
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouterInterfaces(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouterInterfacesResponse)
		if len(response.RouterInterfaceSet.RouterInterfaceType) <= 0 ||
			response.RouterInterfaceSet.RouterInterfaceType[0].RouterInterfaceId != id {
			return WrapErrorf(Error(GetNotFoundMessage("RouterInterface", id)), NotFoundMsg, ProviderERROR)
		}
		ri = response.RouterInterfaceSet.RouterInterfaceType[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeRouterInterfaceConnection(id, regionId string) (ri vpc.RouterInterfaceType, err error) {
	ri, err = s.DescribeRouterInterface(id, regionId)
	if err != nil {
		return ri, WrapError(err)
	}

	if ri.OppositeInterfaceId == "" || ri.OppositeRouterType == "" ||
		ri.OppositeRouterId == "" || ri.OppositeInterfaceOwnerId == "" {
		return ri, WrapErrorf(Error(GetNotFoundMessage("RouterInterface", id)), NotFoundMsg, ProviderERROR)
	}
	return ri, nil
}

func (s *VpcService) DescribeCenInstanceGrant(id string) (rule vpc.CbnGrantRule, err error) {
	request := vpc.CreateDescribeGrantRulesToCenRequest()
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return rule, WrapError(err)
	}
	cenId := parts[0]
	instanceId := parts[1]
	instanceType, err := GetCenChildInstanceType(instanceId)
	if err != nil {
		return rule, WrapError(err)
	}

	request.RegionId = s.client.RegionId
	request.InstanceId = instanceId
	request.InstanceType = instanceType

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeGrantRulesToCen(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeGrantRulesToCenResponse)
		ruleList := response.CenGrantRules.CbnGrantRule
		if len(ruleList) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("GrantRules", id)), NotFoundMsg, ProviderERROR)
		}

		for ruleNum := 0; ruleNum <= len(response.CenGrantRules.CbnGrantRule)-1; ruleNum++ {
			if ruleList[ruleNum].CenInstanceId == cenId {
				rule = ruleList[ruleNum]
				return nil
			}
		}

		return WrapErrorf(Error(GetNotFoundMessage("GrantRules", id)), NotFoundMsg, ProviderERROR)
	})
	return
}

func (s *VpcService) WaitForCenInstanceGrant(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[1]
	ownerId := parts[2]
	for {
		object, err := s.DescribeCenInstanceGrant(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.CenInstanceId == instanceId && fmt.Sprint(object.CenOwnerId) == ownerId && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.CenInstanceId, instanceId, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *VpcService) DescribeCommonBandwidthPackage(id string) (v vpc.CommonBandwidthPackage, err error) {
	request := vpc.CreateDescribeCommonBandwidthPackagesRequest()
	request.RegionId = s.client.RegionId
	request.BandwidthPackageId = id
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCommonBandwidthPackages(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeCommonBandwidthPackagesResponse)
		//Finding the commonBandwidthPackageId
		for _, bandPackage := range response.CommonBandwidthPackages.CommonBandwidthPackage {
			if bandPackage.BandwidthPackageId == id {
				v = bandPackage
				return nil
			}
		}
		return WrapErrorf(Error(GetNotFoundMessage("CommonBandWidthPackage", id)), NotFoundMsg, ProviderERROR)
	})
	return
}

func (s *VpcService) DescribeCommonBandwidthPackageAttachment(id string) (v vpc.CommonBandwidthPackage, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return v, WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]

	object, err := s.DescribeCommonBandwidthPackage(bandwidthPackageId)
	if err != nil {
		return v, WrapError(err)
	}

	for _, ipAddresse := range object.PublicIpAddresses.PublicIpAddresse {
		if ipAddresse.AllocationId == ipInstanceId {
			v = object
			return
		}
	}
	return v, WrapErrorf(Error(GetNotFoundMessage("CommonBandWidthPackageAttachment", id)), NotFoundMsg, ProviderERROR)
}

func (s *VpcService) DescribeRouteTable(id string) (v vpc.RouterTableListType, err error) {
	request := vpc.CreateDescribeRouteTableListRequest()
	request.RegionId = s.client.RegionId
	request.RouteTableId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTableList(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouteTableListResponse)
		//Finding the routeTableId
		for _, routerTableType := range response.RouterTableList.RouterTableListType {
			if routerTableType.RouteTableId == id {
				v = routerTableType
				return nil
			}
		}
		return WrapErrorf(Error(GetNotFoundMessage("RouteTable", id)), NotFoundMsg, ProviderERROR)
	})
	return v, WrapError(err)
}

func (s *VpcService) DescribeRouteTableAttachment(id string) (v vpc.RouterTableListType, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return v, WrapError(err)
	}
	invoker := NewInvoker()
	routeTableId := parts[0]
	vSwitchId := parts[1]

	err = invoker.Run(func() error {
		object, err := s.DescribeRouteTable(routeTableId)
		if err != nil {
			return WrapError(err)
		}

		for _, id := range object.VSwitchIds.VSwitchId {
			if id == vSwitchId {
				v = object
				return nil
			}
		}
		return WrapErrorf(Error(GetNotFoundMessage("RouteTableAttachment", id)), NotFoundMsg, ProviderERROR)
	})
	return v, WrapError(err)
}

func (s *VpcService) WaitForVSwitch(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeVSwitch(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForNatGateway(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNatGateway(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForRouteEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouteEntry(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForAllRouteEntriesAvailable(routeTableId string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		table, err := s.QueryRouteTableById(routeTableId)
		if err != nil {
			return WrapError(err)
		}
		success := true
		for _, routeEntry := range table.RouteEntrys.RouteEntry {
			if routeEntry.Status != string(Available) {
				success = false
				break
			}
		}
		if success {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, routeTableId, GetFunc(1), timeout, Available, Null, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *VpcService) WaitForRouterInterface(id, regionId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouterInterface(id, regionId)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForRouterInterfaceConnection(id, regionId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouterInterfaceConnection(id, regionId)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForEip(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEip(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForEipAssociation(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEipAssociation(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) DeactivateRouterInterface(interfaceId string) error {
	request := vpc.CreateDeactivateRouterInterfaceRequest()
	request.RegionId = s.client.RegionId
	request.RouterInterfaceId = interfaceId
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeactivateRouterInterface(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "RouterInterface", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *VpcService) ActivateRouterInterface(interfaceId string) error {
	request := vpc.CreateActivateRouterInterfaceRequest()
	request.RegionId = s.client.RegionId
	request.RouterInterfaceId = interfaceId
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ActivateRouterInterface(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "RouterInterface", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *VpcService) WaitForForwardEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeForwardEntry(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForSnatEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeSnatEntry(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}

	}
}

func (s *VpcService) WaitForCommonBandwidthPackage(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCommonBandwidthPackage(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
	}
}

func (s *VpcService) WaitForCommonBandwidthPackageAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCommonBandwidthPackageAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
	}
}

// Flattens an array of vpc.public_ip_addresses into a []map[string]string
func (s *VpcService) FlattenPublicIpAddressesMappings(list []vpc.PublicIpAddresse) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		l := map[string]interface{}{
			"ip_address":    i.IpAddress,
			"allocation_id": i.AllocationId,
		}
		result = append(result, l)
	}

	return result
}

func (s *VpcService) WaitForRouteTable(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouteTable(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.Status == string(status) {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *VpcService) WaitForRouteTableAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouteTableAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.Status == string(status) {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *VpcService) DescribeNetworkAcl(id string) (networkAcl vpc.NetworkAcl, err error) {

	request := vpc.CreateDescribeNetworkAclsRequest()
	request.RegionId = s.client.RegionId
	request.NetworkAclId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeNetworkAcls(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidNetworkAcl.NotFound"}) {
			return networkAcl, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return networkAcl, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*vpc.DescribeNetworkAclsResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(response.NetworkAcls.NetworkAcl) <= 0 || response.NetworkAcls.NetworkAcl[0].NetworkAclId != id {
		return networkAcl, WrapErrorf(Error(GetNotFoundMessage("NetworkAcl", id)), NotFoundMsg, ProviderERROR)
	}
	return response.NetworkAcls.NetworkAcl[0], nil
}

func (s *VpcService) DescribeNetworkAclAttachment(id string, resource []vpc.Resource) (err error) {

	invoker := NewInvoker()
	return invoker.Run(func() error {
		object, err := s.DescribeNetworkAcl(id)
		if err != nil {
			return WrapError(err)
		}
		if len(object.Resources.Resource) < 1 {
			return WrapErrorf(Error(GetNotFoundMessage("Network Acl Attachment", id)), NotFoundMsg, ProviderERROR)
		}
		success := true
		for _, source := range object.Resources.Resource {
			success = false
			for _, res := range resource {
				if source.ResourceId == res.ResourceId {
					success = true
				}
			}
			if success == false {
				return WrapErrorf(Error(GetNotFoundMessage("Network Acl Attachment", id)), NotFoundMsg, ProviderERROR)
			}
		}
		return nil
	})
}

func (s *VpcService) WaitForNetworkAcl(networkAclId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNetworkAcl(networkAclId)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		success := true
		// Check Acl's binding resources
		for _, resource := range object.Resources.Resource {
			if resource.Status != string(BINDED) {
				success = false
			}
		}
		if object.Status == string(status) && success == true {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, networkAclId, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) WaitForNetworkAclAttachment(id string, resource []vpc.Resource, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		err := s.DescribeNetworkAclAttachment(id, resource)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		object, err := s.DescribeNetworkAcl(id)
		success := true
		// Check Acl's binding resources
		for _, resource := range object.Resources.Resource {
			if resource.Status != string(BINDED) {
				success = false
			}
		}
		if object.Status == string(status) && success == true {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *VpcService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []vpc.TagResource, err error) {
	request := vpc.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []vpc.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, vpc.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ListTagResources(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*vpc.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

func (s *VpcService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		conn, err := s.client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}

		if len(removed) > 0 {
			action := "UnTagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceId.1": d.Id(),
				"ResourceType": string(resourceType),
			}
			for i, key := range removed {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceId.1": d.Id(),
				"ResourceType": string(resourceType),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value	", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		d.SetPartial("tags")
	}

	return nil
}

func (s *VpcService) tagsToMap(tags []vpc.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *VpcService) ignoreTag(t vpc.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

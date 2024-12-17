package alicloud

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeEipAddresses(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeEipAddressesResponse)
		if len(response.EipAddresses.EipAddress) <= 0 || response.EipAddresses.EipAddress[0].AllocationId != id {
			return resource.NonRetryableError(WrapErrorf(Error(GetNotFoundMessage("Eip", id)), NotFoundMsg, ProviderERROR))
		}
		eip = response.EipAddresses.EipAddress[0]
		return nil
	})

	return
}

func (s *VpcService) DescribeEipAssociation(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	object, err = s.DescribeEipAddress(parts[0])
	if err != nil {
		return object, WrapError(err)
	}

	if object["InstanceId"] != parts[1] {
		return object, WrapErrorf(Error(GetNotFoundMessage("Eip:Association", id)), NotFoundWithResponse)
	}

	return
}

func (s *VpcService) DescribeNatGateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeNatGateways"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"NatGatewayId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidNatGatewayId.NotFound", "InvalidRegionId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("NatGateway", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.NatGateways.NatGateway", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NatGateways.NatGateway", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["NatGatewayId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeVpc(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVpcs"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"VpcId":    id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.VpcNotFound", "InvalidVpcID.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("Vpc", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Vpcs.Vpc", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Vpcs.Vpc", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["VpcId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
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
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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

func (s *VpcService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *VpcService) VpcStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpc(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
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

func (s *VpcService) DescribeVSwitchWithTeadsl(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"VSwitchId": id,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	if v, ok := response["VSwitchId"].(string); ok && v != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("vswitch", id)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
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

func (s *VpcService) DescribeSnatEntry(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeSnatTableEntries"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"SnatEntryId": parts[1],
		"SnatTableId": parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSnatEntryId.NotFound", "InvalidSnatTableId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC:SnatEntry", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SnatTableEntries.SnatTableEntry", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SnatTableEntries.SnatTableEntry", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["SnatEntryId"].(string) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeForwardEntry(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeForwardTableEntries"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"ForwardEntryId": parts[1],
		"ForwardTableId": parts[0],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidForwardEntryId.NotFound", "InvalidForwardTableId.NotFound", "InvalidRegionId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ForwardEntry", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ForwardTableEntries.ForwardTableEntry", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ForwardTableEntries.ForwardTableEntry", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["ForwardEntryId"].(string) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) QueryRouteTableById(routeTableId string) (rt vpc.RouteTable, err error) {
	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = s.client.RegionId
	request.RouteTableId = routeTableId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTables(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouteTablesResponse)
		if len(response.RouteTables.RouteTable) == 0 ||
			response.RouteTables.RouteTable[0].RouteTableId != routeTableId {
			return resource.NonRetryableError(WrapErrorf(Error(GetNotFoundMessage("RouteTable", routeTableId)), NotFoundMsg, ProviderERROR))
		}
		rt = response.RouteTables.RouteTable[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeRouteEntry(id string) (*vpc.RouteEntry, error) {
	v := &vpc.RouteEntry{}
	var raw interface{}
	parts, err := ParseResourceId(id, 5)
	if err != nil {
		return v, WrapError(err)
	}
	rtId, cidr, nexthop_type, nexthop_id := parts[0], parts[2], parts[3], parts[4]

	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = s.client.RegionId
	request.RouteTableId = rtId

	if strings.Contains(cidr, "_") {
		cidr = strings.Replace(cidr, "_", ":", -1)
	}

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeRouteTables(request)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			raw = response
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		response, _ := raw.(*vpc.DescribeRouteTablesResponse)
		if len(response.RouteTables.RouteTable) < 1 {
			return v, WrapErrorf(Error(GetNotFoundMessage("RouteEntry", id)), NotFoundWithResponse, response)
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

func (s *VpcService) DescribeCommonBandwidthPackage(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeCommonBandwidthPackages"
	request := map[string]interface{}{
		"RegionId":           s.client.RegionId,
		"BandwidthPackageId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRegionId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("CommonBandwidthPackage", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.CommonBandwidthPackages.CommonBandwidthPackage", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CommonBandwidthPackages.CommonBandwidthPackage", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["BandwidthPackageId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeRouteTable(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeRouteTableList"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"RouteTableId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		err = WrapErrorf(err, ResponseCodeMsg, id, action, response)
		return object, err
	}
	v, err := jsonpath.Get("$.RouterTableList.RouterTableListType", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RouterTableList.RouterTableListType", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["RouteTableId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeRouteTableAttachment(id string) (v map[string]interface{}, err error) {
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

		if val, ok := object["VSwitchIds"].(map[string]interface{}); ok {
			if vs, ok := val["VSwitchId"]; ok {
				for _, id := range vs.([]interface{}) {
					if fmt.Sprint(id) == vSwitchId {
						v = object
						return nil
					}
				}
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

func (s *VpcService) EipAssociationStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEipAssociation(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, object["Status"]))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
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

		if fmt.Sprint(object["Status"]) == string(status) {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, fmt.Sprint(object["Status"]), string(status), ProviderERROR)
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *VpcService) DescribeNetworkAcl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeNetworkAclAttributes"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"NetworkAclId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidNetworkAcl.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC:NetworkAcl", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NetworkAclAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NetworkAclAttribute", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeNetworkAclAttachment(id string, resource []vpc.Resource) (err error) {

	invoker := NewInvoker()
	return invoker.Run(func() error {
		object, err := s.DescribeNetworkAcl(id)
		if err != nil {
			return WrapError(err)
		}
		resources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		if len(resources) < 1 {
			return WrapErrorf(Error(GetNotFoundMessage("Network Acl Attachment", id)), NotFoundMsg, ProviderERROR)
		}
		success := true
		for _, source := range resources {
			success = false
			for _, res := range resource {
				item := source.(map[string]interface{})
				if fmt.Sprint(item["ResourceId"]) == res.ResourceId {
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
		resources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		// Check Acl's binding resources
		for _, res := range resources {
			item := res.(map[string]interface{})
			if fmt.Sprint(item["Status"]) != string(BINDED) {
				success = false
			}
		}
		if fmt.Sprint(object["Status"]) == string(status) && success == true {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, networkAclId, GetFunc(1), timeout, fmt.Sprint(object["Status"]), string(status), ProviderERROR)
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
		resources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		// Check Acl's binding resources
		for _, res := range resources {
			item := res.(map[string]interface{})
			if fmt.Sprint(item["Status"]) != string(BINDED) {
				success = false
			}
		}
		if fmt.Sprint(object["Status"]) == string(status) && success == true {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, fmt.Sprint(object["Status"]), string(status), ProviderERROR)
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
					if NeedRetry(err) {
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
					if NeedRetry(err) {
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

func (s *VpcService) SetInstanceSecondaryCidrBlocks(d *schema.ResourceData) error {
	if d.HasChange("secondary_cidr_blocks") {
		oraw, nraw := d.GetChange("secondary_cidr_blocks")
		removed := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		added := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		conn, err := s.client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		if len(removed) > 0 {
			action := "UnassociateVpcCidrBlock"
			request := map[string]interface{}{
				"RegionId": s.client.RegionId,
				"VpcId":    d.Id(),
			}
			for _, item := range removed {
				request["SecondaryCidrBlock"] = item
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				addDebug(action, response, request)
			}
		}

		if len(added) > 0 {
			action := "AssociateVpcCidrBlock"
			request := map[string]interface{}{
				"RegionId": s.client.RegionId,
				"VpcId":    d.Id(),
			}
			for _, item := range added {
				request["SecondaryCidrBlock"] = item
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				addDebug(action, response, request)
			}
		}
		d.SetPartial("secondary_cidr_blocks")
	}
	return nil
}

func (s *VpcService) DescribeNatGatewayTransform(id string) ([]interface{}, error) {
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}

	action := "GetNatGatewayConvertStatus"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"NatGatewayId": id,
	}
	request["ClientToken"] = buildClientToken("GetNatGatewayConvertStatus")

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)

	response, err1 := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err1 != nil {
		if IsExpectedErrors(err1, []string{"InvalidVpcID.NotFound", "Forbidden.VpcNotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	ob, err2 := jsonpath.Get("$.ConvertSteps", response)
	if err2 != nil {
		return nil, WrapErrorf(err2, FailedGetAttributeMsg, id, "$.ConvertSteps", response)
	}

	natType, err3 := jsonpath.Get("$.DstNatType", response)
	if err3 != nil {
		return nil, WrapErrorf(err3, FailedGetAttributeMsg, id, "$.DstNatType", response)
	}

	object := ob.([]interface{})
	if len(object) < 1 || natType.(string) != "Enhanced" {
		return nil, WrapErrorf(Error(GetNotFoundMessage("NAT", id)), NotFoundWithResponse, response)
	}
	return object, nil

}

func (s *VpcService) WaitForNatGatewayTransform(id string, timeout int64) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNatGatewayTransform(id)
		if err != nil {
			if NotFoundError(err) {
				return err
			}
			if IsExpectedErrors(err, []string{"OperationFailed.NatGwRouteInMiddleStatusError"}) {
				return nil
			}
			return err
		}

		isOk := false
		for _, v := range object {
			val := v.(map[string]interface{})
			if val["StepName"].(string) == "end_success" && val["StepStatus"].(string) == "successful" {
				isOk = true
				break
			}
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", "", ProviderERROR)
		}
		if isOk {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *VpcService) DescribeRouteTableList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeRouteTableList"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"VpcId":      id,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "200" {
			err = fmt.Errorf("DescribeRouteTableList failed, response: %v ", response)
			return object, err
		}
		v, err := jsonpath.Get("$.RouterTableList.RouterTableListType", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RouterTableList.RouterTableListType", response)
		}
		result, _ := v.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if item["RouteTableType"] == "System" {
				object = item
				return object, nil
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return object, nil
}

func (s *VpcService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		conn, err := s.client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UnTagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
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
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
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

func (s *VpcService) DescribeVswitch(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"VSwitchId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVSwitchId.NotFound", "InvalidVswitchID.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("Vswitch", id)), NotFoundWithError, err)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if fmt.Sprint(object["VSwitchId"]) != id {
		return object, WrapErrorf(Error(GetNotFoundMessage("vswitch", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *VpcService) VswitchStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVswitch(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *VpcService) SnatEntryStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSnatEntry(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) ForwardEntryStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeForwardEntry(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *VpcService) RouteTableStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRouteTable(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *VpcService) CommonBandwidthPackageStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCommonBandwidthPackage(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *VpcService) DescribeHavip(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeHaVips"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PageNumber": 1,
		"PageSize":   20,
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidFilterKey.ValueNotSupported", "InvalidHaVipId.NotFound", "InvalidRegionId.NotFound"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("Havip", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(action, response, request)
		v, err := jsonpath.Get("$.HaVips.HaVip", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HaVips.HaVip", response)
		}
		result, _ := v.([]interface{})
		if len(result) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
		for _, v := range result {
			if v.(map[string]interface{})["HaVipId"].(string) == id {
				return v.(map[string]interface{}), nil
			}
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
}

func (s *VpcService) HavipStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHavip(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *VpcService) NatGatewayStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNatGateway(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *VpcService) DescribeVpcFlowLog(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeFlowLogs"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"FlowLogId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.FlowLogs.FlowLog", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.FlowLogs.FlowLog", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["FlowLogId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpcFlowLogStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcFlowLog(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *VpcService) NetworkAclStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNetworkAcl(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DeleteAclResources(id string) (object map[string]interface{}, err error) {
	acl, err := s.DescribeNetworkAcl(id)
	if err != nil {
		return object, WrapError(err)
	}
	deleteResources := make([]map[string]interface{}, 0)
	res, err := jsonpath.Get("$.Resources.Resource", acl)
	if err != nil {
		return object, WrapError(err)
	}
	resources, _ := res.([]interface{})
	if resources != nil && len(resources) < 1 {
		return object, nil
	}
	for _, val := range resources {
		item, _ := val.(map[string]interface{})
		if item["Status"] == "UNBINDING" {
			continue
		}

		deleteResources = append(deleteResources, map[string]interface{}{
			"ResourceId":   item["ResourceId"],
			"ResourceType": item["ResourceType"],
		})
	}
	if len(deleteResources) == 0 {
		return nil, nil
	}

	var response map[string]interface{}
	request := map[string]interface{}{
		"NetworkAclId": id,
		"Resource":     deleteResources,
		"RegionId":     s.client.RegionId,
	}
	action := "UnassociateNetworkAcl"
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	request["ClientToken"] = buildClientToken("UnassociateNetworkAcl")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy", "ResourceStatus.Error"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidResource.NotBinding"}) {
			return object, nil
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, 10*time.Minute, 5*time.Second, s.NetworkAclStateRefreshFunc(id, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return response, WrapErrorf(err, IdMsg, id)
	}
	return object, nil
}

func (s *VpcService) DescribeEipAddress(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeEipAddresses"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"AllocationId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.EipAddresses.EipAddress", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.EipAddresses.EipAddress", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("EIP", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["AllocationId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("EIP", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) EipAddressStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEipAddress(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeExpressConnectPhysicalConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribePhysicalConnections"

	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}

	filterMapList := make([]map[string]interface{}, 0)
	filterMapList = append(filterMapList, map[string]interface{}{
		"Key":   "PhysicalConnectionId",
		"Value": []string{id},
	})
	request["Filter"] = filterMapList

	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.PhysicalConnectionSet.PhysicalConnectionType", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PhysicalConnectionSet.PhysicalConnectionType", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnect:PhysicalConnection", id)), NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["PhysicalConnectionId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnect:PhysicalConnection", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *VpcService) ExpressConnectPhysicalConnectionStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectPhysicalConnection(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeExpressConnectVirtualBorderRouter(id string, includeCrossAccountVbr bool) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVirtualBorderRouters"
	request := map[string]interface{}{
		"RegionId":               s.client.RegionId,
		"PageNumber":             1,
		"PageSize":               50,
		"IncludeCrossAccountVbr": includeCrossAccountVbr,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnect", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["VbrId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnect", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) ExpressConnectVirtualBorderRouterStateRefreshFunc(id string, includeCrossAccountVbr bool, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectVirtualBorderRouter(id, includeCrossAccountVbr)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcDhcpOptionsSet(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetDhcpOptionsSet"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"DhcpOptionsSetId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRegionId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC:DhcpOptionsSet", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if _, ok := object["DhcpOptionsSetId"]; !ok {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC:DhcpOptionsSet", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	return object, nil
}

func (s *VpcService) VpcDhcpOptionsSetStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcDhcpOptionsSet(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcDhcpOptionsSetAttachment(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}
	object, err = s.DescribeVpcDhcpOptionsSet(parts[1])
	if err != nil {
		return object, WrapError(err)
	}
	return object, nil
}

func (s *VpcService) DescribeVpcDhcpOptionsSetAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		parts, err := ParseResourceId(id, 2)
		if err != nil {
			return nil, "", WrapError(err)
		}
		object, err := s.DescribeVpcDhcpOptionsSetAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		status := ""
		if associateVpcsList, ok := object["AssociateVpcs"]; ok {
			for _, associateVpcsListItem := range associateVpcsList.([]interface{}) {
				if associateVpcsListItem != nil {
					associateVpcsListItemMap, ok := associateVpcsListItem.(map[string]interface{})
					if ok && associateVpcsListItemMap["VpcId"] == parts[0] {
						status = associateVpcsListItemMap["AssociateStatus"].(string)
						break
					}
				}
			}
		}

		for _, failState := range failStates {
			if status == failState {
				return object, fmt.Sprint(object["AssociateStatus"]), WrapError(Error(FailedToReachTargetStatus, status))
			}
		}
		return object, status, nil
	}
}

func (s *VpcService) DescribeVpcNatIpCidr(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListNatIpCidrs"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"NatGatewayId": parts[0],
		"NatIpCidr":    parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("ListNatIpCidrs")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NatIpCidrs", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NatIpCidrs", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["NatIpCidr"]) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeVpcNatIp(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListNatIps"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"NatGatewayId": parts[0],
		"NatIpIds":     []string{parts[1]},
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("ListNatIps")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NatIps", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NatIps", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["NatIpId"]) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpcNatIpStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcNatIp(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["NatIpStatus"]) == failState {
				return object, fmt.Sprint(object["NatIpStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["NatIpStatus"])))
			}
		}
		return object, fmt.Sprint(object["NatIpStatus"]), nil
	}
}

func (s *VpcService) DescribeVpcTrafficMirrorFilter(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTrafficMirrorFilters"
	request := map[string]interface{}{
		"RegionId":               s.client.RegionId,
		"TrafficMirrorFilterIds": []string{id},
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TrafficMirrorFilters", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TrafficMirrorFilters", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["TrafficMirrorFilterId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpcTrafficMirrorFilterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcTrafficMirrorFilter(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TrafficMirrorFilterStatus"]) == failState {
				return object, fmt.Sprint(object["TrafficMirrorFilterStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["TrafficMirrorFilterStatus"])))
			}
		}
		return object, fmt.Sprint(object["TrafficMirrorFilterStatus"]), nil
	}
}

func (s *VpcService) DescribeVpcTrafficMirrorFilterEgressRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTrafficMirrorFilters"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}

	request := map[string]interface{}{
		"RegionId":               s.client.RegionId,
		"MaxResults":             10,
		"TrafficMirrorFilterIds": []string{parts[0]},
	}
	idExist := false

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TrafficMirrorFilters", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TrafficMirrorFilters", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		EgressRules := v.(map[string]interface{})["EgressRules"]
		if EgressRulesMap, ok := EgressRules.([]interface{}); ok {
			for _, v := range EgressRulesMap {
				if fmt.Sprint(v.(map[string]interface{})["TrafficMirrorFilterRuleId"]) == parts[1] {
					idExist = true
					return v.(map[string]interface{}), nil
				}
			}
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) VpcTrafficMirrorFilterEgressRuleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcTrafficMirrorFilterEgressRule(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TrafficMirrorFilterRuleStatus"]) == failState {
				return object, fmt.Sprint(object["TrafficMirrorFilterRuleStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["TrafficMirrorFilterRuleStatus"])))
			}
		}
		return object, fmt.Sprint(object["TrafficMirrorFilterRuleStatus"]), nil
	}
}

func (s *VpcService) DescribeVpcTrafficMirrorFilterIngressRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTrafficMirrorFilters"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}

	request := map[string]interface{}{
		"RegionId":               s.client.RegionId,
		"MaxResults":             10,
		"TrafficMirrorFilterIds": []string{parts[0]},
	}
	idExist := false

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TrafficMirrorFilters", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TrafficMirrorFilters", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		EgressRules := v.(map[string]interface{})["IngressRules"]
		if EgressRulesMap, ok := EgressRules.([]interface{}); ok {
			for _, v := range EgressRulesMap {
				if fmt.Sprint(v.(map[string]interface{})["TrafficMirrorFilterRuleId"]) == parts[1] {
					idExist = true
					return v.(map[string]interface{}), nil
				}
			}
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) VpcTrafficMirrorFilterIngressRuleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcTrafficMirrorFilterIngressRule(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TrafficMirrorFilterRuleStatus"]) == failState {
				return object, fmt.Sprint(object["TrafficMirrorFilterRuleStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["TrafficMirrorFilterRuleStatus"])))
			}
		}
		return object, fmt.Sprint(object["TrafficMirrorFilterRuleStatus"]), nil
	}
}

func (s *VpcService) DescribeVpcTrafficMirrorSession(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTrafficMirrorSessions"
	request := map[string]interface{}{
		"RegionId":                s.client.RegionId,
		"TrafficMirrorSessionIds": []string{id},
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TrafficMirrorSessions", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TrafficMirrorSessions", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["TrafficMirrorSessionId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpcTrafficMirrorSessionStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcTrafficMirrorSession(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TrafficMirrorSessionStatus"]) == failState {
				return object, fmt.Sprint(object["TrafficMirrorSessionStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["TrafficMirrorSessionStatus"])))
			}
		}
		return object, fmt.Sprint(object["TrafficMirrorSessionStatus"]), nil
	}
}

func (s *VpcService) DescribeVpcIpv6Gateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeIpv6GatewayAttribute"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"Ipv6GatewayId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if Ipv6GatewayId, ok := object["Ipv6GatewayId"]; !ok || fmt.Sprint(Ipv6GatewayId) == "" {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *VpcService) VpcIpv6GatewayStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpv6Gateway(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcIpv6EgressRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	action := "DescribeIpv6EgressOnlyRules"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"PageSize":      PageSizeLarge,
		"PageNumber":    1,
		"Ipv6GatewayId": parts[0],
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Ipv6EgressOnlyRules.Ipv6EgressOnlyRule", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Ipv6EgressOnlyRules.Ipv6EgressOnlyRule", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Ipv6EgressOnlyRuleId"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) VpcIpv6EgressRuleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpv6EgressRule(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {

			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcIpv6InternetBandwidth(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeIpv6Addresses"
	request := map[string]interface{}{
		"RegionId":                s.client.RegionId,
		"Ipv6InternetBandwidthId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Ipv6Addresses.Ipv6Address", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Ipv6Addresses.Ipv6Address", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["Ipv6InternetBandwidth"].(map[string]interface{})["Ipv6InternetBandwidthId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeVpcBgpGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBgpGroups"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"BgpGroupId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.BgpGroups.BgpGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BgpGroups.BgpGroup", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["BgpGroupId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpcBgpGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcBgpGroup(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcVbrHa(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVbrHa"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"VbrHaId":  id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DescribeVbrHa")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if vbrHaId, ok := object["VbrHaId"]; !ok || fmt.Sprint(vbrHaId) == "" {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *VpcService) VpcVbrHaStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcVbrHa(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcBgpPeer(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBgpPeers"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"BgpPeerId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.BgpPeers.BgpPeer", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BgpPeers.BgpPeer", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["BgpPeerId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpcBgpPeerStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcBgpPeer(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcBgpNetwork(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBgpNetworks"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"RouterId":   parts[0],
		"PageNumber": 1,
		"PageSize":   PageSizeMedium,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.BgpNetworks.BgpNetwork", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BgpNetworks.BgpNetwork", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["DstCidrBlock"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) VpcBgpNetworkStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcBgpNetwork(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpnGateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVpnGateway"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"VpnGatewayId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVpnGatewayInstanceId.NotFound", "InvalidRegionId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VpnGateway", id)), NotFoundWithResponse, response)
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}

	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeVpnGateways(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVpnGateways"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"VpnGatewayId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidCustomerGatewayInstanceId"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VpnCustomerGateway", id)), NotFoundWithResponse, response)
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	totalCount, err := jsonpath.Get("$.TotalCount", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TotalCount", response)
	}
	totalCountInt, err := totalCount.(json.Number).Int64()
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TotalCount", response)
	}
	if totalCountInt == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VpnCustomerGateway", id)), NotFoundWithResponse, response)
	}

	addDebug(action, response, request)
	v, err := jsonpath.Get("$.VpnGateways.VpnGateway[0]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpnGatewayStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpnGateway(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *VpcService) DescribeVpnCustomerGateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeCustomerGateway"
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"CustomerGatewayId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidCustomerGatewayInstanceId"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VpnCustomerGateway", id)), NotFoundWithResponse, response)
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}

	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeVpnRouteEntry(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return object, WrapError(err)
	}
	action := "DescribeVpnRouteEntries"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return object, WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"VpnGatewayId": parts[0],
		"PageNumber":   1,
		"PageSize":     PageSizeMedium,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.VpnRouteEntries.VpnRouteEntry", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VpnRouteEntries.VpnRouteEntry", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["NextHop"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["RouteDest"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) VpnRouteEntryStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpnRouteEntry(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), WrapError(Error(FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *VpcService) DescribeVpnIpsecServer(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListIpsecServers"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"IpsecServerId": []string{id},
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.IpsecServers", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.IpsecServers", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["IpsecServerId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeVpnConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVpnConnection"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"VpnConnectionId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVpnConnectionInstanceId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VpnCustomerGateway", id)), NotFoundWithResponse, response)
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}

	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpnConnectionStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpnConnection(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *VpcService) DescribeVpnPbrRouteEntry(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return object, WrapError(err)
	}
	action := "DescribeVpnPbrRouteEntries"
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return object, WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"VpnGatewayId": parts[0],
		"PageNumber":   1,
		"PageSize":     PageSizeMedium,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.VpnPbrRouteEntries.VpnPbrRouteEntry", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VpnPbrRouteEntries.VpnPbrRouteEntry", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["NextHop"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["RouteSource"]) == parts[2] && fmt.Sprint(v.(map[string]interface{})["RouteDest"]) == parts[3] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) VpnPbrRouteEntryStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpnPbrRouteEntry(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), WrapError(Error(FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *VpcService) DescribeVpnGatewayVpnAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVpnConnection"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"VpnConnectionId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVpnConnectionInstanceId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VpnConnection", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpnGatewayVpnAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpnGatewayVpnAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), WrapError(Error(FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *VpcService) DescribeVpcIpv4Gateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetIpv4GatewayAttribute"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"Ipv4GatewayId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.Ipv4Gateway"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC:Ipv4Gateway", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpcIpv4GatewayStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcIpv4Gateway(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) GetVpcPrefixListEntries(id string) (objects []map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetVpcPrefixListEntries"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"PrefixListId": id,
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return objects, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		if formatInt(response["TotalCount"]) == 0 {
			return objects, nil
		}
		resp, err := jsonpath.Get("$.PrefixListEntry", response)
		if err != nil {
			return objects, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return objects, nil
}

func (s *VpcService) DescribeVpcPrefixList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListPrefixLists"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"PrefixListIds": []string{id},
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.PrefixLists", response)
	if formatInt(response["TotalCount"]) != 0 && err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PrefixLists", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(resp.([]interface{})[0].(map[string]interface{})["PrefixListId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = resp.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpcPrefixListStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcPrefixList(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpnGatewayVcoRoute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVcoRouteEntries"
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"VpnConnectionId": parts[0],
		"PageNumber":      1,
		"PageSize":        PageSizeMedium,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			if IsExpectedErrors(err, []string{"UnknownError"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("VPC:VPNGateway", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.VcoRouteEntries", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VcoRouteEntries", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPNGateway", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			item := v.(map[string]interface{})
			if fmt.Sprint(item["RouteDest"]) == parts[1] && fmt.Sprint(item["NextHop"]) == parts[2] && fmt.Sprint(item["Weight"]) == parts[3] {
				idExist = true
				return item, nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPNGateway", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *VpcService) DescribeVpcIpv4CidrBlock(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeVpcs"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"VpcId":      parts[0],
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Vpcs.Vpc", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Vpcs.Vpc", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VpcIpv4CidrBlock", id)), NotFoundWithResponse, response)
	}

	vpcObject := v.([]interface{})[0].(map[string]interface{})

	v, err = jsonpath.Get("$.SecondaryCidrBlocks.SecondaryCidrBlock", vpcObject)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SecondaryCidrBlocks.SecondaryCidrBlock", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VpcIpv4CidrBlock", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if v.(string) == parts[1] {
			idExist = true
			return object, nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("VpcIpv4CidrBlock", id)), NotFoundWithResponse, response)
	}
	return
}

func (s *VpcService) DescribeVpcPublicIpAddressPool(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListPublicIpAddressPools"
	request := map[string]interface{}{
		"RegionId":               s.client.RegionId,
		"PublicIpAddressPoolIds": []string{id},
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.PublicIpAddressPoolList", response)
	if formatInt(response["TotalCount"]) != 0 && err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PublicIpAddressPoolList", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(resp.([]interface{})[0].(map[string]interface{})["PublicIpAddressPoolId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = resp.([]interface{})[0].(map[string]interface{})

	return object, nil
}

func (s *VpcService) VpcPublicIpAddressPoolStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcPublicIpAddressPool(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcPublicIpAddressPoolCidrBlock(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}
	action := "ListPublicIpAddressPoolCidrBlocks"
	request := map[string]interface{}{
		"RegionId":              s.client.RegionId,
		"PublicIpAddressPoolId": parts[0],
		"CidrBlock":             parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.PublicIpPoolCidrBlockList", response)
	if formatInt(response["TotalCount"]) != 0 && err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PublicIpPoolCidrBlockList", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(resp.([]interface{})[0].(map[string]interface{})["CidrBlock"]) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = resp.([]interface{})[0].(map[string]interface{})

	return object, nil
}

func (s *VpcService) VpcPublicIpAddressPoolCidrBlockStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcPublicIpAddressPoolCidrBlock(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcNetworkAclAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeNetworkAcls"

	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ClientToken":  buildClientToken("DescribeNetworkAcls"),
		"NetworkAclId": parts[0],
		"ResourceId":   parts[1],
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.NetworkAcls.NetworkAcl", response)
	if formatInt(response["TotalCount"]) != 0 && err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NetworkAcls.NetworkAcl", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC:NetworkAclAttachment", id)), NotFoundWithResponse, response)
	}

	for _, networkAcls := range resp.([]interface{}) {
		ruleActionArg := networkAcls.(map[string]interface{})
		if networkAclId, ok := ruleActionArg["NetworkAclId"]; ok && fmt.Sprint(networkAclId) == parts[0] {
			if resources, ok := ruleActionArg["Resources"]; ok {
				resourcesArg := resources.(map[string]interface{})
				if resourcesArgs, ok := resourcesArg["Resource"].([]interface{}); ok {
					for _, networkAclResource := range resourcesArgs {
						resourceArg := networkAclResource.(map[string]interface{})
						if resourceId, ok := resourceArg["ResourceId"]; ok && fmt.Sprint(resourceId) == parts[1] {
							idExist = true
							return resourceArg, nil
						}
					}
				}
			}
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("VPC:NetworkAclAttachment", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *VpcService) VpcNetworkAclAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcNetworkAclAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeVpcGatewayRouteTableAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetIpv4GatewayAttribute"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"Ipv4GatewayId": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.Ipv4Gateway"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC:Ipv4Gateway", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpcGatewayRouteTableAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcGatewayRouteTableAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeExpressConnectGrantRuleToCen(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeGrantRulesToCen"

	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}

	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ClientToken":  buildClientToken("DescribeGrantRulesToCen"),
		"InstanceId":   parts[2],
		"InstanceType": "VBR",
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.CenGrantRules.CbnGrantRule", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CenGrantRules.CbnGrantRule", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnect:GrantRuleToCen", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["CenInstanceId"]) == parts[0] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnect:GrantRuleToCen", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *VpcService) DescribeExpressConnectRouterInterface(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"Filter.1.Key":   "RouterInterfaceId",
		"Filter.1.Value": []string{id},
		"RegionId":       s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeRouterInterfaces"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.RouterInterfaceSet.RouterInterfaceType", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RouterInterfaceSet.RouterInterfaceType", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("RouterInterface", id)), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *VpcService) ExpressConnectRouterInterfaceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterInterface(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *VpcService) DescribeHighDefinitionMonitorLogAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeHighDefinitionMonitorLogAttribute"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.EipAddresses.EipAddress", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *VpcService) DescribeVpcHaVipAttachment(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return object, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"Filter.1.Key":     "HaVipId",
		"Filter.1.Value.1": parts[0],
		"RegionId":         s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeHaVips"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HaVips.HaVip", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HaVips.HaVip", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("HaVipAttachment", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["HaVipId"]) != parts[0] {
			return object, WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *VpcService) VpcHaVipAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVpcHaVipAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		status84 := object["Status"]
		for _, failState := range failStates {
			if fmt.Sprint(status84) == failState {
				return object, fmt.Sprint(status84), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(status84)))
			}
		}
		return object, fmt.Sprint(status84), nil
	}
}

func (s *VpcService) ModifyExpressConnectPhysicalConnectionStatus(d *schema.ResourceData, status string) (err error) {
	var response map[string]interface{}
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	switch status {
	case "Confirmed":
		action := "ConfirmPhysicalConnection"

		confirmedReq := map[string]interface{}{
			"RegionId":             s.client.RegionId,
			"ClientToken":          buildClientToken("ConfirmPhysicalConnection"),
			"PhysicalConnectionId": d.Id(),
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(20*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, confirmedReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, confirmedReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Confirmed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	case "Enabled":
		action := "CreatePhysicalConnectionOccupancyOrder"

		createPhysicalConnectionOccupancyOrderReq := map[string]interface{}{
			"RegionId":             s.client.RegionId,
			"ClientToken":          buildClientToken("CreatePhysicalConnectionOccupancyOrder"),
			"PhysicalConnectionId": d.Id(),
			"InstanceChargeType":   "PrePaid",
			"AutoPay":              true,
		}

		if v, ok := d.GetOkExists("period"); ok {
			createPhysicalConnectionOccupancyOrderReq["Period"] = v
		}

		if v, ok := d.GetOk("pricing_cycle"); ok {
			createPhysicalConnectionOccupancyOrderReq["PricingCycle"] = v
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(20*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, createPhysicalConnectionOccupancyOrderReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, createPhysicalConnectionOccupancyOrderReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if resp, err := jsonpath.Get("$.Data", response); err != nil || resp == nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Data", response)
		} else {
			orderId := resp.(map[string]interface{})["OrderId"]
			d.Set("order_id", orderId)
		}

		stateConf := BuildStateConf([]string{}, []string{"Enabled"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	case "Canceled":
		action := "CancelPhysicalConnection"

		canceledReq := map[string]interface{}{
			"RegionId":             s.client.RegionId,
			"ClientToken":          buildClientToken("CancelPhysicalConnection"),
			"PhysicalConnectionId": d.Id(),
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(20*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, canceledReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, canceledReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Canceled"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	case "Terminated":
		action := "TerminatePhysicalConnection"

		terminatedReq := map[string]interface{}{
			"RegionId":             s.client.RegionId,
			"ClientToken":          buildClientToken("TerminatePhysicalConnection"),
			"PhysicalConnectionId": d.Id(),
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(20*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, terminatedReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, terminatedReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Terminated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return nil
}

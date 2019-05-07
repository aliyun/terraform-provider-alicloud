package alicloud

import (
	"time"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
	addDebug(request.GetActionName(), raw)
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
			if IsExceptedError(err, InvalidNatGatewayIdNotFound) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*vpc.DescribeNatGatewaysResponse)
		if len(response.NatGateways.NatGateway) <= 0 || response.NatGateways.NatGateway[0].NatGatewayId != id {
			return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		nat = response.NatGateways.NatGateway[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeVpc(id string) (v vpc.DescribeVpcAttributeResponse, err error) {
	request := vpc.CreateDescribeVpcAttributeRequest()
	request.VpcId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpcAttribute(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidVpcIDNotFound, ForbiddenVpcNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*vpc.DescribeVpcAttributeResponse)
		if response.VpcId != id {
			return WrapErrorf(Error(GetNotFoundMessage("VPC", id)), NotFoundMsg, ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}

func (s *VpcService) DescribeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.VSwitchId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if err != nil {
			if IsExceptedError(err, InvalidVswitchIDNotFound) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if response.VSwitchId != id {
			return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		v = *response
		return nil
	})
	return
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
			if IsExceptedErrors(err, []string{InvalidSnatTableIdNotFound, InvalidSnatEntryIdNotFound}) {
				return snat, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return snat, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)

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
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return snat, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return snat, WrapErrorf(Error(GetNotFoundMessage("SnatEntry", id)), NotFoundMsg, ProviderERROR)
}

func (s *VpcService) DescribeForwardEntry(forwardTableId string, forwardEntryId string) (entry vpc.ForwardTableEntry, err error) {

	args := vpc.CreateDescribeForwardTableEntriesRequest()
	args.RegionId = string(s.client.Region)
	args.ForwardTableId = forwardTableId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeForwardTableEntries(args)
		})
		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		//so judge the snatEntries length priority
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidForwardEntryIdNotFound, InvalidForwardTableIdNotFound}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("Forward Entry", forwardTableId))
			}
			return err
		}
		resp, _ := raw.(*vpc.DescribeForwardTableEntriesResponse)
		if resp == nil || len(resp.ForwardTableEntries.ForwardTableEntry) <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("Forward Entry", forwardTableId))
		}

		for _, forward := range resp.ForwardTableEntries.ForwardTableEntry {
			if forward.ForwardEntryId == forwardEntryId {
				entry = forward
				return nil
			}
		}

		return GetNotFoundErrorFromString(GetNotFoundMessage("Forward Entry", forwardTableId))
	})
	return
}

func (s *VpcService) QueryRouteTableById(routeTableId string) (rt vpc.RouteTable, err error) {
	request := vpc.CreateDescribeRouteTablesRequest()
	request.RouteTableId = routeTableId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTables(request)
		})
		if err != nil {
			return err
		}
		rts, _ := raw.(*vpc.DescribeRouteTablesResponse)
		if rts == nil || len(rts.RouteTables.RouteTable) == 0 ||
			rts.RouteTables.RouteTable[0].RouteTableId != routeTableId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("Route Table", routeTableId))
		}

		rt = rts.RouteTables.RouteTable[0]
		return nil
	})
	return
}

func (s *VpcService) QueryRouteEntry(routeTableId, cidrBlock, nextHopType, nextHopId string) (rn vpc.RouteEntry, err error) {
	rt, err := s.QueryRouteTableById(routeTableId)
	if err != nil {
		return
	}

	for _, e := range rt.RouteEntrys.RouteEntry {
		if e.DestinationCidrBlock == cidrBlock && e.NextHopType == nextHopType && e.InstanceId == nextHopId {
			return e, nil
		}
	}
	return rn, GetNotFoundErrorFromString(GetNotFoundMessage("Route Entry", routeTableId))
}

func (s *VpcService) DescribeRouterInterface(regionId, interfaceId string) (ri vpc.RouterInterfaceType, err error) {
	request := vpc.CreateDescribeRouterInterfacesRequest()
	request.RegionId = regionId
	values := []string{interfaceId}
	filter := []vpc.DescribeRouterInterfacesFilter{{
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
			return err
		}
		resp, _ := raw.(*vpc.DescribeRouterInterfacesResponse)
		if resp == nil || len(resp.RouterInterfaceSet.RouterInterfaceType) <= 0 ||
			resp.RouterInterfaceSet.RouterInterfaceType[0].RouterInterfaceId != interfaceId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("Router Interface", interfaceId))
		}
		ri = resp.RouterInterfaceSet.RouterInterfaceType[0]
		return nil
	})
	return
}

func (s *VpcService) DescribeGrantRulesToCen(id string) (rule vpc.CbnGrantRule, err error) {
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
		resp, _ := raw.(*vpc.DescribeGrantRulesToCenResponse)
		ruleList := resp.CenGrantRules.CbnGrantRule
		if resp == nil || len(ruleList) <= 0 {
			return WrapErrorf(Error(GetNotFoundMessage("GrantRules", id)), NotFoundMsg, ProviderERROR)
		}

		for ruleNum := 0; ruleNum <= len(resp.CenGrantRules.CbnGrantRule)-1; ruleNum++ {
			if ruleList[ruleNum].CenInstanceId == cenId {
				rule = ruleList[ruleNum]
				return nil
			}
		}

		return WrapErrorf(Error(GetNotFoundMessage("GrantRules", id)), NotFoundMsg, ProviderERROR)
	})
	return
}

func (s *VpcService) WaitForVpc(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeVpc(id)
		if err != nil {
			if NotFoundError(err) && status == Deleted {
				return nil
			}
			return WrapError(err)
		}
		if object.Status == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
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

func (s *VpcService) WaitForAllRouteEntries(routeTableId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {

		table, err := s.QueryRouteTableById(routeTableId)

		if err != nil {
			return err
		}

		success := true

		for _, routeEntry := range table.RouteEntrys.RouteEntry {
			if routeEntry.Status != string(status) {
				success = false
				break
			}
		}
		if success {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("All Route Entries", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *VpcService) WaitForRouterInterface(regionId, interfaceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		result, err := s.DescribeRouterInterface(regionId, interfaceId)
		if err != nil {
			if !NotFoundError(err) {
				return err
			}
		} else if result.Status == string(status) {
			break
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Router Interface", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
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
	req := vpc.CreateDeactivateRouterInterfaceRequest()
	req.RouterInterfaceId = interfaceId
	_, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeactivateRouterInterface(req)
	})
	if err != nil {
		return fmt.Errorf("Deactivating RouterInterface %s got an error: %#v.", interfaceId, err)
	}
	return nil
}

func (s *VpcService) ActivateRouterInterface(interfaceId string) error {
	req := vpc.CreateActivateRouterInterfaceRequest()
	req.RouterInterfaceId = interfaceId
	_, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ActivateRouterInterface(req)
	})
	if err != nil {
		return fmt.Errorf("Activating RouterInterface %s got an error: %#v.", interfaceId, err)
	}
	return nil
}

func (s *VpcService) WaitForForwardEntry(tableId, id string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		forward, err := s.DescribeForwardEntry(tableId, id)
		if err != nil {
			if !NotFoundError(err) {
				return WrapError(err)
			}
		} else if forward.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return WrapError(Error(GetTimeoutMessage("Forward Entry", string(status))))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
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

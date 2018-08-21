package alicloud

import (
	"time"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

func (client *AliyunClient) BuildVpcCommonRequest(region string) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	endpoint := LoadEndpoint(client.RegionId, VPCCode)
	if region == "" {
		region = client.RegionId
	}
	if endpoint == "" {
		endpoint, _ = client.DescribeEndpointByCode(region, VPCCode)
	}
	if endpoint == "" {
		endpoint = fmt.Sprintf("vpc.%s.aliyuncs.com", region)
	}
	request.Domain = endpoint
	request.Version = ApiVersion20160428
	request.RegionId = region
	return request
}

func (client *AliyunClient) DescribeEipAddress(allocationId string) (eip vpc.EipAddress, err error) {

	args := vpc.CreateDescribeEipAddressesRequest()
	args.RegionId = string(client.Region)
	args.AllocationId = allocationId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		eips, err := client.vpcconn.DescribeEipAddresses(args)
		if err != nil {
			return err
		}
		if eips == nil || len(eips.EipAddresses.EipAddress) <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("EIP", allocationId))
		}
		eip = eips.EipAddresses.EipAddress[0]
		return nil
	})
	return
}

func (client *AliyunClient) DescribeNatGateway(natGatewayId string) (nat vpc.NatGateway, err error) {

	args := vpc.CreateDescribeNatGatewaysRequest()
	args.RegionId = string(client.Region)
	args.NatGatewayId = natGatewayId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		resp, err := client.vpcconn.DescribeNatGateways(args)
		if err != nil {
			if IsExceptedError(err, InvalidNatGatewayIdNotFound) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("Nat Gateway", natGatewayId))
			}
			return err
		}

		if resp == nil || len(resp.NatGateways.NatGateway) <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("Nat Gateway", natGatewayId))
		}

		nat = resp.NatGateways.NatGateway[0]
		return nil
	})
	return
}

func (client *AliyunClient) DescribeVpc(vpcId string) (v vpc.DescribeVpcAttributeResponse, err error) {
	request := vpc.CreateDescribeVpcAttributeRequest()
	request.VpcId = vpcId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		resp, err := client.vpcconn.DescribeVpcAttribute(request)
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidVpcIDNotFound, ForbiddenVpcNotFound}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("VPC", vpcId))
			}
			return err
		}
		if resp == nil || resp.VpcId != vpcId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("VPC", vpcId))
		}
		v = *resp
		return nil
	})
	return
}

func (client *AliyunClient) DescribeVswitch(vswitchId string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.VSwitchId = vswitchId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		resp, err := client.vpcconn.DescribeVSwitchAttributes(request)
		if err != nil {
			if IsExceptedError(err, InvalidVswitchIDNotFound) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("VSwitch", vswitchId))
			}
			return err
		}
		if resp == nil || resp.VSwitchId != vswitchId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("VSwitch", vswitchId))
		}
		v = *resp
		return nil
	})
	return
}

func (client *AliyunClient) DescribeSnatEntry(snatTableId string, snatEntryId string) (snat vpc.SnatTableEntry, err error) {

	request := vpc.CreateDescribeSnatTableEntriesRequest()
	request.RegionId = string(client.Region)
	request.SnatTableId = snatTableId
	request.PageSize = requests.NewInteger(PageSizeLarge)

	for {
		invoker := NewInvoker()
		var snatEntries *vpc.DescribeSnatTableEntriesResponse
		err = invoker.Run(func() error {
			resp, err := client.vpcconn.DescribeSnatTableEntries(request)
			snatEntries = resp
			return err
		})

		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		//so judge the snatEntries length priority
		if err != nil {
			if IsExceptedError(err, InvalidSnatTableIdNotFound) {
				return snat, GetNotFoundErrorFromString(GetNotFoundMessage("Snat Entry", snatEntryId))
			}
			return snat, err
		}

		if snatEntries == nil || len(snatEntries.SnatTableEntries.SnatTableEntry) < 1 {
			break
		}

		for _, snat := range snatEntries.SnatTableEntries.SnatTableEntry {
			if snat.SnatEntryId == snatEntryId {
				return snat, nil
			}
		}

		if len(snatEntries.SnatTableEntries.SnatTableEntry) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return snat, err
		} else {
			request.PageNumber = page
		}
	}

	return snat, GetNotFoundErrorFromString(GetNotFoundMessage("Snat Entry", snatEntryId))
}

func (client *AliyunClient) DescribeForwardEntry(forwardTableId string, forwardEntryId string) (entry vpc.ForwardTableEntry, err error) {

	args := vpc.CreateDescribeForwardTableEntriesRequest()
	args.RegionId = string(client.Region)
	args.ForwardTableId = forwardTableId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		resp, err := client.vpcconn.DescribeForwardTableEntries(args)
		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		//so judge the snatEntries length priority
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidForwardEntryIdNotFound, InvalidForwardTableIdNotFound}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("Forward Entry", forwardTableId))
			}
			return err
		}

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

func (client *AliyunClient) QueryRouteTableById(routeTableId string) (rt vpc.RouteTable, err error) {
	request := vpc.CreateDescribeRouteTablesRequest()
	request.RouteTableId = routeTableId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		rts, err := client.vpcconn.DescribeRouteTables(request)
		if err != nil {
			return err
		}

		if rts == nil || len(rts.RouteTables.RouteTable) == 0 ||
			rts.RouteTables.RouteTable[0].RouteTableId != routeTableId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("Route Table", routeTableId))
		}

		rt = rts.RouteTables.RouteTable[0]
		return nil
	})
	return
}

func (client *AliyunClient) QueryRouteEntry(routeTableId, cidrBlock, nextHopType, nextHopId string) (rn vpc.RouteEntry, err error) {
	rt, err := client.QueryRouteTableById(routeTableId)
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

func (client *AliyunClient) DescribeRouterInterface(regionId, interfaceId string) (ri vpc.RouterInterfaceTypeInDescribeRouterInterfaces, err error) {
	request := vpc.CreateDescribeRouterInterfacesRequest()
	request.RegionId = regionId
	values := []string{interfaceId}
	filter := []vpc.DescribeRouterInterfacesFilter{vpc.DescribeRouterInterfacesFilter{
		Key:   "RouterInterfaceId",
		Value: &values,
	},
	}
	request.Filter = &filter

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		resp, err := client.vpcconn.DescribeRouterInterfaces(request)
		if err != nil {
			return err
		}
		if resp == nil || len(resp.RouterInterfaceSet.RouterInterfaceType) <= 0 ||
			resp.RouterInterfaceSet.RouterInterfaceType[0].RouterInterfaceId != interfaceId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("Router Interface", interfaceId))
		}
		ri = resp.RouterInterfaceSet.RouterInterfaceType[0]
		return nil
	})
	return
}

func (client *AliyunClient) WaitForVpc(vpcId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		vpc, err := client.DescribeVpc(vpcId)
		if err != nil {
			return err
		}
		if vpc.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("VPC", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (client *AliyunClient) WaitForVSwitch(vswitchId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		vswitch, err := client.DescribeVswitch(vswitchId)
		if err != nil {
			return err
		}
		if vswitch.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("VSwitch", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (client *AliyunClient) WaitForAllRouteEntries(routeTableId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {

		table, err := client.QueryRouteTableById(routeTableId)

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

func (client *AliyunClient) WaitForRouterInterface(regionId, interfaceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		result, err := client.DescribeRouterInterface(regionId, interfaceId)
		if err != nil {
			return err
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

func (client *AliyunClient) WaitForEip(allocationId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		eip, err := client.DescribeEipAddress(allocationId)
		if err != nil {
			if !NotFoundError(err) {
				return err
			}
		} else if eip.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("EIP", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (client *AliyunClient) DeactivateRouterInterface(interfaceId string) error {
	req := vpc.CreateDeactivateRouterInterfaceRequest()
	req.RouterInterfaceId = interfaceId
	if _, err := client.vpcconn.DeactivateRouterInterface(req); err != nil {
		return fmt.Errorf("Deactivating RouterInterface %s got an error: %#v.", interfaceId, err)
	}
	return nil
}

func (client *AliyunClient) ActivateRouterInterface(interfaceId string) error {
	req := vpc.CreateActivateRouterInterfaceRequest()
	req.RouterInterfaceId = interfaceId
	if _, err := client.vpcconn.ActivateRouterInterface(req); err != nil {
		return fmt.Errorf("Activating RouterInterface %s got an error: %#v.", interfaceId, err)
	}
	return nil
}

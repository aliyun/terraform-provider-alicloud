package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/ecs"
)

const Negative = ecs.Spec("Negative")

func (client *AliyunClient) DescribeEipAddress(allocationId string) (eip vpc.EipAddress, err error) {

	args := vpc.CreateDescribeEipAddressesRequest()
	args.RegionId = string(client.Region)
	args.AllocationId = allocationId

	eips, err := client.vpcconn.DescribeEipAddresses(args)
	if err != nil {
		return
	}
	if eips == nil || len(eips.EipAddresses.EipAddress) <= 0 {
		return eip, GetNotFoundErrorFromString(GetNotFoundMessage("EIP", allocationId))
	}

	return eips.EipAddresses.EipAddress[0], nil
}

func (client *AliyunClient) DescribeNatGateway(natGatewayId string) (nat vpc.NatGateway, err error) {

	args := vpc.CreateDescribeNatGatewaysRequest()
	args.RegionId = string(client.Region)
	args.NatGatewayId = natGatewayId

	resp, err := client.vpcconn.DescribeNatGateways(args)
	if err != nil {
		if IsExceptedError(err, InvalidNatGatewayIdNotFound) {
			return nat, GetNotFoundErrorFromString(GetNotFoundMessage("Nat Gateway", natGatewayId))
		}
		return
	}

	if resp == nil || len(resp.NatGateways.NatGateway) <= 0 {
		return nat, GetNotFoundErrorFromString(GetNotFoundMessage("Nat Gateway", natGatewayId))
	}

	return resp.NatGateways.NatGateway[0], nil
}

func (client *AliyunClient) DescribeVpc(vpcId string) (v vpc.DescribeVpcAttributeResponse, err error) {
	request := vpc.CreateDescribeVpcAttributeRequest()
	request.VpcId = vpcId

	resp, err := client.vpcconn.DescribeVpcAttribute(request)
	if err != nil {
		if IsExceptedError(err, InvalidVpcIDNotFound) || IsExceptedError(err, ForbiddenVpcNotFound) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPC", vpcId))
		}
		return
	}
	if resp == nil || resp.VpcId != vpcId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPC", vpcId))
	}
	return *resp, nil
}

func (client *AliyunClient) DescribeVswitch(vswitchId string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.VSwitchId = vswitchId

	resp, err := client.vpcconn.DescribeVSwitchAttributes(request)
	if err != nil {
		if IsExceptedError(err, InvalidVswitchIDNotFound) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VSwitch", vswitchId))
		}
		return
	}
	if resp == nil || resp.VSwitchId != vswitchId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VSwitch", vswitchId))
	}
	return *resp, nil
}

func (client *AliyunClient) DescribeSnatEntry(snatTableId string, snatEntryId string) (snat vpc.SnatTableEntry, err error) {

	request := vpc.CreateDescribeSnatTableEntriesRequest()
	request.RegionId = string(client.Region)
	request.SnatTableId = snatTableId
	request.PageSize = requests.NewInteger(PageSizeLarge)

	for {
		snatEntries, err := client.vpcconn.DescribeSnatTableEntries(request)

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
		request.PageNumber = request.PageNumber + requests.NewInteger(1)
	}

	return snat, GetNotFoundErrorFromString(GetNotFoundMessage("Snat Entry", snatEntryId))
}

func (client *AliyunClient) DescribeForwardEntry(forwardTableId string, forwardEntryId string) (entry vpc.ForwardTableEntry, err error) {

	args := vpc.CreateDescribeForwardTableEntriesRequest()
	args.RegionId = string(client.Region)
	args.ForwardTableId = forwardTableId

	resp, err := client.vpcconn.DescribeForwardTableEntries(args)
	//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
	//so judge the snatEntries length priority
	if err != nil {
		if IsExceptedError(err, InvalidForwardEntryIdNotFound) ||
			IsExceptedError(err, InvalidForwardTableIdNotFound) {
			return entry, GetNotFoundErrorFromString(GetNotFoundMessage("Forward Entry", forwardTableId))
		}
		return
	}

	if resp == nil || len(resp.ForwardTableEntries.ForwardTableEntry) <= 0 {
		return entry, GetNotFoundErrorFromString(GetNotFoundMessage("Forward Entry", forwardTableId))
	}

	for _, forward := range resp.ForwardTableEntries.ForwardTableEntry {
		if forward.ForwardEntryId == forwardEntryId {
			return forward, nil
		}
	}

	return entry, GetNotFoundErrorFromString(GetNotFoundMessage("Forward Entry", forwardTableId))
}

func (client *AliyunClient) QueryRouteTableById(routeTableId string) (rt vpc.RouteTable, err error) {
	request := vpc.CreateDescribeRouteTablesRequest()
	request.RouteTableId = routeTableId

	rts, err := client.vpcconn.DescribeRouteTables(request)
	if err != nil {
		return
	}

	if rts == nil || len(rts.RouteTables.RouteTable) == 0 ||
		rts.RouteTables.RouteTable[0].RouteTableId != routeTableId {
		return rt, GetNotFoundErrorFromString(GetNotFoundMessage("Route Table", routeTableId))
	}

	return rts.RouteTables.RouteTable[0], nil
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

func (client *AliyunClient) DescribeRouterInterface(interfaceId string) (ri vpc.RouterInterfaceTypeInDescribeRouterInterfaces, err error) {
	request := vpc.CreateDescribeRouterInterfacesRequest()
	request.RegionId = string(client.Region)
	values := []string{interfaceId}
	filter := []vpc.DescribeRouterInterfacesFilter{vpc.DescribeRouterInterfacesFilter{
		Key:   "RouterInterfaceId",
		Value: &values,
	},
	}
	request.Filter = &filter

	resp, err := client.vpcconn.DescribeRouterInterfaces(request)
	if err != nil {
		return
	}
	if resp == nil || len(resp.RouterInterfaceSet.RouterInterfaceType) <= 0 ||
		resp.RouterInterfaceSet.RouterInterfaceType[0].RouterInterfaceId != interfaceId {
		return ri, GetNotFoundErrorFromString(GetNotFoundMessage("Router Interface", interfaceId))
	}
	return resp.RouterInterfaceSet.RouterInterfaceType[0], nil
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

func (client *AliyunClient) WaitForRouterInterface(interfaceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		result, err := client.DescribeRouterInterface(interfaceId)
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
			return err
		}

		if eip.Status == string(status) {
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

func GetAllRouterInterfaceSpec() (specifications []string) {
	specifications = append(specifications, string(ecs.Large1), string(ecs.Large2),
		string(ecs.Small1), string(ecs.Small2), string(ecs.Small5), string(ecs.Middle1),
		string(ecs.Middle2), string(ecs.Middle5), string(Negative))
	return
}

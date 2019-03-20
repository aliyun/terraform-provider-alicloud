package alicloud

import (
	"fmt"
	"strings"

	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type EcsService struct {
	client *connectivity.AliyunClient
}

func (s *EcsService) JudgeRegionValidation(key, region string) error {
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeRegions(ecs.CreateDescribeRegionsRequest())
	})
	if err != nil {
		return fmt.Errorf("DescribeRegions got an error: %#v", err)
	}
	resp, _ := raw.(*ecs.DescribeRegionsResponse)
	if resp == nil || len(resp.Regions.Region) < 1 {
		return GetNotFoundErrorFromString("There is no any available region.")
	}

	var rs []string
	for _, v := range resp.Regions.Region {
		if v.RegionId == region {
			return nil
		}
		rs = append(rs, v.RegionId)
	}
	return fmt.Errorf("'%s' is invalid. Expected on %v.", key, strings.Join(rs, ", "))
}

// DescribeZone validate zoneId is valid in region
func (s *EcsService) DescribeZone(zoneID string) (zone ecs.Zone, err error) {
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(ecs.CreateDescribeZonesRequest())
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ecs.DescribeZonesResponse)
	if resp == nil || len(resp.Zones.Zone) < 1 {
		return zone, fmt.Errorf("There is no any availability zone in region %s.", s.client.RegionId)
	}

	zoneIds := []string{}
	for _, z := range resp.Zones.Zone {
		if z.ZoneId == zoneID {
			return z, nil
		}
		zoneIds = append(zoneIds, z.ZoneId)
	}
	return zone, fmt.Errorf("availability_zone %s not exists in region %s, all zones are %s", zoneID, s.client.RegionId, zoneIds)
}

func (s *EcsService) DescribeInstanceById(id string) (instance ecs.Instance, err error) {
	req := ecs.CreateDescribeInstancesRequest()
	req.InstanceIds = convertListToJsonString([]interface{}{id})

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstances(req)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ecs.DescribeInstancesResponse)
	if resp == nil || len(resp.Instances.Instance) < 1 {
		return instance, GetNotFoundErrorFromString(GetNotFoundMessage("Instance", id))
	}

	return resp.Instances.Instance[0], nil
}

func (s *EcsService) DescribeInstanceAttribute(id string) (instance ecs.DescribeInstanceAttributeResponse, err error) {
	req := ecs.CreateDescribeInstanceAttributeRequest()
	req.InstanceId = id

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstanceAttribute(req)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ecs.DescribeInstanceAttributeResponse)
	if resp == nil {
		return instance, GetNotFoundErrorFromString(GetNotFoundMessage("Instance", id))
	}

	return *resp, nil
}

func (s *EcsService) QueryInstanceSystemDisk(id string) (disk ecs.Disk, err error) {
	args := ecs.CreateDescribeDisksRequest()
	args.InstanceId = id
	args.DiskType = string(DiskTypeSystem)

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(args)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ecs.DescribeDisksResponse)
	if resp != nil && len(resp.Disks.Disk) < 1 {
		return disk, GetNotFoundErrorFromString(fmt.Sprintf("The specified system disk is not found by instance id %s.", id))
	}

	return resp.Disks.Disk[0], nil
}

// ResourceAvailable check resource available for zone
func (s *EcsService) ResourceAvailable(zone ecs.Zone, resourceType ResourceType) error {
	for _, res := range zone.AvailableResourceCreation.ResourceTypes {
		if res == string(resourceType) {
			return nil
		}
	}
	return fmt.Errorf("%s is not available in %s zone of %s region", resourceType, zone.ZoneId, s.client.Region)
}

func (s *EcsService) DiskAvailable(zone ecs.Zone, diskCategory DiskCategory) error {
	for _, disk := range zone.AvailableDiskCategories.DiskCategories {
		if disk == string(diskCategory) {
			return nil
		}
	}
	return fmt.Errorf("%s is not available in %s zone of %s region", diskCategory, zone.ZoneId, s.client.Region)
}

func (s *EcsService) JoinSecurityGroups(instanceId string, securityGroupIds []string) error {
	req := ecs.CreateJoinSecurityGroupRequest()
	req.InstanceId = instanceId
	for _, sid := range securityGroupIds {
		req.SecurityGroupId = sid
		_, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.JoinSecurityGroup(req)
		})
		if err != nil && IsExceptedErrors(err, []string{InvalidInstanceIdAlreadyExists}) {
			return err
		}
	}

	return nil
}

func (s *EcsService) LeaveSecurityGroups(instanceId string, securityGroupIds []string) error {
	req := ecs.CreateLeaveSecurityGroupRequest()
	req.InstanceId = instanceId
	for _, sid := range securityGroupIds {
		req.SecurityGroupId = sid
		_, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.LeaveSecurityGroup(req)
		})
		if err != nil && IsExceptedErrors(err, []string{InvalidSecurityGroupIdNotFound}) {
			return err
		}
	}

	return nil
}

func (s *EcsService) DescribeSecurityGroupAttribute(securityGroupId string) (group ecs.DescribeSecurityGroupAttributeResponse, err error) {
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	request.SecurityGroupId = securityGroupId

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidSecurityGroupIdNotFound}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	if response == nil {
		err = WrapErrorf(Error(GetNotFoundMessage("Security Group", securityGroupId)), NotFoundMsg, ProviderERROR)
		return
	}

	return *response, nil
}

func (s *EcsService) DescribeSecurityGroupRule(groupId, direction, ipProtocol, portRange, nicType, cidr_ip, policy string, priority int) (rule ecs.Permission, err error) {
	args := ecs.CreateDescribeSecurityGroupAttributeRequest()
	args.SecurityGroupId = groupId
	args.Direction = direction
	args.NicType = nicType

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(args)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	if resp == nil {
		return rule, GetNotFoundErrorFromString(GetNotFoundMessage("Security Group", groupId))
	}

	for _, ru := range resp.Permissions.Permission {
		if strings.ToLower(string(ru.IpProtocol)) == ipProtocol && ru.PortRange == portRange {
			cidr := ru.SourceCidrIp
			if direction == string(DirectionIngress) && cidr == "" {
				cidr = ru.SourceGroupId
			}
			if direction == string(DirectionEgress) {
				if cidr = ru.DestCidrIp; cidr == "" {
					cidr = ru.DestGroupId
				}
			}

			if cidr == cidr_ip && strings.ToLower(string(ru.Policy)) == policy && ru.Priority == strconv.Itoa(priority) {
				return ru, nil
			}
		}
	}

	return rule, GetNotFoundErrorFromString(fmt.Sprintf("Security group rule not found by group id %s.", groupId))

}

func (s *EcsService) DescribeAvailableResources(d *schema.ResourceData, meta interface{}, destination DestinationResource) (zoneId string, validZones []ecs.AvailableZone, err error) {
	client := meta.(*connectivity.AliyunClient)
	// Before creating resources, check input parameters validity according available zone.
	// If availability zone is nil, it will return all of supported resources in the current.
	args := ecs.CreateDescribeAvailableResourceRequest()
	args.DestinationResource = string(destination)
	args.IoOptimized = string(IOOptimized)

	if v, ok := d.GetOk("availability_zone"); ok && strings.TrimSpace(v.(string)) != "" {
		zoneId = strings.TrimSpace(v.(string))
	} else if v, ok := d.GetOk("vswitch_id"); ok && strings.TrimSpace(v.(string)) != "" {
		vpcService := VpcService{s.client}
		if vsw, err := vpcService.DescribeVswitch(strings.TrimSpace(v.(string))); err == nil {
			zoneId = vsw.ZoneId
		}
	}

	if v, ok := d.GetOk("instance_charge_type"); ok && strings.TrimSpace(v.(string)) != "" {
		args.InstanceChargeType = strings.TrimSpace(v.(string))
	}

	if v, ok := d.GetOk("spot_strategy"); ok && strings.TrimSpace(v.(string)) != "" {
		args.SpotStrategy = strings.TrimSpace(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok && strings.TrimSpace(v.(string)) != "" {
		args.NetworkCategory = strings.TrimSpace(v.(string))
	}

	if v, ok := d.GetOk("is_outdated"); ok && v.(bool) == true {
		args.IoOptimized = string(NoneOptimized)
	}

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeAvailableResource(args)
	})
	if err != nil {
		return "", nil, fmt.Errorf("Error DescribeAvailableResource: %#v", err)
	}
	resources, _ := raw.(*ecs.DescribeAvailableResourceResponse)

	if resources == nil || len(resources.AvailableZones.AvailableZone) < 1 {
		err = fmt.Errorf("There are no availability resources in the region: %s.", client.RegionId)
		return
	}

	valid := false
	soldout := false
	var expectedZones []string
	for _, zone := range resources.AvailableZones.AvailableZone {
		if zone.Status == string(SoldOut) {
			if zone.ZoneId == zoneId {
				soldout = true
			}
			continue
		}
		if zoneId != "" && zone.ZoneId == zoneId {
			valid = true
			validZones = append(make([]ecs.AvailableZone, 1), zone)
			break
		}
		expectedZones = append(expectedZones, zone.ZoneId)
		validZones = append(validZones, zone)
	}
	if zoneId != "" {
		if !valid {
			err = fmt.Errorf("Availability zone %s status is not available in the region %s. Expected availability zones: %s.",
				zoneId, client.RegionId, strings.Join(expectedZones, ", "))
			return
		}
		if soldout {
			err = fmt.Errorf("Availability zone %s status is sold out in the region %s. Expected availability zones: %s.",
				zoneId, client.RegionId, strings.Join(expectedZones, ", "))
			return
		}
	}

	if len(validZones) <= 0 {
		err = fmt.Errorf("There is no availability resources in the region %s. Please choose another region.", client.RegionId)
		return
	}

	return
}

func (s *EcsService) InstanceTypeValidation(targetType, zoneId string, validZones []ecs.AvailableZone) error {

	mapInstanceTypeToZones := make(map[string]string)
	var expectedInstanceTypes []string
	for _, zone := range validZones {
		if zoneId != "" && zoneId != zone.ZoneId {
			continue
		}
		for _, r := range zone.AvailableResources.AvailableResource {
			if r.Type == string(InstanceTypeResource) {
				for _, t := range r.SupportedResources.SupportedResource {
					if t.Status == string(SoldOut) {
						continue
					}
					if targetType == t.Value {
						return nil
					}

					if _, ok := mapInstanceTypeToZones[t.Value]; !ok {
						expectedInstanceTypes = append(expectedInstanceTypes, t.Value)
						mapInstanceTypeToZones[t.Value] = t.Value
					}
				}
			}
		}
	}
	if zoneId != "" {
		return fmt.Errorf("The instance type %s is solded out or is not supported in the zone %s. Expected instance types: %s", targetType, zoneId, strings.Join(expectedInstanceTypes, ", "))
	}
	return fmt.Errorf("The instance type %s is solded out or is not supported in the region %s. Expected instance types: %s", targetType, s.client.RegionId, strings.Join(expectedInstanceTypes, ", "))
}

func (s *EcsService) QueryInstancesWithKeyPair(instanceIdsStr, keypair string) (instanceIds []string, instances []ecs.Instance, err error) {

	args := ecs.CreateDescribeInstancesRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	args.InstanceIds = instanceIdsStr
	args.KeyPairName = keypair
	for true {
		raw, e := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(args)
		})
		if e != nil {
			err = e
			return
		}
		resp, _ := raw.(*ecs.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances.Instance) < 0 {
			return
		}
		for _, inst := range resp.Instances.Instance {
			instanceIds = append(instanceIds, inst.InstanceId)
			instances = append(instances, inst)
		}
		if len(instances) < PageSizeLarge {
			break
		}
		if page, e := getNextpageNumber(args.PageNumber); e != nil {
			err = e
			return
		} else {
			args.PageNumber = page
		}
	}
	return
}

func (s *EcsService) DescribeKeyPair(keyName string) (keypair ecs.KeyPair, err error) {
	req := ecs.CreateDescribeKeyPairsRequest()
	req.KeyPairName = keyName
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeKeyPairs(req)
	})

	if err != nil {
		return
	}
	resp, _ := raw.(*ecs.DescribeKeyPairsResponse)
	if resp == nil || len(resp.KeyPairs.KeyPair) < 1 {
		return keypair, GetNotFoundErrorFromString(GetNotFoundMessage("KeyPair", keyName))
	}
	return resp.KeyPairs.KeyPair[0], nil

}

func (s *EcsService) DescribeDiskById(instanceId, diskId string) (disk ecs.Disk, err error) {
	req := ecs.CreateDescribeDisksRequest()
	if instanceId != "" {
		req.InstanceId = instanceId
	}
	req.DiskIds = convertListToJsonString([]interface{}{diskId})

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(req)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ecs.DescribeDisksResponse)
	if resp == nil || len(resp.Disks.Disk) < 1 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("ECS disk", diskId))
		return
	}
	return resp.Disks.Disk[0], nil
}

func (s *EcsService) DescribeDisksByType(instanceId string, diskType DiskType) (disk []ecs.Disk, err error) {
	req := ecs.CreateDescribeDisksRequest()
	if instanceId != "" {
		req.InstanceId = instanceId
	}
	req.DiskType = string(diskType)

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(req)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ecs.DescribeDisksResponse)
	if resp == nil {
		return
	}
	return resp.Disks.Disk, nil
}

func (s *EcsService) DescribeTags(resourceId string, resourceType TagResourceType) (tags []ecs.Tag, err error) {
	req := ecs.CreateDescribeTagsRequest()
	req.ResourceType = string(resourceType)
	req.ResourceId = resourceId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeTags(req)
	})

	if err != nil {
		return
	}
	resp, _ := raw.(*ecs.DescribeTagsResponse)
	if resp == nil || len(resp.Tags.Tag) < 1 {
		err = GetNotFoundErrorFromString(fmt.Sprintf("Describe %s tag by id %s got an error.", resourceType, resourceId))
		return
	}

	return resp.Tags.Tag, nil
}

func (s *EcsService) DescribeImageById(id string) (image ecs.Image, err error) {
	req := ecs.CreateDescribeImagesRequest()
	req.ImageId = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeImages(req)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ecs.DescribeImagesResponse)
	if resp == nil || len(resp.Images.Image) < 1 {
		return image, GetNotFoundErrorFromString(GetNotFoundMessage("Image", id))
	}
	return resp.Images.Image[0], nil
}

func (s *EcsService) DescribeNetworkInterfaceById(instanceId string, eniId string) (networkInterface ecs.NetworkInterfaceSet, err error) {
	req := ecs.CreateDescribeNetworkInterfacesRequest()
	if instanceId != "" {
		req.InstanceId = instanceId
	}
	eniIds := []string{eniId}
	req.NetworkInterfaceId = &eniIds
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeNetworkInterfaces(req)
	})
	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeNetworkInterfacesResponse)
	if resp == nil || len(resp.NetworkInterfaceSets.NetworkInterfaceSet) < 1 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("ECS network interface", eniId))
		return
	}

	return resp.NetworkInterfaceSets.NetworkInterfaceSet[0], nil
}

// WaitForInstance waits for instance to given status
func (s *EcsService) WaitForEcsInstance(instanceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := s.DescribeInstanceById(instanceId)
		if err != nil && !NotFoundError(err) {
			return err
		}
		if instance.Status == string(status) {
			//Sleep one more time for timing issues
			time.Sleep(DefaultIntervalMedium * time.Second)
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("ECS Instance", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
}

// WaitForInstance waits for instance to given status
func (s *EcsService) WaitForEcsDisk(diskId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := s.DescribeDiskById("", diskId)
		if err != nil {
			return err
		}
		if instance.Status == string(status) {
			//Sleep one more time for timing issues
			time.Sleep(DefaultIntervalMedium * time.Second)
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("ECS Disk", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
}

func (s *EcsService) WaitForEcsNetworkInterface(eniId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		eni, err := s.DescribeNetworkInterfaceById("", eniId)
		if err != nil {
			return err
		}
		if eni.Status == string(status) {
			break
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("ECS eni", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *EcsService) QueryPrivateIps(eniId string) ([]string, error) {
	if eni, err := s.DescribeNetworkInterfaceById("", eniId); err != nil {
		return nil, fmt.Errorf("Describe NetworkInterface(%s) failed, %s", eniId, err)
	} else {
		filterIps := make([]string, 0, len(eni.PrivateIpSets.PrivateIpSet))
		for i := range eni.PrivateIpSets.PrivateIpSet {
			if eni.PrivateIpSets.PrivateIpSet[i].Primary {
				continue
			}
			filterIps = append(filterIps, eni.PrivateIpSets.PrivateIpSet[i].PrivateIpAddress)
		}
		return filterIps, nil
	}
}

func (s *EcsService) WaitForVpcAttributesChanged(instanceId, vswitchId, privateIp string) error {
	deadline := time.Now().Add(DefaultTimeout * time.Second)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("Wait for VPC attributes changed timeout")
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		instance, err := s.DescribeInstanceById(instanceId)
		if err != nil {
			return fmt.Errorf("Describe instance(%s) failed, %s", instanceId, err)
		}

		if instance.VpcAttributes.PrivateIpAddress.IpAddress[0] != privateIp {
			continue
		}

		if instance.VpcAttributes.VSwitchId != vswitchId {
			continue
		}

		return nil
	}
}

func (s *EcsService) WaitForPrivateIpsCountChanged(eniId string, count int) error {
	deadline := time.Now().Add(DefaultTimeout * time.Second)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("Wait for private IP addrsses count changed timeout")
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		ips, err := s.QueryPrivateIps(eniId)
		if err != nil {
			return fmt.Errorf("Query private IP failed, %s", err)
		}
		if len(ips) == count {
			return nil
		}
	}
}

func (s *EcsService) WaitForPrivateIpsListChanged(eniId string, ipList []string) error {
	deadline := time.Now().Add(DefaultTimeout * time.Second)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("Wait for private IP addrsses list changed timeout")
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		ips, err := s.QueryPrivateIps(eniId)
		if err != nil {
			return fmt.Errorf("Query private IP failed, %s", err)
		}

		if len(ips) != len(ipList) {
			continue
		}

		diff := false
		for i := range ips {
			exist := false
			for j := range ipList {
				if ips[i] == ipList[j] {
					exist = true
					break
				}
			}
			if !exist {
				diff = true
				break
			}
		}

		if !diff {
			return nil
		}
	}
}

func (s *EcsService) WaitForModifySecurityGroupPolicy(id, target string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSecurityGroupAttribute(id)
		if err != nil {
			return WrapError(err)
		}
		if object.InnerAccessPolicy == target {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InnerAccessPolicy, target, ProviderERROR)
		}
	}
}

func (s *EcsService) AttachKeyPair(keyname string, instanceIds []interface{}) error {
	args := ecs.CreateAttachKeyPairRequest()
	args.KeyPairName = keyname
	args.InstanceIds = convertListToJsonString(instanceIds)
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachKeyPair(args)
		})

		if err != nil {
			if IsExceptedError(err, KeyPairServiceUnavailable) {
				return resource.RetryableError(fmt.Errorf("Attach Key Pair timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Error Attach KeyPair: %#v", err))
		}
		return nil
	})
}

func (s *EcsService) QueryInstanceAllDisks(id string) ([]string, error) {
	args := ecs.CreateDescribeDisksRequest()
	args.InstanceId = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(args)
	})
	if err != nil {
		return nil, fmt.Errorf("describe disk failed, %s\n", err)
	}

	resp, _ := raw.(*ecs.DescribeDisksResponse)
	if resp != nil && len(resp.Disks.Disk) < 1 {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("The specified system disk is not found by instance id %s.", id))
	}

	var ids []string
	for _, disk := range resp.Disks.Disk {
		ids = append(ids, disk.DiskId)
	}
	return ids, nil
}

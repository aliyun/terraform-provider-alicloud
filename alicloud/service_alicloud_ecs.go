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

func (s *EcsService) DescribeInstance(id string) (instance ecs.Instance, err error) {
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

func (s *EcsService) DescribeSecurityGroup(id string) (group ecs.DescribeSecurityGroupAttributeResponse, err error) {
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	request.SecurityGroupId = id

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
	if response.SecurityGroupId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("Security Group", id)), NotFoundMsg, ProviderERROR)
		return
	}

	return *response, nil
}

func (s *EcsService) DescribeSecurityGroupRule(id string) (rule ecs.Permission, err error) {
	parts, err := ParseResourceId(id, 8)
	if err != nil {
		return rule, WrapError(err)
	}
	groupId, direction, ipProtocol, portRange, nicType, cidr_ip, policy := parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6]
	priority, err := strconv.Atoi(parts[7])
	if err != nil {
		return rule, WrapError(err)
	}
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	request.SecurityGroupId = groupId
	request.Direction = direction
	request.NicType = nicType

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(request)
	})
	if err != nil {
		return
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	if response == nil {
		return rule, GetNotFoundErrorFromString(GetNotFoundMessage("Security Group", groupId))
	}

	for _, ru := range response.Permissions.Permission {
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

func (s *EcsService) QueryInstancesWithKeyPair(instanceIdsStr, keyPair string) (instanceIds []string, instances []ecs.Instance, err error) {

	request := ecs.CreateDescribeInstancesRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.InstanceIds = instanceIdsStr
	request.KeyPairName = keyPair
	for {
		raw, e := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(request)
		})
		if e != nil {
			err = WrapErrorf(e, DefaultErrorMsg, keyPair, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return
		}
		addDebug(request.GetActionName(), raw)
		object, _ := raw.(*ecs.DescribeInstancesResponse)
		if len(object.Instances.Instance) < 0 {
			return
		}
		for _, inst := range object.Instances.Instance {
			instanceIds = append(instanceIds, inst.InstanceId)
			instances = append(instances, inst)
		}
		if len(instances) < PageSizeLarge {
			break
		}
		if page, e := getNextpageNumber(request.PageNumber); e != nil {
			err = WrapErrorf(e, DefaultErrorMsg, keyPair, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return
		} else {
			request.PageNumber = page
		}
	}
	return
}

func (s *EcsService) DescribeKeyPair(id string) (keyPair ecs.KeyPair, err error) {
	request := ecs.CreateDescribeKeyPairsRequest()
	request.KeyPairName = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeKeyPairs(request)
	})
	if err != nil {
		return keyPair, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	object, _ := raw.(*ecs.DescribeKeyPairsResponse)
	if len(object.KeyPairs.KeyPair) < 1 || object.KeyPairs.KeyPair[0].KeyPairName != id {
		return keyPair, WrapErrorf(Error(GetNotFoundMessage("KeyPair", id)), NotFoundMsg, ProviderERROR)
	}
	return object.KeyPairs.KeyPair[0], nil

}

func (s *EcsService) DescribeKeyPairAttachment(id string) (keyPair ecs.KeyPair, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return keyPair, WrapError(err)
	}
	keyPairName := parts[0]
	keyPair, err = s.DescribeKeyPair(keyPairName)
	if err != nil {
		return keyPair, WrapError(err)
	}
	if keyPair.KeyPairName != keyPairName {
		err = WrapErrorf(Error(GetNotFoundMessage("KeyPairAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	return keyPair, nil

}

func (s *EcsService) DescribeDisk(id string) (disk ecs.Disk, err error) {
	request := ecs.CreateDescribeDisksRequest()
	request.DiskIds = convertListToJsonString([]interface{}{id})

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(request)
	})
	if err != nil {
		return disk, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ecs.DescribeDisksResponse)
	if len(response.Disks.Disk) < 1 || response.Disks.Disk[0].DiskId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("Disk", id)), NotFoundMsg, ProviderERROR)
		return
	}
	addDebug(request.GetActionName(), raw)
	return response.Disks.Disk[0], nil
}

func (s *EcsService) DescribeDiskAttachment(id string) (disk ecs.Disk, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return disk, WrapError(err)
	}
	disk, err = s.DescribeDisk(parts[0])
	if err != nil {
		return disk, WrapError(err)
	}

	if disk.InstanceId != parts[1] {
		err = WrapErrorf(Error(GetNotFoundMessage("DiskAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	return
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

func (s *EcsService) DescribeNetworkInterface(id string) (networkInterface ecs.NetworkInterfaceSet, err error) {
	request := ecs.CreateDescribeNetworkInterfacesRequest()
	eniIds := []string{id}
	request.NetworkInterfaceId = &eniIds
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeNetworkInterfaces(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw)
	response := raw.(*ecs.DescribeNetworkInterfacesResponse)
	if len(response.NetworkInterfaceSets.NetworkInterfaceSet) < 1 ||
		response.NetworkInterfaceSets.NetworkInterfaceSet[0].NetworkInterfaceId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("NetworkInterface", id)), NotFoundMsg, ProviderERROR)
		return
	}

	return response.NetworkInterfaceSets.NetworkInterfaceSet[0], nil
}

func (s *EcsService) DescribeNetworkInterfaceAttachment(id string) (networkInterface ecs.NetworkInterfaceSet, err error) {
	parts, err := ParseResourceId(id, 2)

	if err != nil {
		return networkInterface, WrapError(err)
	}
	eniId, instanceId := parts[0], parts[1]
	request := ecs.CreateDescribeNetworkInterfacesRequest()
	request.InstanceId = instanceId
	eniIds := []string{eniId}
	request.NetworkInterfaceId = &eniIds
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeNetworkInterfaces(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw)
	response := raw.(*ecs.DescribeNetworkInterfacesResponse)
	if len(response.NetworkInterfaceSets.NetworkInterfaceSet) < 1 ||
		response.NetworkInterfaceSets.NetworkInterfaceSet[0].NetworkInterfaceId != eniId {
		err = WrapErrorf(Error(GetNotFoundMessage("NetworkInterfaceAttachment", id)), NotFoundMsg, ProviderERROR)
		return
	}

	return response.NetworkInterfaceSets.NetworkInterfaceSet[0], nil
}

// WaitForInstance waits for instance to given status
func (s *EcsService) WaitForEcsInstance(instanceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := s.DescribeInstance(instanceId)
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

func (s *EcsService) WaitForDisk(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeDisk(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		// Disk need 3-5 seconds to get ExpiredTime after the status is available
		if object.Status == string(status) && object.ExpiredTime != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}

	}
}

func (s *EcsService) WaitForSecurityGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		_, err := s.DescribeSecurityGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), ProviderERROR)
		}

	}
}

func (s *EcsService) WaitForKeyPair(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		_, err := s.DescribeKeyPair(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), ProviderERROR)
		}

	}
}

func (s *EcsService) WaitForDiskAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDiskAttachment(id)
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

func (s *EcsService) WaitForNetworkInterface(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeNetworkInterface(id)
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
	return nil
}

func (s *EcsService) QueryPrivateIps(eniId string) ([]string, error) {
	if eni, err := s.DescribeNetworkInterface(eniId); err != nil {
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

		instance, err := s.DescribeInstance(instanceId)
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
		object, err := s.DescribeSecurityGroup(id)
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

func (s *EcsService) AttachKeyPair(keyName string, instanceIds []interface{}) error {
	request := ecs.CreateAttachKeyPairRequest()
	request.KeyPairName = keyName
	request.InstanceIds = convertListToJsonString(instanceIds)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachKeyPair(request)
		})
		if err != nil {
			if IsExceptedError(err, KeyPairServiceUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, keyName, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
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

func (s *EcsService) WaitForSnapshot(snapshotId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		if time.Now().After(deadline) {
			return WrapErrorf(Error(GetTimeoutMessage("Snapshot", "timeout")), DefaultTimeoutMsg, snapshotId, "WaitForSnapshot", ProviderERROR)
		}
		snapshot, err := s.DescribeSnapshotById(snapshotId)
		if err != nil {
			if NotFoundError(err) && status == Deleted {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, snapshotId, "WaitForSnapshot", AlibabaCloudSdkGoERROR)
		}
		if snapshot.Status == string(status) {
			return nil
		}
		if snapshot.Status == string(SnapshotCreatingFailed) {
			return WrapErrorf(Error("create snapshot failed"), DefaultErrorMsg, snapshotId, "WaitForSnapshot", ProviderERROR)
		}

		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *EcsService) DescribeSnapshotById(snapshotId string) (*ecs.Snapshot, error) {
	args := ecs.CreateDescribeSnapshotsRequest()
	args.SnapshotIds = fmt.Sprintf("[\"%s\"]", snapshotId)
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSnapshots(args)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, snapshotId, args.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	resp := raw.(*ecs.DescribeSnapshotsResponse)
	if resp == nil || len(resp.Snapshots.Snapshot) < 1 {
		return nil, WrapErrorf(GetNotFoundErrorFromString(GetNotFoundMessage("ECS snapshot", snapshotId)), DefaultErrorMsg, snapshotId, args.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return &resp.Snapshots.Snapshot[0], nil
}

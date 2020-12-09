package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"time"

	"strconv"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type EcsService struct {
	client *connectivity.AliyunClient
}

func (s *EcsService) JudgeRegionValidation(key, region string) error {
	request := ecs.CreateDescribeRegionsRequest()
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeRegions(request)
	})
	if err != nil {
		return fmt.Errorf("DescribeRegions got an error: %#v", err)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
func (s *EcsService) DescribeZone(id string) (zone ecs.Zone, err error) {
	request := ecs.CreateDescribeZonesRequest()
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeZonesResponse)
	if len(response.Zones.Zone) < 1 {
		return zone, WrapError(Error("There is no any availability zone in region %s.", s.client.RegionId))
	}

	zoneIds := []string{}
	for _, z := range response.Zones.Zone {
		if z.ZoneId == id {
			return z, nil
		}
		zoneIds = append(zoneIds, z.ZoneId)
	}
	return zone, WrapError(Error("availability_zone %s not exists in region %s, all zones are %s", id, s.client.RegionId, zoneIds))
}

func (s *EcsService) DescribeZones(d *schema.ResourceData) (zones []ecs.Zone, err error) {
	request := ecs.CreateDescribeZonesRequest()
	request.RegionId = s.client.RegionId
	request.InstanceChargeType = d.Get("instance_charge_type").(string)
	request.SpotStrategy = d.Get("spot_strategy").(string)
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, "alicloud_instance_type_families", request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeZonesResponse)
	if len(response.Zones.Zone) < 1 {
		return zones, WrapError(Error("There is no any availability zone in region %s.", s.client.RegionId))
	}
	if v, ok := d.GetOk("zone_id"); ok {
		zoneIds := []string{}
		for _, z := range response.Zones.Zone {
			if z.ZoneId == v.(string) {
				return []ecs.Zone{z}, nil
			}
			zoneIds = append(zoneIds, z.ZoneId)
		}
		return zones, WrapError(Error("availability_zone %s not exists in region %s, all zones are %s", v.(string), s.client.RegionId, zoneIds))
	} else {
		return response.Zones.Zone, nil
	}
}

func (s *EcsService) DescribeInstance(id string) (instance ecs.Instance, err error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = s.client.RegionId
	request.InstanceIds = convertListToJsonString([]interface{}{id})

	var response *ecs.DescribeInstancesResponse
	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(request)
		})
		if err != nil {
			if IsThrottling(err) {
				wait()
				return resource.RetryableError(err)

			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ecs.DescribeInstancesResponse)
		return nil
	})

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	if len(response.Instances.Instance) < 1 {
		return instance, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	}

	return response.Instances.Instance[0], nil
}

func (s *EcsService) DescribeInstanceAttribute(id string) (instance ecs.DescribeInstanceAttributeResponse, err error) {
	request := ecs.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = id
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstanceAttribute(request)
	})
	if err != nil {
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeInstanceAttributeResponse)
	if response.InstanceId != id {
		return instance, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	}

	return *response, nil
}

func (s *EcsService) DescribeInstanceSystemDisk(id, rg string) (disk ecs.Disk, err error) {
	request := ecs.CreateDescribeDisksRequest()
	request.InstanceId = id
	request.DiskType = string(DiskTypeSystem)
	request.RegionId = s.client.RegionId
	// resource_group_id may cause failure to query the system disk of the instance, because the newly created instance may fail to query through the resource_group_id parameter, so temporarily remove this parameter.
	// request.ResourceGroupId = rg
	var response *ecs.DescribeDisksResponse
	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(request)
		})
		if err != nil {
			if IsThrottling(err) {
				wait()
				return resource.RetryableError(err)

			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ecs.DescribeDisksResponse)
		return nil
	})
	if err != nil {
		return disk, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	if len(response.Disks.Disk) < 1 || response.Disks.Disk[0].InstanceId != id {
		return disk, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	}
	return response.Disks.Disk[0], nil
}

// ResourceAvailable check resource available for zone
func (s *EcsService) ResourceAvailable(zone ecs.Zone, resourceType ResourceType) error {
	for _, res := range zone.AvailableResourceCreation.ResourceTypes {
		if res == string(resourceType) {
			return nil
		}
	}
	return WrapError(Error("%s is not available in %s zone of %s region", resourceType, zone.ZoneId, s.client.Region))
}

func (s *EcsService) DiskAvailable(zone ecs.Zone, diskCategory DiskCategory) error {
	for _, disk := range zone.AvailableDiskCategories.DiskCategories {
		if disk == string(diskCategory) {
			return nil
		}
	}
	return WrapError(Error("%s is not available in %s zone of %s region", diskCategory, zone.ZoneId, s.client.Region))
}

func (s *EcsService) JoinSecurityGroups(instanceId string, securityGroupIds []string) error {
	request := ecs.CreateJoinSecurityGroupRequest()
	request.InstanceId = instanceId
	request.RegionId = s.client.RegionId
	for _, sid := range securityGroupIds {
		request.SecurityGroupId = sid
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.JoinSecurityGroup(request)
		})
		if err != nil && IsExpectedErrors(err, []string{"InvalidInstanceId.AlreadyExists"}) {
			return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func (s *EcsService) LeaveSecurityGroups(instanceId string, securityGroupIds []string) error {
	request := ecs.CreateLeaveSecurityGroupRequest()
	request.InstanceId = instanceId
	request.RegionId = s.client.RegionId
	for _, sid := range securityGroupIds {
		request.SecurityGroupId = sid
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.LeaveSecurityGroup(request)
		})
		if err != nil && IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
			return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func (s *EcsService) DescribeSecurityGroup(id string) (group ecs.DescribeSecurityGroupAttributeResponse, err error) {
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	request.SecurityGroupId = id
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	if response.SecurityGroupId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("Security Group", id)), NotFoundMsg, ProviderERROR, response.RequestId)
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
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSecurityGroupId.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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

	return rule, WrapErrorf(Error(GetNotFoundMessage("Security Group Rule", id)), NotFoundMsg, ProviderERROR, response.RequestId)

}

func (s *EcsService) DescribeAvailableResources(d *schema.ResourceData, meta interface{}, destination DestinationResource) (zoneId string, validZones []ecs.AvailableZone, requestId string, err error) {
	client := meta.(*connectivity.AliyunClient)
	// Before creating resources, check input parameters validity according available zone.
	// If availability zone is nil, it will return all of supported resources in the current.
	request := ecs.CreateDescribeAvailableResourceRequest()
	request.RegionId = s.client.RegionId
	request.DestinationResource = string(destination)
	request.IoOptimized = string(IOOptimized)

	if v, ok := d.GetOk("availability_zone"); ok && strings.TrimSpace(v.(string)) != "" {
		zoneId = strings.TrimSpace(v.(string))
	} else if v, ok := d.GetOk("vswitch_id"); ok && strings.TrimSpace(v.(string)) != "" {
		vpcService := VpcService{s.client}
		if vsw, err := vpcService.DescribeVSwitch(strings.TrimSpace(v.(string))); err == nil {
			zoneId = vsw.ZoneId
		}
	}

	if v, ok := d.GetOk("instance_charge_type"); ok && strings.TrimSpace(v.(string)) != "" {
		request.InstanceChargeType = strings.TrimSpace(v.(string))
	}

	if v, ok := d.GetOk("spot_strategy"); ok && strings.TrimSpace(v.(string)) != "" {
		request.SpotStrategy = strings.TrimSpace(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok && strings.TrimSpace(v.(string)) != "" {
		request.NetworkCategory = strings.TrimSpace(v.(string))
	}

	if v, ok := d.GetOk("is_outdated"); ok && v.(bool) == true {
		request.IoOptimized = string(NoneOptimized)
	}

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeAvailableResource(request)
	})
	if err != nil {
		return "", nil, "", WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeAvailableResourceResponse)
	requestId = response.RequestId

	if len(response.AvailableZones.AvailableZone) < 1 {
		err = WrapError(Error("There are no availability resources in the region: %s. RequestId: %s.", client.RegionId, requestId))
		return
	}

	valid := false
	soldout := false
	var expectedZones []string
	for _, zone := range response.AvailableZones.AvailableZone {
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
			err = WrapError(Error("Availability zone %s status is not available in the region %s. Expected availability zones: %s. RequestId: %s.",
				zoneId, client.RegionId, strings.Join(expectedZones, ", "), requestId))
			return
		}
		if soldout {
			err = WrapError(Error("Availability zone %s status is sold out in the region %s. Expected availability zones: %s. RequestId: %s.",
				zoneId, client.RegionId, strings.Join(expectedZones, ", "), requestId))
			return
		}
	}

	if len(validZones) <= 0 {
		err = WrapError(Error("There is no availability resources in the region %s. Please choose another region. RequestId: %s.", client.RegionId, response.RequestId))
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
		return WrapError(Error("The instance type %s is solded out or is not supported in the zone %s. Expected instance types: %s", targetType, zoneId, strings.Join(expectedInstanceTypes, ", ")))
	}
	return WrapError(Error("The instance type %s is solded out or is not supported in the region %s. Expected instance types: %s", targetType, s.client.RegionId, strings.Join(expectedInstanceTypes, ", ")))
}

func (s *EcsService) QueryInstancesWithKeyPair(instanceIdsStr, keyPair string) (instanceIds []string, instances []ecs.Instance, err error) {

	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = s.client.RegionId
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
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
	request.RegionId = s.client.RegionId
	request.KeyPairName = id
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeKeyPairs(request)
	})
	if err != nil {
		return keyPair, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	object, _ := raw.(*ecs.DescribeKeyPairsResponse)
	if len(object.KeyPairs.KeyPair) < 1 || object.KeyPairs.KeyPair[0].KeyPairName != id {
		return keyPair, WrapErrorf(Error(GetNotFoundMessage("KeyPair", id)), NotFoundMsg, ProviderERROR, object.RequestId)
	}
	return object.KeyPairs.KeyPair[0], nil

}

func (s *EcsService) DescribeKeyPairAttachment(id string) (keyPair ecs.KeyPair, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidKeyPair.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return
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
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(request)
	})
	if err != nil {
		return disk, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ecs.DescribeDisksResponse)
	if len(response.Disks.Disk) < 1 || response.Disks.Disk[0].DiskId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("Disk", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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

	if disk.InstanceId != parts[1] && disk.Status != string(InUse) {
		err = WrapErrorf(Error(GetNotFoundMessage("DiskAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *EcsService) DescribeDisksByType(instanceId string, diskType DiskType) (disk []ecs.Disk, err error) {
	request := ecs.CreateDescribeDisksRequest()
	request.RegionId = s.client.RegionId
	if instanceId != "" {
		request.InstanceId = instanceId
	}
	request.DiskType = string(diskType)

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(request)
	})
	if err != nil {
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeDisksResponse)
	if resp == nil {
		return
	}
	return resp.Disks.Disk, nil
}

func (s *EcsService) DescribeTags(resourceId string, resourceType TagResourceType) (tags []ecs.Tag, err error) {
	request := ecs.CreateDescribeTagsRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = resourceId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeTags(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeTagsResponse)
	if len(response.Tags.Tag) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("Tags", resourceId)), NotFoundMsg, ProviderERROR)
		return
	}

	return response.Tags.Tag, nil
}

func (s *EcsService) DescribeImageById(id string) (image ecs.Image, err error) {
	request := ecs.CreateDescribeImagesRequest()
	request.RegionId = s.client.RegionId
	request.ImageId = id
	request.Status = fmt.Sprintf("%s,%s,%s,%s,%s", "Creating", "Waiting", "Available", "UnAvailable", "CreateFailed")
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeImages(request)
	})
	if err != nil {
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeImagesResponse)
	if resp == nil || len(resp.Images.Image) < 1 {
		return image, GetNotFoundErrorFromString(GetNotFoundMessage("Image", id))
	}
	return resp.Images.Image[0], nil
}

func (s *EcsService) deleteImage(d *schema.ResourceData) error {

	object, err := s.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	request := ecs.CreateDeleteImageRequest()

	if force, ok := d.GetOk("force"); ok {
		request.Force = requests.NewBoolean(force.(bool))
	}
	request.RegionId = s.client.RegionId
	request.ImageId = object.ImageId

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteImage(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutCreate), 3*time.Second, s.ImageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func (s *EcsService) updateImage(d *schema.ResourceData) error {

	d.Partial(true)

	err := setTags(s.client, TagResourceImage, d)
	if err != nil {
		return WrapError(err)
	} else {
		d.SetPartial("tags")
	}

	request := ecs.CreateModifyImageAttributeRequest()
	request.RegionId = s.client.RegionId
	request.ImageId = d.Id()

	if d.HasChange("description") || d.HasChange("name") || d.HasChange("image_name") {
		if description, ok := d.GetOk("description"); ok {
			request.Description = description.(string)
		}
		if imageName, ok := d.GetOk("image_name"); ok {
			request.ImageName = imageName.(string)
		} else {
			if imageName, ok := d.GetOk("name"); ok {
				request.ImageName = imageName.(string)
			}
		}
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyImageAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("name")
		d.SetPartial("image_name")
		d.SetPartial("description")
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return nil
}

func (s *EcsService) DescribeNetworkInterface(id string) (networkInterface ecs.NetworkInterfaceSet, err error) {
	request := ecs.CreateDescribeNetworkInterfacesRequest()
	request.RegionId = s.client.RegionId
	eniIds := []string{id}
	request.NetworkInterfaceId = &eniIds
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeNetworkInterfaces(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeNetworkInterfacesResponse)
	if len(response.NetworkInterfaceSets.NetworkInterfaceSet) < 1 ||
		response.NetworkInterfaceSets.NetworkInterfaceSet[0].NetworkInterfaceId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("NetworkInterface", id)), NotFoundMsg, ProviderERROR, response.RequestId)
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
	request.RegionId = s.client.RegionId
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
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeNetworkInterfacesResponse)
	if len(response.NetworkInterfaceSets.NetworkInterfaceSet) < 1 ||
		response.NetworkInterfaceSets.NetworkInterfaceSet[0].NetworkInterfaceId != eniId {
		err = WrapErrorf(Error(GetNotFoundMessage("NetworkInterfaceAttachment", id)), NotFoundMsg, ProviderERROR, response.RequestId)
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

// WaitForInstance waits for instance to given status
func (s *EcsService) InstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeInstance(id)
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
			return WrapError(Error("Wait for VPC attributes changed timeout"))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		instance, err := s.DescribeInstance(instanceId)
		if err != nil {
			return WrapError(err)
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
	request.RegionId = s.client.RegionId
	request.KeyPairName = keyName
	request.InstanceIds = convertListToJsonString(instanceIds)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachKeyPair(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, keyName, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *EcsService) SnapshotStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSnapshot(id)
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

func (s *EcsService) ImageStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeImageById(id)
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

func (s *EcsService) TaskStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTaskById(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if object.TaskStatus == failState {
				return object, object.TaskStatus, WrapError(Error(FailedToReachTargetStatus, object.TaskStatus))
			}
		}
		return object, object.TaskStatus, nil
	}
}

func (s *EcsService) DescribeTaskById(id string) (task *ecs.DescribeTaskAttributeResponse, err error) {
	request := ecs.CreateDescribeTaskAttributeRequest()
	request.TaskId = id

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeTaskAttribute(request)
	})
	if err != nil {
		return task, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	task, _ = raw.(*ecs.DescribeTaskAttributeResponse)

	if task.TaskId == "" {
		return task, GetNotFoundErrorFromString(GetNotFoundMessage("task", id))
	}
	return task, nil
}

func (s *EcsService) DescribeSnapshot(id string) (*ecs.Snapshot, error) {
	snapshot := &ecs.Snapshot{}
	request := ecs.CreateDescribeSnapshotsRequest()
	request.RegionId = s.client.RegionId
	request.SnapshotIds = fmt.Sprintf("[\"%s\"]", id)
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSnapshots(request)
	})
	if err != nil {
		return snapshot, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeSnapshotsResponse)
	if len(response.Snapshots.Snapshot) != 1 || response.Snapshots.Snapshot[0].SnapshotId != id {
		return snapshot, WrapErrorf(Error(GetNotFoundMessage("Snapshot", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	}
	return &response.Snapshots.Snapshot[0], nil
}

func (s *EcsService) DescribeSnapshotPolicy(id string) (*ecs.AutoSnapshotPolicy, error) {
	policy := &ecs.AutoSnapshotPolicy{}
	request := ecs.CreateDescribeAutoSnapshotPolicyExRequest()
	request.AutoSnapshotPolicyId = id
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeAutoSnapshotPolicyEx(request)
	})
	if err != nil {
		return policy, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response := raw.(*ecs.DescribeAutoSnapshotPolicyExResponse)
	if len(response.AutoSnapshotPolicies.AutoSnapshotPolicy) != 1 ||
		response.AutoSnapshotPolicies.AutoSnapshotPolicy[0].AutoSnapshotPolicyId != id {
		return policy, WrapErrorf(Error(GetNotFoundMessage("SnapshotPolicy", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	}

	return &response.AutoSnapshotPolicies.AutoSnapshotPolicy[0], nil
}

func (s *EcsService) DescribeReservedInstance(id string) (reservedInstance ecs.ReservedInstance, err error) {
	request := ecs.CreateDescribeReservedInstancesRequest()
	var balance = &[]string{id}
	request.ReservedInstanceId = balance
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeReservedInstances(request)
	})
	if err != nil {
		return reservedInstance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeReservedInstancesResponse)

	if len(response.ReservedInstances.ReservedInstance) != 1 ||
		response.ReservedInstances.ReservedInstance[0].ReservedInstanceId != id {
		return reservedInstance, GetNotFoundErrorFromString(GetNotFoundMessage("PurchaseReservedInstance", id))
	}
	return response.ReservedInstances.ReservedInstance[0], nil
}

func (s *EcsService) WaitForReservedInstance(id string, status Status, timeout int) error {
	deadLine := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		reservedInstance, err := s.DescribeReservedInstance(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if reservedInstance.Status == string(status) {
			return nil
		}

		if time.Now().After(deadLine) {
			return WrapErrorf(GetTimeErrorFromString("ECS WaitForSnapshotPolicy"), WaitTimeoutMsg, id, GetFunc(1), timeout, reservedInstance.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *EcsService) WaitForSnapshotPolicy(id string, status Status, timeout int) error {
	deadLine := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		snapshotPolicy, err := s.DescribeSnapshotPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}

		if snapshotPolicy.Status == string(status) {
			return nil
		}

		if time.Now().After(deadLine) {
			return WrapErrorf(GetTimeErrorFromString("ECS WaitForSnapshotPolicy"), WaitTimeoutMsg, id, GetFunc(1), timeout, snapshotPolicy.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *EcsService) DescribeLaunchTemplate(id string) (set ecs.LaunchTemplateSet, err error) {

	request := ecs.CreateDescribeLaunchTemplatesRequest()
	request.RegionId = s.client.RegionId
	request.LaunchTemplateId = &[]string{id}

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeLaunchTemplates(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeLaunchTemplatesResponse)
	if len(response.LaunchTemplateSets.LaunchTemplateSet) != 1 ||
		response.LaunchTemplateSets.LaunchTemplateSet[0].LaunchTemplateId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("LaunchTemplate", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}

	return response.LaunchTemplateSets.LaunchTemplateSet[0], nil

}

func (s *EcsService) DescribeLaunchTemplateVersion(id string, version int) (set ecs.LaunchTemplateVersionSet, err error) {

	request := ecs.CreateDescribeLaunchTemplateVersionsRequest()
	request.RegionId = s.client.RegionId
	request.LaunchTemplateId = id
	request.LaunchTemplateVersion = &[]string{strconv.FormatInt(int64(version), 10)}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeLaunchTemplateVersions(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidLaunchTemplate.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeLaunchTemplateVersionsResponse)
	if len(response.LaunchTemplateVersionSets.LaunchTemplateVersionSet) != 1 ||
		response.LaunchTemplateVersionSets.LaunchTemplateVersionSet[0].LaunchTemplateId != id {
		err = WrapErrorf(Error(GetNotFoundMessage("LaunchTemplateVersion", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}

	return response.LaunchTemplateVersionSets.LaunchTemplateVersionSet[0], nil

}

func (s *EcsService) WaitForLaunchTemplate(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeLaunchTemplate(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.LaunchTemplateId == id && string(status) != string(Deleted) {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *EcsService) DescribeImageShareByImageId(id string) (imageShare *ecs.DescribeImageSharePermissionResponse, err error) {
	request := ecs.CreateDescribeImageSharePermissionRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return imageShare, WrapError(err)
	}
	request.ImageId = parts[0]
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeImageSharePermission(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidImageId.NotFound"}) {
			return imageShare, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return imageShare, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeImageSharePermissionResponse)
	if len(resp.Accounts.Account) == 0 {
		return imageShare, WrapErrorf(Error(GetNotFoundMessage("ModifyImageSharePermission", id)), NotFoundMsg, ProviderERROR, resp.RequestId)
	}
	return resp, nil
}

func (s *EcsService) DescribeAutoProvisioningGroup(id string) (group ecs.AutoProvisioningGroup, err error) {
	request := ecs.CreateDescribeAutoProvisioningGroupsRequest()
	ids := []string{id}
	request.AutoProvisioningGroupId = &ids
	request.RegionId = s.client.RegionId
	raw, e := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeAutoProvisioningGroups(request)
	})
	if e != nil {
		err = WrapErrorf(e, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeAutoProvisioningGroupsResponse)
	for _, v := range response.AutoProvisioningGroups.AutoProvisioningGroup {
		if v.AutoProvisioningGroupId == id {
			return v, nil
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("AutoProvisioningGroup", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	return
}

func (s *EcsService) ListTagResources(resourceId, resourceType string) (object interface{}, err error) {
	conn, err := s.client.NewEcsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTagResources"

	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceId.1": resourceId,
		"ResourceType": resourceType,
	}
	tags := make([]interface{}, 0)

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

func (s *EcsService) tagsToMap(tags []ecs.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ecsTagIgnored(t) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func (s *EcsService) ecsTagIgnored(t ecs.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *EcsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]ecs.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, ecs.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := ecs.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidRegionId.NotFound", "InvalidResourceId.NotFound", "InvalidResourceType.NotFound", "MissingParameter.RegionId", "MissingParameter.ResourceIds", "MissingParameter.ResourceType", "MissingParameter.TagOwnerBid", "MissingParameter.TagOwnerUid", "MissingParameter.Tags"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := ecs.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidRegionId.NotFound", "InvalidResourceId.NotFound", "InvalidResourceType.NotFound", "MissingParameter.RegionId", "MissingParameter.ResourceIds", "MissingParameter.ResourceType", "MissingParameter.TagOwnerBid", "MissingParameter.TagOwnerUid", "MissingParameter.Tags"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *EcsService) DescribeEcsDedicatedHost(id string) (object ecs.DedicatedHost, err error) {
	request := ecs.CreateDescribeDedicatedHostsRequest()
	request.RegionId = s.client.RegionId

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	response := new(ecs.DescribeDedicatedHostsResponse)
	for {

		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDedicatedHosts(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidLockReason.NotFound"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("EcsDedicatedHost", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ecs.DescribeDedicatedHostsResponse)

		if len(response.DedicatedHosts.DedicatedHost) < 1 {
			err = WrapErrorf(Error(GetNotFoundMessage("EcsDedicatedHost", id)), NotFoundMsg, ProviderERROR, response.RequestId)
			return object, err
		}
		for _, object := range response.DedicatedHosts.DedicatedHost {
			if object.DedicatedHostId == id {
				return object, nil
			}
		}
		if len(response.DedicatedHosts.DedicatedHost) < PageSizeMedium {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return object, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("EcsDedicatedHost", id)), NotFoundMsg, ProviderERROR, response.RequestId)
	return
}

func (s *EcsService) EcsDedicatedHostStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcsDedicatedHost(id)
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

func (s *EcsService) WaitForAutoProvisioningGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAutoProvisioningGroup(id)
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
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AutoProvisioningGroupId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

package alicloud

import (
	"fmt"
	"strings"

	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

// DescribeZone validate zoneId is valid in region
func (client *AliyunClient) DescribeZone(zoneID string) (zone ecs.Zone, err error) {
	resp, err := client.ecsconn.DescribeZones(ecs.CreateDescribeZonesRequest())
	if err != nil {
		return
	}
	if resp == nil || len(resp.Zones.Zone) < 1 {
		return zone, fmt.Errorf("There is no any availability zone in region %s.", client.RegionId)
	}

	zoneIds := []string{}
	for _, z := range resp.Zones.Zone {
		if z.ZoneId == zoneID {
			return z, nil
		}
		zoneIds = append(zoneIds, z.ZoneId)
	}
	return zone, fmt.Errorf("availability_zone not exists in range %s, all zones are %s", client.RegionId, zoneIds)
}

func (client *AliyunClient) DescribeInstanceById(id string) (instance ecs.Instance, err error) {
	req := ecs.CreateDescribeInstancesRequest()
	req.InstanceIds = convertListToJsonString([]interface{}{id})

	resp, err := client.ecsconn.DescribeInstances(req)
	if err != nil {
		return
	}
	if resp == nil || len(resp.Instances.Instance) < 1 {
		return instance, GetNotFoundErrorFromString(GetNotFoundMessage("Instance", id))
	}

	return resp.Instances.Instance[0], nil
}

func (client *AliyunClient) QueryInstanceSystemDisk(id string) (disk ecs.Disk, err error) {
	args := ecs.CreateDescribeDisksRequest()
	args.InstanceId = id
	args.DiskType = string(DiskTypeSystem)

	resp, err := client.ecsconn.DescribeDisks(args)
	if err != nil {
		return
	}
	if resp != nil && len(resp.Disks.Disk) < 1 {
		return disk, GetNotFoundErrorFromString(fmt.Sprintf("The specified system disk is not found by instance id %s.", id))
	}

	return resp.Disks.Disk[0], nil
}

// ResourceAvailable check resource available for zone
func (client *AliyunClient) ResourceAvailable(zone ecs.Zone, resourceType ResourceType) error {
	for _, res := range zone.AvailableResourceCreation.ResourceTypes {
		if res == string(resourceType) {
			return nil
		}
	}
	return fmt.Errorf("%s is not available in %s zone of %s region", resourceType, zone.ZoneId, client.Region)
}

func (client *AliyunClient) DiskAvailable(zone ecs.Zone, diskCategory DiskCategory) error {
	for _, disk := range zone.AvailableDiskCategories.DiskCategories {
		if disk == string(diskCategory) {
			return nil
		}
	}
	return fmt.Errorf("%s is not available in %s zone of %s region", diskCategory, zone.ZoneId, client.Region)
}

func (client *AliyunClient) JoinSecurityGroups(instanceId string, securityGroupIds []string) error {
	req := ecs.CreateJoinSecurityGroupRequest()
	req.InstanceId = instanceId
	for _, sid := range securityGroupIds {
		req.SecurityGroupId = sid
		_, err := client.ecsconn.JoinSecurityGroup(req)
		if err != nil && IsExceptedErrors(err, []string{InvalidInstanceIdAlreadyExists}) {
			return err
		}
	}

	return nil
}

func (client *AliyunClient) LeaveSecurityGroups(instanceId string, securityGroupIds []string) error {
	req := ecs.CreateLeaveSecurityGroupRequest()
	req.InstanceId = instanceId
	for _, sid := range securityGroupIds {
		req.SecurityGroupId = sid
		_, err := client.ecsconn.LeaveSecurityGroup(req)
		if err != nil && IsExceptedErrors(err, []string{InvalidSecurityGroupIdNotFound}) {
			return err
		}
	}

	return nil
}

func (client *AliyunClient) DescribeSecurityGroupAttribute(securityGroupId string) (group ecs.DescribeSecurityGroupAttributeResponse, err error) {
	args := ecs.CreateDescribeSecurityGroupAttributeRequest()
	args.SecurityGroupId = securityGroupId

	resp, err := client.ecsconn.DescribeSecurityGroupAttribute(args)
	if err != nil {
		return
	}
	if resp == nil {
		return group, GetNotFoundErrorFromString(GetNotFoundMessage("Security Group", securityGroupId))
	}

	return *resp, nil
}

func (client *AliyunClient) DescribeSecurityGroupRule(groupId, direction, ipProtocol, portRange, nicType, cidr_ip, policy string, priority int) (rule ecs.Permission, err error) {
	args := ecs.CreateDescribeSecurityGroupAttributeRequest()
	args.SecurityGroupId = groupId
	args.Direction = direction
	args.NicType = nicType

	resp, err := client.ecsconn.DescribeSecurityGroupAttribute(args)
	if err != nil {
		return
	}
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

func (client *AliyunClient) DescribeAvailableResources(d *schema.ResourceData, meta interface{}, destination DestinationResource) (zoneId string, validZones []ecs.AvailableZone, err error) {
	// Before creating resources, check input parameters validity according available zone.
	// If availability zone is nil, it will return all of supported resources in the current.
	conn := meta.(*AliyunClient).ecsconn
	args := ecs.CreateDescribeAvailableResourceRequest()
	args.DestinationResource = string(destination)
	args.IoOptimized = string(IOOptimized)

	if v, ok := d.GetOk("availability_zone"); ok && strings.TrimSpace(v.(string)) != "" {
		zoneId = strings.TrimSpace(v.(string))
	} else if v, ok := d.GetOk("vswitch_id"); ok && strings.TrimSpace(v.(string)) != "" {
		if vsw, err := meta.(*AliyunClient).DescribeVswitch(strings.TrimSpace(v.(string))); err == nil {
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

	resources, e := conn.DescribeAvailableResource(args)
	if e != nil {
		return "", nil, fmt.Errorf("Error DescribeAvailableResource: %#v", e)
	}

	if resources == nil || len(resources.AvailableZones.AvailableZone) < 1 {
		err = fmt.Errorf("There are no availability resources in the region: %s.", getRegionId(d, meta))
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
				zoneId, getRegionId(d, meta), strings.Join(expectedZones, ", "))
			return
		}
		if soldout {
			err = fmt.Errorf("Availability zone %s status is sold out in the region %s. Expected availability zones: %s.",
				zoneId, getRegionId(d, meta), strings.Join(expectedZones, ", "))
			return
		}
	}

	if len(validZones) <= 0 {
		err = fmt.Errorf("There is no availability resources in the region %s. Please choose another region.", getRegionId(d, meta))
		return
	}

	return
}

func (client *AliyunClient) InstanceTypeValidation(targetType, zoneId string, validZones []ecs.AvailableZone) error {

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
	return fmt.Errorf("The instance type %s is solded out or is not supported in the region %s. Expected instance types: %s", targetType, client.RegionId, strings.Join(expectedInstanceTypes, ", "))
}

func (client *AliyunClient) QueryInstancesWithKeyPair(instanceIdsStr, keypair string) (instanceIds []string, instances []ecs.Instance, err error) {

	conn := client.ecsconn
	args := ecs.CreateDescribeInstancesRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	args.InstanceIds = instanceIdsStr
	args.KeyPairName = keypair
	for true {
		resp, e := conn.DescribeInstances(args)
		if e != nil {
			err = e
			return
		}
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
		args.PageNumber = args.PageNumber + requests.NewInteger(1)
	}
	return
}

func (client *AliyunClient) DescribeKeyPair(keyName string) (keypair ecs.KeyPair, err error) {
	req := ecs.CreateDescribeKeyPairsRequest()
	req.KeyPairName = keyName
	resp, err := client.ecsconn.DescribeKeyPairs(req)

	if err != nil {
		return
	}

	if resp == nil || len(resp.KeyPairs.KeyPair) < 1 {
		return keypair, GetNotFoundErrorFromString(GetNotFoundMessage("KeyPair", keyName))
	}
	return resp.KeyPairs.KeyPair[0], nil

}

func (client *AliyunClient) DescribeDiskById(instanceId, diskId string) (disk ecs.Disk, err error) {
	req := ecs.CreateDescribeDisksRequest()
	if instanceId != "" {
		req.InstanceId = instanceId
	}
	req.DiskIds = convertListToJsonString([]interface{}{diskId})

	resp, err := client.ecsconn.DescribeDisks(req)
	if err != nil {
		return
	}
	if resp == nil || len(resp.Disks.Disk) < 1 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("ECS disk", diskId))
		return
	}
	return resp.Disks.Disk[0], nil
}

func (client *AliyunClient) DescribeDisksByType(instanceId string, diskType DiskType) (disk []ecs.Disk, err error) {
	req := ecs.CreateDescribeDisksRequest()
	if instanceId != "" {
		req.InstanceId = instanceId
	}
	req.DiskType = string(diskType)

	resp, err := client.ecsconn.DescribeDisks(req)
	if err != nil {
		return
	}
	if resp == nil {
		return
	}
	return resp.Disks.Disk, nil
}

func (client *AliyunClient) DescribeTags(resourceId string, resourceType TagResourceType) (tags []ecs.Tag, err error) {
	req := ecs.CreateDescribeTagsRequest()
	req.ResourceType = string(resourceType)
	req.ResourceId = resourceId
	resp, err := client.ecsconn.DescribeTags(req)

	if err != nil {
		return
	}
	if resp == nil || len(resp.Tags.Tag) < 1 {
		err = GetNotFoundErrorFromString(fmt.Sprintf("Describe %s tag by id %s got an error.", resourceType, resourceId))
		return
	}

	return resp.Tags.Tag, nil
}

func (client *AliyunClient) DescribeImageById(id string) (image ecs.Image, err error) {
	req := ecs.CreateDescribeImagesRequest()
	req.ImageId = id
	resp, err := client.ecsconn.DescribeImages(req)
	if err != nil {
		return
	}
	if resp == nil || len(resp.Images.Image) < 1 {
		return image, GetNotFoundErrorFromString(GetNotFoundMessage("Image", id))
	}
	return resp.Images.Image[0], nil
}

// WaitForInstance waits for instance to given status
func (client *AliyunClient) WaitForEcsInstance(instanceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := client.DescribeInstanceById(instanceId)
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
			return GetTimeErrorFromString(GetTimeoutMessage("ECS Instance", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
}

// WaitForInstance waits for instance to given status
func (client *AliyunClient) WaitForEcsDisk(diskId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := client.DescribeDiskById("", diskId)
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

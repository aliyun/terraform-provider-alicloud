package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func (client *AliyunClient) DescribeImage(imageId string) (*ecs.ImageType, error) {

	pagination := common.Pagination{
		PageNumber: 1,
	}
	args := ecs.DescribeImagesArgs{
		Pagination: pagination,
		RegionId:   client.Region,
		Status:     ecs.ImageStatusAvailable,
	}

	var allImages []ecs.ImageType

	for {
		images, _, err := client.ecsconn.DescribeImages(&args)
		if err != nil {
			break
		}

		if len(images) == 0 {
			break
		}

		allImages = append(allImages, images...)

		args.Pagination.PageNumber++
	}

	if len(allImages) == 0 {
		return nil, common.GetClientErrorFromString("Not found")
	}

	var image *ecs.ImageType
	imageIds := []string{}
	for _, im := range allImages {
		if im.ImageId == imageId {
			image = &im
		}
		imageIds = append(imageIds, im.ImageId)
	}

	if image == nil {
		return nil, fmt.Errorf("image_id %s not exists in range %s, all images are %s", imageId, client.Region, imageIds)
	}

	return image, nil
}

// DescribeZone validate zoneId is valid in region
func (client *AliyunClient) DescribeZone(zoneID string) (*ecs.ZoneType, error) {
	zones, err := client.ecsconn.DescribeZones(client.Region)
	if err != nil {
		return nil, fmt.Errorf("error to list zones not found")
	}

	var zone *ecs.ZoneType
	zoneIds := []string{}
	for _, z := range zones {
		if z.ZoneId == zoneID {
			zone = &ecs.ZoneType{
				ZoneId:                    z.ZoneId,
				LocalName:                 z.LocalName,
				AvailableResourceCreation: z.AvailableResourceCreation,
				AvailableDiskCategories:   z.AvailableDiskCategories,
			}
		}
		zoneIds = append(zoneIds, z.ZoneId)
	}

	if zone == nil {
		return nil, fmt.Errorf("availability_zone not exists in range %s, all zones are %s", client.Region, zoneIds)
	}

	return zone, nil
}

func (client *AliyunClient) QueryInstancesByIds(ids []string) (instances []ecs.InstanceAttributesType, err error) {
	idsStr, jerr := json.Marshal(ids)
	if jerr != nil {
		return nil, jerr
	}

	args := ecs.DescribeInstancesArgs{
		RegionId:    client.Region,
		InstanceIds: string(idsStr),
	}

	instances, _, errs := client.ecsconn.DescribeInstances(&args)

	if errs != nil {
		return nil, errs
	}

	return instances, nil
}

func (client *AliyunClient) QueryInstancesById(id string) (instance *ecs.InstanceAttributesType, err error) {
	ids := []string{id}

	instances, errs := client.QueryInstancesByIds(ids)
	if errs != nil {
		return nil, errs
	}

	if len(instances) == 0 {
		return nil, GetNotFoundErrorFromString(InstanceNotFound)
	}

	return &instances[0], nil
}

func (client *AliyunClient) QueryInstanceSystemDisk(id string) (disk *ecs.DiskItemType, err error) {
	args := ecs.DescribeDisksArgs{
		RegionId:   client.Region,
		InstanceId: string(id),
		DiskType:   ecs.DiskTypeAllSystem,
	}
	disks, _, err := client.ecsconn.DescribeDisks(&args)
	if err != nil {
		return nil, err
	}
	if len(disks) == 0 {
		return nil, GetNotFoundErrorFromString(SystemDiskNotFound)
	}

	return &disks[0], nil
}

// ResourceAvailable check resource available for zone
func (client *AliyunClient) ResourceAvailable(zone *ecs.ZoneType, resourceType ecs.ResourceType) error {
	available := false
	for _, res := range zone.AvailableResourceCreation.ResourceTypes {
		if res == resourceType {
			available = true
		}
	}
	if !available {
		return fmt.Errorf("%s is not available in %s zone of %s region", resourceType, zone.ZoneId, client.Region)
	}

	return nil
}

func (client *AliyunClient) DiskAvailable(zone *ecs.ZoneType, diskCategory ecs.DiskCategory) error {
	available := false
	for _, dist := range zone.AvailableDiskCategories.DiskCategories {
		if dist == diskCategory {
			available = true
		}
	}
	if !available {
		return fmt.Errorf("%s is not available in %s zone of %s region", diskCategory, zone.ZoneId, client.Region)
	}
	return nil
}

// todo: support syc
func (client *AliyunClient) JoinSecurityGroups(instanceId string, securityGroupIds []string) error {
	for _, sid := range securityGroupIds {
		err := client.ecsconn.JoinSecurityGroup(instanceId, sid)
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code != InvalidInstanceIdAlreadyExists {
				return err
			}
		}
	}

	return nil
}

func (client *AliyunClient) LeaveSecurityGroups(instanceId string, securityGroupIds []string) error {
	for _, sid := range securityGroupIds {
		err := client.ecsconn.LeaveSecurityGroup(instanceId, sid)
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code != InvalidSecurityGroupIdNotFound {
				return err
			}
		}
	}

	return nil
}

func (client *AliyunClient) DescribeSecurity(securityGroupId string) (*ecs.DescribeSecurityGroupAttributeResponse, error) {

	args := &ecs.DescribeSecurityGroupAttributeArgs{
		RegionId:        client.Region,
		SecurityGroupId: securityGroupId,
	}

	return client.ecsconn.DescribeSecurityGroupAttribute(args)
}

func (client *AliyunClient) DescribeSecurityGroupRule(groupId, direction, ipProtocol, portRange, nicType, cidr_ip, policy string, priority int) (*ecs.PermissionType, error) {
	rules, err := client.ecsconn.DescribeSecurityGroupAttribute(&ecs.DescribeSecurityGroupAttributeArgs{
		RegionId:        client.Region,
		SecurityGroupId: groupId,
		Direction:       ecs.Direction(direction),
		NicType:         ecs.NicType(nicType),
	})

	if err != nil {
		return nil, err
	}

	for _, ru := range rules.Permissions.Permission {
		if strings.ToLower(string(ru.IpProtocol)) == ipProtocol && ru.PortRange == portRange {
			cidr := ru.SourceCidrIp
			if ecs.Direction(direction) == ecs.DirectionIngress && cidr == "" {
				cidr = ru.SourceGroupId
			}
			if ecs.Direction(direction) == ecs.DirectionEgress {
				if cidr = ru.DestCidrIp; cidr == "" {
					cidr = ru.DestGroupId
				}
			}

			if cidr == cidr_ip && strings.ToLower(string(ru.Policy)) == policy && ru.Priority == priority {
				return &ru, nil
			}
		}
	}

	return nil, GetNotFoundErrorFromString("Security group rule not found")

}

func (client *AliyunClient) RevokeSecurityGroup(args *ecs.RevokeSecurityGroupArgs) error {
	//when the rule is not exist, api will return success(200)
	return client.ecsconn.RevokeSecurityGroup(args)
}

func (client *AliyunClient) RevokeSecurityGroupEgress(args *ecs.RevokeSecurityGroupEgressArgs) error {
	//when the rule is not exist, api will return success(200)
	return client.ecsconn.RevokeSecurityGroupEgress(args)
}

func (client *AliyunClient) DescribeAvailableResources(d *schema.ResourceData, meta interface{}, destination DestinationResource) (zoneId string, validZones []ecs.AvailableZoneType, err error) {
	// Before creating resources, check input parameters validity according available zone.
	// If availability zone is nil, it will return all of supported resources in the current.
	conn := meta.(*AliyunClient).ecsconn
	args := ecs.DescribeAvailableResourceArgs{
		RegionId:            getRegionId(d, meta),
		DestinationResource: string(destination),
		IoOptimized:         string(IOOptimized),
	}
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

	resources, e := conn.DescribeAvailableResource(&args)
	if e != nil {
		return "", nil, fmt.Errorf("Error DescribeAvailableResource: %#v", e)
	}

	if resources == nil || len(resources.AvailableZones.AvailableZone) < 1 {
		err = fmt.Errorf("There are no availability resources in the region: %s.", getRegionId(d, meta))
		return
	}

	valid := false
	var expectedZones []string
	for _, zone := range resources.AvailableZones.AvailableZone {
		if zone.Status == string(SoldOut) {
			continue
		}
		if zoneId != "" && zone.ZoneId == zoneId {
			valid = true
			validZones = append(make([]ecs.AvailableZoneType, 1), zone)
			break
		}
		expectedZones = append(expectedZones, zone.ZoneId)
		validZones = append(validZones, zone)
	}
	if zoneId != "" && !valid {
		err = fmt.Errorf("Availability zone %s status is sold out in the region %s. Expected availability zones: %s.",
			zoneId, getRegionId(d, meta), strings.Join(expectedZones, ", "))
		return
	}

	if len(validZones) <= 0 {
		err = fmt.Errorf("There is no availability resources in the region %s. Please choose another region.", getRegionId(d, meta))
		return
	}

	return
}

func (client *AliyunClient) InstanceTypeValidation(targetType, zoneId string, validZones []ecs.AvailableZoneType) error {

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

func (client *AliyunClient) QueryInstancesWithKeyPair(region common.Region, instanceIds, keypair string) ([]interface{}, []ecs.InstanceAttributesType, error) {
	var instance_ids []interface{}
	var instanceList []ecs.InstanceAttributesType

	conn := client.ecsconn
	args := &ecs.DescribeInstancesArgs{
		RegionId: region,
	}
	pagination := getPagination(1, 50)
	for true {
		if instanceIds != "" {
			args.InstanceIds = instanceIds
		}
		args.Pagination = pagination
		instances, _, err := conn.DescribeInstances(args)
		if err != nil {
			return nil, nil, fmt.Errorf("Error DescribeInstances: %#v", err)
		}
		for _, inst := range instances {
			if inst.KeyPairName == keypair {
				instance_ids = append(instance_ids, inst.InstanceId)
				instanceList = append(instanceList, inst)
			}
		}
		if len(instances) < pagination.PageSize {
			break
		}
		pagination.PageNumber += 1
	}
	return instance_ids, instanceList, nil
}

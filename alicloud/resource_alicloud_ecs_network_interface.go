package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudEcsNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsNetworkInterfaceCreate,
		Read:   resourceAliCloudEcsNetworkInterfaceRead,
		Update: resourceAliCloudEcsNetworkInterfaceUpdate,
		Delete: resourceAliCloudEcsNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.123.1. New field 'network_interface_name' instead",
				ConflictsWith: []string{"network_interface_name"},
			},
			"primary_ip_address": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"private_ip"},
			},
			"private_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'private_ip' has been deprecated from provider version 1.123.1. New field 'primary_ip_address' instead",
				ConflictsWith: []string{"primary_ip_address"},
			},
			"private_ip_addresses": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"private_ips", "secondary_private_ip_address_count", "private_ips_count", "ipv4_prefixes", "ipv4_prefix_count"},
			},
			"private_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Deprecated:    "Field 'private_ips' has been deprecated from provider version 1.123.1. New field 'private_ip_addresses' instead",
				ConflictsWith: []string{"private_ip_addresses", "secondary_private_ip_address_count", "private_ips_count", "ipv4_prefixes", "ipv4_prefix_count"},
			},
			"queue_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secondary_private_ip_address_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"private_ips_count", "private_ip_addresses", "private_ips", "ipv4_prefixes", "ipv4_prefix_count"},
			},
			"private_ips_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'private_ips_count' has been deprecated from provider version 1.123.1. New field 'secondary_private_ip_address_count' instead",
				ConflictsWith: []string{"secondary_private_ip_address_count", "private_ip_addresses", "private_ips", "ipv4_prefixes", "ipv4_prefix_count"},
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems:      1,
				ConflictsWith: []string{"security_groups"},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems:      1,
				Deprecated:    "Field 'security_groups' has been deprecated from provider version 1.123.1. New field 'security_group_ids' instead",
				ConflictsWith: []string{"security_group_ids"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv6_address_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.IntBetween(1, 10),
				ConflictsWith: []string{"ipv6_addresses"},
			},
			"ipv6_addresses": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				MaxItems:      10,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"ipv6_address_count"},
			},
			"ipv4_prefix_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.IntBetween(1, 10),
				ConflictsWith: []string{"ipv4_prefixes", "private_ip_addresses", "secondary_private_ip_address_count", "private_ip_addresses", "private_ips_count"},
			},
			"ipv4_prefixes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				MaxItems:      10,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"ipv4_prefix_count", "private_ip_addresses", "secondary_private_ip_address_count", "private_ip_addresses", "private_ips_count"},
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"network_interface_traffic_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEcsNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "CreateNetworkInterface"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("network_interface_name"); ok {
		request["NetworkInterfaceName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["NetworkInterfaceName"] = v
	}

	if v, ok := d.GetOk("primary_ip_address"); ok {
		request["PrimaryIpAddress"] = v
	} else if v, ok := d.GetOk("private_ip"); ok {
		request["PrimaryIpAddress"] = v
	}

	if v, ok := d.GetOk("private_ip_addresses"); ok {
		request["PrivateIpAddress"] = v.(*schema.Set).List()
	} else if v, ok := d.GetOk("private_ips"); ok {
		request["PrivateIpAddress"] = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("queue_number"); ok {
		request["QueueNumber"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("secondary_private_ip_address_count"); ok {
		request["SecondaryPrivateIpAddressCount"] = v
	} else if v, ok := d.GetOk("private_ips_count"); ok {
		request["SecondaryPrivateIpAddressCount"] = v
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		request["SecurityGroupIds"] = v
	} else if v, ok := d.GetOk("security_groups"); ok {
		request["SecurityGroupIds"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	request["VSwitchId"] = d.Get("vswitch_id")
	if v, ok := d.GetOk("ipv6_addresses"); ok {
		request["Ipv6Address"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOkExists("ipv6_address_count"); ok {
		request["Ipv6AddressCount"] = v
	}

	if v, ok := d.GetOk("ipv4_prefixes"); ok {
		request["Ipv4Prefix"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOkExists("ipv4_prefix_count"); ok {
		request["Ipv4PrefixCount"] = v
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}

	if v, ok := d.GetOk("network_interface_traffic_mode"); ok {
		request["NetworkInterfaceTrafficMode"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"IncorrectVSwitchStatus"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_network_interface", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["NetworkInterfaceId"]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsNetworkInterfaceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcsNetworkInterfaceRead(d, meta)
}

func resourceAliCloudEcsNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsNetworkInterface(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_network_interface ecsService.DescribeEcsNetworkInterface Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("mac", object["MacAddress"])
	d.Set("network_interface_name", object["NetworkInterfaceName"])
	d.Set("name", object["NetworkInterfaceName"])
	d.Set("primary_ip_address", object["PrivateIpAddress"])
	d.Set("private_ip", object["PrivateIpAddress"])
	privateIps := make([]interface{}, 0, len(object["PrivateIpSets"].(map[string]interface{})["PrivateIpSet"].([]interface{})))
	for _, v := range object["PrivateIpSets"].(map[string]interface{})["PrivateIpSet"].([]interface{}) {
		if !v.(map[string]interface{})["Primary"].(bool) {
			privateIps = append(privateIps, v.(map[string]interface{})["PrivateIpAddress"])
		}
	}
	d.Set("private_ips", privateIps)
	d.Set("private_ip_addresses", privateIps)
	d.Set("private_ips_count", len(privateIps))
	d.Set("secondary_private_ip_address_count", len(privateIps))
	ipv6SetList := make([]interface{}, 0, len(object["Ipv6Sets"].(map[string]interface{})["Ipv6Set"].([]interface{})))
	for _, v := range object["Ipv6Sets"].(map[string]interface{})["Ipv6Set"].([]interface{}) {
		ipv6Set := v.(map[string]interface{})
		ipv6SetList = append(ipv6SetList, ipv6Set["Ipv6Address"])
	}
	d.Set("ipv6_addresses", ipv6SetList)
	d.Set("ipv6_address_count", len(ipv6SetList))
	ipv4PrefixSetList := make([]interface{}, 0, len(object["Ipv4PrefixSets"].(map[string]interface{})["Ipv4PrefixSet"].([]interface{})))
	for _, v := range object["Ipv4PrefixSets"].(map[string]interface{})["Ipv4PrefixSet"].([]interface{}) {
		ipv4PrefixSet := v.(map[string]interface{})
		ipv4PrefixSetList = append(ipv4PrefixSetList, ipv4PrefixSet["Ipv4Prefix"])
	}
	d.Set("ipv4_prefixes", ipv4PrefixSetList)
	d.Set("ipv4_prefix_count", len(ipv4PrefixSetList))
	d.Set("queue_number", formatInt(object["QueueNumber"]))
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("security_group_ids", object["SecurityGroupIds"].(map[string]interface{})["SecurityGroupId"])
	d.Set("security_groups", object["SecurityGroupIds"].(map[string]interface{})["SecurityGroupId"])
	d.Set("status", object["Status"])

	tags, err := ecsService.ListTagResources(d.Id(), "eni")
	if err != nil {
		return WrapError(err)
	} else {
		d.Set("tags", tagsToMap(tags))
	}

	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("instance_type", object["Type"])
	d.Set("network_interface_traffic_mode", object["NetworkInterfaceTrafficMode"])

	return nil
}

func resourceAliCloudEcsNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	ecsService := EcsService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "eni"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	updateSecurityGroup := false
	request := map[string]interface{}{
		"NetworkInterfaceId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("network_interface_name") {
		update = true
		request["NetworkInterfaceName"] = d.Get("network_interface_name")
	}
	if d.HasChange("name") {
		update = true
		request["NetworkInterfaceName"] = d.Get("name")
	}
	if d.HasChange("queue_number") {
		update = true
		request["QueueNumber"] = d.Get("queue_number")
	}
	if d.HasChange("security_group_ids") {
		update = true
		request["SecurityGroupId"] = d.Get("security_group_ids")
		updateSecurityGroup = true
	}
	if d.HasChange("security_groups") {
		update = true
		request["SecurityGroupId"] = d.Get("security_groups")
		updateSecurityGroup = true
	}
	if update {
		action := "ModifyNetworkInterfaceAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
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
		d.SetPartial("description")
		d.SetPartial("name")
		d.SetPartial("network_interface_name")
		d.SetPartial("queue_number")
		d.SetPartial("security_groups")
		d.SetPartial("security_group_ids")
		if updateSecurityGroup {
			time.Sleep(500 * time.Millisecond)
		}
	}
	d.Partial(false)
	if d.HasChange("private_ip_addresses") {
		oldPrivateIpAddresses, newPrivateIpAddresses := d.GetChange("private_ip_addresses")
		oldPrivateIpAddressesSet := oldPrivateIpAddresses.(*schema.Set)
		newPrivateIpAddressesSet := newPrivateIpAddresses.(*schema.Set)

		removed := oldPrivateIpAddressesSet.Difference(newPrivateIpAddressesSet)
		added := newPrivateIpAddressesSet.Difference(oldPrivateIpAddressesSet)
		if removed.Len() > 0 {
			unassignprivateipaddressesrequest := map[string]interface{}{
				"NetworkInterfaceId": d.Id(),
			}
			unassignprivateipaddressesrequest["PrivateIpAddress"] = removed.List()
			unassignprivateipaddressesrequest["RegionId"] = client.RegionId
			action := "UnassignPrivateIpAddresses"
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, unassignprivateipaddressesrequest, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"InternalError", "InvalidOperation.InvalidEcsState", "InvalidOperation.InvalidEniState", "OperationConflict", "ServiceUnavailable"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, unassignprivateipaddressesrequest)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			d.SetPartial("private_ip_addresses")
		}
		if added.Len() > 0 {
			assignprivateipaddressesrequest := map[string]interface{}{
				"NetworkInterfaceId": d.Id(),
			}
			assignprivateipaddressesrequest["PrivateIpAddress"] = added.List()
			assignprivateipaddressesrequest["RegionId"] = client.RegionId
			action := "AssignPrivateIpAddresses"
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, assignprivateipaddressesrequest, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"InternalError", "InvalidOperation.InvalidEcsState", "InvalidOperation.InvalidEniState", "OperationConflict", "ServiceUnavailable"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, assignprivateipaddressesrequest)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			d.SetPartial("private_ip_addresses")
		}
		if err := ecsService.WaitForPrivateIpsListChanged(d.Id(), expandStringList(newPrivateIpAddressesSet.List())); err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("private_ips") {
		oldPrivateIps, newPrivateIps := d.GetChange("private_ips")
		oldPrivateIpsSet := oldPrivateIps.(*schema.Set)
		newPrivateIpsSet := newPrivateIps.(*schema.Set)

		removed := oldPrivateIpsSet.Difference(newPrivateIpsSet)
		added := newPrivateIpsSet.Difference(oldPrivateIpsSet)
		if removed.Len() > 0 {
			unassignprivateipaddressesrequest := map[string]interface{}{
				"NetworkInterfaceId": d.Id(),
			}
			unassignprivateipaddressesrequest["PrivateIpAddress"] = removed.List()
			unassignprivateipaddressesrequest["RegionId"] = client.RegionId
			action := "UnassignPrivateIpAddresses"
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, unassignprivateipaddressesrequest, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"InternalError", "InvalidOperation.InvalidEcsState", "InvalidOperation.InvalidEniState", "OperationConflict", "ServiceUnavailable"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, unassignprivateipaddressesrequest)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			d.SetPartial("private_ips")
		}
		if added.Len() > 0 {
			assignprivateipaddressesrequest := map[string]interface{}{
				"NetworkInterfaceId": d.Id(),
			}
			assignprivateipaddressesrequest["PrivateIpAddress"] = added.List()
			assignprivateipaddressesrequest["RegionId"] = client.RegionId
			action := "AssignPrivateIpAddresses"
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, assignprivateipaddressesrequest, false)
				if err != nil {
					if IsExpectedErrors(err, []string{"InternalError", "InvalidOperation.InvalidEcsState", "InvalidOperation.InvalidEniState", "OperationConflict", "ServiceUnavailable"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, assignprivateipaddressesrequest)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			d.SetPartial("private_ips")
		}
		if err := ecsService.WaitForPrivateIpsListChanged(d.Id(), expandStringList(newPrivateIpsSet.List())); err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("private_ips_count") {
		privateIpList := expandStringList(d.Get("private_ips").(*schema.Set).List())
		oldIpsCount, newIpsCount := d.GetChange("private_ips_count")
		if oldIpsCount != nil && newIpsCount != nil && newIpsCount != len(privateIpList) {
			diff := newIpsCount.(int) - oldIpsCount.(int)
			if diff > 0 {
				assignPrivateIpsCountrequest := map[string]interface{}{
					"NetworkInterfaceId": d.Id(),
				}
				assignPrivateIpsCountrequest["RegionId"] = client.RegionId
				assignPrivateIpsCountrequest["SecondaryPrivateIpAddressCount"] = requests.NewInteger(diff)
				action := "AssignPrivateIpAddresses"
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, assignPrivateIpsCountrequest, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"InternalError", "InvalidOperation.InvalidEcsState", "InvalidOperation.InvalidEniState", "OperationConflict", "ServiceUnavailable"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, assignPrivateIpsCountrequest)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}

			if diff < 0 {
				diff *= -1
				unAssignIps := privateIpList[:diff]
				unAssignPrivateIpsCountRequest := map[string]interface{}{
					"NetworkInterfaceId": d.Id(),
				}
				unAssignPrivateIpsCountRequest["RegionId"] = client.RegionId
				unAssignPrivateIpsCountRequest["PrivateIpAddress"] = &unAssignIps
				action := "UnassignPrivateIpAddresses"
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {

					response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, unAssignPrivateIpsCountRequest, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"InternalError", "InvalidOperation.InvalidEcsState", "InvalidOperation.InvalidEniState", "OperationConflict", "ServiceUnavailable"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, unAssignPrivateIpsCountRequest)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
			err := ecsService.WaitForPrivateIpsCountChanged(d.Id(), newIpsCount.(int))
			if err != nil {
				return WrapError(err)
			}
			d.SetPartial("private_ips_count")
		}
	}
	if d.HasChange("secondary_private_ip_address_count") {
		privateIpList := expandStringList(d.Get("private_ip_addresses").(*schema.Set).List())
		oldIpsCount, newIpsCount := d.GetChange("secondary_private_ip_address_count")
		if oldIpsCount != nil && newIpsCount != nil && newIpsCount != len(privateIpList) {
			diff := newIpsCount.(int) - oldIpsCount.(int)
			if diff > 0 {
				assignSecondaryPrivateIpAddressCountrequest := map[string]interface{}{
					"NetworkInterfaceId": d.Id(),
				}
				assignSecondaryPrivateIpAddressCountrequest["RegionId"] = client.RegionId
				assignSecondaryPrivateIpAddressCountrequest["SecondaryPrivateIpAddressCount"] = requests.NewInteger(diff)
				action := "AssignPrivateIpAddresses"
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, assignSecondaryPrivateIpAddressCountrequest, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"InternalError", "InvalidOperation.InvalidEcsState", "InvalidOperation.InvalidEniState", "OperationConflict", "ServiceUnavailable"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, assignSecondaryPrivateIpAddressCountrequest)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}

			if diff < 0 {
				diff *= -1
				unAssignIps := privateIpList[:diff]
				unassignSecondaryPrivateIpAddressCountrequest := map[string]interface{}{
					"NetworkInterfaceId": d.Id(),
				}
				unassignSecondaryPrivateIpAddressCountrequest["RegionId"] = client.RegionId
				unassignSecondaryPrivateIpAddressCountrequest["PrivateIpAddress"] = &unAssignIps
				action := "UnassignPrivateIpAddresses"
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, unassignSecondaryPrivateIpAddressCountrequest, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"InternalError", "InvalidOperation.InvalidEcsState", "InvalidOperation.InvalidEniState", "OperationConflict", "ServiceUnavailable"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, unassignSecondaryPrivateIpAddressCountrequest)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}

			err := ecsService.WaitForPrivateIpsCountChanged(d.Id(), newIpsCount.(int))
			if err != nil {
				return WrapError(err)
			}
			d.SetPartial("secondary_private_ip_address_count")
		}
	}
	if d.HasChange("ipv6_address_count") {
		ipv6List := expandStringList(d.Get("ipv6_addresses").(*schema.Set).List())
		oldIpv6Count, newIpv6Count := d.GetChange("ipv6_address_count")
		if oldIpv6Count != nil && newIpv6Count != nil && newIpv6Count != len(ipv6List) {
			diff := newIpv6Count.(int) - oldIpv6Count.(int)
			if diff > 0 {
				action := "AssignIpv6Addresses"
				request := map[string]interface{}{
					"RegionId":           client.RegionId,
					"NetworkInterfaceId": d.Id(),
					"ClientToken":        buildClientToken(action),
					"Ipv6AddressCount":   diff,
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
					response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
					if err != nil {
						if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
			if diff < 0 {
				diff *= -1
				action := "UnassignIpv6Addresses"
				request := map[string]interface{}{
					"RegionId":           client.RegionId,
					"NetworkInterfaceId": d.Id(),
					"ClientToken":        buildClientToken(action),
				}
				for index, val := range ipv6List[:diff] {
					request[fmt.Sprintf("Ipv6Address.%d", index+1)] = val
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
					response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
					if err != nil {
						if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
			err = ecsService.WaitForModifyIpv6AddressCount(d.Id(), newIpv6Count.(int), DefaultTimeoutMedium)
			if err != nil {
				return WrapError(err)
			}
			d.SetPartial("ipv6_address_count")
			d.SetPartial("ipv6_addresses")
		}
	}

	if d.HasChange("ipv6_addresses") {
		oraw, nraw := d.GetChange("ipv6_addresses")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		if len(remove) > 0 {
			action := "UnassignIpv6Addresses"
			request := map[string]interface{}{
				"RegionId":           client.RegionId,
				"NetworkInterfaceId": d.Id(),
				"ClientToken":        buildClientToken(action),
			}
			for index, val := range remove {
				request[fmt.Sprintf("Ipv6Address.%d", index+1)] = val
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if len(create) > 0 {
			action := "AssignIpv6Addresses"
			request := map[string]interface{}{
				"RegionId":           client.RegionId,
				"NetworkInterfaceId": d.Id(),
				"ClientToken":        buildClientToken(action),
			}
			for index, val := range create {
				request[fmt.Sprintf("Ipv6Address.%d", index+1)] = val
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		err = ecsService.WaitForModifyIpv6Address(d.Id(), expandStringList(nraw.(*schema.Set).List()), DefaultTimeoutMedium)
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("ipv6_address_count")
		d.SetPartial("ipv6_addresses")
	}

	if d.HasChange("ipv4_prefix_count") {
		ipv4PrefixList := expandStringList(d.Get("ipv4_prefixes").(*schema.Set).List())
		oldIpv4PrefixCount, newIpv4PrefixCount := d.GetChange("ipv4_prefix_count")
		if oldIpv4PrefixCount != nil && newIpv4PrefixCount != nil && newIpv4PrefixCount != len(ipv4PrefixList) {
			diff := newIpv4PrefixCount.(int) - oldIpv4PrefixCount.(int)
			if diff > 0 {
				action := "AssignPrivateIpAddresses"
				request := map[string]interface{}{
					"RegionId":           client.RegionId,
					"NetworkInterfaceId": d.Id(),
					"ClientToken":        buildClientToken(action),
					"Ipv4PrefixCount":    diff,
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
					response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
					if err != nil {
						if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
			if diff < 0 {
				diff *= -1
				action := "UnassignPrivateIpAddresses"
				request := map[string]interface{}{
					"RegionId":           client.RegionId,
					"NetworkInterfaceId": d.Id(),
					"ClientToken":        buildClientToken(action),
				}
				for index, val := range ipv4PrefixList[:diff] {
					request[fmt.Sprintf("Ipv4Prefix.%d", index+1)] = val
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
					response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
					if err != nil {
						if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
			err = ecsService.WaitForModifyIpv4PrefixCount(d.Id(), newIpv4PrefixCount.(int), DefaultTimeoutMedium)
			if err != nil {
				return WrapError(err)
			}
			d.SetPartial("ipv4_prefix_count")
			d.SetPartial("ipv4_prefixes")
		}
	}

	if d.HasChange("ipv4_prefixes") {
		oraw, nraw := d.GetChange("ipv4_prefixes")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		if len(remove) > 0 {
			action := "UnassignPrivateIpAddresses"
			request := map[string]interface{}{
				"RegionId":           client.RegionId,
				"NetworkInterfaceId": d.Id(),
				"ClientToken":        buildClientToken(action),
			}
			for index, val := range remove {
				request[fmt.Sprintf("Ipv4Prefix.%d", index+1)] = val
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if len(create) > 0 {
			action := "AssignPrivateIpAddresses"
			request := map[string]interface{}{
				"RegionId":           client.RegionId,
				"NetworkInterfaceId": d.Id(),
				"ClientToken":        buildClientToken(action),
			}
			for index, val := range create {
				request[fmt.Sprintf("Ipv4Prefix.%d", index+1)] = val
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		err = ecsService.WaitForModifyIpv4Prefix(d.Id(), expandStringList(nraw.(*schema.Set).List()), DefaultTimeoutMedium)
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("ipv4_prefix_count")
		d.SetPartial("ipv4_prefixes")
	}
	return resourceAliCloudEcsNetworkInterfaceRead(d, meta)
}

func resourceAliCloudEcsNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	action := "DeleteNetworkInterface"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"NetworkInterfaceId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidOperation.Conflict", "InternalError", "InvalidOperation.InvalidEcsState", "InvalidOperation.InvalidEniState", "InvalidOperation.InvalidEniType", "OperationConflict", "ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidEcsId.NotFound", "InvalidEniId.NotFound", "InvalidSecurityGroupId.NotFound", "InvalidVSwitchId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecsService.EcsNetworkInterfaceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

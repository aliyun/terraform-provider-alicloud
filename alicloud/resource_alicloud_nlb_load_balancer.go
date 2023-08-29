// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudNlbLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNlbLoadBalancerCreate,
		Read:   resourceAliCloudNlbLoadBalancerRead,
		Update: resourceAliCloudNlbLoadBalancerUpdate,
		Delete: resourceAliCloudNlbLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Ipv4", "DualStack"}, false),
			},
			"address_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Internet", "Intranet"}, false),
			},
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cross_zone_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"deletion_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"deletion_protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.Get("deletion_protection_enabled").(bool)
				},
			},
			"ipv6_address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Internet", "Intranet"}, false),
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_balancer_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Network"}, false),
			},
			"modification_protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("modification_protection_status"); ok && v.(string) == "ConsoleProtection" {
						return false
					}
					return true
				},
			},
			"modification_protection_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"NonProtection", "ConsoleProtection"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_mappings": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"eni_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"allocation_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"public_ipv4_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ipv4_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudNlbLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("load_balancer_type"); ok {
		request["LoadBalancerType"] = v
	}
	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}
	request["AddressType"] = d.Get("address_type")
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIpVersion"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("zone_mappings"); ok {
		zoneMappingsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
			dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
			dataLoopMap["AllocationId"] = dataLoopTmp["allocation_id"]
			dataLoopMap["PrivateIPv4Address"] = dataLoopTmp["private_ipv4_address"]
			zoneMappingsMaps = append(zoneMappingsMaps, dataLoopMap)
		}
		request["ZoneMappings"] = zoneMappingsMaps
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request["BandwidthPackageId"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("deletion_protection_enabled"); ok {
		request["DeletionProtectionConfig.Enabled"] = v
	}
	if v, ok := d.GetOk("deletion_protection_reason"); ok {
		request["DeletionProtectionConfig.Reason"] = v
	}
	if v, ok := d.GetOk("modification_protection_status"); ok {
		request["ModificationProtectionConfig.Status"] = v
	}
	if v, ok := d.GetOk("modification_protection_reason"); ok {
		request["ModificationProtectionConfig.Reason"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_load_balancer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LoadbalancerId"]))

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.NlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNlbLoadBalancerUpdate(d, meta)
}

func resourceAliCloudNlbLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}

	objectRaw, err := nlbServiceV2.DescribeNlbLoadBalancer(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_load_balancer DescribeNlbLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address_ip_version", objectRaw["AddressIpVersion"])
	d.Set("address_type", convertNlbAddressTypeResponse(objectRaw["AddressType"]))
	d.Set("bandwidth_package_id", objectRaw["BandwidthPackageId"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("cross_zone_enabled", objectRaw["CrossZoneEnabled"])
	d.Set("ipv6_address_type", objectRaw["Ipv6AddressType"])
	d.Set("load_balancer_name", objectRaw["LoadBalancerName"])
	d.Set("load_balancer_type", convertNlbLoadBalancerTypeResponse(objectRaw["LoadBalancerType"]))
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["LoadBalancerStatus"])
	d.Set("vpc_id", objectRaw["VpcId"])
	deletionProtectionConfig1RawObj, _ := jsonpath.Get("$.DeletionProtectionConfig", objectRaw)
	deletionProtectionConfig1Raw := make(map[string]interface{})
	if deletionProtectionConfig1RawObj != nil {
		deletionProtectionConfig1Raw = deletionProtectionConfig1RawObj.(map[string]interface{})
	}
	d.Set("deletion_protection_enabled", deletionProtectionConfig1Raw["Enabled"])
	d.Set("deletion_protection_reason", deletionProtectionConfig1Raw["Reason"])
	modificationProtectionConfig1RawObj, _ := jsonpath.Get("$.ModificationProtectionConfig", objectRaw)
	modificationProtectionConfig1Raw := make(map[string]interface{})
	if modificationProtectionConfig1RawObj != nil {
		modificationProtectionConfig1Raw = modificationProtectionConfig1RawObj.(map[string]interface{})
	}
	d.Set("modification_protection_reason", modificationProtectionConfig1Raw["Reason"])
	d.Set("modification_protection_status", modificationProtectionConfig1Raw["Status"])
	securityGroupIds1Raw := make([]interface{}, 0)
	if objectRaw["SecurityGroupIds"] != nil {
		securityGroupIds1Raw = objectRaw["SecurityGroupIds"].([]interface{})
	}

	d.Set("security_group_ids", securityGroupIds1Raw)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	zoneMappings1Raw := objectRaw["ZoneMappings"]
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	if zoneMappings1Raw != nil {
		for _, zoneMappingsChild1Raw := range zoneMappings1Raw.([]interface{}) {
			zoneMappingsMap := make(map[string]interface{})
			zoneMappingsChild1Raw := zoneMappingsChild1Raw.(map[string]interface{})
			zoneMappingsMap["vswitch_id"] = zoneMappingsChild1Raw["VSwitchId"]
			zoneMappingsMap["zone_id"] = zoneMappingsChild1Raw["ZoneId"]
			loadBalancerAddressesChild1RawArrayObj, _ := jsonpath.Get("$.LoadBalancerAddresses[*]", zoneMappingsChild1Raw)
			loadBalancerAddressesChild1RawArray := make([]interface{}, 0)
			if loadBalancerAddressesChild1RawArrayObj != nil {
				loadBalancerAddressesChild1RawArray = loadBalancerAddressesChild1RawArrayObj.([]interface{})
			}
			loadBalancerAddressesChild1Raw := make(map[string]interface{})
			if len(loadBalancerAddressesChild1RawArray) > 0 {
				loadBalancerAddressesChild1Raw = loadBalancerAddressesChild1RawArray[0].(map[string]interface{})
			}

			zoneMappingsMap["allocation_id"] = loadBalancerAddressesChild1Raw["AllocationId"]
			zoneMappingsMap["eni_id"] = loadBalancerAddressesChild1Raw["EniId"]
			zoneMappingsMap["ipv6_address"] = loadBalancerAddressesChild1Raw["Ipv6Address"]
			zoneMappingsMap["private_ipv4_address"] = loadBalancerAddressesChild1Raw["PrivateIPv4Address"]
			zoneMappingsMap["public_ipv4_address"] = loadBalancerAddressesChild1Raw["PublicIPv4Address"]
			zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
		}
	}
	d.Set("zone_mappings", zoneMappingsMaps)

	return nil
}

func resourceAliCloudNlbLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateLoadBalancerAddressTypeConfig"
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("address_type") {
		update = true
	}
	request["AddressType"] = d.Get("address_type")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

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
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("address_type")
	}
	update = false
	action = "UpdateLoadBalancerAttribute"
	conn, err = client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("load_balancer_name") {
		update = true
		request["LoadBalancerName"] = d.Get("load_balancer_name")
	}

	crossZoneEnabled := true
	if v, ok := d.GetOkExists("cross_zone_enabled"); ok {
		crossZoneEnabled = v.(bool)
	}
	if d.HasChange("cross_zone_enabled") || (d.IsNewResource() && !crossZoneEnabled) {
		update = true
		request["CrossZoneEnabled"] = crossZoneEnabled
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

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
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("load_balancer_name")
		d.SetPartial("cross_zone_enabled")
	}
	update = false
	action = "UpdateLoadBalancerZones"
	conn, err = client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("zone_mappings") {
		update = true
		if v, ok := d.GetOk("zone_mappings"); ok {
			zoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range v.(*schema.Set).List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
				dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
				dataLoopMap["AllocationId"] = dataLoopTmp["allocation_id"]
				dataLoopMap["PrivateIPv4Address"] = dataLoopTmp["private_ipv4_address"]
				zoneMappingsMaps = append(zoneMappingsMaps, dataLoopMap)
			}
			request["ZoneMappings"] = zoneMappingsMaps
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

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
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 15*time.Second, nlbServiceV2.NlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateLoadBalancerProtection"
	conn, err = client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("deletion_protection_enabled") {
		update = true
		request["DeletionProtectionEnabled"] = d.Get("deletion_protection_enabled")
	}

	if !d.IsNewResource() && d.HasChange("deletion_protection_reason") {
		update = true
		request["DeletionProtectionReason"] = d.Get("deletion_protection_reason")
	}

	if !d.IsNewResource() && d.HasChange("modification_protection_status") {
		update = true
		request["ModificationProtectionStatus"] = d.Get("modification_protection_status")
	}

	if !d.IsNewResource() && d.HasChange("modification_protection_reason") {
		update = true
		request["ModificationProtectionReason"] = d.Get("modification_protection_reason")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

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
		d.SetPartial("deletion_protection_enabled")
		d.SetPartial("deletion_protection_reason")
		d.SetPartial("modification_protection_status")
		d.SetPartial("modification_protection_reason")
	}
	update = false
	action = "MoveResourceGroup"
	conn, err = client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "loadbalancer"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

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
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("ipv6_address_type") {
		client := meta.(*connectivity.AliyunClient)
		nlbServiceV2 := NlbServiceV2{client}
		object, err := nlbServiceV2.DescribeNlbLoadBalancer(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("ipv6_address_type").(string)
		if object["Ipv6AddressType"].(string) != target {
			if target == "Intranet" {
				action = "DisableLoadBalancerIpv6Internet"
				conn, err = client.NewNlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					request["ClientToken"] = buildClientToken(action)

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
			if target == "Internet" {
				action = "EnableLoadBalancerIpv6Internet"
				conn, err = client.NewNlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					request["ClientToken"] = buildClientToken(action)

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
		}
	}

	if !d.IsNewResource() && d.HasChange("bandwidth_package_id") {
		oldEntry, newEntry := d.GetChange("bandwidth_package_id")
		oldValue := oldEntry.(string)
		newValue := newEntry.(string)

		if oldValue != "" {
			action := "DetachCommonBandwidthPackageFromLoadBalancer"
			conn, err := client.NewNlbClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			request["LoadBalancerId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			request["BandwidthPackageId"] = oldValue
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				request["ClientToken"] = buildClientToken(action)

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
			nlbServiceV2 := NlbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 15*time.Second, nlbServiceV2.NlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("bandwidth_package_id")

		}

		if newValue != "" {
			action := "AttachCommonBandwidthPackageToLoadBalancer"
			conn, err := client.NewNlbClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			request["LoadBalancerId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			if v, ok := d.GetOk("bandwidth_package_id"); ok {
				request["BandwidthPackageId"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				request["ClientToken"] = buildClientToken(action)

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
			nlbServiceV2 := NlbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 15*time.Second, nlbServiceV2.NlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("bandwidth_package_id")

		}
	}
	if d.HasChange("tags") {
		nlbServiceV2 := NlbServiceV2{client}
		if err := nlbServiceV2.SetResourceTags(d, "loadbalancer"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if d.HasChange("security_group_ids") {
		oldEntry, newEntry := d.GetChange("security_group_ids")
		removed := oldEntry
		added := newEntry

		if len(removed.([]interface{})) > 0 {
			action := "LoadBalancerLeaveSecurityGroup"
			conn, err := client.NewNlbClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			request["LoadBalancerId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := removed.([]interface{})
			securityGroupIdsMaps := localData
			request["SecurityGroupIds"] = securityGroupIdsMaps

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				request["ClientToken"] = buildClientToken(action)

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
			nlbServiceV2 := NlbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 15*time.Second, nlbServiceV2.NlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if len(added.([]interface{})) > 0 {
			action := "LoadBalancerJoinSecurityGroup"
			conn, err := client.NewNlbClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			request["LoadBalancerId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := added.([]interface{})
			securityGroupIdsMaps := localData
			request["SecurityGroupIds"] = securityGroupIdsMaps

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				request["ClientToken"] = buildClientToken(action)

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
	}
	d.Partial(false)
	return resourceAliCloudNlbLoadBalancerRead(d, meta)
}

func resourceAliCloudNlbLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"ResourceNotFound.loadBalancer"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nlbServiceV2.NlbLoadBalancerJobStateRefreshFunc(d.Id(), response["JobId"].(string), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertNlbAddressTypeResponse(source interface{}) interface{} {
	switch source {
	case "Ipv4":
		return "ipv4"
	}
	return source
}
func convertNlbLoadBalancerTypeResponse(source interface{}) interface{} {
	switch source {
	}
	return source
}

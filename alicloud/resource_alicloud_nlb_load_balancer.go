// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"hash/crc32"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
			Delete: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address_ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"address_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cps": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 1000000),
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
			"deletion_protection_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"enabled_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_address_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"load_balancer_business_status": {
				Type:     schema.TypeString,
				Computed: true,
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
				ValidateFunc: StringInSlice([]string{"Network", "Gateway", "Inner"}, false),
			},
			"modification_protection_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"NonProtection", "ConsoleProtection"}, false),
						},
						"enabled_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription", "PayAsYouGo"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
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
				Set: func(v interface{}) int {
					return int(crc32.ChecksumIEEE([]byte(v.(map[string]interface{})["zone_id"].(string))))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_local_addresses": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"eni_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_local_addresses": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
							Optional: true,
							ForceNew: true,
						},
						"private_ipv4_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
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
			"modification_protection_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"NonProtection", "ConsoleProtection"}, false),
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
		},
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			for _, key := range diff.GetChangedKeysPrefix("zone_mappings") {
				// If the set contains computed key, there are some diff changes when one of element has been changed,
				// and there aims to ignore the diff
				if o, n := diff.GetChange(key); o == n {
					diff.Clear(key)
				}
			}
			return nil
		},
	}
}

func resourceAliCloudNlbLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
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
		zoneMappingsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
			dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
			dataLoopMap["AllocationId"] = dataLoopTmp["allocation_id"]
			if v, ok := dataLoopTmp["private_ipv4_address"]; ok && v != "" {
				dataLoopMap["PrivateIPv4Address"] = v
			}
			dataLoopMap["Ipv6Address"] = dataLoopTmp["ipv6_address"]
			dataLoopMap["Ipv6LocalAddresses"] = dataLoopTmp["ipv6_local_addresses"].(*schema.Set).List()
			dataLoopMap["Ipv4LocalAddresses"] = dataLoopTmp["ipv4_local_addresses"].(*schema.Set).List()
			zoneMappingsMapsArray = append(zoneMappingsMapsArray, dataLoopMap)
		}
		request["ZoneMappings"] = zoneMappingsMapsArray
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request["BandwidthPackageId"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOkExists("deletion_protection_config"); ok {
		jsonPathResult7, err := jsonpath.Get("$[0].enabled", v)
		if err == nil && jsonPathResult7 != "" {
			request["DeletionProtectionConfig.Enabled"] = jsonPathResult7
		}
	}
	if v, ok := d.GetOk("deletion_protection_config"); ok {
		jsonPathResult8, err := jsonpath.Get("$[0].reason", v)
		if err == nil && jsonPathResult8 != "" {
			request["DeletionProtectionConfig.Reason"] = jsonPathResult8
		}
	}
	if v, ok := d.GetOk("modification_protection_config"); ok {
		jsonPathResult9, err := jsonpath.Get("$[0].reason", v)
		if err == nil && jsonPathResult9 != "" {
			request["ModificationProtectionConfig.Reason"] = jsonPathResult9
		}
	}
	if v, ok := d.GetOk("modification_protection_config"); ok {
		jsonPathResult10, err := jsonpath.Get("$[0].status", v)
		if err == nil && jsonPathResult10 != "" {
			request["ModificationProtectionConfig.Status"] = jsonPathResult10
		}
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["LoadBalancerBillingConfig.PayType"] = convertNlbLoadBalancerLoadBalancerBillingConfigPayTypeRequest(v.(string))
	}

	if _, ok := d.GetOk("deletion_protection_config"); !ok {
		if v, ok := d.GetOkExists("deletion_protection_enabled"); ok {
			request["DeletionProtectionConfig.Enabled"] = v
		}
		if v, ok := d.GetOk("deletion_protection_reason"); ok {
			request["DeletionProtectionConfig.Reason"] = v
		}
	}

	if _, ok := d.GetOk("modification_protection_config"); !ok {
		if v, ok := d.GetOk("modification_protection_status"); ok {
			request["ModificationProtectionConfig.Status"] = v
		}
		if v, ok := d.GetOk("modification_protection_reason"); ok {
			request["ModificationProtectionConfig.Reason"] = v
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_load_balancer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LoadbalancerId"]))

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 0, nlbServiceV2.NlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
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
	d.Set("address_type", objectRaw["AddressType"])
	d.Set("bandwidth_package_id", objectRaw["BandwidthPackageId"])
	d.Set("cps", objectRaw["Cps"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("cross_zone_enabled", objectRaw["CrossZoneEnabled"])
	d.Set("dns_name", objectRaw["DNSName"])
	d.Set("ipv6_address_type", objectRaw["Ipv6AddressType"])
	d.Set("load_balancer_business_status", objectRaw["LoadBalancerBusinessStatus"])
	d.Set("load_balancer_name", objectRaw["LoadBalancerName"])
	d.Set("load_balancer_type", objectRaw["LoadBalancerType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["LoadBalancerStatus"])
	d.Set("vpc_id", objectRaw["VpcId"])

	loadBalancerBillingConfigRawObj, _ := jsonpath.Get("$.LoadBalancerBillingConfig", objectRaw)
	loadBalancerBillingConfigRaw := make(map[string]interface{})
	if loadBalancerBillingConfigRawObj != nil {
		loadBalancerBillingConfigRaw = loadBalancerBillingConfigRawObj.(map[string]interface{})
	}
	d.Set("payment_type", convertNlbLoadBalancerLoadBalancerBillingConfigPayTypeResponse(loadBalancerBillingConfigRaw["PayType"]))

	deletionProtectionConfigMaps := make([]map[string]interface{}, 0)
	deletionProtectionConfigMap := make(map[string]interface{})
	deletionProtectionConfigRaw := make(map[string]interface{})
	if objectRaw["DeletionProtectionConfig"] != nil {
		deletionProtectionConfigRaw = objectRaw["DeletionProtectionConfig"].(map[string]interface{})
	}
	if len(deletionProtectionConfigRaw) > 0 {
		deletionProtectionConfigMap["enabled"] = deletionProtectionConfigRaw["Enabled"]
		deletionProtectionConfigMap["enabled_time"] = deletionProtectionConfigRaw["EnabledTime"]
		deletionProtectionConfigMap["reason"] = deletionProtectionConfigRaw["Reason"]
		d.Set("deletion_protection_enabled", deletionProtectionConfigRaw["Enabled"])
		d.Set("deletion_protection_reason", deletionProtectionConfigRaw["Reason"])
		deletionProtectionConfigMaps = append(deletionProtectionConfigMaps, deletionProtectionConfigMap)
	}
	d.Set("deletion_protection_config", deletionProtectionConfigMaps)
	modificationProtectionConfigMaps := make([]map[string]interface{}, 0)
	modificationProtectionConfigMap := make(map[string]interface{})
	modificationProtectionConfigRaw := make(map[string]interface{})
	if objectRaw["ModificationProtectionConfig"] != nil {
		modificationProtectionConfigRaw = objectRaw["ModificationProtectionConfig"].(map[string]interface{})
	}
	if len(modificationProtectionConfigRaw) > 0 {
		modificationProtectionConfigMap["enabled_time"] = modificationProtectionConfigRaw["EnabledTime"]
		modificationProtectionConfigMap["reason"] = modificationProtectionConfigRaw["Reason"]
		modificationProtectionConfigMap["status"] = modificationProtectionConfigRaw["Status"]
		d.Set("modification_protection_status", modificationProtectionConfigRaw["Status"])
		d.Set("modification_protection_reason", modificationProtectionConfigRaw["Reason"])
		modificationProtectionConfigMaps = append(modificationProtectionConfigMaps, modificationProtectionConfigMap)
	}
	d.Set("modification_protection_config", modificationProtectionConfigMaps)

	securityGroupIdsRaw := make([]interface{}, 0)
	if objectRaw["SecurityGroupIds"] != nil {
		securityGroupIdsRaw = objectRaw["SecurityGroupIds"].([]interface{})
	}

	d.Set("security_group_ids", securityGroupIdsRaw)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	zoneMappingsRaw := objectRaw["ZoneMappings"]
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	if zoneMappingsRaw != nil {
		for _, zoneMappingsChildRaw := range zoneMappingsRaw.([]interface{}) {
			zoneMappingsMap := make(map[string]interface{})
			zoneMappingsChildRaw := zoneMappingsChildRaw.(map[string]interface{})
			zoneMappingsMap["status"] = zoneMappingsChildRaw["Status"]
			zoneMappingsMap["vswitch_id"] = zoneMappingsChildRaw["VSwitchId"]
			zoneMappingsMap["zone_id"] = zoneMappingsChildRaw["ZoneId"]

			loadBalancerAddressesChildRawArrayObj, _ := jsonpath.Get("$.LoadBalancerAddresses[*]", zoneMappingsChildRaw)
			loadBalancerAddressesChildRawArray := make([]interface{}, 0)
			if loadBalancerAddressesChildRawArrayObj != nil {
				loadBalancerAddressesChildRawArray = loadBalancerAddressesChildRawArrayObj.([]interface{})
			}
			loadBalancerAddressesChildRaw := make(map[string]interface{})
			if len(loadBalancerAddressesChildRawArray) > 0 {
				loadBalancerAddressesChildRaw = loadBalancerAddressesChildRawArray[0].(map[string]interface{})
			}

			zoneMappingsMap["allocation_id"] = loadBalancerAddressesChildRaw["AllocationId"]
			zoneMappingsMap["eni_id"] = loadBalancerAddressesChildRaw["EniId"]
			zoneMappingsMap["ipv6_address"] = loadBalancerAddressesChildRaw["Ipv6Address"]
			zoneMappingsMap["private_ipv4_address"] = loadBalancerAddressesChildRaw["PrivateIPv4Address"]
			zoneMappingsMap["public_ipv4_address"] = loadBalancerAddressesChildRaw["PublicIPv4Address"]

			zoneMappingsMap["ipv4_local_addresses"] = loadBalancerAddressesChildRaw["Ipv4LocalAddresses"]
			zoneMappingsMap["ipv6_local_addresses"] = loadBalancerAddressesChildRaw["Ipv6LocalAddresses"]
			zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
		}
	}
	if err := d.Set("zone_mappings", zoneMappingsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudNlbLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("ipv6_address_type") {
		nlbServiceV2 := NlbServiceV2{client}
		object, err := nlbServiceV2.DescribeNlbLoadBalancer(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("ipv6_address_type").(string)
		if object["Ipv6AddressType"].(string) != target {
			if target == "Intranet" {
				action := "DisableLoadBalancerIpv6Internet"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectStatus.Ipv6Gateway"}) || NeedRetry(err) {
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
			if target == "Internet" {
				action := "EnableLoadBalancerIpv6Internet"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectStatus.Ipv6Gateway"}) || NeedRetry(err) {
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
		}
	}

	var err error
	action := "UpdateLoadBalancerAddressTypeConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
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
			response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerStateRefreshFunc(d, response, "$.Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}
	update = false
	action = "UpdateLoadBalancerAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
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
		request["CrossZoneEnabled"] = d.Get("cross_zone_enabled")
	}

	if d.HasChange("cps") {
		update = true
		request["Cps"] = d.Get("cps")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateLoadBalancerZones"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("zone_mappings") {
		update = true
	}
	if v, ok := d.GetOk("zone_mappings"); ok || d.HasChange("zone_mappings") {
		zoneMappingsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
			dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
			dataLoopMap["AllocationId"] = dataLoopTmp["allocation_id"]
			if v, ok := dataLoopTmp["private_ipv4_address"]; ok && v != "" {
				dataLoopMap["PrivateIPv4Address"] = dataLoopTmp["private_ipv4_address"]
			}
			zoneMappingsMapsArray = append(zoneMappingsMapsArray, dataLoopMap)
		}
		request["ZoneMappings"] = zoneMappingsMapsArray
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateLoadBalancerProtection"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("deletion_protection_config.0.enabled") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].enabled", d.Get("deletion_protection_config"))
		if err == nil {
			request["DeletionProtectionEnabled"] = jsonPathResult
		}
	}

	if !d.IsNewResource() && d.HasChange("deletion_protection_config.0.reason") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].reason", d.Get("deletion_protection_config"))
		if err == nil {
			request["DeletionProtectionReason"] = jsonPathResult1
		}
	}

	if !d.IsNewResource() && d.HasChange("modification_protection_config.0.status") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$[0].status", d.Get("modification_protection_config"))
		if err == nil {
			request["ModificationProtectionStatus"] = jsonPathResult2
		}
	}

	if !d.IsNewResource() && d.HasChange("modification_protection_config.0.reason") {
		update = true
		jsonPathResult3, err := jsonpath.Get("$[0].reason", d.Get("modification_protection_config"))
		if err == nil {
			request["ModificationProtectionReason"] = jsonPathResult3
		}
	}

	if !d.IsNewResource() && d.HasChange("deletion_protection_enabled") {
		update = true
		if v, ok := d.GetOkExists("deletion_protection_enabled"); ok {
			request["DeletionProtectionEnabled"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("deletion_protection_reason") {
		update = true
		if v, ok := d.GetOk("deletion_protection_reason"); ok {
			request["DeletionProtectionReason"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("modification_protection_status") {
		update = true
		if v, ok := d.GetOk("modification_protection_status"); ok {
			request["ModificationProtectionStatus"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("modification_protection_reason") {
		update = true
		if v, ok := d.GetOk("modification_protection_reason"); ok {
			request["ModificationProtectionReason"] = v
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "loadbalancer"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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

	if d.HasChange("tags") {
		nlbServiceV2 := NlbServiceV2{client}
		if err := nlbServiceV2.SetResourceTags(d, "loadbalancer"); err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("security_group_ids") {
		oldEntry, newEntry := d.GetChange("security_group_ids")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action := "LoadBalancerLeaveSecurityGroup"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["LoadBalancerId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := removed.List()
			securityGroupIdsMapsArray := localData
			request["SecurityGroupIds"] = securityGroupIdsMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"SystemBusy"}) || NeedRetry(err) {
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
			nlbServiceV2 := NlbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerStateRefreshFunc(d, response, "$.Status", []string{}))
			if jobDetail, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
			}

		}

		if added.Len() > 0 {
			action := "LoadBalancerJoinSecurityGroup"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["LoadBalancerId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := added.List()
			securityGroupIdsMapsArray := localData
			request["SecurityGroupIds"] = securityGroupIdsMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"SystemBusy"}) || NeedRetry(err) {
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
			nlbServiceV2 := NlbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerStateRefreshFunc(d, response, "$.Status", []string{}))
			if jobDetail, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
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
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.loadBalancer"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutDelete), 30*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerStateRefreshFunc(d, response, "$.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}

func convertNlbLoadBalancerLoadBalancerBillingConfigPayTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PayAsYouGo":
		return "PostPay"
	}
	return source
}

func convertNlbLoadBalancerLoadBalancerBillingConfigPayTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PostPay":
		return "PayAsYouGo"
	}
	return source
}

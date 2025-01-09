package alicloud

import (
	"fmt"
	"hash/crc32"
	"log"
	"regexp"
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
			"dns_name": {
				Type:     schema.TypeString,
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[a-zA-Z\u4e00-\u9fa5][-\\.\\w\u4e00-\u9fa5]{1,127}$"), "The name of the network-based load balancing instance.2 to 128 English or Chinese characters in length, which must start with a letter or Chinese, and can contain numbers, half-width periods (.), underscores (_), and dashes (-)."),
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
							Type:     schema.TypeString,
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
							Optional: true,
							Computed: true,
						},
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
			if v, ok := dataLoopTmp["private_ipv4_address"]; ok && v != "" {
				dataLoopMap["PrivateIPv4Address"] = v
			}
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
	if v, ok := d.GetOk("deletion_protection_config"); ok {
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

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
	d.Set("address_type", objectRaw["AddressType"])
	d.Set("bandwidth_package_id", objectRaw["BandwidthPackageId"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("cross_zone_enabled", objectRaw["CrossZoneEnabled"])
	d.Set("ipv6_address_type", objectRaw["Ipv6AddressType"])
	d.Set("load_balancer_name", objectRaw["LoadBalancerName"])
	d.Set("load_balancer_type", objectRaw["LoadBalancerType"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["LoadBalancerStatus"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("dns_name", objectRaw["DNSName"])
	d.Set("load_balancer_business_status", objectRaw["LoadBalancerBusinessStatus"])

	deletionProtectionConfigMaps := make([]map[string]interface{}, 0)
	deletionProtectionConfigMap := make(map[string]interface{})
	deletionProtectionConfig1Raw := make(map[string]interface{})
	if objectRaw["DeletionProtectionConfig"] != nil {
		deletionProtectionConfig1Raw = objectRaw["DeletionProtectionConfig"].(map[string]interface{})
	}
	if len(deletionProtectionConfig1Raw) > 0 {
		deletionProtectionConfigMap["enabled"] = deletionProtectionConfig1Raw["Enabled"]
		deletionProtectionConfigMap["enabled_time"] = deletionProtectionConfig1Raw["EnabledTime"]
		deletionProtectionConfigMap["reason"] = deletionProtectionConfig1Raw["Reason"]
		d.Set("deletion_protection_enabled", deletionProtectionConfig1Raw["Enabled"])
		d.Set("deletion_protection_reason", deletionProtectionConfig1Raw["Reason"])

		deletionProtectionConfigMaps = append(deletionProtectionConfigMaps, deletionProtectionConfigMap)
	}
	d.Set("deletion_protection_config", deletionProtectionConfigMaps)
	modificationProtectionConfigMaps := make([]map[string]interface{}, 0)
	modificationProtectionConfigMap := make(map[string]interface{})
	modificationProtectionConfig1Raw := make(map[string]interface{})
	if objectRaw["ModificationProtectionConfig"] != nil {
		modificationProtectionConfig1Raw = objectRaw["ModificationProtectionConfig"].(map[string]interface{})
	}
	if len(modificationProtectionConfig1Raw) > 0 {
		modificationProtectionConfigMap["enabled_time"] = modificationProtectionConfig1Raw["EnabledTime"]
		modificationProtectionConfigMap["reason"] = modificationProtectionConfig1Raw["Reason"]
		modificationProtectionConfigMap["status"] = modificationProtectionConfig1Raw["Status"]
		d.Set("modification_protection_status", modificationProtectionConfig1Raw["Status"])
		d.Set("modification_protection_reason", modificationProtectionConfig1Raw["Reason"])

		modificationProtectionConfigMaps = append(modificationProtectionConfigMaps, modificationProtectionConfigMap)
	}
	d.Set("modification_protection_config", modificationProtectionConfigMaps)
	securityGroupIds1Raw := make([]interface{}, 0)
	if objectRaw["SecurityGroupIds"] != nil {
		securityGroupIds1Raw = objectRaw["SecurityGroupIds"].([]interface{})
	}

	d.Set("security_group_ids", securityGroupIds1Raw)
	zoneMappings1Raw := objectRaw["ZoneMappings"]
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	if zoneMappings1Raw != nil {
		for _, zoneMappingsChild1Raw := range zoneMappings1Raw.([]interface{}) {
			zoneMappingsMap := make(map[string]interface{})
			zoneMappingsChild1Raw := zoneMappingsChild1Raw.(map[string]interface{})
			zoneMappingsMap["status"] = zoneMappingsChild1Raw["Status"]
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

	objectRaw, err = nlbServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudNlbLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateLoadBalancerAddressTypeConfig"
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
		stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerStateRefreshFunc(d, response, "$.Status", []string{}))
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
		request["CrossZoneEnabled"] = crossZoneEnabled
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
	query = make(map[string]interface{})
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
				if v, ok := dataLoopTmp["private_ipv4_address"]; ok && v != "" {
					dataLoopMap["PrivateIPv4Address"] = v
				}
				zoneMappingsMaps = append(zoneMappingsMaps, dataLoopMap)
			}
			request["ZoneMappings"] = zoneMappingsMaps
		}
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
	}
	update = false
	action = "UpdateLoadBalancerProtection"
	conn, err = client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("deletion_protection_config") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].enabled", d.Get("deletion_protection_config"))
		if err == nil {
			request["DeletionProtectionEnabled"] = jsonPathResult
		}
	}

	if !d.IsNewResource() && d.HasChange("deletion_protection_config") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].reason", d.Get("deletion_protection_config"))
		if err == nil && jsonPathResult1 != "" {
			request["DeletionProtectionReason"] = jsonPathResult1
		}
	}

	if !d.IsNewResource() && d.HasChange("modification_protection_config") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$[0].status", d.Get("modification_protection_config"))
		if err == nil {
			request["ModificationProtectionStatus"] = jsonPathResult2
		}
	}

	if !d.IsNewResource() && d.HasChange("modification_protection_config") {
		update = true
		jsonPathResult3, err := jsonpath.Get("$[0].reason", d.Get("modification_protection_config"))
		if err == nil && jsonPathResult3 != "" {
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
		d.SetPartial("enabled")
		d.SetPartial("reason")
		d.SetPartial("status")
		d.SetPartial("reason")
	}
	update = false
	action = "MoveResourceGroup"
	conn, err = client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "loadbalancer"
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
				query = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
					request["ClientToken"] = buildClientToken(action)

					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectStatus.Ipv6Gateway"}) || NeedRetry(err) {
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
				query = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
					request["ClientToken"] = buildClientToken(action)

					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectStatus.Ipv6Gateway"}) || NeedRetry(err) {
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
	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		nlbServiceV2 := NlbServiceV2{client}
		object, err := nlbServiceV2.DescribeNlbLoadBalancer(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["$.ZoneMappings[0].Status"].(string) != target {
			if target == "Active" {
				action = "CancelShiftLoadBalancerZones"
				conn, err = client.NewNlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
			if target == "Shifted" {
				action = "StartShiftLoadBalancerZones"
				conn, err = client.NewNlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
		nlbService := NlbService{client}
		o, n := d.GetChange("bandwidth_package_id")
		oldBandwidthPackageId := o.(string)
		newBandwidthPackageId := n.(string)

		if oldBandwidthPackageId != "" {
			request = map[string]interface{}{
				"LoadBalancerId":     d.Id(),
				"RegionId":           client.RegionId,
				"BandwidthPackageId": oldBandwidthPackageId,
			}
			action := "DetachCommonBandwidthPackageFromLoadBalancer"
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
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

			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbLoadBalancerStateRefreshFunc(d, []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
		if newBandwidthPackageId != "" {
			request = map[string]interface{}{
				"LoadBalancerId":     d.Id(),
				"RegionId":           client.RegionId,
				"BandwidthPackageId": newBandwidthPackageId,
			}
			action := "AttachCommonBandwidthPackageToLoadBalancer"
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
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

			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbLoadBalancerStateRefreshFunc(d, []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		d.SetPartial("bandwidth_package_id")
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
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action := "LoadBalancerLeaveSecurityGroup"
			conn, err := client.NewNlbClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["LoadBalancerId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := removed.List()
			securityGroupIdsMaps := localData
			request["SecurityGroupIds"] = securityGroupIdsMaps

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
				request["ClientToken"] = buildClientToken(action)

				if err != nil {
					if IsExpectedErrors(err, []string{"SystemBusy"}) || NeedRetry(err) {
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
			stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerStateRefreshFunc(d, response, "$.Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if added.Len() > 0 {
			action := "LoadBalancerJoinSecurityGroup"
			conn, err := client.NewNlbClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["LoadBalancerId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := added.List()
			securityGroupIdsMaps := localData
			request["SecurityGroupIds"] = securityGroupIdsMaps

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
				request["ClientToken"] = buildClientToken(action)

				if err != nil {
					if IsExpectedErrors(err, []string{"SystemBusy"}) || NeedRetry(err) {
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
			stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerStateRefreshFunc(d, response, "$.Status", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
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
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
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
	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, nlbServiceV2.DescribeAsyncNlbLoadBalancerStateRefreshFunc(d, response, "$.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

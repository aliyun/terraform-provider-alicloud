package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudNlbLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNlbLoadBalancerCreate,
		Read:   resourceAlicloudNlbLoadBalancerRead,
		Update: resourceAlicloudNlbLoadBalancerUpdate,
		Delete: resourceAlicloudNlbLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
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
			"ipv6_address_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]{1,127}$`), "The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter."),
			},
			"load_balancer_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Network"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deletion_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"deletion_protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("modification_protection_status"); ok && v.(string) == "ConsoleProtection" {
						return false
					}
					return true
				},
			},
			"tags": tagsSchema(),
			"zone_mappings": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"allocation_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"private_ipv4_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"eni_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ipv4_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudNlbLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}
	request := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIpVersion"] = v
	}
	request["AddressType"] = d.Get("address_type")
	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}
	if v, ok := d.GetOk("load_balancer_type"); ok {
		request["LoadBalancerType"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request["BandwidthPackageId"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	for zoneMappingsPtr, zoneMappings := range d.Get("zone_mappings").(*schema.Set).List() {
		zoneMappingsArg := zoneMappings.(map[string]interface{})

		request["ZoneMappings."+fmt.Sprint(zoneMappingsPtr+1)+".VSwitchId"] = zoneMappingsArg["vswitch_id"]
		request["ZoneMappings."+fmt.Sprint(zoneMappingsPtr+1)+".ZoneId"] = zoneMappingsArg["zone_id"]
		request["ZoneMappings."+fmt.Sprint(zoneMappingsPtr+1)+".AllocationId"] = zoneMappingsArg["allocation_id"]
		request["ZoneMappings."+fmt.Sprint(zoneMappingsPtr+1)+".PrivateIPv4Address"] = zoneMappingsArg["private_ip_address"]

	}
	request["ClientToken"] = buildClientToken("CreateLoadBalancer")
	var response map[string]interface{}
	action := "CreateLoadBalancer"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_load_balancer", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.LoadbalancerId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_nlb_load_balancer")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbLoadBalancerStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudNlbLoadBalancerUpdate(d, meta)
}

func resourceAlicloudNlbLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}

	object, err := nlbService.DescribeNlbLoadBalancer(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_load_balancer nlbService.DescribeNlbLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address_ip_version", object["AddressIpVersion"])
	d.Set("address_type", object["AddressType"])
	d.Set("bandwidth_package_id", object["BandwidthPackageId"])
	d.Set("create_time", object["CreateTime"])
	d.Set("cross_zone_enabled", object["CrossZoneEnabled"])
	d.Set("dns_name", object["DNSName"])
	d.Set("ipv6_address_type", object["Ipv6AddressType"])
	d.Set("load_balancer_business_status", object["LoadBalancerBusinessStatus"])
	d.Set("load_balancer_name", object["LoadBalancerName"])
	d.Set("load_balancer_type", object["LoadBalancerType"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["LoadBalancerStatus"])
	d.Set("vpc_id", object["VpcId"])

	zoneMappingsMaps := make([]map[string]interface{}, 0)
	zoneMappingsRaw := object["ZoneMappings"]
	for _, value0 := range zoneMappingsRaw.([]interface{}) {
		zoneMappings := value0.(map[string]interface{})
		zoneMappingsMap := make(map[string]interface{})
		zoneMappingsMap["vswitch_id"] = zoneMappings["VSwitchId"]
		zoneMappingsMap["zone_id"] = zoneMappings["ZoneId"]
		if v, ok := zoneMappings["LoadBalancerAddresses"]; ok && len(v.([]interface{})) > 0 {
			LoadBalancerAddressesMap := v.([]interface{})[0].(map[string]interface{})
			zoneMappingsMap["allocation_id"] = LoadBalancerAddressesMap["AllocationId"]
			zoneMappingsMap["eni_id"] = LoadBalancerAddressesMap["EniId"]
			zoneMappingsMap["ipv6_address"] = LoadBalancerAddressesMap["Ipv6Address"]
			zoneMappingsMap["private_ipv4_address"] = LoadBalancerAddressesMap["PrivateIPv4Address"]
			zoneMappingsMap["public_ipv4_address"] = LoadBalancerAddressesMap["PublicIPv4Address"]
		}
		zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
	}
	d.Set("zone_mappings", zoneMappingsMaps)

	if deletionProtectionConfig, ok := object["DeletionProtectionConfig"]; ok {
		deletionProtectionConfigArg := deletionProtectionConfig.(map[string]interface{})

		d.Set("deletion_protection_enabled", deletionProtectionConfigArg["Enabled"])
		d.Set("deletion_protection_reason", deletionProtectionConfigArg["Reason"])
	}

	if modificationProtectionConfig, ok := object["ModificationProtectionConfig"]; ok {
		modificationProtectionConfigArg := modificationProtectionConfig.(map[string]interface{})

		d.Set("modification_protection_status", modificationProtectionConfigArg["Status"])
		d.Set("modification_protection_reason", modificationProtectionConfigArg["Reason"])
	}

	listTagResourcesObject, err := nlbService.ListTagResources(d.Id(), "loadbalancer")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAlicloudNlbLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}
	var response map[string]interface{}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := nlbService.SetResourceTags(d, "loadbalancer"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	update := false
	request := map[string]interface{}{
		"LoadBalancerId": d.Id(),
		"RegionId":       client.RegionId,
	}

	crossZoneEnabled := true
	if v, ok := d.GetOkExists("cross_zone_enabled"); ok {
		crossZoneEnabled = v.(bool)
	}
	if d.HasChange("cross_zone_enabled") || (d.IsNewResource() && !crossZoneEnabled) {
		update = true
		request["CrossZoneEnabled"] = crossZoneEnabled
	}
	if !d.IsNewResource() && d.HasChange("load_balancer_name") {
		update = true
		if v, ok := d.GetOk("load_balancer_name"); ok {
			request["LoadBalancerName"] = v
		}
	}

	if update {
		action := "UpdateLoadBalancerAttribute"
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

		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbService.NlbLoadBalancerStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("cross_zone_enabled")
		d.SetPartial("load_balancer_name")
	}

	update = false
	request = map[string]interface{}{
		"LoadBalancerId": d.Id(),
		"RegionId":       client.RegionId,
	}

	if !d.IsNewResource() && d.HasChange("zone_mappings") {
		update = true
		if v, ok := d.GetOk("zone_mappings"); ok {
			zoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, value0 := range v.(*schema.Set).List() {
				zoneMappings := value0.(map[string]interface{})
				zoneMappingsMap := make(map[string]interface{})
				zoneMappingsMap["VSwitchId"] = zoneMappings["vswitch_id"]
				zoneMappingsMap["ZoneId"] = zoneMappings["zone_id"]
				zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
			}
			request["ZoneMappings"] = zoneMappingsMaps
		}
	}

	if update {
		action := "UpdateLoadBalancerZones"
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

		d.SetPartial("zone_mappings")
	}

	update = false
	request = map[string]interface{}{
		"LoadBalancerId": d.Id(),
		"RegionId":       client.RegionId,
	}

	if !d.IsNewResource() && d.HasChange("address_type") {
		update = true
		request["AddressType"] = d.Get("address_type")
		if v, ok := d.GetOk("zone_mappings"); ok {
			zoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, value0 := range v.(*schema.Set).List() {
				zoneMappings := value0.(map[string]interface{})
				zoneMappingsMap := make(map[string]interface{})
				if v, ok := zoneMappings["allocation_id"]; ok && v != "" {
					zoneMappingsMap["AllocationId"] = v
				}
				if v, ok := zoneMappings["private_ipv4_address"]; ok && v != "" {
					zoneMappingsMap["PrivateIPv4Address"] = v
				}
				zoneMappingsMap["VSwitchId"] = zoneMappings["vswitch_id"]
				zoneMappingsMap["ZoneId"] = zoneMappings["zone_id"]
				zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
			}
			request["ZoneMappings"] = zoneMappingsMaps
		}
	}
	if update {
		action := "UpdateLoadBalancerAddressTypeConfig"
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

		d.SetPartial("address_type")
	}

	if !d.IsNewResource() && d.HasChange("bandwidth_package_id") {
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

	update = false
	updateLoadBalancerProtectionReq := map[string]interface{}{
		"RegionId":       client.RegionId,
		"ClientToken":    buildClientToken("UpdateLoadBalancerProtection"),
		"LoadBalancerId": d.Id(),
	}

	if d.HasChange("deletion_protection_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("deletion_protection_enabled"); ok {
		updateLoadBalancerProtectionReq["DeletionProtectionEnabled"] = v
	}

	if d.HasChange("deletion_protection_reason") {
		update = true
	}
	if v, ok := d.GetOk("deletion_protection_reason"); ok {
		updateLoadBalancerProtectionReq["DeletionProtectionReason"] = v
	}

	if d.HasChange("modification_protection_status") {
		update = true
	}
	if v, ok := d.GetOk("modification_protection_status"); ok {
		updateLoadBalancerProtectionReq["ModificationProtectionStatus"] = v
	}

	if d.HasChange("modification_protection_reason") {
		update = true
	}
	if v, ok := d.GetOk("modification_protection_reason"); ok {
		updateLoadBalancerProtectionReq["ModificationProtectionReason"] = v
	}

	if update {
		action := "UpdateLoadBalancerProtection"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, updateLoadBalancerProtectionReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateLoadBalancerProtectionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbLoadBalancerStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("deletion_protection_enabled")
		d.SetPartial("deletion_protection_reason")
		d.SetPartial("modification_protection_status")
		d.SetPartial("modification_protection_reason")
	}

	d.Partial(false)
	return resourceAlicloudNlbLoadBalancerRead(d, meta)
}

func resourceAlicloudNlbLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"LoadBalancerId": d.Id(),
		"RegionId":       client.RegionId,
	}

	request["ClientToken"] = buildClientToken("DeleteLoadBalancer")
	action := "DeleteLoadBalancer"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"ResourceInConfiguring.loadbalancer"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nlbService.NlbLoadBalancerStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

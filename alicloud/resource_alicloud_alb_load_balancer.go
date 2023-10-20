// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAlbLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlbLoadBalancerCreate,
		Read:   resourceAliCloudAlbLoadBalancerRead,
		Update: resourceAliCloudAlbLoadBalancerUpdate,
		Delete: resourceAliCloudAlbLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_log_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_store": {
							Type:     schema.TypeString,
							Required: true,
						},
						"log_project": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"address_allocated_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Fixed", "Dynamic"}, false),
			},
			"address_ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"address_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Internet", "Intranet"}, false),
			},
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"load_balancer_billing_config": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pay_type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"PayAsYouGo"}, false),
						},
					},
				},
			},
			"load_balancer_edition": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Basic", "Standard", "StandardWithWaf"}, false),
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modification_protection_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"ConsoleProtection", "NonProtection"}, false),
						},
						"reason": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: modificationProtectionConfigDiffSuppressFunc,
							ValidateFunc:     StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-.]{1,127}$`), "The reason must be 2 to 128 characters in length, and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-)."),
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"load_balancer_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"allocation_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"eip_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudAlbLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["VpcId"] = d.Get("vpc_id")
	request["AddressType"] = d.Get("address_type")
	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}
	if v, ok := d.GetOk("address_allocated_mode"); ok {
		request["AddressAllocatedMode"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["LoadBalancerEdition"] = d.Get("load_balancer_edition")
	if v, ok := d.GetOk("zone_mappings"); ok {
		zoneMappingsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
			dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
			zoneMappingsMaps = append(zoneMappingsMaps, dataLoopMap)
		}
		request["ZoneMappings"] = zoneMappingsMaps
	}

	for _, loadBalancerBillingConfigs := range d.Get("load_balancer_billing_config").(*schema.Set).List() {
		loadBalancerBillingConfigArg := loadBalancerBillingConfigs.(map[string]interface{})
		request["LoadBalancerBillingConfig.PayType"] = convertAlbLoadBalancerBillingConfigPayTypeRequest(loadBalancerBillingConfigArg["pay_type"].(string))
	}

	for _, modificationProtectionConfigs := range d.Get("modification_protection_config").(*schema.Set).List() {
		modificationProtectionConfigArg := modificationProtectionConfigs.(map[string]interface{})
		request["ModificationProtectionConfig.Reason"] = modificationProtectionConfigArg["reason"]
		request["ModificationProtectionConfig.Status"] = modificationProtectionConfigArg["status"]
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request["LoadBalancerBillingConfig.BandwidthPackageId"] = v
	}
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIpVersion"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_load_balancer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LoadBalancerId"]))

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAlbLoadBalancerUpdate(d, meta)
}

func resourceAliCloudAlbLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}

	objectRaw, err := albServiceV2.DescribeAlbLoadBalancer(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_load_balancer DescribeAlbLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address_allocated_mode", objectRaw["AddressAllocatedMode"])
	d.Set("address_ip_version", objectRaw["AddressIpVersion"])
	d.Set("address_type", objectRaw["AddressType"])
	d.Set("bandwidth_package_id", objectRaw["BandwidthPackageId"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("dns_name", objectRaw["DNSName"])
	d.Set("ipv6_address_type", objectRaw["Ipv6AddressType"])
	d.Set("load_balancer_edition", objectRaw["LoadBalancerEdition"])
	d.Set("load_balancer_name", objectRaw["LoadBalancerName"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["LoadBalancerStatus"])
	d.Set("vpc_id", objectRaw["VpcId"])
	accessLogConfigMaps := make([]map[string]interface{}, 0)
	accessLogConfigMap := make(map[string]interface{})
	accessLogConfig1Raw := make(map[string]interface{})
	if objectRaw["AccessLogConfig"] != nil {
		accessLogConfig1Raw = objectRaw["AccessLogConfig"].(map[string]interface{})
	}
	if len(accessLogConfig1Raw) > 0 {
		accessLogConfigMap["log_project"] = accessLogConfig1Raw["LogProject"]
		accessLogConfigMap["log_store"] = accessLogConfig1Raw["LogStore"]
		accessLogConfigMaps = append(accessLogConfigMaps, accessLogConfigMap)
	}
	d.Set("access_log_config", accessLogConfigMaps)

	loadBalancerBillingConfigMaps := make([]map[string]interface{}, 0)
	loadBalancerBillingConfigMap := make(map[string]interface{})
	loadBalancerBillingConfig1Raw := make(map[string]interface{})
	if objectRaw["LoadBalancerBillingConfig"] != nil {
		loadBalancerBillingConfig1Raw = objectRaw["LoadBalancerBillingConfig"].(map[string]interface{})
	}
	if len(loadBalancerBillingConfig1Raw) > 0 {
		loadBalancerBillingConfigMap["pay_type"] = convertAlbLoadBalancerBillingConfigPayTypeResponse(loadBalancerBillingConfig1Raw["PayType"])
		loadBalancerBillingConfigMaps = append(loadBalancerBillingConfigMaps, loadBalancerBillingConfigMap)
	}
	d.Set("load_balancer_billing_config", loadBalancerBillingConfigMaps)

	modificationProtectionConfigMaps := make([]map[string]interface{}, 0)
	modificationProtectionConfigMap := make(map[string]interface{})
	modificationProtectionConfig1Raw := make(map[string]interface{})
	if objectRaw["ModificationProtectionConfig"] != nil {
		modificationProtectionConfig1Raw = objectRaw["ModificationProtectionConfig"].(map[string]interface{})
	}
	if len(modificationProtectionConfig1Raw) > 0 {
		modificationProtectionConfigMap["reason"] = modificationProtectionConfig1Raw["Reason"]
		modificationProtectionConfigMap["status"] = modificationProtectionConfig1Raw["Status"]
		modificationProtectionConfigMaps = append(modificationProtectionConfigMaps, modificationProtectionConfigMap)
	}
	d.Set("modification_protection_config", modificationProtectionConfigMaps)

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	zoneMappings1Raw := objectRaw["ZoneMappings"]
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	if zoneMappings1Raw != nil {
		for _, zoneMappingsChild1Raw := range zoneMappings1Raw.([]interface{}) {
			zoneMappingsMap := make(map[string]interface{})
			zoneMappingsChild1Raw := zoneMappingsChild1Raw.(map[string]interface{})
			zoneMappingsMap["load_balancer_addresses"] = zoneMappingsChild1Raw["LoadBalancerAddresses"]
			zoneMappingsMap["vswitch_id"] = zoneMappingsChild1Raw["VSwitchId"]
			zoneMappingsMap["zone_id"] = zoneMappingsChild1Raw["ZoneId"]
			loadBalancerAddresses3Raw := zoneMappingsChild1Raw["LoadBalancerAddresses"]
			loadBalancerAddressesMaps := make([]map[string]interface{}, 0)
			if loadBalancerAddresses3Raw != nil {
				for _, loadBalancerAddressesChild1Raw := range loadBalancerAddresses3Raw.([]interface{}) {
					loadBalancerAddressesMap := make(map[string]interface{})
					loadBalancerAddressesChild1Raw := loadBalancerAddressesChild1Raw.(map[string]interface{})
					loadBalancerAddressesMap["address"] = loadBalancerAddressesChild1Raw["Address"]
					loadBalancerAddressesMap["allocation_id"] = loadBalancerAddressesChild1Raw["AllocationId"]
					loadBalancerAddressesMap["eip_type"] = loadBalancerAddressesChild1Raw["EipType"]
					loadBalancerAddressesMap["ipv6_address"] = loadBalancerAddressesChild1Raw["Ipv6Address"]
					loadBalancerAddressesMaps = append(loadBalancerAddressesMaps, loadBalancerAddressesMap)
				}
			}
			zoneMappingsMap["load_balancer_addresses"] = loadBalancerAddressesMaps
			zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
		}
	}
	d.Set("zone_mappings", zoneMappingsMaps)

	return nil
}

func resourceAliCloudAlbLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateLoadBalancerAttribute"
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("load_balancer_name") {
		update = true
		request["LoadBalancerName"] = d.Get("load_balancer_name")
	}

	if !d.IsNewResource() && d.HasChange("modification_protection_config") {
		update = true
		modificationProtectionConfigMap := map[string]interface{}{}
		if v, ok := d.GetOk("modification_protection_config"); ok {
			for _, modificationProtectionConfigs := range v.(*schema.Set).List() {
				modificationProtectionConfigArg := modificationProtectionConfigs.(map[string]interface{})
				modificationProtectionConfigMap["Reason"] = modificationProtectionConfigArg["reason"]
				modificationProtectionConfigMap["Status"] = modificationProtectionConfigArg["status"]
			}
		}
		request["ModificationProtectionConfig"] = modificationProtectionConfigMap
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("load_balancer_name")
		d.SetPartial("reason")
		d.SetPartial("status")
	}
	update = false
	action = "UpdateLoadBalancerEdition"
	conn, err = client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("load_balancer_edition") {
		update = true
	}
	request["LoadBalancerEdition"] = d.Get("load_balancer_edition")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("load_balancer_edition"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerEdition", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("load_balancer_edition")
	}
	update = false
	action = "EnableLoadBalancerAccessLog"
	conn, err = client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("access_log_config") {
		if v, ok := d.GetOk("access_log_config"); ok {
			for _, enableLoadBalancerAccessLogs := range v.(*schema.Set).List() {
				update = true
				enableLoadBalancerAccessArg := enableLoadBalancerAccessLogs.(map[string]interface{})
				request["LogProject"] = enableLoadBalancerAccessArg["log_project"]
				request["LogStore"] = enableLoadBalancerAccessArg["log_store"]
			}
		}
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.AccessLogEnabled", "SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("log_project")
		d.SetPartial("log_store")
	}
	update = false
	action = "DisableLoadBalancerAccessLog"
	conn, err = client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("access_log_config") {
		oraw, _ := d.GetChange("access_log_config")
		if oraw != nil && oraw.(*schema.Set).Len() > 0 {
			update = true
		}
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.AccessLogEnabled", "SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("log_project")
		d.SetPartial("log_store")
	}
	update = false
	action = "MoveResourceGroup"
	conn, err = client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	if v, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = v
	}

	request["ResourceType"] = "loadbalancer"
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)

			if err != nil {
				if IsExpectedErrors(err, []string{"undefined"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
	}
	update = false
	action = "UpdateLoadBalancerAddressTypeConfig"
	conn, err = client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("address_type") {
		update = true
	}
	request["AddressType"] = d.Get("address_type")
	if !d.IsNewResource() && d.HasChange("zone_mappings") {
		update = true
		if v, ok := d.GetOk("zone_mappings"); ok {
			zoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range v.(*schema.Set).List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
				dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
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
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("address_type")
	}

	if d.HasChange("enabled") {
		client := meta.(*connectivity.AliyunClient)
		albServiceV2 := AlbServiceV2{client}
		object, err := albServiceV2.DescribeAlbLoadBalancer(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("enabled").(bool)
		if object["Enabled"].(bool) != target {
			if target == true {
				action = "EnableDeletionProtection"
				conn, err = client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				request["ResourceId"] = d.Id()
				request["ClientToken"] = buildClientToken(action)
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
					request["ClientToken"] = buildClientToken(action)

					if err != nil {
						if IsExpectedErrors(err, []string{"IdempotenceProcessing", "SystemBusy"}) || NeedRetry(err) {
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
				albServiceV2 := AlbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}
	if d.HasChange("ipv6_address_type") {
		client := meta.(*connectivity.AliyunClient)
		albServiceV2 := AlbServiceV2{client}
		object, err := albServiceV2.DescribeAlbLoadBalancer(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("ipv6_address_type").(string)
		if object["Ipv6AddressType"] != nil && object["Ipv6AddressType"].(string) != target {
			if target == "Internet" {
				action = "EnableLoadBalancerIpv6Internet"
				conn, err = client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()
				request["ClientToken"] = buildClientToken(action)
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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
				albServiceV2 := AlbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Internet"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "Ipv6AddressType", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Intranet" {
				action = "DisableLoadBalancerIpv6Internet"
				conn, err = client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()
				request["ClientToken"] = buildClientToken(action)
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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
				albServiceV2 := AlbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Intranet"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "Ipv6AddressType", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	if d.HasChange("tags") {
		albServiceV2 := AlbServiceV2{client}
		if err := albServiceV2.SetResourceTags(d, "loadbalancer"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudAlbLoadBalancerRead(d, meta)
}

func resourceAliCloudAlbLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.LoadBalancer"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albServiceV2.AlbLoadBalancerJobStateRefreshFunc(d.Id(), response["JobId"].(string), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertAlbLoadBalancerPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPay":
		return "PayAsYouGo"
	}
	return source
}

func modificationProtectionConfigDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("modification_protection_config"); ok {
		val := v.(*schema.Set).List()
		if len(val) > 2 {
			// modification_protection_config 为 Object 类型
			return true
		}
		for _, modificationProtectionConfigs := range val {
			modificationProtectionConfigArg := modificationProtectionConfigs.(map[string]interface{})
			return fmt.Sprintf(modificationProtectionConfigArg["status"].(string)) != "ConsoleProtection"
		}
	}
	return true
}

func convertAlbLoadBalancerBillingConfigPayTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPay":
		return "PayAsYouGo"
	}
	return source
}
func convertAlbRegionIdResponse(source interface{}) interface{} {
	switch source {
	case "cn-hangzhou-onebox-nebula":
		return "cn-hangzhou"
	}
	return source
}
func convertAlbLoadBalancerBillingConfigPayTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPay"
	}
	return source
}

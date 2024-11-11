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
			"load_balancer_edition": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Basic", "Standard", "StandardWithWaf"}, false),
			},
			"address_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Internet", "Intranet"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address_allocated_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Fixed", "Dynamic"}, false),
			},
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"IPv4", "DualStack"}, false),
			},
			"ipv6_address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Internet", "Intranet"}, false),
			},
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deletion_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
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
							ValidateFunc:     StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-.]{1,127}$`), "The reason must be 2 to 128 characters in length, and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-)."),
							DiffSuppressFunc: modificationProtectionConfigDiffSuppressFunc,
						},
					},
				},
			},
			"access_log_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_project": {
							Type:     schema.TypeString,
							Required: true,
						},
						"log_store": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"zone_mappings": {
				Type:     schema.TypeSet,
				Required: true,
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
						"load_balancer_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allocation_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"eip_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"address": {
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
			"tags": tagsSchema(),
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAlbLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}
	var response map[string]interface{}
	action := "CreateLoadBalancer"
	request := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}

	request["ClientToken"] = buildClientToken("CreateLoadBalancer")
	request["LoadBalancerEdition"] = d.Get("load_balancer_edition")
	request["AddressType"] = d.Get("address_type")
	request["VpcId"] = d.Get("vpc_id")

	if v, ok := d.GetOk("address_allocated_mode"); ok {
		request["AddressAllocatedMode"] = v
	}

	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIpVersion"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}

	if v, ok := d.GetOkExists("deletion_protection_enabled"); ok {
		request["DeletionProtectionEnabled"] = v
	}

	loadBalancerBillingConfig := d.Get("load_balancer_billing_config")
	loadBalancerBillingConfigMap := map[string]interface{}{}
	for _, loadBalancerBillingConfigList := range loadBalancerBillingConfig.(*schema.Set).List() {
		loadBalancerBillingConfigArg := loadBalancerBillingConfigList.(map[string]interface{})

		loadBalancerBillingConfigMap["PayType"] = convertAlbLoadBalancerBillingConfigPayTypeRequest(loadBalancerBillingConfigArg["pay_type"])
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		loadBalancerBillingConfigMap["BandwidthPackageId"] = v
	}

	request["LoadBalancerBillingConfig"] = loadBalancerBillingConfigMap

	if v, ok := d.GetOk("modification_protection_config"); ok {
		modificationProtectionConfigMap := map[string]interface{}{}
		for _, modificationProtectionConfigList := range v.(*schema.Set).List() {
			modificationProtectionConfigArg := modificationProtectionConfigList.(map[string]interface{})

			if status, ok := modificationProtectionConfigArg["status"]; ok {
				modificationProtectionConfigMap["Status"] = status
			}

			if reason, ok := modificationProtectionConfigArg["reason"]; ok {
				modificationProtectionConfigMap["Reason"] = reason
			}
		}

		request["ModificationProtectionConfig"] = modificationProtectionConfigMap
	}

	zoneMappings := d.Get("zone_mappings")
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	for _, zoneMappingsList := range zoneMappings.(*schema.Set).List() {
		zoneMappingsMap := make(map[string]interface{})
		zoneMappingsArg := zoneMappingsList.(map[string]interface{})

		zoneMappingsMap["VSwitchId"] = zoneMappingsArg["vswitch_id"]
		zoneMappingsMap["ZoneId"] = zoneMappingsArg["zone_id"]

		zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
	}

	request["ZoneMappings"] = zoneMappingsMaps

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tag"] = tagsMap
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_load_balancer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LoadBalancerId"]))

	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{"CreateFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAlbLoadBalancerUpdate(d, meta)
}

func resourceAliCloudAlbLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}

	object, err := albServiceV2.DescribeAlbLoadBalancer(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_load_balancer albServiceV2.DescribeAlbLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("load_balancer_edition", object["LoadBalancerEdition"])
	d.Set("address_type", object["AddressType"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("address_allocated_mode", object["AddressAllocatedMode"])
	d.Set("address_ip_version", convertAlbLoadBalancerAddressIpVersionResponse(object["AddressIpVersion"]))
	d.Set("ipv6_address_type", object["Ipv6AddressType"])
	d.Set("bandwidth_package_id", object["BandwidthPackageId"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("load_balancer_name", object["LoadBalancerName"])
	d.Set("tags", tagsToMap(object["Tags"]))
	d.Set("dns_name", object["DNSName"])
	d.Set("status", object["LoadBalancerStatus"])
	d.Set("create_time", object["CreateTime"])

	if deletionProtectionConfig, ok := object["DeletionProtectionConfig"]; ok {
		deletionProtectionConfigArg := deletionProtectionConfig.(map[string]interface{})

		if enabled, ok := deletionProtectionConfigArg["Enabled"]; ok {
			d.Set("deletion_protection_enabled", enabled)
		}
	}

	if loadBalancerBillingConfig, ok := object["LoadBalancerBillingConfig"]; ok {
		loadBalancerBillingConfigMaps := make([]map[string]interface{}, 0)
		loadBalancerBillingConfigArg := loadBalancerBillingConfig.(map[string]interface{})
		loadBalancerBillingConfigMap := make(map[string]interface{})

		if payType, ok := loadBalancerBillingConfigArg["PayType"]; ok {
			loadBalancerBillingConfigMap["pay_type"] = convertAlbLoadBalancerBillingConfigPayTypeResponse(payType)
		}

		loadBalancerBillingConfigMaps = append(loadBalancerBillingConfigMaps, loadBalancerBillingConfigMap)

		d.Set("load_balancer_billing_config", loadBalancerBillingConfigMaps)
	}

	if modificationProtectionConfig, ok := object["ModificationProtectionConfig"]; ok {
		modificationProtectionConfigMaps := make([]map[string]interface{}, 0)
		modificationProtectionConfigArg := modificationProtectionConfig.(map[string]interface{})
		modificationProtectionConfigMap := make(map[string]interface{})

		if status, ok := modificationProtectionConfigArg["Status"]; ok {
			modificationProtectionConfigMap["status"] = status
		}

		if reason, ok := modificationProtectionConfigArg["Reason"]; ok {
			modificationProtectionConfigMap["reason"] = reason
		}

		modificationProtectionConfigMaps = append(modificationProtectionConfigMaps, modificationProtectionConfigMap)

		d.Set("modification_protection_config", modificationProtectionConfigMaps)
	}

	if accessLogConfig, ok := object["AccessLogConfig"]; ok {
		accessLogConfigMaps := make([]map[string]interface{}, 0)
		accessLogConfigArg := accessLogConfig.(map[string]interface{})
		accessLogConfigMap := make(map[string]interface{})

		if logProject, ok := accessLogConfigArg["LogProject"]; ok {
			accessLogConfigMap["log_project"] = logProject
		}

		if logStore, ok := accessLogConfigArg["LogStore"]; ok {
			accessLogConfigMap["log_store"] = logStore
		}

		accessLogConfigMaps = append(accessLogConfigMaps, accessLogConfigMap)

		d.Set("access_log_config", accessLogConfigMaps)
	}

	if zoneMappings, ok := object["ZoneMappings"]; ok {
		zoneMappingsMaps := make([]map[string]interface{}, 0)
		for _, zoneMappingsList := range zoneMappings.([]interface{}) {
			zoneMappingsArg := zoneMappingsList.(map[string]interface{})
			zoneMappingsMap := map[string]interface{}{}

			if vSwitchId, ok := zoneMappingsArg["VSwitchId"]; ok {
				zoneMappingsMap["vswitch_id"] = vSwitchId
			}

			if zoneId, ok := zoneMappingsArg["ZoneId"]; ok {
				zoneMappingsMap["zone_id"] = zoneId
			}

			if loadBalancerAddresses, ok := zoneMappingsArg["LoadBalancerAddresses"]; ok {
				loadBalancerAddressesMaps := make([]map[string]interface{}, 0)
				for _, loadBalancerAddressesList := range loadBalancerAddresses.([]interface{}) {
					loadBalancerAddressesArg := loadBalancerAddressesList.(map[string]interface{})
					loadBalancerAddressesMap := make(map[string]interface{})

					if allocationId, ok := loadBalancerAddressesArg["AllocationId"]; ok {
						loadBalancerAddressesMap["allocation_id"] = allocationId
					}

					if eipType, ok := loadBalancerAddressesArg["EipType"]; ok {
						loadBalancerAddressesMap["eip_type"] = eipType
					}

					if address, ok := loadBalancerAddressesArg["Address"]; ok {
						loadBalancerAddressesMap["address"] = address
					}

					if ipv6Address, ok := loadBalancerAddressesArg["Ipv6Address"]; ok {
						loadBalancerAddressesMap["ipv6_address"] = ipv6Address
					}

					loadBalancerAddressesMaps = append(loadBalancerAddressesMaps, loadBalancerAddressesMap)
				}

				zoneMappingsMap["load_balancer_addresses"] = loadBalancerAddressesMaps
			}

			zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
		}

		d.Set("zone_mappings", zoneMappingsMaps)
	}

	return nil
}

func resourceAliCloudAlbLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albServiceV2 := AlbServiceV2{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	updateLoadBalancerEditionReq := map[string]interface{}{
		"ClientToken":    buildClientToken("UpdateLoadBalancerEdition"),
		"LoadBalancerId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("load_balancer_edition") {
		update = true
	}
	updateLoadBalancerEditionReq["LoadBalancerEdition"] = d.Get("load_balancer_edition")

	if v, ok := d.GetOkExists("dry_run"); ok {
		updateLoadBalancerEditionReq["DryRun"] = v
	}

	if update {
		action := "UpdateLoadBalancerEdition"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateLoadBalancerEditionReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateLoadBalancerEditionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("load_balancer_edition"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerEdition", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("load_balancer_edition")
	}

	update = false
	updateLoadBalancerAddressTypeConfigReq := map[string]interface{}{
		"ClientToken":    buildClientToken("UpdateLoadBalancerAddressTypeConfig"),
		"LoadBalancerId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("address_type") {
		update = true
	}
	updateLoadBalancerAddressTypeConfigReq["AddressType"] = d.Get("address_type")

	if v, ok := d.GetOkExists("dry_run"); ok {
		updateLoadBalancerAddressTypeConfigReq["DryRun"] = v
	}

	if update {
		action := "UpdateLoadBalancerAddressTypeConfig"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateLoadBalancerAddressTypeConfigReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateLoadBalancerAddressTypeConfigReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("address_type")
	}

	update = false
	moveResourceGroupReq := map[string]interface{}{
		"ResourceType": "loadbalancer",
		"ResourceId":   d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		moveResourceGroupReq["NewResourceGroupId"] = v
	}

	if update {
		action := "MoveResourceGroup"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, moveResourceGroupReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"undefined"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, moveResourceGroupReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("resource_group_id")
	}

	update = false
	updateLoadBalancerAttributeReq := map[string]interface{}{
		"ClientToken":    buildClientToken("UpdateLoadBalancerAttribute"),
		"LoadBalancerId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("load_balancer_name") {
		update = true
	}
	if v, ok := d.GetOk("load_balancer_name"); ok {
		updateLoadBalancerAttributeReq["LoadBalancerName"] = v
	}

	if !d.IsNewResource() && d.HasChange("modification_protection_config") {
		update = true
	}
	if v, ok := d.GetOk("modification_protection_config"); ok {
		modificationProtectionConfigMap := map[string]interface{}{}
		for _, modificationProtectionConfigList := range v.(*schema.Set).List() {
			modificationProtectionConfigArg := modificationProtectionConfigList.(map[string]interface{})

			if status, ok := modificationProtectionConfigArg["status"]; ok {
				modificationProtectionConfigMap["Status"] = status
			}

			if reason, ok := modificationProtectionConfigArg["reason"]; ok {
				modificationProtectionConfigMap["Reason"] = reason
			}
		}

		updateLoadBalancerAttributeReq["ModificationProtectionConfig"] = modificationProtectionConfigMap
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		updateLoadBalancerAttributeReq["DryRun"] = v
	}

	if update {
		action := "UpdateLoadBalancerAttribute"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateLoadBalancerAttributeReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateLoadBalancerAttributeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("load_balancer_name")
		d.SetPartial("modification_protection_config")
	}

	update = false
	updateLoadBalancerZonesReq := map[string]interface{}{
		"ClientToken":    buildClientToken("UpdateLoadBalancerZones"),
		"LoadBalancerId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("zone_mappings") {
		update = true
	}
	zoneMappings := d.Get("zone_mappings")
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	for _, zoneMappingsList := range zoneMappings.(*schema.Set).List() {
		zoneMappingsMap := make(map[string]interface{})
		zoneMappingsArg := zoneMappingsList.(map[string]interface{})

		zoneMappingsMap["VSwitchId"] = zoneMappingsArg["vswitch_id"]
		zoneMappingsMap["ZoneId"] = zoneMappingsArg["zone_id"]

		zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
	}

	updateLoadBalancerZonesReq["ZoneMappings"] = zoneMappingsMaps

	if v, ok := d.GetOkExists("dry_run"); ok {
		updateLoadBalancerZonesReq["DryRun"] = v
	}

	if update {
		action := "UpdateLoadBalancerZones"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateLoadBalancerZonesReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateLoadBalancerZonesReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("zone_mappings")
	}

	if d.HasChange("ipv6_address_type") {
		object, err := albServiceV2.DescribeAlbLoadBalancer(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("ipv6_address_type").(string)
		if object["Ipv6AddressType"] != nil && object["Ipv6AddressType"].(string) != target {
			if target == "Intranet" {
				request := map[string]interface{}{
					"ClientToken":    buildClientToken("DisableLoadBalancerIpv6Internet"),
					"LoadBalancerId": d.Id(),
				}

				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}

				action := "DisableLoadBalancerIpv6Internet"
				conn, err := client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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

				stateConf := BuildStateConf([]string{}, []string{"Intranet"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "Ipv6AddressType", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			if target == "Internet" {
				request := map[string]interface{}{
					"ClientToken":    buildClientToken("EnableLoadBalancerIpv6Internet"),
					"LoadBalancerId": d.Id(),
				}

				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}

				action := "EnableLoadBalancerIpv6Internet"
				conn, err := client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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

				stateConf := BuildStateConf([]string{}, []string{"Internet"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "Ipv6AddressType", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			d.SetPartial("ipv6_address_type")
		}
	}

	if !d.IsNewResource() && d.HasChange("deletion_protection_enabled") {
		current := false

		object, err := albServiceV2.DescribeAlbLoadBalancer(d.Id())
		if err != nil {
			return WrapError(err)
		}

		if deletionProtectionConfig, ok := object["DeletionProtectionConfig"]; ok {
			deletionProtectionConfigArg := deletionProtectionConfig.(map[string]interface{})

			if enabled, ok := deletionProtectionConfigArg["Enabled"]; ok {
				current = enabled.(bool)
			}
		}

		target := d.Get("deletion_protection_enabled").(bool)
		if current != target {
			if target == false {
				request := map[string]interface{}{
					"ClientToken": buildClientToken("DisableDeletionProtection"),
					"ResourceId":  d.Id(),
				}

				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}

				action := "DisableDeletionProtection"
				conn, err := client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
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

				stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			if target == true {
				request := map[string]interface{}{
					"ClientToken": buildClientToken("EnableDeletionProtection"),
					"ResourceId":  d.Id(),
				}

				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}

				action := "EnableDeletionProtection"
				conn, err := client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
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

				stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			d.SetPartial("deletion_protection_enabled")
		}
	}

	if d.HasChange("access_log_config") {
		oldAccessLogConfig, newAccessLogConfig := d.GetChange("access_log_config")
		removed := oldAccessLogConfig.(*schema.Set)
		added := newAccessLogConfig.(*schema.Set)

		if removed.Len() > 0 {
			request := map[string]interface{}{
				"ClientToken":    buildClientToken("DisableLoadBalancerAccessLog"),
				"LoadBalancerId": d.Id(),
			}

			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}

			action := "DisableLoadBalancerAccessLog"
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"OperationDenied.AccessLogEnabled", "SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
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

			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if added.Len() > 0 {
			request := map[string]interface{}{
				"ClientToken":    buildClientToken("EnableLoadBalancerAccessLog"),
				"LoadBalancerId": d.Id(),
			}

			if v, ok := d.GetOk("access_log_config"); ok {
				for _, accessLogConfigList := range v.(*schema.Set).List() {
					accessLogConfigArg := accessLogConfigList.(map[string]interface{})

					request["LogProject"] = accessLogConfigArg["log_project"]
					request["LogStore"] = accessLogConfigArg["log_store"]
				}
			}

			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}

			action := "EnableLoadBalancerAccessLog"
			conn, err := client.NewAlbClient()
			if err != nil {
				return WrapError(err)
			}

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"OperationDenied.AccessLogEnabled", "SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
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

			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		d.SetPartial("access_log_config")
	}

	if !d.IsNewResource() && d.HasChange("tags") {
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
	albServiceV2 := AlbServiceV2{client}
	action := "DeleteLoadBalancer"
	var response map[string]interface{}

	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":    buildClientToken("DeleteLoadBalancer"),
		"LoadBalancerId": d.Id(),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy", "IdempotenceProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.LoadBalancer"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{"Succeeded"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albServiceV2.AlbLoadBalancerJobStateRefreshFunc(d.Id(), response["JobId"].(string), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertAlbLoadBalancerBillingConfigPayTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPay"
	}

	return source
}

func convertAlbLoadBalancerBillingConfigPayTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPay":
		return "PayAsYouGo"
	}

	return source
}

func convertAlbLoadBalancerPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPay":
		return "PayAsYouGo"
	}

	return source
}

func convertAlbLoadBalancerAddressIpVersionResponse(source interface{}) interface{} {
	switch source {
	case "Ipv4":
		return "IPv4"
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

func convertAlbRegionIdResponse(source interface{}) interface{} {
	switch source {
	case "cn-hangzhou-onebox-nebula":
		return "cn-hangzhou"
	}

	return source
}

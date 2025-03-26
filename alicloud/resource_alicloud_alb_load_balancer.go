// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"hash/crc32"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
							Optional: true,
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_project": {
							Optional: true,
							Type:     schema.TypeString,
							Computed: true,
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
						},
						"enabled_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv6_address_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"load_balancer_billing_config": {
				Type:     schema.TypeList,
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
				Type:     schema.TypeString,
				Required: true,
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
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
						},
						"reason": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: modificationProtectionConfigDiffSuppressFunc,
						},
					},
				},
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
				Set: func(v interface{}) int {
					return int(crc32.ChecksumIEEE([]byte(v.(map[string]interface{})["zone_id"].(string))))
				},
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"intranet_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"address": {
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
						"eip_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"load_balancer_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv6_local_addresses": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"intranet_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"intranet_address_hc_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv4_local_addresses": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allocation_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6_address_hc_status": {
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
						"ipv6_address": {
							Type:     schema.TypeString,
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
		},
	}
}

func resourceAliCloudAlbLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("deletion_protection_enabled"); ok {
		request["DeletionProtectionEnabled"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("zone_mappings"); ok {
		zoneMappingsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.(*schema.Set).List() {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["ZoneId"] = dataLoop1Tmp["zone_id"]
			dataLoop1Map["AllocationId"] = dataLoop1Tmp["allocation_id"]
			dataLoop1Map["EipType"] = dataLoop1Tmp["eip_type"]
			dataLoop1Map["VSwitchId"] = dataLoop1Tmp["vswitch_id"]
			dataLoop1Map["IntranetAddress"] = dataLoop1Tmp["intranet_address"]
			zoneMappingsMapsArray = append(zoneMappingsMapsArray, dataLoop1Map)
		}
		request["ZoneMappings"] = zoneMappingsMapsArray
	}

	request["AddressType"] = d.Get("address_type")
	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIpVersion"] = v
	}
	request["LoadBalancerEdition"] = d.Get("load_balancer_edition")
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("modification_protection_config"); !IsNil(v) {
		status1, _ := jsonpath.Get("$[0].status", v)
		if status1 != nil && status1 != "" {
			objectDataLocalMap["Status"] = status1
		}
		reason1, _ := jsonpath.Get("$[0].reason", v)
		if reason1 != nil && reason1 != "" {
			objectDataLocalMap["Reason"] = reason1
		}

		request["ModificationProtectionConfig"] = objectDataLocalMap
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		objectDataLocalMap1["BandwidthPackageId"] = v
	}

	if v, ok := d.GetOk("load_balancer_billing_config"); ok {
		payType1, _ := jsonpath.Get("$[0].pay_type", v)
		if payType1 != nil && payType1 != "" {
			objectDataLocalMap1["PayType"] = convertAlbLoadBalancerBillingConfigPayTypeRequest(payType1)
		}
	}

	request["LoadBalancerBillingConfig"] = objectDataLocalMap1

	if v, ok := d.GetOk("address_allocated_mode"); ok {
		request["AddressAllocatedMode"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("deletion_protection_config"); ok {
		jsonPathResult7, err := jsonpath.Get("$[0].enabled", v)
		if err == nil && jsonPathResult7 != "" {
			request["DeletionProtectionEnabled"] = jsonPathResult7
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
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
	d.Set("address_ip_version", convertAlbLoadBalancerAddressIpVersionResponse(objectRaw["AddressIpVersion"]))
	d.Set("address_type", objectRaw["AddressType"])
	d.Set("bandwidth_package_id", objectRaw["BandwidthPackageId"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("dns_name", objectRaw["DNSName"])
	d.Set("ipv6_address_type", objectRaw["Ipv6AddressType"])
	d.Set("load_balancer_edition", objectRaw["LoadBalancerEdition"])
	d.Set("load_balancer_name", objectRaw["LoadBalancerName"])
	d.Set("region_id", convertAlbLoadBalancerRegionIdResponse(objectRaw["RegionId"]))
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["LoadBalancerStatus"])
	d.Set("vpc_id", objectRaw["VpcId"])

	accessLogConfigMaps := make([]map[string]interface{}, 0)
	accessLogConfigMap := make(map[string]interface{})
	accessLogConfigRaw := make(map[string]interface{})
	if objectRaw["AccessLogConfig"] != nil {
		accessLogConfigRaw = objectRaw["AccessLogConfig"].(map[string]interface{})
	}
	if len(accessLogConfigRaw) > 0 {
		accessLogConfigMap["log_project"] = accessLogConfigRaw["LogProject"]
		accessLogConfigMap["log_store"] = accessLogConfigRaw["LogStore"]

		accessLogConfigMaps = append(accessLogConfigMaps, accessLogConfigMap)
	}
	if err := d.Set("access_log_config", accessLogConfigMaps); err != nil {
		return err
	}
	deletionProtectionConfigMaps := make([]map[string]interface{}, 0)
	deletionProtectionConfigMap := make(map[string]interface{})
	deletionProtectionConfigRaw := make(map[string]interface{})
	if objectRaw["DeletionProtectionConfig"] != nil {
		deletionProtectionConfigRaw = objectRaw["DeletionProtectionConfig"].(map[string]interface{})
	}
	if len(deletionProtectionConfigRaw) > 0 {
		d.Set("deletion_protection_enabled", deletionProtectionConfigRaw["Enabled"])
		deletionProtectionConfigMap["enabled"] = deletionProtectionConfigRaw["Enabled"]
		deletionProtectionConfigMap["enabled_time"] = deletionProtectionConfigRaw["EnabledTime"]

		deletionProtectionConfigMaps = append(deletionProtectionConfigMaps, deletionProtectionConfigMap)
	}
	if err := d.Set("deletion_protection_config", deletionProtectionConfigMaps); err != nil {
		return err
	}
	loadBalancerBillingConfigMaps := make([]map[string]interface{}, 0)
	loadBalancerBillingConfigMap := make(map[string]interface{})
	loadBalancerBillingConfigRaw := make(map[string]interface{})
	if objectRaw["LoadBalancerBillingConfig"] != nil {
		loadBalancerBillingConfigRaw = objectRaw["LoadBalancerBillingConfig"].(map[string]interface{})
	}
	if len(loadBalancerBillingConfigRaw) > 0 {
		loadBalancerBillingConfigMap["pay_type"] = convertAlbLoadBalancerLoadBalancerBillingConfigPayTypeResponse(loadBalancerBillingConfigRaw["PayType"])

		loadBalancerBillingConfigMaps = append(loadBalancerBillingConfigMaps, loadBalancerBillingConfigMap)
	}
	if err := d.Set("load_balancer_billing_config", loadBalancerBillingConfigMaps); err != nil {
		return err
	}
	modificationProtectionConfigMaps := make([]map[string]interface{}, 0)
	modificationProtectionConfigMap := make(map[string]interface{})
	modificationProtectionConfigRaw := make(map[string]interface{})
	if objectRaw["ModificationProtectionConfig"] != nil {
		modificationProtectionConfigRaw = objectRaw["ModificationProtectionConfig"].(map[string]interface{})
	}
	if len(modificationProtectionConfigRaw) > 0 {
		modificationProtectionConfigMap["reason"] = modificationProtectionConfigRaw["Reason"]
		modificationProtectionConfigMap["status"] = modificationProtectionConfigRaw["Status"]

		modificationProtectionConfigMaps = append(modificationProtectionConfigMaps, modificationProtectionConfigMap)
	}
	if err := d.Set("modification_protection_config", modificationProtectionConfigMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	zoneMappingsRaw := objectRaw["ZoneMappings"]
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	if zoneMappingsRaw != nil {
		for _, zoneMappingsChildRaw := range zoneMappingsRaw.([]interface{}) {
			zoneMappingsMap := make(map[string]interface{})
			zoneMappingsChildRaw := zoneMappingsChildRaw.(map[string]interface{})
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

			zoneMappingsMap["address"] = loadBalancerAddressesChildRaw["Address"]
			zoneMappingsMap["allocation_id"] = loadBalancerAddressesChildRaw["AllocationId"]
			zoneMappingsMap["eip_type"] = loadBalancerAddressesChildRaw["EipType"]
			zoneMappingsMap["intranet_address"] = loadBalancerAddressesChildRaw["IntranetAddress"]
			zoneMappingsMap["ipv6_address"] = loadBalancerAddressesChildRaw["Ipv6Address"]

			loadBalancerAddressesRaw := zoneMappingsChildRaw["LoadBalancerAddresses"]
			loadBalancerAddressesMaps := make([]map[string]interface{}, 0)
			if loadBalancerAddressesRaw != nil {
				for _, loadBalancerAddressesChildRaw := range loadBalancerAddressesRaw.([]interface{}) {
					loadBalancerAddressesMap := make(map[string]interface{})
					loadBalancerAddressesChildRaw := loadBalancerAddressesChildRaw.(map[string]interface{})
					loadBalancerAddressesMap["address"] = loadBalancerAddressesChildRaw["Address"]
					loadBalancerAddressesMap["allocation_id"] = loadBalancerAddressesChildRaw["AllocationId"]
					loadBalancerAddressesMap["eip_type"] = loadBalancerAddressesChildRaw["EipType"]
					loadBalancerAddressesMap["intranet_address"] = loadBalancerAddressesChildRaw["IntranetAddress"]
					loadBalancerAddressesMap["intranet_address_hc_status"] = loadBalancerAddressesChildRaw["IntranetAddressHcStatus"]
					loadBalancerAddressesMap["ipv6_address"] = loadBalancerAddressesChildRaw["Ipv6Address"]
					loadBalancerAddressesMap["ipv6_address_hc_status"] = loadBalancerAddressesChildRaw["Ipv6AddressHcStatus"]

					ipv4LocalAddressesRaw := make([]interface{}, 0)
					if loadBalancerAddressesChildRaw["Ipv4LocalAddresses"] != nil {
						ipv4LocalAddressesRaw = loadBalancerAddressesChildRaw["Ipv4LocalAddresses"].([]interface{})
					}

					loadBalancerAddressesMap["ipv4_local_addresses"] = ipv4LocalAddressesRaw
					ipv6LocalAddressesRaw := make([]interface{}, 0)
					if loadBalancerAddressesChildRaw["Ipv6LocalAddresses"] != nil {
						ipv6LocalAddressesRaw = loadBalancerAddressesChildRaw["Ipv6LocalAddresses"].([]interface{})
					}

					loadBalancerAddressesMap["ipv6_local_addresses"] = ipv6LocalAddressesRaw
					loadBalancerAddressesMaps = append(loadBalancerAddressesMaps, loadBalancerAddressesMap)
				}
			}
			zoneMappingsMap["load_balancer_addresses"] = loadBalancerAddressesMaps
			zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
		}
	}
	if err := d.Set("zone_mappings", zoneMappingsMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudAlbLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if !d.IsNewResource() && d.HasChanges("deletion_protection_config.0.enabled", "deletion_protection_enabled") {
		albServiceV2 := AlbServiceV2{client}
		object, err := albServiceV2.DescribeAlbLoadBalancer(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("deletion_protection_config.0.enabled").(bool)
		if d.HasChange("deletion_protection_enabled") {
			target = d.Get("deletion_protection_enabled").(bool)
		}

		currentValue, err := jsonpath.Get("$.DeletionProtectionConfig.Enabled", object)
		if currentValue != nil && currentValue.(bool) != target {
			if target == false {
				action := "DisableDeletionProtection"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ResourceId"] = d.Id()

				request["ClientToken"] = buildClientToken(action)

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"SystemBusy", "IncorrectStatus.LoadBalancer", "IdempotenceProcessing"}) || NeedRetry(err) {
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
				albServiceV2 := AlbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == true {
				action := "EnableDeletionProtection"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ResourceId"] = d.Id()

				request["ClientToken"] = buildClientToken(action)

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"SystemBusy", "IncorrectStatus.LoadBalancer", "IdempotenceProcessing", "undefined"}) || NeedRetry(err) {
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
				albServiceV2 := AlbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}
	if d.HasChange("ipv6_address_type") {
		albServiceV2 := AlbServiceV2{client}
		object, err := albServiceV2.DescribeAlbLoadBalancer(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("ipv6_address_type").(string)
		if object["Ipv6AddressType"].(string) != target {
			if target == "Internet" {
				action := "EnableLoadBalancerIpv6Internet"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()

				request["ClientToken"] = buildClientToken(action)
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectStatus.LoadBalancer"}) || NeedRetry(err) {
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
				albServiceV2 := AlbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Internet"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "Ipv6AddressType", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Intranet" {
				action := "DisableLoadBalancerIpv6Internet"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["LoadBalancerId"] = d.Id()

				request["ClientToken"] = buildClientToken(action)
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectStatus.LoadBalancer"}) || NeedRetry(err) {
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
				albServiceV2 := AlbServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Intranet"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "Ipv6AddressType", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	action := "UpdateLoadBalancerAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("modification_protection_config") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("modification_protection_config"); v != nil {
			status1, _ := jsonpath.Get("$[0].status", v)
			if status1 != nil && (d.HasChange("modification_protection_config.0.status") || status1 != "") {
				objectDataLocalMap["Status"] = status1
			}
			reason1, _ := jsonpath.Get("$[0].reason", v)
			if reason1 != nil && (d.HasChange("modification_protection_config.0.reason") || reason1 != "") {
				objectDataLocalMap["Reason"] = reason1
			}

			request["ModificationProtectionConfig"] = objectDataLocalMap
		}
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if !d.IsNewResource() && d.HasChange("load_balancer_name") {
		update = true
		request["LoadBalancerName"] = d.Get("load_balancer_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"SystemBusy", "IncorrectStatus.LoadBalancer", "ResourceNotFound.LoadBalancer", "IdempotenceProcessing"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateLoadBalancerEdition"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if !d.IsNewResource() && d.HasChange("load_balancer_edition") {
		update = true
	}
	request["LoadBalancerEdition"] = d.Get("load_balancer_edition")

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"SystemBusy", "IncorrectStatus.LoadBalancer", "IdempotenceProcessing"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("load_balancer_edition"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerEdition", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "loadbalancer"

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.LoadBalancer", "undefined"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateLoadBalancerZones"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("zone_mappings") {
		update = true
	}
	if v, ok := d.GetOk("zone_mappings"); ok || d.HasChange("zone_mappings") {
		zoneMappingsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
			dataLoopMap["IntranetAddress"] = dataLoopTmp["intranet_address"]
			dataLoopMap["AllocationId"] = dataLoopTmp["allocation_id"]
			dataLoopMap["EipType"] = dataLoopTmp["eip_type"]
			dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
			zoneMappingsMapsArray = append(zoneMappingsMapsArray, dataLoopMap)
		}
		request["ZoneMappings"] = zoneMappingsMapsArray
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.LoadBalancer"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateLoadBalancerAddressTypeConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("address_type") {
		update = true
	}
	request["AddressType"] = d.Get("address_type")

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.LoadBalancer"}) || NeedRetry(err) {
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
		albServiceV2 := AlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		albServiceV2 := AlbServiceV2{client}
		if err := albServiceV2.SetResourceTags(d, "loadbalancer"); err != nil {
			return WrapError(err)
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

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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

			albServiceV2 := AlbServiceV2{client}
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

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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

			albServiceV2 := AlbServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albServiceV2.AlbLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		d.SetPartial("access_log_config")
	}

	d.Partial(false)
	return resourceAliCloudAlbLoadBalancerRead(d, meta)
}

func resourceAliCloudAlbLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"SystemBusy", "ResourceNotFound.LoadBalancer", "IdempotenceProcessing"}) || NeedRetry(err) {
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

	albServiceV2 := AlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"[Succeeded]"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albServiceV2.DescribeAsyncAlbLoadBalancerStateRefreshFunc(d, response, "$.Jobs[*].Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
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

func convertAlbLoadBalancerBillingConfigPayTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPay"
	}

	return source
}

func modificationProtectionConfigDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("modification_protection_config"); ok {
		val := v.([]interface{})
		if len(val) > 2 {
			// modification_protection_config 为 Object 类型
			return true
		}
		for _, modificationProtectionConfigs := range val {
			modificationProtectionConfigArg := modificationProtectionConfigs.(map[string]interface{})
			return fmt.Sprint(modificationProtectionConfigArg["status"]) != "ConsoleProtection"
		}
	}
	return true
}

func convertAlbLoadBalancerAddressIpVersionResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Ipv4":
		return "IPv4"
	}
	return source
}
func convertAlbLoadBalancerLoadBalancerBillingConfigPayTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PostPay":
		return "PayAsYouGo"
	}
	return source
}
func convertAlbLoadBalancerRegionIdResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "cn-hangzhou-onebox-nebula":
		return "cn-hangzhou"
	}
	return source
}
func convertAlbLoadBalancerLoadBalancerBillingConfigPayTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PayAsYouGo":
		return "PostPay"
	}
	return source
}

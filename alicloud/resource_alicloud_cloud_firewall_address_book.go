package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallAddressBook() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallAddressBookCreate,
		Read:   resourceAliCloudCloudFirewallAddressBookRead,
		Update: resourceAliCloudCloudFirewallAddressBookUpdate,
		Delete: resourceAliCloudCloudFirewallAddressBookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ip", "ipv6", "domain", "port", "tag", "asset", "assetIpv6"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_add_tag_ecs": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1}),
			},
			"tag_relation": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"and", "or"}, false),
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"zh", "en"}, false),
			},
			"address_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ecs_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"asset_member_uids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"asset_region_resource_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asset_region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_type": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"eip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ecs_eip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ecs_public_ip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"slb_eip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"slb_public_ip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"nlb_eip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"alb_eip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"nat_eip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"nat_public_ip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"eni_eip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ga_eip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"api_gateway_eip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ai_gateway_eip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"bastion_host_ip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"bastion_host_ingress_ip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"bastion_host_egress_ip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"havip": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"ipv6": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ecs_ipv6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"slb_ipv6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"nlb_ipv6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"alb_ipv6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"eni_eipv6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ga_eipv6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"api_gateway_eipv6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ai_gateway_eipv6": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"address_list_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"reference_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

var cloudFirewallAddressBookIpv4AssetTypesMap = map[string]string{
	"eip":                     "EIP",
	"ecs_eip":                 "EcsEIP",
	"ecs_public_ip":           "EcsPublicIP",
	"slb_eip":                 "SlbEIP",
	"slb_public_ip":           "SlbPublicIP",
	"nlb_eip":                 "NlbEIP",
	"alb_eip":                 "AlbEIP",
	"nat_eip":                 "NatEIP",
	"nat_public_ip":           "NatPublicIP",
	"eni_eip":                 "EniEIP",
	"ga_eip":                  "GaEIP",
	"api_gateway_eip":         "ApiGatewayEIP",
	"ai_gateway_eip":          "AiGatewayEIP",
	"bastion_host_ip":         "BastionHostIP",
	"bastion_host_ingress_ip": "BastionHostIngressIP",
	"bastion_host_egress_ip":  "BastionHostEgressIP",
	"havip":                   "HAVIP",
}

var cloudFirewallAddressBookIpv6AssetTypesMap = map[string]string{
	"ecs_ipv6":          "EcsIPv6",
	"slb_ipv6":          "SlbIPv6",
	"nlb_ipv6":          "NlbIPv6",
	"alb_ipv6":          "AlbIPv6",
	"eni_eipv6":         "EniEIPv6",
	"ga_eipv6":          "GaEIPv6",
	"api_gateway_eipv6": "ApiGatewayEIPv6",
	"ai_gateway_eipv6":  "AiGatewayEIPv6",
}

func expandCloudFirewallAddressBookAssetRegionResourceTypes(configured []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, item := range configured {
		if item == nil {
			continue
		}

		itemArg := item.(map[string]interface{})
		itemMap := make(map[string]interface{})

		if v, ok := itemArg["asset_region_id"].(string); ok && v != "" {
			itemMap["AssetRegionId"] = v
		}

		if resourceTypeList, ok := itemArg["resource_type"].([]interface{}); ok && len(resourceTypeList) > 0 && resourceTypeList[0] != nil {
			resourceTypeArg := resourceTypeList[0].(map[string]interface{})
			resourceTypeMap := make(map[string]interface{})

			if ipv4List, ok := resourceTypeArg["ipv4"].([]interface{}); ok && len(ipv4List) > 0 && ipv4List[0] != nil {
				ipv4Arg := ipv4List[0].(map[string]interface{})
				ipv4Map := make(map[string]interface{})
				for schemaKey, apiKey := range cloudFirewallAddressBookIpv4AssetTypesMap {
					if v, ok := ipv4Arg[schemaKey].(bool); ok {
						ipv4Map[apiKey] = v
					}
				}

				resourceTypeMap["Ipv4"] = ipv4Map
			}

			if ipv6List, ok := resourceTypeArg["ipv6"].([]interface{}); ok && len(ipv6List) > 0 && ipv6List[0] != nil {
				ipv6Arg := ipv6List[0].(map[string]interface{})
				ipv6Map := make(map[string]interface{})
				for schemaKey, apiKey := range cloudFirewallAddressBookIpv6AssetTypesMap {
					if v, ok := ipv6Arg[schemaKey].(bool); ok {
						ipv6Map[apiKey] = v
					}
				}

				resourceTypeMap["Ipv6"] = ipv6Map
			}

			itemMap["ResourceType"] = resourceTypeMap
		}

		result = append(result, itemMap)
	}

	return result
}

func flattenCloudFirewallAddressBookAssetRegionResourceTypes(src interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	srcList, ok := src.([]interface{})
	if !ok {
		return result
	}

	for _, item := range srcList {
		itemArg, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemMap := make(map[string]interface{})
		itemMap["asset_region_id"] = itemArg["AssetRegionId"]

		if resourceTypeArg, ok := itemArg["ResourceType"].(map[string]interface{}); ok {
			resourceTypeMap := make(map[string]interface{})

			// The service always returns both Ipv4 and Ipv6 maps with explicit booleans,
			// so a sub-block is exposed only when at least one asset type is enabled.
			if ipv4Arg, ok := resourceTypeArg["Ipv4"].(map[string]interface{}); ok {
				ipv4Map := make(map[string]interface{})
				ipv4Enabled := false
				for schemaKey, apiKey := range cloudFirewallAddressBookIpv4AssetTypesMap {
					if v, ok := ipv4Arg[apiKey]; ok {
						ipv4Map[schemaKey] = v
						if enabled, ok := v.(bool); ok && enabled {
							ipv4Enabled = true
						}
					}
				}

				if ipv4Enabled {
					resourceTypeMap["ipv4"] = []map[string]interface{}{ipv4Map}
				}
			}

			if ipv6Arg, ok := resourceTypeArg["Ipv6"].(map[string]interface{}); ok {
				ipv6Map := make(map[string]interface{})
				ipv6Enabled := false
				for schemaKey, apiKey := range cloudFirewallAddressBookIpv6AssetTypesMap {
					if v, ok := ipv6Arg[apiKey]; ok {
						ipv6Map[schemaKey] = v
						if enabled, ok := v.(bool); ok && enabled {
							ipv6Enabled = true
						}
					}
				}

				if ipv6Enabled {
					resourceTypeMap["ipv6"] = []map[string]interface{}{ipv6Map}
				}
			}

			itemMap["resource_type"] = []map[string]interface{}{resourceTypeMap}
		}

		result = append(result, itemMap)
	}

	return result
}

func resourceAliCloudCloudFirewallAddressBookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	var endpoint string
	action := "AddAddressBook"
	request := make(map[string]interface{})

	request["GroupName"] = d.Get("group_name")
	request["GroupType"] = d.Get("group_type")
	request["Description"] = d.Get("description")

	if v, ok := d.GetOkExists("auto_add_tag_ecs"); ok {
		request["AutoAddTagEcs"] = v
	}

	if v, ok := d.GetOk("tag_relation"); ok {
		request["TagRelation"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if v, ok := d.GetOk("address_list"); ok {
		request["AddressList"] = strings.Join(expandStringList(v.([]interface{})), ",")
	}

	if v, ok := d.GetOk("ecs_tags"); ok {
		for i, tagItem := range v.(*schema.Set).List() {
			tagItemArg := tagItem.(map[string]interface{})
			request[fmt.Sprintf("TagList.%d.TagValue", i+1)] = tagItemArg["tag_value"]
			request[fmt.Sprintf("TagList.%d.TagKey", i+1)] = tagItemArg["tag_key"]
		}
	}

	if v, ok := d.GetOk("asset_member_uids"); ok {
		assetMemberUidsJson, err := convertArrayObjectToJsonString(v.([]interface{}))
		if err != nil {
			return WrapError(err)
		}

		request["AssetMemberUids"] = assetMemberUidsJson
	}

	if v, ok := d.GetOk("asset_region_resource_types"); ok {
		assetRegionResourceTypesJson, err := convertArrayObjectToJsonString(expandCloudFirewallAddressBookAssetRegionResourceTypes(v.([]interface{})))
		if err != nil {
			return WrapError(err)
		}

		request["AssetRegionResourceTypes"] = assetRegionResourceTypesJson
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_address_book", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["GroupUuid"]))

	return resourceAliCloudCloudFirewallAddressBookRead(d, meta)
}

func resourceAliCloudCloudFirewallAddressBookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudFirewallAddressBook(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_address_book cloudfwService.DescribeCloudFirewallAddressBook Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("group_name", object["GroupName"])
	d.Set("group_type", object["GroupType"])
	d.Set("description", object["Description"])
	d.Set("tag_relation", object["TagRelation"])

	if v, ok := object["AutoAddTagEcs"]; ok {
		d.Set("auto_add_tag_ecs", formatInt(v))
	}

	addressListItems := make([]string, 0)
	for _, addressListArg := range object["AddressList"].([]interface{}) {
		addressListItems = append(addressListItems, fmt.Sprint(addressListArg))
	}
	d.Set("address_list", addressListItems)

	ecsTags := make([]map[string]interface{}, 0)
	for _, tagListItem := range object["TagList"].([]interface{}) {
		ecsTagItem := make(map[string]interface{})
		ecsTagItem["tag_value"] = tagListItem.(map[string]interface{})["TagValue"]
		ecsTagItem["tag_key"] = tagListItem.(map[string]interface{})["TagKey"]
		ecsTags = append(ecsTags, ecsTagItem)
	}

	d.Set("ecs_tags", ecsTags)

	if v, ok := object["AssetMemberUids"].([]interface{}); ok {
		assetMemberUids := make([]int, 0)
		for _, assetMemberUid := range v {
			assetMemberUids = append(assetMemberUids, formatInt(assetMemberUid))
		}

		d.Set("asset_member_uids", assetMemberUids)
	}

	if v, ok := object["AssetRegionResourceTypes"]; ok {
		d.Set("asset_region_resource_types", flattenCloudFirewallAddressBookAssetRegionResourceTypes(v))
	}

	if v, ok := object["AddressListCount"]; ok {
		d.Set("address_list_count", formatInt(v))
	}

	if v, ok := object["ReferenceCount"]; ok {
		d.Set("reference_count", formatInt(v))
	}

	return nil
}

func resourceAliCloudCloudFirewallAddressBookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"GroupUuid": d.Id(),
	}

	if d.HasChange("group_name") {
		update = true
	}
	request["GroupName"] = d.Get("group_name")

	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")

	if d.HasChange("auto_add_tag_ecs") {
		update = true

		if v, ok := d.GetOkExists("auto_add_tag_ecs"); ok {
			request["AutoAddTagEcs"] = v
		}
	}

	if d.HasChange("tag_relation") {
		update = true

		if v, ok := d.GetOk("tag_relation"); ok {
			request["TagRelation"] = v
		}
	}

	if d.HasChange("address_list") {
		update = true
		if v, ok := d.GetOk("address_list"); ok {
			request["AddressList"] = strings.Join(expandStringList(v.([]interface{})), ",")
		}
	}

	if d.HasChange("ecs_tags") {
		update = true

		if v, ok := d.GetOk("ecs_tags"); ok {
			for i, tagItem := range v.(*schema.Set).List() {
				tagItemArg := tagItem.(map[string]interface{})
				request[fmt.Sprintf("TagList.%d.TagValue", i+1)] = tagItemArg["tag_value"]
				request[fmt.Sprintf("TagList.%d.TagKey", i+1)] = tagItemArg["tag_key"]
			}
		}
	}

	if d.HasChange("asset_member_uids") {
		update = true

		if v, ok := d.GetOk("asset_member_uids"); ok {
			assetMemberUidsJson, err := convertArrayObjectToJsonString(v.([]interface{}))
			if err != nil {
				return WrapError(err)
			}

			request["AssetMemberUids"] = assetMemberUidsJson
		}
	}

	if d.HasChange("asset_region_resource_types") {
		update = true

		if v, ok := d.GetOk("asset_region_resource_types"); ok {
			assetRegionResourceTypesJson, err := convertArrayObjectToJsonString(expandCloudFirewallAddressBookAssetRegionResourceTypes(v.([]interface{})))
			if err != nil {
				return WrapError(err)
			}

			request["AssetRegionResourceTypes"] = assetRegionResourceTypesJson
		}
	}

	if update {
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}

		if v, ok := d.GetOk("source_ip"); ok {
			request["SourceIp"] = v
		}

		action := "ModifyAddressBook"
		var err error
		var endpoint string
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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

	return resourceAliCloudCloudFirewallAddressBookRead(d, meta)
}

func resourceAliCloudCloudFirewallAddressBookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAddressBook"
	var response map[string]interface{}
	var err error
	var endpoint string

	request := map[string]interface{}{
		"GroupUuid": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
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

	return nil
}

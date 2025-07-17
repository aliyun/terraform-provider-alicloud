package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcVpcCreate,
		Read:   resourceAliCloudVpcVpcRead,
		Update: resourceAliCloudVpcVpcUpdate,
		Delete: resourceAliCloudVpcVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"classic_link_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_hostname_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ENABLED", "DISABLED", "MODIFYING"}, false),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv4_cidr_mask": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ipv4_ipam_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ipv6_cidr_blocks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv6_isp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ipv6_isp": {
				Type:     schema.TypeString,
				Optional: true,
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
			"route_table_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"router_table_id"},
				Computed:      true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_cidr_blocks": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'secondary_cidr_blocks' has been deprecated from provider version 1.185.0. Field 'secondary_cidr_blocks' has been deprecated from provider version 1.185.0 and it will be removed in the future version. Please use the new resource 'alicloud_vpc_ipv4_cidr_block'. `secondary_cidr_blocks` attributes and `alicloud_vpc_ipv4_cidr_block` resource cannot be used at the same time.",
				Elem:       &schema.Schema{Type: schema.TypeString},
			},
			"secondary_cidr_mask": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'secondary_cidr_mask' has been deprecated from provider version 1.248.0. Field 'secondary_cidr_blocks' has been deprecated from provider version 1.248.0 and it will be removed in the future version. Please use the new resource 'alicloud_vpc_ipv4_cidr_block'. `secondary_cidr_mask` attributes and `alicloud_vpc_ipv4_cidr_block` resource cannot be used at the same time.",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_route_table_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_route_table_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_route_table_route_propagation_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"user_cidrs": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.119.0. New field 'vpc_name' instead.",
			},
			"router_table_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'router_table_id' has been deprecated since provider version 1.221.0. New field 'route_table_id' instead.",
			},
		},
	}
}

func resourceAliCloudVpcVpcCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	if v := d.Get("is_default"); !v.(bool) {

		action := "CreateVpc"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)

		if v, ok := d.GetOk("cidr_block"); ok {
			request["CidrBlock"] = v
		}
		if v, ok := d.GetOk("name"); ok || d.HasChange("name") {
			request["VpcName"] = v
		}

		if v, ok := d.GetOk("vpc_name"); ok {
			request["VpcName"] = v
		}
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
		if v, ok := d.GetOk("user_cidrs"); ok {
			jsonPathResult4, err := jsonpath.Get("$", v)
			if err == nil && jsonPathResult4 != "" {
				request["UserCidr"] = convertListToCommaSeparate(jsonPathResult4.([]interface{}))
			}
		}
		if v, ok := d.GetOk("ipv6_isp"); ok {
			request["Ipv6Isp"] = v
		}
		if v, ok := d.GetOkExists("enable_ipv6"); ok {
			request["EnableIpv6"] = v
		}
		if v, ok := d.GetOk("ipv4_ipam_pool_id"); ok {
			request["Ipv4IpamPoolId"] = v
		}
		if v, ok := d.GetOk("tags"); ok {
			tagsMap := ConvertTags(v.(map[string]interface{}))
			request = expandTagsToMap(request, tagsMap)
		}

		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		if v, ok := d.GetOk("dns_hostname_status"); ok {
			request["EnableDnsHostname"] = convertVpcVpcEnableDnsHostnameRequest(v.(string))
		}
		if v, ok := d.GetOkExists("ipv4_cidr_mask"); ok {
			request["Ipv4CidrMask"] = v
		}
		if v, ok := d.GetOk("ipv6_cidr_block"); ok {
			request["Ipv6CidrBlock"] = v
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprint(response["VpcId"]))

		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	if v, ok := d.GetOk("is_default"); ok && InArray(fmt.Sprint(v), []string{"true"}) {
		action := "CreateDefaultVpc"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["ClientToken"] = buildClientToken(action)

		if v, ok := d.GetOk("ipv6_cidr_block"); ok {
			request["Ipv6CidrBlock"] = v
		}
		if v, ok := d.GetOkExists("enable_ipv6"); ok {
			request["EnableIpv6"] = v
		}
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"TaskConflict"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprint(response["VpcId"]))

	}

	return resourceAliCloudVpcVpcUpdate(d, meta)
}

func resourceAliCloudVpcVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcVpc(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc DescribeVpcVpc Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cidr_block", objectRaw["CidrBlock"])
	d.Set("classic_link_enabled", objectRaw["ClassicLinkEnabled"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("dns_hostname_status", objectRaw["DnsHostnameStatus"])
	d.Set("enable_ipv6", objectRaw["EnabledIpv6"])
	d.Set("ipv6_cidr_block", objectRaw["Ipv6CidrBlock"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("router_id", objectRaw["VRouterId"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_name", objectRaw["VpcName"])

	ipv6CidrBlockRaw, _ := jsonpath.Get("$.Ipv6CidrBlocks.Ipv6CidrBlock", objectRaw)
	ipv6CidrBlocksMaps := make([]map[string]interface{}, 0)
	if ipv6CidrBlockRaw != nil {
		for _, ipv6CidrBlockChildRaw := range ipv6CidrBlockRaw.([]interface{}) {
			ipv6CidrBlocksMap := make(map[string]interface{})
			ipv6CidrBlockChildRaw := ipv6CidrBlockChildRaw.(map[string]interface{})
			ipv6CidrBlocksMap["ipv6_cidr_block"] = ipv6CidrBlockChildRaw["Ipv6CidrBlock"]
			ipv6CidrBlocksMap["ipv6_isp"] = ipv6CidrBlockChildRaw["Ipv6Isp"]

			ipv6CidrBlocksMaps = append(ipv6CidrBlocksMaps, ipv6CidrBlocksMap)
		}
	}
	if err := d.Set("ipv6_cidr_blocks", ipv6CidrBlocksMaps); err != nil {
		return err
	}
	secondaryCidrBlockRaw, _ := jsonpath.Get("$.SecondaryCidrBlocks.SecondaryCidrBlock", objectRaw)
	d.Set("secondary_cidr_blocks", secondaryCidrBlockRaw)
	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))
	userCidrRaw, _ := jsonpath.Get("$.UserCidrs.UserCidr", objectRaw)
	d.Set("user_cidrs", userCidrRaw)

	objectRaw, err = vpcServiceV2.DescribeVpcDescribeRouteTableList(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("route_table_id", objectRaw["RouteTableId"])
	d.Set("router_id", objectRaw["RouterId"])
	d.Set("system_route_table_description", objectRaw["Description"])
	d.Set("system_route_table_name", objectRaw["RouteTableName"])
	d.Set("system_route_table_route_propagation_enable", objectRaw["RoutePropagationEnable"])

	d.Set("name", d.Get("vpc_name"))
	d.Set("router_table_id", d.Get("route_table_id"))
	return nil
}

func resourceAliCloudVpcVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("classic_link_enabled") {
		var err error
		vpcServiceV2 := VpcServiceV2{client}
		object, err := vpcServiceV2.DescribeVpcVpc(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("classic_link_enabled").(bool)
		if _, ok := object["ClassicLinkEnabled"]; ok && object["ClassicLinkEnabled"].(bool) != target {
			if target == true {
				action := "EnableVpcClassicLink"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["VpcId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectVpcStatus"}) || NeedRetry(err) {
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
			if target == false {
				action := "DisableVpcClassicLink"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["VpcId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"InternalError", "IncorrectVpcStatus"}) || NeedRetry(err) {
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
	action := "ModifyVpcAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VpcId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["VpcName"] = d.Get("name")
	}

	if !d.IsNewResource() && d.HasChange("vpc_name") {
		update = true
		request["VpcName"] = d.Get("vpc_name")
	}

	if !d.IsNewResource() && d.HasChange("cidr_block") {
		update = true
		request["CidrBlock"] = d.Get("cidr_block")
	}

	if !d.IsNewResource() && d.HasChange("enable_ipv6") {
		update = true
		request["EnableIPv6"] = d.Get("enable_ipv6")
	}

	if v, ok := d.GetOk("ipv6_isp"); ok {
		request["Ipv6Isp"] = v
	}
	if !d.IsNewResource() && d.HasChange("dns_hostname_status") {
		update = true
		request["EnableDnsHostname"] = convertVpcVpcEnableDnsHostnameRequest(d.Get("dns_hostname_status").(string))
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationFailed.LastTokenProcessing", "LastTokenProcessing", "OperationFailed.QueryCenIpv6Status", "IncorrectStatus", "OperationConflict", "SystemBusy", "ServiceUnavailable", "IncorrectVpcStatus"}) || NeedRetry(err) {
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
		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		if d.HasChange("dns_hostname_status") {
			stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("dns_hostname_status"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "DnsHostnameStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "VPC"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
	action = "ModifyRouteTableAttributes"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	vpcServiceV2 := VpcServiceV2{client}
	objectRaw, err := vpcServiceV2.DescribeVpcDescribeRouteTableList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if objectRaw["RouteTableId"] != nil {
		request["RouteTableId"] = objectRaw["RouteTableId"]
	}

	if d.HasChange("system_route_table_description") {
		update = true
		request["Description"] = d.Get("system_route_table_description")
	}

	if d.HasChange("router_table_id") {
		update = true
		request["RouteTableId"] = d.Get("router_table_id")
	}

	if d.HasChange("route_table_id") {
		update = true
		request["RouteTableId"] = d.Get("route_table_id")
	}

	if d.HasChange("system_route_table_name") {
		update = true
		request["RouteTableName"] = d.Get("system_route_table_name")
	}

	if v, ok := d.GetOkExists("system_route_table_route_propagation_enable"); (d.IsNewResource() && ok && !v.(bool)) || (!d.IsNewResource() && d.HasChange("system_route_table_route_propagation_enable")) {
		update = true
		request["RoutePropagationEnable"] = d.Get("system_route_table_route_propagation_enable")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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

	if d.HasChange("secondary_cidr_blocks") {
		oldEntry, newEntry := d.GetChange("secondary_cidr_blocks")
		removed := oldEntry
		added := newEntry

		if len(removed.([]interface{})) > 0 {
			secondaryCidrBlocks := removed.([]interface{})

			for _, item := range secondaryCidrBlocks {
				action := "UnassociateVpcCidrBlock"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["VpcId"] = d.Id()
				request["RegionId"] = client.RegionId
				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					request["SecondaryCidrBlock"] = convertListToCommaSeparate(jsonPathResult.([]interface{}))
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
				vpcServiceV2 := VpcServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Created", "Available"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}

		if len(added.([]interface{})) > 0 {
			secondaryCidrBlocks := added.([]interface{})

			for _, item := range secondaryCidrBlocks {
				action := "AssociateVpcCidrBlock"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["VpcId"] = d.Id()
				request["RegionId"] = client.RegionId
				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					request["SecondaryCidrBlock"] = jsonPathResult
				}
				request["IpVersion"] = "IPV4"
				if v, ok := item.(string); ok {
					_request := make(map[string]interface{})
					_request["IpamPoolId"] = v
				}
				if v, ok := d.GetOkExists("secondary_cidr_mask"); ok {
					request["SecondaryCidrMask"] = v
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
				vpcServiceV2 := VpcServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Created", "Available"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}
	if d.HasChange("tags") {
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "VPC"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudVpcVpcRead(d, meta)
}

func resourceAliCloudVpcVpcDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpc"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["VpcId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("force_delete"); ok {
		request["ForceDelete"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation.VSwitch", "DependencyViolation.SecurityGroup", "IncorrectVpcStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidResource.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertVpcVpcEnableDnsHostnameRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "ENABLED":
		return "true"
	case "DISABLED":
		return "false"
	}
	return source
}

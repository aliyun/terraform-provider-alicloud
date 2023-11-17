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
			Create: schema.DefaultTimeout(5 * time.Minute),
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
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
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
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_cidr_blocks": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'SecondaryCidrBlocks' has been deprecated from provider version 1.206.0. Field 'secondary_cidr_blocks' has been deprecated from provider version 1.185.0 and it will be removed in the future version. Please use the new resource 'alicloud_vpc_ipv4_cidr_block'. `secondary_cidr_blocks` attributes and `alicloud_vpc_ipv4_cidr_block` resource cannot be used at the same time.",
				Elem:       &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
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
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated from provider version 1.119.0. New field 'vpc_name' instead.",
			},
			"router_table_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'router_table_id' has been deprecated from provider version 1.206.0. New field 'route_table_id' instead.",
			},
		},
	}
}

func resourceAliCloudVpcVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpc"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("cidr_block"); ok {
		request["CidrBlock"] = v
	}
	if v, ok := d.GetOk("name"); ok {
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
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("user_cidrs"); ok {
		jsonPathResult5, err := jsonpath.Get("$", v)
		if err != nil {
			return WrapError(err)
		}
		request["UserCidr"] = convertListToCommaSeparate(jsonPathResult5.([]interface{}))
	}
	if v, ok := d.GetOk("ipv6_isp"); ok {
		request["Ipv6Isp"] = v
	}
	if v, ok := d.GetOkExists("enable_ipv6"); ok {
		request["EnableIpv6"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VpcId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 0, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
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
	d.Set("ipv6_cidr_block", objectRaw["Ipv6CidrBlock"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("router_id", objectRaw["VRouterId"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_name", objectRaw["VpcName"])
	ipv6CidrBlock7Raw, _ := jsonpath.Get("$.Ipv6CidrBlocks.Ipv6CidrBlock", objectRaw)
	ipv6CidrBlocksMaps := make([]map[string]interface{}, 0)
	if ipv6CidrBlock7Raw != nil {
		for _, ipv6CidrBlockChild1Raw := range ipv6CidrBlock7Raw.([]interface{}) {
			ipv6CidrBlocksMap := make(map[string]interface{})
			ipv6CidrBlockChild1Raw := ipv6CidrBlockChild1Raw.(map[string]interface{})
			ipv6CidrBlocksMap["ipv6_cidr_block"] = ipv6CidrBlockChild1Raw["Ipv6CidrBlock"]
			ipv6CidrBlocksMap["ipv6_isp"] = ipv6CidrBlockChild1Raw["Ipv6Isp"]
			ipv6CidrBlocksMaps = append(ipv6CidrBlocksMaps, ipv6CidrBlocksMap)
		}
	}
	d.Set("ipv6_cidr_blocks", ipv6CidrBlocksMaps)
	secondaryCidrBlock1Raw, _ := jsonpath.Get("$.SecondaryCidrBlocks.SecondaryCidrBlock", objectRaw)
	d.Set("secondary_cidr_blocks", secondaryCidrBlock1Raw)
	userCidr1Raw, _ := jsonpath.Get("$.UserCidrs.UserCidr", objectRaw)
	d.Set("user_cidrs", userCidr1Raw)

	objectRaw, err = vpcServiceV2.DescribeDescribeRouteTableList(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("route_table_id", objectRaw["RouteTableId"])
	d.Set("router_id", objectRaw["RouterId"])
	d.Set("status", objectRaw["Status"])

	objectRaw, err = vpcServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("name", d.Get("vpc_name"))
	d.Set("router_table_id", d.Get("route_table_id"))
	return nil
}

func resourceAliCloudVpcVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyVpcAttribute"
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

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

	if v, ok := d.GetOkExists("enable_ipv6"); ok {
		request["EnableIPv6"] = v
	}
	if v, ok := d.GetOk("ipv6_isp"); ok {
		request["Ipv6Isp"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

			if err != nil {
				if IsExpectedErrors(err, []string{"OperationFailed.LastTokenProcessing", "LastTokenProcessing", "OperationFailed.QueryCenIpv6Status", "IncorrectStatus", "OperationConflict", "SystemBusy", "ServiceUnavailable", "IncorrectVpcStatus"}) || NeedRetry(err) {
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
		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("description")
		d.SetPartial("vpc_name")
		d.SetPartial("cidr_block")
	}
	update = false
	action = "MoveResourceGroup"
	conn, err = client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "VPC"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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

	if d.HasChange("classic_link_enabled") {
		client := meta.(*connectivity.AliyunClient)
		vpcServiceV2 := VpcServiceV2{client}
		object, err := vpcServiceV2.DescribeVpcVpc(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("classic_link_enabled").(bool)
		if object["ClassicLinkEnabled"].(bool) != target {
			if target == true {
				action = "EnableVpcClassicLink"
				conn, err = client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})

				request["VpcId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					request["ClientToken"] = buildClientToken(action)

					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectVpcStatus"}) || NeedRetry(err) {
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
			if target == false {
				action = "DisableVpcClassicLink"
				conn, err = client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})

				request["VpcId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					request["ClientToken"] = buildClientToken(action)

					if err != nil {
						if IsExpectedErrors(err, []string{"InternalError", "IncorrectVpcStatus"}) || NeedRetry(err) {
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

	update = false
	if d.HasChange("secondary_cidr_blocks") {
		update = true
		oldEntry, newEntry := d.GetChange("secondary_cidr_blocks")
		removed := oldEntry
		added := newEntry

		if len(removed.([]interface{})) > 0 {
			secondaryCidrBlocks := removed.([]interface{})

			for _, item := range secondaryCidrBlocks {
				action = "UnassociateVpcCidrBlock"
				conn, err = client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})

				request["VpcId"] = d.Id()
				request["RegionId"] = client.RegionId
				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					request["SecondaryCidrBlock"] = jsonPathResult
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
				vpcServiceV2 := VpcServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Created", "Available"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
				d.SetPartial("secondary_cidr_blocks")

			}
			d.SetPartial("secondary_cidr_blocks")
		}

		if len(added.([]interface{})) > 0 {
			secondaryCidrBlocks := added.([]interface{})

			for _, item := range secondaryCidrBlocks {
				action = "AssociateVpcCidrBlock"
				conn, err = client.NewVpcClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})

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
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
				vpcServiceV2 := VpcServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Created", "Available"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
				d.SetPartial("secondary_cidr_blocks")

			}
			d.SetPartial("secondary_cidr_blocks")
		}
	}
	update = false
	if d.HasChange("tags") {
		update = true
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "VPC"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudVpcVpcRead(d, meta)
}

func resourceAliCloudVpcVpcDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteVpc"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["VpcId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.LastTokenProcessing", "LastTokenProcessing", "OperationFailed.QueryCenIpv6Status", "IncorrectStatus", "OperationConflict", "SystemBusy", "ServiceUnavailable", "IncorrectVpcStatus"}) || NeedRetry(err) {
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

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 0, vpcServiceV2.VpcVpcStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

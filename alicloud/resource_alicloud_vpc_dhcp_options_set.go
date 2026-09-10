// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcDhcpOptionsSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcDhcpOptionsSetCreate,
		Read:   resourceAliCloudVpcDhcpOptionsSetRead,
		Update: resourceAliCloudVpcDhcpOptionsSetUpdate,
		Delete: resourceAliCloudVpcDhcpOptionsSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"associate_vpcs": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'associate_vpcs' has been deprecated from provider version 1.153.0. Field 'associate_vpcs' has been deprecated from provider version 1.153.0 and it will be removed in the future version. Please use the new resource 'alicloud_vpc_dhcp_options_set_attachment' to attach DhcpOptionsSet and Vpc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"associate_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"dhcp_options_set_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dhcp_options_set_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[a-zA-Z\u4E00-\u9FA5][\u4E00-\u9FA5A-Za-z0-9_-]{2,128}$"), "The name must be 2 to 128 characters in length and can contain letters, Chinese characters, digits, underscores (_), and hyphens (-). It must start with a letter or a Chinese character."),
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name_servers": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ipv6_lease_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lease_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeInt,
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
		},
	}
}

func resourceAliCloudVpcDhcpOptionsSetCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDhcpOptionsSet"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("domain_name_servers"); ok {
		request["DomainNameServers"] = v
	}
	if v, ok := d.GetOk("dhcp_options_set_name"); ok {
		request["DhcpOptionsSetName"] = v
	}
	if v, ok := d.GetOk("dhcp_options_set_description"); ok {
		request["DhcpOptionsSetDescription"] = v
	}
	if v, ok := d.GetOk("domain_name"); ok {
		request["DomainName"] = v
	}
	if v, ok := d.GetOk("lease_time"); ok {
		request["LeaseTime"] = v
	}
	if v, ok := d.GetOk("ipv6_lease_time"); ok {
		request["Ipv6LeaseTime"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_dhcp_options_set", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DhcpOptionsSetId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available", "InUse"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcDhcpOptionsSetStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcDhcpOptionsSetUpdate(d, meta)
}

func resourceAliCloudVpcDhcpOptionsSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcDhcpOptionsSet(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_dhcp_options_set DescribeVpcDhcpOptionsSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("dhcp_options_set_description", objectRaw["DhcpOptionsSetDescription"])
	d.Set("dhcp_options_set_name", objectRaw["DhcpOptionsSetName"])
	d.Set("owner_id", objectRaw["OwnerId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])

	dhcpOptions1RawObj, _ := jsonpath.Get("$.DhcpOptions", objectRaw)
	dhcpOptions1Raw := make(map[string]interface{})
	if dhcpOptions1RawObj != nil {
		dhcpOptions1Raw = dhcpOptions1RawObj.(map[string]interface{})
	}
	d.Set("domain_name", dhcpOptions1Raw["DomainName"])
	d.Set("domain_name_servers", dhcpOptions1Raw["DomainNameServers"])
	d.Set("ipv6_lease_time", dhcpOptions1Raw["Ipv6LeaseTime"])
	d.Set("lease_time", dhcpOptions1Raw["LeaseTime"])

	associateVpcs1Raw := objectRaw["AssociateVpcs"]
	associateVpcsMaps := make([]map[string]interface{}, 0)
	if associateVpcs1Raw != nil {
		for _, associateVpcsChild1Raw := range associateVpcs1Raw.([]interface{}) {
			associateVpcsMap := make(map[string]interface{})
			associateVpcsChild1Raw := associateVpcsChild1Raw.(map[string]interface{})
			associateVpcsMap["associate_status"] = associateVpcsChild1Raw["AssociateStatus"]
			associateVpcsMap["vpc_id"] = associateVpcsChild1Raw["VpcId"]

			associateVpcsMaps = append(associateVpcsMaps, associateVpcsMap)
		}
	}
	d.Set("associate_vpcs", associateVpcsMaps)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudVpcDhcpOptionsSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateDhcpOptionsSetAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DhcpOptionsSetId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("domain_name_servers") {
		update = true
		request["DomainNameServers"] = d.Get("domain_name_servers")
	}

	if !d.IsNewResource() && d.HasChange("domain_name") {
		update = true
		request["DomainName"] = d.Get("domain_name")
	}

	if !d.IsNewResource() && d.HasChange("dhcp_options_set_name") {
		update = true
		request["DhcpOptionsSetName"] = d.Get("dhcp_options_set_name")
	}

	if !d.IsNewResource() && d.HasChange("dhcp_options_set_description") {
		update = true
		request["DhcpOptionsSetDescription"] = d.Get("dhcp_options_set_description")
	}

	if !d.IsNewResource() && d.HasChange("lease_time") {
		update = true
		request["LeaseTime"] = d.Get("lease_time")
	}

	if !d.IsNewResource() && d.HasChange("ipv6_lease_time") {
		update = true
		request["Ipv6LeaseTime"] = d.Get("ipv6_lease_time")
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available", "InUse"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcDhcpOptionsSetStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "DhcpOptionsSet"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, false)

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

	if d.HasChange("associate_vpcs") {
		oldEntry, newEntry := d.GetChange("associate_vpcs")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			associateVpcs := removed.List()

			for _, item := range associateVpcs {
				action := "DetachDhcpOptionsSetFromVpc"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["DhcpOptionsSetId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				if v, ok := item.(map[string]interface{}); ok {
					jsonPathResult, err := jsonpath.Get("$.vpc_id", v)
					if err != nil {
						return WrapError(err)
					}
					request["VpcId"] = jsonPathResult
				}
				if v, ok := item.(map[string]interface{}); ok {
					request["DryRun"] = v
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
				vpcServiceV2 := VpcServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Available", "InUse"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, vpcServiceV2.VpcDhcpOptionsSetStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}

		if added.Len() > 0 {
			associateVpcs := added.List()

			for _, item := range associateVpcs {
				action := "AttachDhcpOptionsSetToVpc"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["DhcpOptionsSetId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				if v, ok := item.(map[string]interface{}); ok {
					jsonPathResult, err := jsonpath.Get("$.vpc_id", v)
					if err != nil {
						return WrapError(err)
					}
					request["VpcId"] = jsonPathResult
				}
				if v, ok := item.(map[string]interface{}); ok {
					request["DryRun"] = v
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
				vpcServiceV2 := VpcServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Available", "InUse"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, vpcServiceV2.VpcDhcpOptionsSetStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}

	}
	if d.HasChange("tags") {
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "DhcpOptionsSet"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudVpcDhcpOptionsSetRead(d, meta)
}

func resourceAliCloudVpcDhcpOptionsSetDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDhcpOptionsSet"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DhcpOptionsSetId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.DhcpOptionsSet", "DependencyViolation.VpcAttachment", "OperationFailed.LastTokenProcessing", "LastTokenProcessing", "OperationConflict", "IncorrectStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDhcpOptionsSetId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcDhcpOptionsSetStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

package ga

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudGaAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaAclCreate,
		Read:   resourceAliCloudGaAclRead,
		Update: resourceAliCloudGaAclUpdate,
		Delete: resourceAliCloudGaAclDelete,
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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"IPv4", "IPv6"}, false),
			},
			"acl_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringMatch(regexp.MustCompile(`^[a-zA-Z][A-Za-z0-9._-]{2,128}$`), "The name must be `2` to `128` characters in length, and can contain letters, digits, periods (.), hyphens (-) and underscores (_). It must start with a letter."),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"acl_entries": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field `acl_entries` has been deprecated from provider version 1.190.0 and it will be removed in the future version. Please use the new resource `alicloud_ga_acl_entry_attachment`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"entry_description": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validators.StringMatch(regexp.MustCompile(`^[A-Za-z0-9._/-]{1,256}$`), "The description of the IP entry. The description must be 1 to 256 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.),and underscores (_)."),
						},
					},
				},
			},
			"tags": tags.TagsSchema(),
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	action := "CreateAcl"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("CreateAcl")
	request["AddressIPVersion"] = d.Get("address_ip_version")

	if v, ok := d.GetOk("acl_name"); ok {
		request["AclName"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if m, ok := d.GetOk("acl_entries"); ok {
		for k, aclEntries := range m.(*schema.Set).List() {
			aclEntriesArg := aclEntries.(map[string]interface{})
			request[fmt.Sprintf("AclEntries.%d.Entry", k+1)] = aclEntriesArg["entry"].(string)
			request[fmt.Sprintf("AclEntries.%d.EntryDescription", k+1)] = aclEntriesArg["entry_description"].(string)
		}
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_acl", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AclId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudGaAclUpdate(d, meta)
}

func resourceAliCloudGaAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	object, err := gaService.DescribeGaAcl(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_acl gaService.DescribeGaAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("address_ip_version", object["AddressIPVersion"])
	d.Set("acl_name", object["AclName"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["AclStatus"])

	if v, ok := object["AclEntries"].([]interface{}); ok {
		aclEntries := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"entry":             item["Entry"],
				"entry_description": item["EntryDescription"],
			}
			aclEntries = append(aclEntries, temp)
		}
		if err := d.Set("acl_entries", aclEntries); err != nil {
			return errs.WrapError(err)
		}
	}

	listTagResourcesObject, err := gaService.ListTagResources(d.Id(), "acl")
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("tags", tags.TagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudGaAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := gaService.SetResourceTags(d, "acl"); err != nil {
			return errs.WrapError(err)
		}
	}

	update := false
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken("UpdateAclAttribute"),
		"AclId":       d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("acl_name") {
		update = true
	}
	if v, ok := d.GetOk("acl_name"); ok {
		request["AclName"] = v
	}

	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}

		action := "UpdateAclAttribute"
		var err error

		wait := helper.IncrementalWait(3*time.Second, 20*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"StateError.Acl"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

	}

	if !d.IsNewResource() && d.HasChange("acl_entries") {
		oraw, nraw := d.GetChange("acl_entries")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()

		if len(remove) > 0 {
			removeEntriesFromAclReq := map[string]interface{}{
				"RegionId":    client.RegionId,
				"ClientToken": helper.BuildClientToken("RemoveEntriesFromAcl"),
				"AclId":       d.Id(),
			}

			for k, aclEntries := range remove {
				aclEntriesArg := aclEntries.(map[string]interface{})
				removeEntriesFromAclReq[fmt.Sprintf("AclEntries.%d.Entry", k+1)] = aclEntriesArg["entry"].(string)
			}

			if v, ok := d.GetOkExists("dry_run"); ok {
				removeEntriesFromAclReq["DryRun"] = v
			}

			action := "RemoveEntriesFromAcl"
			var err error

			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, removeEntriesFromAclReq, true)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"StateError.Acl"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, removeEntriesFromAclReq)

			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

			stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}

		if len(create) > 0 {
			addEntriesToAclReq := map[string]interface{}{
				"RegionId":    client.RegionId,
				"ClientToken": helper.BuildClientToken("AddEntriesToAcl"),
				"AclId":       d.Id(),
			}

			for k, aclEntries := range create {
				aclEntriesArg := aclEntries.(map[string]interface{})
				addEntriesToAclReq[fmt.Sprintf("AclEntries.%d.Entry", k+1)] = aclEntriesArg["entry"].(string)
				addEntriesToAclReq[fmt.Sprintf("AclEntries.%d.EntryDescription", k+1)] = aclEntriesArg["entry_description"].(string)
			}

			if v, ok := d.GetOkExists("dry_run"); ok {
				addEntriesToAclReq["DryRun"] = v
			}

			action := "AddEntriesToAcl"
			var err error

			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, addEntriesToAclReq, true)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"StateError.Acl"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, addEntriesToAclReq)

			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

			stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}

	}

	update = false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ClientToken":  helper.BuildClientToken("ChangeResourceGroup"),
		"ResourceId":   d.Id(),
		"ResourceType": "acl",
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		changeResourceGroupReq["NewResourceGroupId"] = v
	}

	if update {
		action := "ChangeResourceGroup"
		var err error
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, changeResourceGroupReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, changeResourceGroupReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

	}

	d.Partial(false)

	return resourceAliCloudGaAclRead(d, meta)
}

func resourceAliCloudGaAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	action := "DeleteAcl"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken("DeleteAcl"),
		"AclId":       d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 20*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Acl", "AclHasBindListener"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotExist.Acl"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

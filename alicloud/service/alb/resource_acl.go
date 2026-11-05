package alb

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAclCreate,
		Read:   resourceAclRead,
		Update: resourceAclUpdate,
		Delete: resourceAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(16 * time.Minute),
			Delete: schema.DefaultTimeout(16 * time.Minute),
			Update: schema.DefaultTimeout(16 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_entries": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'acl_entries' has been deprecated from provider version 1.166.0 and it will be removed in the future version. Please use the new resource 'alicloud_alb_acl_entry_attachment'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
						"entry": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"acl_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.TagsSchemaWithElements(),
		},
	}
}

func resourceAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAcl"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("acl_name"); ok {
		request["AclName"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ClientToken"] = helper.BuildClientToken("CreateAcl")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "SystemBusy", "Throttling"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alb_acl", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AclId"]))
	albService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAclUpdate(d, meta)
}

func resourceAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := NewService(client)
	object, err := albService.DescribeAcl(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_acl albService.DescribeAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("acl_name", object["AclName"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["AclStatus"])
	listAclEntriesObject, err := albService.ListAclEntries(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	aclEntriesMaps := make([]map[string]interface{}, 0)
	for _, aclEntriesListItem := range listAclEntriesObject {
		aclEntriesArg := make(map[string]interface{}, 0)
		aclEntriesArg["description"] = aclEntriesListItem["Description"]
		aclEntriesArg["entry"] = aclEntriesListItem["Entry"]
		aclEntriesArg["status"] = aclEntriesListItem["Status"]
		aclEntriesMaps = append(aclEntriesMaps, aclEntriesArg)
	}
	d.Set("acl_entries", aclEntriesMaps)

	listTagResourcesObject, err := albService.ListTagResources(d.Id(), "acl")
	d.Set("tags", tags.TagsToMap(listTagResourcesObject))
	return nil
}

func resourceAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := NewService(client)
	var response map[string]interface{}
	var err error
	d.Partial(true)

	if d.HasChange("tags") {
		if err := albService.SetResourceTags(d, "acl"); err != nil {
			return errs.WrapError(err)
		}
	}
	update := false
	request := map[string]interface{}{
		"ResourceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["NewResourceGroupId"] = v
	}
	request["ResourceType"] = "acl"
	if update {
		action := "MoveResourceGroup"
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"NotExist.ResourceGroup"}) || helper.NeedRetry(err) {
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
	update = false
	updateAclAttributeReq := map[string]interface{}{
		"AclId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("acl_name") {
		update = true
	}
	if v, ok := d.GetOk("acl_name"); ok {
		updateAclAttributeReq["AclName"] = v
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			updateAclAttributeReq["DryRun"] = v
		}
		action := "UpdateAclAttribute"
		request["ClientToken"] = helper.BuildClientToken("UpdateAclAttribute")
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, updateAclAttributeReq, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "SystemBusy", "Throttling", "AclStatusNotStable"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, updateAclAttributeReq)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AclStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	update = false

	if d.HasChange("acl_entries") {

		oraw, nraw := d.GetChange("acl_entries")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()

		if len(remove) > 0 {
			removeList := helper.SplitSlice(remove, 20)
			for _, item := range removeList {
				removeEntriesFromAclReq := map[string]interface{}{
					"AclId": d.Id(),
				}

				aclEntriesMaps := make([]string, 0)
				for _, aclEntries := range item {
					aclEntriesArg := aclEntries.(map[string]interface{})
					aclEntriesMaps = append(aclEntriesMaps, aclEntriesArg["entry"].(string))
				}
				removeEntriesFromAclReq["Entries"] = aclEntriesMaps

				if v, ok := d.GetOkExists("dry_run"); ok {
					removeEntriesFromAclReq["DryRun"] = v
				}
				action := "RemoveEntriesFromAcl"
				request["ClientToken"] = helper.BuildClientToken("RemoveEntriesFromAcl")
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Alb", "2020-06-16", action, nil, removeEntriesFromAclReq, true)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"IncorrectStatus.Acl", "OperationFailed.ResourceGroupStatusCheckFail", "SystemBusy", "Throttling", "AclStatusNotStable"}) || helper.NeedRetry(err) {
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
				stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AclStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}
			}
		}

		if len(create) > 0 {
			createList := helper.SplitSlice(create, 20)
			for _, item := range createList {
				addEntriesToAclReq := map[string]interface{}{
					"AclId": d.Id(),
				}
				aclEntriesMaps := make([]map[string]interface{}, 0)
				for _, aclEntries := range item {
					aclEntriesArg := aclEntries.(map[string]interface{})
					aclEntriesMap := map[string]interface{}{}
					aclEntriesMap["Description"] = aclEntriesArg["description"]
					aclEntriesMap["Entry"] = aclEntriesArg["entry"]
					aclEntriesMaps = append(aclEntriesMaps, aclEntriesMap)
				}
				addEntriesToAclReq["AclEntries"] = aclEntriesMaps

				if v, ok := d.GetOkExists("dry_run"); ok {
					addEntriesToAclReq["DryRun"] = v
				}
				action := "AddEntriesToAcl"
				request["ClientToken"] = helper.BuildClientToken("AddEntriesToAcl")
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Alb", "2020-06-16", action, nil, addEntriesToAclReq, true)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "SystemBusy", "Throttling", "AclStatusNotStable"}) || helper.NeedRetry(err) {
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
				stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AclStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}
			}
		}

	}
	d.Partial(false)
	return resourceAclRead(d, meta)
}

func resourceAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAcl"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"AclId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = helper.BuildClientToken("DeleteAcl")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "SystemBusy", "ResourceInUse.Acl", "IncorrectStatus.Acl", "AclStatusNotStable"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.Acl"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

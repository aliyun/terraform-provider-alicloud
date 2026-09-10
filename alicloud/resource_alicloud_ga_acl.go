package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGaAcl() *schema.Resource {
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
				ValidateFunc: StringInSlice([]string{"IPv4", "IPv6"}, false),
			},
			"acl_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^[a-zA-Z][A-Za-z0-9._-]{2,128}$`), "The name must be `2` to `128` characters in length, and can contain letters, digits, periods (.), hyphens (-) and underscores (_). It must start with a letter."),
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
							ValidateFunc: StringMatch(regexp.MustCompile(`^[A-Za-z0-9._/-]{1,256}$`), "The description of the IP entry. The description must be 1 to 256 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.),and underscores (_)."),
						},
					},
				},
			},
			"tags": tagsSchema(),
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
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateAcl"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateAcl")
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_acl", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AclId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaAclUpdate(d, meta)
}

func resourceAliCloudGaAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	object, err := gaService.DescribeGaAcl(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_acl gaService.DescribeGaAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
			return WrapError(err)
		}
	}

	listTagResourcesObject, err := gaService.ListTagResources(d.Id(), "acl")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudGaAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := gaService.SetResourceTags(d, "acl"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	update := false
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("UpdateAclAttribute"),
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

		wait := incrementalWait(3*time.Second, 20*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Acl"}) || NeedRetry(err) {
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

		d.SetPartial("acl_name")
	}

	if !d.IsNewResource() && d.HasChange("acl_entries") {
		oraw, nraw := d.GetChange("acl_entries")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()

		if len(remove) > 0 {
			removeEntriesFromAclReq := map[string]interface{}{
				"RegionId":    client.RegionId,
				"ClientToken": buildClientToken("RemoveEntriesFromAcl"),
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

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, removeEntriesFromAclReq, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"StateError.Acl"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, removeEntriesFromAclReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if len(create) > 0 {
			addEntriesToAclReq := map[string]interface{}{
				"RegionId":    client.RegionId,
				"ClientToken": buildClientToken("AddEntriesToAcl"),
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

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, addEntriesToAclReq, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"StateError.Acl"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, addEntriesToAclReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

			stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		d.SetPartial("acl_entries")
	}

	update = false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ClientToken":  buildClientToken("ChangeResourceGroup"),
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, changeResourceGroupReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, changeResourceGroupReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("resource_group_id")
	}

	d.Partial(false)

	return resourceAliCloudGaAclRead(d, meta)
}

func resourceAliCloudGaAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteAcl"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("DeleteAcl"),
		"AclId":       d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := incrementalWait(3*time.Second, 20*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Acl", "AclHasBindListener"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.Acl"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

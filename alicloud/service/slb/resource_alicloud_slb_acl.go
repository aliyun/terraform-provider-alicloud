package slb

import (
	"fmt"
	"strings"
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ipv4",
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
			},
			"entry_list": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'entry_list' has been deprecated from provider version 1.162.0 and it will be removed in the future version. Please use the new resource 'alicloud_slb_acl_entry_attachment'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				MaxItems: 300,
				MinItems: 0,
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAccessControlList"
	request := make(map[string]interface{})
	var err error
	request["RegionId"] = client.RegionId
	if v := d.Get("resource_group_id").(string); v != "" {
		request["ResourceGroupId"] = v
	}
	request["AclName"] = strings.TrimSpace(d.Get("name").(string))
	if v, ok := d.GetOk("ip_version"); ok {
		request["AddressIPVersion"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_slb_acl", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AclId"]))
	return resourceAclUpdate(d, meta)
}

func resourceAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)

	tagsMap, err := slbService.ListTagResources(d.Id(), "acl")
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("tags", tags.TagsToMap(tagsMap))

	object, err := slbService.DescribeSlbAcl(d.Id())
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AclNotExist"}) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("name", object["AclName"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("ip_version", object["AddressIPVersion"])

	if aclEntrys, ok := object["AclEntrys"]; ok {
		if v, ok := aclEntrys.(map[string]interface{})["AclEntry"].([]interface{}); ok {
			aclEntry := make([]map[string]interface{}, 0)
			for _, val := range v {
				item := val.(map[string]interface{})
				temp := map[string]interface{}{
					"comment": item["AclEntryComment"],
					"entry":   item["AclEntryIP"],
				}

				aclEntry = append(aclEntry, temp)
			}
			if err := d.Set("entry_list", aclEntry); err != nil {
				return errs.WrapError(err)
			}
		}
	}
	return nil
}

func resourceAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	slbService := NewService(client)
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := slbService.setInstanceTags(d, tags.TagResourceAcl); err != nil {
			return errs.WrapError(err)
		}
	}

	if !d.IsNewResource() && d.HasChange("name") {
		request := map[string]interface{}{
			"AclId": d.Id(),
		}
		if v, ok := d.GetOk("name"); ok {
			request["AclName"] = v
		}
		request["RegionId"] = client.RegionId
		action := "SetAccessControlListAttribute"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("entry_list") {
		o, n := d.GetChange("entry_list")
		oe := o.(*schema.Set)
		ne := n.(*schema.Set)
		remove := oe.Difference(ne).List()
		add := ne.Difference(oe).List()

		if len(remove) > 0 {
			removeList := helper.SplitSlice(remove, 50)
			for _, item := range removeList {
				removedRequest := map[string]interface{}{
					"AclId":    d.Id(),
					"RegionId": client.RegionId,
				}
				aclEntries, err := slbService.convertAclEntriesToString(item)
				if err != nil {
					return errs.WrapError(err)
				}
				removedRequest["AclEntrys"] = aclEntries
				action := "RemoveAccessControlListEntry"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Slb", "2014-05-15", action, nil, removedRequest, false)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, removedRequest)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}

		if len(add) > 0 {
			addList := helper.SplitSlice(add, 50)
			for _, item := range addList {
				addedRequest := map[string]interface{}{
					"AclId":    d.Id(),
					"RegionId": client.RegionId,
				}
				aclEntries, err := slbService.convertAclEntriesToString(item)
				if err != nil {
					return errs.WrapError(err)
				}
				addedRequest["AclEntrys"] = aclEntries
				action := "AddAccessControlListEntry"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("Slb", "2014-05-15", action, nil, addedRequest, false)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, addedRequest)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}
	d.Partial(false)

	return resourceAclRead(d, meta)
}

func resourceAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAccessControlList"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"AclId":    d.Id(),
		"RegionId": client.RegionId,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"AclNotExist"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

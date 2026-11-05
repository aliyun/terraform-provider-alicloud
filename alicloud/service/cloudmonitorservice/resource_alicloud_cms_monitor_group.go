package cloudmonitorservice

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudCmsMonitorGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsMonitorGroupCreate,
		Read:   resourceAlicloudCmsMonitorGroupRead,
		Update: resourceAlicloudCmsMonitorGroupUpdate,
		Delete: resourceAlicloudCmsMonitorGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"contact_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"monitor_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCmsMonitorGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	request := make(map[string]interface{})
	if v, exist := d.GetOk("resource_group_id"); exist {
		action := "CreateMonitorGroupByResourceGroupId"
		request["RegionId"] = client.RegionId
		request["ResourceGroupId"] = v.(string)
		request["ResourceGroupName"] = d.Get("resource_group_name")
		for k, v := range d.Get("contact_groups").([]interface{}) {
			request[fmt.Sprintf("ContactGroupList.%d", k+1)] = v.(string)
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cms_monitor_group", action, errs.AlibabaCloudSdkGoERROR)
		}
		d.SetId(fmt.Sprint(response["Id"]))
		return resourceAlicloudCmsMonitorGroupUpdate(d, meta)
	}

	action := "CreateMonitorGroup"
	if err != nil {
		return errs.WrapError(err)
	}
	if v, ok := d.GetOk("contact_groups"); ok && v != nil {
		request["ContactGroups"] = convert.ConvertListToCommaSeparate(v.([]interface{}))
	}
	request["GroupName"] = d.Get("monitor_group_name")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cms_monitor_group", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return errs.WrapError(errs.Error("CreateMonitorGroup failed for %s", response["Message"]))
	}

	groupId, err := convert.FormatInt(response["GroupId"])
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(groupId))

	return resourceAlicloudCmsMonitorGroupUpdate(d, meta)
}
func resourceAlicloudCmsMonitorGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsMonitorGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_monitor_group cmsService.DescribeCmsMonitorGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	var contactGroups []interface{}
	if v, ok := object["ContactGroups"].(map[string]interface{})["ContactGroup"]; ok {
		for _, contactGroup := range v.([]interface{}) {
			contactGroups = append(contactGroups, contactGroup.(map[string]interface{})["Name"])
		}
	}
	d.Set("contact_groups", contactGroups)
	d.Set("monitor_group_name", object["GroupName"])
	d.Set("tags", tags.TagsToMap(object["Tags"].(map[string]interface{})["Tag"]))
	return nil
}
func resourceAlicloudCmsMonitorGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	var err error
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"GroupId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("contact_groups") {
		update = true
		request["ContactGroups"] = convert.ConvertListToCommaSeparate(d.Get("contact_groups").([]interface{}))
	}
	if d.HasChange("monitor_group_name") {
		update = true
		request["GroupName"] = d.Get("monitor_group_name")
	}
	if update {
		action := "ModifyMonitorGroup"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	if d.HasChange("tags") {
		if err := cmsService.SetResourceTags(d, ""); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAlicloudCmsMonitorGroupRead(d, meta)
}
func resourceAlicloudCmsMonitorGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMonitorGroup"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"GroupId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	if errs.IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"400", "403", "404", "ResourceNotFound"}) {
		return nil
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return errs.WrapError(errs.Error("DeleteMonitorGroup failed for %s", response["Message"]))
	}
	return nil
}

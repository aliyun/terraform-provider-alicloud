package cloudmonitorservice

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudCmsAlarmContactGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsAlarmContactGroupCreate,
		Read:   resourceAlicloudCmsAlarmContactGroupRead,
		Update: resourceAlicloudCmsAlarmContactGroupUpdate,
		Delete: resourceAlicloudCmsAlarmContactGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alarm_contact_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contacts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"describe": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_subscribed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCmsAlarmContactGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cms.CreatePutContactGroupRequest()
	request.ContactGroupName = d.Get("alarm_contact_group_name").(string)
	if v, ok := d.GetOk("contacts"); ok {
		contactNames := convert.ExpandStringList(v.(*schema.Set).List())
		request.ContactNames = &contactNames
	}

	if v, ok := d.GetOk("describe"); ok {
		request.Describe = v.(string)
	}

	if v, ok := d.GetOkExists("enable_subscribed"); ok {
		request.EnableSubscribed = requests.NewBoolean(v.(bool))
	}
	var response *cms.PutContactGroupResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := client.LegacySdkRpc(request, "/", cms.CreatePutContactGroupResponse())
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw)
		response = &cms.PutContactGroupResponse{}
		if err := legacyCmsResponse(raw, response); err != nil {
			return retry.NonRetryableError(errs.WrapError(err))
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cms_alarm_contact_group", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if response.Code != "200" {
		return errs.WrapError(errs.Error("PutContactGroup failed for %s", response.Message))
	}
	d.SetId(fmt.Sprintf("%v", request.ContactGroupName))

	return resourceAlicloudCmsAlarmContactGroupRead(d, meta)
}
func resourceAlicloudCmsAlarmContactGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsAlarmContactGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_alarm_contact_group cmsService.DescribeCmsAlarmContactGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("alarm_contact_group_name", d.Id())
	d.Set("contacts", object.Contacts.Contact)
	d.Set("describe", object.Describe)
	d.Set("enable_subscribed", object.EnableSubscribed)
	return nil
}
func resourceAlicloudCmsAlarmContactGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := cms.CreatePutContactGroupRequest()
	request.ContactGroupName = d.Id()
	if d.HasChange("contacts") {
		update = true
		contactNames := convert.ExpandStringList(d.Get("contacts").(*schema.Set).List())
		request.ContactNames = &contactNames

	}
	if d.HasChange("describe") {
		update = true
		request.Describe = d.Get("describe").(string)
	}
	if d.HasChange("enable_subscribed") {
		update = true
		request.EnableSubscribed = requests.NewBoolean(d.Get("enable_subscribed").(bool))
	}
	if update {
		var response *cms.PutContactGroupResponse
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			raw, err := client.LegacySdkRpc(request, "/", cms.CreatePutContactGroupResponse())
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw)
			response = &cms.PutContactGroupResponse{}
			if err := legacyCmsResponse(raw, response); err != nil {
				return retry.NonRetryableError(errs.WrapError(err))
			}
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		if response.Code != "200" {
			return errs.WrapError(errs.Error("PutContactGroup failed for %s", response.Message))
		}
	}
	return resourceAlicloudCmsAlarmContactGroupRead(d, meta)
}
func resourceAlicloudCmsAlarmContactGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cms.CreateDeleteContactGroupRequest()
	request.ContactGroupName = d.Id()
	var response *cms.DeleteContactGroupResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := client.LegacySdkRpc(request, "/", cms.CreateDeleteContactGroupResponse())
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw)
		response = &cms.DeleteContactGroupResponse{}
		if err := legacyCmsResponse(raw, response); err != nil {
			return retry.NonRetryableError(errs.WrapError(err))
		}
		return nil
	})
	if err != nil {
		// DeleteContactGroup reports ContactNotExists for a missing group;
		// ContactGroupNotExists is kept as a defensive alias.
		if cmsContactAbsent(err, "ContactNotExists", "ContactGroupNotExists") {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if response.Code != "200" {
		return errs.WrapError(errs.Error("DeleteContactGroup failed for %s", response.Message))
	}
	return nil
}

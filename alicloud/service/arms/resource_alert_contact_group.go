package arms

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAlertContactGroup() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlertContactGroupCreate,
		Read:   ResourceAlertContactGroupRead,
		Update: ResourceAlertContactGroupUpdate,
		Delete: ResourceAlertContactGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alert_contact_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contact_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ResourceAlertContactGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAlertContactGroup"
	request := make(map[string]interface{})
	var err error
	request["ContactGroupName"] = d.Get("alert_contact_group_name")
	if v, ok := d.GetOk("contact_ids"); ok {
		request["ContactIds"] = convert.ConvertArrayToString(v.(*schema.Set).List(), " ")
	}
	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_alert_contact_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ContactGroupId"]))

	return ResourceAlertContactGroupRead(d, meta)
}
func ResourceAlertContactGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := Service{client}
	object, err := armsService.DescribeArmsAlertContactGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_alert_contact_group armsService.DescribeArmsAlertContactGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("alert_contact_group_name", object["ContactGroupName"])
	contactIdsItems := make([]string, 0)
	if contacts, ok := object["Contacts"]; ok && contacts != nil {
		for _, contactsItem := range contacts.([]interface{}) {
			if contactId, ok := contactsItem.(map[string]interface{})["ContactId"]; ok && contactId != nil {
				contactIdsItems = append(contactIdsItems, fmt.Sprint(contactId))
			}
		}
	}
	d.Set("contact_ids", contactIdsItems)
	return nil
}
func ResourceAlertContactGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ContactGroupId": d.Id(),
	}
	if d.HasChange("alert_contact_group_name") {
		update = true
	}
	request["ContactGroupName"] = d.Get("alert_contact_group_name")
	request["RegionId"] = client.RegionId
	if d.HasChange("contact_ids") {
		update = true
	}
	if v, ok := d.GetOk("contact_ids"); ok {
		request["ContactIds"] = convert.ConvertArrayToString(v.(*schema.Set).List(), " ")
	}
	if update {
		action := "UpdateAlertContactGroup"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
	return ResourceAlertContactGroupRead(d, meta)
}
func ResourceAlertContactGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAlertContactGroup"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ContactGroupId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
	return nil
}

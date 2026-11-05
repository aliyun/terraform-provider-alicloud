package mscopensubscription

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceContactCreate,
		Read:   resourceContactRead,
		Update: resourceContactUpdate,
		Delete: resourceContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"contact_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[\u4E00-\u9FA5a-zA-Z]{2,12}$"), "The name must be 2 to 12 characters in length."),
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z0-9+_.-]+@[a-zA-Z0-9_-]+(.[\u4e00-\u9fa5a-zA-Z0-9_-]+)+$"), "The email must has correct format."),
			},
			"mobile": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9]+$`), "The mobile only has digit."),
			},
			"position": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CEO", "Finance Director", "Maintenance Director", "Other", "Project Director", "Technical Director"}, false),
			},
		},
	}
}

func resourceContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	request["ContactName"] = d.Get("contact_name")
	request["Email"] = d.Get("email")
	request["Mobile"] = d.Get("mobile")
	request["Position"] = d.Get("position")
	request["ClientToken"] = helper.BuildClientToken("CreateContact")
	request["Locale"] = "en"
	var response map[string]interface{}
	var err error
	action := "CreateContact"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_msc_sub_contact", action, errs.AlibabaCloudSdkGoERROR)
	}

	contactId, err := convert.FormatInt(response["ContactId"])
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(contactId))

	return resourceContactRead(d, meta)
}
func resourceContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mscOpenSubscriptionService := NewService(client)
	object, err := mscOpenSubscriptionService.DescribeContact(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_msc_sub_contact mscOpenSubscriptionService.DescribeContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("contact_name", object["ContactName"])
	d.Set("email", object["Email"])
	d.Set("mobile", object["Mobile"])
	d.Set("position", convertMscSubContactPositionResponse(fmt.Sprint(object["Position"])))
	return nil
}
func resourceContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error

	update := false
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}
	if d.HasChange("contact_name") {
		update = true
	}
	request["ContactName"] = d.Get("contact_name")
	if d.HasChange("email") {
		update = true
	}
	request["Email"] = d.Get("email")
	if d.HasChange("mobile") {
		update = true
	}
	request["Mobile"] = d.Get("mobile")
	if update {
		action := "UpdateContact"
		request["ClientToken"] = helper.BuildClientToken("UpdateContact")
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, true)
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
	return resourceContactRead(d, meta)
}
func resourceContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}

	action := "DeleteContact"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertMscSubContactPositionResponse(source string) string {
	switch source {
	case "Others":
		return "Other"
	}
	return source
}

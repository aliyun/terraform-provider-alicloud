package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMscSubContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMscSubContactCreate,
		Read:   resourceAlicloudMscSubContactRead,
		Update: resourceAlicloudMscSubContactUpdate,
		Delete: resourceAlicloudMscSubContactDelete,
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

func resourceAlicloudMscSubContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	request["ContactName"] = d.Get("contact_name")
	request["Email"] = d.Get("email")
	request["Mobile"] = d.Get("mobile")
	request["Position"] = d.Get("position")
	request["ClientToken"] = buildClientToken("CreateContact")
	request["Locale"] = "en"
	var response map[string]interface{}
	var err error
	action := "CreateContact"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_msc_sub_contact", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(formatInt(response["ContactId"])))

	return resourceAlicloudMscSubContactRead(d, meta)
}
func resourceAlicloudMscSubContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mscOpenSubscriptionService := MscOpenSubscriptionService{client}
	object, err := mscOpenSubscriptionService.DescribeMscSubContact(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_msc_sub_contact mscOpenSubscriptionService.DescribeMscSubContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("contact_name", object["ContactName"])
	d.Set("email", object["Email"])
	d.Set("mobile", object["Mobile"])
	d.Set("position", convertMscSubContactPositionResponse(fmt.Sprint(object["Position"])))
	return nil
}
func resourceAlicloudMscSubContactUpdate(d *schema.ResourceData, meta interface{}) error {
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
		request["ClientToken"] = buildClientToken("UpdateContact")
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, true)
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
	return resourceAlicloudMscSubContactRead(d, meta)
}
func resourceAlicloudMscSubContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}

	action := "DeleteContact"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
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

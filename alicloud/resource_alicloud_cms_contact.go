package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCmsContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsContactCreate,
		Read:   resourceAlicloudCmsContactRead,
		Update: resourceAlicloudCmsContactUpdate,
		Delete: resourceAlicloudCmsContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"contact_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh_CN", "en_US"}, false),
			},
			"im_user_ids": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"workspace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contact_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCmsContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/contact"
	query := make(map[string]*string)
	body := make(map[string]interface{})
	body["name"] = d.Get("contact_name")
	if v, ok := d.GetOk("email"); ok {
		body["email"] = v
	}
	if v, ok := d.GetOk("phone"); ok {
		body["phone"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		body["lang"] = v
	}
	if v, ok := d.GetOk("im_user_ids"); ok {
		body["imUserIds"] = v
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceExist", "InvalidParameter", "FrequencyLimit", "SmsQuotaExceeded"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, body)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_contact", action, AlibabaCloudSdkGoERROR)
	}
	contactId, ok := response["contactId"]
	if !ok || contactId == nil || fmt.Sprint(contactId) == "" {
		return WrapError(Error("CreateContact response missing contactId: %v", response))
	}
	d.SetId(fmt.Sprint(contactId))
	return resourceAlicloudCmsContactRead(d, meta)
}

func resourceAlicloudCmsContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}
	object, err := cmsServiceV2.DescribeCmsContact(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_contact DescribeCmsContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("contact_name", object["name"])
	d.Set("email", object["email"])
	d.Set("phone", object["phone"])
	d.Set("lang", object["lang"])
	d.Set("workspace", object["workspace"])
	d.Set("contact_id", object["contactId"])
	d.Set("region_id", client.RegionId)
	if imUserIds, ok := object["imUserIds"]; ok && imUserIds != nil {
		d.Set("im_user_ids", imUserIds)
	}
	return nil
}

func resourceAlicloudCmsContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	contactId := d.Id()
	action := fmt.Sprintf("/contact/%s", contactId)
	query := make(map[string]*string)
	body := make(map[string]interface{})
	update := false
	if d.HasChange("contact_name") {
		update = true
	}
	body["name"] = d.Get("contact_name")
	if d.HasChange("email") {
		update = true
	}
	body["email"] = d.Get("email")
	if d.HasChange("phone") {
		update = true
	}
	body["phone"] = d.Get("phone")
	if d.HasChange("lang") {
		update = true
	}
	body["lang"] = d.Get("lang")
	if d.HasChange("im_user_ids") {
		update = true
	}
	body["imUserIds"] = d.Get("im_user_ids")
	if update {
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			var err error
			response, err = client.RoaPatch("Cms", "2024-03-30", action, query, nil, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidParameter", "ResourceNotFound"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, body)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCmsContactRead(d, meta)
}

func resourceAlicloudCmsContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/contacts"
	query := make(map[string]*string)
	contactIdsJson, _ := json.Marshal([]string{d.Id()})
	contactIdsStr := string(contactIdsJson)
	query["contactIds"] = &contactIdsStr
	source := "OBS"
	query["source"] = &source
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceNotFound", "404", "InvalidParameter"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound", "404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

// DescribeCmsContact queries a single CMS alarm contact by id via ListContacts.
func (s *CmsServiceV2) DescribeCmsContact(id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "/contact"
	query := make(map[string]*string)
	contactIdsJson, _ := json.Marshal([]string{id})
	contactIdsStr := string(contactIdsJson)
	query["contactIds"] = &contactIdsStr
	source := "OBS"
	query["source"] = &source
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound", "404"}) {
			return object, WrapErrorf(NotFoundErr("Contact", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	var contactsArr []interface{}
	if dataRaw, ok := response["data"].(map[string]interface{}); ok {
		if arr, ok := dataRaw["data"].([]interface{}); ok {
			contactsArr = arr
		} else if arr, ok := dataRaw["contacts"].([]interface{}); ok {
			contactsArr = arr
		}
	} else if arr, ok := response["contacts"].([]interface{}); ok {
		contactsArr = arr
	}
	if len(contactsArr) == 0 {
		return object, WrapErrorf(NotFoundErr("Contact", id), NotFoundMsg, response)
	}
	return contactsArr[0].(map[string]interface{}), nil
}

package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudHbrHanaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrHanaInstanceCreate,
		Read:   resourceAlicloudHbrHanaInstanceRead,
		Update: resourceAlicloudHbrHanaInstanceUpdate,
		Delete: resourceAlicloudHbrHanaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alert_setting": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"INHERITED"}, false),
			},
			"ecs_instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hana_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hana_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"validate_certificate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudHbrHanaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHanaInstance"
	request := make(map[string]interface{})
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("alert_setting"); ok {
		request["AlertSetting"] = v
	}
	if v, ok := d.GetOk("ecs_instance_ids"); ok {
		request["EcsInstanceId"] = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("hana_name"); ok {
		request["HanaName"] = v
	}
	if v, ok := d.GetOk("host"); ok {
		request["Host"] = v
	}
	if v, ok := d.GetOk("instance_number"); ok {
		request["InstanceNumber"] = v
	}
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("sid"); ok {
		request["Sid"] = v
	}
	if v, ok := d.GetOkExists("use_ssl"); ok {
		request["UseSsl"] = v
	}
	if v, ok := d.GetOk("user_name"); ok {
		request["UserName"] = v
	}
	if v, ok := d.GetOkExists("validate_certificate"); ok {
		request["ValidateCertificate"] = v
	}
	request["VaultId"] = d.Get("vault_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_hana_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["VaultId"], ":", response["ClusterId"]))

	return resourceAlicloudHbrHanaInstanceRead(d, meta)
}
func resourceAlicloudHbrHanaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrHanaInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_hana_instance hbrService.DescribeHbrHanaInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("hana_instance_id", parts[1])
	d.Set("vault_id", parts[0])
	d.Set("alert_setting", object["AlertSetting"])
	d.Set("hana_name", object["HanaName"])
	d.Set("host", object["Host"])
	if v, ok := object["InstanceNumber"]; ok && fmt.Sprint(v) != "0" {
		d.Set("instance_number", formatInt(v))
	}
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["Status"])
	d.Set("use_ssl", object["UseSsl"])
	d.Set("user_name", object["UserName"])
	d.Set("validate_certificate", object["ValidateCertificate"])
	return nil
}
func resourceAlicloudHbrHanaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ClusterId": parts[1],
		"VaultId":   parts[0],
	}
	if d.HasChange("instance_number") {
		update = true
	}
	if v, ok := d.GetOk("instance_number"); ok {
		request["InstanceNumber"] = v
	}
	if d.HasChange("use_ssl") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("use_ssl"); ok {
		request["UseSsl"] = v
	}
	if d.HasChange("validate_certificate") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("validate_certificate"); ok {
		request["ValidateCertificate"] = v
	}
	if d.HasChange("alert_setting") {
		update = true
		if v, ok := d.GetOk("alert_setting"); ok {
			request["AlertSetting"] = v
		}
	}
	if d.HasChange("hana_name") {
		update = true
		if v, ok := d.GetOk("hana_name"); ok {
			request["HanaName"] = v
		}
	}
	if d.HasChange("host") {
		update = true
		if v, ok := d.GetOk("host"); ok {
			request["Host"] = v
		}
	}
	if d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
	}
	if d.HasChange("user_name") {
		update = true
		if v, ok := d.GetOk("user_name"); ok {
			request["UserName"] = v
		}
	}
	if update {
		action := "UpdateHanaInstance"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	return resourceAlicloudHbrHanaInstanceRead(d, meta)
}
func resourceAlicloudHbrHanaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteHanaInstance"
	var response map[string]interface{}
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ClusterId": parts[1],
		"VaultId":   parts[0],
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("sid"); ok {
		request["Sid"] = v
	} else {
		return WrapError(Error(`[ERROR] To destroy the resource, you must configure the 'sid' attribute.`))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}

package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEdasNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasNamespaceCreate,
		Read:   resourceAlicloudEdasNamespaceRead,
		Update: resourceAlicloudEdasNamespaceUpdate,
		Delete: resourceAlicloudEdasNamespaceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"debug_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"namespace_logical_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 63),
			},
		},
	}
}

func resourceAlicloudEdasNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/pop/v5/user_region_def"
	request := make(map[string]*string)
	conn, err := client.NewEdasClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("debug_enable"); ok {
		request["DebugEnable"] = StringPointer(strconv.FormatBool(v.(bool)))
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = StringPointer(v.(string))
	}
	request["RegionTag"] = StringPointer(d.Get("namespace_logical_id").(string))
	request["RegionName"] = StringPointer(d.Get("namespace_name").(string))
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2017-08-01"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_namespace", action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	responseUserDefineRegionEntity := response["UserDefineRegionEntity"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseUserDefineRegionEntity["Id"]))

	return resourceAlicloudEdasNamespaceRead(d, meta)
}
func resourceAlicloudEdasNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}
	object, err := edasService.DescribeEdasNamespace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_edas_namespace edasService.DescribeEdasNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("debug_enable", object["DebugEnable"])
	d.Set("description", object["Description"])
	d.Set("namespace_logical_id", object["RegionId"])
	d.Set("namespace_name", object["RegionName"])
	return nil
}
func resourceAlicloudEdasNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewEdasClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]*string{
		"Id": StringPointer(d.Id()),
	}

	request["RegionTag"] = StringPointer(d.Get("namespace_logical_id").(string))
	if d.HasChange("namespace_name") {
		update = true
	}
	request["RegionName"] = StringPointer(d.Get("namespace_name").(string))
	if v, ok := d.GetOkExists("debug_enable"); ok {
		request["DebugEnable"] = StringPointer(strconv.FormatBool(v.(bool)))
	}
	if d.HasChange("debug_enable") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = StringPointer(v.(string))
	}
	if d.HasChange("description") {
		update = true
	}
	if update {
		action := "/pop/v5/user_region_def"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2017-08-01"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
		if fmt.Sprint(response["Code"]) != "200" {
			return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
		}
	}
	return resourceAlicloudEdasNamespaceRead(d, meta)
}
func resourceAlicloudEdasNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v5/user_region_def"
	var response map[string]interface{}
	conn, err := client.NewEdasClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]*string{
		"Id": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2017-08-01"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "DELETE "+action, response))
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "DELETE "+action, response))
	}
	return nil
}

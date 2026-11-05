package edas

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	var err error
	if v, ok := d.GetOkExists("debug_enable"); ok {
		request["DebugEnable"] = helper.StringPointer(strconv.FormatBool(v.(bool)))
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = helper.StringPointer(v.(string))
	}
	request["RegionTag"] = helper.StringPointer(d.Get("namespace_logical_id").(string))
	request["RegionName"] = helper.StringPointer(d.Get("namespace_name").(string))
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("Edas", "2017-08-01", action, request, nil, nil, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_edas_namespace", action, errs.AlibabaCloudSdkGoERROR)
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
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_edas_namespace edasService.DescribeEdasNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("debug_enable", object["DebugEnable"])
	d.Set("description", object["Description"])
	d.Set("namespace_logical_id", object["RegionId"])
	d.Set("namespace_name", object["RegionName"])
	return nil
}
func resourceAlicloudEdasNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]*string{
		"Id": helper.StringPointer(d.Id()),
	}

	request["RegionTag"] = helper.StringPointer(d.Get("namespace_logical_id").(string))
	if d.HasChange("namespace_name") {
		update = true
	}
	request["RegionName"] = helper.StringPointer(d.Get("namespace_name").(string))
	if v, ok := d.GetOkExists("debug_enable"); ok {
		request["DebugEnable"] = helper.StringPointer(strconv.FormatBool(v.(bool)))
	}
	if d.HasChange("debug_enable") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = helper.StringPointer(v.(string))
	}
	if d.HasChange("description") {
		update = true
	}
	if update {
		action := "/pop/v5/user_region_def"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("Edas", "2017-08-01", action, request, nil, nil, false)
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
	return resourceAlicloudEdasNamespaceRead(d, meta)
}
func resourceAlicloudEdasNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v5/user_region_def"
	var response map[string]interface{}
	var err error
	request := map[string]*string{
		"Id": helper.StringPointer(d.Id()),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("Edas", "2017-08-01", action, request, nil, nil, false)
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

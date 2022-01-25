package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudSddpDataLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSddpDataLimitCreate,
		Read:   resourceAlicloudSddpDataLimitRead,
		Update: resourceAlicloudSddpDataLimitUpdate,
		Delete: resourceAlicloudSddpDataLimitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"audit_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"engine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "SQLServer"}, false),
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_store_day": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{180, 30, 365, 90}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("audit_status"); ok && fmt.Sprint(v) == "1" {
						return false
					}
					return true
				},
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("resource_type"); ok && fmt.Sprint(v) == "RDS" {
						return false
					}
					return true
				},
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("resource_type"); ok && fmt.Sprint(v) == "RDS" {
						return false
					}
					return true
				},
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MaxCompute", "OSS", "RDS"}, false),
			},
			"service_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("resource_type"); ok && fmt.Sprint(v) == "RDS" {
						return false
					}
					return true
				},
			},
		},
	}
}

func resourceAlicloudSddpDataLimitCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDataLimit"
	request := make(map[string]interface{})
	conn, err := client.NewSddpClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("audit_status"); ok {
		request["AuditStatus"] = v
	}
	if v, ok := d.GetOk("engine_type"); ok {
		request["EngineType"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("log_store_day"); ok {
		request["LogStoreDay"] = v
	}

	request["ParentId"] = d.Get("parent_id")
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOk("port"); ok {
		request["Port"] = v
	}
	request["ResourceType"] = convertSddpDataLimitResourceTypeRequest(d.Get("resource_type").(string))
	if v, ok := d.GetOk("service_region_id"); ok {
		request["ServiceRegionId"] = v
	}
	if v, ok := d.GetOk("user_name"); ok {
		request["UserName"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sddp_data_limit", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAlicloudSddpDataLimitRead(d, meta)
}
func resourceAlicloudSddpDataLimitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sddpService := SddpService{client}
	object, err := sddpService.DescribeSddpDataLimit(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sddp_data_limit sddpService.DescribeSddpDataLimit Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["AuditStatus"]; ok {
		d.Set("audit_status", formatInt(v))
	}
	d.Set("engine_type", object["EngineType"])
	if v, ok := object["LogStoreDay"]; ok {
		d.Set("log_store_day", formatInt(v))
	}
	d.Set("parent_id", object["ParentId"])
	if v, ok := object["Port"]; ok {
		d.Set("port", formatInt(v))
	}
	d.Set("service_region_id", object["RegionId"])
	d.Set("resource_type", object["ResourceTypeCode"])
	d.Set("user_name", object["UserName"])
	return nil
}
func resourceAlicloudSddpDataLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Id": d.Id(),
	}
	request["ResourceType"] = convertSddpDataLimitResourceTypeRequest(d.Get("resource_type").(string))
	if d.HasChange("audit_status") {
		update = true
		if v, ok := d.GetOkExists("audit_status"); ok {
			request["AuditStatus"] = v
		}
	}
	if d.HasChange("log_store_day") {
		update = true
		if v, ok := d.GetOk("log_store_day"); ok {
			request["LogStoreDay"] = v
		}
	}

	if update {
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}
		action := "ModifyDataLimit"
		conn, err := client.NewSddpClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}
	return resourceAlicloudSddpDataLimitRead(d, meta)
}
func resourceAlicloudSddpDataLimitDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDataLimit"
	var response map[string]interface{}
	conn, err := client.NewSddpClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Id": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
func convertSddpDataLimitResourceTypeRequest(source interface{}) interface{} {
	switch source {
	case "MaxCompute":
		return 1
	case "OSS":
		return 2
	case "RDS":
		return 5
	}
	return 0
}

package cloudfirewall

import (
	json "encoding/json"
	"fmt"
	"log"
	"time"

	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	convert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	retry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sjson "github.com/tidwall/sjson"
)

func ResourceCloudFirewallThreatIntelligenceSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceThreatIntelligenceSwitchCreate,
		Read:   resourceThreatIntelligenceSwitchRead,
		Update: resourceThreatIntelligenceSwitchUpdate,
		Delete: resourceThreatIntelligenceSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_status": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceThreatIntelligenceSwitchCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ModifyThreatIntelligenceSwitch"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	dataList := make(map[string]interface{})

	if v, ok := d.GetOkExists("enable_status"); ok {
		dataList["EnableStatus"] = v
	}

	if v, ok := d.GetOkExists("action"); ok {
		dataList["Action"] = v
	}

	if v, ok := d.GetOkExists("category_id"); ok {
		dataList["CategoryId"] = v
	}

	CategoryListMap := make([]interface{}, 0)
	CategoryListMap = append(CategoryListMap, dataList)
	request["CategoryList"] = CategoryListMap

	jsonString := convert.ConvertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "CategoryList.0.CategoryId", d.Get("category_id"))
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_firewall_threat_intelligence_switch", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(d.Get("category_id")))

	return resourceThreatIntelligenceSwitchRead(d, meta)
}

func resourceThreatIntelligenceSwitchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := NewService(client)

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallThreatIntelligenceSwitch(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_threat_intelligence_switch DescribeCloudFirewallThreatIntelligenceSwitch Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("action", objectRaw["Action"])
	d.Set("enable_status", objectRaw["EnableStatus"])
	d.Set("category_id", objectRaw["CategoryId"])

	d.Set("category_id", d.Id())

	return nil
}

func resourceThreatIntelligenceSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyThreatIntelligenceSwitch"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	dataList := make(map[string]interface{})

	if d.HasChange("enable_status") {
		update = true
	}
	if v, ok := d.GetOk("enable_status"); ok {
		dataList["EnableStatus"] = v
	}

	if d.HasChange("action") {
		update = true
	}
	if v, ok := d.GetOk("action"); ok {
		dataList["Action"] = v
	}

	if d.HasChange("category_id") {
		update = true
	}
	if v, ok := d.GetOk("category_id"); ok {
		dataList["CategoryId"] = v
	}

	CategoryListMap := make([]interface{}, 0)
	CategoryListMap = append(CategoryListMap, dataList)
	request["CategoryList"] = CategoryListMap

	jsonString := convert.ConvertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "CategoryList.0.CategoryId", d.Id())
	_ = json.Unmarshal([]byte(jsonString), &request)

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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

	return resourceThreatIntelligenceSwitchRead(d, meta)
}

func resourceThreatIntelligenceSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Threat Intelligence Switch. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

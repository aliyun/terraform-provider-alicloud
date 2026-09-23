package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudCloudFirewallThreatIntelligenceSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallThreatIntelligenceSwitchCreate,
		Read:   resourceAliCloudCloudFirewallThreatIntelligenceSwitchRead,
		Update: resourceAliCloudCloudFirewallThreatIntelligenceSwitchUpdate,
		Delete: resourceAliCloudCloudFirewallThreatIntelligenceSwitchDelete,
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

func resourceAliCloudCloudFirewallThreatIntelligenceSwitchCreate(d *schema.ResourceData, meta interface{}) error {

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

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "CategoryList.0.CategoryId", d.Get("category_id"))
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_threat_intelligence_switch", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(d.Get("category_id")))

	return resourceAliCloudCloudFirewallThreatIntelligenceSwitchRead(d, meta)
}

func resourceAliCloudCloudFirewallThreatIntelligenceSwitchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallThreatIntelligenceSwitch(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_threat_intelligence_switch DescribeCloudFirewallThreatIntelligenceSwitch Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("action", objectRaw["Action"])
	d.Set("enable_status", objectRaw["EnableStatus"])
	d.Set("category_id", objectRaw["CategoryId"])

	d.Set("category_id", d.Id())

	return nil
}

func resourceAliCloudCloudFirewallThreatIntelligenceSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
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

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "CategoryList.0.CategoryId", d.Id())
	_ = json.Unmarshal([]byte(jsonString), &request)

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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

	return resourceAliCloudCloudFirewallThreatIntelligenceSwitchRead(d, meta)
}

func resourceAliCloudCloudFirewallThreatIntelligenceSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Threat Intelligence Switch. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

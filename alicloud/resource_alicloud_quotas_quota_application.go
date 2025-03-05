// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudQuotasQuotaApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudQuotasQuotaApplicationCreate,
		Read:   resourceAliCloudQuotasQuotaApplicationRead,
		Delete: resourceAliCloudQuotasQuotaApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"approve_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"audit_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Async", "Sync"}, false),
			},
			"audit_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"desire_value": {
				Type:     schema.TypeFloat,
				Required: true,
				ForceNew: true,
			},
			"dimensions": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"effective_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"env_language": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"zh", "en"}, false),
			},
			"expire_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"notice_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: IntInSlice([]int{0, 3}),
			},
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_action_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"CommonQuota", "FlowControl", "WhiteListLabel"}, false),
			},
			"quota_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudQuotasQuotaApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateQuotaApplication"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["ProductCode"] = d.Get("product_code")
	request["QuotaActionCode"] = d.Get("quota_action_code")
	request["DesireValue"] = d.Get("desire_value")
	request["Reason"] = d.Get("reason")
	if v, ok := d.GetOk("notice_type"); ok {
		request["NoticeType"] = v
	}
	if v, ok := d.GetOk("effective_time"); ok {
		request["EffectiveTime"] = v
	}
	if v, ok := d.GetOk("expire_time"); ok {
		request["ExpireTime"] = v
	}
	if v, ok := d.GetOk("dimensions"); ok {
		dimensionsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Key"] = dataLoopTmp["key"]
			dataLoopMap["Value"] = dataLoopTmp["value"]
			dimensionsMaps = append(dimensionsMaps, dataLoopMap)
		}
		request["Dimensions"] = dimensionsMaps
	}

	if v, ok := d.GetOk("quota_category"); ok {
		request["QuotaCategory"] = v
	}
	if v, ok := d.GetOk("audit_mode"); ok {
		request["AuditMode"] = v
	}
	if v, ok := d.GetOk("env_language"); ok {
		request["EnvLanguage"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("quotas", rpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_quotas_quota_application", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ApplicationId"]))

	return resourceAliCloudQuotasQuotaApplicationRead(d, meta)
}

func resourceAliCloudQuotasQuotaApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	quotasServiceV2 := QuotasServiceV2{client}

	objectRaw, err := quotasServiceV2.DescribeQuotasQuotaApplication(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_quotas_quota_application DescribeQuotasQuotaApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("approve_value", objectRaw["ApproveValue"])
	d.Set("audit_reason", objectRaw["AuditReason"])
	d.Set("create_time", objectRaw["ApplyTime"])
	d.Set("desire_value", objectRaw["DesireValue"])
	d.Set("effective_time", objectRaw["EffectiveTime"])
	d.Set("expire_time", objectRaw["ExpireTime"])
	d.Set("notice_type", objectRaw["NoticeType"])
	d.Set("product_code", objectRaw["ProductCode"])
	d.Set("quota_action_code", objectRaw["QuotaActionCode"])
	d.Set("quota_description", objectRaw["QuotaDescription"])
	d.Set("quota_name", objectRaw["QuotaName"])
	d.Set("quota_unit", objectRaw["QuotaUnit"])
	d.Set("reason", objectRaw["Reason"])
	d.Set("status", objectRaw["Status"])

	e := jsonata.MustCompile("$each($.Dimension, function($v, $k) {{\"value\":$v, \"key\": $k}})[]")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("dimensions", evaluation)

	return nil
}

func resourceAliCloudQuotasQuotaApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Quota Application. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

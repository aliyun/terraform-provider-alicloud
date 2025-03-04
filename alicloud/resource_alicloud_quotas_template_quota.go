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

func resourceAliCloudQuotasTemplateQuota() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudQuotasTemplateQuotaCreate,
		Read:   resourceAliCloudQuotasTemplateQuotaRead,
		Update: resourceAliCloudQuotasTemplateQuotaUpdate,
		Delete: resourceAliCloudQuotasTemplateQuotaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"desire_value": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"dimensions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"effective_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"env_language": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notice_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
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
				ValidateFunc: StringInSlice([]string{"CommonQuota", "WhiteListLabel", "FlowControl"}, false),
			},
		},
	}
}

func resourceAliCloudQuotasTemplateQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateTemplateQuotaItem"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["ProductCode"] = d.Get("product_code")
	request["QuotaActionCode"] = d.Get("quota_action_code")
	request["DesireValue"] = d.Get("desire_value")
	if v, ok := d.GetOk("notice_type"); ok {
		request["NoticeType"] = v
	}
	if v, ok := d.GetOk("env_language"); ok {
		request["EnvLanguage"] = v
	}
	if v, ok := d.GetOk("quota_category"); ok {
		request["QuotaCategory"] = v
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

	if v, ok := d.GetOk("effective_time"); ok {
		request["EffectiveTime"] = v
	}
	if v, ok := d.GetOk("expire_time"); ok {
		request["ExpireTime"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("quotas", rpc("POST", "2020-05-10", action), nil, request, nil, nil, false)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_quotas_template_quota", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudQuotasTemplateQuotaUpdate(d, meta)
}

func resourceAliCloudQuotasTemplateQuotaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	quotasServiceV2 := QuotasServiceV2{client}

	objectRaw, err := quotasServiceV2.DescribeQuotasTemplateQuota(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_quotas_template_quota DescribeQuotasTemplateQuota Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("desire_value", objectRaw["DesireValue"])
	d.Set("effective_time", objectRaw["EffectiveTime"])
	d.Set("env_language", objectRaw["EnvLanguage"])
	d.Set("expire_time", objectRaw["ExpireTime"])
	d.Set("notice_type", objectRaw["NoticeType"])
	d.Set("product_code", objectRaw["ProductCode"])
	d.Set("quota_action_code", objectRaw["QuotaActionCode"])
	d.Set("quota_category", objectRaw["QuotaCategory"])

	e := jsonata.MustCompile("$each($.Dimensions, function($v, $k) {{\"value\":$v, \"key\": $k}})[]")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("dimensions", evaluation)

	return nil
}

func resourceAliCloudQuotasTemplateQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "ModifyTemplateQuotaItem"
	var err error
	request = make(map[string]interface{})

	request["Id"] = d.Id()
	if !d.IsNewResource() && d.HasChange("desire_value") {
		update = true
	}
	request["DesireValue"] = d.Get("desire_value")
	if !d.IsNewResource() && d.HasChange("notice_type") {
		update = true
		request["NoticeType"] = d.Get("notice_type")
	}

	if !d.IsNewResource() && d.HasChange("env_language") {
		update = true
		request["EnvLanguage"] = d.Get("env_language")
	}

	if !d.IsNewResource() && d.HasChange("quota_category") {
		update = true
		request["QuotaCategory"] = d.Get("quota_category")
	}

	if !d.IsNewResource() && d.HasChange("dimensions") {
		update = true
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
	}

	if !d.IsNewResource() && d.HasChange("expire_time") {
		update = true
		request["ExpireTime"] = d.Get("expire_time")
	}

	if !d.IsNewResource() && d.HasChange("effective_time") {
		update = true
		request["EffectiveTime"] = d.Get("effective_time")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("quotas", rpc("POST", "2020-05-10", action), nil, request, nil, nil, false)

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
	return resourceAliCloudQuotasTemplateQuotaRead(d, meta)
}

func resourceAliCloudQuotasTemplateQuotaDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteTemplateQuotaItem"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["Id"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("quotas", rpc("POST", "2020-05-10", action), nil, request, nil, nil, false)

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

	return nil
}

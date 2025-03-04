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

func resourceAliCloudQuotasTemplateApplications() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudQuotasTemplateApplicationsCreate,
		Read:   resourceAliCloudQuotasTemplateApplicationsRead,
		Update: resourceAliCloudQuotasTemplateApplicationsUpdate,
		Delete: resourceAliCloudQuotasTemplateApplicationsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aliyun_uids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"notice_type": {
				Type:     schema.TypeInt,
				Optional: true,
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
			"quota_application_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"audit_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"period": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"period_value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"period_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dimensions": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"approve_value": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"env_language": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aliyun_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notice_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"quota_category": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"CommonQuota", "FlowControl", "WhiteListLabel"}, false),
			},
			"reason": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudQuotasTemplateApplicationsCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateQuotaApplicationsForTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["QuotaCategory"] = d.Get("quota_category")
	request["ProductCode"] = d.Get("product_code")
	request["QuotaActionCode"] = d.Get("quota_action_code")
	if v, ok := d.GetOk("effective_time"); ok {
		request["EffectiveTime"] = v
	}
	if v, ok := d.GetOk("expire_time"); ok {
		request["ExpireTime"] = v
	}
	if v, ok := d.GetOk("aliyun_uids"); ok {
		aliyunUidsMaps := v.([]interface{})
		request["AliyunUids"] = aliyunUidsMaps
	}

	request["DesireValue"] = d.Get("desire_value")
	if v, ok := d.GetOk("notice_type"); ok {
		request["NoticeType"] = v
	}
	if v, ok := d.GetOk("env_language"); ok {
		request["EnvLanguage"] = v
	}
	request["Reason"] = d.Get("reason")
	if v, ok := d.GetOk("dimensions"); ok {
		dimensionsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop1 := range v.(*schema.Set).List() {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["Key"] = dataLoop1Tmp["key"]
			dataLoop1Map["Value"] = dataLoop1Tmp["value"]
			dimensionsMaps = append(dimensionsMaps, dataLoop1Map)
		}
		request["Dimensions"] = dimensionsMaps
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_quotas_template_applications", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BatchQuotaApplicationId"]))

	quotasServiceV2 := QuotasServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(response["BatchQuotaApplicationId"])}, d.Timeout(schema.TimeoutCreate), 5*time.Second, quotasServiceV2.QuotasTemplateApplicationsStateRefreshFunc(d.Id(), "BatchQuotaApplicationId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudQuotasTemplateApplicationsRead(d, meta)
}

func resourceAliCloudQuotasTemplateApplicationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	quotasServiceV2 := QuotasServiceV2{client}

	objectRaw, err := quotasServiceV2.DescribeQuotasTemplateApplications(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_quotas_template_applications DescribeQuotasTemplateApplications Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("desire_value", objectRaw["DesireValue"])
	d.Set("effective_time", objectRaw["EffectiveTime"])
	d.Set("expire_time", objectRaw["ExpireTime"])
	d.Set("product_code", objectRaw["ProductCode"])
	d.Set("quota_action_code", objectRaw["QuotaActionCode"])
	d.Set("quota_category", objectRaw["QuotaCategory"])
	d.Set("reason", objectRaw["Reason"])

	aliyunUids1Raw := make([]interface{}, 0)
	if objectRaw["AliyunUids"] != nil {
		aliyunUids1Raw = objectRaw["AliyunUids"].([]interface{})
	}

	d.Set("aliyun_uids", aliyunUids1Raw)

	e := jsonata.MustCompile("$each($.Dimensions, function($v, $k) {{\"value\":$v, \"key\": $k}})[]")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("dimensions", evaluation)
	objectRaw, err = quotasServiceV2.DescribeListQuotaApplicationsDetailForTemplate(d.Id())
	if err != nil {
		return WrapError(err)
	}

	quotaApplications1Raw := objectRaw["QuotaApplications"]
	quotaApplicationDetailsMaps := make([]map[string]interface{}, 0)
	if quotaApplications1Raw != nil {
		for _, quotaApplicationsChild1Raw := range quotaApplications1Raw.([]interface{}) {
			quotaApplicationDetailsMap := make(map[string]interface{})
			quotaApplicationsChild1Raw := quotaApplicationsChild1Raw.(map[string]interface{})
			quotaApplicationDetailsMap["aliyun_uid"] = quotaApplicationsChild1Raw["AliyunUid"]
			quotaApplicationDetailsMap["application_id"] = quotaApplicationsChild1Raw["ApplicationId"]
			quotaApplicationDetailsMap["approve_value"] = quotaApplicationsChild1Raw["ApproveValue"]
			quotaApplicationDetailsMap["audit_reason"] = quotaApplicationsChild1Raw["AuditReason"]
			quotaApplicationDetailsMap["dimensions"] = quotaApplicationsChild1Raw["QuotaDimension"]
			quotaApplicationDetailsMap["env_language"] = quotaApplicationsChild1Raw["EnvLanguage"]
			quotaApplicationDetailsMap["notice_type"] = quotaApplicationsChild1Raw["NoticeType"]
			quotaApplicationDetailsMap["quota_arn"] = quotaApplicationsChild1Raw["QuotaArn"]
			quotaApplicationDetailsMap["quota_description"] = quotaApplicationsChild1Raw["QuotaDescription"]
			quotaApplicationDetailsMap["quota_name"] = quotaApplicationsChild1Raw["QuotaName"]
			quotaApplicationDetailsMap["quota_unit"] = quotaApplicationsChild1Raw["QuotaUnit"]
			quotaApplicationDetailsMap["reason"] = quotaApplicationsChild1Raw["Reason"]
			quotaApplicationDetailsMap["status"] = quotaApplicationsChild1Raw["Status"]

			periodMaps := make([]map[string]interface{}, 0)
			periodMap := make(map[string]interface{})
			period1Raw := quotaApplicationsChild1Raw["Period"].(map[string]interface{})
			periodMap["period_unit"] = period1Raw["PeriodUnit"]
			periodMap["period_value"] = formatInt(period1Raw["PeriodValue"])

			periodMaps = append(periodMaps, periodMap)
			quotaApplicationDetailsMap["period"] = periodMaps
			quotaApplicationDetailsMaps = append(quotaApplicationDetailsMaps, quotaApplicationDetailsMap)
		}
	}
	d.Set("quota_application_details", quotaApplicationDetailsMaps)

	return nil
}

func resourceAliCloudQuotasTemplateApplicationsUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Template Applications.")
	return nil
}

func resourceAliCloudQuotasTemplateApplicationsDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Template Applications. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

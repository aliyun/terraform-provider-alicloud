package alicloud

import (
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceAlicloudQuotasTemplateQuota() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudQuotasTemplateQuotaCreate,
		Read:   resourceAlicloudQuotasTemplateQuotaRead,
		Update: resourceAlicloudQuotasTemplateQuotaUpdate,
		Delete: resourceAlicloudQuotasTemplateQuotaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"quota_action_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"desire_value": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"notice_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 3}),
			},
			"dimensions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"env_language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en"}, false),
			},
		},
	}
}

func resourceAlicloudQuotasTemplateQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTemplateQuotaItem"
	request := make(map[string]interface{})
	conn, err := client.NewQuotasClient()
	if err != nil {
		return WrapError(err)
	}

	request["ProductCode"] = d.Get("product_code")
	request["QuotaActionCode"] = d.Get("quota_action_code")
	request["DesireValue"] = d.Get("desire_value")

	if v, ok := d.GetOkExists("notice_type"); ok {
		request["NoticeType"] = v
	}

	if v, ok := d.GetOk("dimensions"); ok {
		dimensionsMaps := make([]map[string]interface{}, 0)
		for _, dimensions := range v.(*schema.Set).List() {
			dimensionsArg := dimensions.(map[string]interface{})
			dimensionsMap := map[string]interface{}{}

			if dimensionsKey, ok := dimensionsArg["key"]; ok && dimensionsKey != "" {
				dimensionsMap["Key"] = dimensionsKey
			}

			if dimensionsValue, ok := dimensionsArg["value"]; ok && dimensionsValue != "" {
				dimensionsMap["Value"] = dimensionsValue
			}

			dimensionsMaps = append(dimensionsMaps, dimensionsMap)
		}
		request["Dimensions"] = dimensionsMaps
	}

	if v, ok := d.GetOk("env_language"); ok {
		request["EnvLanguage"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_quotas_template_quota", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAlicloudQuotasTemplateQuotaRead(d, meta)
}

func resourceAlicloudQuotasTemplateQuotaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	quotasService := QuotasService{client}

	object, err := quotasService.DescribeQuotasTemplateQuota(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("product_code", object["ProductCode"])
	d.Set("quota_action_code", object["QuotaActionCode"])
	d.Set("desire_value", object["DesireValue"])
	d.Set("notice_type", object["NoticeType"])
	//d.Set("dimensions", object["Dimensions"])
	d.Set("env_language", object["EnvLanguage"])

	return nil
}

func resourceAlicloudQuotasTemplateQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"Id": d.Id(),
	}

	if d.HasChange("desire_value") {
		update = true
	}
	if v, ok := d.GetOkExists("desire_value"); ok {
		request["DesireValue"] = v
	}

	if d.HasChange("notice_type") {
		update = true
	}
	if v, ok := d.GetOkExists("notice_type"); ok {
		request["NoticeType"] = v
	}

	if d.HasChange("env_language") {
		update = true
	}
	if v, ok := d.GetOk("env_language"); ok {
		request["EnvLanguage"] = v
	}

	if update {
		action := "ModifyTemplateQuotaItem"
		conn, err := client.NewQuotasClient()
		if err != nil {
			return WrapError(err)
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	return resourceAlicloudQuotasTemplateQuotaRead(d, meta)
}

func resourceAlicloudQuotasTemplateQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTemplateQuotaItem"
	var response map[string]interface{}

	conn, err := client.NewQuotasClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"Id": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

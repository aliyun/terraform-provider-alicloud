package alb

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceHealthCheckTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHealthCheckTemplatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"health_check_template_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"health_check_template_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\\_\\.\\-]{1,127}$`), "The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter."),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_codes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"health_check_connect_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_http_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceHealthCheckTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListHealthCheckTemplates"
	request := make(map[string]interface{})
	if m, ok := d.GetOk("health_check_template_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("HealthCheckTemplateIds.%d", k+1)] = v.(string)
		}
	}
	if v, ok := d.GetOk("health_check_template_name"); ok {
		request["HealthCheckTemplateNames.1"] = v
	}
	request["MaxResults"] = helper.PageSizeLarge
	var objects []map[string]interface{}
	var healthCheckTemplateNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		healthCheckTemplateNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	var err error
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_alb_health_check_templates", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.HealthCheckTemplates", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.HealthCheckTemplates", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if healthCheckTemplateNameRegex != nil && !healthCheckTemplateNameRegex.MatchString(fmt.Sprint(item["HealthCheckTemplateName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["HealthCheckTemplateId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		health_check_connect_portVal, err := convert.FormatInt(object["HealthCheckConnectPort"])
		if err != nil {
			return err
		}

		health_check_intervalVal, err := convert.FormatInt(object["HealthCheckInterval"])
		if err != nil {
			return err
		}

		health_check_timeoutVal, err := convert.FormatInt(object["HealthCheckTimeout"])
		if err != nil {
			return err
		}
		healthy_thresholdVal, err := convert.FormatInt(object["HealthyThreshold"])
		if err != nil {
			return err
		}
		unhealthy_thresholdVal, err := convert.FormatInt(object["UnhealthyThreshold"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"health_check_codes":         object["HealthCheckCodes"],
			"health_check_connect_port":  health_check_connect_portVal,
			"health_check_host":          object["HealthCheckHost"],
			"health_check_http_version":  object["HealthCheckHttpVersion"],
			"health_check_interval":      health_check_intervalVal,
			"health_check_method":        object["HealthCheckMethod"],
			"health_check_path":          object["HealthCheckPath"],
			"health_check_protocol":      object["HealthCheckProtocol"],
			"id":                         fmt.Sprint(object["HealthCheckTemplateId"]),
			"health_check_template_id":   fmt.Sprint(object["HealthCheckTemplateId"]),
			"health_check_template_name": object["HealthCheckTemplateName"],
			"health_check_timeout":       health_check_timeoutVal,
			"healthy_threshold":          healthy_thresholdVal,
			"unhealthy_threshold":        unhealthy_thresholdVal,
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["HealthCheckTemplateName"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("templates", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

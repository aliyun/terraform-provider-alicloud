package sae

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudSaeGreyTagRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSaeGreyTagRouteCreate,
		Read:   resourceAliCloudSaeGreyTagRouteRead,
		Update: resourceAliCloudSaeGreyTagRouteUpdate,
		Delete: resourceAliCloudSaeGreyTagRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"grey_tag_route_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sc_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"condition": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringInSlice([]string{"AND", "OR"}, false),
						},
						"items": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"param", "cookie", "header"}, false),
									},
									"cond": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{">", "<", ">=", "<=", "==", "!="}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"operator": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"rawvalue", "list", "mod", "deterministic_proportional_steaming_division"}, false),
									},
								},
							},
						},
					},
				},
			},
			"dubbo_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"group": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"condition": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringInSlice([]string{"AND", "OR"}, false),
						},
						"items": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"expr": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cond": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{">", "<", ">=", "<=", "==", "!="}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"operator": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"rawvalue", "list", "mod", "deterministic_proportional_steaming_division"}, false),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudSaeGreyTagRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/pop/v1/sam/tagroute/greyTagRoute"
	request := make(map[string]*string)
	var err error
	request["AppId"] = helper.StringPointer(d.Get("app_id").(string))
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = helper.StringPointer(v.(string))
	}
	request["Name"] = helper.StringPointer(d.Get("grey_tag_route_name").(string))
	if v, ok := d.GetOk("sc_rules"); ok {
		scRulesMaps := make([]map[string]interface{}, 0)
		for _, scRules := range v.(*schema.Set).List() {
			scRulesArg := scRules.(map[string]interface{})
			scRulesMap := map[string]interface{}{}
			scRulesMap["path"] = scRulesArg["path"]
			scRulesMap["condition"] = scRulesArg["condition"]
			itemsMaps := make([]map[string]interface{}, 0)
			for _, items := range scRulesArg["items"].(*schema.Set).List() {
				itemsArg := items.(map[string]interface{})
				itemsMap := map[string]interface{}{}
				itemsMap["name"] = itemsArg["name"]
				itemsMap["cond"] = itemsArg["cond"]
				itemsMap["type"] = itemsArg["type"]
				itemsMap["value"] = itemsArg["value"]
				itemsMap["operator"] = itemsArg["operator"]
				itemsMaps = append(itemsMaps, itemsMap)
			}
			scRulesMap["items"] = itemsMaps
			scRulesMaps = append(scRulesMaps, scRulesMap)
		}
		scRulesMapsStrting, _ := convert.ConvertListMapToJsonString(scRulesMaps)
		request["ScRules"] = helper.StringPointer(scRulesMapsStrting)
	}
	if v, ok := d.GetOk("dubbo_rules"); ok {
		dubboRulesMaps := make([]map[string]interface{}, 0)
		for _, dubboRules := range v.(*schema.Set).List() {
			dubboRulesArg := dubboRules.(map[string]interface{})
			dubboRulesMap := map[string]interface{}{}
			dubboRulesMap["condition"] = dubboRulesArg["condition"]
			dubboRulesMap["methodName"] = dubboRulesArg["method_name"]
			dubboRulesMap["serviceName"] = dubboRulesArg["service_name"]
			dubboRulesMap["version"] = dubboRulesArg["version"]
			dubboRulesMap["group"] = dubboRulesArg["group"]
			itemsMaps := make([]map[string]interface{}, 0)
			for _, items := range dubboRulesArg["items"].(*schema.Set).List() {
				itemsArg := items.(map[string]interface{})
				itemsMap := map[string]interface{}{}
				itemsMap["index"] = itemsArg["index"]
				itemsMap["expr"] = itemsArg["expr"]
				itemsMap["cond"] = itemsArg["cond"]
				itemsMap["value"] = itemsArg["value"]
				itemsMap["operator"] = itemsArg["operator"]
				itemsMaps = append(itemsMaps, itemsMap)
			}
			dubboRulesMap["items"] = itemsMaps
			dubboRulesMaps = append(dubboRulesMaps, dubboRulesMap)
		}
		dubboRulesMapsStrting, _ := convert.ConvertListMapToJsonString(dubboRulesMaps)
		request["DubboRules"] = helper.StringPointer(dubboRulesMapsStrting)
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_sae_grey_tag_route", "POST "+action, errs.AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["GreyTagRouteId"]))

	return resourceAliCloudSaeGreyTagRouteRead(d, meta)
}
func resourceAliCloudSaeGreyTagRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := NewService(client)
	object, err := saeService.DescribeSaeGreyTagRoute(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_grey_tag_route saeService.DescribeSaeGreyTagRoute Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("grey_tag_route_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("app_id", object["AppId"])
	if scRulesList, ok := object["ScRules"]; ok {
		scRulesMaps := make([]map[string]interface{}, 0)
		for _, scRules := range scRulesList.([]interface{}) {
			scRulesArg := scRules.(map[string]interface{})
			scRulesMap := map[string]interface{}{}
			scRulesMap["path"] = scRulesArg["path"]
			scRulesMap["condition"] = scRulesArg["condition"]
			itemsMaps := make([]map[string]interface{}, 0)
			for _, items := range scRulesArg["items"].([]interface{}) {
				itemsArg := items.(map[string]interface{})
				itemsMap := map[string]interface{}{}
				itemsMap["name"] = itemsArg["name"]
				itemsMap["cond"] = itemsArg["cond"]
				itemsMap["type"] = itemsArg["type"]
				itemsMap["value"] = itemsArg["value"]
				itemsMap["operator"] = itemsArg["operator"]
				itemsMaps = append(itemsMaps, itemsMap)
			}
			scRulesMap["items"] = itemsMaps
			scRulesMaps = append(scRulesMaps, scRulesMap)
		}

		d.Set("sc_rules", scRulesMaps)
	}

	if v, ok := object["DubboRules"]; ok {
		dubboRulesMaps := make([]map[string]interface{}, 0)
		for _, dubboRules := range v.([]interface{}) {
			dubboRulesArg := dubboRules.(map[string]interface{})
			dubboRulesMap := map[string]interface{}{}
			dubboRulesMap["condition"] = dubboRulesArg["condition"]
			dubboRulesMap["method_name"] = dubboRulesArg["methodName"]
			dubboRulesMap["service_name"] = dubboRulesArg["serviceName"]
			dubboRulesMap["version"] = dubboRulesArg["version"]
			dubboRulesMap["group"] = dubboRulesArg["group"]
			itemsMaps := make([]map[string]interface{}, 0)
			for _, items := range dubboRulesArg["items"].([]interface{}) {
				itemsArg := items.(map[string]interface{})
				itemsMap := map[string]interface{}{}
				itemsMap["index"] = itemsArg["index"]
				itemsMap["expr"] = itemsArg["expr"]
				itemsMap["cond"] = itemsArg["cond"]
				itemsMap["value"] = itemsArg["value"]
				itemsMap["operator"] = itemsArg["operator"]
				itemsMaps = append(itemsMaps, itemsMap)
			}
			dubboRulesMap["items"] = itemsMaps
			dubboRulesMaps = append(dubboRulesMaps, dubboRulesMap)
		}

		d.Set("dubbo_rules", dubboRulesMaps)
	}

	return nil
}
func resourceAliCloudSaeGreyTagRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]*string{
		"GreyTagRouteId": helper.StringPointer(d.Id()),
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = helper.StringPointer(v.(string))
	}

	if d.HasChange("sc_rules") {
		update = true
	}
	if v, ok := d.GetOk("sc_rules"); ok {
		scRulesMaps := make([]map[string]interface{}, 0)
		for _, scRules := range v.(*schema.Set).List() {
			scRulesArg := scRules.(map[string]interface{})
			scRulesMap := map[string]interface{}{}
			scRulesMap["path"] = scRulesArg["path"]
			scRulesMap["condition"] = scRulesArg["condition"]
			itemsMaps := make([]map[string]interface{}, 0)
			for _, items := range scRulesArg["items"].(*schema.Set).List() {
				itemsArg := items.(map[string]interface{})
				itemsMap := map[string]interface{}{}
				itemsMap["name"] = itemsArg["name"]
				itemsMap["cond"] = itemsArg["cond"]
				itemsMap["type"] = itemsArg["type"]
				itemsMap["value"] = itemsArg["value"]
				itemsMap["operator"] = itemsArg["operator"]
				itemsMaps = append(itemsMaps, itemsMap)
			}
			scRulesMap["items"] = itemsMaps
			scRulesMaps = append(scRulesMaps, scRulesMap)
		}
		scRulesMapsStrting, _ := convert.ConvertListMapToJsonString(scRulesMaps)
		request["ScRules"] = helper.StringPointer(scRulesMapsStrting)
	}

	if d.HasChange("dubbo_rules") {
		update = true
	}
	if v, ok := d.GetOk("dubbo_rules"); ok {
		dubboRulesMaps := make([]map[string]interface{}, 0)
		for _, dubboRules := range v.(*schema.Set).List() {
			dubboRulesArg := dubboRules.(map[string]interface{})
			dubboRulesMap := map[string]interface{}{}
			dubboRulesMap["condition"] = dubboRulesArg["condition"]
			dubboRulesMap["methodName"] = dubboRulesArg["method_name"]
			dubboRulesMap["serviceName"] = dubboRulesArg["service_name"]
			dubboRulesMap["version"] = dubboRulesArg["version"]
			dubboRulesMap["group"] = dubboRulesArg["group"]
			itemsMaps := make([]map[string]interface{}, 0)
			for _, items := range dubboRulesArg["items"].(*schema.Set).List() {
				itemsArg := items.(map[string]interface{})
				itemsMap := map[string]interface{}{}
				itemsMap["index"] = itemsArg["index"]
				itemsMap["expr"] = itemsArg["expr"]
				itemsMap["cond"] = itemsArg["cond"]
				itemsMap["value"] = itemsArg["value"]
				itemsMap["operator"] = itemsArg["operator"]
				itemsMaps = append(itemsMaps, itemsMap)
			}
			dubboRulesMap["items"] = itemsMaps
			dubboRulesMaps = append(dubboRulesMaps, dubboRulesMap)
		}
		dubboRulesMapsStrting, _ := convert.ConvertListMapToJsonString(dubboRulesMaps)
		request["DubboRules"] = helper.StringPointer(dubboRulesMapsStrting)
	}

	if update {
		action := "/pop/v1/sam/tagroute/greyTagRoute"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("sae", "2019-05-06", action, request, nil, nil, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "PUT "+action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)
	}
	return resourceAliCloudSaeGreyTagRouteRead(d, meta)
}
func resourceAliCloudSaeGreyTagRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/tagroute/greyTagRoute"
	var response map[string]interface{}
	var err error
	request := map[string]*string{
		"GreyTagRouteId": helper.StringPointer(d.Id()),
	}

	wait := helper.IncrementalWait(3*time.Second, 1*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DELETE "+action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	return nil
}

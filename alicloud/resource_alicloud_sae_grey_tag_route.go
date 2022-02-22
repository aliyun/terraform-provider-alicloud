package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSaeGreyTagRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeGreyTagRouteCreate,
		Read:   resourceAlicloudSaeGreyTagRouteRead,
		Update: resourceAlicloudSaeGreyTagRouteUpdate,
		Delete: resourceAlicloudSaeGreyTagRouteDelete,
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
			"dubbo_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
										ValidateFunc: validation.StringInSlice([]string{">", "<", ">=", "<=", "==", "!="}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"operator": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"rawvalue", "list", "mod", "deterministic_proportional_steaming_division"}, false),
									},
								},
							},
						},
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
							ValidateFunc: validation.StringInSlice([]string{"AND", "OR"}, false),
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"grey_tag_route_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sc_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"items": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cond": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{">", "<", ">=", "<=", "==", "!="}, false),
									},
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"param", "cookie", "header"}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"operator": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"rawvalue", "list", "mod", "deterministic_proportional_steaming_division"}, false),
									},
								},
							},
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"condition": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"AND", "OR"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudSaeGreyTagRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/pop/v1/sam/tagroute/greyTagRoute"
	request := make(map[string]*string)
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	request["AppId"] = StringPointer(d.Get("app_id").(string))
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = StringPointer(v.(string))
	}
	request["Name"] = StringPointer(d.Get("grey_tag_route_name").(string))
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
		scRulesMapsStrting, _ := convertListMapToJsonString(scRulesMaps)
		request["ScRules"] = StringPointer(scRulesMapsStrting)
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
		dubboRulesMapsStrting, _ := convertListMapToJsonString(dubboRulesMaps)
		request["DubboRules"] = StringPointer(dubboRulesMapsStrting)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_grey_tag_route", "POST "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	addDebug(action, response, request)
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["GreyTagRouteId"]))

	return resourceAlicloudSaeGreyTagRouteRead(d, meta)
}
func resourceAlicloudSaeGreyTagRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	object, err := saeService.DescribeSaeGreyTagRoute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_grey_tag_route saeService.DescribeSaeGreyTagRoute Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
func resourceAlicloudSaeGreyTagRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]*string{
		"GreyTagRouteId": StringPointer(d.Id()),
	}
	if d.HasChange("sc_rules") {
		update = true
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
			scRulesMapsStrting, _ := convertListMapToJsonString(scRulesMaps)
			request["ScRules"] = StringPointer(scRulesMapsStrting)
		}
	}
	if d.HasChange("dubbo_rules") {
		update = true
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
			dubboRulesMapsStrting, _ := convertListMapToJsonString(dubboRulesMaps)
			request["DubboRules"] = StringPointer(dubboRulesMapsStrting)
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = StringPointer(v.(string))
		}
	}
	if update {
		action := "/pop/v1/sam/tagroute/greyTagRoute"
		conn, err := client.NewServerlessClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PUT "+action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
		}
		addDebug(action, response, request)
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
		}
	}
	return resourceAlicloudSaeGreyTagRouteRead(d, meta)
}
func resourceAlicloudSaeGreyTagRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/tagroute/greyTagRoute"
	var response map[string]interface{}
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]*string{
		"GreyTagRouteId": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DELETE "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "DELETE "+action, response))
	}
	addDebug(action, response, request)
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "AlicloudSaeGreyTagRouteDelete", response))
	}
	return nil
}

package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSaeGreyTagRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSaeGreyTagRoutesRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dubbo_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"items": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"index": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"expr": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cond": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"method_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"condition": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grey_tag_route_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sc_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"items": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cond": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"condition": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceAlicloudSaeGreyTagRoutesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/tagroute/greyTagRouteList"
	request := make(map[string]*string)
	request["PageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	request["CurrentPage"] = StringPointer("1")
	request["AppId"] = StringPointer(d.Get("app_id").(string))
	var objects []map[string]interface{}
	var namespaceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		namespaceNameRegex = r
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
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sae_grey_tag_routes", action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "AlicloudSaeGreyTagRoutesRead", response))
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.Data.Result", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Result", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if namespaceNameRegex != nil && !namespaceNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["GreyTagRouteId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		currntPage, err := strconv.Atoi(*request["CurrentPage"])
		if err != nil {
			return WrapError(err)
		}
		request["CurrentPage"] = StringPointer(strconv.Itoa(currntPage + 1))
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"grey_tag_route_name": object["Name"],
			"id":                  fmt.Sprint(object["GreyTagRouteId"]),
			"description":         object["Description"],
		}

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
			mapping["sc_rules"] = scRulesMaps
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

			mapping["dubbo_rules"] = dubboRulesMaps
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("routes", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

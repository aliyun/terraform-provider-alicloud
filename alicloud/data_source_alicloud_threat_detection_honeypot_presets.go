package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudThreatDetectionHoneypotPresets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionHoneypotPresetsRead,
		Schema: map[string]*schema.Schema{
			"current_page": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeInt,
			},
			"honeypot_image_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"lang": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"node_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"node_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"preset_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  20,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"presets": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"honeypot_image_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"honeypot_preset_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"meta": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"portrait_option": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"burp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"trojan_git": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"node_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"preset_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionHoneypotPresetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("current_page"); ok {
		request["CurrentPage"] = v
	}
	if v, ok := d.GetOk("honeypot_image_name"); ok {
		request["HoneypotImageName"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("node_name"); ok {
		request["NodeName"] = v
	}
	if v, ok := d.GetOk("preset_name"); ok {
		request["PresetName"] = v
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["CurrentPage"] = v.(int)
	} else {
		request["CurrentPage"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}

	nodeId, nodeIdOk := d.GetOk("node_id")

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListHoneypotPreset"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_honeypot_presets", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.List", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.List", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["HoneypotPresetId"])]; !ok {
					continue
				}
			}
			if nodeIdOk && nodeId.(string) != "" && nodeId.(string) != item["NodeId"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                  fmt.Sprint(object["HoneypotPresetId"]),
			"honeypot_image_name": object["HoneypotImageName"],
			"honeypot_preset_id":  object["HoneypotPresetId"],
			"node_id":             object["NodeId"],
			"preset_name":         object["PresetName"],
		}

		ids = append(ids, fmt.Sprint(object["HoneypotPresetId"]))

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		threatDetectionService := ThreatDetectionService{client}
		id := fmt.Sprint(object["HoneypotPresetId"])
		getResp, err := threatDetectionService.DescribeThreatDetectionHoneypotPreset(id)
		if err != nil {
			return WrapError(err)
		}
		if v, ok := getResp["Meta"]; ok {
			metaMap, err := convertJsonStringToMap(v.(string))
			if err != nil {
				return WrapError(err)
			}
			metaList := make([]map[string]interface{}, 0)
			metaList = append(metaList, map[string]interface{}{
				"portrait_option": metaMap["portrait_option"],
				"burp":            metaMap["burp"],
				"trojan_git":      metaMap["trojan_git"],
			})
			mapping["meta"] = metaList
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("presets", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

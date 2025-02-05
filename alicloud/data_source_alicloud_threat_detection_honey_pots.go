package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudThreatDetectionHoneyPots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionHoneyPotsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"honeypot_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"honeypot_name": {
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
			"pots": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"honeypot_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"honeypot_image_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"honeypot_image_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"honeypot_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"node_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"preset_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"state": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionHoneyPotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	var honeypotNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		honeypotNameRegex = r
	}
	if v, ok := d.GetOk("honeypot_id"); ok {
		request["HoneypotIds"] = v
	}
	if v, ok := d.GetOk("honeypot_name"); ok {
		request["HoneypotName"] = v
	}
	if v, ok := d.GetOk("node_id"); ok {
		request["NodeId"] = v
	}
	if v, ok := d.GetOk("node_name"); ok {
		request["NodeName"] = v
	}
	request["CurrentPage"] = 1
	request["PageSize"] = PageSizeMedium

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
		action := "ListHoneypot"
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_honey_pots", action, AlibabaCloudSdkGoERROR)
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
				if _, ok := idsMap[fmt.Sprint(item["HoneypotId"])]; !ok {
					continue
				}
			}
			if honeypotNameRegex != nil && !honeypotNameRegex.MatchString(fmt.Sprint(item["HoneypotName"])) {
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
			"id":                  fmt.Sprint(object["HoneypotId"]),
			"honeypot_id":         object["HoneypotId"],
			"honeypot_image_id":   object["HoneypotImageId"],
			"honeypot_image_name": object["HoneypotImageName"],
			"honeypot_name":       object["HoneypotName"],
			"node_id":             object["NodeId"],
			"preset_id":           object["PresetId"],
			"state":               object["State"].([]interface{}),
			"status":              object["State"].([]interface{})[0],
		}

		ids = append(ids, fmt.Sprint(object["HoneypotId"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("pots", s); err != nil {
		return WrapError(err)
	}
	return nil
}

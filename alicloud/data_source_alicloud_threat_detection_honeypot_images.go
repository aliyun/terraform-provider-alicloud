package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudThreatDetectionHoneypotImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionHoneypotImagesRead,
		Schema: map[string]*schema.Schema{
			"node_id": {
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
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
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
			"images": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"honeypot_image_display_name": {
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
						"honeypot_image_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"honeypot_image_version": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"multiports": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"proto": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"service_port": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"template": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionHoneypotImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("node_id"); ok {
		request["NodeId"] = v
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

	var honeypotImageNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		honeypotImageNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}
	action := "ListAvailableHoneypot"
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_honeypot_images", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["HoneypotImageId"])]; !ok {
				continue
			}
		}

		if honeypotImageNameRegex != nil && !honeypotImageNameRegex.MatchString(fmt.Sprint(item["HoneypotImageName"])) {
			continue
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                          fmt.Sprint(object["HoneypotImageId"]),
			"honeypot_image_display_name": object["HoneypotImageDisplayName"],
			"honeypot_image_id":           object["HoneypotImageId"],
			"honeypot_image_name":         object["HoneypotImageName"],
			"honeypot_image_type":         object["HoneypotImageType"],
			"honeypot_image_version":      object["HoneypotImageVersion"],
			"multiports":                  object["Multiports"],
			"proto":                       object["Proto"],
			"service_port":                object["ServicePort"],
			"template":                    object["Template"],
		}

		ids = append(ids, fmt.Sprint(object["HoneypotImageId"]))
		names = append(names, object["HoneypotImageName"])

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("images", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

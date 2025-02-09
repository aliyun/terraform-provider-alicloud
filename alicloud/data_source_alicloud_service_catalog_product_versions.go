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

func dataSourceAlicloudServiceCatalogProductVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudServiceCatalogProductVersionsRead,
		Schema: map[string]*schema.Schema{
			"product_id": {
				Required: true,
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
			"enable_details": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"versions": {
				Deprecated: "Field 'versions' has been deprecated from provider version 1.197.0.",
				Computed:   true,
				Type:       schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"active": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"guidance": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_version_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_version_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"template_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"template_url": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"product_versions": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"active": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"guidance": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_version_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_version_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"template_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"template_url": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudServiceCatalogProductVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("product_id"); ok {
		request["ProductId"] = v
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

	var productVersionNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		productVersionNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}
	action := "ListProductVersions"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("servicecatalog", "2021-09-01", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_service_catalog_product_versions", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.ProductVersionDetails", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ProductVersionDetails", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["ProductVersionId"])]; !ok {
				continue
			}
		}

		if productVersionNameRegex != nil && !productVersionNameRegex.MatchString(fmt.Sprint(item["ProductVersionName"])) {
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
			"id":                   fmt.Sprint(object["ProductVersionId"]),
			"active":               object["Active"],
			"create_time":          object["CreateTime"],
			"description":          object["Description"],
			"guidance":             object["Guidance"],
			"product_version_id":   object["ProductVersionId"],
			"product_version_name": object["ProductVersionName"],
			"template_type":        object["TemplateType"],
			"template_url":         object["TemplateUrl"],
			"product_id":           object["ProductId"],
		}

		ids = append(ids, fmt.Sprint(object["ProductVersionId"]))
		names = append(names, object["ProductVersionName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("versions", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("product_versions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

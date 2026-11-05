// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package sae

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAliCloudSaeApplicationGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSaeApplicationGroupRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"package_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"web_container": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"jdk": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"running_instances": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"edas_container_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudSaeApplicationGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	appId := d.Get("app_id")
	// DescribeApplicationGroups
	action := fmt.Sprintf("/pop/v1/sam/app/describeApplicationGroups")
	request := make(map[string]*string)
	request["PageSize"] = helper.StringPointer(strconv.Itoa(helper.PageSizeLarge))
	request["CurrentPage"] = helper.StringPointer("1")
	request["AppId"] = helper.StringPointer(appId.(string))

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var groupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}

		groupNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_sae_application_groups", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["GroupId"])]; !ok {
					continue
				}
			}

			if groupNameRegex != nil {
				if !groupNameRegex.MatchString(fmt.Sprint(item["GroupName"])) {
					continue
				}
			}

			objects = append(objects, item)
		}

		if len(result) < helper.PageSizeLarge {
			break
		}

		currentPage, err := strconv.Atoi(*request["CurrentPage"])
		if err != nil {
			return errs.WrapError(err)
		}

		request["CurrentPage"] = helper.StringPointer(strconv.Itoa(currentPage + 1))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                     fmt.Sprint(object["GroupId"]),
			"group_id":               fmt.Sprint(object["GroupId"]),
			"group_name":             object["GroupName"],
			"group_type":             object["GroupType"],
			"replicas":               object["Replicas"],
			"package_type":           object["PackageType"],
			"package_version":        object["PackageVersion"],
			"package_version_id":     object["PackageVersionId"],
			"package_url":            object["PackageUrl"],
			"image_url":              object["ImageUrl"],
			"web_container":          object["WebContainer"],
			"jdk":                    object["Jdk"],
			"running_instances":      object["RunningInstances"],
			"edas_container_version": object["EdasContainerVersion"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["GroupName"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

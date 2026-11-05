// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package arms

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceEnvironments() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
			"environment_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// lintignore: S006
			"tags": tags.TagsSchemaWithElements(),
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bind_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bind_vpc_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"environment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"environment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"environment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grafana_datasource_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grafana_folder_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"managed_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prometheus_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func DataSourceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListEnvironments"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("environment_type"); ok {
		request["EnvironmentType"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		tagsJson, err := convert.ConvertArrayObjectToJsonString(tagsMap)
		if err != nil {
			return errs.WrapError(err)
		}

		request["Tag"] = tagsJson
	}

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

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}

		nameRegex = r
	}

	var response map[string]interface{}
	var err error

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)

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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_arms_environments", action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data.Environments", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.Environments", response)
	}

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["EnvironmentName"])) {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["EnvironmentId"])]; !ok {
				continue
			}
		}

		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}
		mapping["id"] = objectRaw["EnvironmentId"]
		mapping["bind_resource_id"] = objectRaw["BindResourceId"]
		mapping["bind_resource_type"] = objectRaw["BindResourceType"]
		mapping["bind_vpc_cidr"] = objectRaw["BindVpcCidr"]
		mapping["environment_id"] = objectRaw["EnvironmentId"]
		mapping["environment_name"] = objectRaw["EnvironmentName"]
		mapping["environment_type"] = objectRaw["EnvironmentType"]
		mapping["grafana_datasource_uid"] = objectRaw["GrafanaDatasourceUid"]
		mapping["grafana_folder_uid"] = objectRaw["GrafanaFolderUid"]
		mapping["managed_type"] = objectRaw["ManagedType"]
		mapping["prometheus_instance_id"] = objectRaw["PrometheusInstanceId"]
		mapping["region_id"] = objectRaw["RegionId"]
		mapping["resource_group_id"] = objectRaw["ResourceGroupId"]
		mapping["user_id"] = objectRaw["UserId"]

		mapping["tags"] = tags.TagsToMap(objectRaw["Tags"])

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["EnvironmentName"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("environments", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

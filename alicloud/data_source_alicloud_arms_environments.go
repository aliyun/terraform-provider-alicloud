// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
	"time"
)

func dataSourceAliCloudArmsEnvironments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudArmsEnvironmentRead,
		Schema: map[string]*schema.Schema{
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
			"environment_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
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

func dataSourceAliCloudArmsEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListEnvironments"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("environment_type"); ok {
		request["EnvironmentType"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		tagsJson, err := convertArrayObjectToJsonString(tagsMap)
		if err != nil {
			return WrapError(err)
		}

		//request["Tags"] = tagsJson
		request["Tag"] = tagsJson
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data.Environments", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Environments", response)
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

		mapping["tags"] = tagsToMap(objectRaw["Tags"])

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["EnvironmentName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("environments", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

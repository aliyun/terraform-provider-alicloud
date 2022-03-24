package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcsImagePipelines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsImagePipelinesRead,
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
			"name": {
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pipelines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_image": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"base_image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"build_content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_instance_on_failure": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_pipeline_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_max_bandwidth_out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
						"add_account": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"to_region_id": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcsImagePipelinesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeImagePipelines"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	request["MaxResults"] = PageSizeLarge
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
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_image_pipelines", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ImagePipeline.ImagePipelineSet", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ImagePipeline.ImagePipelineSet", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ImagePipelineId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"base_image":                 object["BaseImage"],
			"base_image_type":            object["BaseImageType"],
			"build_content":              object["BuildContent"],
			"delete_instance_on_failure": object["DeleteInstanceOnFailure"],
			"description":                object["Description"],
			"image_name":                 object["ImageName"],
			"creation_time":              object["CreationTime"],
			"id":                         fmt.Sprint(object["ImagePipelineId"]),
			"image_pipeline_id":          fmt.Sprint(object["ImagePipelineId"]),
			"instance_type":              object["InstanceType"],
			"internet_max_bandwidth_out": formatInt(object["InternetMaxBandwidthOut"]),
			"name":                       object["Name"],
			"resource_group_id":          object["ResourceGroupId"],
			"system_disk_size":           formatInt(object["SystemDiskSize"]),
			"vswitch_id":                 object["VSwitchId"],
		}
		if v, ok := object["Tags"].(map[string]interface{}); ok {
			mapping["tags"] = tagsToMap(v["Tag"])
		}
		addAccountsList := make([]string, 0)
		if addAccounts, ok := object["AddAccounts"]; ok {
			addAccountsMap := addAccounts.(map[string]interface{})
			if addAccount, ok := addAccountsMap["AddAccount"]; ok {
				for _, item := range addAccount.([]interface{}) {
					addAccountsList = append(addAccountsList, item.(string))
				}
			}
		}
		mapping["add_account"] = addAccountsList

		toRegionIdList := make([]string, 0)
		if toRegionIds, ok := object["ToRegionIds"]; ok {
			toRegionIdsMap := toRegionIds.(map[string]interface{})
			if toRegionId, ok := toRegionIdsMap["ToRegionId"]; ok {
				for _, item := range toRegionId.([]interface{}) {
					toRegionIdList = append(toRegionIdList, item.(string))
				}
			}
		}
		mapping["to_region_id"] = toRegionIdList
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

	if err := d.Set("pipelines", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudSaeApplicationInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSaeApplicationInstancesRead,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_container_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_container_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_health_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"main_container_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSaeApplicationInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	appId := d.Get("application_id").(string)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	groupIds := make([]string, 0)
	if v, ok := d.GetOk("group_id"); ok {
		groupIds = append(groupIds, v.(string))
	} else {
		action := "/pop/v1/sam/app/describeApplicationGroups"
		request := map[string]*string{
			"AppId":       StringPointer(appId),
			"PageSize":    StringPointer(strconv.Itoa(PageSizeLarge)),
			"CurrentPage": StringPointer("1"),
		}
		var response map[string]interface{}
		var err error
		for {
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sae_application_instances", "GET "+action, AlibabaCloudSdkGoERROR)
			}
			resp, err := jsonpath.Get("$.Data", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				groupIds = append(groupIds, fmt.Sprint(item["GroupId"]))
			}
			if len(result) < PageSizeLarge {
				break
			}
			currentPage, err := strconv.Atoi(*request["CurrentPage"])
			if err != nil {
				return WrapError(err)
			}
			request["CurrentPage"] = StringPointer(strconv.Itoa(currentPage + 1))
		}
	}

	action := "/pop/v1/sam/app/describeApplicationInstances"
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, groupId := range groupIds {
		request := map[string]*string{
			"AppId":       StringPointer(appId),
			"GroupId":     StringPointer(groupId),
			"PageSize":    StringPointer(strconv.Itoa(PageSizeLarge)),
			"CurrentPage": StringPointer("1"),
		}
		var response map[string]interface{}
		var err error
		for {
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sae_application_instances", "GET "+action, AlibabaCloudSdkGoERROR)
			}
			resp, err := jsonpath.Get("$.Data.Instances", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Instances", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				if len(idsMap) > 0 {
					if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
						continue
					}
				}
				mapping := map[string]interface{}{
					"id":                        fmt.Sprint(item["InstanceId"]),
					"instance_id":               item["InstanceId"],
					"group_id":                  item["GroupId"],
					"instance_container_ip":     item["InstanceContainerIp"],
					"instance_container_status": item["InstanceContainerStatus"],
					"instance_health_status":    item["InstanceHealthStatus"],
					"main_container_status":     item["MainContainerStatus"],
					"image_url":                 item["ImageUrl"],
					"package_version":           item["PackageVersion"],
					"vswitch_id":                item["VSwitchId"],
				}
				ids = append(ids, fmt.Sprint(mapping["id"]))
				s = append(s, mapping)
			}
			if len(result) < PageSizeLarge {
				break
			}
			currentPage, err := strconv.Atoi(*request["CurrentPage"])
			if err != nil {
				return WrapError(err)
			}
			request["CurrentPage"] = StringPointer(strconv.Itoa(currentPage + 1))
		}
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

package sae

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceApplicationInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
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
						"instance_container_restarts": {
							Type:     schema.TypeInt,
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
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"debug_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unhealthy_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"main_container_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time_stamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"finish_time_stamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sidecar_containers_status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"container_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"container_status": {
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

func dataSourceApplicationInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	appId := d.Get("app_id")
	GroupId := d.Get("group_id")

	// DescribeApplicationInstances
	action := fmt.Sprintf("/pop/v1/sam/app/describeApplicationInstances")
	request := make(map[string]*string)
	request["PageSize"] = helper.StringPointer(strconv.Itoa(helper.PageSizeLarge))
	request["CurrentPage"] = helper.StringPointer("1")
	request["AppId"] = helper.StringPointer(appId.(string))
	request["GroupId"] = helper.StringPointer(GroupId.(string))

	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = helper.StringPointer(v.(string))
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
			if errs.IsExpectedErrors(err, []string{"sae.errorcode.app.instance.not.found.message"}) {
				return nil
			}
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_sae_application_instances", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Data.Instances", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.Instances", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                          fmt.Sprint(object["InstanceId"]),
			"instance_id":                 fmt.Sprint(object["InstanceId"]),
			"group_id":                    object["GroupId"],
			"instance_container_ip":       object["InstanceContainerIp"],
			"instance_container_restarts": object["InstanceContainerRestarts"],
			"instance_container_status":   object["InstanceContainerStatus"],
			"instance_health_status":      object["InstanceHealthStatus"],
			"vswitch_id":                  object["VSwitchId"],
			"package_version":             object["PackageVersion"],
			"image_url":                   object["ImageUrl"],
			"debug_status":                object["DebugStatus"],
			"eip":                         object["Eip"],
			"unhealthy_message":           object["UnhealthyMessage"],
			"traffic_status":              object["TrafficStatus"],
			"main_container_status":       object["MainContainerStatus"],
			"timestamp":                   object["Timestamp"],
			"create_time_stamp":           object["CreateTimeStamp"],
			"finish_time_stamp":           object["FinishTimeStamp"],
		}

		if sidecarContainersStatus, ok := object["SidecarContainersStatus"]; ok {
			sidecarContainersStatusMaps := make([]map[string]interface{}, 0)
			for _, sidecarContainersStatusList := range convert.ConvertToInterfaceArray(sidecarContainersStatus) {
				sidecarContainersStatusMap := map[string]interface{}{}
				sidecarContainersStatusArg := sidecarContainersStatusList.(map[string]interface{})

				if containerId, ok := sidecarContainersStatusArg["ContainerId"]; ok {
					sidecarContainersStatusMap["container_id"] = containerId
				}

				if imageUrl, ok := sidecarContainersStatusArg["ImageUrl"]; ok {
					sidecarContainersStatusMap["image_url"] = imageUrl
				}

				if containerStatus, ok := sidecarContainersStatusArg["ContainerStatus"]; ok {
					sidecarContainersStatusMap["container_status"] = containerStatus
				}

				sidecarContainersStatusMaps = append(sidecarContainersStatusMaps, sidecarContainersStatusMap)
			}

			mapping["sidecar_containers_status"] = sidecarContainersStatusMaps
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

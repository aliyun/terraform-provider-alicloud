package ecd

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDesktops() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDesktopsRead,
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
			"office_site_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desktop_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Starting", "Pending", "Running", "Rebuilding", "Stopping", "Stopped", "Deleted", "Expired"}, false),
			},
			"end_user_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"desktops": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_user_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"disks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"performance_level": {
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

func dataSourceDesktopsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDesktops"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["MaxResults"] = helper.PageSizeLarge

	if v, ok := d.GetOk("office_site_id"); ok {
		request["OfficeSiteId"] = v
	}

	if v, ok := d.GetOk("policy_group_id"); ok {
		request["PolicyGroupId"] = v
	}

	if v, ok := d.GetOk("desktop_name"); ok {
		request["DesktopName"] = v
	}

	if v, ok := d.GetOk("status"); ok {
		request["DesktopStatus"] = v
	}

	if v, ok := d.GetOk("end_user_ids"); ok {
		request["EndUserId"] = v
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

	var desktopNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}

		desktopNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
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
			if errs.IsExpectedErrors(err, []string{"InvalidDirectoryId.NotFound"}) {
				return nil
			}

			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ecd_desktops", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Desktops", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Desktops", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DesktopId"])]; !ok {
					continue
				}
			}

			if desktopNameRegex != nil && !desktopNameRegex.MatchString(fmt.Sprint(item["DesktopName"])) {
				continue
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
			"id":                   fmt.Sprint(object["DesktopId"]),
			"desktop_id":           fmt.Sprint(object["DesktopId"]),
			"desktop_type":         object["DesktopType"],
			"desktop_name":         object["DesktopName"],
			"policy_group_id":      object["PolicyGroupId"],
			"directory_id":         object["DirectoryId"],
			"network_interface_id": object["NetworkInterfaceId"],
			"image_id":             object["ImageId"],
			"memory":               fmt.Sprint(object["Memory"]),
			"payment_type":         convertEcdDesktopPaymentTypeResponse(object["ChargeType"]),
			"status":               object["DesktopStatus"],
			"create_time":          object["CreationTime"],
			"expired_time":         object["ExpiredTime"],
			"end_user_ids":         object["EndUserIds"],
		}

		cpu, err := convert.FormatInt(object["Cpu"])
		if err != nil {
			return err
		}

		mapping["cpu"] = cpu

		systemDiskSize, err := convert.FormatInt(object["SystemDiskSize"])
		if err != nil {
			return err
		}

		mapping["system_disk_size"] = systemDiskSize

		if disks, ok := object["Disks"]; ok {
			disksMaps := make([]map[string]interface{}, 0)
			for _, disksList := range convert.ConvertToInterfaceArray(disks) {
				disksMap := map[string]interface{}{}

				if disksArg, ok := disksList.(map[string]interface{}); ok {
					if diskId, ok := disksArg["DiskId"]; ok {
						disksMap["disk_id"] = diskId
					}

					if diskType, ok := disksArg["DiskType"]; ok {
						disksMap["disk_type"] = diskType
					}

					if diskCategory, ok := disksArg["DiskCategory"]; ok {
						disksMap["disk_category"] = diskCategory
					}

					if diskSize, ok := disksArg["DiskSize"]; ok {
						disksMap["disk_size"] = diskSize
					}

					if performanceLevel, ok := disksArg["PerformanceLevel"]; ok {
						disksMap["performance_level"] = performanceLevel
					}
				}

				disksMaps = append(disksMaps, disksMap)
			}

			mapping["disks"] = disksMaps
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DesktopName"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("desktops", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

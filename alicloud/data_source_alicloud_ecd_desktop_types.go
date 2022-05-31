package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcdDesktopTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdDesktopTypesRead,
		Schema: map[string]*schema.Schema{
			"cpu_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"gpu_count": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"instance_type_family": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"eds.graphics", "eds.hf", "eds.general", "ecd.graphics", "ecd.performance", "ecd.advanced", "ecd.basic"}, false),
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"SUFFICIENT"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_disk_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_type_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu_count": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"gpu_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcdDesktopTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDesktopTypes"
	request := make(map[string]interface{})
	if v, ok := d.GetOkExists("cpu_count"); ok {
		request["CpuCount"] = v
	}
	if v, ok := d.GetOkExists("gpu_count"); ok {
		request["GpuCount"] = v
	}
	if v, ok := d.GetOk("instance_type_family"); ok {
		request["InstanceTypeFamily"] = v
	}
	if v, ok := d.GetOkExists("memory_size"); ok {
		request["MemorySize"] = v
	}
	request["RegionId"] = client.RegionId
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_desktop_types", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.DesktopTypes", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DesktopTypes", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["DesktopTypeId"])]; !ok {
				continue
			}
		}
		if statusOk && status.(string) != "" && status.(string) != item["DesktopTypeStatus"].(string) {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"cpu_count":            object["CpuCount"],
			"data_disk_size":       object["DataDiskSize"],
			"id":                   fmt.Sprint(object["DesktopTypeId"]),
			"desktop_type_id":      fmt.Sprint(object["DesktopTypeId"]),
			"gpu_count":            object["GpuCount"],
			"gpu_spec":             object["GpuSpec"],
			"instance_type_family": object["InstanceTypeFamily"],
			"memory_size":          object["MemorySize"],
			"status":               object["DesktopTypeStatus"],
			"system_disk_size":     object["SystemDiskSize"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("types", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

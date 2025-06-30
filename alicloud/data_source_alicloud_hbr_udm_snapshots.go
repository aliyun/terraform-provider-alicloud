// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudHbrUdmSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudHbrUdmSnapshotRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"disk_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"udm_snapshot_id": {
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
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudHbrUdmSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeUdmSnapshots"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["UdmRegionId"] = client.RegionId
	if v, ok := d.GetOk("disk_id"); ok {
		request["DiskId"] = v
	}
	request["EndTime"] = d.Get("end_time")
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("job_id"); ok {
		request["JobId"] = v
	}
	request["SourceType"] = d.Get("source_type")
	request["StartTime"] = d.Get("start_time")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	resp, _ := jsonpath.Get("$.Snapshots[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(response["$.Snapshots[0].SnapshotId"])]; !ok {
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

		mapping["id"] = objectRaw["SnapshotId"]

		mapping["disk_id"] = objectRaw["DiskId"]
		mapping["instance_id"] = objectRaw["InstanceId"]
		mapping["job_id"] = objectRaw["JobId"]
		mapping["source_type"] = objectRaw["SourceType"]
		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["start_time"] = objectRaw["StartTime"]
		mapping["udm_snapshot_id"] = objectRaw["SnapshotId"]

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw[""])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("snapshots", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

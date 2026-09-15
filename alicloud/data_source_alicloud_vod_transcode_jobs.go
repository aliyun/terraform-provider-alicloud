package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudVodTranscodeJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVodTranscodeJobsRead,
		Schema: map[string]*schema.Schema{
			"video_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snapshot_type": {
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
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transcode_job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"video_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
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
		},
	}
}

func dataSourceAlicloudVodTranscodeJobsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ListSnapshots"
	request := map[string]interface{}{
		"VideoId":  d.Get("video_id"),
		"PageSize": PageSizeLarge,
		"PageNo":   1,
	}
	if v, ok := d.GetOk("snapshot_type"); ok {
		request["SnapshotType"] = v
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
	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vod_transcode_jobs", action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.MediaSnapshot.Snapshots.Snapshot", response)
		if err != nil {
			break
		}
		result, ok := v.([]interface{})
		if !ok || len(result) == 0 {
			break
		}
		for _, item := range result {
			snapshot, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			jobId := fmt.Sprint(snapshot["JobId"])
			if len(idsMap) > 0 {
				if _, ok := idsMap[jobId]; !ok {
					continue
				}
			}
			objects = append(objects, snapshot)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNo"] = request["PageNo"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		jobId := fmt.Sprint(object["JobId"])
		mapping := map[string]interface{}{
			"transcode_job_id": jobId,
			"video_id":         d.Get("video_id"),
			"create_time":      object["CreateTime"],
			"id":               jobId,
		}
		ids = append(ids, jobId)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("jobs", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

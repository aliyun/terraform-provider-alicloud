package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVodTranscodeJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVodTranscodeJobCreate,
		Read:   resourceAlicloudVodTranscodeJobRead,
		Delete: resourceAlicloudVodTranscodeJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"video_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reference_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"height": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"width": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"specified_offset_time": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"sprite_snapshot_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"specified_offset_times": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"transcode_job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVodTranscodeJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "SubmitSnapshotJob"
	request := make(map[string]interface{})
	request["VideoId"] = d.Get("video_id")
	if v, ok := d.GetOk("reference_id"); ok {
		request["ReferenceId"] = v
	}
	if v, ok := d.GetOk("height"); ok {
		request["Height"] = v
	}
	if v, ok := d.GetOk("width"); ok {
		request["Width"] = v
	}
	if v, ok := d.GetOk("specified_offset_time"); ok {
		request["SpecifiedOffsetTime"] = v
	}
	if v, ok := d.GetOk("count"); ok {
		request["Count"] = v
	}
	if v, ok := d.GetOk("interval"); ok {
		request["Interval"] = v
	}
	if v, ok := d.GetOk("sprite_snapshot_config"); ok {
		request["SpriteSnapshotConfig"] = v
	}
	if v, ok := d.GetOk("snapshot_template_id"); ok {
		request["SnapshotTemplateId"] = v
	}
	if v, ok := d.GetOk("user_data"); ok {
		request["UserData"] = v
	}
	if v, ok := d.GetOk("specified_offset_times"); ok {
		times := v.([]interface{})
		for i, t := range times {
			request[fmt.Sprintf("SpecifiedOffsetTimes.%d", i+1)] = t
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	var err error
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vod_transcode_job", action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SnapshotJob", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_vod_transcode_job", "$.SnapshotJob", response)
	}
	snapshotJob, ok := v.(map[string]interface{})
	if !ok || snapshotJob == nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_vod_transcode_job", "$.SnapshotJob", response)
	}
	d.SetId(fmt.Sprint(snapshotJob["JobId"]))

	return resourceAlicloudVodTranscodeJobRead(d, meta)
}

func resourceAlicloudVodTranscodeJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := VodService{client}
	object, err := vodService.DescribeVodTranscodeJob(d.Id())
	if err != nil {
		if NotFoundError(err) {
			// Fallback: try ListSnapshots to find the snapshot job by VideoId
			videoId := d.Get("video_id").(string)
			if videoId != "" {
				object, err = vodService.DescribeVodSnapshot(videoId, d.Id())
				if err != nil {
					if NotFoundError(err) {
						log.Printf("[DEBUG] Resource alicloud_vod_transcode_job not found via GetTranscodeTask or ListSnapshots: %s", d.Id())
						d.SetId("")
						return nil
					}
					return WrapError(err)
				}
				d.Set("transcode_job_id", d.Id())
				if v, ok := object["CreateTime"]; ok && v != nil {
					d.Set("create_time", fmt.Sprint(v))
				}
				d.Set("video_id", videoId)
				d.Set("region_id", client.RegionId)
				return nil
			}
			log.Printf("[DEBUG] Resource alicloud_vod_transcode_job not found: %s", d.Id())
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("transcode_job_id", d.Id())
	if v, ok := object["CreationTime"]; ok && v != nil {
		d.Set("create_time", fmt.Sprint(v))
	}
	if v, ok := object["VideoId"]; ok && v != nil {
		d.Set("video_id", fmt.Sprint(v))
	}
	d.Set("region_id", client.RegionId)
	return nil
}

func resourceAlicloudVodTranscodeJobDelete(d *schema.ResourceData, meta interface{}) error {
	// Snapshot jobs are submitted asynchronously and produce snapshot results
	// stored in OSS. There is no API to delete a snapshot job; removing the
	// resource from Terraform state is sufficient. The snapshot results persist
	// in OSS until the video itself is deleted.
	return nil
}

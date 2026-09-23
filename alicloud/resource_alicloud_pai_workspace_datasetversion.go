// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPaiWorkspaceDatasetversion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspaceDatasetversionCreate,
		Read:   resourceAliCloudPaiWorkspaceDatasetversionRead,
		Update: resourceAliCloudPaiWorkspaceDatasetversionUpdate,
		Delete: resourceAliCloudPaiWorkspaceDatasetversionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"data_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"data_source_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dataset_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"property": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"FILE", "DIRECTORY"}, false),
			},
			"source_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudPaiWorkspaceDatasetversionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	DatasetId := d.Get("dataset_id")
	action := fmt.Sprintf("/api/v1/datasets/%s/versions", DatasetId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["DataSourceType"] = d.Get("data_source_type")
	request["Uri"] = d.Get("uri")
	request["Property"] = d.Get("property")
	if v, ok := d.GetOk("source_type"); ok {
		request["SourceType"] = v
	}
	if v, ok := d.GetOk("source_id"); ok {
		request["SourceId"] = v
	}
	if v, ok := d.GetOk("options"); ok {
		request["Options"] = v
	}
	if v, ok := d.GetOkExists("data_size"); ok {
		request["DataSize"] = v
	}
	if v, ok := d.GetOkExists("data_count"); ok {
		request["DataCount"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("labels"); ok {
		labelsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Key"] = dataLoopTmp["key"]
			dataLoopMap["Value"] = dataLoopTmp["value"]
			labelsMapsArray = append(labelsMapsArray, dataLoopMap)
		}
		request["Labels"] = labelsMapsArray
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_workspace_datasetversion", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", DatasetId, response["VersionName"]))

	return resourceAliCloudPaiWorkspaceDatasetversionRead(d, meta)
}

func resourceAliCloudPaiWorkspaceDatasetversionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspaceDatasetversion(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_datasetversion DescribePaiWorkspaceDatasetversion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["GmtCreateTime"] != nil {
		d.Set("create_time", objectRaw["GmtCreateTime"])
	}
	if objectRaw["DataCount"] != nil {
		d.Set("data_count", objectRaw["DataCount"])
	}
	if objectRaw["DataSize"] != nil {
		d.Set("data_size", objectRaw["DataSize"])
	}
	if objectRaw["DataSourceType"] != nil {
		d.Set("data_source_type", objectRaw["DataSourceType"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["Options"] != nil {
		d.Set("options", objectRaw["Options"])
	}
	if objectRaw["Property"] != nil {
		d.Set("property", objectRaw["Property"])
	}
	if objectRaw["SourceId"] != nil {
		d.Set("source_id", objectRaw["SourceId"])
	}
	if objectRaw["SourceType"] != nil {
		d.Set("source_type", objectRaw["SourceType"])
	}
	if objectRaw["Uri"] != nil {
		d.Set("uri", objectRaw["Uri"])
	}
	if objectRaw["DatasetId"] != nil {
		d.Set("dataset_id", objectRaw["DatasetId"])
	}
	if objectRaw["VersionName"] != nil {
		d.Set("version_name", objectRaw["VersionName"])
	}

	labels1Raw := objectRaw["Labels"]
	labelsMaps := make([]map[string]interface{}, 0)
	if labels1Raw != nil {
		for _, labelsChild1Raw := range labels1Raw.([]interface{}) {
			labelsMap := make(map[string]interface{})
			labelsChild1Raw := labelsChild1Raw.(map[string]interface{})
			labelsMap["key"] = labelsChild1Raw["Key"]
			labelsMap["value"] = labelsChild1Raw["Value"]

			labelsMaps = append(labelsMaps, labelsMap)
		}
	}
	if objectRaw["Labels"] != nil {
		if err := d.Set("labels", labelsMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceAliCloudPaiWorkspaceDatasetversionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	parts := strings.Split(d.Id(), ":")
	DatasetId := parts[0]
	VersionName := parts[1]
	action := fmt.Sprintf("/api/v1/datasets/%s/versions/%s", DatasetId, VersionName)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("options") {
		update = true
	}
	if v, ok := d.GetOk("options"); ok || d.HasChange("options") {
		request["Options"] = v
	}
	if d.HasChange("data_size") {
		update = true
	}
	if v, ok := d.GetOkExists("data_size"); ok || d.HasChange("data_size") {
		request["DataSize"] = v
	}
	if d.HasChange("data_count") {
		update = true
	}
	if v, ok := d.GetOkExists("data_count"); ok || d.HasChange("data_count") {
		request["DataCount"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["Description"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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
	}

	return resourceAliCloudPaiWorkspaceDatasetversionRead(d, meta)
}

func resourceAliCloudPaiWorkspaceDatasetversionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	DatasetId := parts[0]
	VersionName := parts[1]
	action := fmt.Sprintf("/api/v1/datasets/%s/versions/%s", DatasetId, VersionName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AIWorkSpace", "2021-02-04", action, query, nil, nil, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

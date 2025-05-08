// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPaiWorkspaceModelVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspaceModelVersionCreate,
		Read:   resourceAliCloudPaiWorkspaceModelVersionRead,
		Update: resourceAliCloudPaiWorkspaceModelVersionUpdate,
		Delete: resourceAliCloudPaiWorkspaceModelVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"approval_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"extra_info": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"format_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"framework_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"inference_spec": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"metrics": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"model_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"training_spec": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudPaiWorkspaceModelVersionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	ModelId := d.Get("model_id")
	action := fmt.Sprintf("/api/v1/models/%s/versions", ModelId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("version_name"); ok {
		request["VersionName"] = v
	}

	if v, ok := d.GetOk("extra_info"); ok {
		request["ExtraInfo"] = v
	}
	request["Uri"] = d.Get("uri")
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

	if v, ok := d.GetOk("version_description"); ok {
		request["VersionDescription"] = v
	}
	if v, ok := d.GetOk("format_type"); ok {
		request["FormatType"] = v
	}
	if v, ok := d.GetOk("framework_type"); ok {
		request["FrameworkType"] = v
	}
	if v, ok := d.GetOk("options"); ok {
		request["Options"] = v
	}
	if v, ok := d.GetOk("training_spec"); ok {
		request["TrainingSpec"] = v
	}
	if v, ok := d.GetOk("inference_spec"); ok {
		request["InferenceSpec"] = v
	}
	if v, ok := d.GetOk("source_type"); ok {
		request["SourceType"] = v
	}
	if v, ok := d.GetOk("source_id"); ok {
		request["SourceId"] = v
	}
	if v, ok := d.GetOk("approval_status"); ok {
		request["ApprovalStatus"] = v
	}
	if v, ok := d.GetOk("metrics"); ok {
		request["Metrics"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_workspace_model_version", action, AlibabaCloudSdkGoERROR)
	}

	VersionNameVar, _ := jsonpath.Get("$.VersionName", response)
	d.SetId(fmt.Sprintf("%v:%v", ModelId, VersionNameVar))

	return resourceAliCloudPaiWorkspaceModelVersionRead(d, meta)
}

func resourceAliCloudPaiWorkspaceModelVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspaceModelVersion(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_model_version DescribePaiWorkspaceModelVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("approval_status", objectRaw["ApprovalStatus"])
	d.Set("extra_info", objectRaw["ExtraInfo"])
	d.Set("format_type", objectRaw["FormatType"])
	d.Set("framework_type", objectRaw["FrameworkType"])
	d.Set("inference_spec", objectRaw["InferenceSpec"])
	d.Set("metrics", objectRaw["Metrics"])
	d.Set("options", objectRaw["Options"])
	d.Set("source_id", objectRaw["SourceId"])
	d.Set("source_type", objectRaw["SourceType"])
	d.Set("training_spec", objectRaw["TrainingSpec"])
	d.Set("uri", objectRaw["Uri"])
	d.Set("version_description", objectRaw["VersionDescription"])
	d.Set("version_name", objectRaw["VersionName"])

	labelsRaw := objectRaw["Labels"]
	labelsMaps := make([]map[string]interface{}, 0)
	if labelsRaw != nil {
		for _, labelsChildRaw := range labelsRaw.([]interface{}) {
			labelsMap := make(map[string]interface{})
			labelsChildRaw := labelsChildRaw.(map[string]interface{})
			labelsMap["key"] = labelsChildRaw["Key"]
			labelsMap["value"] = labelsChildRaw["Value"]

			labelsMaps = append(labelsMaps, labelsMap)
		}
	}
	if err := d.Set("labels", labelsMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("model_id", parts[0])

	return nil
}

func resourceAliCloudPaiWorkspaceModelVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	ModelId := parts[0]
	VersionName := parts[1]
	action := fmt.Sprintf("/api/v1/models/%s/versions/%s", ModelId, VersionName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("version_description") {
		update = true
	}
	if v, ok := d.GetOk("version_description"); ok {
		request["VersionDescription"] = v
	}
	if d.HasChange("metrics") {
		update = true
	}
	if v, ok := d.GetOk("metrics"); ok {
		request["Metrics"] = v
	}
	if d.HasChange("training_spec") {
		update = true
	}
	if v, ok := d.GetOk("training_spec"); ok {
		request["TrainingSpec"] = v
	}
	if d.HasChange("inference_spec") {
		update = true
	}
	if v, ok := d.GetOk("inference_spec"); ok {
		request["InferenceSpec"] = v
	}
	if d.HasChange("options") {
		update = true
	}
	if v, ok := d.GetOk("options"); ok {
		request["Options"] = v
	}
	if d.HasChange("source_type") {
		update = true
	}
	if v, ok := d.GetOk("source_type"); ok {
		request["SourceType"] = v
	}
	if d.HasChange("source_id") {
		update = true
	}
	if v, ok := d.GetOk("source_id"); ok {
		request["SourceId"] = v
	}
	if d.HasChange("approval_status") {
		update = true
	}
	if v, ok := d.GetOk("approval_status"); ok {
		request["ApprovalStatus"] = v
	}
	if d.HasChange("extra_info") {
		update = true
	}
	if v, ok := d.GetOk("extra_info"); ok {
		request["ExtraInfo"] = v
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

	if d.HasChange("labels") {
		oldEntry, newEntry := d.GetChange("labels")
		removed := oldEntry
		added := newEntry

		if len(removed.([]interface{})) > 0 {
			action = fmt.Sprintf("/api/v1/models/%s/versions/%s/labels", ModelId, VersionName)
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})

			localData := removed.([]interface{})
			labelsMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				labelsMapsArray = append(labelsMapsArray, dataLoopTmp["key"])
			}
			query["LabelKeys"] = StringPointer(convertListToCommaSeparate(labelsMapsArray))

			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaDelete("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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

		if len(added.([]interface{})) > 0 {
			action = fmt.Sprintf("/api/v1/models/%s/versions/%s/labels", ModelId, VersionName)
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})

			localData := added.([]interface{})
			labelsMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Value"] = dataLoopTmp["value"]
				dataLoopMap["Key"] = dataLoopTmp["key"]
				labelsMapsArray = append(labelsMapsArray, dataLoopMap)
			}
			request["Labels"] = labelsMapsArray

			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}

	d.Partial(false)
	return resourceAliCloudPaiWorkspaceModelVersionRead(d, meta)
}

func resourceAliCloudPaiWorkspaceModelVersionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	ModelId := parts[0]
	VersionName := parts[1]
	action := fmt.Sprintf("/api/v1/models/%s/versions/%s", ModelId, VersionName)
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

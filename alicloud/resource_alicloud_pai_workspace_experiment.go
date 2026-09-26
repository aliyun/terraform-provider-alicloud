// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPaiWorkspaceExperiment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspaceExperimentCreate,
		Read:   resourceAliCloudPaiWorkspaceExperimentRead,
		Update: resourceAliCloudPaiWorkspaceExperimentUpdate,
		Delete: resourceAliCloudPaiWorkspaceExperimentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accessibility": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PRIVATE", "PUBLIC"}, false),
			},
			"artifact_uri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"experiment_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^\\d+$"), "WorkspaceId is the workspace id which contains the experiment"),
			},
		},
	}
}

func resourceAliCloudPaiWorkspaceExperimentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/experiments")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ArtifactUri"] = d.Get("artifact_uri")
	request["WorkspaceId"] = d.Get("workspace_id")
	if v, ok := d.GetOk("accessibility"); ok {
		request["Accessibility"] = v
	}
	request["Name"] = d.Get("experiment_name")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_workspace_experiment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ExperimentId"]))

	return resourceAliCloudPaiWorkspaceExperimentRead(d, meta)
}

func resourceAliCloudPaiWorkspaceExperimentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspaceExperiment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_experiment DescribePaiWorkspaceExperiment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Accessibility"] != nil {
		d.Set("accessibility", objectRaw["Accessibility"])
	}
	if objectRaw["ArtifactUri"] != nil {
		d.Set("artifact_uri", objectRaw["ArtifactUri"])
	}
	if objectRaw["GmtCreateTime"] != nil {
		d.Set("create_time", objectRaw["GmtCreateTime"])
	}
	if objectRaw["Name"] != nil {
		d.Set("experiment_name", objectRaw["Name"])
	}
	if objectRaw["WorkspaceId"] != nil {
		d.Set("workspace_id", objectRaw["WorkspaceId"])
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

func resourceAliCloudPaiWorkspaceExperimentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	ExperimentId := d.Id()
	action := fmt.Sprintf("/api/v1/experiments/%s", ExperimentId)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["ExperimentId"] = d.Id()

	if d.HasChange("accessibility") {
		update = true
	}
	if v, ok := d.GetOk("accessibility"); ok || d.HasChange("accessibility") {
		request["Accessibility"] = v
	}
	if d.HasChange("experiment_name") {
		update = true
	}
	request["Name"] = d.Get("experiment_name")
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
			ExperimentId := d.Id()
			localData := removed.([]interface{})
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				key := dataLoopTmp["key"].(string)
				action := fmt.Sprintf("/api/v1/experiments/%s/labels/%s", ExperimentId, key)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["ExperimentId"] = d.Id()
				request["Key"] = key
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		if len(added.([]interface{})) > 0 {
			ExperimentId := d.Id()
			action := fmt.Sprintf("/api/v1/experiments/%s/labels", ExperimentId)
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			request["ExperimentId"] = d.Id()
			localData := added.([]interface{})
			labelsMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Key"] = dataLoopTmp["key"]
				dataLoopMap["Value"] = dataLoopTmp["value"]
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

	return resourceAliCloudPaiWorkspaceExperimentRead(d, meta)
}

func resourceAliCloudPaiWorkspaceExperimentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	ExperimentId := d.Id()
	action := fmt.Sprintf("/api/v1/experiments/%s", ExperimentId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["ExperimentId"] = d.Id()

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

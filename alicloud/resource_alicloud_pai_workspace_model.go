// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudPaiWorkspaceModel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspaceModelCreate,
		Read:   resourceAliCloudPaiWorkspaceModelRead,
		Update: resourceAliCloudPaiWorkspaceModelUpdate,
		Delete: resourceAliCloudPaiWorkspaceModelDelete,
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
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extra_info": {
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
			"model_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"model_doc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"model_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"model_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudPaiWorkspaceModelCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/models")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ModelName"] = d.Get("model_name")
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

	if v, ok := d.GetOk("model_description"); ok {
		request["ModelDescription"] = v
	}
	if v, ok := d.GetOk("workspace_id"); ok {
		request["WorkspaceId"] = v
	}
	if v, ok := d.GetOk("accessibility"); ok {
		request["Accessibility"] = v
	}
	if v, ok := d.GetOk("origin"); ok {
		request["Origin"] = v
	}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	if v, ok := d.GetOk("task"); ok {
		request["Task"] = v
	}
	if v, ok := d.GetOk("model_doc"); ok {
		request["ModelDoc"] = v
	}
	if v, ok := d.GetOkExists("order_number"); ok {
		request["OrderNumber"] = v
	}
	if v, ok := d.GetOk("model_type"); ok {
		request["ModelType"] = v
	}
	if v, ok := d.GetOk("extra_info"); ok {
		request["ExtraInfo"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_workspace_model", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.ModelId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudPaiWorkspaceModelRead(d, meta)
}

func resourceAliCloudPaiWorkspaceModelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspaceModel(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_model DescribePaiWorkspaceModel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accessibility", objectRaw["Accessibility"])
	d.Set("domain", objectRaw["Domain"])
	d.Set("extra_info", objectRaw["ExtraInfo"])
	d.Set("model_description", objectRaw["ModelDescription"])
	d.Set("model_doc", objectRaw["ModelDoc"])
	d.Set("model_name", objectRaw["ModelName"])
	d.Set("model_type", objectRaw["ModelType"])
	d.Set("order_number", objectRaw["OrderNumber"])
	d.Set("origin", objectRaw["Origin"])
	d.Set("task", objectRaw["Task"])
	d.Set("workspace_id", objectRaw["WorkspaceId"])

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

	return nil
}

func resourceAliCloudPaiWorkspaceModelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	ModelId := d.Id()
	action := fmt.Sprintf("/api/v1/models/%s", ModelId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["ModelId"] = d.Id()

	if d.HasChange("model_name") {
		update = true
	}
	request["ModelName"] = d.Get("model_name")
	if d.HasChange("model_description") {
		update = true
	}
	if v, ok := d.GetOk("model_description"); ok {
		request["ModelDescription"] = v
	}
	if d.HasChange("accessibility") {
		update = true
	}
	if v, ok := d.GetOk("accessibility"); ok {
		request["Accessibility"] = v
	}
	if d.HasChange("origin") {
		update = true
	}
	if v, ok := d.GetOk("origin"); ok {
		request["Origin"] = v
	}
	if d.HasChange("domain") {
		update = true
	}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	if d.HasChange("task") {
		update = true
	}
	if v, ok := d.GetOk("task"); ok {
		request["Task"] = v
	}
	if d.HasChange("model_doc") {
		update = true
	}
	if v, ok := d.GetOk("model_doc"); ok {
		request["ModelDoc"] = v
	}
	if d.HasChange("order_number") {
		update = true
	}
	if v, ok := d.GetOkExists("order_number"); ok {
		request["OrderNumber"] = v
	}
	if d.HasChange("model_type") {
		update = true
	}
	if v, ok := d.GetOk("model_type"); ok {
		request["ModelType"] = v
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
			ModelId := d.Id()
			action := fmt.Sprintf("/api/v1/models/%s/labels", ModelId)
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			request["ModelId"] = d.Id()

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
			ModelId := d.Id()
			action := fmt.Sprintf("/api/v1/models/%s/labels", ModelId)
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			request["ModelId"] = d.Id()

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
	return resourceAliCloudPaiWorkspaceModelRead(d, meta)
}

func resourceAliCloudPaiWorkspaceModelDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	ModelId := d.Id()
	action := fmt.Sprintf("/api/v1/models/%s", ModelId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["ModelId"] = d.Id()

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

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

func resourceAliCloudPaiWorkspaceDataset() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspaceDatasetCreate,
		Read:   resourceAliCloudPaiWorkspaceDatasetRead,
		Update: resourceAliCloudPaiWorkspaceDatasetUpdate,
		Delete: resourceAliCloudPaiWorkspaceDatasetDelete,
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
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"OSS", "NAS"}, false),
			},
			"data_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"COMMON", "PIC", "TEXT"}, false),
			},
			"dataset_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"USER", "ITAG"}, false),
			},
			"uri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^\\d+$"), "The ID of the dataset owner."),
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudPaiWorkspaceDatasetCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/datasets")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Property"] = d.Get("property")
	request["DataSourceType"] = d.Get("data_source_type")
	request["Uri"] = d.Get("uri")
	if v, ok := d.GetOk("data_type"); ok {
		request["DataType"] = v
	}
	if v, ok := d.GetOk("source_type"); ok {
		request["SourceType"] = v
	}
	if v, ok := d.GetOk("source_id"); ok {
		request["SourceId"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["WorkspaceId"] = d.Get("workspace_id")
	if v, ok := d.GetOk("options"); ok {
		request["Options"] = v
	}
	if v, ok := d.GetOk("accessibility"); ok {
		request["Accessibility"] = v
	}
	request["Name"] = d.Get("dataset_name")
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

	if v, ok := d.GetOk("user_id"); ok {
		request["UserId"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_workspace_dataset", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DatasetId"]))

	return resourceAliCloudPaiWorkspaceDatasetRead(d, meta)
}

func resourceAliCloudPaiWorkspaceDatasetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspaceDataset(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_dataset DescribePaiWorkspaceDataset Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Accessibility"] != nil {
		d.Set("accessibility", objectRaw["Accessibility"])
	}
	if objectRaw["GmtCreateTime"] != nil {
		d.Set("create_time", objectRaw["GmtCreateTime"])
	}
	if objectRaw["DataSourceType"] != nil {
		d.Set("data_source_type", objectRaw["DataSourceType"])
	}
	if objectRaw["DataType"] != nil {
		d.Set("data_type", objectRaw["DataType"])
	}
	if objectRaw["Name"] != nil {
		d.Set("dataset_name", objectRaw["Name"])
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
	if objectRaw["UserId"] != nil {
		d.Set("user_id", objectRaw["UserId"])
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

func resourceAliCloudPaiWorkspaceDatasetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	if d.HasChange("accessibility") {
		paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}
		object, err := paiWorkspaceServiceV2.DescribePaiWorkspaceDataset(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("accessibility").(string)
		if object["Accessibility"].(string) != target {
			if target == "PUBLIC" {
				DatasetId := d.Id()
				action := fmt.Sprintf("/api/v1/datasets/%s/publish", DatasetId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["DatasetId"] = d.Id()

				body = request
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
		}
	}

	DatasetId := d.Id()
	action := fmt.Sprintf("/api/v1/datasets/%s", DatasetId)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["DatasetId"] = d.Id()

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["Description"] = v
	}
	if d.HasChange("options") {
		update = true
	}
	if v, ok := d.GetOk("options"); ok || d.HasChange("options") {
		request["Options"] = v
	}
	if d.HasChange("dataset_name") {
		update = true
	}
	request["Name"] = d.Get("dataset_name")
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
			DatasetId := d.Id()
			action := fmt.Sprintf("/api/v1/datasets/%s/labels", DatasetId)
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			request["DatasetId"] = d.Id()
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
			DatasetId := d.Id()
			action := fmt.Sprintf("/api/v1/datasets/%s/labels", DatasetId)
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			request["DatasetId"] = d.Id()

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
	return resourceAliCloudPaiWorkspaceDatasetRead(d, meta)
}

func resourceAliCloudPaiWorkspaceDatasetDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	DatasetId := d.Id()
	action := fmt.Sprintf("/api/v1/datasets/%s", DatasetId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["DatasetId"] = d.Id()

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

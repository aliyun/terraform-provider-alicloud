package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudOosExecution() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosExecutionCreate,
		Read:   resourceAlicloudOosExecutionRead,
		Delete: resourceAlicloudOosExecutionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"counters": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"end_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"executed_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_parent": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"loop_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Automatic", "Debug"}, false),
				Default:      "Automatic",
			},
			"outputs": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "{}",
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"parent_execution_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ram_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"safety_check": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"start_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tagsSchemaForceNew(),
			"update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// oosExecutionStartExecutionRequest builds the StartExecution payload for the
// configured execution. It is kept separate from Create so the exact payload
// handed to the transport (notably the target ResourceGroupId) can be verified
// without cloud access.
func oosExecutionStartExecutionRequest(d *schema.ResourceData, regionId string) (map[string]interface{}, error) {
	request := make(map[string]interface{})
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("loop_mode"); ok {
		request["LoopMode"] = v
	}

	if v, ok := d.GetOk("mode"); ok {
		request["Mode"] = v
	}

	if v, ok := d.GetOk("parameters"); ok {
		request["Parameters"] = v
	}

	if v, ok := d.GetOk("parent_execution_id"); ok {
		request["ParentExecutionId"] = v
	}

	request["RegionId"] = regionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("safety_check"); ok {
		request["SafetyCheck"] = v
	}

	if v, ok := d.GetOk("template_content"); ok {
		request["TemplateContent"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return nil, WrapError(err)
		}
		request["Tags"] = respJson
	}

	request["TemplateName"] = d.Get("template_name")
	if v, ok := d.GetOk("template_version"); ok {
		request["TemplateVersion"] = v
	}

	return request, nil
}

func resourceAlicloudOosExecutionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "StartExecution"
	request, err := oosExecutionStartExecutionRequest(d, client.RegionId)
	if err != nil {
		return WrapError(err)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_execution", action, AlibabaCloudSdkGoERROR)
	}
	responseExecution := response["Execution"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseExecution["ExecutionId"]))
	oosService := OosService{client}
	stateConf := BuildStateConf([]string{}, []string{"Success", "Failed", "Cancelled"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, oosService.OosExecutionStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudOosExecutionRead(d, meta)
}
func resourceAlicloudOosExecutionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosExecution(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_execution oosService.DescribeOosExecution Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	setOosExecutionState(d, object)
	return nil
}

// setOosExecutionState maps a ListExecutions execution object into state. It is
// kept separate from Read so the mapping of the response to resource_group_id
// and the derived resource_ids is verifiable without cloud access.
func setOosExecutionState(d *schema.ResourceData, object map[string]interface{}) {
	d.Set("create_date", object["CreateDate"])
	if v, ok := object["Counters"].(map[string]interface{}); ok {
		counters, _ := convertMaptoJsonString(v)
		d.Set("counters", counters)
	} else {
		d.Set("counters", object["Counters"])
	}
	d.Set("description", object["Description"])
	d.Set("end_date", object["EndDate"])
	d.Set("executed_by", object["ExecutedBy"])
	d.Set("is_parent", object["IsParent"])
	d.Set("loop_mode", object["LoopMode"])
	d.Set("mode", object["Mode"])
	d.Set("outputs", object["Outputs"])
	if v, ok := object["Parameters"].(map[string]interface{}); ok {
		params, _ := convertMaptoJsonString(v)
		d.Set("parameters", params)
	} else {
		d.Set("parameters", object["Parameters"])
	}
	d.Set("parent_execution_id", object["ParentExecutionId"])
	d.Set("ram_role", object["RamRole"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("safety_check", object["SafetyCheck"])
	d.Set("start_date", object["StartDate"])
	d.Set("status", object["Status"])
	d.Set("status_message", object["StatusMessage"])
	d.Set("template_content", object["TemplateContent"])
	d.Set("template_id", object["TemplateId"])
	d.Set("template_name", object["TemplateName"])
	d.Set("template_version", object["TemplateVersion"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v))
	}
	// $.ResourceIds is exposed as a computed attribute. ListExecutions returns
	// the execution's tags but not its ResourceIds, and ListTagResources only
	// echoes the queried execution id once per bound tag, so the value is the
	// execution itself while it carries tags and empty otherwise.
	resourceIds := make([]string, 0)
	if v, ok := object["Tags"].(map[string]interface{}); ok && len(v) > 0 {
		resourceIds = append(resourceIds, d.Id())
	}
	d.Set("resource_ids", resourceIds)
	d.Set("update_date", object["UpdateDate"])
}

func resourceAlicloudOosExecutionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteExecutions"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ExecutionIds": convertListToJsonString(convertListStringToListInterface([]string{d.Id()})),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

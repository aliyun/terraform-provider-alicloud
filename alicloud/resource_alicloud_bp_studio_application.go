package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudBpStudioApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBpStudioApplicationCreate,
		Read:   resourceAlicloudBpStudioApplicationRead,
		Delete: resourceAlicloudBpStudioApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(120 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"area_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"node_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"node_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"configuration": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"variables": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudBpStudioApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bpStudioService := BpStudioService{client}

	//1. CreateApplication
	action := "CreateApplication"
	var response map[string]interface{}
	request := make(map[string]interface{})
	var err error
	request["Name"] = d.Get("application_name")
	request["TemplateId"] = d.Get("template_id")
	request["ClientToken"] = buildClientToken("CreateApplication")

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("area_id"); ok {
		request["AreaId"] = v
	}

	if v, ok := d.GetOk("instances"); ok {
		instancesMaps := make([]map[string]interface{}, 0)
		for _, instances := range v.([]interface{}) {
			instancesMap := map[string]interface{}{}
			instancesArg := instances.(map[string]interface{})

			if id, ok := instancesArg["id"].(string); ok && id != "" {
				instancesMap["Id"] = id
			}

			if nodeName, ok := instancesArg["node_name"].(string); ok && nodeName != "" {
				instancesMap["NodeName"] = nodeName
			}

			if nodeType, ok := instancesArg["node_type"].(string); ok && nodeType != "" {
				instancesMap["NodeType"] = nodeType
			}

			instancesMaps = append(instancesMaps, instancesMap)
		}

		instancesMapsJson, err := convertListMapToJsonString(instancesMaps)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_bp_studio_application", action, AlibabaCloudSdkGoERROR)
		}
		request["Instances"] = instancesMapsJson
	}

	if v, ok := d.GetOk("configuration"); ok {
		configuration := v.(map[string]interface{})
		configurationJson, err := convertMaptoJsonString(configuration)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_bp_studio_application", action, AlibabaCloudSdkGoERROR)
		}
		request["Configuration"] = configurationJson
	}

	if v, ok := d.GetOk("variables"); ok {
		variables := v.(map[string]interface{})
		variablesJson, err := convertMaptoJsonString(variables)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_bp_studio_application", action, AlibabaCloudSdkGoERROR)
		}
		request["Variables"] = variablesJson
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("BPStudio", "2021-09-31", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bp_studio_application", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Data"]))

	stateConf := BuildStateConf([]string{}, []string{"Modified"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, bpStudioService.BpStudioApplicationStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	//2. ValidateApplication
	validateApplicationAction := "ValidateApplication"
	validateApplicationResponse := make(map[string]interface{})
	deployReq := make(map[string]interface{})
	deployReq["ApplicationId"] = response["Data"]
	if v, ok := d.GetOk("resource_group_id"); ok {
		deployReq["ResourceGroupId"] = v
	}

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		validateApplicationResponse, err = client.RpcPost("BPStudio", "2021-09-31", validateApplicationAction, nil, deployReq, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(validateApplicationAction, validateApplicationResponse, deployReq)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), validateApplicationAction, AlibabaCloudSdkGoERROR)
	}

	validateApplicationStateConf := BuildStateConf([]string{}, []string{"Verified_Success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, bpStudioService.BpStudioApplicationStateRefreshFunc(d.Id(), []string{}))
	if _, err := validateApplicationStateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	//3. ValuateApplication
	valuateApplicationAction := "ValuateApplication"
	valuateApplicationResponse := make(map[string]interface{})

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		valuateApplicationResponse, err = client.RpcPost("BPStudio", "2021-09-31", valuateApplicationAction, nil, deployReq, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(valuateApplicationAction, valuateApplicationResponse, deployReq)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), valuateApplicationAction, AlibabaCloudSdkGoERROR)
	}

	valuateApplicationStateConf := BuildStateConf([]string{}, []string{"Valuating_Success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, bpStudioService.BpStudioApplicationStateRefreshFunc(d.Id(), []string{}))
	if _, err := valuateApplicationStateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	//4. DeployApplication
	deployApplicationAction := "DeployApplication"
	deployApplicationResponse := make(map[string]interface{})

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		deployApplicationResponse, err = client.RpcPost("BPStudio", "2021-09-31", deployApplicationAction, nil, deployReq, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(deployApplicationAction, deployApplicationResponse, deployReq)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), deployApplicationAction, AlibabaCloudSdkGoERROR)
	}

	deployApplicationStateConf := BuildStateConf([]string{}, []string{"Deployed_Success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, bpStudioService.BpStudioApplicationStateRefreshFunc(d.Id(), []string{}))
	if _, err := deployApplicationStateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudBpStudioApplicationRead(d, meta)
}

func resourceAlicloudBpStudioApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bpStudioService := BpStudioService{client}

	object, err := bpStudioService.DescribeBpStudioApplication(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("application_name", object["Name"])
	d.Set("template_id", object["TemplateId"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAlicloudBpStudioApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bpStudioService := BpStudioService{client}
	var err error

	request := make(map[string]interface{})
	request["ApplicationId"] = d.Id()

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	object, err := bpStudioService.DescribeBpStudioApplication(d.Id())
	if err != nil {
		return WrapError(err)
	}

	//1. ReleaseApplication
	wait := incrementalWait(3*time.Second, 5*time.Second)

	if fmt.Sprint(object["Status"]) == "Deployed_Failure" || fmt.Sprint(object["Status"]) == "PartiallyDeployedSuccess" || fmt.Sprint(object["Status"]) == "Deployed_Success" || fmt.Sprint(object["Status"]) == "Destroyed_Failure" || fmt.Sprint(object["Status"]) == "PartiallyDestroyedSuccess" {
		releaseApplicationAction := "ReleaseApplication"
		releaseApplicationResponse := make(map[string]interface{})

		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			releaseApplicationResponse, err = client.RpcPost("BPStudio", "2021-09-31", releaseApplicationAction, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(releaseApplicationAction, releaseApplicationResponse, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), releaseApplicationAction, AlibabaCloudSdkGoERROR)
		}

		releaseApplicationStateConf := BuildStateConf([]string{}, []string{"Destroyed_Success"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, bpStudioService.BpStudioApplicationStateRefreshFunc(d.Id(), []string{}))
		if _, err := releaseApplicationStateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	//2. DeleteApplication
	deleteApplicationAction := "DeleteApplication"
	deleteApplicationResponse := make(map[string]interface{})

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		deleteApplicationResponse, err = client.RpcPost("BPStudio", "2021-09-31", deleteApplicationAction, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(deleteApplicationAction, deleteApplicationResponse, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteApplicationAction, AlibabaCloudSdkGoERROR)
	}

	deleteApplicationStateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, bpStudioService.BpStudioApplicationStateRefreshFunc(d.Id(), []string{}))
	if _, err := deleteApplicationStateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

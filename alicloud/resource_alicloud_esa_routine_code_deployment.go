package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudEsaRoutineCodeDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaRoutineCodeDeploymentCreate,
		Read:   resourceAliCloudEsaRoutineCodeDeploymentRead,
		Delete: resourceAliCloudEsaRoutineCodeDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"routine_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"env": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"staging", "production"}, false),
			},
			"strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "percentage",
				ValidateFunc: StringInSlice([]string{"percentage"}, false),
			},
			"code_versions": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code_version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"percentage": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
					},
				},
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaRoutineCodeDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if err := esaRoutineCreateCodeDeployment(d, client, d.Timeout(schema.TimeoutCreate)); err != nil {
		return WrapError(err)
	}

	d.SetId(fmt.Sprintf("%v:%v", d.Get("routine_name"), d.Get("env")))

	return resourceAliCloudEsaRoutineCodeDeploymentRead(d, meta)
}

func resourceAliCloudEsaRoutineCodeDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaRoutineCodeDeployment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_routine_code_deployment DescribeEsaRoutineCodeDeployment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("routine_name", parts[0])
	d.Set("env", parts[1])
	d.Set("deployment_id", objectRaw["DeployId"])
	if v, ok := objectRaw["Strategy"]; ok && fmt.Sprint(v) != "" && fmt.Sprint(v) != "<nil>" {
		d.Set("strategy", v)
	}

	codeVersions := make([]map[string]interface{}, 0)
	if versionsRaw, ok := objectRaw["CodeVersions"].([]interface{}); ok {
		for _, item := range versionsRaw {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			codeVersions = append(codeVersions, map[string]interface{}{
				"code_version": fmt.Sprint(itemMap["CodeVersion"]),
				"percentage":   formatInt(itemMap["Percentage"]),
			})
		}
	}
	if err := d.Set("code_versions", codeVersions); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliCloudEsaRoutineCodeDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	// ESA has no DeleteRoutineCodeDeployment API. A code deployment is an immutable
	// historical deployment record on the routine; the resource is removed from
	// Terraform state only. Deleting the parent alicloud_esa_routine removes all of
	// its deployment records on the server side.
	log.Printf("[WARN] Cannot destroy resource alicloud_esa_routine_code_deployment %s. Terraform will remove this resource from the state file, however resources may remain.", d.Id())
	return nil
}

// esaRoutineCreateCodeDeployment issues CreateRoutineCodeDeployment and records the
// resulting deployment id. It is shared by Create and Update since a new percentage
// configuration is applied by creating a new deployment.
func esaRoutineCreateCodeDeployment(d *schema.ResourceData, client *connectivity.AliyunClient, timeout time.Duration) error {
	action := "CreateRoutineCodeDeployment"
	request := map[string]interface{}{
		"Name": d.Get("routine_name"),
		"Env":  d.Get("env"),
	}
	if v, ok := d.GetOk("strategy"); ok {
		request["Strategy"] = v
	}

	codeVersionsArray := make([]interface{}, 0)
	for _, dataLoop := range d.Get("code_versions").([]interface{}) {
		dataLoopTmp := dataLoop.(map[string]interface{})
		codeVersionsArray = append(codeVersionsArray, map[string]interface{}{
			"CodeVersion": dataLoopTmp["code_version"],
			"Percentage":  dataLoopTmp["percentage"],
		})
	}
	codeVersionsJson, err := json.Marshal(codeVersionsArray)
	if err != nil {
		return WrapError(err)
	}
	request["CodeVersions"] = string(codeVersionsJson)

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(timeout, func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "LockFailed"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_routine_code_deployment", action, AlibabaCloudSdkGoERROR)
	}

	if v, ok := response["DeploymentId"]; ok {
		d.Set("deployment_id", v)
	}

	return nil
}

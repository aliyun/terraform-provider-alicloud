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
				MinItems: 1,
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

	action := "CreateRoutineCodeDeployment"
	request := map[string]interface{}{
		"Name":     d.Get("routine_name"),
		"Env":      d.Get("env"),
		"Strategy": d.Get("strategy"),
	}

	codeVersionsArray := make([]map[string]interface{}, 0)
	for _, v := range d.Get("code_versions").([]interface{}) {
		item := v.(map[string]interface{})
		codeVersionsArray = append(codeVersionsArray, map[string]interface{}{
			"CodeVersion": item["code_version"],
			"Percentage":  item["percentage"],
		})
	}
	codeVersionsJson, err := json.Marshal(codeVersionsArray)
	if err != nil {
		return WrapError(err)
	}
	request["CodeVersions"] = string(codeVersionsJson)

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		d.Set("deployment_id", fmt.Sprint(v))
	}

	d.SetId(fmt.Sprintf("%s:%s", d.Get("routine_name"), d.Get("env")))

	return resourceAliCloudEsaRoutineCodeDeploymentRead(d, meta)
}

func resourceAliCloudEsaRoutineCodeDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	routineName := parts[0]
	env := parts[1]

	object, err := esaServiceV2.DescribeEsaRoutineCodeDeployment(routineName, env)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_routine_code_deployment DescribeEsaRoutineCodeDeployment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("routine_name", routineName)
	d.Set("env", env)

	if v, ok := object["Strategy"]; ok && v != nil {
		d.Set("strategy", v)
	}
	if v, ok := object["DeployId"]; ok && v != nil {
		d.Set("deployment_id", fmt.Sprint(v))
	}

	codeVersions := make([]map[string]interface{}, 0)
	if raw, ok := object["CodeVersions"].([]interface{}); ok {
		for _, item := range raw {
			cv, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			codeVersions = append(codeVersions, map[string]interface{}{
				"code_version": fmt.Sprint(cv["CodeVersion"]),
				"percentage":   formatInt(cv["Percentage"]),
			})
		}
	}
	if err := d.Set("code_versions", codeVersions); err != nil {
		return WrapError(err)
	}

	return nil
}

// resourceAliCloudEsaRoutineCodeDeploymentDelete is a no-op: the ESA OpenAPI
// exposes no operation to remove a code deployment. A deployment is a rollout
// event bound to an environment; the currently active version can only be
// replaced by a new deployment, never withdrawn. The resource is therefore
// removed from state without calling the backend.
func resourceAliCloudEsaRoutineCodeDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource alicloud_esa_routine_code_deployment. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

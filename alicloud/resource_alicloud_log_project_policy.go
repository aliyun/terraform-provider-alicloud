package alicloud

import (
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudLogProjectPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogProjectPolicyCreate,
		Read:   resourceAlicloudLogProjectPolicyRead,
		Update: resourceAlicloudLogProjectPolicyUpdate,
		Delete: resourceAlicloudLogProjectPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
		},
	}
}

func resourceAlicloudLogProjectPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *sls.Client

	policy := d.Get("policy").(string)
	projectName := d.Get("project").(string)

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			if policy == "" {
				return nil, slsClient.DeleteProjectPolicy(projectName)
			}
			return nil, slsClient.UpdateProjectPolicy(projectName, policy)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateProjectPolicy", policy, requestInfo, map[string]interface{}{
			"policy": policy,
		})
		d.SetId(projectName)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_project_policy", "CreateProjectPolicy", AliyunLogGoSdkERROR)
	}
	return resourceAlicloudLogProjectPolicyRead(d, meta)
}

func resourceAlicloudLogProjectPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	policy, err := logService.DescribeLogProjectPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("project", d.Id())
	d.Set("policy", policy)
	return nil
}

func resourceAlicloudLogProjectPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	projectName := d.Id()
	client := meta.(*connectivity.AliyunClient)
	policy := d.Get("policy").(string)
	_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		if policy == "" {
			return nil, slsClient.DeleteProjectPolicy(projectName)
		}
		return nil, slsClient.UpdateProjectPolicy(projectName, policy)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateProjectPolicy", AliyunLogGoSdkERROR)
	}
	return resourceAlicloudLogProjectPolicyRead(d, meta)
}

func resourceAlicloudLogProjectPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *sls.Client
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteProjectPolicy(d.Id())
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout, "RequestTimeout"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteProjectPolicy", raw, requestInfo, map[string]interface{}{
			"project": d.Id(),
		})
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ProjectNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteProjectPolicy", AliyunLogGoSdkERROR)
	}
	return nil
}

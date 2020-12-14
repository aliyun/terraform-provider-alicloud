package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudResourceManagerPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerPolicyCreate,
		Read:   resourceAlicloudResourceManagerPolicyRead,
		Update: resourceAlicloudResourceManagerPolicyUpdate,
		Delete: resourceAlicloudResourceManagerPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"default_version": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'default_version' has been deprecated from provider version 1.90.0",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_document": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateCreatePolicyRequest()
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	request.PolicyDocument = d.Get("policy_document").(string)
	request.PolicyName = d.Get("policy_name").(string)
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.CreatePolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_policy", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*resourcemanager.CreatePolicyResponse)
	d.SetId(fmt.Sprintf("%v", response.Policy.PolicyName))

	return resourceAlicloudResourceManagerPolicyRead(d, meta)
}
func resourceAlicloudResourceManagerPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_policy resourcemanagerService.DescribeResourceManagerPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("policy_name", d.Id())
	d.Set("default_version", object.DefaultVersion)
	d.Set("description", object.Description)
	d.Set("policy_type", object.PolicyType)

	getPolicyVersionObject, err := resourcemanagerService.GetPolicyVersion(d.Id(), d)
	if err != nil {
		return WrapError(err)
	}
	d.Set("policy_document", getPolicyVersionObject.PolicyDocument)
	return nil
}
func resourceAlicloudResourceManagerPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("default_version") {
		request := resourcemanager.CreateSetDefaultPolicyVersionRequest()
		request.PolicyName = d.Id()
		request.VersionId = d.Get("default_version").(string)
		raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
			return resourcemanagerClient.SetDefaultPolicyVersion(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudResourceManagerPolicyRead(d, meta)
}
func resourceAlicloudResourceManagerPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := resourcemanager.CreateDeletePolicyRequest()
	request.PolicyName = d.Id()
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.DeletePolicy(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

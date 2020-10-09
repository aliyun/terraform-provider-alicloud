package alicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudResourceManagerPolicyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerPolicyVersionCreate,
		Read:   resourceAlicloudResourceManagerPolicyVersionRead,
		Update: resourceAlicloudResourceManagerPolicyVersionUpdate,
		Delete: resourceAlicloudResourceManagerPolicyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"is_default_version": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    false,
				Deprecated: "Field 'is_default_version' has been deprecated from provider version 1.90.0",
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
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Field 'create_date' has been removed from provider version 1.100.0.",
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Removed:  "Field 'version_id' has been removed from provider version 1.100.0.",
			},
		},
	}
}

func resourceAlicloudResourceManagerPolicyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateCreatePolicyVersionRequest()
	if v, ok := d.GetOkExists("is_default_version"); ok {
		request.SetAsDefault = requests.NewBoolean(v.(bool))
	}

	request.PolicyDocument = d.Get("policy_document").(string)
	request.PolicyName = d.Get("policy_name").(string)
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.CreatePolicyVersion(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_policy_version", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*resourcemanager.CreatePolicyVersionResponse)
	d.SetId(fmt.Sprintf("%v:%v", request.PolicyName, response.PolicyVersion.VersionId))

	return resourceAlicloudResourceManagerPolicyVersionRead(d, meta)
}
func resourceAlicloudResourceManagerPolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerPolicyVersion(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_policy_version resourcemanagerService.DescribeResourceManagerPolicyVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("policy_name", parts[0])
	d.Set("is_default_version", object.IsDefaultVersion)
	d.Set("policy_document", object.PolicyDocument)
	return nil
}
func resourceAlicloudResourceManagerPolicyVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	if d.HasChange("is_default_version") {
		object, err := resourcemanagerService.DescribeResourceManagerPolicyVersion(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := strconv.FormatBool(d.Get("is_default_version").(bool))
		if strconv.FormatBool(object.IsDefaultVersion) != target {
			if target == "true" {
				request := resourcemanager.CreateSetDefaultPolicyVersionRequest()
				request.PolicyName = parts[0]
				request.VersionId = parts[1]
				raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
					return resourcemanagerClient.SetDefaultPolicyVersion(request)
				})
				addDebug(request.GetActionName(), raw)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
		}
	}
	return resourceAlicloudResourceManagerPolicyVersionRead(d, meta)
}
func resourceAlicloudResourceManagerPolicyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := resourcemanager.CreateDeletePolicyVersionRequest()
	request.PolicyName = parts[0]
	request.VersionId = parts[1]
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.DeletePolicyVersion(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

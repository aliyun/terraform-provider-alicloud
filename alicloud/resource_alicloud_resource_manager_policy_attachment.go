package alicloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudResourceManagerPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerPolicyAttachmentCreate,
		Read:   resourceAlicloudResourceManagerPolicyAttachmentRead,
		Delete: resourceAlicloudResourceManagerPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Custom", "System"}, false),
			},
			"principal_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"principal_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IMSGroup", "IMSUser", "ServiceRole"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateAttachPolicyRequest()
	request.PolicyName = d.Get("policy_name").(string)
	request.PolicyType = d.Get("policy_type").(string)
	request.PrincipalName = d.Get("principal_name").(string)
	request.PrincipalType = d.Get("principal_type").(string)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.AttachPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_policy_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	d.SetId(fmt.Sprintf("%v:%v:%v:%v:%v", request.PolicyName, request.PolicyType, request.PrincipalName, request.PrincipalType, request.ResourceGroupId))

	return resourceAlicloudResourceManagerPolicyAttachmentRead(d, meta)
}
func resourceAlicloudResourceManagerPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	if len(strings.Split(d.Id(), ":")) != 5 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), d.Get("resource_group_id").(string)))
	}
	_, err := resourcemanagerService.DescribeResourceManagerPolicyAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_policy_attachment resourcemanagerService.DescribeResourceManagerPolicyAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 5)
	if err != nil {
		return WrapError(err)
	}
	d.Set("policy_name", parts[0])
	d.Set("policy_type", parts[1])
	d.Set("principal_name", parts[2])
	d.Set("principal_type", parts[3])
	d.Set("resource_group_id", parts[4])
	return nil
}
func resourceAlicloudResourceManagerPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if len(strings.Split(d.Id(), ":")) != 5 {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), d.Get("resource_group_id").(string)))
	}
	parts, err := ParseResourceId(d.Id(), 5)
	if err != nil {
		return WrapError(err)
	}
	request := resourcemanager.CreateDetachPolicyRequest()
	request.PolicyName = parts[0]
	request.PolicyType = parts[1]
	request.PrincipalName = parts[2]
	request.PrincipalType = parts[3]
	request.ResourceGroupId = parts[4]
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.DetachPolicy(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExists.ResourceGroup"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

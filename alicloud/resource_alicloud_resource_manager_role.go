package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudResourceManagerRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerRoleCreate,
		Read:   resourceAlicloudResourceManagerRoleRead,
		Update: resourceAlicloudResourceManagerRoleUpdate,
		Delete: resourceAlicloudResourceManagerRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assume_role_policy_document": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
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
			"max_session_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateCreateRoleRequest()
	request.AssumeRolePolicyDocument = d.Get("assume_role_policy_document").(string)
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("max_session_duration"); ok {
		request.MaxSessionDuration = requests.NewInteger(v.(int))
	}
	request.RoleName = d.Get("role_name").(string)
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.CreateRole(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_role", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*resourcemanager.CreateRoleResponse)
	d.SetId(fmt.Sprintf("%v", response.Role.RoleName))

	return resourceAlicloudResourceManagerRoleRead(d, meta)
}
func resourceAlicloudResourceManagerRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerRole(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("role_name", d.Id())
	d.Set("arn", object.Arn)
	d.Set("assume_role_policy_document", object.AssumeRolePolicyDocument)
	d.Set("create_date", object.CreateDate)
	d.Set("description", object.Description)
	d.Set("max_session_duration", object.MaxSessionDuration)
	d.Set("role_id", object.RoleId)
	d.Set("update_date", object.UpdateDate)
	return nil
}
func resourceAlicloudResourceManagerRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := resourcemanager.CreateUpdateRoleRequest()
	request.RoleName = d.Id()
	if d.HasChange("assume_role_policy_document") {
		update = true
	}
	request.NewAssumeRolePolicyDocument = d.Get("assume_role_policy_document").(string)
	if d.HasChange("max_session_duration") {
		update = true
		request.NewMaxSessionDuration = requests.NewInteger(d.Get("max_session_duration").(int))
	}
	if update {
		raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
			return resourcemanagerClient.UpdateRole(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudResourceManagerRoleRead(d, meta)
}
func resourceAlicloudResourceManagerRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := resourcemanager.CreateDeleteRoleRequest()
	request.RoleName = d.Id()
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.DeleteRole(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

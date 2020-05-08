package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudResourceManagerHandshake() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerHandshakeCreate,
		Read:   resourceAlicloudResourceManagerHandshakeRead,
		Delete: resourceAlicloudResourceManagerHandshakeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_directory_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_entity": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Account", "Email"}, false),
			},
		},
	}
}

func resourceAlicloudResourceManagerHandshakeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateInviteAccountToResourceDirectoryRequest()
	if v, ok := d.GetOk("note"); ok {
		request.Note = v.(string)
	}
	request.TargetEntity = d.Get("target_entity").(string)
	request.TargetType = d.Get("target_type").(string)
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.InviteAccountToResourceDirectory(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_handshake", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*resourcemanager.InviteAccountToResourceDirectoryResponse)
	d.SetId(fmt.Sprintf("%v", response.Handshake.HandshakeId))

	return resourceAlicloudResourceManagerHandshakeRead(d, meta)
}
func resourceAlicloudResourceManagerHandshakeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerHandshake(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("expire_time", object.ExpireTime)
	d.Set("master_account_id", object.MasterAccountId)
	d.Set("master_account_name", object.MasterAccountName)
	d.Set("modify_time", object.ModifyTime)
	d.Set("note", object.Note)
	d.Set("resource_directory_id", object.ResourceDirectoryId)
	d.Set("status", object.Status)
	d.Set("target_entity", object.TargetEntity)
	d.Set("target_type", object.TargetType)
	return nil
}
func resourceAlicloudResourceManagerHandshakeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := resourcemanager.CreateCancelHandshakeRequest()
	request.HandshakeId = d.Id()
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.CancelHandshake(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Handshake", "HandshakeStatusMismatch"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudResourceManagerAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerAccountCreate,
		Read:   resourceAlicloudResourceManagerAccountRead,
		Update: resourceAlicloudResourceManagerAccountUpdate,
		Delete: resourceAlicloudResourceManagerAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"folder_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"join_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"join_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payer_account_id": {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateCreateResourceAccountRequest()
	request.DisplayName = d.Get("display_name").(string)
	if v, ok := d.GetOk("folder_id"); ok {
		request.ParentFolderId = v.(string)
	}
	if v, ok := d.GetOk("payer_account_id"); ok {
		request.PayerAccountId = v.(string)
	}
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.CreateResourceAccount(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_account", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*resourcemanager.CreateResourceAccountResponse)
	d.SetId(fmt.Sprintf("%v", response.Account.AccountId))

	return resourceAlicloudResourceManagerAccountRead(d, meta)
}
func resourceAlicloudResourceManagerAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("display_name", object.DisplayName)
	d.Set("folder_id", object.FolderId)
	d.Set("join_method", object.JoinMethod)
	d.Set("join_time", object.JoinTime)
	d.Set("modify_time", object.ModifyTime)
	d.Set("resource_directory_id", object.ResourceDirectoryId)
	d.Set("status", object.Status)
	d.Set("type", object.Type)
	return nil
}
func resourceAlicloudResourceManagerAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	if d.HasChange("folder_id") {
		request := resourcemanager.CreateMoveAccountRequest()
		request.AccountId = d.Id()
		request.DestinationFolderId = d.Get("folder_id").(string)
		raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
			return resourcemanagerClient.MoveAccount(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("folder_id")
	}
	update := false
	request := resourcemanager.CreateUpdateAccountRequest()
	request.AccountId = d.Id()
	if d.HasChange("display_name") {
		update = true
	}
	request.NewDisplayName = d.Get("display_name").(string)
	if update {
		raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
			return resourcemanagerClient.UpdateAccount(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("display_name")
	}
	d.Partial(false)
	return resourceAlicloudResourceManagerAccountRead(d, meta)
}
func resourceAlicloudResourceManagerAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudResourceManagerAccount. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

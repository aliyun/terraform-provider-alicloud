package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudResourceManagerFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerFolderCreate,
		Read:   resourceAlicloudResourceManagerFolderRead,
		Update: resourceAlicloudResourceManagerFolderUpdate,
		Delete: resourceAlicloudResourceManagerFolderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"folder_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_folder_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerFolderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateCreateFolderRequest()
	request.FolderName = d.Get("folder_name").(string)
	if v, ok := d.GetOk("parent_folder_id"); ok {
		request.ParentFolderId = v.(string)
	}
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.CreateFolder(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_folder", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*resourcemanager.CreateFolderResponse)
	d.SetId(fmt.Sprintf("%v", response.Folder.FolderId))

	return resourceAlicloudResourceManagerFolderRead(d, meta)
}
func resourceAlicloudResourceManagerFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerFolder(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("folder_name", object.FolderName)
	d.Set("parent_folder_id", object.ParentFolderId)
	return nil
}
func resourceAlicloudResourceManagerFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("folder_name") {
		request := resourcemanager.CreateUpdateFolderRequest()
		request.FolderId = d.Id()
		request.NewFolderName = d.Get("folder_name").(string)
		raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
			return resourcemanagerClient.UpdateFolder(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudResourceManagerFolderRead(d, meta)
}
func resourceAlicloudResourceManagerFolderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := resourcemanager.CreateDeleteFolderRequest()
	request.FolderId = d.Id()
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.DeleteFolder(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Folder", "EntityNotExists.ResourceDirectory"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

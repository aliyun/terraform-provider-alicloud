package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudResourceManagerResourceDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerResourceDirectoryCreate,
		Read:   resourceAlicloudResourceManagerResourceDirectoryRead,
		Delete: resourceAlicloudResourceManagerResourceDirectoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"master_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_folder_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerResourceDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateInitResourceDirectoryRequest()
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.InitResourceDirectory(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_resource_directory", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*resourcemanager.InitResourceDirectoryResponse)
	d.SetId(fmt.Sprintf("%v", response.ResourceDirectory.ResourceDirectoryId))

	return resourceAlicloudResourceManagerResourceDirectoryRead(d, meta)
}
func resourceAlicloudResourceManagerResourceDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerResourceDirectory(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("master_account_id", object.MasterAccountId)
	d.Set("master_account_name", object.MasterAccountName)
	d.Set("root_folder_id", object.RootFolderId)
	return nil
}
func resourceAlicloudResourceManagerResourceDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := resourcemanager.CreateDestroyResourceDirectoryRequest()

	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.DestroyResourceDirectory(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

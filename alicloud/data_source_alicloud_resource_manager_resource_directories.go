package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudResourceManagerResourceDirectories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudResourceManagerResourceDirectorysRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_folder_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudResourceManagerResourceDirectorysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateGetResourceDirectoryRequest()
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.GetResourceDirectory(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_resource_directories", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*resourcemanager.GetResourceDirectoryResponse)

	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"master_account_id":     response.ResourceDirectory.MasterAccountId,
		"master_account_name":   response.ResourceDirectory.MasterAccountName,
		"id":                    response.ResourceDirectory.ResourceDirectoryId,
		"resource_directory_id": response.ResourceDirectory.ResourceDirectoryId,
		"root_folder_id":        response.ResourceDirectory.RootFolderId,
	}
	s = append(s, mapping)
	d.SetId(response.ResourceDirectory.ResourceDirectoryId)

	if err := d.Set("directories", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

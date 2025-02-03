package alicloud

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudResourceManagerResourceDirectories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudResourceManagerResourceDirectoriesRead,
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudResourceManagerResourceDirectoriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetResourceDirectory"
	request := make(map[string]interface{})
	var response map[string]interface{}
	var err error
	response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_resource_directories", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"master_account_id":     response["ResourceDirectory"].(map[string]interface{})["MasterAccountId"],
		"master_account_name":   response["ResourceDirectory"].(map[string]interface{})["MasterAccountName"],
		"id":                    fmt.Sprint(response["ResourceDirectory"].(map[string]interface{})["ResourceDirectoryId"]),
		"resource_directory_id": fmt.Sprint(response["ResourceDirectory"].(map[string]interface{})["ResourceDirectoryId"]),
		"root_folder_id":        response["ResourceDirectory"].(map[string]interface{})["RootFolderId"],
		"status":                response["ResourceDirectory"].(map[string]interface{})["ScpStatus"],
	}
	s = append(s, mapping)

	d.SetId(fmt.Sprint(response["ResourceDirectory"].(map[string]interface{})["ResourceDirectoryId"]))

	if err := d.Set("directories", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

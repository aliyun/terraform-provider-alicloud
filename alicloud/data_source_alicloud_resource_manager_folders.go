package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudResourceManagerFolders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudResourceManagerFoldersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"parent_folder_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"folders": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"folder_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"folder_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudResourceManagerFoldersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateListFoldersForParentRequest()
	if v, ok := d.GetOk("parent_folder_id"); ok {
		request.ParentFolderId = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []resourcemanager.Folder
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for {
		raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
			return resourcemanagerClient.ListFoldersForParent(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_folders", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*resourcemanager.ListFoldersForParentResponse)

		for _, item := range response.Folders.Folder {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.FolderName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.FolderId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.Folders.Folder) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, len(objects))
	names := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))

	for i, object := range objects {
		mapping := map[string]interface{}{
			"id":          object.FolderId,
			"folder_id":   object.FolderId,
			"folder_name": object.FolderName,
		}
		ids[i] = object.FolderId
		names[i] = object.FolderName
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("folders", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

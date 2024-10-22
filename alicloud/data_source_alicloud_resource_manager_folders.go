package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudResourceManagerFolders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudResourceManagerFoldersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"parent_folder_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"query_keyword": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
						"parent_folder_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudResourceManagerFoldersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListFoldersForParent"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	if v, ok := d.GetOk("parent_folder_id"); ok {
		request["ParentFolderId"] = v
	}

	if v, ok := d.GetOk("query_keyword"); ok {
		request["QueryKeyword"] = v
	}

	var objects []map[string]interface{}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var folderNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		folderNameRegex = r
	}

	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_folders", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Folders.Folder", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Folders.Folder", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["FolderId"])]; !ok {
					continue
				}
			}

			if folderNameRegex != nil {
				if !folderNameRegex.MatchString(fmt.Sprint(item["FolderName"])) {
					continue
				}
			}

			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(object["FolderId"]),
			"folder_id":   fmt.Sprint(object["FolderId"]),
			"folder_name": object["FolderName"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["FolderName"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["FolderId"])
		resourceManagerService := ResourcemanagerService{client}

		resourceManagerFolderDetail, err := resourceManagerService.DescribeResourceManagerFolder(id)
		if err != nil {
			return WrapError(err)
		}

		mapping["parent_folder_id"] = resourceManagerFolderDetail["ParentFolderId"]

		s = append(s, mapping)
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

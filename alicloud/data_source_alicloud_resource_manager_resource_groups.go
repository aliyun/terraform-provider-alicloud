package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudResourceManagerResourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudResourceManagerResourceGroupsRead,
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Creating", "Deleting", "OK", "PendingDelete"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
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

func dataSourceAlicloudResourceManagerResourceGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateListResourceGroupsRequest()
	if v, ok := d.GetOk("status"); ok {
		request.Status = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []resourcemanager.ResourceGroup
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
			return resourcemanagerClient.ListResourceGroups(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_resource_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*resourcemanager.ListResourceGroupsResponse)

		for _, item := range response.ResourceGroups.ResourceGroup {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.Name) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.Id]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.ResourceGroups.ResourceGroup) < PageSizeLarge {
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
			"account_id":   object.AccountId,
			"create_date":  object.CreateDate,
			"display_name": object.DisplayName,
			"id":           object.Id,
			"name":         object.Name,
			"status":       object.Status,
		}
		ids[i] = object.Id
		names[i] = object.Name
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

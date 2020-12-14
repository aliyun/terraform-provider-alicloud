package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudResourceManagerRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudResourceManagerRolesRead,
		Schema: map[string]*schema.Schema{
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_session_duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudResourceManagerRolesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateListRolesRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []resourcemanager.Role
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
			return resourcemanagerClient.ListRoles(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_roles", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*resourcemanager.ListRolesResponse)

		for _, item := range response.Roles.Role {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.RoleName) {
					continue
				}
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[item.RoleName]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.Roles.Role) < PageSizeLarge {
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
			"arn":                  object.Arn,
			"create_date":          object.CreateDate,
			"description":          object.Description,
			"max_session_duration": object.MaxSessionDuration,
			"role_id":              object.RoleId,
			"id":                   object.RoleName,
			"role_name":            object.RoleName,
			"update_date":          object.UpdateDate,
		}
		ids[i] = object.RoleName
		names[i] = object.RoleName
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("roles", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

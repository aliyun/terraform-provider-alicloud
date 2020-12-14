package alicloud

import (
	"regexp"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	dms_enterprise "github.com/aliyun/alibaba-cloud-sdk-go/services/dms-enterprise"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDmsEnterpriseUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDmsEnterpriseUsersRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"search_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLE", "NORMAL", "DELETE"}, false),
			},
			"tid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mobile": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"role_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"role_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nick_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDmsEnterpriseUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := dms_enterprise.CreateListUsersRequest()
	if v, ok := d.GetOk("role"); ok {
		request.Role = v.(string)
	}
	if v, ok := d.GetOk("search_key"); ok {
		request.SearchKey = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		request.UserState = v.(string)
	}
	if v, ok := d.GetOk("tid"); ok {
		request.Tid = requests.NewInteger(v.(int))
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []dms_enterprise.User
	var userNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		userNameRegex = r
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
	var response *dms_enterprise.ListUsersResponse
	for {
		raw, err := client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
			return dms_enterpriseClient.ListUsers(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_users", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*dms_enterprise.ListUsersResponse)

		for _, item := range response.UserList.User {
			if userNameRegex != nil {
				if !userNameRegex.MatchString(item.NickName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[strconv.FormatInt(item.Uid, 10)]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.UserList.User) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"mobile":     object.Mobile,
			"parent_uid": object.ParentUid,
			"role_ids":   object.RoleIdList.RoleIds,
			"role_names": object.RoleNameList.RoleNames,
			"status":     object.State,
			"id":         strconv.FormatInt(object.Uid, 10),
			"user_id":    object.UserId,
			"user_name":  object.NickName,
			"nick_name":  object.NickName,
		}
		ids = append(ids, strconv.FormatInt(object.Uid, 10))
		names = append(names, object.NickName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("users", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

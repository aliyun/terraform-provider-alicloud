package alicloud

import (
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudRamUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamUsersRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyType,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Computed values
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_login_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	allUsers := []interface{}{}

	allUsersMap := make(map[string]interface{})
	groupFilterUsersMap := make(map[string]interface{})
	policyFilterUsersMap := make(map[string]interface{})

	dataMap := []map[string]interface{}{}

	groupName, groupNameOk := d.GetOk("group_name")
	policyName, policyNameOk := d.GetOk("policy_name")
	policyType, policyTypeOk := d.GetOk("policy_type")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if policyTypeOk && !policyNameOk {
		return WrapError(Error("You must set 'policy_name' at one time when you set 'policy_type'."))
	}

	// all users
	request := ram.CreateListUsersRequest()
	for {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsers(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_users", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*ram.ListUsersResponse)
		for _, v := range response.Users.User {
			if nameRegexOk {
				r := regexp.MustCompile(nameRegex.(string))
				if !r.MatchString(v.UserName) {
					continue
				}
			}
			allUsersMap[v.UserName] = v
		}
		if !response.IsTruncated {
			break
		}
		request.Marker = response.Marker
	}

	// users for group
	if groupNameOk {
		request := ram.CreateListUsersForGroupRequest()
		request.GroupName = groupName.(string)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsersForGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_users", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*ram.ListUsersForGroupResponse)
		for _, v := range response.Users.User {
			groupFilterUsersMap[v.UserName] = v
		}
		dataMap = append(dataMap, groupFilterUsersMap)
	}

	// users which attach with this policy
	if policyNameOk {
		pType := "System"
		if policyTypeOk {
			pType = policyType.(string)
		}
		request := ram.CreateListEntitiesForPolicyRequest()
		request.PolicyName = policyName.(string)
		request.PolicyType = pType
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListEntitiesForPolicy(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_users", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*ram.ListEntitiesForPolicyResponse)
		for _, v := range response.Users.User {
			policyFilterUsersMap[v.UserName] = v
		}
		dataMap = append(dataMap, policyFilterUsersMap)
	}

	// GetIntersection of each map
	allUsers = ramService.GetIntersection(dataMap, allUsersMap)

	return ramUsersDescriptionAttributes(d, allUsers)
}

func ramUsersDescriptionAttributes(d *schema.ResourceData, users []interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, v := range users {
		user := v.(ram.User)
		mapping := map[string]interface{}{
			"id":              user.UserId,
			"name":            user.UserName,
			"create_date":     user.CreateDate,
			"last_login_date": user.LastLoginDate,
		}
		log.Printf("[DEBUG] alicloud_ram_users - adding user: %v", mapping)
		ids = append(ids, user.UserId)
		names = append(names, user.UserName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("users", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

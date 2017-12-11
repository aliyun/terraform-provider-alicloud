package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamGroupName,
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
	conn := meta.(*AliyunClient).ramconn
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
		return fmt.Errorf("You must set 'policy_name' at one time when you set 'policy_type'.")
	}

	// all users
	args := ram.ListUserRequest{}
	for {
		resp, err := conn.ListUsers(args)
		if err != nil {
			return fmt.Errorf("ListUsers got an error: %#v", err)
		}
		for _, v := range resp.Users.User {
			if nameRegexOk {
				r := regexp.MustCompile(nameRegex.(string))
				if !r.MatchString(v.UserName) {
					continue
				}
			}
			allUsersMap[v.UserName] = v
		}
		if !resp.IsTruncated {
			break
		}
		args.Marker = resp.Marker
	}

	// users for group
	if groupNameOk {
		resp, err := conn.ListUsersForGroup(ram.GroupQueryRequest{GroupName: groupName.(string)})
		if err != nil {
			return fmt.Errorf("ListUsersForGroup got an error: %#v", err)
		}

		for _, v := range resp.Users.User {
			groupFilterUsersMap[v.UserName] = v
		}
		dataMap = append(dataMap, groupFilterUsersMap)
	}

	// users which attach with this policy
	if policyNameOk {
		pType := ram.System
		if policyTypeOk {
			pType = ram.Type(policyType.(string))
		}
		resp, err := conn.ListEntitiesForPolicy(ram.PolicyRequest{PolicyName: policyName.(string), PolicyType: pType})
		if err != nil {
			return fmt.Errorf("ListEntitiesForPolicy got an error: %#v", err)
		}

		for _, v := range resp.Users.User {
			policyFilterUsersMap[v.UserName] = v
		}
		dataMap = append(dataMap, policyFilterUsersMap)
	}

	// GetIntersection of each map
	allUsers = GetIntersection(dataMap, allUsersMap)

	if len(allUsers) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_ram_users - Users found: %#v", allUsers)

	return ramUsersDescriptionAttributes(d, allUsers)
}

func ramUsersDescriptionAttributes(d *schema.ResourceData, users []interface{}) error {
	var ids []string
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
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("users", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

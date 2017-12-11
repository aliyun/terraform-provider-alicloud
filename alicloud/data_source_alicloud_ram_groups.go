package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"regexp"
)

func dataSourceAlicloudRamGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
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
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamGroupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	allGroups := []interface{}{}

	allGroupsMap := make(map[string]interface{})
	userFilterGroupsMap := make(map[string]interface{})
	policyFilterGroupsMap := make(map[string]interface{})

	dataMap := []map[string]interface{}{}

	userName, userNameOk := d.GetOk("user_name")
	policyName, policyNameOk := d.GetOk("policy_name")
	policyType, policyTypeOk := d.GetOk("policy_type")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if policyTypeOk && !policyNameOk {
		return fmt.Errorf("You must set 'policy_name' at one time when you set 'policy_type'.")
	}

	// groups filtered by name_regex
	args := ram.GroupListRequest{}
	for {
		resp, err := conn.ListGroup(args)
		if err != nil {
			return fmt.Errorf("ListGroup got an error: %#v", err)
		}
		for _, v := range resp.Groups.Group {
			if nameRegexOk {
				r := regexp.MustCompile(nameRegex.(string))
				if !r.MatchString(v.GroupName) {
					continue
				}
			}
			allGroupsMap[v.GroupName] = v
		}
		if !resp.IsTruncated {
			break
		}
		args.Marker = resp.Marker
	}

	// groups for user
	if userNameOk {
		resp, err := conn.ListGroupsForUser(ram.UserQueryRequest{UserName: userName.(string)})
		if err != nil {
			return fmt.Errorf("ListGroupsForUser got an error: %#v", err)
		}

		for _, v := range resp.Groups.Group {
			userFilterGroupsMap[v.GroupName] = v
		}
		dataMap = append(dataMap, userFilterGroupsMap)
	}

	// groups which attach with this policy
	if policyNameOk {
		pType := ram.System
		if policyTypeOk {
			pType = ram.Type(policyType.(string))
		}
		resp, err := conn.ListEntitiesForPolicy(ram.PolicyRequest{PolicyName: policyName.(string), PolicyType: pType})
		if err != nil {
			return fmt.Errorf("ListEntitiesForPolicy got an error: %#v", err)
		}

		for _, v := range resp.Groups.Group {
			policyFilterGroupsMap[v.GroupName] = v
		}
		dataMap = append(dataMap, policyFilterGroupsMap)
	}

	// GetIntersection of each map
	allGroups = GetIntersection(dataMap, allGroupsMap)

	if len(allGroups) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_ram_groups - Groups found: %#v", allGroups)

	return ramGroupsDescriptionAttributes(d, allGroups)
}

func ramGroupsDescriptionAttributes(d *schema.ResourceData, groups []interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, v := range groups {
		group := v.(ram.Group)
		mapping := map[string]interface{}{
			"name":     group.GroupName,
			"comments": group.Comments,
		}
		log.Printf("[DEBUG] alicloud_ram_groups - adding group: %v", mapping)
		ids = append(ids, v.(ram.Group).GroupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

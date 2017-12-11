package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudRamPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamPoliciesRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyType,
			},
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
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"role_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attachment_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"document": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	allPolicies := []interface{}{}

	allPoliciesMap := make(map[string]interface{})
	userFilterPoliciesMap := make(map[string]interface{})
	groupFilterPoliciesMap := make(map[string]interface{})
	roleFilterPoliciesMap := make(map[string]interface{})

	dataMap := []map[string]interface{}{}

	userName, userNameOk := d.GetOk("user_name")
	groupName, groupNameOk := d.GetOk("group_name")
	roleName, roleNameOk := d.GetOk("role_name")
	policyType, policyTypeOk := d.GetOk("type")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	// policies filtered by name_regex and type
	args := ram.PolicyQueryRequest{}
	for {
		resp, err := conn.ListPolicies(args)
		if err != nil {
			return fmt.Errorf("ListPolicies got an error: %#v", err)
		}
		for _, v := range resp.Policies.Policy {
			if policyTypeOk && policyType.(string) != v.PolicyType {
				continue
			}
			if nameRegexOk {
				r := regexp.MustCompile(nameRegex.(string))
				if !r.MatchString(v.PolicyName) {
					continue
				}
			}
			allPoliciesMap[v.PolicyType+v.PolicyName] = v
		}
		if !resp.IsTruncated {
			break
		}
		args.Marker = resp.Marker
	}

	// policies for user
	if userNameOk {
		resp, err := conn.ListPoliciesForUser(ram.UserQueryRequest{UserName: userName.(string)})
		if err != nil {
			return fmt.Errorf("ListPoliciesForUser got an error: %#v", err)
		}

		for _, v := range resp.Policies.Policy {
			userFilterPoliciesMap[v.PolicyType+v.PolicyName] = v
		}
		dataMap = append(dataMap, userFilterPoliciesMap)
	}

	// policies for group
	if groupNameOk {
		resp, err := conn.ListPoliciesForGroup(ram.GroupQueryRequest{GroupName: groupName.(string)})
		if err != nil {
			return fmt.Errorf("ListPoliciesForGroup got an error: %#v", err)
		}

		for _, v := range resp.Policies.Policy {
			groupFilterPoliciesMap[v.PolicyType+v.PolicyName] = v
		}
		dataMap = append(dataMap, groupFilterPoliciesMap)
	}

	// policies for role
	if roleNameOk {
		resp, err := conn.ListPoliciesForRole(ram.RoleQueryRequest{RoleName: roleName.(string)})
		if err != nil {
			return fmt.Errorf("ListPoliciesForRole got an error: %#v", err)
		}

		for _, v := range resp.Policies.Policy {
			roleFilterPoliciesMap[v.PolicyType+v.PolicyName] = v
		}
		dataMap = append(dataMap, roleFilterPoliciesMap)
	}

	// GetIntersection of each map
	allPolicies = GetIntersection(dataMap, allPoliciesMap)

	if len(allPolicies) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_ram_policies - Policies found: %#v", allPolicies)

	return ramPoliciesDescriptionAttributes(d, allPolicies, meta)
}

func ramPoliciesDescriptionAttributes(d *schema.ResourceData, policies []interface{}, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	var ids []string
	var s []map[string]interface{}
	for _, v := range policies {
		policy := v.(ram.Policy)
		resp, err := conn.GetPolicyVersionNew(ram.PolicyRequest{
			PolicyName: policy.PolicyName,
			PolicyType: ram.Type(policy.PolicyType),
			VersionId:  policy.DefaultVersion,
		})
		if err != nil {
			return err
		}

		mapping := map[string]interface{}{
			"name":             policy.PolicyName,
			"type":             policy.PolicyType,
			"description":      policy.Description,
			"default_version":  policy.DefaultVersion,
			"attachment_count": int(policy.AttachmentCount),
			"create_date":      policy.CreateDate,
			"update_date":      policy.UpdateDate,
			"document":         resp.PolicyVersion.PolicyDocument,
		}

		log.Printf("[DEBUG] alicloud_ram_policies - adding policy: %v", mapping)
		ids = append(ids, policy.PolicyName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("policies", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

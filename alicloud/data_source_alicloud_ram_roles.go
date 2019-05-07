package alicloud

import (
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudRamRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamRolesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
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
			"roles": {
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
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assume_role_policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document": {
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
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamRolesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	allRoles := []interface{}{}

	allRolesMap := make(map[string]interface{})

	dataMap := []interface{}{}

	policyName, policyNameOk := d.GetOk("policy_name")
	policyType, policyTypeOk := d.GetOk("policy_type")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	// all roles

	request := ram.CreateListRolesRequest()
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListRoles(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_roles", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw)

	response, _ := raw.(*ram.ListRolesResponse)
	for _, v := range response.Roles.Role {
		if nameRegexOk {
			r := regexp.MustCompile(nameRegex.(string))
			if !r.MatchString(v.RoleName) {
				continue
			}
		}
		allRolesMap[v.RoleName] = v
		allRoles = append(allRoles, v)
	}

	// roles which attach with this policy
	if policyNameOk {
		pType := "System"
		if policyTypeOk {
			pType = policyType.(string)
		}
		request := ram.CreateListEntitiesForPolicyRequest()
		request.PolicyType = pType
		request.PolicyName = policyName.(string)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListEntitiesForPolicy(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_roles", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*ram.ListEntitiesForPolicyResponse)
		for _, v := range response.Roles.Role {
			role, ok := allRolesMap[v.RoleName]
			if ok {
				dataMap = append(dataMap, role)
			}
		}
		allRoles = dataMap
	}
	return ramRolesDescriptionAttributes(d, meta, allRoles)
}

func ramRolesDescriptionAttributes(d *schema.ResourceData, meta interface{}, roles []interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, v := range roles {
		role := v.(ram.Role)
		request := ram.CreateGetRoleRequest()
		request.RoleName = role.RoleName
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetRole(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_roles", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*ram.GetRoleResponse)
		mapping := map[string]interface{}{
			"id":                          role.RoleId,
			"name":                        role.RoleName,
			"arn":                         role.Arn,
			"description":                 response.Role.Description,
			"create_date":                 role.CreateDate,
			"update_date":                 role.UpdateDate,
			"assume_role_policy_document": response.Role.AssumeRolePolicyDocument,
			"document":                    response.Role.AssumeRolePolicyDocument,
		}
		log.Printf("[DEBUG] alicloud_ram_roles - adding role: %v", mapping)
		ids = append(ids, role.RoleId)
		names = append(names, role.RoleName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("roles", s); err != nil {
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

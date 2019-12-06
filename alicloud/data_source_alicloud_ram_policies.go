package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudRamPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamPoliciesRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"System", "Custom"}, false),
			},
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
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"role_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
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
	request := ram.CreateListPoliciesRequest()
	request.RegionId = client.RegionId
	for {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPolicies(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ram.ListPoliciesResponse)
		for _, v := range response.Policies.Policy {
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
		if !response.IsTruncated {
			break
		}
		request.Marker = response.Marker
	}

	// policies for user
	if userNameOk {
		request := ram.CreateListPoliciesForUserRequest()
		request.UserName = userName.(string)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForUser(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ram.ListPoliciesForUserResponse)
		for _, v := range response.Policies.Policy {
			userFilterPoliciesMap[v.PolicyType+v.PolicyName] = v
		}
		dataMap = append(dataMap, userFilterPoliciesMap)
	}

	// policies for group
	if groupNameOk {
		request := ram.CreateListPoliciesForGroupRequest()
		request.GroupName = groupName.(string)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ram.ListPoliciesForGroupResponse)
		for _, v := range response.Policies.Policy {
			groupFilterPoliciesMap[v.PolicyType+v.PolicyName] = v
		}
		dataMap = append(dataMap, groupFilterPoliciesMap)
	}

	// policies for role
	if roleNameOk {
		request := ram.CreateListPoliciesForRoleRequest()
		request.RoleName = roleName.(string)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForRole(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ram.ListPoliciesForRoleResponse)
		for _, v := range response.Policies.Policy {
			roleFilterPoliciesMap[v.PolicyType+v.PolicyName] = v
		}
		dataMap = append(dataMap, roleFilterPoliciesMap)
	}

	// GetIntersection of each map
	allPolicies = ramService.GetIntersection(dataMap, allPoliciesMap)

	return ramPoliciesDescriptionAttributes(d, allPolicies, meta)
}

func ramPoliciesDescriptionAttributes(d *schema.ResourceData, policies []interface{}, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var ids []string
	var s []map[string]interface{}
	for _, v := range policies {
		policy := v.(ram.Policy)
		request := ram.CreateGetPolicyVersionRequest()
		request.PolicyName = policy.PolicyName
		request.PolicyType = policy.PolicyType
		request.VersionId = policy.DefaultVersion
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetPolicyVersion(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ram.GetPolicyVersionResponse)
		mapping := map[string]interface{}{
			"name":             policy.PolicyName,
			"type":             policy.PolicyType,
			"description":      policy.Description,
			"default_version":  policy.DefaultVersion,
			"attachment_count": int(policy.AttachmentCount),
			"create_date":      policy.CreateDate,
			"update_date":      policy.UpdateDate,
			"document":         response.PolicyVersion.PolicyDocument,
		}

		ids = append(ids, policy.PolicyName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

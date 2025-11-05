package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudRamPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudRamPoliciesRead,
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
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"System", "Custom"}, false),
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(0, 64),
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"role_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(0, 64),
			},
			"tags": tagsSchemaForceNew(),
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"default_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attachment_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_id": {
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
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
							Removed:  "Field `user_name` has been removed from provider version 1.262.1.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudRamPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListPolicies"
	request := make(map[string]interface{})

	request["MaxItems"] = PageSizeLarge

	if v, ok := d.GetOk("type"); ok {
		request["PolicyType"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMaps := ConvertTags(v.(map[string]interface{}))

		tagsMapsJson, err := convertListMapToJsonString(tagsMaps)
		if err != nil {
			return WrapError(err)
		}

		request["Tag"] = tagsMapsJson
	}

	userName, userNameOk := d.GetOk("user_name")
	groupName, groupNameOk := d.GetOk("group_name")
	roleName, roleNameOk := d.GetOk("role_name")

	var policyMaps []map[string]interface{}
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

	var policyNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}

		policyNameRegex = r
	}

	var response map[string]interface{}
	var err error

	if userNameOk {
		userPolicyAction := "ListPoliciesForUser"

		listPoliciesForUserRequest := map[string]interface{}{
			"UserName": userName,
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", userPolicyAction, nil, listPoliciesForUserRequest, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(userPolicyAction, response, listPoliciesForUserRequest)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", userPolicyAction, AlibabaCloudSdkGoERROR)
		}

		userResp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, userPolicyAction, "$.Policies.Policy", response)
		}

		userPolicyMap := make(map[string]interface{}, 0)
		if userResp != nil && len(userResp.([]interface{})) > 0 {
			for _, v := range userResp.([]interface{}) {
				item := v.(map[string]interface{})

				userPolicyMap[fmt.Sprintf("%v%v", item["PolicyName"], item["PolicyType"])] = item
			}
		}

		policyMaps = append(policyMaps, userPolicyMap)
	}

	if groupNameOk {
		groupPolicyAction := "ListPoliciesForGroup"

		listPoliciesForGroupRequest := map[string]interface{}{
			"GroupName": groupName,
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", groupPolicyAction, nil, listPoliciesForGroupRequest, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(groupPolicyAction, response, listPoliciesForGroupRequest)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", groupPolicyAction, AlibabaCloudSdkGoERROR)
		}

		groupResp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, groupPolicyAction, "$.Policies.Policy", response)
		}

		groupPolicyMap := make(map[string]interface{}, 0)
		if groupResp != nil && len(groupResp.([]interface{})) > 0 {
			for _, v := range groupResp.([]interface{}) {
				item := v.(map[string]interface{})

				groupPolicyMap[fmt.Sprintf("%v%v", item["PolicyName"], item["PolicyType"])] = item
			}
		}

		policyMaps = append(policyMaps, groupPolicyMap)
	}

	if roleNameOk {
		rolePolicyAction := "ListPoliciesForRole"

		listPoliciesForRoleRequest := map[string]interface{}{
			"RoleName": roleName,
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", rolePolicyAction, nil, listPoliciesForRoleRequest, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(rolePolicyAction, response, listPoliciesForRoleRequest)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", rolePolicyAction, AlibabaCloudSdkGoERROR)
		}

		roleResp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, rolePolicyAction, "$.Policies.Policy", response)
		}

		rolePolicyMap := make(map[string]interface{}, 0)
		if roleResp != nil && len(roleResp.([]interface{})) > 0 {
			for _, v := range roleResp.([]interface{}) {
				item := v.(map[string]interface{})

				rolePolicyMap[fmt.Sprintf("%v%v", item["PolicyName"], item["PolicyType"])] = item
			}
		}

		policyMaps = append(policyMaps, rolePolicyMap)
	}

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies.Policy", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PolicyName"])]; !ok {
					continue
				}
			}

			if policyNameRegex != nil && !policyNameRegex.MatchString(fmt.Sprint(item["PolicyName"])) {
				continue
			}

			if len(policyMaps) > 0 {
				isExist := false

				for _, policyMap := range policyMaps {
					if _, ok := policyMap[fmt.Sprintf("%v%v", item["PolicyName"], item["PolicyType"])]; ok {
						isExist = true
						break
					}
				}

				if !isExist {
					continue
				}
			}

			objects = append(objects, item)
		}

		if !response["IsTruncated"].(bool) {
			break
		}

		request["Marker"] = response["Marker"]
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":               fmt.Sprint(object["PolicyName"]),
			"policy_name":      fmt.Sprint(object["PolicyName"]),
			"name":             fmt.Sprint(object["PolicyName"]),
			"type":             object["PolicyType"],
			"description":      object["Description"],
			"default_version":  object["DefaultVersion"],
			"attachment_count": formatInt(object["AttachmentCount"]),
			"create_date":      object["CreateDate"],
			"update_date":      object["UpdateDate"],
		}

		if v, ok := object["Tags"]; ok {
			tags := v.(map[string]interface{})
			if tagMaps, ok := tags["Tag"]; ok {
				mapping["tags"] = tagsToMap(tagMaps)
			}
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["PolicyName"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["PolicyName"])

		getPolicyAction := "GetPolicy"

		getPolicyRequest := map[string]interface{}{
			"PolicyName": id,
			"PolicyType": object["PolicyType"],
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", getPolicyAction, nil, getPolicyRequest, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(getPolicyAction, response, getPolicyRequest)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", getPolicyAction, AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, getPolicyAction, "$", response)
		}

		if ramPolicyDocumentDetail, ok := v.(map[string]interface{})["DefaultPolicyVersion"].(map[string]interface{}); ok {
			mapping["policy_document"] = ramPolicyDocumentDetail["PolicyDocument"]
			mapping["document"] = ramPolicyDocumentDetail["PolicyDocument"]
			mapping["version_id"] = ramPolicyDocumentDetail["VersionId"]
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

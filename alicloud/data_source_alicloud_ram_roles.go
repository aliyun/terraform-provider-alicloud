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

func dataSourceAliCloudRamRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudRamRolesRead,
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
			"policy_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(0, 128),
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"System", "Custom"}, false),
			},
			"tags": tagsSchemaForceNew(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
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
						"assume_role_policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document": {
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
						"arn": {
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

func dataSourceAliCloudRamRolesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListRoles"
	request := make(map[string]interface{})

	request["MaxItems"] = PageSizeLarge

	if v, ok := d.GetOk("tags"); ok {
		tagsMaps := ConvertTags(v.(map[string]interface{}))

		tagsMapsJson, err := convertListMapToJsonString(tagsMaps)
		if err != nil {
			return WrapError(err)
		}

		request["Tag"] = tagsMapsJson
	}

	policyName, policyNameOk := d.GetOk("policy_name")
	policyType, policyTypeOk := d.GetOk("policy_type")

	var objects []map[string]interface{}
	allRolesMap := make(map[string]interface{})

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}

			idsMap[vv.(string)] = vv.(string)
		}
	}

	var roleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}

		roleNameRegex = r
	}

	var response map[string]interface{}
	var err error

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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_roles", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Roles.Role", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Roles.Role", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RoleId"])]; !ok {
					continue
				}
			}

			if roleNameRegex != nil && !roleNameRegex.MatchString(fmt.Sprint(item["RoleName"])) {
				continue
			}

			allRolesMap[fmt.Sprint(item["RoleName"])] = item
			objects = append(objects, item)
		}

		if !response["IsTruncated"].(bool) {
			break
		}

		request["Marker"] = response["Marker"]
	}

	if policyNameOk {
		action = "ListEntitiesForPolicy"

		listEntitiesForPolicyRequest := map[string]interface{}{
			"PolicyName": fmt.Sprint(policyName),
			"PolicyType": "System",
		}

		if policyTypeOk {
			listEntitiesForPolicyRequest["PolicyType"] = fmt.Sprint(policyType)
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, nil, listEntitiesForPolicyRequest, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, listEntitiesForPolicyRequest)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_roles", action, AlibabaCloudSdkGoERROR)
		}

		roleResp, err := jsonpath.Get("$.Roles.Role", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Roles.Role", response)
		}

		policyRoleMaps := make([]map[string]interface{}, 0)
		if roleResp != nil && len(roleResp.([]interface{})) > 0 {
			for _, v := range roleResp.([]interface{}) {
				item := v.(map[string]interface{})
				policyRoleMap, ok := allRolesMap[fmt.Sprint(item["RoleName"])].(map[string]interface{})
				if ok {
					policyRoleMaps = append(policyRoleMaps, policyRoleMap)
				}
			}
		}

		objects = policyRoleMaps
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(object["RoleId"]),
			"name":        fmt.Sprint(object["RoleName"]),
			"description": object["Description"],
			"arn":         object["Arn"],
			"create_date": object["CreateDate"],
			"update_date": object["UpdateDate"],
		}

		if v, ok := object["Tags"]; ok {
			tags := v.(map[string]interface{})
			if tagMaps, ok := tags["Tag"]; ok {
				mapping["tags"] = tagsToMap(tagMaps)
			}
		}

		ramServiceV2 := RamServiceV2{client}
		ramRoleDetail, err := ramServiceV2.DescribeRamRole(fmt.Sprint(object["RoleName"]))
		if err != nil {
			return WrapError(err)
		}

		mapping["assume_role_policy_document"] = ramRoleDetail["AssumeRolePolicyDocument"]
		mapping["document"] = ramRoleDetail["AssumeRolePolicyDocument"]

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["RoleName"])
		s = append(s, mapping)
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

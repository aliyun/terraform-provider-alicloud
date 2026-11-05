package resourcemanager

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourcePolicyVersions() *schema.Resource {
	return &schema.Resource{
		Read: DataSourcePolicyVersionsRead,
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Custom", "System"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default_version": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func DataSourcePolicyVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListPolicyVersions"
	request := make(map[string]interface{})
	request["PolicyName"] = d.Get("policy_name")
	request["PolicyType"] = d.Get("policy_type")
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_resource_manager_policy_versions", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)

	resp, err := jsonpath.Get("$.PolicyVersions.PolicyVersion", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.PolicyVersions.PolicyVersion", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(request["PolicyName"], ":", item["VersionId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                 fmt.Sprint(request["PolicyName"], ":", object["VersionId"]),
			"is_default_version": object["IsDefaultVersion"],
			"version_id":         object["VersionId"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(request["PolicyName"], ":", object["VersionId"]))
			s = append(s, mapping)
			continue
		}

		resourcemanagerService := ResourcemanagerService{client}
		id := fmt.Sprint(request["PolicyName"], ":", object["VersionId"])
		getResp, err := resourcemanagerService.DescribeResourceManagerPolicyVersion(id)
		if err != nil {
			return errs.WrapError(err)
		}
		mapping["policy_document"] = getResp["PolicyDocument"]
		ids = append(ids, fmt.Sprint(request["PolicyName"], ":", object["VersionId"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("versions", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}

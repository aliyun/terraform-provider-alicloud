package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCrEndpointAclPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCrEndpointAclPoliciesRead,
		Schema: map[string]*schema.Schema{
			"endpoint_type": {
				Type:         schema.TypeString,
				Required:     true,
				StateFunc:    func(v interface{}) string { return strings.ToLower(v.(string)) },
				ValidateFunc: validation.StringInSlice([]string{"internet", "Internet"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"module_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Registry"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entry": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCrEndpointAclPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "GetInstanceEndpoint"
	request := make(map[string]interface{})
	request["EndpointType"] = convertCrEndpointType(d.Get("endpoint_type").(string))
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("module_name"); ok {
		request["ModuleName"] = v
	}
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
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cr_instance_endpoint_acl_policies", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.AclEntries", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AclEntries", response)
	}
	result, _ := resp.([]interface{})
	endpointType := d.Get("endpoint_type").(string)
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(request["InstanceId"], ":", endpointType, ":", item["Entry"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"description":   object["Comment"],
			"endpoint_type": endpointType,
			"id":            fmt.Sprint(request["InstanceId"], ":", endpointType, ":", object["Entry"]),
			"entry":         fmt.Sprint(object["Entry"]),
			"instance_id":   request["InstanceId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
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

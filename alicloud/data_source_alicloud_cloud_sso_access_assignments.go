package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCloudSsoAccessAssignments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudSsoAccessAssignmentsRead,
		Schema: map[string]*schema.Schema{
			"access_configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"principal_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Group", "User"}, false),
			},
			"target_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RD-Account"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assignments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_configuration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_configuration_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_path_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCloudSsoAccessAssignmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListAccessAssignments"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("access_configuration_id"); ok {
		request["AccessConfigurationId"] = v
	}
	request["DirectoryId"] = d.Get("directory_id")
	if v, ok := d.GetOk("principal_type"); ok {
		request["PrincipalType"] = v
	}
	if v, ok := d.GetOk("target_id"); ok {
		request["TargetId"] = v
	}
	if v, ok := d.GetOk("target_type"); ok {
		request["TargetType"] = v
	}
	request["MaxResults"] = PageSizeSmall
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
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_sso_access_assignments", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AccessAssignments", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AccessAssignments", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(request["DirectoryId"], ":", item["AccessConfigurationId"], ":", item["TargetType"], ":", item["TargetId"], ":", item["PrincipalType"], ":", item["PrincipalId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"access_configuration_id":   object["AccessConfigurationId"],
			"access_configuration_name": object["AccessConfigurationName"],
			"directory_id":              request["DirectoryId"],
			"id":                        fmt.Sprint(request["DirectoryId"], ":", object["AccessConfigurationId"], ":", object["TargetType"], ":", object["TargetId"], ":", object["PrincipalType"], ":", object["PrincipalId"]),
			"principal_id":              fmt.Sprint(object["PrincipalId"]),
			"principal_name":            object["PrincipalName"],
			"principal_type":            object["PrincipalType"],
			"target_id":                 object["TargetId"],
			"target_name":               object["TargetName"],
			"target_path_name":          object["TargetPathName"],
			"target_type":               object["TargetType"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("assignments", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

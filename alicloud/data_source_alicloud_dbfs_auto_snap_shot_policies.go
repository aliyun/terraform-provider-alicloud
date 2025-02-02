package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDbfsAutoSnapShotPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDbfsAutoSnapShotPoliciesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  10,
			},
			"auto_snap_shot_policies": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"applied_dbfs_number": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"last_modified": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"policy_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"policy_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"repeat_weekdays": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"retention_days": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status_detail": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"time_points": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDbfsAutoSnapShotPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	setPagingRequest(d, request, PageSizeLarge)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListAutoSnapshotPolicies"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("DBFS", "2020-04-18", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dbfs_auto_snap_shot_policies", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.SnapshotPolicies", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.SnapshotPolicies", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil {
				if !nameRegex.MatchString(fmt.Sprint(item["PolicyName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PolicyId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})

		mapping := map[string]interface{}{
			"id": fmt.Sprint(object["PolicyId"]),
		}
		mapping["applied_dbfs_number"] = object["AppliedDbfsNumber"]
		mapping["create_time"] = object["CreatedTime"]
		mapping["last_modified"] = object["LastModified"]
		mapping["policy_id"] = object["PolicyId"]
		mapping["policy_name"] = object["PolicyName"]
		mapping["repeat_weekdays"] = object["RepeatWeekdays"].([]interface{})
		mapping["retention_days"] = object["RetentionDays"]
		mapping["status"] = object["Status"]
		mapping["status_detail"] = object["StatusDetail"]
		mapping["time_points"] = object["TimePoints"].([]interface{})

		ids = append(ids, fmt.Sprint(object["PolicyId"]))
		names = append(names, fmt.Sprint(object["PolicyName"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("auto_snap_shot_policies", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

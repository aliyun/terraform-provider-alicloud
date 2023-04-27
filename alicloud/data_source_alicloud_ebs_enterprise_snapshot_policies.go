package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEbsEnterpriseSnapshotPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEbsEnterpriseSnapshotPoliciesRead,
		Schema: map[string]*schema.Schema{
			"enterprise_snapshot_policy_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			},
			"enterprise_snapshot_policies": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"cross_region_copy_info": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"regions": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"retain_days": {
													Computed: true,
													Type:     schema.TypeInt,
												},
											},
										},
									},
								},
							},
						},
						"desc": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"enterprise_snapshot_policy_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"enterprise_snapshot_policy_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"retain_rule": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"number": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"time_interval": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"time_unit": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"schedule": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cron_expression": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"storage_rule": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_immediate_access": {
										Computed: true,
										Type:     schema.TypeBool,
									},
								},
							},
						},
						"tags": tagsSchema(),
						"target_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEbsEnterpriseSnapshotPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("enterprise_snapshot_policy_ids"); ok {
		request["PolicyIds"] = v.([]interface{})
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
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

	var enterpriseSnapshotPolicyNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		enterpriseSnapshotPolicyNameRegex = r
	}

	conn, err := client.NewEbsClient()
	if err != nil {
		return WrapError(err)
	}
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeEnterpriseSnapshotPolicy"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-07-30"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ebs_enterprise_snapshot_policies", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Policies", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PolicyId"])]; !ok {
					continue
				}
			}

			if enterpriseSnapshotPolicyNameRegex != nil && !enterpriseSnapshotPolicyNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
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

		createTime84 := object["CreateTime"]
		mapping["create_time"] = createTime84

		desc47 := object["Desc"]
		mapping["desc"] = desc47

		policyId10 := object["PolicyId"]
		mapping["enterprise_snapshot_policy_id"] = policyId10

		name65 := object["Name"]
		mapping["enterprise_snapshot_policy_name"] = name65

		resourceGroupId62 := object["ResourceGroupId"]
		mapping["resource_group_id"] = resourceGroupId62

		state29 := object["State"]
		mapping["status"] = state29

		targetType20 := object["TargetType"]
		mapping["target_type"] = targetType20

		crossRegionCopyInfo48Maps := make([]map[string]interface{}, 0)
		crossRegionCopyInfo48Map := make(map[string]interface{})
		crossRegionCopyInfo48Raw := object["CrossRegionCopyInfo"].(map[string]interface{})
		crossRegionCopyInfo48Map["enabled"] = crossRegionCopyInfo48Raw["Enabled"]
		regions48Maps := make([]map[string]interface{}, 0)
		regions48Raw := crossRegionCopyInfo48Raw["Regions"]
		for _, value1 := range regions48Raw.([]interface{}) {
			regions48 := value1.(map[string]interface{})
			regions48Map := make(map[string]interface{})
			regions48Map["region_id"] = regions48["RegionId"]
			regions48Map["retain_days"] = regions48["RetainDays"]
			regions48Maps = append(regions48Maps, regions48Map)
		}
		crossRegionCopyInfo48Map["regions"] = regions48Maps
		crossRegionCopyInfo48Maps = append(crossRegionCopyInfo48Maps, crossRegionCopyInfo48Map)
		mapping["cross_region_copy_info"] = crossRegionCopyInfo48Maps
		retainRule56Maps := make([]map[string]interface{}, 0)
		retainRule56Map := make(map[string]interface{})
		retainRule56Raw := object["RetainRule"].(map[string]interface{})
		retainRule56Map["number"] = retainRule56Raw["Number"]
		retainRule56Map["time_interval"] = retainRule56Raw["TimeInterval"]
		retainRule56Map["time_unit"] = retainRule56Raw["TimeUnit"]
		retainRule56Maps = append(retainRule56Maps, retainRule56Map)
		mapping["retain_rule"] = retainRule56Maps
		schedule95Maps := make([]map[string]interface{}, 0)
		schedule95Map := make(map[string]interface{})
		schedule95Raw := object["Schedule"].(map[string]interface{})
		schedule95Map["cron_expression"] = schedule95Raw["CronExpression"]
		schedule95Maps = append(schedule95Maps, schedule95Map)
		mapping["schedule"] = schedule95Maps
		storageRule66Maps := make([]map[string]interface{}, 0)
		storageRule66Map := make(map[string]interface{})
		storageRule66Raw := object["StorageRule"].(map[string]interface{})
		storageRule66Map["enable_immediate_access"] = storageRule66Raw["EnableImmediateAccess"]
		storageRule66Maps = append(storageRule66Maps, storageRule66Map)
		mapping["storage_rule"] = storageRule66Maps
		tagsMap := make(map[string]interface{})
		tagsRaw, _ := jsonpath.Get("$.Tags", object)
		if tagsRaw != nil {
			for _, value0 := range tagsRaw.([]interface{}) {
				tags := value0.(map[string]interface{})
				key := tags["TagKey"].(string)
				value := tags["TagValue"]
				if !ignoredTags(key, value) {
					tagsMap[key] = value
				}
			}
		}
		if len(tagsMap) > 0 {
			mapping["tags"] = tagsMap
		}

		ids = append(ids, fmt.Sprint(object["PolicyId"]))
		names = append(names, object["Name"])

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("enterprise_snapshot_policies", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

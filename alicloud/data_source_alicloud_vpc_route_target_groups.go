// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudVpcRouteTargetGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudVpcRouteTargetGroupRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_target_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_target_member_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"member_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enable_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_target_group_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_target_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_target_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// route_target_member_list is a datasource OUTPUT: every nested
						// field is Computed. terraform-plugin-sdk v1's HashResource
						// (the default Set func) only hashes Required/Optional fields and
						// SKIPS Computed ones (helper/schema/serialize.go), so a TypeSet
						// whose Elem resource is all-Computed hashes every element to the
						// same empty value and deduplicates N members down to 1. TypeList
						// preserves all elements (no hashing), which is the correct shape
						// for a read-only datasource output; the API returns the members in
						// a stable order. The resource keeps TypeSet because its
						// member_id/member_type/weight are Required (hashed, distinct).
						"route_target_member_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"member_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"member_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enable_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudVpcRouteTargetGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListRouteTargetGroups"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if v, ok := d.GetOk("route_target_group_id"); ok {
		request["RouteTargetGroupIds.1"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	routeTargetMemberListMemberIdJsonPath, err := jsonpath.Get("$.member_id", d.Get("route_target_member_list"))
	if err == nil {
		request["MemberId"] = convertToInterfaceArray(routeTargetMemberListMemberIdJsonPath)
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	request["VpcId"] = d.Get("vpc_id")
	request["MaxResults"] = PageSizeLarge
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if IsExpectedErrors(err, []string{"TaskConflict"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.RouteTargetGroups[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["RouteTargetGroupName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RouteTargetGroupId"])]; !ok {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["RouteTargetGroupId"]

		mapping["config_mode"] = objectRaw["ConfigMode"]
		mapping["create_time"] = objectRaw["CreateTime"]
		// ListRouteTargetGroups does not return RegionId (it documents it but
		// the backend omits it). A route target group is regional and is always
		// queried in the client's region, so derive region_id from the client
		// rather than leaving it empty.
		mapping["region_id"] = client.RegionId
		mapping["resource_group_id"] = objectRaw["ResourceGroupId"]
		mapping["route_target_group_description"] = objectRaw["RouteTargetGroupDescription"]
		mapping["route_target_group_name"] = objectRaw["RouteTargetGroupName"]
		mapping["status"] = objectRaw["Status"]
		mapping["vpc_id"] = objectRaw["VpcId"]
		mapping["route_target_group_id"] = objectRaw["RouteTargetGroupId"]

		routeTargetMemberListRaw := objectRaw["RouteTargetMemberList"]
		routeTargetMemberListMaps := make([]map[string]interface{}, 0)
		if routeTargetMemberListRaw != nil {
			for _, routeTargetMemberListChildRaw := range convertToInterfaceArray(routeTargetMemberListRaw) {
				routeTargetMemberListMap := make(map[string]interface{})
				routeTargetMemberListChildRaw := routeTargetMemberListChildRaw.(map[string]interface{})
				routeTargetMemberListMap["enable_status"] = routeTargetMemberListChildRaw["EnableStatus"]
				routeTargetMemberListMap["health_check_status"] = routeTargetMemberListChildRaw["HealthCheckStatus"]
				routeTargetMemberListMap["member_id"] = routeTargetMemberListChildRaw["MemberId"]
				routeTargetMemberListMap["member_type"] = routeTargetMemberListChildRaw["MemberType"]
				routeTargetMemberListMap["weight"] = routeTargetMemberListChildRaw["Weight"]

				routeTargetMemberListMaps = append(routeTargetMemberListMaps, routeTargetMemberListMap)
			}
		}
		mapping["route_target_member_list"] = routeTargetMemberListMaps
		tagsMaps := objectRaw["Tags"]
		mapping["tags"] = tagsToMap(tagsMaps)

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["RouteTargetGroupName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

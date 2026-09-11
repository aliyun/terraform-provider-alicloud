// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCloudFirewallVpcCenTrFirewallPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudFirewallVpcCenTrFirewallPolicyRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"firewall_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dest_candidate_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"candidate_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"candidate_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"policy_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_candidate_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"candidate_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"candidate_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tr_firewall_route_policy_id": {
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

func dataSourceAliCloudCloudFirewallVpcCenTrFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeTrFirewallV2RoutePolicyList"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("firewall_id"); ok {
		request["FirewallId"] = v
	}
	request["FirewallId"] = d.Get("firewall_id")
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["CurrentPage"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
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

		resp, _ := jsonpath.Get("$.TrFirewallRoutePolicies[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(request["FirewallId"], ":", item["TrFirewallRoutePolicyId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(request["FirewallId"], ":", objectRaw["TrFirewallRoutePolicyId"])

		mapping["policy_description"] = objectRaw["PolicyDescription"]
		mapping["policy_name"] = objectRaw["PolicyName"]
		mapping["policy_type"] = objectRaw["PolicyType"]
		mapping["status"] = objectRaw["PolicyStatus"]
		mapping["tr_firewall_route_policy_id"] = objectRaw["TrFirewallRoutePolicyId"]

		destCandidateListRaw := objectRaw["DestCandidateList"]
		destCandidateListMaps := make([]map[string]interface{}, 0)
		if destCandidateListRaw != nil {
			for _, destCandidateListChildRaw := range convertToInterfaceArray(destCandidateListRaw) {
				destCandidateListMap := make(map[string]interface{})
				destCandidateListChildRaw := destCandidateListChildRaw.(map[string]interface{})
				destCandidateListMap["candidate_id"] = destCandidateListChildRaw["CandidateId"]
				destCandidateListMap["candidate_type"] = destCandidateListChildRaw["CandidateType"]

				destCandidateListMaps = append(destCandidateListMaps, destCandidateListMap)
			}
		}
		mapping["dest_candidate_list"] = destCandidateListMaps
		srcCandidateListRaw := objectRaw["SrcCandidateList"]
		srcCandidateListMaps := make([]map[string]interface{}, 0)
		if srcCandidateListRaw != nil {
			for _, srcCandidateListChildRaw := range convertToInterfaceArray(srcCandidateListRaw) {
				srcCandidateListMap := make(map[string]interface{})
				srcCandidateListChildRaw := srcCandidateListChildRaw.(map[string]interface{})
				srcCandidateListMap["candidate_id"] = srcCandidateListChildRaw["CandidateId"]
				srcCandidateListMap["candidate_type"] = srcCandidateListChildRaw["CandidateType"]

				srcCandidateListMaps = append(srcCandidateListMaps, srcCandidateListMap)
			}
		}
		mapping["src_candidate_list"] = srcCandidateListMaps

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

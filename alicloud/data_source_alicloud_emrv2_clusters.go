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

func dataSourceAlicloudEmrV2Clusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEmrV2ClustersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cluster_states": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"STARTING", "START_FAILED", "RUNNING", "TERMINATED", "TERMINATING", "TERMINATE_FAILED", "TERMINATED_WITH_ERRORS"}, false),
				},
			},
			"payment_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				},
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"next_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_results": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ready_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"release_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_change_reason": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"emr_default_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudEmrV2ClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListClusters"
	request := map[string]interface{}{
		"NextToken":  "0",
		"MaxResults": PageSizeXLarge,
		"RegionId":   client.RegionId,
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		request["ClusterName"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("ids"); ok {
		request["ClusterIds"] = v
	}
	if v, ok := d.GetOk("cluster_types"); ok {
		request["ClusterTypes"] = v
	}
	if v, ok := d.GetOk("cluster_states"); ok {
		request["ClusterStates"] = v
	}
	if v, ok := d.GetOk("payment_types"); ok {
		request["PaymentTypes"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags []map[string]interface{}
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value,
			})
		}
		if len(tags) > 0 {
			request["Tags"] = tags
		}
	}
	if v, ok := d.GetOk("next_token"); ok {
		request["NextToken"] = v
	}
	if v, ok := d.GetOk("max_results"); ok {
		request["MaxResults"] = formatInt(v)
	}

	var clusterNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		clusterNameRegex = r
	}

	var response map[string]interface{}
	var err error
	var objects []interface{}

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Emr", "2021-03-20", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emrv2_clusters", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Clusters", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Clusters", response)
		}
		result, _ := resp.([]interface{})
		if _, ok := d.GetOkExists("next_token"); ok {
			objects = resp.([]interface{})
			break
		}

		for _, v := range result {
			item := v.(map[string]interface{})
			if clusterNameRegex != nil && !clusterNameRegex.MatchString(fmt.Sprint(item["ClusterName"])) {
				continue
			}
			objects = append(objects, item)
		}

		if len(result) < request["MaxResults"].(int) {
			break
		}
		nextToken, err := jsonpath.Get("$.NextToken", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NextToken", response)
		}
		request["NextToken"] = nextToken
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"cluster_id":        fmt.Sprint(object["ClusterId"]),
			"cluster_name":      object["ClusterName"],
			"cluster_type":      object["ClusterType"],
			"cluster_state":     fmt.Sprint(object["ClusterState"]),
			"payment_type":      object["PaymentType"],
			"create_time":       fmt.Sprint(object["CreateTime"]),
			"ready_time":        fmt.Sprint(object["ReadyTime"]),
			"expire_time":       fmt.Sprint(object["ExpireTime"]),
			"end_time":          fmt.Sprint(object["EndTime"]),
			"release_version":   object["ReleaseVersion"],
			"resource_group_id": object["ResourceGroupId"],
			"emr_default_role":  object["EmrDefaultRole"],
		}

		var tags []map[string]interface{}
		t, _ := jsonpath.Get("$.Tags", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["Key"].(string)
				value := t.(map[string]interface{})["Value"].(string)
				if !ignoredTags(key, value) {
					tags = append(tags, map[string]interface{}{
						"key":   key,
						"value": value,
					})
				}
			}
		}
		mapping["tags"] = tags
		if object["StateChangeReason"] != nil {
			mapping["state_change_reason"] = object["StateChangeReason"].(map[string]interface{})
		}

		ids = append(ids, fmt.Sprint(object["ClusterId"]))
		names = append(names, object["ClusterName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("total_count", formatInt(response["TotalCount"])); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

package alicloud

import (
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEmrV2ClusterInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEmrV2ClusterInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"node_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"node_group_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_states": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"next_token": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"max_results": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_renew": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_renew_duration_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_renew_duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
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

func dataSourceAlicloudEmrV2ClusterInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	emrService := EmrService{client}
	action := "ListNodes"
	maxResults := PageSizeSmall
	if v, ok := d.GetOk("max_results"); ok {
		maxResults = v.(int)
	}
	request := map[string]interface{}{
		"NextToken":  "0",
		"MaxResults": maxResults,
		"RegionId":   client.RegionId,
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		var instanceIds []string
		for _, item := range v.([]interface{}) {
			instanceIds = append(instanceIds, item.(string))
		}
		request["NodeIds"] = instanceIds
	}
	if v, ok := d.GetOk("instance_states"); ok && len(v.([]interface{})) > 0 {
		var instanceStates []string
		for _, item := range v.([]interface{}) {
			instanceStates = append(instanceStates, item.(string))
		}
		request["NodeStates"] = instanceStates
	}
	if v, ok := d.GetOk("next_token"); ok && v.(string) != "" {
		request["NextToken"] = v.(string)
	}
	var mergedNodeGroupIds []string
	if v, ok := d.GetOk("node_group_ids"); ok && len(v.([]interface{})) > 0 {
		for _, item := range v.([]interface{}) {
			mergedNodeGroupIds = append(mergedNodeGroupIds, item.(string))
		}
		request["NodeGroupIds"] = mergedNodeGroupIds
	}

	if v, ok := d.GetOk("node_group_names"); ok && len(v.([]interface{})) > 0 {
		if clusterId, exists := request["ClusterId"]; exists && clusterId.(string) != "" {
			nodeGroups, err := emrService.ListEmrV2NodeGroups(clusterId.(string), []string{})
			if err == nil {
				ngnSet := map[string]struct{}{}
				for _, ngn := range v.([]interface{}) {
					ngnSet[ngn.(string)] = struct{}{}
				}
				var nodeGroupIds []string
				if len(nodeGroups) > 0 {
					for _, item := range nodeGroups {
						nodeGroupMap := item.(map[string]interface{})
						if "TERMINATED" == nodeGroupMap["NodeGroupState"].(string) {
							continue
						}
						if _, exists := ngnSet[nodeGroupMap["NodeGroupName"].(string)]; exists {
							nodeGroupIds = append(nodeGroupIds, nodeGroupMap["NodeGroupId"].(string))
						}
					}
					if len(nodeGroupIds) > 0 {
						request["NodeGroupIds"] = nodeGroupIds
					}
				}
			} else {
				if !strings.Contains(err.Error(), "The specified parameter ClusterId is not valid") {
					return WrapError(err)
				}
			}
		}
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

	response := map[string]interface{}{
		"TotalCount": 0,
	}
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
			if strings.Contains(err.Error(), "The specified parameter ClusterId is not valid") {
				break
			}
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emrv2_cluster_instances", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Nodes", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Nodes", response)
		}
		result, _ := resp.([]interface{})
		if _, ok := d.GetOk("next_token"); !ok {
			objects = resp.([]interface{})
			break
		}

		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}

		if nt, exists := d.GetOk("next_token"); exists && nt.(string) != "" {
			break
		}

		_, nextTokenExists := response["NextToken"]
		if len(result) < request["MaxResults"].(int) {
			break
		} else if len(result) == request["MaxResults"].(int) && !nextTokenExists {
			break
		}
		nextToken, err := jsonpath.Get("$.NextToken", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NextToken", response)
		}
		request["NextToken"] = nextToken
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"instance_id":              object["NodeId"],
			"instance_name":            object["NodeName"],
			"instance_type":            object["InstanceType"],
			"instance_state":           object["NodeState"],
			"node_group_id":            object["NodeGroupId"],
			"node_group_type":          object["NodeGroupType"],
			"zone_id":                  object["ZoneId"],
			"public_ip":                object["PublicIp"],
			"private_ip":               object["PrivateIp"],
			"auto_renew":               object["AutoRenew"],
			"auto_renew_duration_unit": object["AutoRenewDurationUnit"],
			"auto_renew_duration":      object["AutoRenewDuration"],
			"create_time":              object["CreateTime"],
			"expire_time":              object["ExpireTime"],
		}

		ids = append(ids, object["NodeId"].(string))
		names = append(names, object["NodeName"].(string))
		s = append(s, mapping)
	}

	if len(ids) > 0 {
		d.SetId(dataResourceIdHash(ids))
	} else {
		d.SetId(request["ClusterId"].(string))
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
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

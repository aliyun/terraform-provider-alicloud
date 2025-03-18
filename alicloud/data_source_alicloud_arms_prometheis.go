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

func dataSourceAliCloudArmsPrometheis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudArmsPrometheisRead,
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
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchemaForceNew(),
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"prometheis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_clusters_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grafana_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_read_intra_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_read_inter_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_write_intra_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_write_inter_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"push_gate_way_intra_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"push_gate_way_inter_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http_api_intra_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http_api_inter_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudArmsPrometheisRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListPrometheusInstanceByTagAndResourceGroupId"
	request := make(map[string]interface{})

	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		request["Tag"] = tagsFromMap(v.(map[string]interface{}))
	}

	var objects []map[string]interface{}
	var prometheusNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		prometheusNameRegex = r
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

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_arms_prometheus", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data.PrometheusInstances", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.PrometheusInstances", response)
	}

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if prometheusNameRegex != nil && !prometheusNameRegex.MatchString(fmt.Sprint(item["ClusterName"])) {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["ClusterId"])]; !ok {
				continue
			}
		}

		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                  fmt.Sprint(object["ClusterId"]),
			"cluster_id":          fmt.Sprint(object["ClusterId"]),
			"cluster_type":        object["ClusterType"],
			"cluster_name":        object["ClusterName"],
			"vpc_id":              object["VpcId"],
			"vswitch_id":          object["VSwitchId"],
			"security_group_id":   object["SecurityGroupId"],
			"sub_clusters_json":   object["SubClustersJson"],
			"grafana_instance_id": object["GrafanaInstanceId"],
			"resource_group_id":   object["ResourceGroupId"],
			"tags":                tagsToMap(object["Tags"]),
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ClusterName"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(fmt.Sprint(object["ClusterId"]))
		armsService := ArmsService{client}

		armsPrometheus, err := armsService.DescribeArmsPrometheus(id)
		if err != nil {
			return WrapError(err)
		}

		mapping["remote_read_intra_url"] = armsPrometheus["RemoteReadIntraUrl"]
		mapping["remote_read_inter_url"] = armsPrometheus["RemoteReadInterUrl"]
		mapping["remote_write_intra_url"] = armsPrometheus["RemoteWriteIntraUrl"]
		mapping["remote_write_inter_url"] = armsPrometheus["RemoteWriteInterUrl"]
		mapping["push_gate_way_intra_url"] = armsPrometheus["PushGatewayIntraUrl"]
		mapping["push_gate_way_inter_url"] = armsPrometheus["PushGatewayInterUrl"]
		mapping["http_api_intra_url"] = armsPrometheus["HttpApiIntraUrl"]
		mapping["http_api_inter_url"] = armsPrometheus["HttpApiInterUrl"]
		mapping["auth_token"] = armsPrometheus["AuthToken"]

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("prometheis", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

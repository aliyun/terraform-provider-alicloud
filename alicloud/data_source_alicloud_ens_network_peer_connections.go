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

func dataSourceAlicloudEnsNetworkPeerConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEnsNetworkPeerConnectionsRead,
		Schema: map[string]*schema.Schema{
			"network_peer_connection_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ens_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accepting_network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ens_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_peer_connection_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEnsNetworkPeerConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeNetworkPeerConnections"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("ens_region_id"); ok {
		request["EnsRegionId"] = v
	}
	if v, ok := d.GetOk("network_peer_connection_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("network_id"); ok {
		request["NetworkIds.1"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ens_network_peer_connections", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.NetworkPeerConnects", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.NetworkPeerConnects", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil {
				name := fmt.Sprint(item["Name"])
				if !nameRegex.MatchString(name) {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"instance_id":                  fmt.Sprint(object["InstanceId"]),
			"network_id":                   fmt.Sprint(object["NetworkId"]),
			"accepting_network_id":         fmt.Sprint(object["AcceptingNetworkId"]),
			"ens_region_id":                fmt.Sprint(object["EnsRegionId"]),
			"status":                       fmt.Sprint(object["Status"]),
			"create_time":                  fmt.Sprint(object["CreationTime"]),
			"description":                  fmt.Sprint(object["Description"]),
			"network_peer_connection_name": fmt.Sprint(object["Name"]),
		}
		ids = append(ids, fmt.Sprint(mapping["instance_id"]))
		names = append(names, mapping["network_peer_connection_name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("connections", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

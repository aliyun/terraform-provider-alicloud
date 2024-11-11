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

func dataSourceAlicloudVpcPeerConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcPeerConnectionsRead,
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
			"peer_connection_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Accepting", "Activated", "Creating", "Deleted", "Deleting", "Expired", "Rejected", "Updating"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepting_ali_uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"accepting_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accepting_vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_connection_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpcPeerConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListVpcPeerConnections"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("peer_connection_name"); ok {
		request["Name"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var peerConnectionNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		peerConnectionNameRegex = r
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("VpcPeer", "2022-01-01", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_peer_connections", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.VpcPeerConnects", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.VpcPeerConnects", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if peerConnectionNameRegex != nil && !peerConnectionNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
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
	for _, object := range objects {
		mapping := map[string]interface{}{
			"accepting_ali_uid":    formatInt(object["AcceptingOwnerUid"]),
			"accepting_region_id":  object["AcceptingRegionId"],
			"accepting_vpc_id":     object["AcceptingVpc"].(map[string]interface{})["VpcId"],
			"bandwidth":            formatInt(object["Bandwidth"]),
			"create_time":          object["GmtCreate"],
			"description":          object["Description"],
			"id":                   fmt.Sprint(object["InstanceId"]),
			"peer_connection_id":   fmt.Sprint(object["InstanceId"]),
			"peer_connection_name": object["Name"],
			"status":               object["Status"],
			"vpc_id":               object["Vpc"].(map[string]interface{})["VpcId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
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

	if err := d.Set("connections", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

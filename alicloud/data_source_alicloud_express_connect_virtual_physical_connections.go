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

func dataSourceAlicloudExpressConnectVirtualPhysicalConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudExpressConnectVirtualPhysicalConnectionsRead,
		Schema: map[string]*schema.Schema{
			"business_status": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "FinancialLocked", "SecurityLocked"}, false),
			},
			"is_confirmed": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeBool,
			},
			"parent_physical_connection_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"virtual_physical_connection_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"virtual_physical_connection_status": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Confirmed", "UnConfirmed", "Deleted"}, false),
			},
			"vlan_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"vpconn_ali_uid": {
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
			"connections": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"access_point_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"ad_location": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"bandwidth": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"business_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"circuit_code": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"enabled_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"end_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"expect_spec": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"line_operator": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"loa_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"order_mode": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"parent_physical_connection_ali_uid": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"parent_physical_connection_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"peer_location": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"port_number": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"port_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"redundant_physical_connection_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"spec": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"virtual_physical_connection_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"virtual_physical_connection_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"virtual_physical_connection_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vlan_id": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"vpconn_ali_uid": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudExpressConnectVirtualPhysicalConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("business_status"); ok {
		request["VirtualPhysicalConnectionBusinessStatus"] = v
	}
	if v, ok := d.GetOk("is_confirmed"); ok {
		request["IsConfirmed"] = v
	}
	if v, ok := d.GetOk("parent_physical_connection_id"); ok {
		request["PhysicalConnectionId"] = v
	}
	if v, ok := d.GetOk("virtual_physical_connection_ids"); ok {
		request["VirtualPhysicalConnectionIds"] = v.([]interface{})
	}
	if v, ok := d.GetOk("virtual_physical_connection_status"); ok {
		request["VirtualPhysicalConnectionStatuses"] = v
	}
	if v, ok := d.GetOk("vlan_ids"); ok {
		request["VlanIds"] = v.([]interface{})
	}
	if v, ok := d.GetOk("vpconn_ali_uid"); ok {
		request["VirtualPhysicalConnectionAliUids"] = v
	}
	request["MaxResults"] = PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var virtualPhysicalConnectionNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		virtualPhysicalConnectionNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListVirtualPhysicalConnections"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_express_connect_virtual_physical_connections", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.VirtualPhysicalConnections", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.VirtualPhysicalConnections", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PhysicalConnectionId"])]; !ok {
					continue
				}
			}

			if virtualPhysicalConnectionNameRegex != nil && !virtualPhysicalConnectionNameRegex.MatchString(fmt.Sprint(item["Name"])) {
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
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                                 fmt.Sprint(object["PhysicalConnectionId"]),
			"access_point_id":                    object["AccessPointId"],
			"ad_location":                        object["AdLocation"],
			"bandwidth":                          object["Bandwidth"],
			"business_status":                    object["BusinessStatus"],
			"circuit_code":                       object["CircuitCode"],
			"create_time":                        object["CreationTime"],
			"description":                        object["Description"],
			"enabled_time":                       object["EnabledTime"],
			"end_time":                           object["EndTime"],
			"expect_spec":                        object["ExpectSpec"],
			"line_operator":                      object["LineOperator"],
			"loa_status":                         object["LoaStatus"],
			"order_mode":                         object["OrderMode"],
			"parent_physical_connection_ali_uid": object["ParentPhysicalConnectionAliUid"],
			"parent_physical_connection_id":      object["ParentPhysicalConnectionId"],
			"peer_location":                      object["PeerLocation"],
			"port_number":                        object["PortNumber"],
			"port_type":                          object["PortType"],
			"redundant_physical_connection_id":   object["RedundantPhysicalConnectionId"],
			"resource_group_id":                  object["ResourceGroupId"],
			"spec":                               object["Spec"],
			"status":                             object["Status"],
			"virtual_physical_connection_id":     object["PhysicalConnectionId"],
			"virtual_physical_connection_name":   object["Name"],
			"virtual_physical_connection_status": object["VirtualPhysicalConnectionStatus"],
			"vlan_id":                            formatInt(object["VlanId"]),
			"vpconn_ali_uid":                     object["AliUid"],
		}

		ids = append(ids, fmt.Sprint(object["PhysicalConnectionId"]))
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

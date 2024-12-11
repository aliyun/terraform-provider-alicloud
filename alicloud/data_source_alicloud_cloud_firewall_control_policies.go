package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCloudFirewallControlPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudFirewallControlPoliciesRead,
		Schema: map[string]*schema.Schema{
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"in", "out"}, false),
			},
			"acl_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"acl_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"accept", "drop", "log"}, false),
			},
			"destination": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"proto": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"TCP", " UDP", "ANY", "ICMP"}, false),
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"en", "zh"}, false),
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
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_port_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_port_group_ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dest_port_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_group_cidrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"destination_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_result_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hit_times": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"order": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"proto": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"release": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_group_cidrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"source_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field `upgrade_type` has been removed from provider version 1.213.0.",
			},
		},
	}
}

func dataSourceAliCloudCloudFirewallControlPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeControlPolicy"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["CurrentPage"] = 1

	request["Direction"] = d.Get("direction")

	if v, ok := d.GetOk("acl_action"); ok {
		request["AclAction"] = v
	}

	if v, ok := d.GetOk("acl_uuid"); ok {
		request["AclUuid"] = v
	}

	if v, ok := d.GetOk("destination"); ok {
		request["Destination"] = v
	}

	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}

	if v, ok := d.GetOk("proto"); ok {
		request["Proto"] = v
	}

	if v, ok := d.GetOk("source"); ok {
		request["Source"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	var endpoint string

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_firewall_control_policies", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Policys", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policys", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}

		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                      fmt.Sprintf("%v:%v", object["AclUuid"], object["Direction"]),
			"acl_uuid":                fmt.Sprint(object["AclUuid"]),
			"direction":               fmt.Sprint(object["Direction"]),
			"acl_action":              object["AclAction"],
			"application_id":          object["ApplicationId"],
			"application_name":        object["ApplicationName"],
			"description":             object["Description"],
			"dest_port":               object["DestPort"],
			"dest_port_group":         object["DestPortGroup"],
			"dest_port_group_ports":   object["DestPortGroupPorts"],
			"dest_port_type":          object["DestPortType"],
			"destination":             object["Destination"],
			"destination_group_cidrs": object["DestinationGroupCidrs"],
			"destination_group_type":  object["DestinationGroupType"],
			"destination_type":        object["DestinationType"],
			"dns_result":              object["DnsResult"],
			"dns_result_time":         fmt.Sprint(object["DnsResultTime"]),
			"hit_times":               fmt.Sprint(object["HitTimes"]),
			"order":                   formatInt(object["Order"]),
			"proto":                   object["Proto"],
			"release":                 formatBool(object["Release"]),
			"source":                  object["Source"],
			"source_group_cidrs":      object["SourceGroupCidrs"],
			"source_group_type":       object["SourceGroupType"],
			"source_type":             object["SourceType"],
		}

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

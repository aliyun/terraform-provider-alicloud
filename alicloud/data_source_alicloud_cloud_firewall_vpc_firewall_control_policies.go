package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCloudFirewallVpcFirewallControlPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudFirewallVpcFirewallControlPoliciesRead,
		Schema: map[string]*schema.Schema{
			"acl_action": {
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop", "log"}, false),
				Type:         schema.TypeString,
			},
			"acl_uuid": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"description": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"destination": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"lang": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en"}, false),
			},
			"member_uid": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"proto": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"ANY", "TCP", "UDP", "ICMP"}, false),
			},
			"release": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeBool,
			},
			"source": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"vpc_firewall_id": {
				Required: true,
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
				Default:  50,
			},
			"policies": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"acl_action": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"acl_uuid": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"application_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"application_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"description": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dest_port": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dest_port_group": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"dest_port_group_ports": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dest_port_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"destination": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"destination_group_cidrs": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"destination_group_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"destination_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"hit_times": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"member_uid": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"order": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"proto": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"release": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"source": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source_group_cidrs": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"source_group_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"source_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vpc_firewall_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCloudFirewallVpcFirewallControlPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("acl_action"); ok {
		request["AclAction"] = v
	}
	if v, ok := d.GetOk("acl_uuid"); ok {
		request["AclUuid"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("destination"); ok {
		request["Destination"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}
	if v, ok := d.GetOk("proto"); ok {
		request["Proto"] = v
	}
	if v, ok := d.GetOkExists("release"); ok {
		request["Release"] = v
	}
	if v, ok := d.GetOk("source"); ok {
		request["Source"] = v
	}
	request["VpcFirewallId"] = d.Get("vpc_firewall_id")
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["CurrentPage"] = v.(int)
	} else {
		request["CurrentPage"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
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

	var err error
	var endpoint string
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeVpcFirewallControlPolicy"
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
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_control_policies", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Policys", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policys", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(request["VpcFirewallId"], ":", item["AclUuid"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                     fmt.Sprint(request["VpcFirewallId"], ":", object["AclUuid"]),
			"acl_action":             object["AclAction"],
			"acl_uuid":               object["AclUuid"],
			"application_id":         object["ApplicationId"],
			"application_name":       object["ApplicationName"],
			"description":            object["Description"],
			"dest_port":              object["DestPort"],
			"dest_port_group":        object["DestPortGroup"],
			"dest_port_type":         object["DestPortType"],
			"destination":            object["Destination"],
			"destination_group_type": object["DestinationGroupType"],
			"destination_type":       object["DestinationType"],
			"hit_times":              object["HitTimes"],
			"member_uid":             object["MemberUid"],
			"order":                  object["Order"],
			"proto":                  object["Proto"],
			"release":                Interface2Bool(object["Release"]),
			"source":                 object["Source"],
			"source_group_type":      object["SourceGroupType"],
			"source_type":            object["SourceType"],
			"vpc_firewall_id":        request["VpcFirewallId"],
		}
		if v, ok := object["SourceGroupCidrs"]; ok {
			mapping["source_group_cidrs"] = v.([]interface{})
		}
		if v, ok := object["DestinationGroupCidrs"]; ok {
			mapping["destination_group_cidrs"] = v.([]interface{})
		}
		if v, ok := object["DestPortGroupPorts"]; ok {
			mapping["dest_port_group_ports"] = v.([]interface{})
		}

		ids = append(ids, fmt.Sprint(request["VpcFirewallId"], ":", object["AclUuid"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}
	return nil
}

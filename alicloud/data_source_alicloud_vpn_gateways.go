package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpnGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnGatewaysRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},

			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Init", "Provisioning", "Active", "Updating", "Deleting"}, false),
			},

			"business_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "FinancialLocked"}, false),
			},
			"enable_ipsec": {
				Type:       schema.TypeBool,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'enable_ipsec' has been deprecated from provider version 1.193.0 and it will be removed in the future version.",
			},
			"ssl_vpn": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, false),
			},
			"include_reservation_data": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"business_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ipsec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ssl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_vpn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_propagate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disaster_recovery_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disaster_recovery_internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"ssl_vpn_internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpnGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeVpnGateways"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request["VpcId"] = v.(string)
	}

	if v, ok := d.GetOk("status"); ok && v.(string) != "" {
		request["Status"] = strings.ToLower(v.(string))
	}
	if includeReservationData, includeReservationDataOk := d.GetOkExists("include_reservation_data"); includeReservationDataOk {
		request["IncludeReservationData"] = includeReservationData.(bool)
	}

	if v, ok := d.GetOk("business_status"); ok && v.(string) != "" {
		request["BusinessStatus"] = v.(string)
	}

	var objects []map[string]interface{}
	var vpnGatewayNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		vpnGatewayNameRegex = r
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
	request["PageNumber"] = 1
	request["PageSize"] = PageSizeLarge

	for {
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpn_gateways", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.VpnGateways.VpnGateway", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.VpnGateways.VpnGateway", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if vpnGatewayNameRegex != nil && !vpnGatewayNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v", item["VpnGatewayId"])]; !ok {
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

	sslVpn := d.Get("ssl_vpn").(string)
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		if sslVpn != "" && sslVpn != fmt.Sprint(object["SslVpn"]) {
			continue
		}
		mapping := map[string]interface{}{
			"id":                            object["VpnGatewayId"],
			"vpc_id":                        object["VpcId"],
			"internet_ip":                   object["InternetIp"],
			"specification":                 object["Spec"],
			"name":                          object["Name"],
			"description":                   object["Description"],
			"status":                        convertStatus(object["Status"].(string)),
			"business_status":               object["BusinessStatus"],
			"instance_charge_type":          convertChargeType(object["ChargeType"].(string)),
			"enable_ipsec":                  object["IpsecVpn"],
			"enable_ssl":                    object["SslVpn"],
			"ssl_vpn":                       object["SslVpn"],
			"ssl_connections":               object["SslMaxConnections"],
			"network_type":                  object["NetworkType"],
			"disaster_recovery_vswitch_id":  object["DisasterRecoveryVSwitchId"],
			"disaster_recovery_internet_ip": object["DisasterRecoveryInternetIp"],
			"vpn_type":                      object["VpnType"],
			"ssl_vpn_internet_ip":           object["SslVpnInternetIp"],
			"vswitch_id":                    object["VSwitchId"],
			"resource_group_id":             object["ResourceGroupId"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["Key"].(string)
				value := t.(map[string]interface{})["Value"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags

		if v, ok := object["CreateTime"]; ok {
			createTime, err := v.(json.Number).Int64()
			if err != nil {
				log.Println(WrapError(err))
			} else {
				mapping["create_time"] = TimestampToStr(createTime)
			}
		}
		if v, ok := object["EndTime"]; ok {
			endTime, err := v.(json.Number).Int64()
			if err != nil {
				log.Println(WrapError(err))
			} else {
				mapping["end_time"] = TimestampToStr(endTime)
			}
		}
		if v, ok := object["AutoPropagate"]; ok {
			if valueBool, ok := v.(bool); ok {
				mapping["auto_propagate"] = strconv.FormatBool(valueBool)
			}
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

	if err := d.Set("gateways", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

func convertStatus(lower string) string {
	upStr := strings.ToUpper(lower)

	wholeStr := string(upStr[0]) + lower[1:]
	return wholeStr
}

func convertChargeType(originType string) string {
	if string("PostpayByFlow") == originType {
		return string(PostPaid)
	} else {
		return string(PrePaid)
	}
}

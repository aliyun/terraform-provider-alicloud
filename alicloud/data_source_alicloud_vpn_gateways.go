package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpnGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnsRead,

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
						"ssl_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpnsRead(d *schema.ResourceData, meta interface{}) error {
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
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	request["PageNumber"] = 1
	request["PageSize"] = PageSizeLarge

	for {
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		createTime, _ := object["CreateTime"].(json.Number).Int64()
		endTime, _ := object["EndTime"].(json.Number).Int64()
		mapping := map[string]interface{}{
			"id":                   object["VpnGatewayId"],
			"vpc_id":               object["VpcId"],
			"internet_ip":          object["InternetIp"],
			"create_time":          TimestampToStr(createTime),
			"end_time":             TimestampToStr(endTime),
			"specification":        object["Spec"],
			"name":                 object["Name"],
			"description":          object["Description"],
			"status":               convertStatus(object["Status"].(string)),
			"business_status":      object["BusinessStatus"],
			"instance_charge_type": convertChargeType(object["ChargeType"].(string)),
			"enable_ipsec":         object["IpsecVpn"],
			"enable_ssl":           object["SslVpn"],
			"ssl_connections":      object["SslMaxConnections"],
			"network_type":         object["NetworkType"],
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

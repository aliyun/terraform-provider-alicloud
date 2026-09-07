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

func dataSourceAlicloudThreatDetectionSasPrivateLinkEndpoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionSasPrivateLinkEndpointsRead,
		Schema: map[string]*schema.Schema{
			"node_name": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_details": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"endpoints": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"node_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"region_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vpc_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"security_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"update_domain": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"jsrv_domain": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"zones": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"v_switch_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"zone_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionSasPrivateLinkEndpointsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("node_name"); ok {
		request["NodeName"] = v
	}
	request["CurrentPage"] = 1
	request["PageSize"] = PageSizeMedium

	var endpointNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		endpointNameRegex = r
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
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListSasPrivateLinkEndpoint"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_sas_private_link_endpoints", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.List", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.List", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
			}
			if endpointNameRegex != nil && !endpointNameRegex.MatchString(fmt.Sprint(item["NodeName"])) {
				continue
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
	sasService := SasService{client}
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["Id"]),
			"node_name":         object["NodeName"],
			"region_id":         object["RegionId"],
			"vpc_id":            object["VpcId"],
			"security_group_id": object["SecurityGroupId"],
			"status":            object["Status"],
			"update_domain":     object["UpdateDomain"],
			"jsrv_domain":       object["JsrvDomain"],
		}

		ids = append(ids, fmt.Sprint(object["Id"]))

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["Id"])
		obj, err := sasService.DescribeThreatDetectionSasPrivateLinkEndpoint(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["node_name"] = obj["NodeName"]
		mapping["region_id"] = obj["RegionId"]
		mapping["vpc_id"] = obj["VpcId"]
		mapping["security_group_id"] = obj["SecurityGroupId"]
		mapping["status"] = obj["Status"]
		mapping["update_domain"] = obj["UpdateDomain"]
		mapping["jsrv_domain"] = obj["JsrvDomain"]
		if zonesRaw, ok := obj["Zones"]; ok && zonesRaw != nil {
			zonesMaps := make([]map[string]interface{}, 0)
			if zonesList, ok := zonesRaw.([]interface{}); ok {
				for _, value := range zonesList {
					zone := value.(map[string]interface{})
					zoneMap := map[string]interface{}{
						"v_switch_id": zone["VSwitchId"],
						"zone_id":     zone["ZoneId"],
					}
					zonesMaps = append(zonesMaps, zoneMap)
				}
			}
			mapping["zones"] = zonesMaps
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("endpoints", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

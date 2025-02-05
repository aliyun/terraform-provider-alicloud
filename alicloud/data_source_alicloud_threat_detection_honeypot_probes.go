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

func dataSourceAlicloudThreatDetectionHoneypotProbes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionHoneypotProbesRead,
		Schema: map[string]*schema.Schema{
			"display_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"probe_status": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"probe_type": {
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
			"enable_details": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"probes": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"arp": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"control_node_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"display_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"honeypot_bind_list": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bind_port_list": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bind_port": {
													Computed: true,
													Type:     schema.TypeBool,
												},
												"end_port": {
													Computed: true,
													Type:     schema.TypeInt,
												},
												"fixed": {
													Computed: true,
													Type:     schema.TypeBool,
												},
												"start_port": {
													Computed: true,
													Type:     schema.TypeInt,
												},
												"target_port": {
													Computed: true,
													Type:     schema.TypeInt,
												},
											},
										},
									},
									"honeypot_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"honeypot_probe_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"ping": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"probe_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"uuid": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vpc_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"service_ip_list": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionHoneypotProbesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("display_name"); ok {
		request["DisplayName"] = v
	}
	if v, ok := d.GetOk("probe_status"); ok {
		request["ProbeStatus"] = v
	}
	if v, ok := d.GetOk("probe_type"); ok {
		request["ProbeType"] = v
	}
	request["CurrentPage"] = 1
	request["PageSize"] = PageSizeMedium

	var probeNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		probeNameRegex = r
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
		action := "ListHoneypotProbe"
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_honeypot_probes", action, AlibabaCloudSdkGoERROR)
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
				if _, ok := idsMap[fmt.Sprint(item["ProbeId"])]; !ok {
					continue
				}
			}
			if probeNameRegex != nil && !probeNameRegex.MatchString(fmt.Sprint(item["DisplayName"])) {
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
			"id":                fmt.Sprint(object["ProbeId"]),
			"control_node_id":   object["ControlNode.NodeId"],
			"display_name":      object["DisplayName"],
			"honeypot_probe_id": object["ProbeId"],
			"probe_type":        object["ProbeType"],
			"status":            object["Status"],
			"uuid":              object["Uuid"],
			"vpc_id":            object["VpcId"],
		}

		ids = append(ids, fmt.Sprint(object["ProbeId"]))

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["ProbeId"])
		object, err = sasService.DescribeThreatDetectionHoneypotProbe(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["arp"] = object["Arp"]
		controlNodeId, _ := jsonpath.Get("$.ControlNode.NodeId", object)
		mapping["control_node_id"] = controlNodeId
		mapping["display_name"] = object["DisplayName"]
		honeypotBindList67Maps := make([]map[string]interface{}, 0)
		honeypotBindList67Raw := object["HoneypotProbeBindList"]
		for _, value0 := range honeypotBindList67Raw.([]interface{}) {
			honeypotBindList67 := value0.(map[string]interface{})
			honeypotBindList67Map := make(map[string]interface{})
			bindPortList67Maps := make([]map[string]interface{}, 0)
			bindPortList67Raw := honeypotBindList67["BindPortList"]
			for _, value1 := range bindPortList67Raw.([]interface{}) {
				bindPortList67 := value1.(map[string]interface{})
				bindPortList67Map := make(map[string]interface{})
				bindPortList67Map["bind_port"] = bindPortList67["BindPort"]
				bindPortList67Map["end_port"] = bindPortList67["EndPort"]
				bindPortList67Map["fixed"] = bindPortList67["Fixed"]
				bindPortList67Map["start_port"] = bindPortList67["StartPort"]
				bindPortList67Map["target_port"] = bindPortList67["TargetPort"]
				bindPortList67Maps = append(bindPortList67Maps, bindPortList67Map)
			}
			honeypotBindList67Map["bind_port_list"] = bindPortList67Maps
			honeypotBindList67Map["honeypot_id"] = honeypotBindList67["HoneypotId"]
			honeypotBindList67Maps = append(honeypotBindList67Maps, honeypotBindList67Map)
		}
		mapping["honeypot_bind_list"] = honeypotBindList67Maps
		mapping["honeypot_probe_id"] = object["ProbeId"]
		mapping["ping"] = object["Ping"]
		mapping["probe_type"] = object["ProbeType"]
		serviceIpList, _ := jsonpath.Get("$.ListenIpList", object)
		mapping["service_ip_list"] = serviceIpList
		mapping["status"] = object["Status"]
		mapping["uuid"] = object["Uuid"]
		mapping["vpc_id"] = object["VpcId"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("probes", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

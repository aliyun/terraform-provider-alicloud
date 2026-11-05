package ga

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudGaListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaListenerCreate,
		Read:   resourceAliCloudGaListenerRead,
		Update: resourceAliCloudGaListenerUpdate,
		Delete: resourceAliCloudGaListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "TCP",
				ValidateFunc: validators.StringInSlice([]string{"TCP", "UDP", "HTTP", "HTTPS"}, false),
			},
			"proxy_protocol": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"security_policy_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"tls_cipher_policy_1_0", "tls_cipher_policy_1_1", "tls_cipher_policy_1_2", "tls_cipher_policy_1_2_strict", "tls_cipher_policy_1_2_strict_with_1_3"}, false),
			},
			"listener_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"Standard", "CustomRouting"}, false),
			},
			"http_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"http1.1", "http2", "http3"}, false),
			},
			"idle_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"request_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"client_affinity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"NONE", "SOURCE_IP"}, false),
				Default:      "NONE",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("protocol") == "UDP" && new == "SOURCE_IP" {
						return true
					}
					return false
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"port_ranges": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"to_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"forwarded_for_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"forwarded_for_ga_id_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"forwarded_for_ga_ap_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"forwarded_for_proto_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"forwarded_for_port_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"real_ip_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"backend_ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"to_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"endpoint_group_configurations": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_group_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringLenBetween(2, 128),
						},
						"endpoint_group_description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringLenBetween(0, 256),
						},
						"endpoint_group_region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"traffic_percentage": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validators.IntBetween(0, 100),
						},
						"health_check_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"health_check_interval_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validators.IntBetween(1, 60),
						},
						"health_check_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"health_check_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validators.IntBetween(1, 65535),
						},
						"health_check_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringInSlice([]string{"TCP", "HTTP", "HTTPS", "tcp", "http", "https"}, false),
						},
						"threshold_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validators.IntBetween(2, 10),
						},
						"endpoint_configurations": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validators.StringInSlice([]string{"Domain", "Ip", "IpTarget", "PublicIp", "ECS", "SLB", "ALB", "NLB", "ENI", "OSS"}, false),
									},
									"weight": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validators.IntBetween(0, 255),
									},
									"endpoint": {
										Type:     schema.TypeString,
										Required: true,
									},
									"api_keys": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"vswitch_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"sub_address": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enable_proxy_protocol": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"provider": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enable_client_ip_preservation": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"endpoint_request_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringInSlice([]string{"HTTP", "HTTPS"}, false),
						},
						"endpoint_group_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringInSlice([]string{"default", "virtual"}, false),
						},
						"port_overrides": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"endpoint_port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"enable_client_ip_preservation_proxy_protocol": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_client_ip_preservation_toa": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"endpoint_protocol_version": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringInSlice([]string{"HTTP1.1", "HTTP2"}, false),
						},
						"endpoint_ip_version": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringInSlice([]string{"IPv4", "IPv6", "ProtocolAffinity"}, false),
						},
						"health_check_host": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"custom_routing_endpoint_group_configurations": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_group_region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringLenBetween(2, 128),
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validators.StringLenBetween(0, 256),
						},
						"destination_configurations": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocols": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"from_port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"to_port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"custom_routing_endpoint_configurations": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"PrivateSubNet"}, false),
									},
									"endpoint": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"traffic_to_endpoint_policy": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validators.StringInSlice([]string{"AllowAll", "DenyAll", "AllowCustom"}, false),
									},
									"policy_configurations": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"custom_routing_port_ranges": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"from_port": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"to_port": {
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var response map[string]interface{}
	action := "CreateListener"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("CreateListener")
	request["AcceleratorId"] = d.Get("accelerator_id")

	if v, ok := d.GetOk("protocol"); ok {
		request["Protocol"] = v
	}

	if v, ok := d.GetOkExists("proxy_protocol"); ok {
		request["ProxyProtocol"] = v
	}

	if v, ok := d.GetOk("security_policy_id"); ok {
		request["SecurityPolicyId"] = v
	}

	if v, ok := d.GetOk("listener_type"); ok {
		request["Type"] = v
	}

	if v, ok := d.GetOk("http_version"); ok {
		request["HttpVersion"] = v
	}

	if v, ok := d.GetOkExists("idle_timeout"); ok {
		request["IdleTimeout"] = v
	}

	if v, ok := d.GetOkExists("request_timeout"); ok {
		request["RequestTimeout"] = v
	}

	if v, ok := d.GetOk("client_affinity"); ok {
		request["ClientAffinity"] = v
	}

	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("certificates"); ok {
		Certificates := make([]map[string]interface{}, len(v.([]interface{})))
		for i, CertificatesValue := range v.([]interface{}) {
			CertificatesMap := CertificatesValue.(map[string]interface{})
			Certificates[i] = make(map[string]interface{})
			Certificates[i]["Id"] = CertificatesMap["id"]
		}
		request["Certificates"] = Certificates

	}

	PortRanges := make([]map[string]interface{}, len(d.Get("port_ranges").([]interface{})))
	for i, PortRangesValue := range d.Get("port_ranges").([]interface{}) {
		PortRangesMap := PortRangesValue.(map[string]interface{})
		PortRanges[i] = make(map[string]interface{})
		PortRanges[i]["FromPort"] = PortRangesMap["from_port"]
		PortRanges[i]["ToPort"] = PortRangesMap["to_port"]
	}
	request["PortRanges"] = PortRanges

	if v, ok := d.GetOk("forwarded_for_config"); ok {
		forwardedForConfigMap := map[string]interface{}{}
		for _, forwardedForConfigList := range v.([]interface{}) {
			forwardedForConfigArg := forwardedForConfigList.(map[string]interface{})

			forwardedForConfigMap["XForwardedForGaIdEnabled"] = forwardedForConfigArg["forwarded_for_ga_id_enabled"]
			forwardedForConfigMap["XForwardedForGaApEnabled"] = forwardedForConfigArg["forwarded_for_ga_ap_enabled"]
			forwardedForConfigMap["XForwardedForProtoEnabled"] = forwardedForConfigArg["forwarded_for_proto_enabled"]
			forwardedForConfigMap["XForwardedForPortEnabled"] = forwardedForConfigArg["forwarded_for_port_enabled"]
			forwardedForConfigMap["XRealIpEnabled"] = forwardedForConfigArg["real_ip_enabled"]
		}

		request["XForwardedForConfig"] = forwardedForConfigMap
	}

	if v, ok := d.GetOk("backend_ports"); ok {
		request["BackendPorts"] = buildGaListenerBackendPorts(v)
	}

	if v, ok := d.GetOk("endpoint_group_configurations"); ok {
		request["EndpointGroupConfigurations"] = buildGaListenerEndpointGroupConfigurations(v)
	}

	if v, ok := d.GetOk("custom_routing_endpoint_group_configurations"); ok {
		request["CustomRoutingEndpointGroupConfigurations"] = buildGaListenerCustomRoutingEndpointGroupConfigurations(v)
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "NotExist.BasicBandwidthPackage"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_listener", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ListenerId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaListenerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudGaListenerRead(d, meta)
}

func resourceAliCloudGaListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}

	object, err := gaService.DescribeGaListener(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_listener gaService.DescribeGaListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("protocol", object["Protocol"])
	d.Set("proxy_protocol", object["ProxyProtocol"])
	d.Set("security_policy_id", object["SecurityPolicyId"])
	d.Set("listener_type", object["Type"])
	d.Set("http_version", object["HttpVersion"])
	d.Set("idle_timeout", object["IdleTimeout"])
	d.Set("request_timeout", object["RequestTimeout"])
	d.Set("client_affinity", object["ClientAffinity"])
	d.Set("name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("status", object["State"])

	certificates := make([]map[string]interface{}, 0)
	if certificatesList, ok := object["Certificates"].([]interface{}); ok {
		for _, v := range certificatesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"id": m1["Id"],
				}
				certificates = append(certificates, temp1)

			}
		}
	}
	if err := d.Set("certificates", certificates); err != nil {
		return errs.WrapError(err)
	}

	portRanges := make([]map[string]interface{}, 0)
	if portRangesList, ok := object["PortRanges"].([]interface{}); ok {
		for _, v := range portRangesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"from_port": m1["FromPort"],
					"to_port":   m1["ToPort"],
				}
				portRanges = append(portRanges, temp1)

			}
		}
	}
	if err := d.Set("port_ranges", portRanges); err != nil {
		return errs.WrapError(err)
	}

	if backendPortsList, ok := object["BackendPorts"].([]interface{}); ok {
		backendPorts := make([]map[string]interface{}, 0)
		for _, v := range backendPortsList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"from_port": m1["FromPort"],
					"to_port":   m1["ToPort"],
				}
				backendPorts = append(backendPorts, temp1)
			}
		}
		if err := d.Set("backend_ports", backendPorts); err != nil {
			return errs.WrapError(err)
		}
	}

	if forwardedForConfig, ok := object["XForwardedForConfig"]; ok {
		forwardedForConfigMaps := make([]map[string]interface{}, 0)
		forwardedForConfigArg := forwardedForConfig.(map[string]interface{})
		forwardedForConfigMap := map[string]interface{}{}

		if forwardedForGaIdEnabled, ok := forwardedForConfigArg["XForwardedForGaIdEnabled"]; ok {
			forwardedForConfigMap["forwarded_for_ga_id_enabled"] = forwardedForGaIdEnabled
		}

		if forwardedForGaApEnabled, ok := forwardedForConfigArg["XForwardedForGaApEnabled"]; ok {
			forwardedForConfigMap["forwarded_for_ga_ap_enabled"] = forwardedForGaApEnabled
		}

		if forwardedForProtoEnabled, ok := forwardedForConfigArg["XForwardedForProtoEnabled"]; ok {
			forwardedForConfigMap["forwarded_for_proto_enabled"] = forwardedForProtoEnabled
		}

		if forwardedForPortEnabled, ok := forwardedForConfigArg["XForwardedForPortEnabled"]; ok {
			forwardedForConfigMap["forwarded_for_port_enabled"] = forwardedForPortEnabled
		}

		if realIpEnabled, ok := forwardedForConfigArg["XRealIpEnabled"]; ok {
			forwardedForConfigMap["real_ip_enabled"] = realIpEnabled
		}

		forwardedForConfigMaps = append(forwardedForConfigMaps, forwardedForConfigMap)

		d.Set("forwarded_for_config", forwardedForConfigMaps)
	}

	return nil
}

func resourceAliCloudGaListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken("UpdateListener"),
		"ListenerId":  d.Id(),
	}

	if d.HasChange("certificates") {
		update = true
		if v, ok := d.GetOk("certificates"); ok {
			Certificates := make([]map[string]interface{}, len(v.([]interface{})))
			for i, CertificatesValue := range v.([]interface{}) {
				CertificatesMap := CertificatesValue.(map[string]interface{})
				Certificates[i] = make(map[string]interface{})
				Certificates[i]["Id"] = CertificatesMap["id"]
			}
			request["Certificates"] = Certificates
		}

	}

	if d.HasChange("idle_timeout") {
		update = true

		if v, ok := d.GetOkExists("idle_timeout"); ok {
			request["IdleTimeout"] = v
		}
	}

	if d.HasChange("request_timeout") {
		update = true

		if v, ok := d.GetOkExists("request_timeout"); ok {
			request["RequestTimeout"] = v
		}
	}

	if d.HasChange("client_affinity") {
		update = true
		request["ClientAffinity"] = d.Get("client_affinity")
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if d.HasChange("port_ranges") {
		update = true
		PortRanges := make([]map[string]interface{}, len(d.Get("port_ranges").([]interface{})))
		for i, PortRangesValue := range d.Get("port_ranges").([]interface{}) {
			PortRangesMap := PortRangesValue.(map[string]interface{})
			PortRanges[i] = make(map[string]interface{})
			PortRanges[i]["FromPort"] = PortRangesMap["from_port"]
			PortRanges[i]["ToPort"] = PortRangesMap["to_port"]
		}
		request["PortRanges"] = PortRanges

	}

	if d.HasChange("protocol") {
		update = true
		request["Protocol"] = d.Get("protocol")
	}

	if d.HasChange("security_policy_id") {
		update = true
		request["SecurityPolicyId"] = d.Get("security_policy_id")
	}

	if d.HasChange("proxy_protocol") {
		update = true
	}
	if v, ok := d.GetOkExists("proxy_protocol"); ok {
		request["ProxyProtocol"] = v
	}

	if d.HasChange("http_version") {
		update = true

		if v, ok := d.GetOk("http_version"); ok {
			request["HttpVersion"] = v
		}
	}

	if d.HasChange("forwarded_for_config") {
		update = true
	}
	if v, ok := d.GetOk("forwarded_for_config"); ok {
		forwardedForConfigMap := map[string]interface{}{}
		for _, forwardedForConfigList := range v.([]interface{}) {
			forwardedForConfigArg := forwardedForConfigList.(map[string]interface{})

			forwardedForConfigMap["XForwardedForGaIdEnabled"] = forwardedForConfigArg["forwarded_for_ga_id_enabled"]
			forwardedForConfigMap["XForwardedForGaApEnabled"] = forwardedForConfigArg["forwarded_for_ga_ap_enabled"]
			forwardedForConfigMap["XForwardedForProtoEnabled"] = forwardedForConfigArg["forwarded_for_proto_enabled"]
			forwardedForConfigMap["XForwardedForPortEnabled"] = forwardedForConfigArg["forwarded_for_port_enabled"]
			forwardedForConfigMap["XRealIpEnabled"] = forwardedForConfigArg["real_ip_enabled"]
		}

		request["XForwardedForConfig"] = forwardedForConfigMap
	}

	if d.HasChange("backend_ports") {
		update = true
		if v, ok := d.GetOk("backend_ports"); ok {
			request["BackendPorts"] = buildGaListenerBackendPorts(v)
		}
	}

	if update {
		action := "UpdateListener"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "NotActive.Listener"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaListenerStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudGaListenerRead(d, meta)
}

func resourceAliCloudGaListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}
	action := "DeleteListener"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"ClientToken":   helper.BuildClientToken("DeleteListener"),
		"AcceleratorId": d.Get("accelerator_id"),
		"ListenerId":    d.Id(),
	}

	// Endpoint groups created inline through CreateListener are not managed by
	// alicloud_ga_endpoint_group, so they must be removed before the listener.
	// Only do this when this resource actually created them; endpoint groups
	// managed by separate resources keep the previous failure behavior.
	if len(d.Get("endpoint_group_configurations").(*schema.Set).List()) > 0 || len(d.Get("custom_routing_endpoint_group_configurations").(*schema.Set).List()) > 0 {
		if err := deleteGaListenerInlineEndpointGroups(client, gaService, fmt.Sprint(d.Get("accelerator_id")), d.Id()); err != nil {
			return errs.WrapError(err)
		}
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "NotActive.Listener", "Exist.ForwardingRule", "Exist.EndpointGroup"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaListenerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

func gaListenerNestedList(v interface{}) []interface{} {
	switch value := v.(type) {
	case *schema.Set:
		return value.List()
	case []interface{}:
		return value
	default:
		return nil
	}
}

func gaListenerStringSlice(v interface{}) []string {
	items := gaListenerNestedList(v)
	result := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok && s != "" {
			result = append(result, s)
		}
	}
	return result
}

// setGaListenerParam omits unset scalar values so that the service keeps its
// documented defaults instead of receiving explicit zero values.
func setGaListenerParam(target map[string]interface{}, key string, value interface{}) {
	switch v := value.(type) {
	case nil:
		return
	case string:
		if v != "" {
			target[key] = v
		}
	case bool:
		target[key] = v
	case int:
		if v != 0 {
			target[key] = v
		}
	case []string:
		if len(v) > 0 {
			target[key] = v
		}
	case []map[string]interface{}:
		if len(v) > 0 {
			target[key] = v
		}
	}
}

func buildGaListenerBackendPorts(v interface{}) []map[string]interface{} {
	items := gaListenerNestedList(v)
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"FromPort": m["from_port"],
			"ToPort":   m["to_port"],
		})
	}
	return result
}

func buildGaListenerEndpointConfigurations(v interface{}) []map[string]interface{} {
	items := gaListenerNestedList(v)
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		endpoint := map[string]interface{}{
			"Type":     m["type"],
			"Weight":   m["weight"],
			"Endpoint": m["endpoint"],
		}
		setGaListenerParam(endpoint, "ApiKeys", gaListenerStringSlice(m["api_keys"]))
		setGaListenerParam(endpoint, "VpcId", m["vpc_id"])
		setGaListenerParam(endpoint, "VSwitchIds", gaListenerStringSlice(m["vswitch_ids"]))
		setGaListenerParam(endpoint, "SubAddress", m["sub_address"])
		setGaListenerParam(endpoint, "EnableProxyProtocol", m["enable_proxy_protocol"])
		setGaListenerParam(endpoint, "Provider", m["provider"])
		setGaListenerParam(endpoint, "EnableClientIPPreservation", m["enable_client_ip_preservation"])
		result = append(result, endpoint)
	}
	return result
}

func buildGaListenerPortOverrides(v interface{}) []map[string]interface{} {
	items := gaListenerNestedList(v)
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"ListenerPort": m["listener_port"],
			"EndpointPort": m["endpoint_port"],
		})
	}
	return result
}

func buildGaListenerEndpointGroupConfigurations(v interface{}) []map[string]interface{} {
	items := gaListenerNestedList(v)
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		config := map[string]interface{}{
			"EndpointGroupRegion": m["endpoint_group_region"],
		}
		setGaListenerParam(config, "EndpointGroupName", m["endpoint_group_name"])
		setGaListenerParam(config, "EndpointGroupDescription", m["endpoint_group_description"])
		setGaListenerParam(config, "TrafficPercentage", m["traffic_percentage"])
		setGaListenerParam(config, "HealthCheckEnabled", m["health_check_enabled"])
		setGaListenerParam(config, "HealthCheckIntervalSeconds", m["health_check_interval_seconds"])
		setGaListenerParam(config, "HealthCheckPath", m["health_check_path"])
		setGaListenerParam(config, "HealthCheckPort", m["health_check_port"])
		setGaListenerParam(config, "HealthCheckProtocol", m["health_check_protocol"])
		setGaListenerParam(config, "ThresholdCount", m["threshold_count"])
		setGaListenerParam(config, "EndpointConfigurations", buildGaListenerEndpointConfigurations(m["endpoint_configurations"]))
		setGaListenerParam(config, "EndpointRequestProtocol", m["endpoint_request_protocol"])
		setGaListenerParam(config, "EndpointGroupType", m["endpoint_group_type"])
		setGaListenerParam(config, "PortOverrides", buildGaListenerPortOverrides(m["port_overrides"]))
		setGaListenerParam(config, "EnableClientIpPreservationProxyProtocol", m["enable_client_ip_preservation_proxy_protocol"])
		setGaListenerParam(config, "EnableClientIpPreservationToa", m["enable_client_ip_preservation_toa"])
		setGaListenerParam(config, "EndpointProtocolVersion", m["endpoint_protocol_version"])
		setGaListenerParam(config, "EndpointIpVersion", m["endpoint_ip_version"])
		setGaListenerParam(config, "HealthCheckHost", m["health_check_host"])
		result = append(result, config)
	}
	return result
}

func buildGaListenerCustomRoutingPortRanges(v interface{}) []map[string]interface{} {
	items := gaListenerNestedList(v)
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"FromPort": m["from_port"],
			"ToPort":   m["to_port"],
		})
	}
	return result
}

func buildGaListenerPolicyConfigurations(v interface{}) []map[string]interface{} {
	items := gaListenerNestedList(v)
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		policy := map[string]interface{}{}
		setGaListenerParam(policy, "Address", m["address"])
		setGaListenerParam(policy, "PortRanges", buildGaListenerCustomRoutingPortRanges(m["custom_routing_port_ranges"]))
		result = append(result, policy)
	}
	return result
}

func buildGaListenerCustomRoutingEndpointConfigurations(v interface{}) []map[string]interface{} {
	items := gaListenerNestedList(v)
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		endpoint := map[string]interface{}{}
		setGaListenerParam(endpoint, "Type", m["type"])
		setGaListenerParam(endpoint, "Endpoint", m["endpoint"])
		setGaListenerParam(endpoint, "TrafficToEndpointPolicy", m["traffic_to_endpoint_policy"])
		setGaListenerParam(endpoint, "PolicyConfigurations", buildGaListenerPolicyConfigurations(m["policy_configurations"]))
		result = append(result, endpoint)
	}
	return result
}

func buildGaListenerDestinationConfigurations(v interface{}) []map[string]interface{} {
	items := gaListenerNestedList(v)
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		destination := map[string]interface{}{}
		setGaListenerParam(destination, "Protocols", gaListenerStringSlice(m["protocols"]))
		setGaListenerParam(destination, "FromPort", m["from_port"])
		setGaListenerParam(destination, "ToPort", m["to_port"])
		result = append(result, destination)
	}
	return result
}

func buildGaListenerCustomRoutingEndpointGroupConfigurations(v interface{}) []map[string]interface{} {
	items := gaListenerNestedList(v)
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		config := map[string]interface{}{
			"EndpointGroupRegion": m["endpoint_group_region"],
		}
		setGaListenerParam(config, "Name", m["name"])
		setGaListenerParam(config, "Description", m["description"])
		setGaListenerParam(config, "DestinationConfigurations", buildGaListenerDestinationConfigurations(m["destination_configurations"]))
		setGaListenerParam(config, "CustomRoutingEndpointConfigurations", buildGaListenerCustomRoutingEndpointConfigurations(m["custom_routing_endpoint_configurations"]))
		result = append(result, config)
	}
	return result
}

func listGaListenerEndpointGroupIds(client *connectivity.AliyunClient, action, acceleratorId, listenerId string) ([]string, error) {
	ids := make([]string, 0)
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"AcceleratorId": acceleratorId,
		"ListenerId":    listenerId,
		"PageSize":      helper.PageSizeLarge,
		"PageNumber":    1,
	}
	for {
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		logs.AddDebug(action, response, request)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotExist.Listener", "NotExist.Accelerator"}) {
				return ids, nil
			}
			return ids, err
		}
		items, _ := response["EndpointGroups"].([]interface{})
		for _, item := range items {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if id := fmt.Sprint(m["EndpointGroupId"]); id != "" && id != "<nil>" {
				ids = append(ids, id)
			}
		}
		if len(items) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return ids, nil
}

func deleteGaListenerStandardEndpointGroup(client *connectivity.AliyunClient, gaService Service, acceleratorId, endpointGroupId string) error {
	action := "DeleteEndpointGroup"
	request := map[string]interface{}{
		"ClientToken":     helper.BuildClientToken("DeleteEndpointGroup"),
		"AcceleratorId":   acceleratorId,
		"EndpointGroupId": endpointGroupId,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(client.GetRetryTimeout(10*time.Minute), func() *retry.RetryError {
		_, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
				return nil
			}
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup", "NotActive.Listener"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, endpointGroupId, action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{}, 10*time.Minute, 5*time.Second, gaService.GaEndpointGroupStateRefreshFunc(endpointGroupId, []string{}))
	_, err = stateConf.WaitForState()
	return err
}

func deleteGaListenerCustomRoutingEndpointGroup(client *connectivity.AliyunClient, gaService Service, acceleratorId, endpointGroupId string) error {
	action := "DeleteCustomRoutingEndpointGroups"
	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"ClientToken":      helper.BuildClientToken("DeleteCustomRoutingEndpointGroups"),
		"AcceleratorId":    acceleratorId,
		"EndpointGroupIds": []string{endpointGroupId},
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(client.GetRetryTimeout(10*time.Minute), func() *retry.RetryError {
		_, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
				return nil
			}
			if errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "NotActive.Listener"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, endpointGroupId, action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{}, 10*time.Minute, 5*time.Second, gaService.GaCustomRoutingEndpointGroupStateRefreshFunc(endpointGroupId, []string{}))
	_, err = stateConf.WaitForState()
	return err
}

// deleteGaListenerInlineEndpointGroups removes the endpoint groups that were
// created together with the listener. DeleteListener rejects the request with
// Exist.EndpointGroup while any of them still exist.
func deleteGaListenerInlineEndpointGroups(client *connectivity.AliyunClient, gaService Service, acceleratorId, listenerId string) error {
	standardIds, err := listGaListenerEndpointGroupIds(client, "ListEndpointGroups", acceleratorId, listenerId)
	if err != nil {
		return errs.WrapError(err)
	}
	for _, endpointGroupId := range standardIds {
		if err := deleteGaListenerStandardEndpointGroup(client, gaService, acceleratorId, endpointGroupId); err != nil {
			return errs.WrapError(err)
		}
	}

	customRoutingIds, err := listGaListenerEndpointGroupIds(client, "ListCustomRoutingEndpointGroups", acceleratorId, listenerId)
	if err != nil {
		return errs.WrapError(err)
	}
	for _, endpointGroupId := range customRoutingIds {
		if err := deleteGaListenerCustomRoutingEndpointGroup(client, gaService, acceleratorId, endpointGroupId); err != nil {
			return errs.WrapError(err)
		}
	}

	return nil
}

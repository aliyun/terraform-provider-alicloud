package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGaListener() *schema.Resource {
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
				ValidateFunc: StringInSlice([]string{"TCP", "UDP", "HTTP", "HTTPS"}, false),
			},
			"proxy_protocol": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"security_policy_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"tls_cipher_policy_1_0", "tls_cipher_policy_1_1", "tls_cipher_policy_1_2", "tls_cipher_policy_1_2_strict", "tls_cipher_policy_1_2_strict_with_1_3"}, false),
			},
			"listener_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Standard", "CustomRouting"}, false),
			},
			"http_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"http1.1", "http2", "http3"}, false),
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
				ValidateFunc: StringInSlice([]string{"NONE", "SOURCE_IP"}, false),
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGaListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateListener"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateListener")
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "NotExist.BasicBandwidthPackage"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_listener", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ListenerId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaListenerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaListenerRead(d, meta)
}

func resourceAliCloudGaListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaListener(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_listener gaService.DescribeGaListener Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		return WrapError(err)
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
		return WrapError(err)
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
	gaService := GaService{client}
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("UpdateListener"),
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

	if update {
		action := "UpdateListener"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator", "NotActive.Listener"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaListenerStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudGaListenerRead(d, meta)
}

func resourceAliCloudGaListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteListener"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"ClientToken":   buildClientToken("DeleteListener"),
		"AcceleratorId": d.Get("accelerator_id"),
		"ListenerId":    d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "NotActive.Listener", "Exist.ForwardingRule", "Exist.EndpointGroup"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaListenerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

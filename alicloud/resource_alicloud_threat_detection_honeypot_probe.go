package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudThreatDetectionHoneypotProbe() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionHoneypotProbeCreate,
		Read:   resourceAlicloudThreatDetectionHoneypotProbeRead,
		Update: resourceAlicloudThreatDetectionHoneypotProbeUpdate,
		Delete: resourceAlicloudThreatDetectionHoneypotProbeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"arp": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"control_node_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"display_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"honeypot_bind_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_port_list": {
							Optional: true,
							Type:     schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bind_port": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeBool,
									},
									"end_port": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeInt,
									},
									"fixed": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeBool,
									},
									"start_port": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeInt,
									},
									"target_port": {
										Optional: true,
										Computed: true,
										Type:     schema.TypeInt,
									},
								},
							},
						},
						"honeypot_id": {
							Optional: true,
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
				Optional: true,
				Type:     schema.TypeBool,
			},
			"probe_type": {
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"host_probe", "vpc_black_hole_probe"}, false),
				Type:         schema.TypeString,
			},
			"probe_version": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"service_ip_list": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"uuid": {
				Optional:      true,
				ForceNew:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"vpc_id"},
			},
			"vpc_id": {
				Optional:      true,
				ForceNew:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"uuid"},
			},
			"proxy_ip": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudThreatDetectionHoneypotProbeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOkExists("arp"); ok {
		request["Arp"] = v
	}
	if v, ok := d.GetOk("control_node_id"); ok {
		request["ControlNodeId"] = v
	}
	if v, ok := d.GetOk("display_name"); ok {
		request["DisplayName"] = v
	}
	if v, ok := d.GetOk("proxy_ip"); ok {
		request["ProxyIp"] = v
	}
	if v, ok := d.GetOk("honeypot_bind_list"); ok {
		honeypotBindListMaps := make([]map[string]interface{}, 0)
		for _, value0 := range v.(*schema.Set).List() {
			honeypotBindList := value0.(map[string]interface{})
			honeypotBindListMap := make(map[string]interface{})
			bindPortListMaps := make([]map[string]interface{}, 0)
			for _, value1 := range honeypotBindList["bind_port_list"].(*schema.Set).List() {
				bindPortList := value1.(map[string]interface{})
				bindPortListMap := make(map[string]interface{})
				bindPortListMap["BindPort"] = bindPortList["bind_port"]
				bindPortListMap["EndPort"] = bindPortList["end_port"]
				bindPortListMap["Fixed"] = bindPortList["fixed"]
				bindPortListMap["StartPort"] = bindPortList["start_port"]
				bindPortListMap["TargetPort"] = bindPortList["target_port"]
				bindPortListMaps = append(bindPortListMaps, bindPortListMap)
			}
			honeypotBindListMap["BindPortList"] = bindPortListMaps
			honeypotBindListMap["HoneypotId"] = honeypotBindList["honeypot_id"]
			honeypotBindListMaps = append(honeypotBindListMaps, honeypotBindListMap)
		}
		request["HoneypotBindList"] = honeypotBindListMaps
	}
	if v, ok := d.GetOkExists("ping"); ok {
		request["Ping"] = v
	}
	if v, ok := d.GetOk("probe_type"); ok {
		request["ProbeType"] = v
	}
	if v, ok := d.GetOk("probe_version"); ok {
		request["ProbeVersion"] = v
	}
	if v, ok := d.GetOk("uuid"); ok {
		request["Uuid"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	var response map[string]interface{}
	action := "CreateHoneypotProbe"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_honeypot_probe", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.HoneypotProbe.ProbeId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_honeypot_probe")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"online", "offline", "unnormal", "unprobe", "installed"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, sasService.ThreatDetectionHoneypotProbeStateRefreshFunc(d, []string{"install_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudThreatDetectionHoneypotProbeUpdate(d, meta)
}

func resourceAlicloudThreatDetectionHoneypotProbeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}

	object, err := sasService.DescribeThreatDetectionHoneypotProbe(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_honeypot_probe sasService.DescribeThreatDetectionHoneypotProbe Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("arp", object["Arp"])
	controlNodeId, _ := jsonpath.Get("$.ControlNode.NodeId", object)
	d.Set("control_node_id", controlNodeId)
	d.Set("display_name", object["DisplayName"])
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
	d.Set("honeypot_bind_list", honeypotBindList67Maps)
	d.Set("probe_version", object["ProbeVersion"])
	d.Set("honeypot_probe_id", object["ProbeId"])
	d.Set("ping", object["Ping"])
	d.Set("probe_type", object["ProbeType"])
	serviceIpList, _ := jsonpath.Get("$.ListenIpList", object)
	d.Set("service_ip_list", serviceIpList)
	d.Set("status", object["Status"])
	d.Set("uuid", object["Uuid"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("proxy_ip", object["ProxyIp"])

	return nil
}

func resourceAlicloudThreatDetectionHoneypotProbeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{
		"ProbeId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("arp") {
		update = true
		if v, ok := d.GetOkExists("arp"); ok {
			request["Arp"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("display_name") {
		update = true
		if v, ok := d.GetOk("display_name"); ok {
			request["DisplayName"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("ping") {
		update = true
		if v, ok := d.GetOkExists("ping"); ok {
			request["Ping"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("service_ip_list") {
		update = true
		if v, ok := d.GetOk("service_ip_list"); ok {
			request["ServiceIpList"] = v
		}
	}

	if update {
		action := "UpdateHoneypotProbe"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudThreatDetectionHoneypotProbeRead(d, meta)
}

func resourceAlicloudThreatDetectionHoneypotProbeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}
	var err error

	request := map[string]interface{}{
		"ProbeId": d.Id(),
	}

	action := "DeleteHoneypotProbe"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, sasService.ThreatDetectionHoneypotProbeStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

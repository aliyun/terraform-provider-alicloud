package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudThreatDetectionHoneypotNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionHoneypotNodeCreate,
		Read:   resourceAlicloudThreatDetectionHoneypotNodeRead,
		Update: resourceAlicloudThreatDetectionHoneypotNodeUpdate,
		Delete: resourceAlicloudThreatDetectionHoneypotNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allow_honeypot_access_internet": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeBool,
			},
			"available_probe_num": {
				Required: true,
				Type:     schema.TypeInt,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"node_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"security_group_probe_ip_list": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Computed: true,
				Type:     schema.TypeInt,
			},
		},
	}
}

func resourceAlicloudThreatDetectionHoneypotNodeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("allow_honeypot_access_internet"); ok {
		request["AllowHoneypotAccessInternet"] = v
	}
	request["AvailableProbeNum"] = d.Get("available_probe_num")
	request["NodeName"] = d.Get("node_name")
	if v, ok := d.GetOk("security_group_probe_ip_list"); ok {
		request["SecurityGroupProbeIpList"] = v.([]interface{})
	}

	var response map[string]interface{}
	action := "CreateHoneypotNode"
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_honeypot_node", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.HoneypotNode.NodeId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_honeypot_node")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, threatDetectionService.ThreatDetectionHoneypotNodeStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudThreatDetectionHoneypotNodeRead(d, meta)
}

func resourceAlicloudThreatDetectionHoneypotNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}

	object, err := threatDetectionService.DescribeThreatDetectionHoneypotNode(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_honeypot_node threatDetectionService.DescribeThreatDetectionHoneypotNode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("allow_honeypot_access_internet", object["AllowHoneypotAccessInternet"])
	d.Set("available_probe_num", object["ProbeTotalCount"])
	d.Set("create_time", object["CreateTime"])
	d.Set("node_name", object["NodeName"])
	securityGroupProbeIpList, _ := jsonpath.Get("$.SecurityGroupProbeIpList", object)
	d.Set("security_group_probe_ip_list", securityGroupProbeIpList)
	d.Set("status", object["TotalStatus"])

	return nil
}

func resourceAlicloudThreatDetectionHoneypotNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	threatDetectionService := ThreatDetectionService{client}
	update := false
	request := map[string]interface{}{
		"NodeId": d.Id(),
	}

	if d.HasChange("available_probe_num") {
		update = true
		request["AvailableProbeNum"] = d.Get("available_probe_num")
	}
	if d.HasChange("node_name") {
		update = true
		request["NodeName"] = d.Get("node_name")
	}
	if d.HasChange("security_group_probe_ip_list") {
		update = true
		if v, ok := d.GetOk("security_group_probe_ip_list"); ok {
			request["SecurityGroupProbeIpList"] = v.([]interface{})
		}
	}

	if update {
		action := "UpdateHoneypotNode"
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
		stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, threatDetectionService.ThreatDetectionHoneypotNodeStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudThreatDetectionHoneypotNodeRead(d, meta)
}

func resourceAlicloudThreatDetectionHoneypotNodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}
	var err error

	request := map[string]interface{}{

		"NodeId": d.Id(),
	}

	action := "DeleteHoneypotNode"
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutCreate), 5*time.Second, threatDetectionService.ThreatDetectionHoneypotNodeStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

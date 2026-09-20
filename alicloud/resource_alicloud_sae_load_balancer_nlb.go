package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSaeLoadBalancerNlb() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeLoadBalancerNlbCreate,
		Read:   resourceAlicloudSaeLoadBalancerNlbRead,
		Update: resourceAlicloudSaeLoadBalancerNlbUpdate,
		Delete: resourceAlicloudSaeLoadBalancerNlbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nlb_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"address_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"listeners": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"target_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cert_ids": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"zone_mappings": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by_sae": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudSaeLoadBalancerNlbCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	var response map[string]interface{}
	action := "/pop/v1/sam/app/nlb"
	request := make(map[string]*string)
	request["AppId"] = StringPointer(d.Get("app_id").(string))
	if v, ok := d.GetOk("nlb_id"); ok {
		request["NlbId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("address_type"); ok {
		request["AddressType"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("listeners"); ok {
		listenersReq := make([]interface{}, 0)
		for _, listener := range v.(*schema.Set).List() {
			listenerMap := listener.(map[string]interface{})
			listenersReq = append(listenersReq, map[string]interface{}{
				"port":       listenerMap["port"],
				"targetPort": listenerMap["target_port"],
				"protocol":   listenerMap["protocol"],
				"certIds":    listenerMap["cert_ids"],
			})
		}
		obj, err := json.Marshal(listenersReq)
		if err != nil {
			return WrapError(err)
		}
		request["Listeners"] = StringPointer(string(obj))
	}
	if v, ok := d.GetOk("zone_mappings"); ok {
		zoneMappingsReq := make([]interface{}, 0)
		for _, zoneMapping := range v.(*schema.Set).List() {
			zoneMappingMap := zoneMapping.(map[string]interface{})
			zoneMappingsReq = append(zoneMappingsReq, map[string]interface{}{
				"vSwitchId": zoneMappingMap["vswitch_id"],
				"zoneId":    zoneMappingMap["zone_id"],
			})
		}
		obj, err := json.Marshal(zoneMappingsReq)
		if err != nil {
			return WrapError(err)
		}
		request["ZoneMappings"] = StringPointer(string(obj))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var err error
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Get("app_id").(string), "POST "+action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(d.Get("app_id").(string)))
	stateConf := BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, saeService.SaeApplicationStateRefreshFunc(d.Get("app_id").(string), []string{"FAIL", "AUTO_BATCH_WAIT", "APPROVED", "WAIT_APPROVAL", "WAIT_BATCH_CONFIRM", "ABORT", "SYSTEM_FAIL"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudSaeLoadBalancerNlbRead(d, meta)
}

func resourceAlicloudSaeLoadBalancerNlbRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	object, err := saeService.DescribeApplicationNlb(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_load_balancer_nlb DescribeApplicationNlb Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	instancesRaw, err := jsonpath.Get("$.Instances", object)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Instances", object)
	}
	instances := extractInstanceList(instancesRaw)
	if len(instances) == 0 {
		if !d.IsNewResource() {
			d.SetId("")
			return nil
		}
		return nil
	}
	instance := instances[0]
	d.Set("dns_name", instance["DnsName"])
	d.Set("created_by_sae", instance["CreatedBySae"])
	listenersArray := make([]interface{}, 0)
	if listenersRaw, ok := instance["Listeners"]; ok && listenersRaw != nil {
		listenersList := extractInstanceList(listenersRaw)
		for _, listener := range listenersList {
			listenersArray = append(listenersArray, map[string]interface{}{
				"port":        listener["Port"],
				"target_port": listener["TargetPort"],
				"protocol":    listener["Protocol"],
				"cert_ids":    listener["CertIds"],
			})
		}
	}
	if err := d.Set("listeners", listenersArray); err != nil {
		return WrapError(err)
	}
	d.Set("app_id", d.Id())
	return nil
}

func resourceAlicloudSaeLoadBalancerNlbUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	var response map[string]interface{}
	action := "/pop/v1/sam/app/nlb"
	request := make(map[string]*string)
	request["AppId"] = StringPointer(d.Id())
	if v, ok := d.GetOk("nlb_id"); ok {
		request["NlbId"] = StringPointer(v.(string))
	}
	update := false
	if d.HasChange("listeners") {
		update = true
	}
	if v, ok := d.GetOk("listeners"); ok {
		listenersReq := make([]interface{}, 0)
		for _, listener := range v.(*schema.Set).List() {
			listenerMap := listener.(map[string]interface{})
			listenersReq = append(listenersReq, map[string]interface{}{
				"port":       listenerMap["port"],
				"targetPort": listenerMap["target_port"],
				"protocol":   listenerMap["protocol"],
				"certIds":    listenerMap["cert_ids"],
			})
		}
		obj, err := json.Marshal(listenersReq)
		if err != nil {
			return WrapError(err)
		}
		request["Listeners"] = StringPointer(string(obj))
	}
	if d.HasChange("zone_mappings") {
		update = true
	}
	if v, ok := d.GetOk("zone_mappings"); ok {
		zoneMappingsReq := make([]interface{}, 0)
		for _, zoneMapping := range v.(*schema.Set).List() {
			zoneMappingMap := zoneMapping.(map[string]interface{})
			zoneMappingsReq = append(zoneMappingsReq, map[string]interface{}{
				"vSwitchId": zoneMappingMap["vswitch_id"],
				"zoneId":    zoneMappingMap["zone_id"],
			})
		}
		obj, err := json.Marshal(zoneMappingsReq)
		if err != nil {
			return WrapError(err)
		}
		request["ZoneMappings"] = StringPointer(string(obj))
	}
	if update {
		// SAE BindNlb (POST /pop/v1/sam/app/nlb) is additive: it appends new
		// listeners/zone mappings without removing the existing binding, so a
		// binding change must unbind-then-rebind. First clear the current binding
		// with RoaDelete (UnbindNlb) on the same path using AppId + NlbId — the
		// same call the Delete function relies on — then RoaPost (BindNlb) with
		// the new configuration built above.
		unbindRequest := map[string]*string{
			"AppId": StringPointer(d.Id()),
		}
		if v, ok := d.GetOk("nlb_id"); ok {
			unbindRequest["NlbId"] = StringPointer(v.(string))
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		var unbindResponse map[string]interface{}
		var unbindErr error
		unbindErr = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			unbindResponse, unbindErr = client.RoaDelete("sae", "2019-05-06", action, unbindRequest, nil, nil, false)
			if unbindErr != nil {
				if IsExpectedErrors(unbindErr, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(unbindErr) {
					wait()
					return resource.RetryableError(unbindErr)
				}
				return resource.NonRetryableError(unbindErr)
			}
			return nil
		})
		addDebug("DELETE "+action, unbindResponse, unbindRequest)
		// A NotFound on unbind means the previous binding is already gone; that
		// is safe and we proceed to rebind. Any other error is fatal.
		if unbindErr != nil && !NotFoundError(unbindErr) {
			return WrapErrorf(unbindErr, DefaultErrorMsg, d.Id(), "DELETE "+action, AlibabaCloudSdkGoERROR)
		}
		// UnbindNlb's change order terminates in FAIL while the app stays
		// RUNNING and healthy, so do not gate it like BindNlb; proceed
		// straight to rebind with the updated listeners / zone mappings —
		// its RoaPost retries on Application.ChangerOrderRunning/InvalidStatus.
		var err error
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, saeService.SaeApplicationStateRefreshFunc(d.Id(), []string{"FAIL", "AUTO_BATCH_WAIT", "APPROVED", "WAIT_APPROVAL", "WAIT_BATCH_CONFIRM", "ABORT", "SYSTEM_FAIL"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudSaeLoadBalancerNlbRead(d, meta)
}

func resourceAlicloudSaeLoadBalancerNlbDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/app/nlb"
	request := map[string]*string{
		"AppId": StringPointer(d.Id()),
	}
	if v, ok := d.GetOk("nlb_id"); ok {
		request["NlbId"] = StringPointer(v.(string))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	var err error
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DELETE "+action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

// extractInstanceList normalises a JSON value (array or map-of-values) into a
// slice of map[string]interface{} entries. SAE DescribeApplicationNlbs returns
// "Instances" and per-instance "listeners" which the API may render as either
// a JSON array or a JSON object keyed by identifier; both forms are handled.
func extractInstanceList(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if raw == nil {
		return result
	}
	switch v := raw.(type) {
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, m)
			}
		}
	case map[string]interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, m)
			}
		}
	}
	return result
}

package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSaeIngress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeIngressCreate,
		Read:   resourceAlicloudSaeIngressRead,
		Update: resourceAlicloudSaeIngressUpdate,
		Delete: resourceAlicloudSaeIngressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slb_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"cert_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cert_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_balance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"listener_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"app_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"container_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rewrite_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"backend_protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"default_rule": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"app_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"container_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudSaeIngressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/pop/v1/sam/ingress/Ingress"
	request := make(map[string]*string)

	var err error

	request["SlbId"] = StringPointer(d.Get("slb_id").(string))
	request["ListenerPort"] = StringPointer(strconv.Itoa(d.Get("listener_port").(int)))
	request["NamespaceId"] = StringPointer(d.Get("namespace_id").(string))

	if v, ok := d.GetOk("cert_id"); ok {
		request["CertId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("cert_ids"); ok {
		request["CertIds"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("load_balance_type"); ok {
		request["LoadBalanceType"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("listener_protocol"); ok {
		request["ListenerProtocol"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = StringPointer(v.(string))
	}

	defaultrulesMap := map[string]interface{}{}

	for _, rules := range d.Get("default_rule").(*schema.Set).List() {
		rulesArg := rules.(map[string]interface{})
		defaultrulesMap["appId"] = rulesArg["app_id"]
		defaultrulesMap["containerPort"] = rulesArg["container_port"]
	}
	if v, err := convertArrayObjectToJsonString(defaultrulesMap); err != nil {
		return WrapError(err)
	} else {
		request["DefaultRule"] = StringPointer(v)
	}

	rulesMaps := make([]map[string]interface{}, 0)
	for _, rules := range d.Get("rules").(*schema.Set).List() {
		rulesArg := rules.(map[string]interface{})
		rulesMap := map[string]interface{}{}
		rulesMap["appId"] = rulesArg["app_id"]
		rulesMap["appName"] = rulesArg["app_name"]
		rulesMap["containerPort"] = rulesArg["container_port"]
		rulesMap["domain"] = rulesArg["domain"]
		rulesMap["path"] = rulesArg["path"]

		if rewritePath, ok := rulesArg["rewrite_path"]; ok {
			rulesMap["rewritePath"] = rewritePath
		}

		if backendProtocol, ok := rulesArg["backend_protocol"]; ok {
			rulesMap["backendProtocol"] = backendProtocol
		}

		rulesMaps = append(rulesMaps, rulesMap)
	}

	if v, err := convertArrayObjectToJsonString(rulesMaps); err != nil {
		return WrapError(err)
	} else {
		request["Rules"] = StringPointer(v)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action+"-Create", response, fmt.Sprint(request))

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_ingress", action, AlibabaCloudSdkGoERROR)
	}

	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["IngressId"]))

	return resourceAlicloudSaeIngressRead(d, meta)
}

func resourceAlicloudSaeIngressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	object, err := saeService.DescribeSaeIngress(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_ingress saeService.DescribeSaeIngress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("namespace_id", object["NamespaceId"])
	d.Set("slb_id", object["SlbId"])
	d.Set("load_balance_type", object["LoadBalanceType"])
	d.Set("listener_protocol", object["ListenerProtocol"])
	d.Set("description", object["Description"])

	if v, ok := object["CertId"]; ok && v.(string) != "" {
		d.Set("cert_id", v)
	}

	if v, ok := object["CertIds"]; ok && v.(string) != "" {
		d.Set("cert_ids", v)
	}

	if v, ok := object["ListenerPort"]; ok && fmt.Sprint(v) != "0" {
		d.Set("listener_port", formatInt(v))
	}

	defaultRuleConfig := make([]map[string]interface{}, 0)
	if defaultRule, ok := object["DefaultRule"]; ok {
		defaultRule_convert := defaultRule.(map[string]interface{})
		defaultRuleData := make(map[string]interface{}, 0)
		defaultRuleData["app_id"] = defaultRule_convert["AppId"]
		defaultRuleData["container_port"] = defaultRule_convert["ContainerPort"]
		defaultRuleConfig = append(defaultRuleConfig, defaultRuleData)
		d.Set("default_rule", defaultRuleConfig)
	}

	config := make([]map[string]interface{}, 0)
	if quotaDimension, ok := object["Rules"]; ok {
		quotaDimension_convert := quotaDimension.([]interface{})
		for _, obj := range quotaDimension_convert {
			obj_convert := obj.(map[string]interface{})
			data := make(map[string]interface{}, 0)
			data["app_id"] = obj_convert["AppId"]
			data["app_name"] = obj_convert["AppName"]
			data["container_port"] = obj_convert["ContainerPort"]
			data["domain"] = obj_convert["Domain"]
			data["path"] = obj_convert["Path"]

			if rewritePath, ok := obj_convert["RewritePath"]; ok {
				data["rewrite_path"] = rewritePath
			}

			if backendProtocol, ok := obj_convert["BackendProtocol"]; ok {
				data["backend_protocol"] = backendProtocol
			}

			config = append(config, data)
		}

		d.Set("rules", config)
	}

	return nil
}

func resourceAlicloudSaeIngressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]*string{
		"IngressId": StringPointer(d.Id()),
	}

	if d.HasChange("cert_id") {
		update = true
	}
	if v, ok := d.GetOk("cert_id"); ok && v.(string) != "" {
		request["CertId"] = StringPointer(v.(string))
	}

	if d.HasChange("cert_ids") {
		update = true
	}
	if v, ok := d.GetOk("cert_ids"); ok && v.(string) != "" {
		request["CertIds"] = StringPointer(v.(string))
	}

	if d.HasChange("load_balance_type") {
		update = true
		if v, ok := d.GetOk("load_balance_type"); ok {
			request["LoadBalanceType"] = StringPointer(v.(string))
		}
	}

	if d.HasChange("listener_protocol") {
		update = true
		if v, ok := d.GetOk("listener_protocol"); ok {
			request["ListenerProtocol"] = StringPointer(v.(string))
		}
	}

	if d.HasChange("default_rule") {
		update = true
		defaultrulesMap := map[string]interface{}{}
		for _, rules := range d.Get("default_rule").(*schema.Set).List() {
			rulesArg := rules.(map[string]interface{})
			defaultrulesMap["appId"] = rulesArg["app_id"]
			defaultrulesMap["containerPort"] = rulesArg["container_port"]
		}
		if v, err := convertArrayObjectToJsonString(defaultrulesMap); err == nil {
			request["DefaultRule"] = StringPointer(v)
		} else {
			return WrapError(err)
		}
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = StringPointer(v.(string))
		}
	}

	if d.HasChange("listener_port") {
		update = true
		request["ListenerPort"] = StringPointer(strconv.Itoa(d.Get("listener_port").(int)))
	}

	if d.HasChange("rules") {
		update = true
		rulesMaps := make([]map[string]interface{}, 0)
		for _, rules := range d.Get("rules").(*schema.Set).List() {
			rulesArg := rules.(map[string]interface{})
			rulesMap := map[string]interface{}{}
			rulesMap["appId"] = rulesArg["app_id"]
			rulesMap["appName"] = rulesArg["app_name"]
			rulesMap["containerPort"] = rulesArg["container_port"]
			rulesMap["domain"] = rulesArg["domain"]
			rulesMap["path"] = rulesArg["path"]

			if rewritePath, ok := rulesArg["rewrite_path"]; ok {
				rulesMap["rewritePath"] = rewritePath
			}

			if backendProtocol, ok := rulesArg["backend_protocol"]; ok {
				rulesMap["backendProtocol"] = backendProtocol
			}

			rulesMaps = append(rulesMaps, rulesMap)
		}
		if v, err := convertArrayObjectToJsonString(rulesMaps); err == nil {
			request["Rules"] = StringPointer(v)
		} else {
			return WrapError(err)
		}
	}

	if update {
		action := "/pop/v1/sam/ingress/Ingress"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RoaPut("sae", "2019-05-06", action, request, nil, nil, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action+"-Update", response, fmt.Sprint(request))

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudSaeIngressRead(d, meta)
}

func resourceAlicloudSaeIngressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/ingress/Ingress"
	var response map[string]interface{}
	var err error
	request := map[string]*string{
		"IngressId": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action+"-Delete", response, fmt.Sprint(request))

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

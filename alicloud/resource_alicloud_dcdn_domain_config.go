package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDcdnDomainConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDcdnDomainConfigCreate,
		Read:   resourceAliCloudDcdnDomainConfigRead,
		Update: resourceAliCloudDcdnDomainConfigUpdate,
		Delete: resourceAliCloudDcdnDomainConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringLenBetween(5, 67),
			},
			"function_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"function_args": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arg_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"arg_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"config_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDcdnDomainConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	var response map[string]interface{}
	action := "BatchSetDcdnDomainConfigs"
	request := make(map[string]interface{})
	var err error

	request["DomainNames"] = d.Get("domain_name").(string)

	config := make([]map[string]interface{}, 1)
	functionArgs := d.Get("function_args").(*schema.Set).List()
	args := make([]map[string]interface{}, len(functionArgs))
	for key, value := range functionArgs {
		arg := value.(map[string]interface{})
		args[key] = map[string]interface{}{
			"argName":  arg["arg_name"],
			"argValue": arg["arg_value"],
		}
	}

	config[0] = map[string]interface{}{
		"functionArgs": args,
		"functionName": d.Get("function_name").(string),
	}

	if v, ok := d.GetOk("parent_id"); ok {
		config[0]["parentId"] = v.(string)
	}

	bytConfig, err := json.Marshal(config)
	if err != nil {
		return WrapError(err)
	}

	request["Functions"] = string(bytConfig)

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"FlowControlError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dcdn_domain", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.DomainConfigList.DomainConfigModel", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_dcdn_domain_config", "$.DomainConfigList.DomainConfigModel", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_dcdn_domain_config", "$.DomainConfigList.DomainConfigModel", response)
	}

	for _, v := range resp.([]interface{}) {
		dcdnConfigModel := v.(map[string]interface{})
		if fmt.Sprint(dcdnConfigModel["DomainName"]) == fmt.Sprint(request["DomainNames"]) && fmt.Sprint(dcdnConfigModel["FunctionName"]) == fmt.Sprint(d.Get("function_name").(string)) {
			d.SetId(fmt.Sprintf("%s:%s:%s", dcdnConfigModel["DomainName"], dcdnConfigModel["FunctionName"], dcdnConfigModel["ConfigId"]))
		} else {
			return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_dcdn_domain_config", "$.DomainConfigList.DomainConfigModel", response)
		}
	}

	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 2*time.Second, dcdnService.DcdnDomainConfigStateRefreshFunc(d.Id(), []string{"failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudDcdnDomainConfigRead(d, meta)
}

func resourceAliCloudDcdnDomainConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := &DcdnService{client: client}

	object, err := dcdnService.DescribeDcdnDomainConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	d.Set("domain_name", parts[0])
	d.Set("function_name", object["FunctionName"])
	d.Set("config_id", object["ConfigId"])
	d.Set("status", object["Status"])
	d.Set("parent_id", object["ParentId"])

	ignoreFunctionArg := []string{"dsl", "disable_l2_log", "dynamic_batch_route", "dynamic_enable_cpool_chash",
		"dynamic_mux_ecn_enable", "dynamic_mux_keepalive_enable", "dynamic_mux_share_enable", "dynamic_pk",
		"dynamic_retry_status", "dynamic_route_cdn_v2", "dynamic_route_cpool", "dynamic_route_magic",
		"dynamic_route_round_robin", "dynamic_route_session", "dynamic_route_tunnel", "dynamie_route_alllink_log",
		"keepalive_sni", "l7tol4", "partition_back_to_origin", "dynamic_route_adapt_cache",
		"dynamic_route_http_methods", "dynamic_mux_tls_enable", "dynamic_route_cdn_v2"}

	if functionArgs, ok := object["FunctionArgs"]; ok {
		if functionArgList, ok := functionArgs.(map[string]interface{})["FunctionArg"]; ok {
			functionArgMaps := make([]map[string]interface{}, 0)
			for _, functionArg := range functionArgList.([]interface{}) {
				functionArgItem := functionArg.(map[string]interface{})
				functionArgMap := map[string]interface{}{}

				// This function args is extra, filter them to pass test check.
				if InArray(fmt.Sprint(functionArgItem["ArgName"]), ignoreFunctionArg) {
					continue
				}

				if argName, ok := functionArgItem["ArgName"]; ok {
					functionArgMap["arg_name"] = argName
				}

				if argValue, ok := functionArgItem["ArgValue"]; ok {
					functionArgMap["arg_value"] = argValue
				}

				functionArgMaps = append(functionArgMaps, functionArgMap)
			}

			d.Set("function_args", functionArgMaps)
		}
	}

	return nil
}

func resourceAliCloudDcdnDomainConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	var response map[string]interface{}

	if d.HasChange("function_args") || d.HasChange("parent_id") {
		action := "BatchSetDcdnDomainConfigs"
		request := make(map[string]interface{})

		parts, err := ParseResourceId(d.Id(), 3)
		if err != nil {
			return WrapError(err)
		}

		request["DomainNames"] = parts[0]

		config := make([]map[string]interface{}, 1)
		functionArgs := d.Get("function_args").(*schema.Set).List()

		args := make([]map[string]interface{}, len(functionArgs))
		for key, value := range functionArgs {
			arg := value.(map[string]interface{})
			args[key] = map[string]interface{}{
				"argName":  arg["arg_name"],
				"argValue": arg["arg_value"],
			}
		}

		config[0] = map[string]interface{}{
			"functionArgs": args,
			"functionName": parts[1],
			"configId":     parts[2],
		}

		if v, ok := d.GetOk("parent_id"); ok {
			config[0]["parentId"] = v
		}

		bytconfig, _ := json.Marshal(config)

		request["Functions"] = string(bytconfig)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"FlowControlError"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 2*time.Second, dcdnService.DcdnDomainConfigStateRefreshFunc(d.Id(), []string{"failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudDcdnDomainConfigRead(d, meta)
}

func resourceAliCloudDcdnDomainConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	action := "DeleteDcdnSpecificConfig"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"DomainName": parts[0],
		"ConfigId":   parts[2],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Second, dcdnService.DcdnDomainConfigStateRefreshFunc(d.Id(), []string{"failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

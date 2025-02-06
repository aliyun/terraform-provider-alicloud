// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEnsLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsLoadBalancerCreate,
		Read:   resourceAliCloudEnsLoadBalancerRead,
		Update: resourceAliCloudEnsLoadBalancerUpdate,
		Delete: resourceAliCloudEnsLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backend_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"ens"}, false),
						},
						"server_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(0, 65535),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 100),
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ens_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_balancer_spec": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEnsLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["EnsRegionId"] = d.Get("ens_region_id")
	request["LoadBalancerSpec"] = d.Get("load_balancer_spec")
	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}
	request["PayType"] = convertEnsLoadBalancerPayTypeRequest(d.Get("payment_type").(string))
	request["NetworkId"] = d.Get("network_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_load_balancer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LoadBalancerId"]))

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, ensServiceV2.EnsLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEnsLoadBalancerUpdate(d, meta)
}

func resourceAliCloudEnsLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsLoadBalancer(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_load_balancer DescribeEnsLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["EnsRegionId"] != nil {
		d.Set("ens_region_id", objectRaw["EnsRegionId"])
	}
	if objectRaw["LoadBalancerName"] != nil {
		d.Set("load_balancer_name", objectRaw["LoadBalancerName"])
	}
	if objectRaw["LoadBalancerSpec"] != nil {
		d.Set("load_balancer_spec", objectRaw["LoadBalancerSpec"])
	}
	if objectRaw["NetworkId"] != nil {
		d.Set("network_id", objectRaw["NetworkId"])
	}
	if convertEnsLoadBalancerPayTypeResponse(objectRaw["PayType"]) != nil {
		d.Set("payment_type", convertEnsLoadBalancerPayTypeResponse(objectRaw["PayType"]))
	}
	if objectRaw["LoadBalancerStatus"] != nil {
		d.Set("status", objectRaw["LoadBalancerStatus"])
	}
	if objectRaw["VSwitchId"] != nil {
		d.Set("vswitch_id", objectRaw["VSwitchId"])
	}

	backendServers1Raw := objectRaw["BackendServers"]
	backendServersMaps := make([]map[string]interface{}, 0)
	if backendServers1Raw != nil {
		for _, backendServersChild1Raw := range backendServers1Raw.([]interface{}) {
			backendServersMap := make(map[string]interface{})
			backendServersChild1Raw := backendServersChild1Raw.(map[string]interface{})
			backendServersMap["ip"] = backendServersChild1Raw["Ip"]
			backendServersMap["port"] = backendServersChild1Raw["Port"]
			backendServersMap["server_id"] = backendServersChild1Raw["ServerId"]
			backendServersMap["type"] = backendServersChild1Raw["Type"]
			backendServersMap["weight"] = backendServersChild1Raw["Weight"]

			backendServersMaps = append(backendServersMaps, backendServersMap)
		}
	}
	if objectRaw["BackendServers"] != nil {
		if err := d.Set("backend_servers", backendServersMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceAliCloudEnsLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyLoadBalancerAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["LoadBalancerId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("load_balancer_name") {
		update = true
	}
	request["LoadBalancerName"] = d.Get("load_balancer_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("backend_servers") {
		oldEntry, newEntry := d.GetChange("backend_servers")
		removed := oldEntry
		added := newEntry

		if len(removed.([]interface{})) > 0 {
			action := "RemoveBackendServers"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["LoadBalancerId"] = d.Id()

			localData := removed.([]interface{})
			backendServersMaps := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["ServerId"] = dataLoopTmp["server_id"]
				dataLoopMap["Weight"] = dataLoopTmp["weight"]
				dataLoopMap["Type"] = dataLoopTmp["type"]
				dataLoopMap["Ip"] = dataLoopTmp["ip"]
				dataLoopMap["Port"] = dataLoopTmp["port"]
				backendServersMaps = append(backendServersMaps, dataLoopMap)
			}
			backendServersMapsJson, err := json.Marshal(backendServersMaps)
			if err != nil {
				return WrapError(err)
			}
			request["BackendServers"] = string(backendServersMapsJson)

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added.([]interface{})) > 0 {
			action := "AddBackendServers"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["LoadBalancerId"] = d.Id()

			localData := added.([]interface{})
			backendServersMaps := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["ServerId"] = dataLoopTmp["server_id"]
				dataLoopMap["Weight"] = dataLoopTmp["weight"]
				dataLoopMap["Type"] = dataLoopTmp["type"]
				dataLoopMap["Ip"] = dataLoopTmp["ip"]
				dataLoopMap["Port"] = dataLoopTmp["port"]
				backendServersMaps = append(backendServersMaps, dataLoopMap)
			}
			backendServersMapsJson, err := json.Marshal(backendServersMaps)
			if err != nil {
				return WrapError(err)
			}
			request["BackendServers"] = string(backendServersMapsJson)

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}
	return resourceAliCloudEnsLoadBalancerRead(d, meta)
}

func resourceAliCloudEnsLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["LoadBalancerId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, ensServiceV2.EnsLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertEnsLoadBalancerPayTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
func convertEnsLoadBalancerPayTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}

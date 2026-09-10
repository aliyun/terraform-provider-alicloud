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
				Type:     schema.TypeSet,
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
							Computed: true,
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
							Computed:     true,
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

	request["ClientToken"] = buildClientToken(action)

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
		return nil
	})
	addDebug(action, response, request)

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

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("ens_region_id", objectRaw["EnsRegionId"])
	d.Set("load_balancer_name", objectRaw["LoadBalancerName"])
	d.Set("load_balancer_spec", objectRaw["LoadBalancerSpec"])
	d.Set("network_id", objectRaw["NetworkId"])
	d.Set("payment_type", convertEnsLoadBalancerPayTypeResponse(objectRaw["PayType"]))
	d.Set("status", objectRaw["LoadBalancerStatus"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])

	backendServersRaw := objectRaw["BackendServers"]
	backendServersMaps := make([]map[string]interface{}, 0)
	if backendServersRaw != nil {
		for _, backendServersChildRaw := range convertToInterfaceArray(backendServersRaw) {
			backendServersMap := make(map[string]interface{})
			backendServersChildRaw := backendServersChildRaw.(map[string]interface{})
			backendServersMap["ip"] = backendServersChildRaw["Ip"]
			backendServersMap["port"] = formatInt(backendServersChildRaw["Port"])
			backendServersMap["server_id"] = backendServersChildRaw["ServerId"]
			backendServersMap["type"] = backendServersChildRaw["Type"]
			backendServersMap["weight"] = backendServersChildRaw["Weight"]

			backendServersMaps = append(backendServersMaps, backendServersMap)
		}
	}
	if err := d.Set("backend_servers", backendServersMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudEnsLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyLoadBalancerAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()

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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("backend_servers") {
		oldEntry, newEntry := d.GetChange("backend_servers")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action := "RemoveBackendServers"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["LoadBalancerId"] = d.Id()

			localData := removed.List()
			backendServersMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["ServerId"] = dataLoopTmp["server_id"]
				dataLoopMap["Weight"] = dataLoopTmp["weight"]
				dataLoopMap["Type"] = dataLoopTmp["type"]
				dataLoopMap["Ip"] = dataLoopTmp["ip"]
				dataLoopMap["Port"] = dataLoopTmp["port"]
				backendServersMapsArray = append(backendServersMapsArray, dataLoopMap)
			}
			backendServersMapsJson, err := json.Marshal(backendServersMapsArray)
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
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if added.Len() > 0 {
			action := "AddBackendServers"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["LoadBalancerId"] = d.Id()

			localData := added.List()
			backendServersMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["ServerId"] = dataLoopTmp["server_id"]
				dataLoopMap["Weight"] = dataLoopTmp["weight"]
				dataLoopMap["Type"] = dataLoopTmp["type"]
				dataLoopMap["Ip"] = dataLoopTmp["ip"]
				dataLoopMap["Port"] = dataLoopTmp["port"]
				backendServersMapsArray = append(backendServersMapsArray, dataLoopMap)
			}
			backendServersMapsJson, err := json.Marshal(backendServersMapsArray)
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
				return nil
			})
			addDebug(action, response, request)
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
	request["LoadBalancerId"] = d.Id()

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
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, ensServiceV2.EnsLoadBalancerStateRefreshFunc(d.Id(), "$.LoadBalancerId", []string{}))
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

// Package ens. This file is generated automatically. Please do not modify it manually, thank you!
package ens

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEnsLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnsLoadBalancerCreate,
		Read:   resourceEnsLoadBalancerRead,
		Update: resourceEnsLoadBalancerUpdate,
		Delete: resourceEnsLoadBalancerDelete,
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
							ValidateFunc: validators.StringInSlice([]string{"ens"}, false),
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
							ValidateFunc: validators.IntBetween(0, 65535),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validators.IntBetween(0, 100),
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
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo"}, false),
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

func resourceEnsLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = helper.BuildClientToken(action)

	request["EnsRegionId"] = d.Get("ens_region_id")
	request["LoadBalancerSpec"] = d.Get("load_balancer_spec")
	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}
	request["PayType"] = convertEnsLoadBalancerPayTypeRequest(d.Get("payment_type").(string))
	request["NetworkId"] = d.Get("network_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ens_load_balancer", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["LoadBalancerId"]))

	ensServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, ensServiceV2.EnsLoadBalancerStateRefreshFunc(d.Id(), "LoadBalancerStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceEnsLoadBalancerUpdate(d, meta)
}

func resourceEnsLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := NewService(client)

	objectRaw, err := ensServiceV2.DescribeEnsLoadBalancer(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_load_balancer DescribeEnsLoadBalancer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
		for _, backendServersChildRaw := range convert.ConvertToInterfaceArray(backendServersRaw) {
			backendServersMap := make(map[string]interface{})
			backendServersChildRaw := backendServersChildRaw.(map[string]interface{})
			backendServersMap["ip"] = backendServersChildRaw["Ip"]
			backendServersChildRawPort, err := convert.FormatInt(backendServersChildRaw["Port"])
			if err != nil {
				return err
			}
			backendServersMap["port"] = backendServersChildRawPort
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

func resourceEnsLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
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
				return errs.WrapError(err)
			}
			request["BackendServers"] = string(backendServersMapsJson)

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
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
				return errs.WrapError(err)
			}
			request["BackendServers"] = string(backendServersMapsJson)

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
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

		}
	}
	return resourceEnsLoadBalancerRead(d, meta)
}

func resourceEnsLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["LoadBalancerId"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	ensServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, ensServiceV2.EnsLoadBalancerStateRefreshFunc(d.Id(), "$.LoadBalancerId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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

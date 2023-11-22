// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
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
	conn, err := client.NewEnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["EnsRegionId"] = d.Get("ens_region_id")
	request["LoadBalancerSpec"] = d.Get("load_balancer_spec")
	if v, ok := d.GetOk("load_balancer_name"); ok {
		request["LoadBalancerName"] = v
	}
	request["PayType"] = convertEnsPayTypeRequest(d.Get("payment_type").(string))
	request["NetworkId"] = d.Get("network_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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

	return resourceAliCloudEnsLoadBalancerRead(d, meta)
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
	d.Set("payment_type", convertEnsPayTypeResponse(objectRaw["PayType"]))
	d.Set("status", objectRaw["LoadBalancerStatus"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])

	return nil
}

func resourceAliCloudEnsLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyLoadBalancerAttribute"
	conn, err := client.NewEnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["LoadBalancerId"] = d.Id()
	if d.HasChange("load_balancer_name") {
		update = true
		request["LoadBalancerName"] = d.Get("load_balancer_name")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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

	return resourceAliCloudEnsLoadBalancerRead(d, meta)
}

func resourceAliCloudEnsLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLoadBalancer"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewEnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["LoadBalancerId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-10"), StringPointer("AK"), query, request, &runtime)

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

func convertEnsPayTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
func convertEnsPayTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}

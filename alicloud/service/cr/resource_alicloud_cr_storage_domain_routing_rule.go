// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cr

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCrStorageDomainRoutingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrStorageDomainRoutingRuleCreate,
		Read:   resourceAliCloudCrStorageDomainRoutingRuleRead,
		Update: resourceAliCloudCrStorageDomainRoutingRuleUpdate,
		Delete: resourceAliCloudCrStorageDomainRoutingRuleDelete,
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
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_domain": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return strings.TrimPrefix(old, "https://") == strings.TrimPrefix(new, "https://")
							},
						},
						"endpoint_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"instance_domain": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCrStorageDomainRoutingRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateStorageDomainRoutingRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("routes"); ok {
		routesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convert.ConvertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["EndpointType"] = dataLoopTmp["endpoint_type"]
			dataLoopMap["InstanceDomain"] = dataLoopTmp["instance_domain"]
			dataLoopMap["StorageDomain"] = dataLoopTmp["storage_domain"]
			routesMapsArray = append(routesMapsArray, dataLoopMap)
		}
		routesMapsJson, err := json.Marshal(routesMapsArray)
		if err != nil {
			return errs.WrapError(err)
		}
		request["Routes"] = string(routesMapsJson)
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cr_storage_domain_routing_rule", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], response["RuleId"]))

	return resourceAliCloudCrStorageDomainRoutingRuleRead(d, meta)
}

func resourceAliCloudCrStorageDomainRoutingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crServiceV2 := CrServiceV2{client}

	objectRaw, err := crServiceV2.DescribeCrStorageDomainRoutingRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_storage_domain_routing_rule DescribeCrStorageDomainRoutingRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("rule_id", objectRaw["RuleId"])

	routesRaw := objectRaw["Routes"]
	routesMaps := make([]map[string]interface{}, 0)
	if routesRaw != nil {
		for _, routesChildRaw := range convert.ConvertToInterfaceArray(routesRaw) {
			routesMap := make(map[string]interface{})
			routesChildRaw := routesChildRaw.(map[string]interface{})
			routesMap["endpoint_type"] = routesChildRaw["EndpointType"]
			routesMap["instance_domain"] = routesChildRaw["InstanceDomain"]
			routesMap["storage_domain"] = routesChildRaw["StorageDomain"]

			routesMaps = append(routesMaps, routesMap)
		}
	}
	if err := d.Set("routes", routesMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudCrStorageDomainRoutingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateStorageDomainRoutingRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RuleId"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("routes") {
		update = true
	}
	if v, ok := d.GetOk("routes"); ok || d.HasChange("routes") {
		routesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convert.ConvertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["EndpointType"] = dataLoopTmp["endpoint_type"]
			dataLoopMap["InstanceDomain"] = dataLoopTmp["instance_domain"]
			dataLoopMap["StorageDomain"] = dataLoopTmp["storage_domain"]
			routesMapsArray = append(routesMapsArray, dataLoopMap)
		}
		routesMapsJson, err := json.Marshal(routesMapsArray)
		if err != nil {
			return errs.WrapError(err)
		}
		request["Routes"] = string(routesMapsJson)
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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

	return resourceAliCloudCrStorageDomainRoutingRuleRead(d, meta)
}

func resourceAliCloudCrStorageDomainRoutingRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteStorageDomainRoutingRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RuleId"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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

	return nil
}

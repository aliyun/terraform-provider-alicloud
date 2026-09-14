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

func resourceAliCloudRocketmqDisasterRecoveryPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRocketmqDisasterRecoveryPlanCreate,
		Read:   resourceAliCloudRocketmqDisasterRecoveryPlanRead,
		Update: resourceAliCloudRocketmqDisasterRecoveryPlanUpdate,
		Delete: resourceAliCloudRocketmqDisasterRecoveryPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_sync_checkpoint": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"consumer_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"endpoint_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_role": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"message_property": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"property_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"property_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"network_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"plan_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plan_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_checkpoint_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRocketmqDisasterRecoveryPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/disaster_recovery"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("plan_name"); ok {
		request["planName"] = v
	}
	if v, ok := d.GetOk("plan_desc"); ok {
		request["planDesc"] = v
	}
	if v, ok := d.GetOk("plan_type"); ok {
		request["planType"] = v
	}
	if v, ok := d.GetOkExists("auto_sync_checkpoint"); ok {
		request["autoSyncCheckpoint"] = v
	}
	if v, ok := d.GetOkExists("sync_checkpoint_enabled"); ok {
		request["syncCheckpointEnabled"] = v
	}
	if v, ok := d.GetOk("instances"); ok && len(v.([]interface{})) > 0 {
		request["instances"] = expandRocketmqDisasterRecoveryPlanInstances(v.([]interface{}))
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("RocketMQ", "2022-08-01", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rocketmq_disaster_recovery_plan", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data", response)
	d.SetId(fmt.Sprint(id))

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"CREATED"}, d.Timeout(schema.TimeoutCreate), 4*time.Minute, rocketmqServiceV2.RocketmqDisasterRecoveryPlanStateRefreshFunc(d.Id(), "planStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRocketmqDisasterRecoveryPlanRead(d, meta)
}

func resourceAliCloudRocketmqDisasterRecoveryPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rocketmqServiceV2 := RocketmqServiceV2{client}

	object, err := rocketmqServiceV2.DescribeRocketmqDisasterRecoveryPlan(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rocketmq_disaster_recovery_plan DescribeRocketmqDisasterRecoveryPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("plan_id", fmt.Sprint(object["planId"]))
	d.Set("plan_name", object["planName"])
	d.Set("plan_desc", object["planDesc"])
	d.Set("plan_type", object["planType"])
	d.Set("status", object["planStatus"])
	d.Set("create_time", object["createTime"])
	d.Set("update_time", object["updateTime"])
	d.Set("region_id", client.RegionId)
	d.Set("auto_sync_checkpoint", object["autoSyncCheckpoint"])
	d.Set("sync_checkpoint_enabled", object["syncCheckpointEnabled"])

	instancesMaps := flattenRocketmqDisasterRecoveryPlanInstances(object["instances"])
	if err := d.Set("instances", instancesMaps); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliCloudRocketmqDisasterRecoveryPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	planId := d.Id()
	action := fmt.Sprintf("/disaster_recovery/%s", planId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	update := false
	d.Partial(true)
	request = make(map[string]interface{})

	if d.HasChange("plan_name") {
		update = true
	}
	if v, ok := d.GetOk("plan_name"); ok {
		request["planName"] = v
	}
	if d.HasChange("plan_desc") {
		update = true
	}
	if v, ok := d.GetOk("plan_desc"); ok {
		request["planDesc"] = v
	}
	if d.HasChange("plan_type") {
		update = true
	}
	if v, ok := d.GetOk("plan_type"); ok {
		request["planType"] = v
	}
	if d.HasChange("auto_sync_checkpoint") {
		update = true
	}
	if v, ok := d.GetOkExists("auto_sync_checkpoint"); ok {
		request["autoSyncCheckpoint"] = v
	}
	if d.HasChange("sync_checkpoint_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("sync_checkpoint_enabled"); ok {
		request["syncCheckpointEnabled"] = v
	}
	if d.HasChange("instances") {
		update = true
	}
	if v, ok := d.GetOk("instances"); ok && len(v.([]interface{})) > 0 {
		request["instances"] = expandRocketmqDisasterRecoveryPlanInstances(v.([]interface{}))
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("RocketMQ", "2022-08-01", action, query, nil, body, true)
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
	d.Partial(false)
	return resourceAliCloudRocketmqDisasterRecoveryPlanRead(d, meta)
}

func resourceAliCloudRocketmqDisasterRecoveryPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	planId := d.Id()
	action := fmt.Sprintf("/disaster_recovery/%s", planId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("RocketMQ", "2022-08-01", action, query, nil, nil, true)
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

	rocketmqServiceV2 := RocketmqServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Minute, rocketmqServiceV2.RocketmqDisasterRecoveryPlanStateRefreshFunc(d.Id(), "$.planId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func expandRocketmqDisasterRecoveryPlanInstances(instancesList []interface{}) []interface{} {
	result := make([]interface{}, 0)
	for _, instance := range instancesList {
		instanceMap := instance.(map[string]interface{})
		item := make(map[string]interface{})
		if v, ok := instanceMap["instance_type"]; ok && v.(string) != "" {
			item["instanceType"] = v
		}
		if v, ok := instanceMap["instance_role"]; ok && v.(string) != "" {
			item["instanceRole"] = v
		}
		if v, ok := instanceMap["security_group_id"]; ok && v.(string) != "" {
			item["securityGroupId"] = v
		}
		if v, ok := instanceMap["password"]; ok && v.(string) != "" {
			item["password"] = v
		}
		if v, ok := instanceMap["instance_id"]; ok && v.(string) != "" {
			item["instanceId"] = v
		}
		if v, ok := instanceMap["consumer_group_id"]; ok && v.(string) != "" {
			item["consumerGroupId"] = v
		}
		if v, ok := instanceMap["region_id"]; ok && v.(string) != "" {
			item["regionId"] = v
		}
		if v, ok := instanceMap["vswitch_id"]; ok && v.(string) != "" {
			item["vSwitchId"] = v
		}
		if v, ok := instanceMap["endpoint_url"]; ok && v.(string) != "" {
			item["endpointUrl"] = v
		}
		if v, ok := instanceMap["vpc_id"]; ok && v.(string) != "" {
			item["vpcId"] = v
		}
		if v, ok := instanceMap["auth_type"]; ok && v.(string) != "" {
			item["authType"] = v
		}
		if v, ok := instanceMap["network_type"]; ok && v.(string) != "" {
			item["networkType"] = v
		}
		if v, ok := instanceMap["username"]; ok && v.(string) != "" {
			item["username"] = v
		}
		if v, ok := instanceMap["message_property"]; ok && len(v.([]interface{})) > 0 {
			mpMap := v.([]interface{})[0].(map[string]interface{})
			messageProperty := make(map[string]interface{})
			if pk, ok := mpMap["property_key"]; ok && pk.(string) != "" {
				messageProperty["propertyKey"] = pk
			}
			if pv, ok := mpMap["property_value"]; ok && pv.(string) != "" {
				messageProperty["propertyValue"] = pv
			}
			if len(messageProperty) > 0 {
				item["messageProperty"] = messageProperty
			}
		}
		result = append(result, item)
	}
	return result
}

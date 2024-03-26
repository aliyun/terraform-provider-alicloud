// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudRedisTairInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRedisTairInstanceCreate,
		Read:   resourceAliCloudRedisTairInstanceRead,
		Update: resourceAliCloudRedisTairInstanceUpdate,
		Delete: resourceAliCloudRedisTairInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"effective_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Immediately", "MaintainTime"}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"force_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"tair_rdb", "tair_scm", "tair_essd"}, false),
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 60}),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(1024, 65535),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secondary_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"shard_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 256),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PL1", "PL2", "PL3"}, false),
			},
			"storage_size_gb": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"tair_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudRedisTairInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTairInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewRedisClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["InstanceClass"] = d.Get("instance_class")
	request["VpcId"] = d.Get("vpc_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["InstanceType"] = d.Get("instance_type")
	if v, ok := d.GetOk("secondary_zone_id"); ok {
		request["SecondaryZoneId"] = v
	}
	if v, ok := d.GetOk("port"); ok {
		request["Port"] = v
	}
	if v, ok := d.GetOk("tair_instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertRedisChargeTypeRequest(v.(string))
	}
	request["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	request["AutoPay"] = "true"
	if v, ok := d.GetOk("shard_count"); ok {
		request["ShardCount"] = v
	}
	if v, ok := d.GetOk("engine_version"); ok {
		request["EngineVersion"] = v
	}
	if v, ok := d.GetOk("storage_size_gb"); ok {
		request["Storage"] = v
	}
	if v, ok := d.GetOk("storage_performance_level"); ok {
		request["StorageType"] = convertRedisStorageTypeRequest(v.(string))
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_redis_tair_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	redisServiceV2 := RedisServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudRedisTairInstanceUpdate(d, meta)
}

func resourceAliCloudRedisTairInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	redisServiceV2 := RedisServiceV2{client}

	objectRaw, err := redisServiceV2.DescribeRedisTairInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_redis_tair_instance DescribeRedisTairInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("engine_version", objectRaw["EngineVersion"])
	d.Set("instance_class", objectRaw["RealInstanceClass"])
	d.Set("instance_type", objectRaw["InstanceType"])
	d.Set("payment_type", convertRedisInstancesDBInstanceAttributeChargeTypeResponse(objectRaw["ChargeType"]))
	d.Set("port", objectRaw["Port"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("secondary_zone_id", objectRaw["SecondaryZoneId"])
	if objectRaw["ShardCount"] != 0 {
		d.Set("shard_count", objectRaw["ShardCount"])
	}
	d.Set("status", objectRaw["InstanceStatus"])
	d.Set("storage_performance_level", convertRedisInstancesDBInstanceAttributeStorageTypeResponse(objectRaw["StorageType"]))
	d.Set("storage_size_gb", formatInt(objectRaw["Storage"]))
	d.Set("tair_instance_name", objectRaw["InstanceName"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("zone_id", objectRaw["ZoneId"])
	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudRedisTairInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyInstanceAttribute"
	conn, err := client.NewRedisClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("tair_instance_name") {
		update = true
		request["InstanceName"] = d.Get("tair_instance_name")
	}

	if !d.IsNewResource() && d.HasChange("password") {
		update = true
		request["NewPassword"] = d.Get("password")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)

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
	update = false
	action = "ModifyInstanceSpec"
	conn, err = client.NewRedisClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("instance_class") {
		update = true
	}
	request["InstanceClass"] = d.Get("instance_class")
	if v, ok := d.GetOkExists("force_upgrade"); ok {
		request["ForceUpgrade"] = v
	}
	if v, ok := d.GetOk("effective_time"); ok {
		request["EffectiveTime"] = v
	}
	request["AutoPay"] = "true"
	if !d.IsNewResource() && d.HasChange("engine_version") {
		update = true
		request["MajorVersion"] = d.Get("engine_version")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
			request["ClientToken"] = buildClientToken(action)

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
		redisServiceV2 := RedisServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "IsOrderCompleted", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}
	update = false
	action = "ModifyResourceGroup"
	conn, err = client.NewRedisClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
			request["ClientToken"] = buildClientToken(action)

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

	if !d.IsNewResource() && d.HasChange("shard_count") {
		oldEntry, newEntry := d.GetChange("shard_count")
		oldEntryValue := oldEntry.(int)
		newEntryValue := newEntry.(int)
		removed := oldEntryValue - newEntryValue
		added := newEntryValue - oldEntryValue

		if removed > 0 {
			action = "DeleteShardingNode"
			conn, err = client.NewRedisClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			request["InstanceId"] = d.Id()
			request["ShardCount"] = removed
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)

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
			redisServiceV2 := RedisServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "IsOrderCompleted", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if added > 0 {
			action = "AddShardingNode"
			conn, err = client.NewRedisClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			request["InstanceId"] = d.Id()
			request["ClientToken"] = buildClientToken(action)
			request["ShardCount"] = added
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)
				request["ClientToken"] = buildClientToken(action)

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
			redisServiceV2 := RedisServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, redisServiceV2.RedisTairInstanceStateRefreshFunc(d.Id(), "IsOrderCompleted", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

	}
	if d.HasChange("tags") {
		redisServiceV2 := RedisServiceV2{client}
		if err := redisServiceV2.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}

	}
	d.Partial(false)
	return resourceAliCloudRedisTairInstanceRead(d, meta)
}

func resourceAliCloudRedisTairInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v == "Subscription" {
			log.Printf("[WARN] Cannot destroy resource alicloud_redis_tair_instance which payment_type valued Subscription. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewRedisClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &runtime)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertRedisInstancesDBInstanceAttributeChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "PrePaid":
		return "Subscription"
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
func convertRedisInstancesDBInstanceAttributeStorageTypeResponse(source interface{}) interface{} {
	switch source {
	case "essd_pl1":
		return "PL1"
	case "essd_pl2":
		return "PL2"
	case "essd_pl3":
		return "PL3"
	}
	return source
}
func convertRedisChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
func convertRedisStorageTypeRequest(source interface{}) interface{} {
	switch source {
	case "PL1":
		return "essd_pl1"
	case "PL2":
		return "essd_pl2"
	case "PL3":
		return "essd_pl3"
	}
	return source
}

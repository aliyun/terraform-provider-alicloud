package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMongodbShardingNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongodbShardingNodeCreate,
		Read:   resourceAlicloudMongodbShardingNodeRead,
		Update: resourceAlicloudMongodbShardingNodeUpdate,
		Delete: resourceAlicloudMongodbShardingNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_storage": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(10, 2000),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudMongodbShardingNodeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNode"
	request := make(map[string]interface{})
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	request["AutoPay"] = true
	request["DBInstanceId"] = d.Get("db_instance_id")
	request["NodeClass"] = d.Get("node_class")
	request["NodeStorage"] = d.Get("node_storage")
	request["NodeType"] = "shard"
	request["ClientToken"] = buildClientToken("CreateNode")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_sharding_node", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", response["NodeId"]))

	MongoDBService := MongoDBService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, MongoDBService.MongodbShardingNodeStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudMongodbShardingNodeRead(d, meta)
}
func resourceAlicloudMongodbShardingNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	MongoDBService := MongoDBService{client}
	object, err := MongoDBService.DescribeMongodbShardingNode(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_sharding_node MongoDBService.DescribeMongodbShardingNode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("db_instance_id", parts[0])
	d.Set("node_id", parts[1])
	d.Set("node_class", object["NodeClass"])
	if v, ok := object["NodeStorage"]; ok && fmt.Sprint(v) != "0" {
		d.Set("node_storage", formatInt(v))
	}
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudMongodbShardingNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	MongoDBService := MongoDBService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"NodeId":       parts[1],
	}
	request["AutoPay"] = true
	request["NodeClass"] = d.Get("node_class")
	if d.HasChange("node_class") {
		update = true
	}
	if d.HasChange("node_storage") {
		update = true
	}
	request["NodeStorage"] = d.Get("node_storage")
	if update {
		request["EffectiveTime"] = "Immediately"
		request["FromApp"] = "OpenApi"
		action := "ModifyNodeSpec"
		conn, err := client.NewDdsClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ModifyNodeSpec")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &runtime)
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
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, MongoDBService.MongodbShardingNodeStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudMongodbShardingNodeRead(d, meta)
}
func resourceAlicloudMongodbShardingNodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	MongoDBService := MongoDBService{client}
	action := "DeleteNode"
	var response map[string]interface{}
	conn, err := client.NewDdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"NodeId":       parts[1],
	}

	request["ClientToken"] = buildClientToken("DeleteNode")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-12-01"), StringPointer("AK"), nil, request, &runtime)
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, MongoDBService.MongodbShardingNodeStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

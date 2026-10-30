// Package mongodb. This file is generated automatically. Please do not modify it manually, thank you!
package mongodb

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudMongodbNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongodbNodeCreate,
		Read:   resourceAliCloudMongodbNodeRead,
		Update: resourceAliCloudMongodbNodeUpdate,
		Delete: resourceAliCloudMongodbNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"business_info": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"effective_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Immediately", "MaintainTime"}, false),
			},
			"from_app": {
				Type:     schema.TypeString,
				Optional: true,
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
				Optional:     true,
				ValidateFunc: validators.IntBetween(10, 2000),
			},
			"node_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"order_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"readonly_replicas": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.IntBetween(0, 5),
			},
			"shard_direct": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudMongodbNodeCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateNode"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("business_info"); ok {
		request["BusinessInfo"] = v
	}
	if v, ok := d.GetOk("account_name"); ok {
		request["AccountName"] = v
	}
	if v, ok := d.GetOkExists("shard_direct"); ok {
		request["ShardDirect"] = v
	}
	if v, ok := d.GetOkExists("node_storage"); ok && v.(int) > 0 {
		request["NodeStorage"] = v
	}
	if v, ok := d.GetOkExists("readonly_replicas"); ok {
		request["ReadonlyReplicas"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	request["NodeType"] = d.Get("node_type")
	if v, ok := d.GetOk("account_password"); ok {
		request["AccountPassword"] = v
	}
	request["NodeClass"] = d.Get("node_class")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_mongodb_node", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], response["NodeId"]))

	mongodbServiceV2 := NewMongodbServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 20*time.Minute, mongodbServiceV2.MongodbNodeStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudMongodbNodeRead(d, meta)
}

func resourceAliCloudMongodbNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := NewMongodbServiceV2(client)

	objectRaw, err := mongodbServiceV2.DescribeMongodbNode(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_node DescribeMongodbNode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("node_class", objectRaw["NodeClass"])
	d.Set("node_storage", objectRaw["NodeStorage"])
	d.Set("node_type", objectRaw["NodeType"])
	d.Set("readonly_replicas", objectRaw["ReadonlyReplicas"])
	d.Set("status", objectRaw["Status"])
	d.Set("db_instance_id", objectRaw["DbInstanceId"])
	d.Set("node_id", objectRaw["NodeId"])

	return nil
}

func resourceAliCloudMongodbNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyNodeSpec"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["NodeId"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if v, ok := d.GetOk("business_info"); ok {
		request["BusinessInfo"] = v
	}
	if d.HasChange("node_storage") {
		update = true
		request["NodeStorage"] = d.Get("node_storage")
	}

	if d.HasChange("readonly_replicas") {
		update = true
		request["ReadonlyReplicas"] = d.Get("readonly_replicas")
	}

	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("order_type"); ok {
		request["OrderType"] = v
	}
	if v, ok := d.GetOk("switch_time"); ok {
		request["SwitchTime"] = v
	}
	if v, ok := d.GetOk("from_app"); ok {
		request["FromApp"] = v
	}
	if v, ok := d.GetOk("effective_time"); ok {
		request["EffectiveTime"] = v
	}
	if d.HasChange("node_class") {
		update = true
	}
	request["NodeClass"] = d.Get("node_class")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || helper.NeedRetry(err) {
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
		mongodbServiceV2 := NewMongodbServiceV2(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 30*time.Minute, mongodbServiceV2.MongodbNodeStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudMongodbNodeRead(d, meta)
}

func resourceAliCloudMongodbNodeDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteNode"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["NodeId"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
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

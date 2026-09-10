package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAmqpExchange() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAmqpExchangeCreate,
		Read:   resourceAliCloudAmqpExchangeRead,
		Update: resourceAliCloudAmqpExchangeUpdate,
		Delete: resourceAliCloudAmqpExchangeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alternate_exchange": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_delete_state": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"exchange_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"exchange_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"FANOUT", "DIRECT", "TOPIC", "HEADERS", "X_DELAYED_MESSAGE", "X_CONSISTENT_HASH"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internal": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"virtual_host_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"x_delayed_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"DIRECT", "TOPIC", "FANOUT", "HEADERS", "X_JMS_TOPIC"}, false),
			},
		},
	}
}

func resourceAliCloudAmqpExchangeCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateExchange"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("exchange_name"); ok {
		request["ExchangeName"] = v
	}
	if v, ok := d.GetOk("virtual_host_name"); ok {
		request["VirtualHost"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("x_delayed_type"); ok {
		request["XDelayedType"] = v
	}
	if v, ok := d.GetOk("alternate_exchange"); ok {
		request["AlternateExchange"] = v
	}
	request["ExchangeType"] = d.Get("exchange_type")
	request["AutoDeleteState"] = d.Get("auto_delete_state")
	request["Internal"] = d.Get("internal")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_amqp_exchange", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["InstanceId"], request["VirtualHost"], request["ExchangeName"]))

	return resourceAliCloudAmqpExchangeRead(d, meta)
}

func resourceAliCloudAmqpExchangeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpServiceV2 := AmqpServiceV2{client}

	objectRaw, err := amqpServiceV2.DescribeAmqpExchange(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_amqp_exchange DescribeAmqpExchange Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auto_delete_state", objectRaw["AutoDeleteState"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("exchange_type", objectRaw["ExchangeType"])
	d.Set("exchange_name", objectRaw["Name"])
	d.Set("virtual_host_name", objectRaw["VHostName"])

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudAmqpExchangeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Exchange.")
	return nil
}

func resourceAliCloudAmqpExchangeDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteExchange"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["ExchangeName"] = parts[2]
	request["VirtualHost"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, query, request, true)

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

	return nil
}

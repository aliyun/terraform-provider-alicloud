package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAmqpBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAmqpBindingCreate,
		Read:   resourceAliCloudAmqpBindingRead,
		Delete: resourceAliCloudAmqpBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_host_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_exchange": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"binding_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"EXCHANGE", "QUEUE"}, false),
			},
			"binding_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"argument": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"x-match:all", "x-match:any"}, false),
			},
		},
	}
}

func resourceAliCloudAmqpBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBinding"
	request := make(map[string]interface{})
	var err error

	request["InstanceId"] = d.Get("instance_id")
	request["VirtualHost"] = d.Get("virtual_host_name")
	request["SourceExchange"] = d.Get("source_exchange")
	request["DestinationName"] = d.Get("destination_name")
	request["BindingType"] = d.Get("binding_type")
	request["BindingKey"] = d.Get("binding_key")

	if v, ok := d.GetOk("argument"); ok {
		request["Argument"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_amqp_binding", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", request["InstanceId"], request["VirtualHost"], request["SourceExchange"], request["DestinationName"]))

	return resourceAliCloudAmqpBindingRead(d, meta)
}

func resourceAliCloudAmqpBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpOpenService := AmqpOpenService{client}

	object, err := amqpOpenService.DescribeAmqpBinding(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_amqp_binding amqpOpenService.DescribeAmqpBinding Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("virtual_host_name", parts[1])
	d.Set("source_exchange", object["SourceExchange"])
	d.Set("destination_name", object["DestinationName"])
	d.Set("binding_type", object["BindingType"])
	d.Set("binding_key", object["BindingKey"])
	d.Set("argument", object["Argument"])

	return nil
}

func resourceAliCloudAmqpBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBinding"
	var response map[string]interface{}
	var err error

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId":      parts[0],
		"VirtualHost":     parts[1],
		"SourceExchange":  parts[2],
		"DestinationName": parts[3],
		"BindingType":     d.Get("binding_type"),
		"BindingKey":      d.Get("binding_key"),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
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
		if NotFoundError(err) || IsExpectedErrors(err, []string{"ExchangeNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

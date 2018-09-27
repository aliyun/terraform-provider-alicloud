package alicloud

import (
	"fmt"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudMNSQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMNSQueueCreate,
		Read:   resourceAlicloudMNSQueueRead,
		Update: resourceAlicloudMNSQueueUpdate,
		Delete: resourceAlicloudMNSQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delay_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"maximum_message_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  65536,
			},
			"message_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  259200,
			},
			"visibility_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"polling_wait_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceAlicloudMNSQueueCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf(" creating alicoudMNSQueue  error: %#v", err)
	}
	queueManager := ali_mns.NewMNSQueueManager(*mnsClient)
	name := d.Get("name").(string)
	var delaySeconds, maximumMessageSize, messageRetentionPeriod, visibilityTimeout, pollingWaitSeconds int
	delaySeconds = 0
	maximumMessageSize = 65536
	messageRetentionPeriod = 259200
	visibilityTimeout = 30
	pollingWaitSeconds = 0
	if v, ok := d.GetOk("delay_seconds"); ok {
		delaySeconds = v.(int)
	}
	if v, ok := d.GetOk("maximum_message_size"); ok {
		maximumMessageSize = v.(int)
	}
	if v, ok := d.GetOk("message_retention_period"); ok {
		messageRetentionPeriod = v.(int)
	}
	if v, ok := d.GetOk("visibility_timeout"); ok {
		visibilityTimeout = v.(int)
	}
	if v, ok := d.GetOk("polling_wait_seconds"); ok {
		pollingWaitSeconds = v.(int)
	}

	err = queueManager.CreateQueue(name, int32(delaySeconds), int32(maximumMessageSize), int32(messageRetentionPeriod), int32(visibilityTimeout), int32(pollingWaitSeconds), 3)
	if err != nil {
		return fmt.Errorf("Create queue got an error: %#v", err)
	}
	d.SetId(name)
	return nil
}

func resourceAlicloudMNSQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf(" creating alicoudMNSQueue  error: %#v", err)
	}
	queueManager := ali_mns.NewMNSQueueManager(*mnsClient)
	attr, err := queueManager.GetQueueAttributes(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", attr.QueueName)
	d.Set("delay_seconds", attr.DelaySeconds)
	d.Set("maximum_message_size", attr.MaxMessageSize)
	d.Set("message_retention_period", attr.MessageRetentionPeriod)
	d.Set("visibility_timeout", attr.VisibilityTimeout)
	d.Set("polling_wait_seconds", attr.PollingWaitSeconds)

	return nil
}

func resourceAlicloudMNSQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf(" creating alicoudMNSQueue  error: %#v", err)
	}
	queueManager := ali_mns.NewMNSQueueManager(*mnsClient)
	d.Partial(true)
	attributeUpdate := false

	attr, err1 := queueManager.GetQueueAttributes(d.Get("name").(string))
	if err1 != nil {
		return err1
	}
	var delaySeconds, maximumMessageSize, messageRetentionPeriod, visibilityTimeouts, pollingWaitSeconds int
	delaySeconds = int(attr.DelaySeconds)
	maximumMessageSize = int(attr.MaxMessageSize)
	messageRetentionPeriod = int(attr.MessageRetentionPeriod)
	visibilityTimeouts = int(attr.VisibilityTimeout)
	pollingWaitSeconds = int(attr.PollingWaitSeconds)
	name := d.Get("name").(string)
	if d.HasChange("delay_seconds") {
		d.SetPartial("delay_seconds")
		delaySeconds = d.Get("delay_seconds").(int)
		attributeUpdate = true
	}

	if d.HasChange("maximum_message_size") {
		d.SetPartial("maximum_message_size")
		maximumMessageSize = d.Get("maximum_message_size").(int)
		attributeUpdate = true
	}

	if d.HasChange("message_retention_period") {
		d.SetPartial("message_retention_period")
		messageRetentionPeriod = d.Get("message_retention_period").(int)
		attributeUpdate = true
	}
	if d.HasChange("visibility_timeout") {
		d.SetPartial("visibility_timeout")
		visibilityTimeouts = d.Get("visibility_timeout").(int)
		attributeUpdate = true
	}
	if d.HasChange("polling_wait_seconds") {
		d.SetPartial("polling_wait_seconds")
		pollingWaitSeconds = d.Get("polling_wait_seconds").(int)
		attributeUpdate = true
	}
	d.Partial(false)

	if attributeUpdate {
		err = queueManager.SetQueueAttributes(name, int32(delaySeconds), int32(maximumMessageSize), int32(messageRetentionPeriod), int32(visibilityTimeouts), int32(pollingWaitSeconds), 3)

		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlicloudMNSQueueDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf(" creating alicoudMNSQueue  error: %#v", err)
	}
	queueManager := ali_mns.NewMNSQueueManager(*mnsClient)
	name := d.Get("name").(string)
	err = queueManager.DeleteQueue(name)
	if err != nil {
		return err
	}
	return nil
}

package alicloud

import (
	"fmt"
	"strings"

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(3, 256),
			},
			"delay_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntegerInRange(0, 604800),
			},
			"maximum_message_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      65536,
				ValidateFunc: validateIntegerInRange(1024, 65536),
			},
			"message_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      345600,
				ValidateFunc: validateIntegerInRange(60, 604800),
			},
			"visibility_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validateIntegerInRange(1, 43200),
			},
			"polling_wait_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntegerInRange(0, 1800),
			},
		},
	}
}

func resourceAlicloudMNSQueueCreate(d *schema.ResourceData, meta interface{}) error {
	queueManager, err := getQueueManager(meta)
	if err != nil {
		return err
	}
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
	return resourceAlicloudMNSQueueRead(d, meta)
}

func resourceAlicloudMNSQueueRead(d *schema.ResourceData, meta interface{}) error {
	queueManager, err := getQueueManager(meta)
	if err != nil {
		return err
	}
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
	queueManager, err := getQueueManager(meta)
	if err != nil {
		return err
	}
	attributeUpdate := false
	var delaySeconds, maximumMessageSize, messageRetentionPeriod, visibilityTimeouts, pollingWaitSeconds int
	delaySeconds = d.Get("delay_seconds").(int)
	maximumMessageSize = d.Get("maximum_message_size").(int)
	messageRetentionPeriod = d.Get("message_retention_period").(int)
	visibilityTimeouts = d.Get("visibility_timeout").(int)
	pollingWaitSeconds = d.Get("polling_wait_seconds").(int)
	name := d.Get("name").(string)
	if d.HasChange("delay_seconds") {
		attributeUpdate = true
	}

	if d.HasChange("maximum_message_size") {
		attributeUpdate = true
	}

	if d.HasChange("message_retention_period") {
		attributeUpdate = true
	}
	if d.HasChange("visibility_timeout") {
		attributeUpdate = true
	}
	if d.HasChange("polling_wait_seconds") {
		attributeUpdate = true
	}

	if attributeUpdate {
		err = queueManager.SetQueueAttributes(name, int32(delaySeconds), int32(maximumMessageSize), int32(messageRetentionPeriod), int32(visibilityTimeouts), int32(pollingWaitSeconds), 3)
		if err != nil {
			return err
		}
	}
	return resourceAlicloudMNSQueueRead(d, meta)
}

func resourceAlicloudMNSQueueDelete(d *schema.ResourceData, meta interface{}) error {
	queueManager, err := getQueueManager(meta)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	err = queueManager.DeleteQueue(name)
	if err != nil {
		return err
	}
	err = resourceAlicloudMNSQueueRead(d, meta)
	if err != nil && strings.Contains(err.Error(), "QueueNotExist") {
		return nil
	}
	return err
}

func getQueueManager(meta interface{}) (ali_mns.AliQueueManager, error) {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return nil, fmt.Errorf(" creating MNSQueueManager  error: %#v", err)
	}
	queueManager := ali_mns.NewMNSQueueManager(*mnsClient)
	return queueManager, nil
}

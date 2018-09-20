package alicloud

import (
	"fmt"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudMNSTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMNSTopicCreate,
		Read:   resourceAlicloudMNSTopicRead,
		Update: resourceAlicloudMNSTopicUpdate,
		Delete: resourceAlicloudMNSTopicDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"maximum_message_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  65536,
			},

			"logging_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlicloudMNSTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf(" creating alicoudMNSTopic  error: %#v", err)
	}
	topicManager := ali_mns.NewMNSTopicManager(*mnsClient)
	name := d.Get("name").(string)
	var maximumMessageSize int
	maximumMessageSize = 65536
	var logginEnabled bool

	if v, ok := d.GetOk("maximum_message_size"); ok {
		maximumMessageSize = v.(int)
	}
	if v, ok := d.GetOk("logging_enabled"); ok {
		logginEnabled = v.(bool)
	}

	err = topicManager.CreateTopic(name, int32(maximumMessageSize), logginEnabled)
	if err != nil {
		return fmt.Errorf("Create topic got an error: %#v", err)
	}
	d.SetId(name)
	return nil
}

func resourceAlicloudMNSTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf(" creating alicoudMNSTopic  error: %#v", err)
	}
	topicManager := ali_mns.NewMNSTopicManager(*mnsClient)
	attr, err := topicManager.GetTopicAttributes(d.Get("name").(string))
	if err != nil {
		return err
	}
	d.Set("name", attr.TopicName)
	d.Set("maximum_message_size", attr.MaxMessageSize)
	d.Set("logging_enabled", attr.LoggingEnabled)

	return nil
}

func resourceAlicloudMNSTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf(" creating alicoudMNSTopic  error: %#v", err)
	}
	topicManager := ali_mns.NewMNSTopicManager(*mnsClient)
	d.Partial(true)
	attributeUpdate := false

	var maximumMessageSize int
	var loggingEnabled bool

	name := d.Get("name").(string)

	if d.HasChange("maximum_message_size") {
		d.SetPartial("maximum_message_size")
		maximumMessageSize = d.Get("maximum_message_size").(int)
		attributeUpdate = true
	}

	if d.HasChange("logging_enabled") {
		d.SetPartial("logging_enabled")
		loggingEnabled = d.Get("logging_enabled").(bool)
		attributeUpdate = true
	}
	d.Partial(false)

	if attributeUpdate {
		err = topicManager.SetTopicAttributes(name, int32(maximumMessageSize), loggingEnabled)

		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlicloudMNSTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf(" creating alicoudMNSTopic  error: %#v", err)
	}
	topicManager := ali_mns.NewMNSTopicManager(*mnsClient)
	name := d.Get("name").(string)
	err = topicManager.DeleteTopic(name)
	if err != nil {
		return err
	}
	return nil
}

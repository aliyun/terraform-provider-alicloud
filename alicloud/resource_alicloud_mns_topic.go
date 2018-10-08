package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudMNSTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMNSTopicCreate,
		Read:   resourceAlicloudMNSTopicRead,
		Update: resourceAlicloudMNSTopicUpdate,
		Delete: resourceAlicloudMNSTopicDelete,
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

			"maximum_message_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      65536,
				ValidateFunc: validateIntegerInRange(1024, 65536),
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
	topicManager, err := client.MnsTopicManager()
	if err != nil {
		return fmt.Errorf("Creating alicoudMNSTopicManager error: %#v", err)
	}
	name := d.Get("name").(string)
	maximumMessageSize := d.Get("maximum_message_size").(int)
	loggingEnabled := d.Get("logging_enabled").(bool)
	err = topicManager.CreateTopic(name, int32(maximumMessageSize), loggingEnabled)
	if err != nil {
		return fmt.Errorf("Create topic got an error: %#v", err)
	}
	d.SetId(name)
	return resourceAlicloudMNSTopicRead(d, meta)
}

func resourceAlicloudMNSTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	topicManager, err := client.MnsTopicManager()
	if err != nil {
		return fmt.Errorf("Creating alicoudMNSTopicManager error: %#v", err)
	}
	attr, err := topicManager.GetTopicAttributes(d.Id())
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
	topicManager, err := client.MnsTopicManager()
	if err != nil {
		return fmt.Errorf(" Creating alicoudMNSTopicManager  error: %#v", err)
	}
	attributeUpdate := false
	name := d.Id()
	maximumMessageSize := d.Get("maximum_message_size").(int)
	loggingEnabled := d.Get("logging_enabled").(bool)
	if d.HasChange("maximum_message_size") {
		attributeUpdate = true
	}

	if d.HasChange("logging_enabled") {
		attributeUpdate = true
	}

	if attributeUpdate {
		err = topicManager.SetTopicAttributes(name, int32(maximumMessageSize), loggingEnabled)
		if err != nil {
			return err
		}
	}
	return resourceAlicloudMNSTopicRead(d, meta)
}

func resourceAlicloudMNSTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	topicManager, err := client.MnsTopicManager()
	if err != nil {
		return fmt.Errorf("Creating alicoudMNSTopicManager error: %#v", err)
	}
	name := d.Id()
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		err = topicManager.DeleteTopic(name)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting mns topic got an error: %#v", err))
		}

		attr, err := topicManager.GetTopicAttributes(name)
		if err != nil {
			if TopicNotExistFunc(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe mns topic got an error: %#v", err))
		}
		if attr.TopicName != name {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting mns topic got an error: %#v", err))
	})
}

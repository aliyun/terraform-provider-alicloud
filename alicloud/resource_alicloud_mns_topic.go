package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
	client := meta.(*connectivity.AliyunClient)
	name := d.Get("name").(string)
	maximumMessageSize := d.Get("maximum_message_size").(int)
	loggingEnabled := d.Get("logging_enabled").(bool)
	_, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
		return nil, topicManager.CreateTopic(name, int32(maximumMessageSize), loggingEnabled)
	})
	if err != nil {
		return fmt.Errorf("Create topic got an error: %#v", err)
	}
	d.SetId(name)
	return resourceAlicloudMNSTopicRead(d, meta)
}

func resourceAlicloudMNSTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	raw, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
		return topicManager.GetTopicAttributes(d.Id())
	})
	mnsService := MnsService{}
	if err != nil {
		if mnsService.TopicNotExistFunc(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Read mns topic error: %#v", err)
	}
	attr, _ := raw.(ali_mns.TopicAttribute)
	d.Set("name", attr.TopicName)
	d.Set("maximum_message_size", attr.MaxMessageSize)
	d.Set("logging_enabled", attr.LoggingEnabled)
	return nil
}

func resourceAlicloudMNSTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
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
		_, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
			return nil, topicManager.SetTopicAttributes(name, int32(maximumMessageSize), loggingEnabled)
		})
		if err != nil {
			return err
		}
	}
	return resourceAlicloudMNSTopicRead(d, meta)
}

func resourceAlicloudMNSTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mnsService := MnsService{}
	name := d.Id()
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
			return nil, topicManager.DeleteTopic(name)
		})

		if err != nil {
			if mnsService.TopicNotExistFunc(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting mns topic got an error: %#v", err))
		}

		raw, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
			return topicManager.GetTopicAttributes(name)
		})
		if err != nil {
			if mnsService.TopicNotExistFunc(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe mns topic got an error: %#v", err))
		}
		attr, _ := raw.(ali_mns.TopicAttribute)
		if attr.TopicName != name {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting mns topic got an error: %#v", err))
	})
}

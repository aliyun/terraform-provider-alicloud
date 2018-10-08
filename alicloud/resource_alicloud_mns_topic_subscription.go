package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudMNSSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMNSSubscriptionCreate,
		Read:   resourceAlicloudMNSSubscriptionRead,
		Update: resourceAlicloudMNSSubscriptionUpdate,
		Delete: resourceAlicloudMNSSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(3, 256),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(3, 256),
			},

			"endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateEndpoint,
			},

			"filter_tag": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(0, 16),
			},

			"notify_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ali_mns.BACKOFF_RETRY),
				ValidateFunc: validateAllowedStringValue([]string{
					string(ali_mns.BACKOFF_RETRY), string(ali_mns.EXPONENTIAL_DECAY_RETRY),
				}),
			},

			"notify_content_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ali_mns.SIMPLIFIED),
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ali_mns.SIMPLIFIED), string(ali_mns.XML),
				}),
			},
		},
	}
}

func resourceAlicloudMNSSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	topicName := d.Get("topic_name").(string)
	subscriptionManager, err := client.MnsSubscriptionManager(topicName)
	if err != nil {
		return fmt.Errorf("Creating mns subscription client  error: %#v", err)
	}
	name := d.Get("name").(string)
	endpoint := d.Get("endpoint").(string)
	notifyStrategyStr := d.Get("notify_strategy").(string)
	notifyContentFormatStr := d.Get("notify_content_format").(string)
	var filterTag string
	if v, ok := d.GetOk("filter_tag"); ok {
		filterTag = v.(string)
	}
	notifyStrategy := ali_mns.NotifyStrategyType(notifyStrategyStr)
	notifyContentFormat := ali_mns.NotifyContentFormatType(notifyContentFormatStr)
	subRequest := ali_mns.MessageSubsribeRequest{
		Endpoint:            endpoint,
		FilterTag:           filterTag,
		NotifyStrategy:      notifyStrategy,
		NotifyContentFormat: notifyContentFormat,
	}
	err = subscriptionManager.Subscribe(name, subRequest)
	if err != nil {
		return fmt.Errorf("Create Subscription got an error: %#v", err)
	}
	d.SetId(fmt.Sprintf("%s%s%s", topicName, COLON_SEPARATED, name))
	return resourceAlicloudMNSSubscriptionRead(d, meta)
}

func resourceAlicloudMNSSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	topicName, name := GetTopicNameAndSubscriptionName(d.Id())
	subscriptionManager, err := client.MnsSubscriptionManager(topicName)
	if err != nil {
		return fmt.Errorf("Creating mns subscription client  error: %#v", err)
	}
	attr, err := subscriptionManager.GetSubscriptionAttributes(name)
	if err != nil {
		if SubscriptionNotExistFunc(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Get mns subscription attr   error: %#v", err)
	}
	d.Set("topic_name", attr.TopicName)
	d.Set("name", attr.SubscriptionName)
	d.Set("endpoint", attr.Endpoint)
	d.Set("filter_tag", attr.FilterTag)
	d.Set("notify_strategy", attr.NotifyStrategy)
	d.Set("notify_content_format", attr.NotifyContentFormat)
	return nil
}

func resourceAlicloudMNSSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("notify_strategy") {
		client := meta.(*AliyunClient)
		topicName, name := GetTopicNameAndSubscriptionName(d.Id())
		subscriptionManager, err := client.MnsSubscriptionManager(topicName)
		if err != nil {
			return fmt.Errorf("Creating mns subscription client  error: %#v", err)
		}
		notifyStrategy := ali_mns.NotifyStrategyType(d.Get("notify_strategy").(string))
		err = subscriptionManager.SetSubscriptionAttributes(name, notifyStrategy)
		if err != nil {
			return fmt.Errorf("update mns subscription client  error: %#v", err)
		}
	}
	return resourceAlicloudMNSSubscriptionRead(d, meta)
}

func resourceAlicloudMNSSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	topicName, name := GetTopicNameAndSubscriptionName(d.Id())
	subscriptionManager, err := client.MnsSubscriptionManager(topicName)
	if err != nil {
		return fmt.Errorf("Creating mns subscription client  error: %#v", err)
	}
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		err = subscriptionManager.Unsubscribe(name)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting mns subscription %s got an error: %#v", name, err))
		}
		attr, err := subscriptionManager.GetSubscriptionAttributes(name)
		if err != nil {
			if SubscriptionNotExistFunc(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe mns subscription %s got an error: %#v", name, err))
		}
		if attr.SubscriptionName != name {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting mns subscription %s timeout.", name))
	})

}

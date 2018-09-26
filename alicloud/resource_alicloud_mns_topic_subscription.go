package alicloud

import (
	"fmt"
	"strings"

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"filter_tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"notify_strategy": {
				Type:     schema.TypeString,
				Required: true,
			},

			"notify_content_format": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudMNSSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf(" creating alicoudMNSSubscription  error: %#v", err)
	}
	topicName := d.Get("topic_name").(string)
	subscriptionManager := ali_mns.NewMNSTopic(topicName, *mnsClient)

	name := d.Get("name").(string)
	endpoint := d.Get("endpoint").(string)
	notifyStrategyStr := d.Get("notify_strategy").(string)
	notifyContentFormatStr := d.Get("notify_content_format").(string)
	var filterFlag string
	if v, ok := d.GetOk("filter_tag"); ok {
		filterFlag = v.(string)
	}
	var notifyStrategy ali_mns.NotifyStrategyType
	switch notifyStrategyStr {
	case "BACKOFF_RETRY":
		notifyStrategy = ali_mns.BACKOFF_RETRY
	case "EXPONENTIAL_DECAY_RETRY":
		notifyStrategy = ali_mns.EXPONENTIAL_DECAY_RETRY
	default:
		panic("please input BACKOFF_RETRY or EXPONENTIAL_DECAY_RETRY for notify_strategy")

	}
	var notifyContentFormat ali_mns.NotifyContentFormatType
	switch notifyContentFormatStr {
	case "XML":
		notifyContentFormat = ali_mns.XML
	case "SIMPLIFIED":
		notifyContentFormat = ali_mns.SIMPLIFIED
	default:
		panic("please input BACKOFF_RETRY or EXPONENTIAL_DECAY_RETRY for notify_content_format")

	}
	subRequest := ali_mns.MessageSubsribeRequest{
		Endpoint:            endpoint,
		FilterTag:           filterFlag,
		NotifyStrategy:      notifyStrategy,
		NotifyContentFormat: notifyContentFormat,
	}
	err = subscriptionManager.Subscribe(name, subRequest)
	if err != nil {
		return fmt.Errorf("Create Subscription got an error: %#v", err)
	}
	d.SetId(topicName + "#" + name)
	return nil
}

func resourceAlicloudMNSSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf(" creating alicoudMNSSubscription  error: %#v", err)
	}
	arr := strings.Split(d.Id(), "#")

	subscriptionManager := ali_mns.NewMNSTopic(arr[0], *mnsClient)


	attr, err1 := subscriptionManager.GetSubscriptionAttributes(arr[1])
	if err1 != nil {
		return err1
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
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf("creating mnsclient  error: %#v", err)
	}
	topicName := d.Get("topic_name").(string)
	subscriptionManager := ali_mns.NewMNSTopic(topicName, *mnsClient)
	name := d.Get("name").(string)
	attributeUpdate := false

	var notifyStrategy ali_mns.NotifyStrategyType

	if d.HasChange("notify_strategy") {
		d.SetPartial("notify_strategy")
		notifyStrategyStr := d.Get("notify_strategy").(string)
		switch notifyStrategyStr {
		case "BACKOFF_RETRY":
			notifyStrategy = ali_mns.BACKOFF_RETRY
		case "EXPONENTIAL_DECAY_RETRY":
			notifyStrategy = ali_mns.EXPONENTIAL_DECAY_RETRY
		default:
			panic("please input BACKOFF_RETRY or EXPONENTIAL_DECAY_RETRY for notify_strategy")

		}
		attributeUpdate = true
	}

	if attributeUpdate {
		err = subscriptionManager.SetSubscriptionAttributes(name, notifyStrategy)

		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlicloudMNSSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	mnsClient, err := client.Mnsconn()
	if err != nil {
		return fmt.Errorf("creating mnsclient  error: %#v", err)
	}
	topicName := d.Get("topic_name").(string)
	subscriptionManager := ali_mns.NewMNSTopic(topicName, *mnsClient)
	name := d.Get("name").(string)
	err = subscriptionManager.Unsubscribe(name)
	if err != nil {
		return err
	}
	return nil
}

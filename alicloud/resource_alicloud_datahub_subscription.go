package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudDatahubSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDatahubSubscriptionCreate,
		Read:   resourceAliyunDatahubSubscriptionRead,
		Update: resourceAliyunDatahubSubscriptionUpdate,
		Delete: resourceAliyunDatahubSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDatahubProjectName,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"topic_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDatahubTopicName,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"comment": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "subscription added by terraform",
				ValidateFunc: validateStringLengthInRange(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"sub_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
		},
	}
}

func resourceAliyunDatahubSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	dh := meta.(*AliyunClient).dhconn

	projectName := d.Get("project_name").(string)
	topicName := d.Get("topic_name").(string)
	subComment := d.Get("comment").(string)

	subId, err := dh.CreateSubscription(projectName, topicName, subComment)
	if err != nil {
		return fmt.Errorf("failed to create subscription under '%s/%s' with error: %s", projectName, topicName, err)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", strings.ToLower(projectName), COLON_SEPARATED, strings.ToLower(topicName), COLON_SEPARATED, subId))
	return resourceAliyunDatahubSubscriptionRead(d, meta)
}

func parseId3(d *schema.ResourceData, meta interface{}) (projectName, topicName, subId string, err error) {
	split := strings.Split(d.Id(), COLON_SEPARATED)
	if len(split) != 3 {
		err = fmt.Errorf("you should use resource alicloud_datahub_subscription's new field 'project_name' and 'topic_name' to re-import this resource.")
		return
	} else {
		projectName = split[0]
		topicName = split[1]
		subId = split[2]
	}
	return
}

func resourceAliyunDatahubSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, subId, err := parseId3(d, meta)
	if err != nil {
		return err
	}

	dh := meta.(*AliyunClient).dhconn

	sub, err := dh.GetSubscription(projectName, topicName, subId)
	if err != nil {
		if isDatahubNotExistError(err) {
			d.SetId("")
		}
		return fmt.Errorf("failed to get subscription %s with error: %s", subId, err)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", strings.ToLower(projectName), COLON_SEPARATED, strings.ToLower(sub.TopicName), COLON_SEPARATED, sub.SubId))

	d.Set("project_name", projectName)
	d.Set("topic_name", sub.TopicName)
	d.Set("sub_id", sub.SubId)
	d.Set("comment", sub.Comment)
	d.Set("create_time", datahub.Uint64ToTimeString(sub.CreateTime))
	d.Set("last_modify_time", datahub.Uint64ToTimeString(sub.LastModifyTime))
	return nil
}

func resourceAliyunDatahubSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, subId, err := parseId3(d, meta)
	if err != nil {
		return err
	}

	dh := meta.(*AliyunClient).dhconn

	if d.HasChange("comment") {
		subComment := d.Get("comment").(string)

		err := dh.UpdateSubscription(projectName, topicName, subId, subComment)
		if err != nil {
			return fmt.Errorf("failed to update subscription %s's comment with error: %s", subId, err)
		}
	}

	return resourceAliyunDatahubSubscriptionRead(d, meta)
}

func resourceAliyunDatahubSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, subId, err := parseId3(d, meta)
	if err != nil {
		return err
	}

	dh := meta.(*AliyunClient).dhconn

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := dh.GetSubscription(projectName, topicName, subId)
		if err != nil {
			if isDatahubNotExistError(err) {
				return nil
			}
			if isRetryableDatahubError(err) {
				return resource.RetryableError(fmt.Errorf("while deleting subscription '%s', failed to get it with error: %s", subId, err))
			}
			return resource.NonRetryableError(fmt.Errorf("while deleting subscription '%s', failed to get it with error: %s", subId, err))
		}

		err = dh.DeleteSubscription(projectName, topicName, subId)
		if err == nil || isDatahubNotExistError(err) {
			return nil
		}

		if isRetryableDatahubError(err) {
			return resource.RetryableError(fmt.Errorf("Deleting subscription '%s' timeout and got an error: %#v.", subId, err))
		}
		return resource.NonRetryableError(fmt.Errorf("Deleting subscription '%s' timeout.", subId))
	})
}

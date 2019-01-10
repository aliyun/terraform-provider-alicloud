package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDatahubTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDatahubTopicCreate,
		Read:   resourceAliyunDatahubTopicRead,
		Update: resourceAliyunDatahubTopicUpdate,
		Delete: resourceAliyunDatahubTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDatahubProjectName,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDatahubTopicName,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"shard_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateIntegerInRange(1, 10),
			},
			"life_cycle": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validateIntegerInRange(1, 7),
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "topic added by terraform",
				ValidateFunc: validateStringLengthInRange(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"record_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "TUPLE",
				ValidateFunc: validateAllowedStringValue([]string{string(datahub.TUPLE), string(datahub.BLOB)}),
			},
			"record_schema": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("record_type") != string(datahub.TUPLE)
				},
			},
			"create_time": {
				Type:     schema.TypeString, //converted from UTC(uint64) value
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString, //converted from UTC(uint64) value
				Computed: true,
			},
		},
	}
}

func resourceAliyunDatahubTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	t := &datahub.Topic{
		ProjectName: d.Get("project_name").(string),
		TopicName:   d.Get("name").(string),
		ShardCount:  d.Get("shard_count").(int),
		Lifecycle:   d.Get("life_cycle").(int),
		Comment:     d.Get("comment").(string),
	}

	recordType := d.Get("record_type").(string)
	if recordType == string(datahub.TUPLE) {
		t.RecordType = datahub.TUPLE

		recordSchema := d.Get("record_schema").(map[string]interface{})
		if len(recordSchema) == 0 {
			recordSchema = getDefaultRecordSchemainMap()
		}
		t.RecordSchema = getRecordSchema(recordSchema)
	} else if recordType == string(datahub.BLOB) {
		t.RecordType = datahub.BLOB
	}

	_, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
		return nil, dataHubClient.CreateTopic(t)
	})
	if err != nil {
		return fmt.Errorf("failed to create topic'%s/%s' with error: %s", t.ProjectName, t.TopicName, err)
	}

	d.SetId(strings.ToLower(fmt.Sprintf("%s%s%s", t.ProjectName, COLON_SEPARATED, t.TopicName)))
	return resourceAliyunDatahubTopicRead(d, meta)
}

func parseId2(d *schema.ResourceData, meta interface{}) (projectName, topicName string, err error) {
	split := strings.Split(d.Id(), COLON_SEPARATED)
	if len(split) != 2 {
		err = fmt.Errorf("you should use resource alicloud_datahub_topic's new field 'project_name' and 'name' to re-import this resource.")
		return
	} else {
		projectName = split[0]
		topicName = split[1]
		return
	}
}

func resourceAliyunDatahubTopicRead(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, err := parseId2(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*connectivity.AliyunClient)

	raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
		return dataHubClient.GetTopic(projectName, topicName)
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to access topic '%s/%s' with error: %s", projectName, topicName, err)
	}
	topic, _ := raw.(*datahub.Topic)

	d.SetId(strings.ToLower(fmt.Sprintf("%s%s%s", topic.ProjectName, COLON_SEPARATED, topic.TopicName)))

	d.Set("name", topic.TopicName)
	d.Set("project_name", topic.ProjectName)
	d.Set("shard_count", topic.ShardCount)
	d.Set("life_cycle", topic.Lifecycle)
	d.Set("comment", topic.Comment)
	d.Set("record_type", topic.RecordType.String())
	d.Set("record_schema", topic.RecordSchema.String())
	d.Set("create_time", datahub.Uint64ToTimeString(topic.CreateTime))
	d.Set("last_modify_time", datahub.Uint64ToTimeString(topic.LastModifyTime))
	return nil
}

func resourceAliyunDatahubTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, err := parseId2(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("life_cycle") || d.HasChange("comment") {
		lifeCycle := d.Get("life_cycle").(int)
		topicComment := d.Get("comment").(string)

		_, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			return nil, dataHubClient.UpdateTopic(projectName, topicName, lifeCycle, topicComment)
		})
		if err != nil {
			return fmt.Errorf("failed to update topic '%s/%s' with error: %s", projectName, topicName, err)
		}
	}

	return resourceAliyunDatahubTopicRead(d, meta)
}

func resourceAliyunDatahubTopicDelete(d *schema.ResourceData, meta interface{}) error {
	projectName, topicName, err := parseId2(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*connectivity.AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			return dataHubClient.GetTopic(projectName, topicName)
		})

		if err != nil {
			if isDatahubNotExistError(err) {
				return nil
			}
			if isRetryableDatahubError(err) {
				return resource.RetryableError(fmt.Errorf("while deleting '%s/%s', failed to access it with error: %s", projectName, topicName, err))
			}
			return resource.NonRetryableError(fmt.Errorf("while deleting '%s/%s', failed to access it with error: %s", projectName, topicName, err))
		}

		_, err = client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			return nil, dataHubClient.DeleteTopic(projectName, topicName)
		})
		if err == nil || isDatahubNotExistError(err) {
			return nil
		}

		if isRetryableDatahubError(err) {
			return resource.RetryableError(fmt.Errorf("Deleting topic '%s/%s' timeout and got an error: %#v.", projectName, topicName, err))
		}

		return resource.NonRetryableError(fmt.Errorf("Deleting topic '%s/%s' timeout and got an error: %#v.", projectName, topicName, err))
	})
}

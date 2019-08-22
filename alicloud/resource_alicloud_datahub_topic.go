package alicloud

import (
	"fmt"

	"strconv"
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
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validateIntegerInRange(1, 10),
			},
			"life_cycle": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validateIntegerInRange(1, 7),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(old) != "" && strings.ToLower(new) != strings.ToLower(old)
				},
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
				ForceNew:     true,
				Default:      "TUPLE",
				ValidateFunc: validateAllowedStringValue([]string{string(datahub.TUPLE), string(datahub.BLOB)}),
			},
			"record_schema": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

	var requestInfo *datahub.DataHub

	raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
		requestInfo = dataHubClient
		return nil, dataHubClient.CreateTopic(t)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_datahub_topic", "CreateTopic", AliyunDatahubSdkGo)
	}
	addDebug("CreateTopic", raw, requestInfo, t)

	d.SetId(strings.ToLower(fmt.Sprintf("%s%s%s", t.ProjectName, COLON_SEPARATED, t.TopicName)))
	return resourceAliyunDatahubTopicRead(d, meta)
}

func resourceAliyunDatahubTopicRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	datahubService := DatahubService{client}
	object, err := datahubService.DescribeDatahubTopic(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.SetId(strings.ToLower(fmt.Sprintf("%s%s%s", object.ProjectName, COLON_SEPARATED, object.TopicName)))

	d.Set("name", object.TopicName)
	d.Set("project_name", object.ProjectName)
	d.Set("shard_count", object.ShardCount)
	d.Set("life_cycle", object.Lifecycle)
	d.Set("comment", object.Comment)
	d.Set("record_type", object.RecordType.String())
	if object.RecordSchema != nil {
		d.Set("record_schema", recordSchemaToMap(object.RecordSchema.Fields))
	}
	d.Set("create_time", datahub.Uint64ToTimeString(object.CreateTime))
	d.Set("last_modify_time", datahub.Uint64ToTimeString(object.LastModifyTime))
	return nil
}

func resourceAliyunDatahubTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	projectName, topicName := parts[0], parts[1]

	client := meta.(*connectivity.AliyunClient)
	// Currently, life_cycle can not be modified and it will be fixed in the next future.
	if d.HasChange("life_cycle") || d.HasChange("comment") {
		lifeCycle := d.Get("life_cycle").(int)
		topicComment := d.Get("comment").(string)

		var requestInfo *datahub.DataHub

		raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			requestInfo = dataHubClient
			return nil, dataHubClient.UpdateTopic(projectName, topicName, lifeCycle, topicComment)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateTopic", AliyunDatahubSdkGo)
		}
		if debugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			requestMap["TopicName"] = topicName
			requestMap["LifeCycle"] = strconv.Itoa(lifeCycle)
			requestMap["TopicComment"] = topicComment
			addDebug("UpdateTopic", raw, requestInfo, requestMap)
		}
	}

	return resourceAliyunDatahubTopicRead(d, meta)
}

func resourceAliyunDatahubTopicDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	projectName, topicName := parts[0], parts[1]

	client := meta.(*connectivity.AliyunClient)
	datahubService := DatahubService{client}
	var requestInfo *datahub.DataHub

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			requestInfo = dataHubClient
			return nil, dataHubClient.DeleteTopic(projectName, topicName)
		})
		if err != nil {
			if isRetryableDatahubError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			requestMap["TopicName"] = topicName
			addDebug("DeleteTopic", raw, requestInfo, requestMap)
		}
		return nil
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteTopic", AliyunDatahubSdkGo)
	}
	return WrapError(datahubService.WaitForDatahubTopic(d.Id(), Deleted, DefaultTimeout))
}

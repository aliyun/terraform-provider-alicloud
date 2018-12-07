package alicloud

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/schema"
)

// log is append only, so it only supoort create
func resourceAlicloudLogLogs() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogLogsPost,
		Read:   resourceAlicloudLogLogsNil,
		Delete: resourceAlicloudLogLogsNil,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logs": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"contents": &schema.Schema{
							Type:     schema.TypeMap,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"retry_seconds": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"shard_hash_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"result": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudLogLogsNil(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlicloudLogLogsPost(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	logs := &sls.LogGroup{}

	if topic, ok := d.GetOk("topic"); ok {
		logs.Topic = proto.String(topic.(string))
	}

	if source, ok := d.GetOk("source"); ok {
		logs.Source = proto.String(source.(string))
	}

	var trySeconds int64
	if trySecondsObj, ok := d.GetOk("retry_seconds"); ok {
		trySeconds = int64(trySecondsObj.(int))
	}

	if tags, ok := d.GetOk("tags"); ok {
		for key, value := range tags.(map[string]interface{}) {
			tag := &sls.LogTag{
				Key:   proto.String(key),
				Value: proto.String(value.(string)),
			}
			logs.LogTags = append(logs.LogTags, tag)
		}
	}

	logDetails := d.Get("logs")
	for _, logD := range logDetails.([]interface{}) {
		logDetail := logD.(map[string]interface{})
		log := &sls.Log{}
		if t, ok := logDetail["time"]; ok && t.(int) > 0 {
			log.Time = proto.Uint32(uint32(t.(int)))
		} else {
			log.Time = proto.Uint32(uint32(time.Now().Unix()))
		}

		contents := logDetail["contents"].(map[string]interface{})
		for key, value := range contents {
			content := &sls.LogContent{
				Key:   proto.String(key),
				Value: proto.String(value.(string)),
			}
			log.Contents = append(log.Contents, content)
		}
		logs.Logs = append(logs.Logs, log)
	}

	var shardHash *string
	if hashVal, ok := d.GetOk("shard_hash_key"); ok {
		shardHash = proto.String(hashVal.(string))
	}

	var err error
	if trySeconds > 0 {
		_, err = client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			startTime := time.Now()
			for {
				err := slsClient.PostLogStoreLogs(
					d.Get("project").(string),
					d.Get("logstore").(string),
					logs,
					shardHash,
				)
				if err == nil {
					return nil, err
				}
				if time.Now().Sub(startTime).Seconds() > (float64)(trySeconds) {
					return nil, err
				}
				time.Sleep(time.Millisecond * 200)
			}
		})
	} else {
		_, err = client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.PostLogStoreLogs(
				d.Get("project").(string),
				d.Get("logstore").(string),
				logs,
				shardHash,
			)
		})
	}

	if err != nil {
		return fmt.Errorf("post logs got an error: %#v", err)
	}
	d.Set("result", true)
	return nil
}

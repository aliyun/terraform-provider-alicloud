package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudChatbotPublishTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudChatbotPublishTaskCreate,
		Read:   resourceAlicloudChatbotPublishTaskRead,
		Update: resourceAlicloudChatbotPublishTaskUpdate,
		Delete: resourceAlicloudChatbotPublishTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agent_key": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"biz_type": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"data_id_list": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"modify_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudChatbotPublishTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	var response map[string]interface{}
	action := "CreatePublishTask"
	request["BizType"] = d.Get("biz_type")
	if v, ok := d.GetOk("data_id_list"); ok {
		request["DataIdList"] = v
	}
	if v, ok := d.GetOk("agent_key"); ok {
		request["AgentKey"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Chatbot", "2022-04-08", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_chatbot_publish_task", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.Id", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_chatbot_publish_task")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudChatbotPublishTaskRead(d, meta)
}

func resourceAlicloudChatbotPublishTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	chatbotService := ChatbotService{client}

	object, err := chatbotService.DescribeChatbotPublishTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_chatbot_publish_task chatbotService.DescribeChatbotPublishTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["BizTypeList"]; ok {
		d.Set("biz_type", v.([]interface{})[0])
	}
	d.Set("create_time", object["CreateTime"])
	d.Set("modify_time", object["ModifyTime"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAlicloudChatbotPublishTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudChatbotPublishTaskRead(d, meta)
}

func resourceAlicloudChatbotPublishTaskDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudChatbotPublishTask. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

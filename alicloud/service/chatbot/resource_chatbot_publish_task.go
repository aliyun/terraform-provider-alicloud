package chatbot

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePublishTask() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublishTaskCreate,
		Read:   resourcePublishTaskRead,
		Update: resourcePublishTaskUpdate,
		Delete: resourcePublishTaskDelete,
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

func resourcePublishTaskCreate(d *schema.ResourceData, meta interface{}) error {
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		resp, err := client.RpcPost("Chatbot", "2022-04-08", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_chatbot_publish_task", action, errs.AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.Id", response); err != nil || v == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_chatbot_publish_task")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourcePublishTaskRead(d, meta)
}

func resourcePublishTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	chatbotService := NewService(client)

	object, err := chatbotService.DescribeChatbotPublishTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_chatbot_publish_task chatbotService.DescribeChatbotPublishTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	if v, ok := object["BizTypeList"]; ok {
		d.Set("biz_type", v.([]interface{})[0])
	}
	d.Set("create_time", object["CreateTime"])
	d.Set("modify_time", object["ModifyTime"])
	d.Set("status", object["Status"])

	return nil
}

func resourcePublishTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourcePublishTaskRead(d, meta)
}

func resourcePublishTaskDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourcePublishTask. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

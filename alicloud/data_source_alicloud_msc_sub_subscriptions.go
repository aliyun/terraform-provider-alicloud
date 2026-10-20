package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlicloudMscSubSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMscSubSubscriptionsRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"item_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"item_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pmsg_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sms_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tts_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"webhook_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"webhook_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMscSubSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListSubscriptionItems"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("filter"); ok {
		request["Filter"] = v
	}
	request["Locale"] = "en"
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("MscOpenSubscription", "2021-07-13", action, request, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_msc_sub_subscriptions", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.SubscriptionItems", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.SubscriptionItems", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		email_statusVal, err := formatInt(object["EmailStatus"])
		if err != nil {
			return err
		}

		pmsg_statusVal, err := formatInt(object["PmsgStatus"])
		if err != nil {
			return err
		}
		sms_statusVal, err := formatInt(object["SmsStatus"])
		if err != nil {
			return err
		}
		tts_statusVal, err := formatInt(object["TtsStatus"])
		if err != nil {
			return err
		}

		webhook_statusVal, err := formatInt(object["WebhookStatus"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"channel":        object["Channel"],
			"contact_ids":    object["ContactIds"],
			"description":    object["Description"],
			"email_status":   email_statusVal,
			"id":             fmt.Sprint(object["ItemId"]),
			"item_id":        fmt.Sprint(object["ItemId"]),
			"item_name":      object["ItemName"],
			"pmsg_status":    pmsg_statusVal,
			"sms_status":     sms_statusVal,
			"tts_status":     tts_statusVal,
			"webhook_ids":    object["WebhookIds"],
			"webhook_status": webhook_statusVal,
		}
		s = append(s, mapping)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("subscriptions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

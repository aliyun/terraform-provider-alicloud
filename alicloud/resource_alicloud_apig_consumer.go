package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApigConsumer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigConsumerCreate,
		Read:   resourceAliCloudApigConsumerRead,
		Update: resourceAliCloudApigConsumerUpdate,
		Delete: resourceAliCloudApigConsumerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"consumer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"gateway_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "API",
				ValidateFunc: StringInSlice([]string{"API", "AI"}, false),
			},
			"credential_generate_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "System",
				ValidateFunc: StringInSlice([]string{"System"}, false),
			},
		},
	}
}

func resourceAliCloudApigConsumerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/v1/consumers"
	body := apigConsumerCreateRequest(d)
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, nil, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_consumer", action, AlibabaCloudSdkGoERROR)
	}
	data, err := apigResponseData(response)
	if err != nil {
		return err
	}
	consumerID, ok := data["consumerId"].(string)
	if !ok || consumerID == "" {
		return fmt.Errorf("APIG CreateConsumer response does not contain consumerId")
	}
	d.SetId(consumerID)
	if err := apigWaitForConsumerState(client, d, d.Timeout(schema.TimeoutCreate), func(consumer map[string]interface{}) bool {
		return consumer != nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, consumerID, action, AlibabaCloudSdkGoERROR)
	}
	return resourceAliCloudApigConsumerRead(d, meta)
}

func apigConsumerCreateRequest(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("consumer_name"),
		"description": d.Get("description"),
		"enable":      d.Get("enable"),
		"gatewayType": d.Get("gateway_type"),
		"akSkIdentityConfigs": []map[string]interface{}{
			{
				"type":         "AkSk",
				"generateMode": d.Get("credential_generate_mode"),
			},
		},
	}
}

func resourceAliCloudApigConsumerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	consumer, err := apigFindConsumerByID(client, d.Id(), d.Get("consumer_name").(string))
	if err != nil {
		return WrapError(err)
	}
	if consumer == nil {
		if !d.IsNewResource() {
			log.Printf("[DEBUG] Resource alicloud_apig_consumer %s was not found", d.Id())
		}
		d.SetId("")
		return nil
	}
	if err := d.Set("consumer_name", consumer["name"]); err != nil {
		return err
	}
	if err := d.Set("description", consumer["description"]); err != nil {
		return err
	}
	if err := d.Set("enable", consumer["enable"]); err != nil {
		return err
	}
	return nil
}

func apigFindConsumerByID(client *connectivity.AliyunClient, consumerID, nameLike string) (map[string]interface{}, error) {
	const pageSize = 100
	for page := 1; ; page++ {
		query := map[string]*string{
			"pageNumber": StringPointer(strconv.Itoa(page)),
			"pageSize":   StringPointer(strconv.Itoa(pageSize)),
		}
		if nameLike != "" {
			query["nameLike"] = StringPointer(nameLike)
		}
		response, err := client.RoaGet("APIG", "2024-03-27", "/v1/consumers", query, nil, nil)
		if err != nil {
			return nil, err
		}
		data, err := apigResponseData(response)
		if err != nil {
			return nil, err
		}
		items := apigObjectSlice(data["items"])
		for _, item := range items {
			if item["consumerId"] == consumerID {
				return item, nil
			}
		}
		if len(items) < pageSize {
			return nil, nil
		}
	}
}

func apigWaitForConsumerState(client *connectivity.AliyunClient, d *schema.ResourceData, timeout time.Duration, matches func(map[string]interface{}) bool) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		consumer, err := apigFindConsumerByID(client, d.Id(), d.Get("consumer_name").(string))
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if !matches(consumer) {
			return resource.RetryableError(fmt.Errorf("APIG consumer %s has not reached the expected state", d.Id()))
		}
		return nil
	})
}

func resourceAliCloudApigConsumerUpdate(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChanges("description", "enable") {
		return resourceAliCloudApigConsumerRead(d, meta)
	}
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/v1/consumers/%s", d.Id())
	body := map[string]interface{}{
		"description": d.Get("description"),
		"enable":      d.Get("enable"),
	}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err = client.RoaPut("APIG", "2024-03-27", action, nil, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if err := apigWaitForConsumerState(client, d, d.Timeout(schema.TimeoutUpdate), func(consumer map[string]interface{}) bool {
		if consumer == nil || consumer["enable"] != d.Get("enable") {
			return false
		}
		return fmt.Sprint(consumer["description"]) == d.Get("description").(string)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return resourceAliCloudApigConsumerRead(d, meta)
}

func resourceAliCloudApigConsumerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/v1/consumers/%s", d.Id())
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = client.RoaDelete("APIG", "2024-03-27", action, nil, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFound.ConsumerNotFound"}) || NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if err := apigWaitForConsumerState(client, d, d.Timeout(schema.TimeoutDelete), func(consumer map[string]interface{}) bool {
		return consumer == nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

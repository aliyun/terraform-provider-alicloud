package amqpopen

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAmqpStaticAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAmqpStaticAccountCreate,
		Read:   resourceAliCloudAmqpStaticAccountRead,
		Delete: resourceAliCloudAmqpStaticAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAmqpStaticAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateOrGetAccount"
	request := make(map[string]interface{})
	var err error

	timestamp := time.Now().UnixMilli()
	request["createTimestamp"] = timestamp

	request["instanceId"] = d.Get("instance_id")
	request["accountAccessKey"] = d.Get("access_key")

	if v, ok := d.GetOk("secret_key"); ok {
		mac := hmac.New(sha1.New, []byte(v.(string)))
		mac.Write([]byte(strconv.FormatInt(timestamp, 10)))
		signature := mac.Sum(nil)
		request["signature"] = hex.EncodeToString(signature)

		macSecret := hmac.New(sha1.New, []byte(strconv.FormatInt(timestamp, 10)))
		macSecret.Write([]byte(v.(string)))
		secretSign := macSecret.Sum(nil)
		request["secretSign"] = hex.EncodeToString(secretSign)
	}

	stringToBase64Encode := "2:" + request["instanceId"].(string) + ":" + request["accountAccessKey"].(string)
	request["userName"] = base64.StdEncoding.EncodeToString([]byte(stringToBase64Encode))

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_amqp_static_account", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["instanceId"], ":", request["accountAccessKey"]))

	return resourceAliCloudAmqpStaticAccountRead(d, meta)
}

func resourceAliCloudAmqpStaticAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpOpenService := NewAmqpOpenService(client)

	object, err := amqpOpenService.DescribeAmqpStaticAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_amqp_static_account amqpOpenService.DescribeAmqpStaticAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("access_key", parts[1])
	d.Set("user_name", object["userName"])
	d.Set("password", object["password"])
	d.Set("master_uid", object["masterUid"])
	d.Set("create_time", object["createTimestamp"])

	return nil
}

func resourceAliCloudAmqpStaticAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAccount"
	var response map[string]interface{}
	var err error

	request := make(map[string]interface{})

	if v, ok := d.GetOk("user_name"); ok {
		request["UserName"] = v
	}
	if v, ok := d.GetOk("create_time"); ok {
		request["CreateTimestamp"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

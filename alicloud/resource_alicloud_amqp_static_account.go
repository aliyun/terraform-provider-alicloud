package alicloud

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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAmqpStaticAccount() *schema.Resource {
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
				Type:     schema.TypeInt,
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_amqp_static_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["instanceId"], ":", request["accountAccessKey"]))

	return resourceAliCloudAmqpStaticAccountRead(d, meta)
}

func resourceAliCloudAmqpStaticAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpOpenService := AmqpOpenService{client}

	object, err := amqpOpenService.DescribeAmqpStaticAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_amqp_static_account amqpOpenService.DescribeAmqpStaticAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

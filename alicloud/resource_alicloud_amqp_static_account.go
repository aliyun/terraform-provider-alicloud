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

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudAmqpStaticAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAmqpStaticAccountCreate,
		Read:   resourceAlicloudAmqpStaticAccountRead,
		Delete: resourceAlicloudAmqpStaticAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_key": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"instance_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"master_uid": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"password": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"secret_key": {
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
				Type:      schema.TypeString,
			},
			"user_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudAmqpStaticAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewOnsproxyClient()
	if err != nil {
		return WrapError(err)
	}
	timestamp := time.Now().UnixMilli()
	request["createTimestamp"] = timestamp

	if v, ok := d.GetOk("instance_id"); ok {
		request["instanceId"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["instanceId"] = v
	}

	if v, ok := d.GetOk("access_key"); ok {
		request["accountAccessKey"] = v
	}
	stringToBase64Encode := "2:" + request["instanceId"].(string) + ":" + request["accountAccessKey"].(string)
	request["userName"] = base64.StdEncoding.EncodeToString([]byte(stringToBase64Encode))

	if v, ok := d.GetOk("secret_key"); ok {
		mac := hmac.New(sha1.New, []byte(v.(string)))
		mac.Write([]byte(strconv.FormatInt(timestamp, 10)))
		signature := mac.Sum(nil)
		request["signature"] = hex.EncodeToString(signature)

		macSecret := hmac.New(sha1.New, []byte(strconv.FormatInt(timestamp, 10)))
		macSecret.Write([]byte(request["accountAccessKey"].(string)))
		secretSign := macSecret.Sum(nil)
		request["secretSign"] = hex.EncodeToString(secretSign)
	}

	var response map[string]interface{}
	action := "CreateOrGetAccount"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_amqp_static_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["instanceId"], ":", request["accountAccessKey"]))

	return resourceAlicloudAmqpStaticAccountRead(d, meta)
}

func resourceAlicloudAmqpStaticAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpOpenService := AmqpOpenService{client}

	object, err := amqpOpenService.DescribeAmqpStaticAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
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
	d.Set("master_uid", object["masterUid"])
	d.Set("user_name", object["userName"])
	d.Set("password", object["password"])
	d.Set("create_time", object["createTimestamp"])

	return nil
}

func resourceAlicloudAmqpStaticAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewOnsproxyClient()
	if err != nil {
		return WrapError(err)
	}
	request := make(map[string]interface{})

	if v, ok := d.GetOk("user_name"); ok {
		request["UserName"] = v
	}
	if v, ok := d.GetOk("create_time"); ok {
		request["CreateTimestamp"] = v
	}

	action := "DeleteAccount"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlikafkaSaslUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlikafkaSaslUserCreate,
		Read:   resourceAlicloudAlikafkaSaslUserRead,
		Update: resourceAlicloudAlikafkaSaslUserUpdate,
		Delete: resourceAlicloudAlikafkaSaslUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"plain", "scram"}, false),
			},
		},
	}
}

func resourceAlicloudAlikafkaSaslUserCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSaslUser"
	request := make(map[string]interface{})
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}
	request["InstanceId"] = d.Get("instance_id")
	request["Username"] = d.Get("username")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}

	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}

	if password != "" {
		request["Password"] = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Password"] = decryptResp
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_sasl_user", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	// Server may have cache, sleep a while.
	time.Sleep(2 * time.Second)
	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Username"]))
	return resourceAlicloudAlikafkaSaslUserRead(d, meta)
}

func resourceAlicloudAlikafkaSaslUserRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAliKafkaSaslUser(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ali_kafka_consumer_group alikafkaService.DescribeAlikafkaSaslUser Failed!!! %s", err)
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
	d.Set("username", parts[1])
	d.Set("type", object["Type"])
	return nil
}

func resourceAlicloudAlikafkaSaslUserUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": parts[0],
		"Username":   parts[1],
		"RegionId":   client.RegionId,
	}

	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}

	if d.HasChange("password") || d.HasChange("kms_encrypted_password") {

		password := d.Get("password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password == "" && kmsPassword == "" {
			return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			request["Password"] = password
		} else {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request["Password"] = decryptResp
		}

		action := "CreateSaslUser"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_sasl_user", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		// Server may have cache, sleep a while.
		time.Sleep(1000)
	}

	return resourceAlicloudAlikafkaSaslUserRead(d, meta)
}

func resourceAlicloudAlikafkaSaslUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteSaslUser"
	var response map[string]interface{}
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"Username":   parts[1],
		"InstanceId": parts[0],
	}

	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}

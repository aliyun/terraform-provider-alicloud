package alicloud

import (
	"encoding/json"
	"fmt"

	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudKmsSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKmsSecretCreate,
		Read:   resourceAlicloudKmsSecretRead,
		Update: resourceAlicloudKmsSecretUpdate,
		Delete: resourceAlicloudKmsSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"force_delete_without_recovery": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"planned_delete_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"recovery_window_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("force_delete_without_recovery").(bool)
				},
			},
			"secret_data": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"secret_data_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"binary", "text"}, false),
				Default:      "text",
			},
			"secret_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_stages": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudKmsSecretCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := kms.CreateCreateSecretRequest()
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("encryption_key_id"); ok {
		request.EncryptionKeyId = v.(string)
	}

	request.SecretData = d.Get("secret_data").(string)
	if v, ok := d.GetOk("secret_data_type"); ok {
		request.SecretDataType = v.(string)
	}

	request.SecretName = d.Get("secret_name").(string)
	if v, ok := d.GetOk("tags"); ok {
		addTags := make([]JsonTag, 0)
		for key, value := range v.(map[string]interface{}) {
			addTags = append(addTags, JsonTag{
				TagKey:   key,
				TagValue: value.(string),
			})
		}
		tags, err := json.Marshal(addTags)
		if err != nil {
			return WrapError(err)
		}
		request.Tags = string(tags)
	}
	request.VersionId = d.Get("version_id").(string)
	wait := incrementalWait(3*time.Second, 1*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		args := *request
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.CreateSecret(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Rejected.Throttling"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*kms.CreateSecretResponse)
		d.SetId(fmt.Sprintf("%v", response.SecretName))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_secret", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudKmsSecretRead(d, meta)
}
func resourceAlicloudKmsSecretRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	object, err := kmsService.DescribeKmsSecret(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_secret kmsService.DescribeKmsSecret Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("secret_name", d.Id())
	d.Set("arn", object.Arn)
	d.Set("description", object.Description)
	d.Set("encryption_key_id", object.EncryptionKeyId)
	d.Set("planned_delete_time", object.PlannedDeleteTime)

	tags := make(map[string]string)
	for _, t := range object.Tags.Tag {
		if !ignoredTags(t.TagKey, t.TagValue) {
			tags[t.TagKey] = t.TagValue
		}
	}
	d.Set("tags", tags)

	getSecretValueObject, err := kmsService.GetSecretValue(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("secret_data", getSecretValueObject.SecretData)
	d.Set("secret_data_type", getSecretValueObject.SecretDataType)
	d.Set("version_id", getSecretValueObject.VersionId)
	d.Set("version_stages", getSecretValueObject.VersionStages.VersionStage)
	return nil
}
func resourceAlicloudKmsSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := kmsService.SetResourceTags(d, "secret"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if d.HasChange("description") {
		request := kms.CreateUpdateSecretRequest()
		request.SecretName = d.Id()
		request.Description = d.Get("description").(string)
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.UpdateSecret(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
	}
	update := false
	request := kms.CreatePutSecretValueRequest()
	request.SecretName = d.Id()
	if d.HasChange("secret_data") {
		update = true
	}
	request.SecretData = d.Get("secret_data").(string)
	if d.HasChange("version_id") {
		update = true
	}
	request.VersionId = d.Get("version_id").(string)
	if d.HasChange("secret_data_type") {
		update = true
		request.SecretDataType = d.Get("secret_data_type").(string)
	}
	if d.HasChange("version_stages") {
		update = true
		request.VersionStages = convertListToJsonString(d.Get("version_stages").(*schema.Set).List())
	}
	if update {
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.PutSecretValue(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("secret_data")
		d.SetPartial("version_id")
		d.SetPartial("secret_data_type")
		d.SetPartial("version_stages")
	}
	d.Partial(false)
	return resourceAlicloudKmsSecretRead(d, meta)
}
func resourceAlicloudKmsSecretDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := kms.CreateDeleteSecretRequest()
	request.SecretName = d.Id()
	if v, ok := d.GetOkExists("force_delete_without_recovery"); ok {
		request.ForceDeleteWithoutRecovery = fmt.Sprintf("%v", v.(bool))
	}
	if v, ok := d.GetOk("recovery_window_in_days"); ok {
		request.RecoveryWindowInDays = fmt.Sprintf("%v", v.(int))
	}
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DeleteSecret(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.ResourceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

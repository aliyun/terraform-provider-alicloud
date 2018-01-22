package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/denverdino/aliyungo/kms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudKmsKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKmsKeyCreate,
		Read:   resourceAlicloudKmsKeyRead,
		Update: resourceAlicloudKmsKeyUpdate,
		Delete: resourceAlicloudKmsKeyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "From Terraform",
				ValidateFunc: validateStringLengthInRange(0, 8192),
			},
			"key_usage": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !(kms.KeyUsage(value) == kms.KEY_USAGE_ENCRYPT_DECRYPT) {
						es = append(es, fmt.Errorf(
							"%q must be %s", k, kms.KEY_USAGE_ENCRYPT_DECRYPT))
					}
					return
				},
				Default: kms.KEY_USAGE_ENCRYPT_DECRYPT,
			},
			"is_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"deletion_window_in_days": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(7, 30),
				Default:      30,
			},
		},
	}
}

func resourceAlicloudKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).kmsconn

	args := kms.CreateKeyArgs{
		KeyUsage: kms.KeyUsage(d.Get("key_usage").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		args.Description = v.(string)
	}
	resp, err := conn.CreateKey(&args)
	if err != nil {
		return fmt.Errorf("CreateKey got an error: %#v.", err)
	}

	d.SetId(resp.KeyMetadata.KeyId)

	return resourceAlicloudKmsKeyUpdate(d, meta)
}

func resourceAlicloudKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).kmsconn

	key, err := conn.DescribeKey(d.Id())
	if err != nil {
		if IsExceptedError(err, ForbiddenKeyNotFound) {
			return nil
		}
		return fmt.Errorf("DescribeKey got an error: %#v.", err)
	}

	if KeyState(key.KeyMetadata.KeyState) == PendingDeletion {
		log.Printf("[WARN] Removing KMS key %s because it's already gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("description", key.KeyMetadata.Description)
	d.Set("key_usage", key.KeyMetadata.KeyUsage)
	d.Set("is_enabled", KeyState(key.KeyMetadata.KeyState) == Enabled)
	d.Set("deletion_window_in_days", d.Get("deletion_window_in_days").(int))

	return nil
}

func resourceAlicloudKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).kmsconn

	d.Partial(true)

	if d.HasChange("is_enabled") {
		key, err := conn.DescribeKey(d.Id())
		if err != nil {
			return fmt.Errorf("DescribeKey got an error: %#v.", err)
		}

		if d.Get("is_enabled").(bool) && KeyState(key.KeyMetadata.KeyState) == Disabled {
			if _, err := conn.EnableKey(d.Id()); err != nil {
				return fmt.Errorf("Enable key got an error: %#v.", err)
			}
		}

		if !d.Get("is_enabled").(bool) && KeyState(key.KeyMetadata.KeyState) == Enabled {
			if _, err := conn.DisableKey(d.Id()); err != nil {
				return fmt.Errorf("Disable key got an error: %#v.", err)
			}
		}
		d.SetPartial("is_enabled")
	}

	d.Partial(false)

	return resourceAlicloudKmsKeyRead(d, meta)
}

func resourceAlicloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).kmsconn

	if _, err := conn.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionArgs{
		KeyId:               d.Id(),
		PendingWindowInDays: d.Get("deletion_window_in_days").(int),
	}); err != nil {
		return err
	}

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		key, err := conn.DescribeKey(d.Id())
		if err != nil {
			if IsExceptedError(err, ForbiddenKeyNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("DescribeKey got an error: %#v.", err))
		}

		if key == nil || KeyState(key.KeyMetadata.KeyState) == PendingDeletion {
			log.Printf("[WARN] Removing KMS key %s because it's already gone", d.Id())
			d.SetId("")
			return nil
		}
		return resource.RetryableError(fmt.Errorf("ScheduleKeyDeletion timeout."))
	})
}

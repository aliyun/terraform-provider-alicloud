package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamAccessKeyCreate,
		Read:   resourceAlicloudRamAccessKeyRead,
		Update: resourceAlicloudRamAccessKeyUpdate,
		Delete: resourceAlicloudRamAccessKeyDelete,

		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"secret_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      Active,
				ValidateFunc: validateRamAKStatus,
			},
		},
	}
}

func resourceAlicloudRamAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.UserQueryRequest{}
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		args.UserName = v.(string)
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.CreateAccessKey(args)
	})
	if err != nil {
		return fmt.Errorf("CreateAccessKey got an error: %#v", err)
	}
	response, _ := raw.(ram.AccessKeyResponse)

	// create a secret_file and write access key to it.
	if output, ok := d.GetOk("secret_file"); ok && output != nil {
		writeToFile(output.(string), response.AccessKey)
	}

	d.SetId(response.AccessKey.AccessKeyId)
	return resourceAlicloudRamAccessKeyUpdate(d, meta)
}

func resourceAlicloudRamAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	args := ram.UpdateAccessKeyRequest{
		UserAccessKeyId: d.Id(),
		Status:          ram.State(d.Get("status").(string)),
	}
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		args.UserName = v.(string)
	}

	if d.HasChange("status") {
		d.SetPartial("status")
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.UpdateAccessKey(args)
		})
		if err != nil {
			return fmt.Errorf("UpdateAccessKey got an error: %#v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamAccessKeyRead(d, meta)
}

func resourceAlicloudRamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.UserQueryRequest{}
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		args.UserName = v.(string)
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.ListAccessKeys(args)
	})
	if err != nil {
		return fmt.Errorf("Get list access keys got an error: %#v", err)
	}
	response, _ := raw.(ram.AccessKeyListResponse)
	accessKeys := response.AccessKeys.AccessKey
	if len(accessKeys) < 1 {
		return fmt.Errorf("No access keys found.")
	}

	for _, v := range accessKeys {
		if v.AccessKeyId == d.Id() {
			d.Set("status", v.Status)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceAlicloudRamAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.UpdateAccessKeyRequest{
		UserAccessKeyId: d.Id(),
	}

	queryArgs := ram.UserQueryRequest{}

	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		args.UserName = v.(string)
		queryArgs.UserName = v.(string)
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.DeleteAccessKey(args)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting access key: %#v", err))
		}

		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.ListAccessKeys(queryArgs)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		response, _ := raw.(ram.AccessKeyListResponse)

		if len(response.AccessKeys.AccessKey) < 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Error deleting access key - trying again while it is deleted."))
	})
}

package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudHbrVault() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudHbrVaultCreate,
		Read:   resourceAliCloudHbrVaultRead,
		Update: resourceAliCloudHbrVaultUpdate,
		Delete: resourceAliCloudHbrVaultDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vault_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vault_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"STANDARD", "OTS_BACKUP"}, false),
			},
			"vault_storage_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"STANDARD"}, false),
			},
			"encrypt_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HBR_PRIVATE", "KMS"}, false),
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("encrypt_type"); ok && v.(string) == "HBR_PRIVATE" {
						return true
					}
					return false
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"redundancy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"LRS", "ZRS"}, false),
				Removed:      "Field `redundancy_type` has been removed from provider version 1.209.1.",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudHbrVaultCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	var response map[string]interface{}
	action := "CreateVault"
	request := make(map[string]interface{})
	var err error

	request["VaultRegionId"] = client.RegionId
	request["VaultName"] = d.Get("vault_name")

	if v, ok := d.GetOk("vault_type"); ok {
		request["VaultType"] = v
	}

	if v, ok := d.GetOk("vault_storage_class"); ok {
		request["VaultStorageClass"] = v
	}

	if v, ok := d.GetOk("encrypt_type"); ok {
		request["EncryptType"] = v
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KmsKeyId"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_vault", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VaultId"]))

	stateConf := BuildStateConf([]string{}, []string{"CREATED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbrService.HbrVaultStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudHbrVaultRead(d, meta)
}

func resourceAliCloudHbrVaultRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}

	object, err := hbrService.DescribeHbrVault(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_vault hbrService.DescribeHbrVault Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("vault_name", object["VaultName"])
	d.Set("vault_type", object["VaultType"])
	d.Set("vault_storage_class", object["VaultStorageClass"])
	d.Set("encrypt_type", object["EncryptType"])
	d.Set("kms_key_id", object["KmsKeyId"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAliCloudHbrVaultUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"VaultId": d.Id(),
	}

	if d.HasChange("vault_name") {
		update = true
		request["VaultName"] = d.Get("vault_name")
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if update {
		action := "UpdateVault"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudHbrVaultRead(d, meta)
}

func resourceAliCloudHbrVaultDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVault"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"VaultId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"CannotDeleteReplicationSourceVault"}) || NeedRetry(err) {
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

	return nil
}

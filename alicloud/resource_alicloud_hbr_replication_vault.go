// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudHbrReplicationVault() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudHbrReplicationVaultCreate,
		Read:   resourceAliCloudHbrReplicationVaultRead,
		Update: resourceAliCloudHbrReplicationVaultUpdate,
		Delete: resourceAliCloudHbrReplicationVaultDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypt_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"HBR_PRIVATE", "KMS"}, false),
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_source_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"replication_source_vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vault_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vault_storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudHbrReplicationVaultCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateReplicationVault"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["VaultRegionId"] = client.RegionId

	if v, ok := d.GetOk("vault_storage_class"); ok {
		request["VaultStorageClass"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["VaultName"] = d.Get("vault_name")
	if v, ok := d.GetOk("encrypt_type"); ok {
		request["EncryptType"] = v
	}
	request["ReplicationSourceVaultId"] = d.Get("replication_source_vault_id")
	request["ReplicationSourceRegionId"] = d.Get("replication_source_region_id")
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KmsKeyId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_replication_vault", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VaultId"]))

	return resourceAliCloudHbrReplicationVaultRead(d, meta)
}

func resourceAliCloudHbrReplicationVaultRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrServiceV2 := HbrServiceV2{client}

	objectRaw, err := hbrServiceV2.DescribeHbrReplicationVault(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_replication_vault DescribeHbrReplicationVault Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("encrypt_type", objectRaw["EncryptType"])
	d.Set("kms_key_id", objectRaw["KmsKeyId"])
	d.Set("region_id", objectRaw["VaultRegionId"])
	d.Set("replication_source_region_id", objectRaw["ReplicationSourceRegionId"])
	d.Set("replication_source_vault_id", objectRaw["ReplicationSourceVaultId"])
	d.Set("vault_name", objectRaw["VaultName"])
	d.Set("vault_storage_class", objectRaw["VaultStorageClass"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudHbrReplicationVaultUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateVault"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VaultId"] = d.Id()

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("vault_name") {
		update = true
	}
	request["VaultName"] = d.Get("vault_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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

	return resourceAliCloudHbrReplicationVaultRead(d, meta)
}

func resourceAliCloudHbrReplicationVaultDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVault"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["VaultId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

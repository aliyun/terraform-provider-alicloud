package alicloud

import (
	"fmt"
	"time"

	roacs "github.com/alibabacloud-go/cs-20151215/v7/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const ResourceAlicloudCSKMSEncryption = "resourceAlicloudCSKMSEncryption"

func resourceAlicloudCSKMSEncryption() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKMSEncryptionCreate,
		Read:   resourceAlicloudCSKMSEncryptionRead,
		Update: resourceAlicloudCSKMSEncryptionUpdate,
		Delete: resourceAlicloudCSKMSEncryptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disable_encryption": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCSKMSEncryptionCreate(d *schema.ResourceData, meta interface{}) error {
	// Create and Update use the same API
	return resourceAlicloudCSKMSEncryptionUpdate(d, meta)
}

func resourceAlicloudCSKMSEncryptionRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(*connectivity.AliyunClient).NewRoaCsV7Client()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKMSEncryption, "InitializeClient", err)
	}

	clusterId := d.Id()

	// Get cluster detail to check encryption status
	response, err := csClient.DescribeClusterDetail(tea.String(clusterId))
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKMSEncryption, "DescribeClusterDetail", err)
	}

	d.Set("cluster_id", clusterId)

	// Extract encryption info from MetaData using fetchClusterCapabilities
	capabilities := fetchClusterCapabilities(tea.StringValue(response.Body.MetaData))
	if v, ok := capabilities["EncryptionKMSKeyId"]; ok {
		d.Set("kms_key_id", Interface2String(v))
	} else {
		d.Set("kms_key_id", "")
	}
	if v, ok := capabilities["DisableEncryption"]; ok {
		d.Set("disable_encryption", Interface2Bool(v))
	} else {
		d.Set("disable_encryption", true)
	}

	return nil
}

func resourceAlicloudCSKMSEncryptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csClient, err := client.NewRoaCsV7Client()
	if err != nil {
		return err
	}

	clusterId := d.Get("cluster_id").(string)

	request := &roacs.UpdateKMSEncryptionRequest{
		DisableEncryption: tea.Bool(d.Get("disable_encryption").(bool)),
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request.KmsKeyId = tea.String(v.(string))
	}

	csService := CsService{client}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = csClient.UpdateKMSEncryption(tea.String(clusterId), request)
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
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKMSEncryption, "UpdateKMSEncryption", err)
	}

	// Wait for cluster to return to running state
	stateConf := BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(clusterId, []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	d.SetId(clusterId)
	return resourceAlicloudCSKMSEncryptionRead(d, meta)
}

func resourceAlicloudCSKMSEncryptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csClient, err := client.NewRoaCsV7Client()
	if err != nil {
		return err
	}

	clusterId := d.Id()

	// Disable KMS encryption
	request := &roacs.UpdateKMSEncryptionRequest{
		DisableEncryption: tea.Bool(true),
	}

	csService := CsService{client}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = csClient.UpdateKMSEncryption(tea.String(clusterId), request)
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
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKMSEncryption, "UpdateKMSEncryption", err)
	}

	// Wait for cluster to return to running state
	stateConf := BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutDelete), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(clusterId, []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	// Verify that KMS encryption is disabled by checking cluster capabilities
	response, err := csClient.DescribeClusterDetail(tea.String(clusterId))
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKMSEncryption, "DescribeClusterDetail", err)
	}

	capabilities := fetchClusterCapabilities(tea.StringValue(response.Body.MetaData))
	if v, ok := capabilities["DisableEncryption"]; ok {
		if !Interface2Bool(v) {
			return WrapError(fmt.Errorf("KMS encryption is still enabled after delete operation"))
		}
	}

	return nil
}

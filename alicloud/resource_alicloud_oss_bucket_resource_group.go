package alicloud

import (
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOssBucketResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOssBucketResourceGroupCreate,
		Read:   resourceAlicloudOssBucketResourceGroupRead,
		Update: resourceAlicloudOssBucketResourceGroupUpdate,
		Delete: resourceAlicloudOssBucketResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudOssBucketResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bucket := d.Get("bucket").(string)
	resourceGroupId := d.Get("resource_group_id").(string)

	resourceGroup := oss.PutBucketResourceGroup{
		ResourceGroupId: resourceGroupId,
	}

	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.PutBucketResourceGroup(bucket, resourceGroup)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PutBucketResourceGroup", AliyunOssGoSdk)
	}
	addDebug("PutBucketResourceGroup", raw, requestInfo, map[string]interface{}{
		"bucketName":      bucket,
		"resourceGroupId": resourceGroupId,
	})

	d.SetId(bucket)
	return resourceAlicloudOssBucketResourceGroupRead(d, meta)
}

func resourceAlicloudOssBucketResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bucket := d.Id()

	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.GetBucketResourceGroup(bucket)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketResourceGroup", AliyunOssGoSdk)
	}
	addDebug("GetBucketResourceGroup", raw, requestInfo, map[string]interface{}{
		"bucketName": bucket,
	})

	resourceGroup, _ := raw.(oss.GetBucketResourceGroupResult)

	// Update to resouce data
	d.Set("bucket", bucket)
	d.Set("resource_group_id", resourceGroup.ResourceGroupId)

	return nil
}

func resourceAlicloudOssBucketResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bucket := d.Id()
	resourceGroupId := d.Get("resource_group_id").(string)

	resourceGroup := oss.PutBucketResourceGroup{
		ResourceGroupId: resourceGroupId,
	}

	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.PutBucketResourceGroup(bucket, resourceGroup)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PutBucketResourceGroup", AliyunOssGoSdk)
	}
	addDebug("PutBucketResourceGroup", raw, requestInfo, map[string]interface{}{
		"bucketName":      bucket,
		"resourceGroupId": resourceGroupId,
	})

	// Update to resouce data
	d.Set("bucket", bucket)
	d.Set("resource_group_id", resourceGroupId)
	return nil
}

func resourceAlicloudOssBucketResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudOssBucketResourceGroup. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDataWorksCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataWorksCertificateCreate,
		Read:   resourceAliCloudDataWorksCertificateRead,
		Delete: resourceAliCloudDataWorksCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"certificate_file": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"file_size_in_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudDataWorksCertificateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ImportCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["CertificateFile"] = d.Get("certificate_file")
	request["ProjectId"] = d.Get("project_id")
	request["Name"] = d.Get("name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_certificate", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", d.Get("project_id"), response["Id"]))

	return resourceAliCloudDataWorksCertificateRead(d, meta)
}

func resourceAliCloudDataWorksCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksCertificate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_certificate DescribeDataWorksCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", dataWorksCertificateCreateTimeFormat(objectRaw["CreateTime"]))
	d.Set("create_user", objectRaw["CreateUser"])
	d.Set("description", objectRaw["Description"])
	d.Set("file_size_in_bytes", objectRaw["FileSizeInBytes"])
	d.Set("name", objectRaw["Name"])
	d.Set("project_id", objectRaw["ProjectId"])

	return nil
}

func resourceAliCloudDataWorksCertificateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCertificate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length 2, got %d", d.Id(), len(parts)))
	}
	request = make(map[string]interface{})
	request["Id"] = parts[1]
	request["RegionId"] = client.RegionId

	request["ProjectId"] = parts[0]
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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

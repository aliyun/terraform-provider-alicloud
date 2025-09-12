package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudOssBucketHttpsConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketHttpsConfigCreate,
		Read:   resourceAliCloudOssBucketHttpsConfigRead,
		Update: resourceAliCloudOssBucketHttpsConfigUpdate,
		Delete: resourceAliCloudOssBucketHttpsConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"tls_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudOssBucketHttpsConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?httpsConfig")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})

	tLS := make(map[string]interface{})
	nodeNative1, _ := jsonpath.Get("$", d.Get("tls_versions"))
	tLS["TLSVersion"] = make([]interface{}, 0)
	if nodeNative1 != nil && nodeNative1 != "" {
		tLS["TLSVersion"] = nodeNative1.(*schema.Set).List()
	}
	tLS["Enable"] = d.Get("enable")

	objectDataLocalMap["TLS"] = tLS
	request["HttpsConfiguration"] = objectDataLocalMap

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketHttpsConfig", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_https_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	ossServiceV2 := OssServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketHttpsConfigStateRefreshFunc(d.Id(), "#Enable", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudOssBucketHttpsConfigRead(d, meta)
}

func resourceAliCloudOssBucketHttpsConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketHttpsConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_https_config DescribeOssBucketHttpsConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("enable", objectRaw["Enable"])

	tLSVersion1Raw := make([]interface{}, 0)
	if objectRaw["TLSVersion"] != nil {
		tLSVersion1Raw = objectRaw["TLSVersion"].([]interface{})
	}
	if len(tLSVersion1Raw) > 0 {
		d.Set("tls_versions", tLSVersion1Raw)

	}
	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketHttpsConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/?httpsConfig")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())

	objectDataLocalMap := make(map[string]interface{})

	if d.HasChanges("enable", "tls_versions") {
		update = true
	}
	tLS := make(map[string]interface{})
	tLS["TLSVersion"] = make([]interface{}, 0)
	nodeNative1, _ := jsonpath.Get("$", d.Get("tls_versions"))
	if nodeNative1 != nil && nodeNative1 != "" {
		tLS["TLSVersion"] = nodeNative1.(*schema.Set).List()
	}
	tLS["Enable"] = d.Get("enable")

	objectDataLocalMap["TLS"] = tLS
	request["HttpsConfiguration"] = objectDataLocalMap
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketHttpsConfig", action), query, body, nil, hostMap, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		ossServiceV2 := OssServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketHttpsConfigStateRefreshFunc(d.Id(), "#Enable", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudOssBucketHttpsConfigRead(d, meta)
}

func resourceAliCloudOssBucketHttpsConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Bucket Https Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

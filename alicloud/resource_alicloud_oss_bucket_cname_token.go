package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudOssBucketCnameToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketCnameTokenCreate,
		Read:   resourceAliCloudOssBucketCnameTokenRead,
		Delete: resourceAliCloudOssBucketCnameTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudOssBucketCnameTokenCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?cname&comp=token")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	conn, err := client.NewOssClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "BucketCnameConfiguration.Cname.Domain", d.Get("domain"))
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.Execute(genXmlParam("CreateCnameToken", "POST", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_cname_token", action, AlibabaCloudSdkGoERROR)
	}

	CnameTokenBucketVar, _ := jsonpath.Get("$.body.CnameToken.Bucket", response)
	CnameTokenCnameVar, _ := jsonpath.Get("$.body.CnameToken.Cname", response)
	d.SetId(fmt.Sprintf("%v:%v", CnameTokenBucketVar, CnameTokenCnameVar))

	return resourceAliCloudOssBucketCnameTokenRead(d, meta)
}

func resourceAliCloudOssBucketCnameTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketCnameToken(d.Id())
	if err != nil {
		if IsExpectedErrors(err, []string{"NoNeedCreateCnameToken"}) {
			return nil
		}
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_cname_token DescribeOssBucketCnameToken Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Token"] != nil {
		d.Set("token", objectRaw["Token"])
	}
	if objectRaw["Bucket"] != nil {
		d.Set("bucket", objectRaw["Bucket"])
	}
	if objectRaw["Cname"] != nil {
		d.Set("domain", objectRaw["Cname"])
	}

	return nil
}

func resourceAliCloudOssBucketCnameTokenDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Bucket Cname Token. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

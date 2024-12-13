package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcpKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcpKeyPairCreate,
		Read:   resourceAlicloudEcpKeyPairRead,
		Update: resourceAlicloudEcpKeyPairUpdate,
		Delete: resourceAlicloudEcpKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_pair_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_key_body": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudEcpKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/api/v1/providers/Aliyun/products/ECP/resourceTypes/KeyPair/resources"
	body := make(map[string]interface{})
	conn, err := client.NewIaCServiceClient()
	if err != nil {
		return WrapError(err)
	}
	body["KeyPairName"] = d.Get("key_pair_name")
	body["PublicKeyBody"] = d.Get("public_key_body")
	body["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2021-07-22"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), nil, nil, body, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, body)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecp_key_pair", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(body["KeyPairName"]))

	return resourceAlicloudEcpKeyPairRead(d, meta)
}
func resourceAlicloudEcpKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	iaCServiceService := IaCServiceService{client}
	_, err := iaCServiceService.GetKeyPairResource(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecp_key_pair iaCServiceService.GetKeyPairResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("key_pair_name", d.Id())
	return nil
}
func resourceAlicloudEcpKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudEcpKeyPairRead(d, meta)
}
func resourceAlicloudEcpKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/api/v1/providers/Aliyun/products/ECP/resourceTypes/KeyPair/resources/" + d.Id()
	var response map[string]interface{}
	conn, err := client.NewIaCServiceClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]*string{
		"regionId": StringPointer(client.RegionId),
	}

	body := map[string]interface{}{
		"KeyPairName": []string{d.Id()},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2021-07-22"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, body, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"KeyPair.WithInstance"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, body)
	if err != nil {
		if IsExpectedErrors(err, []string{"KeyPairsNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudOssAccessPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssAccessPointCreate,
		Read:   resourceAliCloudOssAccessPointRead,
		Update: resourceAliCloudOssAccessPointUpdate,
		Delete: resourceAliCloudOssAccessPointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_point_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_origin": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_access_block_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_public_access": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudOssAccessPointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?accessPoint")
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "CreateAccessPointConfiguration.AccessPointName", d.Get("access_point_name"))
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return WrapError(err)
	}

	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOk("network_origin"); ok {
		objectDataLocalMap["NetworkOrigin"] = v
	}

	if v := d.Get("vpc_configuration"); !IsNil(v) {
		vpcConfiguration := make(map[string]interface{})
		vpcId1, _ := jsonpath.Get("$[0].vpc_id", v)
		if vpcId1 != nil && vpcId1 != "" {
			vpcConfiguration["VpcId"] = vpcId1
		}
		objectDataLocalMap["VpcConfiguration"] = vpcConfiguration
	}

	if v, ok := d.GetOk("access_point_name"); ok {
		objectDataLocalMap["AccessPointName"] = v
	}

	request["CreateAccessPointConfiguration"] = objectDataLocalMap
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", genXmlParam("PUT", "2019-05-17", "CreateAccessPoint", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_access_point", action, AlibabaCloudSdkGoERROR)
	}

	CreateAccessPointConfigurationAccessPointNameVar, _ := jsonpath.Get("CreateAccessPointConfiguration.AccessPointName", request)
	d.SetId(fmt.Sprintf("%v:%v", *hostMap["bucket"], CreateAccessPointConfigurationAccessPointNameVar))

	return resourceAliCloudOssAccessPointUpdate(d, meta)
}

func resourceAliCloudOssAccessPointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssAccessPoint(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_access_point DescribeOssAccessPoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["NetworkOrigin"] != nil {
		d.Set("network_origin", objectRaw["NetworkOrigin"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["AccessPointName"] != nil {
		d.Set("access_point_name", objectRaw["AccessPointName"])
	}
	if objectRaw["Bucket"] != nil {
		d.Set("bucket", objectRaw["Bucket"])
	}

	publicAccessBlockConfigurationMaps := make([]map[string]interface{}, 0)
	publicAccessBlockConfigurationMap := make(map[string]interface{})
	publicAccessBlockConfiguration1Raw := make(map[string]interface{})
	if objectRaw["PublicAccessBlockConfiguration"] != nil {
		publicAccessBlockConfiguration1Raw = objectRaw["PublicAccessBlockConfiguration"].(map[string]interface{})
	}
	if len(publicAccessBlockConfiguration1Raw) > 0 {
		publicAccessBlockConfigurationMap["block_public_access"] = publicAccessBlockConfiguration1Raw["BlockPublicAccess"]

		publicAccessBlockConfigurationMaps = append(publicAccessBlockConfigurationMaps, publicAccessBlockConfigurationMap)
	}
	if objectRaw["PublicAccessBlockConfiguration"] != nil {
		if err := d.Set("public_access_block_configuration", publicAccessBlockConfigurationMaps); err != nil {
			return err
		}
	}
	vpcConfigurationMaps := make([]map[string]interface{}, 0)
	vpcConfigurationMap := make(map[string]interface{})
	vpcConfiguration1Raw := make(map[string]interface{})
	if objectRaw["VpcConfiguration"] != nil {
		vpcConfiguration1Raw = objectRaw["VpcConfiguration"].(map[string]interface{})
	}
	if len(vpcConfiguration1Raw) > 0 {
		vpcConfigurationMap["vpc_id"] = vpcConfiguration1Raw["VpcId"]

		vpcConfigurationMaps = append(vpcConfigurationMaps, vpcConfigurationMap)
	}
	if objectRaw["VpcConfiguration"] != nil {
		if err := d.Set("vpc_configuration", vpcConfigurationMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceAliCloudOssAccessPointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?publicAccessBlock")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(parts[0])
	query["x-oss-access-point-name"] = StringPointer(parts[1])

	if d.HasChange("public_access_block_configuration") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("public_access_block_configuration"); v != nil {
		blockPublicAccess1, _ := jsonpath.Get("$[0].block_public_access", v)
		if blockPublicAccess1 != nil && (d.HasChange("public_access_block_configuration.0.block_public_access") || blockPublicAccess1 != "") {
			objectDataLocalMap["BlockPublicAccess"] = blockPublicAccess1
		}

		request["PublicAccessBlockConfiguration"] = objectDataLocalMap
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", genXmlParam("PUT", "2019-05-17", "PutAccessPointPublicAccessBlock", action), query, body, nil, hostMap, false)
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

	return resourceAliCloudOssAccessPointRead(d, meta)
}

func resourceAliCloudOssAccessPointDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?accessPoint")
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(parts[0])
	query["x-oss-access-point-name"] = StringPointer(parts[1])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", genXmlParam("DELETE", "2019-05-17", "DeleteAccessPoint", action), query, nil, nil, hostMap, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"AccessPointCreatingConflict"}) || NeedRetry(err) {
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

	ossServiceV2 := OssServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ossServiceV2.OssAccessPointStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

package oss

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

func ResourceAccessPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessPointCreate,
		Read:   resourceAccessPointRead,
		Update: resourceAccessPointUpdate,
		Delete: resourceAccessPointDelete,
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
				Computed: true,
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

func resourceAccessPointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?accessPoint")
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(d.Get("bucket").(string))
	jsonString := "{}"
	jsonString, _ = sjson.Set(jsonString, "CreateAccessPointConfiguration.AccessPointName", d.Get("access_point_name"))
	err = json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		return errs.WrapError(err)
	}

	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOk("network_origin"); ok {
		objectDataLocalMap["NetworkOrigin"] = v
	}

	if v := d.Get("vpc_configuration"); !helper.IsNil(v) {
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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("PUT", "2019-05-17", "CreateAccessPoint", action), query, body, nil, hostMap, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_oss_access_point", action, errs.AlibabaCloudSdkGoERROR)
	}

	createAccessPointConfigurationAccessPointNameVar, _ := jsonpath.Get("CreateAccessPointConfiguration.AccessPointName", request)
	id := fmt.Sprintf("%s:%s", *hostMap["bucket"], createAccessPointConfigurationAccessPointNameVar)
	ossService := NewService(client)
	stateConf := helper.BuildStateConf([]string{"creating"}, []string{"enable"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossService.AccessPointStateRefreshFunc(id, "Status", []string{}))
	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, id)
	}

	d.SetId(id)

	return resourceAccessPointUpdate(d, meta)
}

func resourceAccessPointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossService := NewService(client)

	objectRaw, err := ossService.DescribeAccessPoint(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_access_point DescribeAccessPoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
	if len(vpcConfiguration1Raw) > 0 && vpcConfiguration1Raw["VpcId"] != nil && vpcConfiguration1Raw["VpcId"] != "" {
		vpcConfigurationMap["vpc_id"] = vpcConfiguration1Raw["VpcId"]

		vpcConfigurationMaps = append(vpcConfigurationMaps, vpcConfigurationMap)
	}
	// Only a VPC-origin access point carries a VpcId; any other origin reports an
	// empty or absent one. Write the empty list back for those origins so a block
	// persisted by an older provider version is cleared instead of forcing a
	// needless replacement, while a VPC-origin point keeps its state when the
	// response omits the VpcId.
	if len(vpcConfigurationMaps) > 0 || objectRaw["NetworkOrigin"] != "vpc" {
		if err := d.Set("vpc_configuration", vpcConfigurationMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceAccessPointUpdate(d *schema.ResourceData, meta interface{}) error {
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
	hostMap["bucket"] = helper.StringPointer(parts[0])
	query["x-oss-access-point-name"] = helper.StringPointer(parts[1])

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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.Do("Oss", helper.XmlParam("PUT", "2019-05-17", "PutAccessPointPublicAccessBlock", action), query, body, nil, hostMap, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAccessPointRead(d, meta)
}

func resourceAccessPointDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?accessPoint")
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(parts[0])
	query["x-oss-access-point-name"] = helper.StringPointer(parts[1])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("DELETE", "2019-05-17", "DeleteAccessPoint", action), query, nil, nil, hostMap, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"AccessPointCreatingConflict"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	ossService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ossService.AccessPointStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

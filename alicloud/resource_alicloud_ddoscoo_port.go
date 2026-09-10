// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDdosCooPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDdosCooPortCreate,
		Read:   resourceAliCloudDdosCooPortRead,
		Update: resourceAliCloudDdosCooPortUpdate,
		Delete: resourceAliCloudDdosCooPortDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backend_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"persistence_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"frontend_port": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"frontend_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"real_servers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudDdosCooPortCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePort"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["InstanceId"] = d.Get("instance_id")
	query["FrontendPort"] = d.Get("frontend_port")
	query["FrontendProtocol"] = d.Get("frontend_protocol")

	if v, ok := d.GetOk("backend_port"); ok {
		request["BackendPort"] = v
	}
	if v, ok := d.GetOk("real_servers"); ok {
		realServersMaps := v.([]interface{})
		request["RealServers"] = realServersMaps
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_port", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", query["InstanceId"], query["FrontendPort"], query["FrontendProtocol"]))

	return resourceAliCloudDdosCooPortUpdate(d, meta)
}

func resourceAliCloudDdosCooPortRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosCooServiceV2 := DdosCooServiceV2{client}

	objectRaw, err := ddosCooServiceV2.DescribeDdosCooPort(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddoscoo_port DescribeDdosCooPort Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["BackendPort"] != nil {
		d.Set("backend_port", objectRaw["BackendPort"])
	}
	if objectRaw["FrontendPort"] != nil {
		d.Set("frontend_port", objectRaw["FrontendPort"])
	}
	if objectRaw["FrontendProtocol"] != nil {
		d.Set("frontend_protocol", objectRaw["FrontendProtocol"])
	}
	if objectRaw["InstanceId"] != nil {
		d.Set("instance_id", objectRaw["InstanceId"])
	}

	realServers1Raw := make([]interface{}, 0)
	if objectRaw["RealServers"] != nil {
		realServers1Raw = objectRaw["RealServers"].([]interface{})
	}

	d.Set("real_servers", realServers1Raw)

	objectRaw, err = ddosCooServiceV2.DescribeDescribeNetworkRuleAttributes(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["FrontendPort"] != nil {
		d.Set("frontend_port", objectRaw["FrontendPort"])
	}
	if objectRaw["Protocol"] != nil {
		d.Set("frontend_protocol", objectRaw["Protocol"])
	}
	if objectRaw["InstanceId"] != nil {
		d.Set("instance_id", objectRaw["InstanceId"])
	}

	configMaps := make([]map[string]interface{}, 0)
	configMap := make(map[string]interface{})
	config1Raw := make(map[string]interface{})
	if objectRaw["Config"] != nil {
		config1Raw = objectRaw["Config"].(map[string]interface{})
	}
	if len(config1Raw) > 0 {
		configMap["persistence_timeout"] = config1Raw["PersistenceTimeout"]

		configMaps = append(configMaps, configMap)
	}
	if objectRaw["Config"] != nil {
		if err := d.Set("config", configMaps); err != nil {
			return err
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])
	d.Set("frontend_port", parts[1])
	d.Set("frontend_protocol", parts[2])

	return nil
}

func resourceAliCloudDdosCooPortUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	parts := strings.Split(d.Id(), ":")
	action := "ModifyPort"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FrontendPort"] = parts[1]
	query["InstanceId"] = parts[0]
	query["FrontendProtocol"] = parts[2]

	if !d.IsNewResource() && d.HasChange("real_servers") {
		update = true
	}
	if v, ok := d.GetOk("real_servers"); ok || d.HasChange("real_servers") {
		realServersMaps := v.([]interface{})
		request["RealServers"] = realServersMaps
	}

	if !d.IsNewResource() && d.HasChange("backend_port") {
		update = true
	}
	request["BackendPort"] = d.Get("backend_port")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, false)
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
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyNetworkRuleAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = parts[0]
	query["ForwardProtocol"] = parts[2]
	query["FrontendPort"] = parts[1]

	if d.HasChange("config") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("config"); v != nil {
		persistenceTimeout1, _ := jsonpath.Get("$[0].persistence_timeout", v)
		if persistenceTimeout1 != nil && (d.HasChange("config.0.persistence_timeout") || persistenceTimeout1 != "") {
			objectDataLocalMap["PersistenceTimeout"] = persistenceTimeout1
		}

		request["Config"] = convertObjectToJsonString(objectDataLocalMap)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, false)
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
	}

	d.Partial(false)
	return resourceAliCloudDdosCooPortRead(d, meta)
}

func resourceAliCloudDdosCooPortDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeletePort"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["FrontendPort"] = parts[1]
	query["InstanceId"] = parts[0]
	query["FrontendProtocol"] = parts[2]

	if v, ok := d.GetOk("real_servers"); ok {
		realServersMaps := v.([]interface{})
		request["RealServers"] = realServersMaps
	}

	if v, ok := d.GetOk("backend_port"); ok {
		request["BackendPort"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, query, request, false)

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

	return nil
}

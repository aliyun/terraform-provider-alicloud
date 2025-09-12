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

func resourceAliCloudNlbServerGroupServerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNlbServerGroupServerAttachmentCreate,
		Read:   resourceAliCloudNlbServerGroupServerAttachmentRead,
		Update: resourceAliCloudNlbServerGroupServerAttachmentUpdate,
		Delete: resourceAliCloudNlbServerGroupServerAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudNlbServerGroupServerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddServersToServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("server_group_id"); ok {
		request["ServerGroupId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOkExists("server_id"); ok {
		objectDataLocalMap["ServerId"] = v
	}

	if v, ok := d.GetOkExists("server_type"); ok {
		objectDataLocalMap["ServerType"] = v
	}

	if v, ok := d.GetOkExists("server_ip"); ok {
		objectDataLocalMap["ServerIp"] = v
	}

	if v, ok := d.GetOkExists("port"); ok {
		objectDataLocalMap["Port"] = v
	}

	if v, ok := d.GetOkExists("weight"); ok {
		objectDataLocalMap["Weight"] = v
	}

	if v, ok := d.GetOkExists("description"); ok {
		objectDataLocalMap["Description"] = v
	}

	ServersMap := make([]interface{}, 0)
	ServersMap = append(ServersMap, objectDataLocalMap)
	request["Servers"] = ServersMap

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "Servers.0.ServerId", d.Get("server_id"))
	jsonString, _ = sjson.Set(jsonString, "Servers.0.ServerType", d.Get("server_type"))
	jsonString, _ = sjson.Set(jsonString, "Servers.0.ServerIp", d.Get("server_ip"))
	jsonString, _ = sjson.Set(jsonString, "Servers.0.Port", d.Get("port"))
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.serverGroup", "OperationFailed.ResourceIsConfiguring", "Conflict.Lock"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_server_group_server_attachment", action, AlibabaCloudSdkGoERROR)
	}

	ServersServerIdVar, _ := jsonpath.Get("Servers[0].ServerId", request)
	ServersServerIpVar, _ := jsonpath.Get("Servers[0].ServerIp", request)
	ServersServerTypeVar, _ := jsonpath.Get("Servers[0].ServerType", request)
	ServersPortVar, _ := jsonpath.Get("Servers[0].Port", request)
	d.SetId(fmt.Sprintf("%v_%v_%v_%v_%v", response["ServerGroupId"], ServersServerIdVar, ServersServerIpVar, ServersServerTypeVar, ServersPortVar))

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.NlbServerGroupServerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNlbServerGroupServerAttachmentUpdate(d, meta)
}

func resourceAliCloudNlbServerGroupServerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}

	objectRaw, err := nlbServiceV2.DescribeNlbServerGroupServerAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_server_group_server_attachment DescribeNlbServerGroupServerAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("status", objectRaw["Status"])
	d.Set("weight", objectRaw["Weight"])
	d.Set("zone_id", objectRaw["ZoneId"])
	d.Set("port", objectRaw["Port"])
	d.Set("server_group_id", objectRaw["ServerGroupId"])
	d.Set("server_id", objectRaw["ServerId"])
	d.Set("server_ip", objectRaw["ServerIp"])
	d.Set("server_type", objectRaw["ServerType"])

	return nil
}

func resourceAliCloudNlbServerGroupServerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), "_")
	action := "UpdateServerGroupServersAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServerGroupId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	objectDataLocalMap := make(map[string]interface{})

	if d.HasChange("server_id") {
		update = true
	}
	if v, ok := d.GetOk("server_id"); ok {
		objectDataLocalMap["ServerId"] = v
	}

	if d.HasChange("server_type") {
		update = true
	}
	if v, ok := d.GetOk("server_type"); ok {
		objectDataLocalMap["ServerType"] = v
	}

	if d.HasChange("port") {
		update = true
	}
	if v, ok := d.GetOk("port"); ok {
		objectDataLocalMap["Port"] = v
	}

	if d.HasChange("weight") {
		update = true
	}
	if v, ok := d.GetOk("weight"); ok {
		objectDataLocalMap["Weight"] = v
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		objectDataLocalMap["Description"] = v
	}

	if d.HasChange("server_ip") {
		update = true
	}
	if v, ok := d.GetOk("server_ip"); ok {
		objectDataLocalMap["ServerIp"] = v
	}

	ServersMap := make([]interface{}, 0)
	ServersMap = append(ServersMap, objectDataLocalMap)
	request["Servers"] = ServersMap

	if len(parts) == 5 {
		jsonString := convertObjectToJsonString(request)
		request["ServerGroupId"] = parts[0]
		jsonString, _ = sjson.Set(jsonString, "Servers.0.ServerId", parts[1])
		jsonString, _ = sjson.Set(jsonString, "Servers.0.ServerType", parts[3])
		jsonString, _ = sjson.Set(jsonString, "Servers.0.Port", parts[4])
		jsonString, _ = sjson.Set(jsonString, "Servers.0.ServerIp", parts[2])
		_ = json.Unmarshal([]byte(jsonString), &request)
	} else {
		parts := strings.Split(d.Id(), ":")
		request["Servers.1.Port"] = parts[3]
		request["Servers.1.ServerId"] = parts[1]
		request["Servers.1.ServerType"] = parts[2]
		request["ServerGroupId"] = parts[0]
		if v, ok := d.GetOk("server_ip"); ok {
			request["Servers.1.ServerIp"] = v
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.serverGroup", "Conflict.Lock", "SystemBusy", "OperationFailed.ResourceIsConfiguring"}) || NeedRetry(err) {
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
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbServerGroupServerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudNlbServerGroupServerAttachmentRead(d, meta)
}

func resourceAliCloudNlbServerGroupServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), "_")
	action := "RemoveServersFromServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ServerGroupId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if len(parts) == 5 {
		jsonString := convertObjectToJsonString(request)
		jsonString, _ = sjson.Set(jsonString, "Servers.0.ServerId", parts[1])
		jsonString, _ = sjson.Set(jsonString, "Servers.0.ServerType", parts[3])
		jsonString, _ = sjson.Set(jsonString, "Servers.0.ServerIp", parts[2])
		jsonString, _ = sjson.Set(jsonString, "Servers.0.Port", parts[4])
		_ = json.Unmarshal([]byte(jsonString), &request)
	} else {
		parts := strings.Split(d.Id(), ":")
		request["ServerGroupId"] = parts[0]
		request["Servers.1.Port"] = parts[3]
		request["Servers.1.ServerId"] = parts[1]
		request["Servers.1.ServerType"] = parts[2]
		if v, ok := d.GetOk("server_ip"); ok {
			request["Servers.1.ServerIp"] = v
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.serverGroup", "Conflict.Lock"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.serverGroup", "ResourceNotFound.BackendServer"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nlbServiceV2.NlbServerGroupServerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func GetNlbServerGroupServerAttachment(client *connectivity.AliyunClient, id string) (object map[string]interface{}, err error) {
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 4 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 4, len(parts)))
	}
	action := "ListServerGroupServers"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServerGroupId"] = parts[0]
	request["ServerIds.1"] = parts[1]
	query["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Nlb", "2022-04-30", action, query, request, true)

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
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Servers[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Servers[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("ServerGroupServerAttachment", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Port"]) != parts[3] {
			continue
		}
		if item["ServerGroupId"] != parts[0] {
			continue
		}
		if item["ServerId"] != parts[1] {
			continue
		}
		if item["ServerType"] != parts[2] {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(NotFoundErr("ServerGroupServerAttachment", id), NotFoundMsg, response)
}

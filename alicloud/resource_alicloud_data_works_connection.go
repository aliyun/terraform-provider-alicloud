package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDataWorksConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDataWorksConnectionCreate,
		Read:   resourceAlicloudDataWorksConnectionRead,
		Update: resourceAlicloudDataWorksConnectionUpdate,
		Delete: resourceAlicloudDataWorksConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"connection_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"connection_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sub_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"env_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"content": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connect_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"binding_calc_engine_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gmt_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sequence": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_engine": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDataWorksConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateConnection"
	request := map[string]interface{}{
		"RegionId":       client.RegionId,
		"ProjectId":      d.Get("project_id"),
		"Name":           d.Get("connection_name"),
		"ConnectionType": d.Get("connection_type"),
		"EnvType":        d.Get("env_type"),
		"Content":        d.Get("content"),
	}
	if v, ok := d.GetOk("sub_type"); ok && v.(string) != "" {
		request["SubType"] = v
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request["Description"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_connection", action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_data_works_connection", "$.Data", response)
	}
	connectionId := fmt.Sprintf("%v", v)
	d.SetId(fmt.Sprintf("%d:%s", d.Get("project_id").(int), connectionId))
	return resourceAlicloudDataWorksConnectionRead(d, meta)
}

func resourceAlicloudDataWorksConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	projectId, _ := strconv.Atoi(parts[0])
	d.Set("project_id", projectId)
	d.Set("connection_id", object["Id"])
	d.Set("connection_name", object["Name"])
	d.Set("connection_type", object["ConnectionType"])
	d.Set("sub_type", object["SubType"])
	d.Set("env_type", object["EnvType"])
	d.Set("content", object["Content"])
	d.Set("description", object["Description"])
	if object["Status"] != nil {
		d.Set("status", fmt.Sprintf("%v", object["Status"]))
	}
	d.Set("operator", object["Operator"])
	d.Set("connect_status", object["ConnectStatus"])
	d.Set("binding_calc_engine_id", object["BindingCalcEngineId"])
	d.Set("gmt_modified", object["GmtModified"])
	d.Set("sequence", object["Sequence"])
	d.Set("shared", object["Shared"])
	d.Set("default_engine", object["DefaultEngine"])
	d.Set("create_time", object["GmtCreate"])
	d.Set("tenant_id", object["TenantId"])
	d.Set("region_id", client.RegionId)
	return nil
}

func resourceAlicloudDataWorksConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	connectionId := parts[1]
	action := "UpdateConnection"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ConnectionId": connectionId,
		"EnvType":      d.Get("env_type"),
		"Content":      d.Get("content"),
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
	return resourceAlicloudDataWorksConnectionRead(d, meta)
}

func resourceAlicloudDataWorksConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	connectionId := parts[1]
	action := "DeleteDataSource"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DataSourceId": connectionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
		if NotFoundError(err) || IsExpectedErrors(err, []string{"Invalid.Tenant.ConnectionNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *DataworksPublicService) DescribeDataWorksConnection(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}
	projectId, _ := strconv.Atoi(parts[0])
	connectionId := parts[1]
	action := "ListConnections"
	request := map[string]interface{}{
		"ProjectId":  projectId,
		"PageNumber": 1,
		"PageSize":   100,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	client := s.client
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Invalid.Tenant.ConnectionNotExists"}) {
			return object, WrapErrorf(NotFoundErr("dataworks", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Data.Connections", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.Connections", response)
	}
	connections, ok := v.([]interface{})
	if !ok {
		return object, WrapError(fmt.Errorf("Connections is not an array for connection %s", id))
	}
	for _, conn := range connections {
		c, ok := conn.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprintf("%v", c["Id"]) == connectionId {
			object = c
			break
		}
	}
	if len(object) < 1 {
		return object, WrapErrorf(NotFoundErr("dataworks", id), NotFoundWithResponse, response)
	}
	return object, nil
}

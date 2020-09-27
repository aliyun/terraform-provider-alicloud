package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	dms_enterprise "github.com/aliyun/alibaba-cloud-sdk-go/services/dms-enterprise"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDmsEnterpriseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDmsEnterpriseInstanceCreate,
		Read:   resourceAlicloudDmsEnterpriseInstanceRead,
		Update: resourceAlicloudDmsEnterpriseInstanceUpdate,
		Delete: resourceAlicloudDmsEnterpriseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"data_link_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dba_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dba_nick_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dba_uid": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ddl_online": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
			},
			"ecs_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ecs_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"env_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"export_timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_alias"},
			},
			"instance_alias": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'instance_alias' has been deprecated from version 1.100.0. Use 'instance_name' instead.",
				ConflictsWith: []string{"instance_name"},
			},
			"instance_source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"query_timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"safe_rule": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"safe_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"skip_test": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:          schema.TypeString,
				Computed:      true,
				ConflictsWith: []string{"state"},
			},
			"state": {
				Type:          schema.TypeString,
				Computed:      true,
				Deprecated:    "Field 'state' has been deprecated from version 1.100.0. Use 'status' instead.",
				ConflictsWith: []string{"status"},
			},
			"tid": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"use_dsql": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDmsEnterpriseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := dms_enterprise.CreateRegisterInstanceRequest()
	if v, ok := d.GetOk("data_link_name"); ok {
		request.DataLinkName = v.(string)
	}

	request.DatabasePassword = d.Get("database_password").(string)
	request.DatabaseUser = d.Get("database_user").(string)
	request.DbaUid = requests.NewInteger(d.Get("dba_uid").(int))
	if v, ok := d.GetOk("ddl_online"); ok {
		request.DdlOnline = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("ecs_instance_id"); ok {
		request.EcsInstanceId = v.(string)
	}

	if v, ok := d.GetOk("ecs_region"); ok {
		request.EcsRegion = v.(string)
	}

	request.EnvType = d.Get("env_type").(string)
	request.ExportTimeout = requests.NewInteger(d.Get("export_timeout").(int))
	request.Host = d.Get("host").(string)
	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceAlias = v.(string)
	} else if v, ok := d.GetOk("instance_alias"); ok {
		request.InstanceAlias = v.(string)
	}

	request.InstanceSource = d.Get("instance_source").(string)
	request.InstanceType = d.Get("instance_type").(string)
	request.NetworkType = d.Get("network_type").(string)
	request.Port = requests.NewInteger(d.Get("port").(int))
	request.QueryTimeout = requests.NewInteger(d.Get("query_timeout").(int))
	request.SafeRule = d.Get("safe_rule").(string)
	if v, ok := d.GetOk("sid"); ok {
		request.Sid = v.(string)
	}

	if v, ok := d.GetOkExists("skip_test"); ok {
		request.SkipTest = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("tid"); ok {
		request.Tid = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("use_dsql"); ok {
		request.UseDsql = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = v.(string)
	}

	wait := incrementalWait(3*time.Second, 2*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
			return dms_enterpriseClient.RegisterInstance(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RegisterInstanceFailure"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		d.SetId(fmt.Sprintf("%v:%v", request.Host, request.Port))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dms_enterprise_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudDmsEnterpriseInstanceUpdate(d, meta)
}
func resourceAlicloudDmsEnterpriseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dms_enterpriseService := Dms_enterpriseService{client}
	object, err := dms_enterpriseService.DescribeDmsEnterpriseInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_enterprise_instance dms_enterpriseService.DescribeDmsEnterpriseInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("host", parts[0])
	d.Set("port", parts[1])
	d.Set("data_link_name", object.DataLinkName)
	d.Set("database_password", object.DatabasePassword)
	d.Set("database_user", object.DatabaseUser)
	d.Set("dba_id", object.DbaId)
	d.Set("ddl_online", object.DdlOnline)
	d.Set("ecs_instance_id", object.EcsInstanceId)
	d.Set("ecs_region", object.EcsRegion)
	d.Set("env_type", object.EnvType)
	d.Set("export_timeout", object.ExportTimeout)
	d.Set("instance_id", object.InstanceId)
	d.Set("instance_name", object.InstanceAlias)
	d.Set("instance_alias", object.InstanceAlias)
	d.Set("instance_source", object.InstanceSource)
	d.Set("instance_type", object.InstanceType)
	d.Set("query_timeout", object.QueryTimeout)
	d.Set("safe_rule_id", object.SafeRuleId)
	d.Set("sid", object.Sid)
	d.Set("status", object.State)
	d.Set("state", object.State)
	d.Set("use_dsql", object.UseDsql)
	d.Set("vpc_id", object.VpcId)
	d.Set("dba_nick_name", object.DbaNickName)
	return nil
}
func resourceAlicloudDmsEnterpriseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := dms_enterprise.CreateUpdateInstanceRequest()
	request.Host = parts[0]
	if v, err := strconv.Atoi(parts[1]); err == nil {
		request.Port = requests.NewInteger(v)
	} else {
		return WrapError(err)
	}
	if !d.IsNewResource() && d.HasChange("database_password") {
		update = true
	}
	request.DatabasePassword = d.Get("database_password").(string)
	if !d.IsNewResource() && d.HasChange("database_user") {
		update = true
	}
	request.DatabaseUser = d.Get("database_user").(string)
	if d.HasChange("dba_id") {
		update = true
	}
	request.DbaId = d.Get("dba_id").(string)
	if !d.IsNewResource() && d.HasChange("env_type") {
		update = true
	}
	request.EnvType = d.Get("env_type").(string)
	if !d.IsNewResource() && d.HasChange("export_timeout") {
		update = true
	}
	request.ExportTimeout = requests.NewInteger(d.Get("export_timeout").(int))
	if d.HasChange("instance_id") {
		update = true
	}
	request.InstanceId = d.Get("instance_id").(string)
	if !d.IsNewResource() && (d.HasChange("instance_name") || d.HasChange("instance_alias")) {
		update = true
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceAlias = v.(string)
	} else {
		request.InstanceAlias = d.Get("instance_alias").(string)
	}
	if !d.IsNewResource() && d.HasChange("instance_source") {
		update = true
	}
	request.InstanceSource = d.Get("instance_source").(string)
	if !d.IsNewResource() && d.HasChange("instance_type") {
		update = true
	}
	request.InstanceType = d.Get("instance_type").(string)
	if !d.IsNewResource() && d.HasChange("query_timeout") {
		update = true
	}
	request.QueryTimeout = requests.NewInteger(d.Get("query_timeout").(int))
	if d.HasChange("safe_rule_id") {
		update = true
	}
	request.SafeRuleId = d.Get("safe_rule_id").(string)
	if !d.IsNewResource() && d.HasChange("data_link_name") {
		update = true
		request.DataLinkName = d.Get("data_link_name").(string)
	}
	if !d.IsNewResource() && d.HasChange("ddl_online") {
		update = true
		request.DdlOnline = requests.NewInteger(d.Get("ddl_online").(int))
	}
	if !d.IsNewResource() && d.HasChange("ecs_instance_id") {
		update = true
		request.EcsInstanceId = d.Get("ecs_instance_id").(string)
	}
	if !d.IsNewResource() && d.HasChange("ecs_region") {
		update = true
		request.EcsRegion = d.Get("ecs_region").(string)
	}
	if !d.IsNewResource() && d.HasChange("sid") {
		update = true
		request.Sid = d.Get("sid").(string)
	}
	if !d.IsNewResource() && d.HasChange("use_dsql") {
		update = true
		request.UseDsql = requests.NewInteger(d.Get("use_dsql").(int))
	}
	if !d.IsNewResource() && d.HasChange("vpc_id") {
		update = true
		request.VpcId = d.Get("vpc_id").(string)
	}
	if update {
		if _, ok := d.GetOkExists("skip_test"); ok {
			request.SkipTest = requests.NewBoolean(d.Get("skip_test").(bool))
		}
		if _, ok := d.GetOk("tid"); ok {
			request.Tid = requests.NewInteger(d.Get("tid").(int))
		}
		raw, err := client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
			return dms_enterpriseClient.UpdateInstance(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudDmsEnterpriseInstanceRead(d, meta)
}
func resourceAlicloudDmsEnterpriseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := dms_enterprise.CreateDeleteInstanceRequest()
	request.Host = parts[0]
	if v, err := strconv.Atoi(parts[1]); err == nil {
		request.Port = requests.NewInteger(v)
	} else {
		return WrapError(err)
	}
	if v, ok := d.GetOk("sid"); ok {
		request.Sid = v.(string)
	}
	if v, ok := d.GetOk("tid"); ok {
		request.Tid = requests.NewInteger(v.(int))
	}
	raw, err := client.WithDmsEnterpriseClient(func(dms_enterpriseClient *dms_enterprise.Client) (interface{}, error) {
		return dms_enterpriseClient.DeleteInstance(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNoEnoughNumber"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

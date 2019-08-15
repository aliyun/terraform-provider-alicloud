package alicloud

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/fc-go-sdk"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudFCService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCServiceCreate,
		Read:   resourceAlicloudFCServiceRead,
		Update: resourceAlicloudFCServiceUpdate,
		Delete: resourceAlicloudFCServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validateStringLengthInRange(1, 128),
			},
			"name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					// uuid is 26 characters, limit the prefix to 229.
					value := v.(string)
					if len(value) > 122 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 102 characters, name is limited to 128", k))
					}
					return
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"internet_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project": {
							Type:     schema.TypeString,
							Required: true,
						},
						"logstore": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"vpc_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudFCServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		name = resource.PrefixedUniqueId(v.(string))
	} else {
		name = resource.UniqueId()
	}

	project, logstore, err := parseLogConfig(d, meta)
	if err != nil {
		return WrapError(err)
	}
	request := &fc.CreateServiceInput{
		ServiceName:    StringPointer(name),
		Description:    StringPointer(d.Get("description").(string)),
		InternetAccess: BoolPointer(d.Get("internet_access").(bool)),
		Role:           StringPointer(d.Get("role").(string)),
		LogConfig: &fc.LogConfig{
			Project:  StringPointer(project),
			Logstore: StringPointer(logstore),
		},
	}
	vpcconfig, err := parseVpcConfig(d, meta)
	if err != nil {
		return WrapError(err)
	}
	request.VPCConfig = vpcconfig
	var requestInfo *fc.Client
	var response *fc.CreateServiceOutput
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.CreateService(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AccessDenied, "does not exist"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateService", raw, requestInfo, request)
		response, _ = raw.(*fc.CreateServiceOutput)
		return nil

	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_service", "CreateService", FcGoSdk)
	}

	d.SetId(*response.ServiceName)
	return resourceAlicloudFCServiceRead(d, meta)
}

func resourceAlicloudFCServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	object, err := fcService.DescribeFcService(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.ServiceName)
	d.Set("service_id", object.ServiceID)
	d.Set("description", object.Description)
	d.Set("internet_access", object.InternetAccess)
	d.Set("role", object.Role)
	var logConfigs []map[string]interface{}
	if logconfig := object.LogConfig; logconfig != nil && *logconfig.Project != "" {
		logConfigs = append(logConfigs, map[string]interface{}{
			"project":  *logconfig.Project,
			"logstore": *logconfig.Logstore,
		})
	}
	if err := d.Set("log_config", logConfigs); err != nil {
		return WrapError(err)
	}
	var vpcConfigs []map[string]interface{}
	if vpcConfig := object.VPCConfig; vpcConfig != nil && *vpcConfig.VPCID != "" {
		vpcConfigs = append(vpcConfigs, map[string]interface{}{
			"vswitch_ids":       schema.NewSet(schema.HashString, flattenStringList(vpcConfig.VSwitchIDs)),
			"security_group_id": *vpcConfig.SecurityGroupID,
			"vpc_id":            *vpcConfig.VPCID,
		})
	}
	if err := d.Set("vpc_config", vpcConfigs); err != nil {
		return WrapError(err)
	}
	d.Set("last_modified", object.LastModifiedTime)

	return nil
}

func resourceAlicloudFCServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)
	request := &fc.UpdateServiceInput{}

	if d.HasChange("role") {
		request.Role = StringPointer(d.Get("role").(string))
		d.SetPartial("role")
	}
	if d.HasChange("internet_access") {
		request.InternetAccess = BoolPointer(d.Get("internet_access").(bool))
		d.SetPartial("internet_access")
	}
	if d.HasChange("description") {
		request.Description = StringPointer(d.Get("description").(string))
		d.SetPartial("description")
	}
	if d.HasChange("log_config") {
		project, logstore, err := parseLogConfig(d, meta)
		if err != nil {
			return WrapError(err)
		}
		request.LogConfig = &fc.LogConfig{
			Project:  StringPointer(project),
			Logstore: StringPointer(logstore),
		}
		d.SetPartial("log_config")
	}

	if d.HasChange("vpc_config") {
		vpcconfig, err := parseVpcConfig(d, meta)
		if err != nil {
			return WrapError(err)
		}
		request.VPCConfig = vpcconfig
		d.SetPartial("vpc_config")
	}

	if request != nil {
		request.ServiceName = StringPointer(d.Id())
		var requestInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.UpdateService(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateService", FcGoSdk)
		}
		addDebug("UpdateService", raw, requestInfo, request)
	}

	d.Partial(false)
	return resourceAlicloudFCServiceRead(d, meta)
}

func resourceAlicloudFCServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}
	request := &fc.DeleteServiceInput{
		ServiceName: StringPointer(d.Id()),
	}
	var requestInfo *fc.Client
	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.DeleteService(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteService", FcGoSdk)
	}
	addDebug("DeleteService", raw, requestInfo, request)
	return WrapError(fcService.WaitForFcService(d.Id(), Deleted, DefaultTimeout))

}

func parseVpcConfig(d *schema.ResourceData, meta interface{}) (config *fc.VPCConfig, err error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	if v, ok := d.GetOk("vpc_config"); ok {

		confs := v.([]interface{})
		conf, ok := confs[0].(map[string]interface{})

		if !ok {
			return
		}
		if role, ok := d.GetOk("role"); !ok || role.(string) == "" {
			err = WrapError(Error("'role' is required when 'vpc_config' is set."))
			return
		}
		if conf != nil {
			vswitch_ids := conf["vswitch_ids"].(*schema.Set).List()
			vsw, e := vpcService.DescribeVSwitch(vswitch_ids[0].(string))
			if e != nil {
				err = WrapError(e)
				return
			}
			config = &fc.VPCConfig{
				VSwitchIDs:      expandStringList(vswitch_ids),
				SecurityGroupID: StringPointer(conf["security_group_id"].(string)),
				VPCID:           StringPointer(vsw.VpcId),
			}
		}
	}
	return
}

func parseLogConfig(d *schema.ResourceData, meta interface{}) (project, logstore string, err error) {
	client := meta.(*connectivity.AliyunClient)
	if v, ok := d.GetOk("log_config"); ok {

		configs := v.([]interface{})
		config, ok := configs[0].(map[string]interface{})

		if !ok {
			return
		}

		if config != nil {
			project = config["project"].(string)
			logstore = config["logstore"].(string)
		}
	}
	if project != "" {
		var requestInfo *sls.Client
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, e := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				requestInfo = slsClient
				return slsClient.CheckProjectExist(project)
			})
			if e != nil {
				if NotFoundError(e) {
					return resource.RetryableError(e)
				}
				return resource.NonRetryableError(e)
			}
			addDebug("CheckProjectExist", raw, requestInfo, project)
			return nil
		})
	}

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, d.Id(), "CheckProjectExist", FcGoSdk)
		return
	}

	if logstore != "" {
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, e := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				return slsClient.CheckLogstoreExist(project, logstore)
			})
			if e != nil {
				if NotFoundError(e) {
					return resource.RetryableError(e)
				}
				return resource.NonRetryableError(e)
			}
			addDebug("CheckLogstoreExist", raw)
			return nil
		})
	}
	return
}

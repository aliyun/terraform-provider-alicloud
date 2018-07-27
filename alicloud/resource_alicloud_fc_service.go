package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/fc-go-sdk"
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
			"name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validateStringLengthInRange(1, 128),
			},
			"name_prefix": &schema.Schema{
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

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"internet_access": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"role": &schema.Schema{
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
		},
	}
}

func resourceAlicloudFCServiceCreate(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}
	client := meta.(*AliyunClient)
	conn := client.fcconn

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
		return err
	}
	input := &fc.CreateServiceInput{
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
		return err
	}
	input.VPCConfig = vpcconfig

	var service *fc.CreateServiceOutput
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		service, err = conn.CreateService(input)
		if err != nil {
			if IsExceptedErrors(err, []string{AccessDenied, "does not exist"}) {
				return resource.RetryableError(fmt.Errorf("Error creating function compute service got an error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Error creating function compute service got an error: %#v", err))
		}
		return nil

	}); err != nil {
		return err
	}

	if service == nil {
		return fmt.Errorf("Creating function compute service got a empty response: %#v.", service)
	}

	d.SetId(*service.ServiceName)

	return resourceAlicloudFCServiceRead(d, meta)
}

func resourceAlicloudFCServiceRead(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}

	client := meta.(*AliyunClient)

	service, err := client.DescribeFcService(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("DescribeFCService %s got an error: %#v", d.Id(), err)
	}

	d.Set("name", service.ServiceName)
	d.Set("description", service.Description)
	d.Set("internet_access", service.InternetAccess)
	d.Set("role", service.Role)
	var logConfigs []map[string]interface{}
	if logconfig := service.LogConfig; logconfig != nil && *logconfig.Project != "" {
		logConfigs = append(logConfigs, map[string]interface{}{
			"project":  *logconfig.Project,
			"logstore": *logconfig.Logstore,
		})
	}
	if err := d.Set("log_config", logConfigs); err != nil {
		return err
	}
	var vpcConfigs []map[string]interface{}
	if vpcConfig := service.VPCConfig; vpcConfig != nil && *vpcConfig.VPCID != "" {
		vpcConfigs = append(vpcConfigs, map[string]interface{}{
			"vswitch_ids":       schema.NewSet(schema.HashString, flattenStringList(vpcConfig.VSwitchIDs)),
			"security_group_id": *vpcConfig.SecurityGroupID,
			"vpc_id":            *vpcConfig.VPCID,
		})
	}
	if err := d.Set("vpc_config", vpcConfigs); err != nil {
		return err
	}
	d.Set("last_modified", service.LastModifiedTime)

	return nil
}

func resourceAlicloudFCServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}
	client := meta.(*AliyunClient)

	d.Partial(true)
	updateInput := &fc.UpdateServiceInput{}

	if d.HasChange("role") {
		updateInput.Role = StringPointer(d.Get("role").(string))
		d.SetPartial("role")
	}
	if d.HasChange("internet_access") {
		updateInput.InternetAccess = BoolPointer(d.Get("internet_access").(bool))
		d.SetPartial("internet_access")
	}
	if d.HasChange("description") {
		updateInput.Description = StringPointer(d.Get("description").(string))
		d.SetPartial("description")
	}
	if d.HasChange("log_config") {
		project, logstore, err := parseLogConfig(d, meta)
		if err != nil {
			return err
		}
		updateInput.LogConfig.Project = StringPointer(project)
		updateInput.LogConfig.Logstore = StringPointer(logstore)
		d.SetPartial("log_config")
	}

	if d.HasChange("vpc_config") {
		vpcconfig, err := parseVpcConfig(d, meta)
		if err != nil {
			return err
		}
		updateInput.VPCConfig = vpcconfig
		d.SetPartial("vpc_config")
	}

	if updateInput != nil {
		updateInput.ServiceName = StringPointer(d.Id())
		if _, err := client.fcconn.UpdateService(updateInput); err != nil {
			return fmt.Errorf("UpdateService %s got an error: %#v.", d.Id(), err)
		}
	}

	d.Partial(false)
	return resourceAlicloudFCServiceRead(d, meta)
}

func resourceAlicloudFCServiceDelete(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}
	client := meta.(*AliyunClient)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.fcconn.DeleteService(&fc.DeleteServiceInput{
			ServiceName: StringPointer(d.Id()),
		}); err != nil {
			if IsExceptedErrors(err, []string{ServiceNotFound}) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting function service got an error: %#v.", err))
		}

		if _, err := client.DescribeFcService(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("While deleting service, getting service %s got an error: %#v.", d.Id(), err))
		}
		return nil
	})

}

func parseVpcConfig(d *schema.ResourceData, meta interface{}) (config *fc.VPCConfig, err error) {
	if v, ok := d.GetOk("vpc_config"); ok {

		confs := v.([]interface{})
		conf, ok := confs[0].(map[string]interface{})

		if !ok {
			return
		}
		if role, ok := d.GetOk("role"); !ok || role.(string) == "" {
			err = fmt.Errorf("'role' is required when 'vpc_config' is set.")
			return
		}
		if conf != nil {
			vswitch_ids := conf["vswitch_ids"].(*schema.Set).List()
			vsw, e := meta.(*AliyunClient).DescribeVswitch(vswitch_ids[0].(string))
			if e != nil {
				err = fmt.Errorf("While creating fc service, describing vswitch %s got an error: %#v.", vswitch_ids[0].(string), e)
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
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			if _, e := meta.(*AliyunClient).logconn.CheckProjectExist(project); e != nil {
				if NotFoundError(e) {
					return resource.RetryableError(fmt.Errorf("Check log project %s failed: %#v.", project, e))
				}
				return resource.NonRetryableError(fmt.Errorf("Check log project %s failed: %#v.", project, e))
			}
			return nil
		})
	}

	if err != nil {
		return
	}

	if logstore != "" {
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			if _, e := meta.(*AliyunClient).logconn.CheckLogstoreExist(project, logstore); e != nil {
				if NotFoundError(e) {
					return resource.RetryableError(fmt.Errorf("Check logstore %s failed: %#v.", logstore, e))
				}
				return resource.NonRetryableError(fmt.Errorf("Check logstore %s failed: %#v.", logstore, e))
			}
			return nil
		})
	}
	return
}

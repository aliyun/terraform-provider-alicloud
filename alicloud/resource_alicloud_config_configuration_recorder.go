package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudConfigConfigurationRecorder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigConfigurationRecorderCreate,
		Read:   resourceAlicloudConfigConfigurationRecorderRead,
		Update: resourceAlicloudConfigConfigurationRecorderUpdate,
		Delete: resourceAlicloudConfigConfigurationRecorderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enterprise_edition": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"organization_enable_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_master_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudConfigConfigurationRecorderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := config.CreateStartConfigurationRecorderRequest()
	if v, ok := d.GetOkExists("enterprise_edition"); ok {
		request.EnterpriseEdition = requests.NewBoolean(v.(bool))
	}

	raw, err := client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
		return configClient.StartConfigurationRecorder(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_configuration_recorder", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*config.StartConfigurationRecorderResponse)
	d.SetId(fmt.Sprintf("%v", response.ConfigurationRecorder.AccountId))

	return resourceAlicloudConfigConfigurationRecorderUpdate(d, meta)
}
func resourceAlicloudConfigConfigurationRecorderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigConfigurationRecorder(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_configuration_recorder configService.DescribeConfigConfigurationRecorder Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("organization_enable_status", object.OrganizationEnableStatus)
	d.Set("organization_master_id", object.OrganizationMasterId)
	d.Set("resource_types", object.ResourceTypes)
	d.Set("status", object.ConfigurationRecorderStatus)
	return nil
}
func resourceAlicloudConfigConfigurationRecorderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	update := false
	request := config.CreatePutConfigurationRecorderRequest()
	if d.HasChange("resource_types") {
		update = true
	}
	request.ResourceTypes = convertListToCommaSeparate(d.Get("resource_types").(*schema.Set).List())

	if update {
		raw, err := client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
			return configClient.PutConfigurationRecorder(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"REGISTERED"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, configService.ConfigConfigurationRecorderStateRefreshFunc(d.Id(), []string{"REGISTRABLE"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudConfigConfigurationRecorderRead(d, meta)
}
func resourceAlicloudConfigConfigurationRecorderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudConfigConfigurationRecorder. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

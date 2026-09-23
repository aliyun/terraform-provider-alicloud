package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallUserAlarmConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallUserAlarmConfigCreate,
		Read:   resourceAliCloudCloudFirewallUserAlarmConfigRead,
		Update: resourceAliCloudCloudFirewallUserAlarmConfigUpdate,
		Delete: resourceAliCloudCloudFirewallUserAlarmConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alarm_lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_default_contact": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_config": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_period": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"alarm_hour": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"alarm_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"alarm_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"alarm_week_day": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"alarm_notify": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"contact_config": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"notify_config"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mobile_phone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"notify_config": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"contact_config"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notify_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"notify_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudCloudFirewallUserAlarmConfigCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceAliCloudCloudFirewallUserAlarmConfigUpdate(d, meta)
}

func resourceAliCloudCloudFirewallUserAlarmConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallUserAlarmConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_user_alarm_config DescribeCloudFirewallUserAlarmConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alarm_lang", objectRaw["AlarmLang"])

	alarmConfigMaps := make([]map[string]interface{}, 0)
	if cfg, ok := objectRaw["AlarmConfig"]; ok {
		for _, alarmConfigChildRaw := range convertToInterfaceArray(cfg) {
			alarmConfigRaw := alarmConfigChildRaw.(map[string]interface{})
			alarmConfigMap := map[string]interface{}{}
			if v, ok := alarmConfigRaw["AlarmHour"]; ok {
				alarmConfigMap["alarm_hour"] = v
			}
			if v, ok := alarmConfigRaw["AlarmNotify"]; ok {
				alarmConfigMap["alarm_notify"] = v
			}
			if v, ok := alarmConfigRaw["AlarmPeriod"]; ok {
				alarmConfigMap["alarm_period"] = v
			}
			if v, ok := alarmConfigRaw["AlarmType"]; ok {
				alarmConfigMap["alarm_type"] = v
			}
			if v, ok := alarmConfigRaw["AlarmValue"]; ok {
				alarmConfigMap["alarm_value"] = v
			}
			if v, ok := alarmConfigRaw["AlarmWeekDay"]; ok {
				alarmConfigMap["alarm_week_day"] = v
			}
			alarmConfigMaps = append(alarmConfigMaps, alarmConfigMap)
		}
	}
	if err := d.Set("alarm_config", alarmConfigMaps); err != nil {
		return err
	}

	contactConfigMaps := make([]map[string]interface{}, 0)
	if cfg, ok := objectRaw["ContactConfig"]; ok {
		for _, contactConfigChildRaw := range convertToInterfaceArray(cfg) {
			contactConfigMap := make(map[string]interface{})
			contactConfigRaw := contactConfigChildRaw.(map[string]interface{})
			if v, ok := contactConfigRaw["Email"]; ok {
				contactConfigMap["email"] = v
			}
			if v, ok := contactConfigRaw["MobilePhone"]; ok {
				contactConfigMap["mobile_phone"] = v
			}
			if v, ok := contactConfigRaw["Name"]; ok {
				contactConfigMap["name"] = v
			}
			if v, ok := contactConfigRaw["Status"]; ok {
				contactConfigMap["status"] = v
			}
			contactConfigMaps = append(contactConfigMaps, contactConfigMap)
		}
	}
	if err := d.Set("contact_config", contactConfigMaps); err != nil {
		return err
	}

	notifyConfigMaps := make([]map[string]interface{}, 0)
	if cfg, ok := objectRaw["NotifyConfig"]; ok {
		for _, notifyConfigChildRaw := range convertToInterfaceArray(cfg) {
			notifyConfigRaw := notifyConfigChildRaw.(map[string]interface{})
			notifyConfigMap := make(map[string]interface{})
			if v, ok := notifyConfigRaw["NotifyValue"]; ok {
				notifyConfigMap["notify_value"] = v
			}
			if v, ok := notifyConfigRaw["NotifyType"]; ok {
				notifyConfigRaw["notify_type"] = v
			}
			notifyConfigMaps = append(notifyConfigMaps, notifyConfigMap)
		}
	}
	if err := d.Set("notify_config", notifyConfigMaps); err != nil {
		return err
	}

	accountId, _ := client.AccountId()
	d.SetId(accountId)
	return nil
}

func resourceAliCloudCloudFirewallUserAlarmConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyUserAlarmConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("use_default_contact"); ok {
		request["UseDefaultContact"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if d.HasChange("alarm_lang") {
		update = true
		request["AlarmLang"] = d.Get("alarm_lang")
	}

	configDataMap := make(map[string]interface{})
	if d.HasChange("alarm_config") {
		aConfigRaw := d.Get("alarm_config").([]interface{})
		if len(aConfigRaw) > 0 {
			update = true
			cfg := aConfigRaw[0].(map[string]interface{})
			if v, ok := cfg["alarm_value"]; ok {
				configDataMap["AlarmValue"] = v
			}
			if v, ok := cfg["alarm_period"]; ok {
				configDataMap["AlarmPeriod"] = v
			}
			if v, ok := cfg["alarm_week_day"]; ok {
				configDataMap["AlarmWeekDay"] = v
			}
			if v, ok := cfg["alarm_notify"]; ok {
				configDataMap["AlarmNotify"] = v
			}
			if v, ok := cfg["alarm_type"]; ok {
				configDataMap["AlarmType"] = v
			}
			if v, ok := cfg["alarm_hour"]; ok {
				configDataMap["AlarmHour"] = v
			}
		}
	}
	request["AlarmConfig"] = []interface{}{configDataMap}

	contactConfigDataMap := make(map[string]interface{})
	if d.HasChange("contact_config") {
		contactConfigRaw := d.Get("contact_config").([]interface{})
		if len(contactConfigRaw) > 0 {
			update = true
			cfg := contactConfigRaw[0].(map[string]interface{})
			if v, ok := cfg["mobile_phone"]; ok && v != "" {
				contactConfigDataMap["MobilePhone"] = v
			}
			if v, ok := cfg["email"]; ok && v != "" {
				contactConfigDataMap["Email"] = v
			}
			if v, ok := cfg["status"]; ok && v != "" {
				contactConfigDataMap["Status"] = v
			}
			if v, ok := cfg["name"]; ok && v != "" {
				contactConfigDataMap["Name"] = v
			}
		}
	}
	request["ContactConfig"] = []interface{}{contactConfigDataMap}

	notifyConfigMap := make(map[string]interface{})
	if d.HasChange("notify_config") {
		notifyConfigRaw := d.Get("notify_config").([]interface{})
		if len(notifyConfigRaw) > 0 {
			update = true
			cfg := notifyConfigRaw[0].(map[string]interface{})
			if v, ok := cfg["notify_type"]; ok {
				notifyConfigMap["NotifyType"] = v
			}
			if v, ok := cfg["notify_value"]; ok {
				notifyConfigMap["NotifyValue"] = v
			}
		}
	}
	request["NotifyConfig"] = []interface{}{notifyConfigMap}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
	return resourceAliCloudCloudFirewallUserAlarmConfigRead(d, meta)
}

func resourceAliCloudCloudFirewallUserAlarmConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource User Alarm Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

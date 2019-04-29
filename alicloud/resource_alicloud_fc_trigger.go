package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"strings"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudFCTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCTriggerCreate,
		Read:   resourceAlicloudFCTriggerRead,
		Update: resourceAlicloudFCTriggerUpdate,
		Delete: resourceAlicloudFCTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"function": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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

			"role": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_arn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"config": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// The read config is json rawMessage and it does not contains space and enter.
					if d.Get("type").(string) == string(fc.TRIGGER_TYPE_MNS_TOPIC) {
						return true
					}
					return old == removeSpaceAndEnter(new)
				},
				ValidateFunc: validateJsonString,
			},
			//Modifying config is not supported when type is mns_topic
			"config_mns": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// The read config is json rawMessage and it does not contains space and enter.
					return old == removeSpaceAndEnter(new)
				},
				ValidateFunc:  validateJsonString,
				ConflictsWith: []string{"config"},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{string(fc.TRIGGER_TYPE_HTTP), string(fc.TRIGGER_TYPE_LOG),
					string(fc.TRIGGER_TYPE_OSS), string(fc.TRIGGER_TYPE_TIMER), string(fc.TRIGGER_TYPE_MNS_TOPIC)}),
			},

			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudFCTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service").(string)
	fcName := d.Get("function").(string)
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		name = resource.PrefixedUniqueId(v.(string))
	} else {
		name = resource.UniqueId()
	}

	var config interface{}

	if d.Get("type").(string) == string(fc.TRIGGER_TYPE_MNS_TOPIC) {
		if v, ok := d.GetOk("config_mns"); ok {
			if err := json.Unmarshal([]byte(v.(string)), &config); err != nil {
				return WrapError(err)
			}
		}
	} else {
		if v, ok := d.GetOk("config"); ok {
			if err := json.Unmarshal([]byte(v.(string)), &config); err != nil {
				return WrapError(err)
			}
		}
	}

	object := fc.TriggerCreateObject{
		TriggerName:    StringPointer(name),
		TriggerType:    StringPointer(d.Get("type").(string)),
		InvocationRole: StringPointer(d.Get("role").(string)),
		TriggerConfig:  config,
	}
	if v, ok := d.GetOk("source_arn"); ok && v.(string) != "" {
		object.SourceARN = StringPointer(v.(string))
	}
	request := &fc.CreateTriggerInput{
		ServiceName:         StringPointer(serviceName),
		FunctionName:        StringPointer(fcName),
		TriggerCreateObject: object,
	}
	var response *fc.CreateTriggerOutput
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.CreateTrigger(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AccessDenied}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateTrigger", raw)
		response, _ = raw.(*fc.CreateTriggerOutput)
		return nil

	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_trigger", "CreateTrigger", FcGoSdk)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", serviceName, COLON_SEPARATED, fcName, COLON_SEPARATED, *response.TriggerName))

	return resourceAlicloudFCTriggerRead(d, meta)
}

func resourceAlicloudFCTriggerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	trigger, err := fcService.DescribeFcTrigger(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("service", parts[0])
	d.Set("function", parts[1])
	d.Set("name", trigger.TriggerName)
	d.Set("role", trigger.InvocationRole)
	d.Set("source_arn", trigger.SourceARN)

	data, err := trigger.RawTriggerConfig.MarshalJSON()
	if err != nil {
		return WrapError(err)
	}

	if d.Get("type").(string) == string(fc.TRIGGER_TYPE_MNS_TOPIC) {
		if err := d.Set("config_mns", string(data)); err != nil {
			return WrapError(err)
		}
	} else {
		if err := d.Set("config", string(data)); err != nil {
			return WrapError(err)
		}
	}

	d.Set("type", trigger.TriggerType)
	d.Set("last_modified", trigger.LastModifiedTime)

	return nil
}

func resourceAlicloudFCTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	updateInput := &fc.UpdateTriggerInput{}

	if d.HasChange("role") {
		updateInput.InvocationRole = StringPointer(d.Get("role").(string))
	}
	if d.HasChange("config") {
		var config interface{}
		if err := json.Unmarshal([]byte(d.Get("config").(string)), &config); err != nil {
			return WrapError(err)
		}
		updateInput.TriggerConfig = config
	}

	if updateInput != nil {
		split := strings.Split(d.Id(), COLON_SEPARATED)
		if len(split) < 3 {
			return WrapError(Error("Invalid resource ID %s. Please check it and try again.", d.Id()))
		}
		updateInput.ServiceName = StringPointer(split[0])
		updateInput.FunctionName = StringPointer(split[1])
		updateInput.TriggerName = StringPointer(split[2])

		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.UpdateTrigger(updateInput)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateTrigger", FcGoSdk)
		}
		addDebug("UpdateTrigger", raw)
	}

	return resourceAlicloudFCTriggerRead(d, meta)
}

func resourceAlicloudFCTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.DeleteTrigger(&fc.DeleteTriggerInput{
			ServiceName:  StringPointer(parts[0]),
			FunctionName: StringPointer(parts[1]),
			TriggerName:  StringPointer(parts[2]),
		})
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound, TriggerNotFound}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteTrigger", FcGoSdk)
	}
	addDebug("DeleteTrigger", raw)
	return WrapError(fcService.WaitForFcTrigger(d.Id(), Deleted, DefaultTimeoutMedium))
}

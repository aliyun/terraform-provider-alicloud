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

func resourceAlicloudFCFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCFunctionCreate,
		Read:   resourceAlicloudFCFunctionRead,
		Update: resourceAlicloudFCFunctionUpdate,
		Delete: resourceAlicloudFCFunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service": {
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

			"oss_bucket": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"filename"},
			},

			"oss_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"filename"},
			},

			"filename": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"oss_bucket", "oss_key"},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"environment_variables": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"handler": {
				Type:     schema.TypeString,
				Required: true,
			},
			"memory_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      128,
				ValidateFunc: validateIntegerInRange(128, 3072),
			},
			"runtime": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudFCFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service").(string)
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		name = resource.PrefixedUniqueId(v.(string))
	} else {
		name = resource.UniqueId()
	}

	input := &fc.CreateFunctionInput{
		ServiceName: StringPointer(serviceName),
	}
	object := fc.FunctionCreateObject{
		FunctionName: StringPointer(name),
		Description:  StringPointer(d.Get("description").(string)),
		Runtime:      StringPointer(d.Get("runtime").(string)),
		Handler:      StringPointer(d.Get("handler").(string)),
		Timeout:      Int32Pointer(int32(d.Get("timeout").(int))),
		MemorySize:   Int32Pointer(int32(d.Get("memory_size").(int))),
	}
	if variables := d.Get("environment_variables").(map[string]interface{}); len(variables) > 0 {
		byteVar, err := json.Marshal(variables)
		if err != nil {
			return WrapError(err)
		}
		err = json.Unmarshal(byteVar, &object.EnvironmentVariables)
		if err != nil {
			return WrapError(fmt.Errorf("EnvironmentVariables must be type of map[string]string, err is %s", err.Error()))
		}
	}
	code, err := getFunctionCode(d, client)
	if err != nil {
		return WrapError(err)
	}
	object.Code = code
	input.FunctionCreateObject = object

	var function *fc.CreateFunctionOutput
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.CreateFunction(input)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AccessDenied}) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(WrapError(err))
		}
		function, _ = raw.(*fc.CreateFunctionOutput)
		return nil

	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "fc_function", "CreateFunction", AliyunLogGoSdkERROR)
	}

	if function == nil {
		return WrapError(Error("Creating function compute function got a empty response"))
	}

	d.SetId(fmt.Sprintf("%s%s%s", serviceName, COLON_SEPARATED, *function.FunctionName))

	return resourceAlicloudFCFunctionRead(d, meta)
}

func resourceAlicloudFCFunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	if len(split) < 2 {
		return WrapError(fmt.Errorf("Invalid resource ID %s. Please check it and try again.", d.Id()))
	}

	function, err := fcService.DescribeFcFunction(split[0], split[1])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("service", split[0])
	d.Set("name", function.FunctionName)
	d.Set("description", function.Description)
	d.Set("handler", function.Handler)
	d.Set("memory_size", function.MemorySize)
	d.Set("runtime", function.Runtime)
	d.Set("timeout", function.Timeout)
	d.Set("last_modified", function.LastModifiedTime)
	d.Set("environment_variables", function.EnvironmentVariables)

	return nil
}

func resourceAlicloudFCFunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	updateInput := &fc.UpdateFunctionInput{}
	update := false

	if d.HasChange("filename") || d.HasChange("oss_bucket") || d.HasChange("oss_key") {
		update = true
		d.SetPartial("filename")
		d.SetPartial("oss_bucket")
		d.SetPartial("oss_key")
	}
	if d.HasChange("description") {
		updateInput.Description = StringPointer(d.Get("description").(string))
		d.SetPartial("description")
	}
	if d.HasChange("handler") {
		updateInput.Handler = StringPointer(d.Get("handler").(string))
		d.SetPartial("handler")
	}
	if d.HasChange("memory_size") {
		updateInput.MemorySize = Int32Pointer(int32(d.Get("memory_size").(int)))
		d.SetPartial("memory_size")
	}
	if d.HasChange("timeout") {
		updateInput.Timeout = Int32Pointer(int32(d.Get("timeout").(int)))
		d.SetPartial("timeout")
	}
	if d.HasChange("runtime") {
		updateInput.Runtime = StringPointer(d.Get("runtime").(string))
		d.SetPartial("runtime")
	}
	if d.HasChange("environment_variables") {
		byteVar, err := json.Marshal(d.Get("environment_variables").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		err = json.Unmarshal(byteVar, &updateInput.EnvironmentVariables)
		if err != nil {
			return WrapError(fmt.Errorf("EnvironmentVariables must be type of map[string]string, err is %s", err.Error()))
		}
		d.SetPartial("environment_variables")
	}

	if updateInput != nil || update {
		split := strings.Split(d.Id(), COLON_SEPARATED)
		updateInput.ServiceName = StringPointer(split[0])
		updateInput.FunctionName = StringPointer(split[1])
		code, err := getFunctionCode(d, client)
		if err != nil {
			return WrapError(err)
		}
		updateInput.Code = code

		_, err = client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.UpdateFunction(updateInput)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateFunction", AliyunLogGoSdkERROR)
		}
	}

	return resourceAlicloudFCFunctionRead(d, meta)
}

func resourceAlicloudFCFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.DeleteFunction(&fc.DeleteFunctionInput{
				ServiceName:  StringPointer(split[0]),
				FunctionName: StringPointer(split[1]),
			})
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound}) {
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteFunction", AliyunLogGoSdkERROR))
		}

		if _, err := fcService.DescribeFcFunction(split[0], split[1]); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteFunction", AliyunLogGoSdkERROR))
		}
		return nil
	})

}

func getFunctionCode(d *schema.ResourceData, client *connectivity.AliyunClient) (*fc.Code, error) {
	code := fc.NewCode()
	if filename, ok := d.GetOk("filename"); ok && filename.(string) != "" {
		file, err := loadFileContent(filename.(string))
		if err != nil {
			return code, WrapError(fmt.Errorf("Unable to load %q: %s", filename.(string), err))
		}
		code.WithZipFile(file)
	} else {
		bucket, bucketOk := d.GetOk("oss_bucket")
		key, keyOk := d.GetOk("oss_key")
		if !bucketOk || !keyOk {
			return code, WrapError(Error("'oss_bucket' and 'oss_key' must all be set while using OSS code source."))
		}
		code.WithOSSBucketName(bucket.(string)).WithOSSObjectName(key.(string))
	}
	return code, nil
}

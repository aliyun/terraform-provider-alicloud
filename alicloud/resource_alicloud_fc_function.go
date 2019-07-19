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
			"function_id": {
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

	request := &fc.CreateFunctionInput{
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
			return WrapError(err)
		}
	}
	code, err := getFunctionCode(d)
	if err != nil {
		return WrapError(err)
	}
	object.Code = code
	request.FunctionCreateObject = object

	var function *fc.CreateFunctionOutput
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.CreateFunction(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AccessDenied}) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(WrapError(err))
		}
		addDebug("CreateFunction", raw)
		function, _ = raw.(*fc.CreateFunctionOutput)
		return nil

	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_function", "CreateFunction", FcGoSdk)
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

	object, err := fcService.DescribeFcFunction(d.Id())
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
	d.Set("service", parts[0])
	d.Set("name", object.FunctionName)
	d.Set("function_id", object.FunctionID)
	d.Set("description", object.Description)
	d.Set("handler", object.Handler)
	d.Set("memory_size", object.MemorySize)
	d.Set("runtime", object.Runtime)
	d.Set("timeout", object.Timeout)
	d.Set("last_modified", object.LastModifiedTime)
	d.Set("environment_variables", object.EnvironmentVariables)

	return nil
}

func resourceAlicloudFCFunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := &fc.UpdateFunctionInput{}

	if d.HasChange("filename") || d.HasChange("oss_bucket") || d.HasChange("oss_key") {
		d.SetPartial("filename")
		d.SetPartial("oss_bucket")
		d.SetPartial("oss_key")
	}
	if d.HasChange("description") {
		request.Description = StringPointer(d.Get("description").(string))
	}
	if d.HasChange("handler") {
		request.Handler = StringPointer(d.Get("handler").(string))
	}
	if d.HasChange("memory_size") {
		request.MemorySize = Int32Pointer(int32(d.Get("memory_size").(int)))
	}
	if d.HasChange("timeout") {
		request.Timeout = Int32Pointer(int32(d.Get("timeout").(int)))
	}
	if d.HasChange("runtime") {
		request.Runtime = StringPointer(d.Get("runtime").(string))
	}
	if d.HasChange("environment_variables") {
		byteVar, err := json.Marshal(d.Get("environment_variables").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		err = json.Unmarshal(byteVar, &request.EnvironmentVariables)
		if err != nil {
			return WrapError(err)
		}
	}

	if request != nil {
		split := strings.Split(d.Id(), COLON_SEPARATED)
		request.ServiceName = StringPointer(split[0])
		request.FunctionName = StringPointer(split[1])
		code, err := getFunctionCode(d)
		if err != nil {
			return WrapError(err)
		}
		request.Code = code

		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.UpdateFunction(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateFunction", FcGoSdk)
		}
		addDebug("UpdateFunction", raw)
	}

	return resourceAlicloudFCFunctionRead(d, meta)
}

func resourceAlicloudFCFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.DeleteFunction(&fc.DeleteFunctionInput{
			ServiceName:  StringPointer(parts[0]),
			FunctionName: StringPointer(parts[1]),
		})
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteFunction", FcGoSdk)
	}
	addDebug("DeleteFunction", raw)
	return WrapError(fcService.WaitForFcFunction(d.Id(), Deleted, DefaultTimeout))
}

func getFunctionCode(d *schema.ResourceData) (*fc.Code, error) {
	code := fc.NewCode()
	if filename, ok := d.GetOk("filename"); ok && filename.(string) != "" {
		file, err := loadFileContent(filename.(string))
		if err != nil {
			return code, WrapError(err)
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

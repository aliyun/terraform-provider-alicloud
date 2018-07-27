package alicloud

import (
	"fmt"
	"time"

	"strings"

	"github.com/aliyun/fc-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"service": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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

			"oss_bucket": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"filename"},
			},

			"oss_key": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"filename"},
			},

			"filename": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"oss_bucket", "oss_key"},
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"handler": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  true,
			},
			"memory_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      128,
				ValidateFunc: validateIntegerInRange(128, 3072),
			},
			"runtime": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": &schema.Schema{
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
	if err := requireAccountId(meta); err != nil {
		return err
	}
	client := meta.(*AliyunClient)
	conn := client.fcconn

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
	code, err := getFunctionCode(d)
	if err != nil {
		return err
	}
	object.Code = code
	input.FunctionCreateObject = object

	var function *fc.CreateFunctionOutput
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		out, err := conn.CreateFunction(input)
		if err != nil {
			if IsExceptedErrors(err, []string{AccessDenied}) {
				return resource.RetryableError(fmt.Errorf("Error creating function compute service got an error: %#v", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Error creating function compute service got an error: %#v", err))
		}
		function = out
		return nil

	}); err != nil {
		return err
	}

	if function == nil {
		return fmt.Errorf("Creating function compute function got a empty response: %#v.", function)
	}

	d.SetId(fmt.Sprintf("%s%s%s", serviceName, COLON_SEPARATED, *function.FunctionName))

	return resourceAlicloudFCFunctionRead(d, meta)
}

func resourceAlicloudFCFunctionRead(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}

	client := meta.(*AliyunClient)

	split := strings.Split(d.Id(), COLON_SEPARATED)
	if len(split) < 2 {
		return fmt.Errorf("Invalid resource ID %s. Please check it and try again.", d.Id())
	}

	function, err := client.DescribeFcFunction(split[0], split[1])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("DescribeFCFunction %s got an error: %#v", d.Id(), err)
	}

	d.Set("service", split[0])
	d.Set("name", function.FunctionName)
	d.Set("description", function.Description)
	d.Set("handler", function.Handler)
	d.Set("memory_size", function.MemorySize)
	d.Set("runtime", function.Runtime)
	d.Set("timeout", function.Timeout)
	d.Set("last_modified", function.LastModifiedTime)

	return nil
}

func resourceAlicloudFCFunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}
	client := meta.(*AliyunClient)

	d.Partial(true)
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

	if updateInput != nil || update {
		split := strings.Split(d.Id(), COLON_SEPARATED)
		updateInput.ServiceName = StringPointer(split[0])
		updateInput.FunctionName = StringPointer(split[1])
		code, err := getFunctionCode(d)
		if err != nil {
			return err
		}
		updateInput.Code = code

		if _, err := client.fcconn.UpdateFunction(updateInput); err != nil {
			return fmt.Errorf("UpdateFunction %s got an error: %#v.", d.Id(), err)
		}
	}

	d.Partial(false)
	return resourceAlicloudFCFunctionRead(d, meta)
}

func resourceAlicloudFCFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	if err := requireAccountId(meta); err != nil {
		return err
	}
	client := meta.(*AliyunClient)
	split := strings.Split(d.Id(), COLON_SEPARATED)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := client.fcconn.DeleteFunction(&fc.DeleteFunctionInput{
			ServiceName:  StringPointer(split[0]),
			FunctionName: StringPointer(split[1]),
		}); err != nil {
			if IsExceptedErrors(err, []string{ServiceNotFound, FunctionNotFound}) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting function got an error: %#v.", err))
		}

		if _, err := client.DescribeFcFunction(split[0], split[1]); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("While deleting function, getting function %s got an error: %#v.", d.Id(), err))
		}
		return nil
	})

}

func getFunctionCode(d *schema.ResourceData) (*fc.Code, error) {
	code := fc.NewCode()
	if filename, ok := d.GetOk("filename"); ok && filename.(string) != "" {
		file, err := loadFileContent(filename.(string))
		if err != nil {
			return code, fmt.Errorf("Unable to load %q: %s", filename.(string), err)
		}
		code.WithZipFile(file)
	} else {
		bucket, bucketOk := d.GetOk("oss_bucket")
		key, keyOk := d.GetOk("oss_key")
		if !bucketOk || !keyOk {
			return code, fmt.Errorf("'oss_bucket' and 'oss_key' must all be set while using OSS code source.")
		}
		code.WithOSSBucketName(bucket.(string)).WithOSSObjectName(key.(string))
	}
	return code, nil
}

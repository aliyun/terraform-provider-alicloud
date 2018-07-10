package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudDBAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBAccountCreate,
		Read:   resourceAlicloudDBAccountRead,
		Update: resourceAlicloudDBAccountUpdate,
		Delete: resourceAlicloudDBAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(Normal), string(Super)}),
				Default:      "Normal",
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDBAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := rds.CreateCreateAccountRequest()
	request.DBInstanceId = d.Get("instance_id").(string)
	request.AccountName = d.Get("name").(string)
	request.AccountPassword = d.Get("password").(string)
	request.AccountType = d.Get("type").(string)

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.AccountDescription = v.(string)
	}
	// wait instance running before modifying
	if err := client.WaitForDBInstance(request.DBInstanceId, Running, 500); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := request
		if _, err := client.rdsconn.CreateAccount(args); err != nil {
			if IsExceptedError(err, InvalidAccountNameDuplicate) {
				return resource.NonRetryableError(fmt.Errorf("The account %s has already existed. Please import it using ID '%s:%s' or specify a new 'name' and try again.",
					args.AccountName, args.DBInstanceId, args.AccountName))
			} else if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(fmt.Errorf("Create db account got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Create db account got an error: %#v.", err))
		}

		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBInstanceId, COLON_SEPARATED, request.AccountName))

	if err := client.WaitForAccount(request.DBInstanceId, request.AccountName, Available, 500); err != nil {
		return fmt.Errorf("Wait db account %s got an error: %#v.", Available, err)
	}

	return resourceAlicloudDBAccountUpdate(d, meta)
}

func resourceAlicloudDBAccountRead(d *schema.ResourceData, meta interface{}) error {

	parts := strings.Split(d.Id(), COLON_SEPARATED)
	account, err := meta.(*AliyunClient).DescribeDatabaseAccount(parts[0], parts[1])
	if err != nil {
		if NotFoundDBInstance(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Describe db account got an error: %#v", err)
	}

	d.Set("instance_id", account.DBInstanceId)
	d.Set("name", account.AccountName)
	d.Set("type", account.AccountType)
	d.Set("description", account.AccountDescription)

	return nil
}

func resourceAlicloudDBAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if d.HasChange("description") && !d.IsNewResource() {

		request := rds.CreateModifyAccountDescriptionRequest()
		request.DBInstanceId = instanceId
		request.AccountName = accountName
		request.AccountDescription = d.Get("description").(string)

		if _, err := meta.(*AliyunClient).rdsconn.ModifyAccountDescription(request); err != nil {
			return fmt.Errorf("ModifyAccountDescription got an error: %#v", err)
		}
		d.SetPartial("description")
	}

	if d.HasChange("password") && !d.IsNewResource() {

		request := rds.CreateResetAccountPasswordRequest()
		request.DBInstanceId = instanceId
		request.AccountName = accountName
		request.AccountPassword = d.Get("password").(string)

		if _, err := client.rdsconn.ResetAccountPassword(request); err != nil {
			return fmt.Errorf("Error reset db account password error: %#v", err)
		}
		d.SetPartial("password")
	}

	d.Partial(false)
	return resourceAlicloudDBAccountRead(d, meta)
}

func resourceAlicloudDBAccountDelete(d *schema.ResourceData, meta interface{}) error {
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	request := rds.CreateDeleteAccountRequest()
	request.DBInstanceId = parts[0]
	request.AccountName = parts[1]

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := meta.(*AliyunClient).rdsconn.DeleteAccount(request); err != nil {
			if IsExceptedError(err, InvalidAccountNameNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete database account got an error: %#v.", err))
		}

		resp, err := meta.(*AliyunClient).DescribeDatabaseAccount(parts[0], parts[1])
		if err != nil {
			if NotFoundDBInstance(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		} else if resp == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Delete database account got an error: %#v.", err))
	})
}

package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	"time"
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
				ValidateFunc: validateAllowedStringValue([]string{string(rds.Normal), string(rds.Super)}),
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

	args := rds.CreateAccountArgs{
		DBInstanceId:    d.Get("instance_id").(string),
		AccountName:     d.Get("name").(string),
		AccountPassword: d.Get("password").(string),
		AccountType:     rds.AccountType(d.Get("type").(string)),
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.AccountDescription = v.(string)
	}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ag := args
		if _, err := meta.(*AliyunClient).rdsconn.CreateAccount(&ag); err != nil {
			if IsExceptedError(err, InvalidAccountNameDuplicate) {
				return resource.NonRetryableError(fmt.Errorf("The account %s has already existed. Please import it using ID '%s:%s' or specify a new 'name' and try again.",
					args.AccountName, args.DBInstanceId, args.AccountName))
			} else if IsExceptedError(err, OperationDeniedDBInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("Create db account got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Create db account got an error: %#v.", err))
		}

		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", args.DBInstanceId, COLON_SEPARATED, args.AccountName))

	if err := meta.(*AliyunClient).rdsconn.WaitForAccount(args.DBInstanceId, args.AccountName, rds.Available, defaultTimeout); err != nil {
		return fmt.Errorf("Wait db account %s got an error: %#v.", rds.Available, err)
	}

	return resourceAlicloudDBAccountUpdate(d, meta)
}

func resourceAlicloudDBAccountRead(d *schema.ResourceData, meta interface{}) error {

	parts := strings.Split(d.Id(), COLON_SEPARATED)
	account, err := meta.(*AliyunClient).DescribeDatabaseAccount(parts[0], parts[1])
	if err != nil {
		if NotFoundError(err) {
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

		if err := meta.(*AliyunClient).rdsconn.ModifyAccountDescription(&rds.ModifyAccountDescriptionArgs{
			DBInstanceId:       instanceId,
			AccountName:        accountName,
			AccountDescription: d.Get("description").(string),
		}); err != nil {
			return fmt.Errorf("ModifyAccountDescription got an error: %#v", err)
		}
		d.SetPartial("description")
	}

	if d.HasChange("password") && !d.IsNewResource() {
		if _, err := client.rdsconn.ResetAccountPassword(instanceId, accountName, d.Get("password").(string)); err != nil {
			return fmt.Errorf("Error reset db account password error: %#v", err)
		}
		d.SetPartial("password")
	}

	d.Partial(false)
	return resourceAlicloudDBAccountRead(d, meta)
}

func resourceAlicloudDBAccountDelete(d *schema.ResourceData, meta interface{}) error {
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := meta.(*AliyunClient).rdsconn.DeleteAccount(parts[0], parts[1]); err != nil {
			if IsExceptedError(err, InvalidAccountNameNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete database account got an error: %#v.", err))
		}

		resp, err := meta.(*AliyunClient).DescribeDatabaseAccount(parts[0], parts[1])
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidAccountNameNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		} else if resp == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Delete database account got an error: %#v.", err))
	})
}

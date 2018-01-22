package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudDBDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBDatabaseCreate,
		Read:   resourceAlicloudDBDatabaseRead,
		Update: resourceAlicloudDBDatabaseUpdate,
		Delete: resourceAlicloudDBDatabaseDelete,
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

			"character_set": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(rds.CHARACTER_SET_NAME),
				Optional:     true,
				Default:      "utf8",
				ForceNew:     true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDBDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	args := rds.CreateDatabaseArgs{
		DBInstanceId:     d.Get("instance_id").(string),
		DBName:           d.Get("name").(string),
		CharacterSetName: d.Get("character_set").(string),
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.DBDescription = v.(string)
	}

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ag := args
		if _, err := meta.(*AliyunClient).rdsconn.CreateDatabase(&ag); err != nil {
			if IsExceptedError(err, OperationDeniedDBInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("Create database got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Create database got an error: %#v.", err))
		}

		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", args.DBInstanceId, COLON_SEPARATED, args.DBName))

	return resourceAlicloudDBDatabaseUpdate(d, meta)
}

func resourceAlicloudDBDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	db, err := meta.(*AliyunClient).DescribeDatabaseByName(parts[0], parts[1])
	if err != nil {
		return fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err)
	}

	if db == nil {
		d.SetId("")
		return nil
	}

	d.Set("instance_id", db.DBInstanceId)
	d.Set("name", db.DBName)
	d.Set("character_set", db.CharacterSetName)
	d.Set("description", db.DBDescription)

	return nil
}

func resourceAlicloudDBDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	if d.HasChange("description") && !d.IsNewResource() {
		parts := strings.Split(d.Id(), COLON_SEPARATED)
		if err := meta.(*AliyunClient).rdsconn.ModifyDatabaseDescription(&rds.ModifyDatabaseDescriptionArgs{
			DBInstanceId:  parts[0],
			DBName:        parts[1],
			DBDescription: d.Get("description").(string),
		}); err != nil {
			return fmt.Errorf("ModifyDatabaseDescription got an error: %#v", err)
		}
		d.SetPartial("description")
	}

	d.Partial(false)
	return resourceAlicloudDBDatabaseRead(d, meta)
}

func resourceAlicloudDBDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).rdsconn
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DeleteDatabase(parts[0], parts[1])

		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidDBNameNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete database %s timeout and got an error: %#v.", parts[1], err))
		}

		db, err := meta.(*AliyunClient).DescribeDatabaseByName(parts[0], parts[1])
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err))
		}
		if db == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Delete database %s timeout and got an error: %#v.", parts[1], err))
	})
}

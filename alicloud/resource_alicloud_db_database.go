package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
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
				ValidateFunc: validateAllowedStringValue(CHARACTER_SET_NAME),
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

	client := meta.(*AliyunClient)
	request := rds.CreateCreateDatabaseRequest()
	request.DBInstanceId = d.Get("instance_id").(string)
	request.DBName = d.Get("name").(string)
	request.CharacterSetName = d.Get("character_set").(string)

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.DBDescription = v.(string)
	}

	if inst, err := client.DescribeDBInstanceById(request.DBInstanceId); err != nil {
		return fmt.Errorf("DescribeDBInstance got an error: %#v", err)
	} else if inst.Engine == string(PostgreSQL) || inst.Engine == string(PPAS) {
		return fmt.Errorf("At present, it does not support creating 'PostgreSQL' and 'PPAS' database. Please login DB instance to create.")
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		ag := request
		if _, err := client.rdsconn.CreateDatabase(ag); err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(fmt.Errorf("Create database got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Create database got an error: %#v.", err))
		}

		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBInstanceId, COLON_SEPARATED, request.DBName))

	return resourceAlicloudDBDatabaseUpdate(d, meta)
}

func resourceAlicloudDBDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	db, err := meta.(*AliyunClient).DescribeDatabaseByName(parts[0], parts[1])
	if err != nil {
		if NotFoundDBInstance(err) || IsExceptedError(err, InvalidDBNameNotFound) {
			d.SetId("")
			return nil
		}
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
		request := rds.CreateModifyDBDescriptionRequest()
		request.DBInstanceId = parts[0]
		request.DBName = parts[1]
		request.DBDescription = d.Get("description").(string)

		if _, err := meta.(*AliyunClient).rdsconn.ModifyDBDescription(request); err != nil {
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
	request := rds.CreateDeleteDatabaseRequest()
	request.DBInstanceId = parts[0]
	request.DBName = parts[1]
	// wait instance status is running before deleting database
	if err := meta.(*AliyunClient).WaitForDBInstance(parts[0], Running, 1800); err != nil {
		return fmt.Errorf("While deleting database, WaitForInstance %s got error: %#v", Running, err)
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DeleteDatabase(request); err != nil {
			if NotFoundDBInstance(err) || IsExceptedError(err, InvalidDBNameNotFound) {
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

package alicloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudPolarDBDynamoTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBDynamoTableCreate,
		Read:   resourceAlicloudPolarDBDynamoTableRead,
		Update: resourceAlicloudPolarDBDynamoTableUpdate,
		Delete: resourceAlicloudPolarDBDynamoTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		// billing_mode defaults to PROVISIONED while the capacities have no default,
		// so validate them at plan time instead of failing CreateTable with a
		// server-side validation error on nil throughput units.
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			if diff.Get("billing_mode").(string) != "PROVISIONED" {
				return nil
			}
			if diff.Get("read_capacity").(int) <= 0 || diff.Get("write_capacity").(int) <= 0 {
				return fmt.Errorf("read_capacity and write_capacity are required and must be greater than 0 when billing_mode is PROVISIONED")
			}
			if v, ok := diff.GetOk("global_secondary_index"); ok {
				for _, gsiRaw := range v.(*schema.Set).List() {
					gsi := gsiRaw.(map[string]interface{})
					if gsi["read_capacity"].(int) <= 0 || gsi["write_capacity"].(int) <= 0 {
						return fmt.Errorf("global_secondary_index %q requires read_capacity and write_capacity greater than 0 when billing_mode is PROVISIONED", gsi["name"])
					}
				}
			}
			return nil
		},
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The PolarDB DynamoDB-compatible endpoint URL.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The account name for PolarDB DynamoDB authentication. If not set, it is resolved from the cluster's DynamoDB-type account automatically.",
			},
			"account_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The authentication password for PolarDB DynamoDB. If not set, it is resolved from the cluster's DynamoDB-type account automatically.",
			},
			"db_cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The ID of the PolarDB cluster where DynamoDB is enabled.",
			},
			"table_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "The name of the DynamoDB-compatible table.",
			},
			"attribute": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "List of attribute definitions for the table key schema and indexes.",
				// Attribute names are unique within a table, so hash the set on the
				// name only. This keeps membership deterministic regardless of the
				// declaration order and makes a type change show up as an in-place
				// update to the same element instead of a remove+add churn.
				Set: func(v interface{}) int {
					return schema.HashString(v.(map[string]interface{})["name"].(string))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"S", "N", "B"}, false),
						},
					},
				},
			},
			"hash_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"range_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"billing_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PROVISIONED",
				ValidateFunc: validation.StringInSlice([]string{"PROVISIONED", "PAY_PER_REQUEST"}, false),
			},
			"read_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"write_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"global_secondary_index": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				// Index names are unique, so hash the set on the name only for
				// stable, order-independent set membership.
				Set: func(v interface{}) int {
					return schema.HashString(v.(map[string]interface{})["name"].(string))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":               {Type: schema.TypeString, Required: true},
						"hash_key":           {Type: schema.TypeString, Optional: true, Computed: true},
						"range_key":          {Type: schema.TypeString, Optional: true, Computed: true},
						"projection_type":    {Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice([]string{"ALL", "KEYS_ONLY", "INCLUDE"}, false)},
						"non_key_attributes": {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"read_capacity":      {Type: schema.TypeInt, Optional: true, Computed: true},
						"write_capacity":     {Type: schema.TypeInt, Optional: true, Computed: true},
					},
				},
			},
			"local_secondary_index": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				// Index names are unique, so hash the set on the name only for
				// stable, order-independent set membership.
				Set: func(v interface{}) int {
					return schema.HashString(v.(map[string]interface{})["name"].(string))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":               {Type: schema.TypeString, Required: true, ForceNew: true},
						"range_key":          {Type: schema.TypeString, Required: true, ForceNew: true},
						"projection_type":    {Type: schema.TypeString, Required: true, ForceNew: true, ValidateFunc: validation.StringInSlice([]string{"ALL", "KEYS_ONLY", "INCLUDE"}, false)},
						"non_key_attributes": {Type: schema.TypeList, Optional: true, ForceNew: true, Elem: &schema.Schema{Type: schema.TypeString}},
					},
				},
			},
			"ttl": {
				Type: schema.TypeList, Optional: true, Computed: true, MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled":        {Type: schema.TypeBool, Optional: true, Default: false},
						"attribute_name": {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"arn": {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceAlicloudPolarDBDynamoTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tableName := d.Get("table_name").(string)

	conn, err := resolvePolarDBDynamoConn(d, meta, d.Get("db_cluster_id").(string))
	if err != nil {
		return err
	}
	dynamoClient := conn.client

	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
	}

	keySchema := make([]types.KeySchemaElement, 0)
	if v, ok := d.GetOk("hash_key"); ok && v.(string) != "" {
		keySchema = append(keySchema, types.KeySchemaElement{AttributeName: aws.String(v.(string)), KeyType: types.KeyTypeHash})
	}
	if v, ok := d.GetOk("range_key"); ok && v.(string) != "" {
		keySchema = append(keySchema, types.KeySchemaElement{AttributeName: aws.String(v.(string)), KeyType: types.KeyTypeRange})
	}
	if len(keySchema) > 0 {
		input.KeySchema = keySchema
	}

	if v, ok := d.GetOk("attribute"); ok {
		attrs := make([]types.AttributeDefinition, 0)
		for _, attrRaw := range v.(*schema.Set).List() {
			attr := attrRaw.(map[string]interface{})
			attrs = append(attrs, types.AttributeDefinition{
				AttributeName: aws.String(attr["name"].(string)),
				AttributeType: types.ScalarAttributeType(attr["type"].(string)),
			})
		}
		input.AttributeDefinitions = attrs
	}

	if v, ok := d.GetOk("billing_mode"); ok {
		input.BillingMode = types.BillingMode(v.(string))
	}

	if d.Get("billing_mode").(string) == "PROVISIONED" {
		pt := &types.ProvisionedThroughput{}
		if v, ok := d.GetOk("read_capacity"); ok {
			pt.ReadCapacityUnits = aws.Int64(int64(v.(int)))
		}
		if v, ok := d.GetOk("write_capacity"); ok {
			pt.WriteCapacityUnits = aws.Int64(int64(v.(int)))
		}
		input.ProvisionedThroughput = pt
	}

	if v, ok := d.GetOk("global_secondary_index"); ok {
		gsis := make([]types.GlobalSecondaryIndex, 0)
		for _, gsiRaw := range v.(*schema.Set).List() {
			gsi := gsiRaw.(map[string]interface{})
			gsiConfig := types.GlobalSecondaryIndex{IndexName: aws.String(gsi["name"].(string))}
			gsiKs := make([]types.KeySchemaElement, 0)
			if hk, ok := gsi["hash_key"]; ok && hk.(string) != "" {
				gsiKs = append(gsiKs, types.KeySchemaElement{AttributeName: aws.String(hk.(string)), KeyType: types.KeyTypeHash})
			}
			if rk, ok := gsi["range_key"]; ok && rk.(string) != "" {
				gsiKs = append(gsiKs, types.KeySchemaElement{AttributeName: aws.String(rk.(string)), KeyType: types.KeyTypeRange})
			}
			if len(gsiKs) > 0 {
				gsiConfig.KeySchema = gsiKs
			}
			proj := &types.Projection{ProjectionType: types.ProjectionType(gsi["projection_type"].(string))}
			if nka, ok := gsi["non_key_attributes"]; ok {
				nkaList := expandStringList(nka.(*schema.Set).List())
				if len(nkaList) > 0 {
					proj.NonKeyAttributes = nkaList
				}
			}
			gsiConfig.Projection = proj
			if d.Get("billing_mode").(string) == "PROVISIONED" {
				gsiPt := &types.ProvisionedThroughput{}
				if rc, ok := gsi["read_capacity"]; ok && rc.(int) > 0 {
					gsiPt.ReadCapacityUnits = aws.Int64(int64(rc.(int)))
				}
				if wc, ok := gsi["write_capacity"]; ok && wc.(int) > 0 {
					gsiPt.WriteCapacityUnits = aws.Int64(int64(wc.(int)))
				}
				gsiConfig.ProvisionedThroughput = gsiPt
			}
			gsis = append(gsis, gsiConfig)
		}
		if len(gsis) > 0 {
			input.GlobalSecondaryIndexes = gsis
		}
	}

	if v, ok := d.GetOk("local_secondary_index"); ok {
		lsis := make([]types.LocalSecondaryIndex, 0)
		for _, lsiRaw := range v.(*schema.Set).List() {
			lsi := lsiRaw.(map[string]interface{})
			lsiConfig := types.LocalSecondaryIndex{IndexName: aws.String(lsi["name"].(string))}
			lsiKs := make([]types.KeySchemaElement, 0)
			if hk, ok := d.GetOk("hash_key"); ok && hk.(string) != "" {
				lsiKs = append(lsiKs, types.KeySchemaElement{AttributeName: aws.String(hk.(string)), KeyType: types.KeyTypeHash})
			}
			if rk, ok := lsi["range_key"]; ok && rk.(string) != "" {
				lsiKs = append(lsiKs, types.KeySchemaElement{AttributeName: aws.String(rk.(string)), KeyType: types.KeyTypeRange})
			}
			if len(lsiKs) > 0 {
				lsiConfig.KeySchema = lsiKs
			}
			proj := &types.Projection{ProjectionType: types.ProjectionType(lsi["projection_type"].(string))}
			if nka, ok := lsi["non_key_attributes"]; ok {
				nkaList := expandStringList(nka.([]interface{}))
				if len(nkaList) > 0 {
					proj.NonKeyAttributes = nkaList
				}
			}
			lsiConfig.Projection = proj
			lsis = append(lsis, lsiConfig)
		}
		if len(lsis) > 0 {
			input.LocalSecondaryIndexes = lsis
		}
	}

	var output *dynamodb.CreateTableOutput
	// DNS registration of a freshly created public endpoint address can take well
	// over ten minutes, so drive the retry loop by the resource's create timeout
	// instead of a short hard-coded window.
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		output, err = dynamoClient.CreateTable(context.Background(), input)
		if err != nil {
			if isDynamoRetryableError(err) {
				log.Printf("[DEBUG] CreateTable %s failed with retryable error, will retry: %s", tableName, err)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug("CreateTable", output, nil, input)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_dynamo_table", "CreateTable", AlibabaCloudSdkGoERROR)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	d.SetId(fmt.Sprintf("%s%s%s", dbClusterId, COLON_SEPARATED, tableName))

	polarDBService := PolarDBService{client}
	stateConf := BuildStateConf(
		[]string{"CREATING", "UPDATING"},
		[]string{"ACTIVE"},
		d.Timeout(schema.TimeoutCreate),
		5*time.Second,
		polarDBService.PolarDBDynamoTableStateRefreshFunc(dbClusterId, tableName, conn.endpoint, conn.accessKey, conn.secretKey),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	// TTL cannot be specified on CreateTable, so apply it once the table is ACTIVE;
	// otherwise the first apply would silently leave it disabled and only converge
	// on the next apply.
	if ttlList := d.Get("ttl").([]interface{}); len(ttlList) > 0 && ttlList[0] != nil {
		ttlConfig := ttlList[0].(map[string]interface{})
		if ttlConfig["enabled"].(bool) {
			ttlSpec := &types.TimeToLiveSpecification{Enabled: aws.Bool(true)}
			if attrName, ok := ttlConfig["attribute_name"]; ok && attrName.(string) != "" {
				ttlSpec.AttributeName = aws.String(attrName.(string))
			}
			_, err = dynamoClient.UpdateTimeToLive(context.Background(), &dynamodb.UpdateTimeToLiveInput{
				TableName:               aws.String(tableName),
				TimeToLiveSpecification: ttlSpec,
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateTimeToLive", AlibabaCloudSdkGoERROR)
			}
		}
	}

	return resourceAlicloudPolarDBDynamoTableRead(d, meta)
}

func resourceAlicloudPolarDBDynamoTableRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	tableName := parts[1]

	conn, err := resolvePolarDBDynamoConn(d, meta, dbClusterId)
	if err != nil {
		return err
	}
	dynamoClient := conn.client

	var output *dynamodb.DescribeTableOutput
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		output, err = dynamoClient.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
		if err != nil {
			if isDynamoNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			if isDynamoRetryableError(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug("DescribeTable", output, nil, nil)
		return nil
	})
	if err != nil {
		if isDynamoNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DescribeTable", AlibabaCloudSdkGoERROR)
	}

	if output.Table == nil {
		return fmt.Errorf("missing Table field in DescribeTable response for %s", d.Id())
	}

	tableDesc := output.Table
	d.Set("db_cluster_id", dbClusterId)
	d.Set("table_name", tableName)

	for _, ks := range tableDesc.KeySchema {
		if ks.KeyType == types.KeyTypeHash {
			d.Set("hash_key", aws.ToString(ks.AttributeName))
		} else if ks.KeyType == types.KeyTypeRange {
			d.Set("range_key", aws.ToString(ks.AttributeName))
		}
	}

	attrSet := make([]map[string]interface{}, 0)
	for _, attr := range tableDesc.AttributeDefinitions {
		attrSet = append(attrSet, map[string]interface{}{
			"name": aws.ToString(attr.AttributeName),
			"type": string(attr.AttributeType),
		})
	}
	d.Set("attribute", attrSet)

	if tableDesc.BillingModeSummary != nil {
		d.Set("billing_mode", string(tableDesc.BillingModeSummary.BillingMode))
	}
	if tableDesc.ProvisionedThroughput != nil {
		d.Set("read_capacity", int(aws.ToInt64(tableDesc.ProvisionedThroughput.ReadCapacityUnits)))
		d.Set("write_capacity", int(aws.ToInt64(tableDesc.ProvisionedThroughput.WriteCapacityUnits)))
	}

	gsiSet := make([]map[string]interface{}, 0)
	// Build a lookup of existing GSI configs from state so we can preserve
	// read_capacity / write_capacity when the API omits ProvisionedThroughput
	// (e.g. PAY_PER_REQUEST billing mode). Without this, helper/schema fills
	// absent keys with zero values and the resource never converges.
	existingGSIs := make(map[string]map[string]interface{})
	if raw, ok := d.GetOk("global_secondary_index"); ok {
		for _, g := range raw.(*schema.Set).List() {
			if gm, ok := g.(map[string]interface{}); ok {
				if name, ok := gm["name"].(string); ok {
					existingGSIs[name] = gm
				}
			}
		}
	}
	for _, gsi := range tableDesc.GlobalSecondaryIndexes {
		gsiName := aws.ToString(gsi.IndexName)
		gsiConfig := map[string]interface{}{"name": gsiName}
		for _, ks := range gsi.KeySchema {
			if ks.KeyType == types.KeyTypeHash {
				gsiConfig["hash_key"] = aws.ToString(ks.AttributeName)
			} else if ks.KeyType == types.KeyTypeRange {
				gsiConfig["range_key"] = aws.ToString(ks.AttributeName)
			}
		}
		if gsi.Projection != nil {
			gsiConfig["projection_type"] = string(gsi.Projection.ProjectionType)
			if len(gsi.Projection.NonKeyAttributes) > 0 {
				gsiConfig["non_key_attributes"] = gsi.Projection.NonKeyAttributes
			}
		}
		if gsi.ProvisionedThroughput != nil {
			gsiConfig["read_capacity"] = int(aws.ToInt64(gsi.ProvisionedThroughput.ReadCapacityUnits))
			gsiConfig["write_capacity"] = int(aws.ToInt64(gsi.ProvisionedThroughput.WriteCapacityUnits))
		} else if prev, ok := existingGSIs[gsiName]; ok {
			// Preserve capacities from state when the API omits them
			if rc, ok := prev["read_capacity"]; ok {
				gsiConfig["read_capacity"] = rc
			}
			if wc, ok := prev["write_capacity"]; ok {
				gsiConfig["write_capacity"] = wc
			}
		}
		gsiSet = append(gsiSet, gsiConfig)
	}
	d.Set("global_secondary_index", gsiSet)

	lsiSet := make([]map[string]interface{}, 0)
	for _, lsi := range tableDesc.LocalSecondaryIndexes {
		lsiConfig := map[string]interface{}{"name": aws.ToString(lsi.IndexName)}
		for _, ks := range lsi.KeySchema {
			if ks.KeyType == types.KeyTypeRange {
				lsiConfig["range_key"] = aws.ToString(ks.AttributeName)
			}
		}
		if lsi.Projection != nil {
			lsiConfig["projection_type"] = string(lsi.Projection.ProjectionType)
			if len(lsi.Projection.NonKeyAttributes) > 0 {
				lsiConfig["non_key_attributes"] = lsi.Projection.NonKeyAttributes
			}
		}
		lsiSet = append(lsiSet, lsiConfig)
	}
	d.Set("local_secondary_index", lsiSet)

	d.Set("arn", aws.ToString(tableDesc.TableArn))

	ttlOutput, ttlErr := dynamoClient.DescribeTimeToLive(context.Background(), &dynamodb.DescribeTimeToLiveInput{TableName: aws.String(tableName)})
	if ttlErr == nil && ttlOutput.TimeToLiveDescription != nil {
		ttlConfig := map[string]interface{}{"enabled": false}
		if ttlOutput.TimeToLiveDescription.TimeToLiveStatus == types.TimeToLiveStatusEnabled {
			ttlConfig["enabled"] = true
		}
		if attrName := aws.ToString(ttlOutput.TimeToLiveDescription.AttributeName); attrName != "" {
			ttlConfig["attribute_name"] = attrName
		}
		d.Set("ttl", []interface{}{ttlConfig})
	}

	return nil
}

func resourceAlicloudPolarDBDynamoTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	tableName := parts[1]

	conn, err := resolvePolarDBDynamoConn(d, meta, dbClusterId)
	if err != nil {
		return err
	}
	dynamoClient := conn.client

	if d.HasChanges("read_capacity", "write_capacity") && d.Get("billing_mode").(string) == "PROVISIONED" {
		updateInput := &dynamodb.UpdateTableInput{TableName: aws.String(tableName)}
		pt := &types.ProvisionedThroughput{}
		if v, ok := d.GetOk("read_capacity"); ok {
			pt.ReadCapacityUnits = aws.Int64(int64(v.(int)))
		}
		if v, ok := d.GetOk("write_capacity"); ok {
			pt.WriteCapacityUnits = aws.Int64(int64(v.(int)))
		}
		updateInput.ProvisionedThroughput = pt

		err = retry.Retry(8*time.Minute, func() *retry.RetryError {
			_, err = dynamoClient.UpdateTable(context.Background(), updateInput)
			if err != nil {
				if isDynamoRetryableError(err) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateTable", AlibabaCloudSdkGoERROR)
		}

		polarDBService := PolarDBService{client}
		stateConf := BuildStateConf(
			[]string{"UPDATING"}, []string{"ACTIVE"},
			d.Timeout(schema.TimeoutUpdate), 5*time.Second,
			polarDBService.PolarDBDynamoTableStateRefreshFunc(dbClusterId, tableName, conn.endpoint, conn.accessKey, conn.secretKey),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("global_secondary_index") {
		oldGSI, newGSI := d.GetChange("global_secondary_index")
		billingMode := d.Get("billing_mode").(string)
		gsiUpdates := updateDiffDynamoGSI(oldGSI.(*schema.Set).List(), newGSI.(*schema.Set).List(), billingMode)

		// Phase 1: deletions first, one at a time (only one online index operation
		// is allowed per table at any moment).
		for _, gsiUpdate := range gsiUpdates {
			if gsiUpdate.Delete == nil {
				continue
			}
			idxName := aws.ToString(gsiUpdate.Delete.IndexName)
			_, err = dynamoClient.UpdateTable(context.Background(), &dynamodb.UpdateTableInput{
				TableName:                   aws.String(tableName),
				GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{gsiUpdate},
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteGSI:"+idxName, AlibabaCloudSdkGoERROR)
			}
			if err := waitDynamoGSIState(dynamoClient, tableName, idxName, false, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return WrapError(err)
			}
		}

		// Phase 2: capacity-only updates can be sent in a single request.
		capacityUpdates := make([]types.GlobalSecondaryIndexUpdate, 0)
		for _, gsiUpdate := range gsiUpdates {
			if gsiUpdate.Update != nil {
				capacityUpdates = append(capacityUpdates, gsiUpdate)
			}
		}
		if len(capacityUpdates) > 0 {
			_, err = dynamoClient.UpdateTable(context.Background(), &dynamodb.UpdateTableInput{
				TableName:                   aws.String(tableName),
				GlobalSecondaryIndexUpdates: capacityUpdates,
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateGSI", AlibabaCloudSdkGoERROR)
			}
			for _, gsiUpdate := range capacityUpdates {
				if err := waitDynamoGSIState(dynamoClient, tableName, aws.ToString(gsiUpdate.Update.IndexName), true, d.Timeout(schema.TimeoutUpdate)); err != nil {
					return WrapError(err)
				}
			}
		}

		// Phase 3: creations last, one at a time; the request must carry only the
		// attribute definitions referenced by the new index's KeySchema (the API
		// rejects unused attributes with a ValidationException).
		for _, gsiUpdate := range gsiUpdates {
			if gsiUpdate.Create == nil {
				continue
			}
			idxName := aws.ToString(gsiUpdate.Create.IndexName)
			// Collect attribute names used in this GSI's KeySchema.
			usedKeys := make(map[string]bool)
			for _, ks := range gsiUpdate.Create.KeySchema {
				usedKeys[aws.ToString(ks.AttributeName)] = true
			}
			attrs := make([]types.AttributeDefinition, 0)
			for _, attrRaw := range d.Get("attribute").(*schema.Set).List() {
				attr := attrRaw.(map[string]interface{})
				if !usedKeys[attr["name"].(string)] {
					continue
				}
				attrs = append(attrs, types.AttributeDefinition{
					AttributeName: aws.String(attr["name"].(string)),
					AttributeType: types.ScalarAttributeType(attr["type"].(string)),
				})
			}
			_, err = dynamoClient.UpdateTable(context.Background(), &dynamodb.UpdateTableInput{
				TableName:                   aws.String(tableName),
				AttributeDefinitions:        attrs,
				GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{gsiUpdate},
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "CreateGSI:"+idxName, AlibabaCloudSdkGoERROR)
			}
			if err := waitDynamoGSIState(dynamoClient, tableName, idxName, true, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return WrapError(err)
			}
		}
	}

	if d.HasChange("ttl") {
		ttlList := d.Get("ttl").([]interface{})
		if len(ttlList) > 0 && ttlList[0] != nil {
			ttlConfig := ttlList[0].(map[string]interface{})
			ttlSpec := &types.TimeToLiveSpecification{Enabled: aws.Bool(ttlConfig["enabled"].(bool))}
			if attrName, ok := ttlConfig["attribute_name"]; ok && attrName.(string) != "" {
				ttlSpec.AttributeName = aws.String(attrName.(string))
			}
			_, err = dynamoClient.UpdateTimeToLive(context.Background(), &dynamodb.UpdateTimeToLiveInput{
				TableName:               aws.String(tableName),
				TimeToLiveSpecification: ttlSpec,
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateTimeToLive", AlibabaCloudSdkGoERROR)
			}
		}
	}

	return resourceAlicloudPolarDBDynamoTableRead(d, meta)
}

func resourceAlicloudPolarDBDynamoTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	tableName := parts[1]

	conn, err := resolvePolarDBDynamoConn(d, meta, dbClusterId)
	if err != nil {
		return err
	}
	dynamoClient := conn.client

	err = retry.Retry(8*time.Minute, func() *retry.RetryError {
		_, err = dynamoClient.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
		if err != nil {
			if isDynamoNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			if isDynamoRetryableError(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug("DeleteTable", nil, nil, nil)
		return nil
	})
	if err != nil {
		if isDynamoNotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteTable", AlibabaCloudSdkGoERROR)
	}

	polarDBService := PolarDBService{client}
	stateConf := BuildStateConf(
		[]string{"ACTIVE", "DELETING"}, []string{},
		d.Timeout(schema.TimeoutDelete), 5*time.Second,
		polarDBService.PolarDBDynamoTableStateRefreshFunc(dbClusterId, tableName, conn.endpoint, conn.accessKey, conn.secretKey),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DeleteTimeoutMsg, d.Id(), "DeleteTable", ProviderERROR)
	}

	return nil
}

func isDynamoNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "ResourceNotFoundException", "TableNotFound":
			return true
		}
		return false
	}
	// Errors from the PolarDB OpenAPI (not the DynamoDB-compatible endpoint) do
	// not implement smithy.APIError, so keep a narrow string match for them.
	return strings.Contains(err.Error(), "InvalidDBClusterId.NotFound")
}

func isDynamoRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Service-side errors carry an explicit code; only transient ones are worth
	// retrying. Everything else (validation, auth, ...) must fail fast instead of
	// hanging in the retry window and masking the real cause.
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "InternalServerError", "ThrottlingException", "LimitExceededException", "ProvisionedThroughputExceededException":
			return true
		}
		return false
	}

	// The request never got a service response: the public endpoint address/domain
	// returned by alicloud_polardb_endpoint_address may not be resolvable or
	// reachable immediately after creation (DNS propagation delay or the listener
	// not fully up yet), so transport-level errors (DNS failures, timeouts,
	// refused/reset connections) are retryable.
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)
}

// polarDBDynamoConn bundles a ready DynamoDB-compatible client together with the
// resolved connection info, which some callers still need (e.g. for the
// table state refresh function).
type polarDBDynamoConn struct {
	client    *dynamodb.Client
	endpoint  string
	accessKey string
	secretKey string
}

// resolvePolarDBDynamoConn resolves the DynamoDB-compatible endpoint address and
// the DynamoDB-type account credentials for the given cluster, falling back to
// the cluster's own configuration when they are absent from state (e.g. on
// import; the endpoint rejects the provider AK/SK). Resolved values are
// persisted into state so follow-up operations can reuse them.
func resolvePolarDBDynamoConn(d *schema.ResourceData, meta interface{}, dbClusterId string) (*polarDBDynamoConn, error) {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	endpoint := d.Get("endpoint").(string)
	if endpoint == "" {
		addr, err := polarDBService.DescribePolarDBDynamoEndpointAddress(dbClusterId)
		if err != nil {
			return nil, WrapError(err)
		}
		endpoint = addr
		d.Set("endpoint", endpoint)
	}

	accessKey := d.Get("account_name").(string)
	secretKey := d.Get("account_auth").(string)
	if accessKey == "" || secretKey == "" {
		ak, sk, err := polarDBService.DescribePolarDBDynamoAccount(dbClusterId)
		if err != nil {
			return nil, WrapError(err)
		}
		accessKey, secretKey = ak, sk
		d.Set("account_name", accessKey)
		d.Set("account_auth", secretKey)
	}

	dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accessKey, secretKey)
	if err != nil {
		return nil, WrapError(err)
	}
	return &polarDBDynamoConn{client: dynamoClient, endpoint: endpoint, accessKey: accessKey, secretKey: secretKey}, nil
}

// expandDynamoGSIKeySchema builds the key schema of a GSI from its config map.
func expandDynamoGSIKeySchema(gsi map[string]interface{}) []types.KeySchemaElement {
	ks := make([]types.KeySchemaElement, 0, 2)
	if hk, ok := gsi["hash_key"]; ok && hk.(string) != "" {
		ks = append(ks, types.KeySchemaElement{AttributeName: aws.String(hk.(string)), KeyType: types.KeyTypeHash})
	}
	if rk, ok := gsi["range_key"]; ok && rk.(string) != "" {
		ks = append(ks, types.KeySchemaElement{AttributeName: aws.String(rk.(string)), KeyType: types.KeyTypeRange})
	}
	return ks
}

// expandDynamoGSIProjection builds the projection of a GSI from its config map.
func expandDynamoGSIProjection(gsi map[string]interface{}) *types.Projection {
	proj := &types.Projection{ProjectionType: types.ProjectionType(gsi["projection_type"].(string))}
	if nka, ok := gsi["non_key_attributes"]; ok {
		nkaList := expandStringList(nka.(*schema.Set).List())
		if len(nkaList) > 0 {
			proj.NonKeyAttributes = nkaList
		}
	}
	return proj
}

// updateDiffDynamoGSI computes the GSI update operations between the old and new
// configurations, following the same strategy as the AWS provider: indexes only
// present in the new config are created, indexes only present in the old config
// are deleted, same-named indexes with key schema/projection changes are recreated
// (delete + create), and same-named indexes with capacity-only changes are updated.
func updateDiffDynamoGSI(oldGsi, newGsi []interface{}, billingMode string) []types.GlobalSecondaryIndexUpdate {
	oldGsis := make(map[string]map[string]interface{})
	for _, raw := range oldGsi {
		m := raw.(map[string]interface{})
		oldGsis[m["name"].(string)] = m
	}
	newGsis := make(map[string]map[string]interface{})
	for _, raw := range newGsi {
		m := raw.(map[string]interface{})
		newGsis[m["name"].(string)] = m
	}

	var ops []types.GlobalSecondaryIndexUpdate

	buildCreate := func(m map[string]interface{}) *types.CreateGlobalSecondaryIndexAction {
		c := &types.CreateGlobalSecondaryIndexAction{
			IndexName:  aws.String(m["name"].(string)),
			KeySchema:  expandDynamoGSIKeySchema(m),
			Projection: expandDynamoGSIProjection(m),
		}
		if billingMode == "PROVISIONED" {
			pt := &types.ProvisionedThroughput{}
			if rc, ok := m["read_capacity"]; ok && rc.(int) > 0 {
				pt.ReadCapacityUnits = aws.Int64(int64(rc.(int)))
			}
			if wc, ok := m["write_capacity"]; ok && wc.(int) > 0 {
				pt.WriteCapacityUnits = aws.Int64(int64(wc.(int)))
			}
			c.ProvisionedThroughput = pt
		}
		return c
	}

	for _, raw := range newGsi {
		m := raw.(map[string]interface{})
		if _, exists := oldGsis[m["name"].(string)]; !exists {
			ops = append(ops, types.GlobalSecondaryIndexUpdate{Create: buildCreate(m)})
		}
	}

	for _, raw := range oldGsi {
		oldMap := raw.(map[string]interface{})
		name := oldMap["name"].(string)
		newMap, exists := newGsis[name]
		if !exists {
			ops = append(ops, types.GlobalSecondaryIndexUpdate{
				Delete: &types.DeleteGlobalSecondaryIndexAction{IndexName: aws.String(name)},
			})
			continue
		}

		// Key schema, projection type and projected attributes cannot be changed
		// online; such a change requires recreating the index.
		oldNka := oldMap["non_key_attributes"].(*schema.Set)
		newNka := newMap["non_key_attributes"].(*schema.Set)
		needsRecreate := oldMap["hash_key"].(string) != newMap["hash_key"].(string) ||
			oldMap["range_key"].(string) != newMap["range_key"].(string) ||
			oldMap["projection_type"].(string) != newMap["projection_type"].(string) ||
			!oldNka.Equal(newNka)

		if needsRecreate {
			ops = append(ops, types.GlobalSecondaryIndexUpdate{
				Delete: &types.DeleteGlobalSecondaryIndexAction{IndexName: aws.String(name)},
			})
			ops = append(ops, types.GlobalSecondaryIndexUpdate{Create: buildCreate(newMap)})
			continue
		}

		capacityChanged := oldMap["read_capacity"].(int) != newMap["read_capacity"].(int) ||
			oldMap["write_capacity"].(int) != newMap["write_capacity"].(int)
		if capacityChanged && billingMode == "PROVISIONED" {
			ops = append(ops, types.GlobalSecondaryIndexUpdate{
				Update: &types.UpdateGlobalSecondaryIndexAction{
					IndexName: aws.String(name),
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(int64(newMap["read_capacity"].(int))),
						WriteCapacityUnits: aws.Int64(int64(newMap["write_capacity"].(int))),
					},
				},
			})
		}
	}

	return ops
}

// waitDynamoGSIState waits until the given index disappears (exists=false) or
// reaches ACTIVE (exists=true). An empty IndexStatus is treated as ACTIVE since
// the compatible endpoint may not report intermediate states.
func waitDynamoGSIState(dynamoClient *dynamodb.Client, tableName, indexName string, exists bool, timeout time.Duration) error {
	pending := []string{"CREATING", "UPDATING", "DELETING", "MISSING", "PENDING"}
	target := []string{"ACTIVE"}
	if !exists {
		pending = []string{"CREATING", "UPDATING", "DELETING", "ACTIVE", "PENDING"}
		target = []string{"MISSING"}
	}

	stateConf := BuildStateConf(pending, target, timeout, 5*time.Second, func() (interface{}, string, error) {
		output, err := dynamoClient.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
		if err != nil {
			if isDynamoNotFoundError(err) {
				return tableName, "MISSING", nil
			}
			// Transient transport errors must not be mistaken for a missing index.
			if isDynamoRetryableError(err) {
				return tableName, "PENDING", nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, tableName, "DescribeTable", AlibabaCloudSdkGoERROR)
		}
		if output.Table == nil {
			return tableName, "MISSING", nil
		}
		for _, gsi := range output.Table.GlobalSecondaryIndexes {
			if aws.ToString(gsi.IndexName) == indexName {
				if gsi.IndexStatus == "" {
					return tableName, "ACTIVE", nil
				}
				return tableName, string(gsi.IndexStatus), nil
			}
		}
		return tableName, "MISSING", nil
	})
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, tableName)
	}
	return nil
}

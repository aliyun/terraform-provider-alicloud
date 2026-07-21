package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPolarDBDynamo() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBDynamoCreate,
		Read:   resourceAlicloudPolarDBDynamoRead,
		Update: resourceAlicloudPolarDBDynamoUpdate,
		Delete: resourceAlicloudPolarDBDynamoDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The PolarDB DynamoDB-compatible endpoint URL.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The account name for PolarDB DynamoDB authentication. If not set, the provider's access key will be used.",
			},
			"account_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The authentication password for PolarDB DynamoDB. If not set, the provider's secret key will be used.",
			},
			"db_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the PolarDB cluster where DynamoDB is enabled.",
			},
			"table_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the DynamoDB-compatible table.",
			},
			"attribute": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "List of attribute definitions for the table key schema and indexes.",
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":               {Type: schema.TypeString, Required: true, ForceNew: true},
						"range_key":          {Type: schema.TypeString, Required: true, ForceNew: true},
						"projection_type":    {Type: schema.TypeString, Required: true, ForceNew: true, ValidateFunc: validation.StringInSlice([]string{"ALL", "KEYS_ONLY", "INCLUDE"}, false)},
						"non_key_attributes": {Type: schema.TypeList, Optional: true, ForceNew: true, Elem: &schema.Schema{Type: schema.TypeString}},
					},
				},
			},
			"stream_enabled":   {Type: schema.TypeBool, Optional: true},
			"stream_view_type": {Type: schema.TypeString, Optional: true, Computed: true, ValidateFunc: validation.StringInSlice([]string{"NEW_IMAGE", "OLD_IMAGE", "NEW_AND_OLD_IMAGES", "KEYS_ONLY"}, false)},
			"ttl": {
				Type: schema.TypeList, Optional: true, Computed: true, MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled":        {Type: schema.TypeBool, Optional: true, Default: false},
						"attribute_name": {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"point_in_time_recovery": {
				Type: schema.TypeList, Optional: true, Computed: true, MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {Type: schema.TypeBool, Required: true},
					},
				},
			},
			"server_side_encryption": {
				Type: schema.TypeList, Optional: true, Computed: true, MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {Type: schema.TypeBool, Required: true},
					},
				},
			},
			"tags":         {Type: schema.TypeMap, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"arn":          {Type: schema.TypeString, Computed: true},
			"stream_arn":   {Type: schema.TypeString, Computed: true},
			"stream_label": {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceAlicloudPolarDBDynamoCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	endpoint := d.Get("endpoint").(string)
	tableName := d.Get("table_name").(string)
	accessKey := d.Get("account_name").(string)
	secretKey := d.Get("account_auth").(string)

	dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accessKey, secretKey)
	if err != nil {
		return WrapError(err)
	}

	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
	}

	keySchema := make([]*dynamodb.KeySchemaElement, 0)
	if v, ok := d.GetOk("hash_key"); ok && v.(string) != "" {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{AttributeName: aws.String(v.(string)), KeyType: aws.String("HASH")})
	}
	if v, ok := d.GetOk("range_key"); ok && v.(string) != "" {
		keySchema = append(keySchema, &dynamodb.KeySchemaElement{AttributeName: aws.String(v.(string)), KeyType: aws.String("RANGE")})
	}
	if len(keySchema) > 0 {
		input.KeySchema = keySchema
	}

	if v, ok := d.GetOk("attribute"); ok {
		attrs := make([]*dynamodb.AttributeDefinition, 0)
		for _, attrRaw := range v.(*schema.Set).List() {
			attr := attrRaw.(map[string]interface{})
			attrs = append(attrs, &dynamodb.AttributeDefinition{
				AttributeName: aws.String(attr["name"].(string)),
				AttributeType: aws.String(attr["type"].(string)),
			})
		}
		input.AttributeDefinitions = attrs
	}

	if v, ok := d.GetOk("billing_mode"); ok {
		input.BillingMode = aws.String(v.(string))
	}

	if d.Get("billing_mode").(string) == "PROVISIONED" {
		pt := &dynamodb.ProvisionedThroughput{}
		if v, ok := d.GetOk("read_capacity"); ok {
			pt.ReadCapacityUnits = aws.Int64(int64(v.(int)))
		}
		if v, ok := d.GetOk("write_capacity"); ok {
			pt.WriteCapacityUnits = aws.Int64(int64(v.(int)))
		}
		input.ProvisionedThroughput = pt
	}

	if v, ok := d.GetOk("global_secondary_index"); ok {
		gsis := make([]*dynamodb.GlobalSecondaryIndex, 0)
		for _, gsiRaw := range v.(*schema.Set).List() {
			gsi := gsiRaw.(map[string]interface{})
			gsiConfig := &dynamodb.GlobalSecondaryIndex{IndexName: aws.String(gsi["name"].(string))}
			gsiKs := make([]*dynamodb.KeySchemaElement, 0)
			if hk, ok := gsi["hash_key"]; ok && hk.(string) != "" {
				gsiKs = append(gsiKs, &dynamodb.KeySchemaElement{AttributeName: aws.String(hk.(string)), KeyType: aws.String("HASH")})
			}
			if rk, ok := gsi["range_key"]; ok && rk.(string) != "" {
				gsiKs = append(gsiKs, &dynamodb.KeySchemaElement{AttributeName: aws.String(rk.(string)), KeyType: aws.String("RANGE")})
			}
			if len(gsiKs) > 0 {
				gsiConfig.KeySchema = gsiKs
			}
			proj := &dynamodb.Projection{ProjectionType: aws.String(gsi["projection_type"].(string))}
			if nka, ok := gsi["non_key_attributes"]; ok {
				nkaList := expandStringList(nka.(*schema.Set).List())
				if len(nkaList) > 0 {
					proj.NonKeyAttributes = aws.StringSlice(nkaList)
				}
			}
			gsiConfig.Projection = proj
			if d.Get("billing_mode").(string) == "PROVISIONED" {
				gsiPt := &dynamodb.ProvisionedThroughput{}
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
		lsis := make([]*dynamodb.LocalSecondaryIndex, 0)
		for _, lsiRaw := range v.(*schema.Set).List() {
			lsi := lsiRaw.(map[string]interface{})
			lsiConfig := &dynamodb.LocalSecondaryIndex{IndexName: aws.String(lsi["name"].(string))}
			lsiKs := make([]*dynamodb.KeySchemaElement, 0)
			if hk, ok := d.GetOk("hash_key"); ok && hk.(string) != "" {
				lsiKs = append(lsiKs, &dynamodb.KeySchemaElement{AttributeName: aws.String(hk.(string)), KeyType: aws.String("HASH")})
			}
			if rk, ok := lsi["range_key"]; ok && rk.(string) != "" {
				lsiKs = append(lsiKs, &dynamodb.KeySchemaElement{AttributeName: aws.String(rk.(string)), KeyType: aws.String("RANGE")})
			}
			if len(lsiKs) > 0 {
				lsiConfig.KeySchema = lsiKs
			}
			proj := &dynamodb.Projection{ProjectionType: aws.String(lsi["projection_type"].(string))}
			if nka, ok := lsi["non_key_attributes"]; ok {
				nkaList := expandStringList(nka.([]interface{}))
				if len(nkaList) > 0 {
					proj.NonKeyAttributes = aws.StringSlice(nkaList)
				}
			}
			lsiConfig.Projection = proj
			lsis = append(lsis, lsiConfig)
		}
		if len(lsis) > 0 {
			input.LocalSecondaryIndexes = lsis
		}
	}

	if v, ok := d.GetOk("stream_enabled"); ok && v.(bool) {
		streamSpec := &dynamodb.StreamSpecification{StreamEnabled: aws.Bool(true)}
		if svt, ok := d.GetOk("stream_view_type"); ok && svt.(string) != "" {
			streamSpec.StreamViewType = aws.String(svt.(string))
		}
		input.StreamSpecification = streamSpec
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := make([]*dynamodb.Tag, 0)
		for k, val := range v.(map[string]interface{}) {
			tags = append(tags, &dynamodb.Tag{Key: aws.String(k), Value: aws.String(val.(string))})
		}
		if len(tags) > 0 {
			input.Tags = tags
		}
	}

	var output *dynamodb.CreateTableOutput
	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
		output, err = dynamoClient.CreateTable(input)
		if err != nil {
			if isDynamoRetryableError(err) {
				log.Printf("[DEBUG] CreateTable %s failed with retryable error, will retry: %s", tableName, err)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateTable", output, nil, input)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_dynamo", "CreateTable", AlibabaCloudSdkGoERROR)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	d.SetId(fmt.Sprintf("%s%s%s", dbClusterId, COLON_SEPARATED, tableName))

	polarDBService := PolarDBService{client}
	stateConf := BuildStateConf(
		[]string{"CREATING", "UPDATING"},
		[]string{"ACTIVE"},
		d.Timeout(schema.TimeoutCreate),
		5*time.Second,
		polarDBService.PolarDBDynamoTableStateRefreshFunc(dbClusterId, tableName, endpoint, accessKey, secretKey),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudPolarDBDynamoRead(d, meta)
}

func resourceAlicloudPolarDBDynamoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	tableName := parts[1]
	endpoint := d.Get("endpoint").(string)
	accessKey := d.Get("account_name").(string)
	secretKey := d.Get("account_auth").(string)

	// On import only the resource ID is available, so resolve the DynamoDB-compatible
	// endpoint address from the cluster to rebuild the connection info.
	if endpoint == "" {
		polarDBService := PolarDBService{client}
		endpoint, err = polarDBService.DescribePolarDBDynamoEndpointAddress(dbClusterId)
		if err != nil {
			return WrapError(err)
		}
		d.Set("endpoint", endpoint)
	}

	dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accessKey, secretKey)
	if err != nil {
		return WrapError(err)
	}

	var output *dynamodb.DescribeTableOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		output, err = dynamoClient.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
		if err != nil {
			if isDynamoNotFoundError(err) {
				return resource.NonRetryableError(err)
			}
			if isDynamoRetryableError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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

	tableDesc := output.Table
	d.Set("db_cluster_id", dbClusterId)
	d.Set("table_name", tableName)

	for _, ks := range tableDesc.KeySchema {
		if aws.StringValue(ks.KeyType) == "HASH" {
			d.Set("hash_key", aws.StringValue(ks.AttributeName))
		} else if aws.StringValue(ks.KeyType) == "RANGE" {
			d.Set("range_key", aws.StringValue(ks.AttributeName))
		}
	}

	attrSet := make([]map[string]interface{}, 0)
	for _, attr := range tableDesc.AttributeDefinitions {
		attrSet = append(attrSet, map[string]interface{}{
			"name": aws.StringValue(attr.AttributeName),
			"type": aws.StringValue(attr.AttributeType),
		})
	}
	d.Set("attribute", attrSet)

	if tableDesc.BillingModeSummary != nil {
		d.Set("billing_mode", aws.StringValue(tableDesc.BillingModeSummary.BillingMode))
	}
	if tableDesc.ProvisionedThroughput != nil {
		d.Set("read_capacity", int(aws.Int64Value(tableDesc.ProvisionedThroughput.ReadCapacityUnits)))
		d.Set("write_capacity", int(aws.Int64Value(tableDesc.ProvisionedThroughput.WriteCapacityUnits)))
	}

	gsiSet := make([]map[string]interface{}, 0)
	for _, gsi := range tableDesc.GlobalSecondaryIndexes {
		gsiConfig := map[string]interface{}{"name": aws.StringValue(gsi.IndexName)}
		for _, ks := range gsi.KeySchema {
			if aws.StringValue(ks.KeyType) == "HASH" {
				gsiConfig["hash_key"] = aws.StringValue(ks.AttributeName)
			} else if aws.StringValue(ks.KeyType) == "RANGE" {
				gsiConfig["range_key"] = aws.StringValue(ks.AttributeName)
			}
		}
		if gsi.Projection != nil {
			gsiConfig["projection_type"] = aws.StringValue(gsi.Projection.ProjectionType)
			if len(gsi.Projection.NonKeyAttributes) > 0 {
				gsiConfig["non_key_attributes"] = aws.StringValueSlice(gsi.Projection.NonKeyAttributes)
			}
		}
		if gsi.ProvisionedThroughput != nil {
			gsiConfig["read_capacity"] = int(aws.Int64Value(gsi.ProvisionedThroughput.ReadCapacityUnits))
			gsiConfig["write_capacity"] = int(aws.Int64Value(gsi.ProvisionedThroughput.WriteCapacityUnits))
		}
		gsiSet = append(gsiSet, gsiConfig)
	}
	d.Set("global_secondary_index", gsiSet)

	lsiSet := make([]map[string]interface{}, 0)
	for _, lsi := range tableDesc.LocalSecondaryIndexes {
		lsiConfig := map[string]interface{}{"name": aws.StringValue(lsi.IndexName)}
		for _, ks := range lsi.KeySchema {
			if aws.StringValue(ks.KeyType) == "RANGE" {
				lsiConfig["range_key"] = aws.StringValue(ks.AttributeName)
			}
		}
		if lsi.Projection != nil {
			lsiConfig["projection_type"] = aws.StringValue(lsi.Projection.ProjectionType)
			if len(lsi.Projection.NonKeyAttributes) > 0 {
				lsiConfig["non_key_attributes"] = aws.StringValueSlice(lsi.Projection.NonKeyAttributes)
			}
		}
		lsiSet = append(lsiSet, lsiConfig)
	}
	d.Set("local_secondary_index", lsiSet)

	if tableDesc.StreamSpecification != nil {
		d.Set("stream_enabled", aws.BoolValue(tableDesc.StreamSpecification.StreamEnabled))
		d.Set("stream_view_type", aws.StringValue(tableDesc.StreamSpecification.StreamViewType))
	}

	d.Set("arn", aws.StringValue(tableDesc.TableArn))
	d.Set("stream_arn", aws.StringValue(tableDesc.LatestStreamArn))
	d.Set("stream_label", aws.StringValue(tableDesc.LatestStreamLabel))

	ttlOutput, ttlErr := dynamoClient.DescribeTimeToLive(&dynamodb.DescribeTimeToLiveInput{TableName: aws.String(tableName)})
	if ttlErr == nil && ttlOutput.TimeToLiveDescription != nil {
		ttlConfig := map[string]interface{}{"enabled": false}
		if aws.StringValue(ttlOutput.TimeToLiveDescription.TimeToLiveStatus) == "ENABLED" {
			ttlConfig["enabled"] = true
		}
		if attrName := aws.StringValue(ttlOutput.TimeToLiveDescription.AttributeName); attrName != "" {
			ttlConfig["attribute_name"] = attrName
		}
		d.Set("ttl", []interface{}{ttlConfig})
	}

	pitrOutput, pitrErr := dynamoClient.DescribeContinuousBackups(&dynamodb.DescribeContinuousBackupsInput{TableName: aws.String(tableName)})
	if pitrErr == nil && pitrOutput.ContinuousBackupsDescription != nil && pitrOutput.ContinuousBackupsDescription.PointInTimeRecoveryDescription != nil {
		pitrConfig := map[string]interface{}{"enabled": false}
		if aws.StringValue(pitrOutput.ContinuousBackupsDescription.PointInTimeRecoveryDescription.PointInTimeRecoveryStatus) == "ENABLED" {
			pitrConfig["enabled"] = true
		}
		d.Set("point_in_time_recovery", []interface{}{pitrConfig})
	}

	return nil
}

func resourceAlicloudPolarDBDynamoUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	tableName := parts[1]
	endpoint := d.Get("endpoint").(string)
	accessKey := d.Get("account_name").(string)
	secretKey := d.Get("account_auth").(string)

	dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accessKey, secretKey)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChanges("read_capacity", "write_capacity") && d.Get("billing_mode").(string) == "PROVISIONED" {
		updateInput := &dynamodb.UpdateTableInput{TableName: aws.String(tableName)}
		pt := &dynamodb.ProvisionedThroughput{}
		if v, ok := d.GetOk("read_capacity"); ok {
			pt.ReadCapacityUnits = aws.Int64(int64(v.(int)))
		}
		if v, ok := d.GetOk("write_capacity"); ok {
			pt.WriteCapacityUnits = aws.Int64(int64(v.(int)))
		}
		updateInput.ProvisionedThroughput = pt

		err = resource.Retry(8*time.Minute, func() *resource.RetryError {
			_, err = dynamoClient.UpdateTable(updateInput)
			if err != nil {
				if isDynamoRetryableError(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
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
			polarDBService.PolarDBDynamoTableStateRefreshFunc(dbClusterId, tableName, endpoint, accessKey, secretKey),
		)
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChanges("stream_enabled", "stream_view_type") {
		updateInput := &dynamodb.UpdateTableInput{TableName: aws.String(tableName)}
		streamSpec := &dynamodb.StreamSpecification{}
		if v, ok := d.GetOk("stream_enabled"); ok {
			streamSpec.StreamEnabled = aws.Bool(v.(bool))
		}
		if v, ok := d.GetOk("stream_view_type"); ok && v.(string) != "" {
			streamSpec.StreamViewType = aws.String(v.(string))
		}
		updateInput.StreamSpecification = streamSpec

		err = resource.Retry(8*time.Minute, func() *resource.RetryError {
			_, err = dynamoClient.UpdateTable(updateInput)
			if err != nil {
				if isDynamoRetryableError(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateTable", AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("ttl") {
		ttlList := d.Get("ttl").([]interface{})
		if len(ttlList) > 0 && ttlList[0] != nil {
			ttlConfig := ttlList[0].(map[string]interface{})
			ttlSpec := &dynamodb.TimeToLiveSpecification{Enabled: aws.Bool(ttlConfig["enabled"].(bool))}
			if attrName, ok := ttlConfig["attribute_name"]; ok && attrName.(string) != "" {
				ttlSpec.AttributeName = aws.String(attrName.(string))
			}
			_, err = dynamoClient.UpdateTimeToLive(&dynamodb.UpdateTimeToLiveInput{
				TableName:               aws.String(tableName),
				TimeToLiveSpecification: ttlSpec,
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateTimeToLive", AlibabaCloudSdkGoERROR)
			}
		}
	}

	if d.HasChange("point_in_time_recovery") {
		pitrList := d.Get("point_in_time_recovery").([]interface{})
		if len(pitrList) > 0 && pitrList[0] != nil {
			pitrConfig := pitrList[0].(map[string]interface{})
			_, err = dynamoClient.UpdateContinuousBackups(&dynamodb.UpdateContinuousBackupsInput{
				TableName: aws.String(tableName),
				PointInTimeRecoverySpecification: &dynamodb.PointInTimeRecoverySpecification{
					PointInTimeRecoveryEnabled: aws.Bool(pitrConfig["enabled"].(bool)),
				},
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateContinuousBackups", AlibabaCloudSdkGoERROR)
			}
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		oldTagsMap := oldTags.(map[string]interface{})
		newTagsMap := newTags.(map[string]interface{})

		removeKeys := make([]*string, 0)
		for k := range oldTagsMap {
			if _, exists := newTagsMap[k]; !exists {
				removeKeys = append(removeKeys, aws.String(k))
			}
		}
		if len(removeKeys) > 0 {
			arn := d.Get("arn").(string)
			_, err = dynamoClient.UntagResource(&dynamodb.UntagResourceInput{ResourceArn: aws.String(arn), TagKeys: removeKeys})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UntagResource", AlibabaCloudSdkGoERROR)
			}
		}

		addTags := make([]*dynamodb.Tag, 0)
		for k, v := range newTagsMap {
			addTags = append(addTags, &dynamodb.Tag{Key: aws.String(k), Value: aws.String(v.(string))})
		}
		if len(addTags) > 0 {
			arn := d.Get("arn").(string)
			_, err = dynamoClient.TagResource(&dynamodb.TagResourceInput{ResourceArn: aws.String(arn), Tags: addTags})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "TagResource", AlibabaCloudSdkGoERROR)
			}
		}
	}

	return resourceAlicloudPolarDBDynamoRead(d, meta)
}

func resourceAlicloudPolarDBDynamoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	tableName := parts[1]
	endpoint := d.Get("endpoint").(string)
	accessKey := d.Get("account_name").(string)
	secretKey := d.Get("account_auth").(string)

	dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accessKey, secretKey)
	if err != nil {
		return WrapError(err)
	}

	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
		_, err = dynamoClient.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
		if err != nil {
			if isDynamoNotFoundError(err) {
				return resource.NonRetryableError(err)
			}
			if isDynamoRetryableError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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
		polarDBService.PolarDBDynamoTableStateRefreshFunc(dbClusterId, tableName, endpoint, accessKey, secretKey),
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
	errMsg := err.Error()
	return strings.Contains(errMsg, "ResourceNotFoundException") ||
		strings.Contains(errMsg, "TableNotFound") ||
		strings.Contains(errMsg, "InvalidDBClusterId.NotFound")
}

func isDynamoRetryableError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "InternalServerError") ||
		strings.Contains(errMsg, "ProvisionedThroughputExceededException") ||
		strings.Contains(errMsg, "ThrottlingException") ||
		strings.Contains(errMsg, "LimitExceededException") ||
		// The public endpoint address/domain returned by alicloud_polardb_endpoint_address
		// may not be resolvable or reachable immediately after creation (DNS propagation
		// delay or the listener not fully up yet), so treat these transient network
		// errors as retryable instead of failing the whole Create/Read/Update/Delete.
		strings.Contains(errMsg, "no such host") ||
		strings.Contains(errMsg, "dial tcp") ||
		strings.Contains(errMsg, "connection refused") ||
		strings.Contains(errMsg, "i/o timeout") ||
		strings.Contains(errMsg, "EOF") ||
		strings.Contains(errMsg, "RequestError: send request failed")
}

package mongodb

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	sdkSchema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkTerraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func unitMongodbAccountResourceData(t *testing.T, id string, attrs map[string]interface{}, changes map[string]interface{}) *sdkSchema.ResourceData {
	t.Helper()
	state := &sdkTerraform.InstanceState{ID: id, Attributes: map[string]string{}}
	for k, v := range attrs {
		state.Attributes[k] = formatUnitValue(v)
	}
	diff := &sdkTerraform.InstanceDiff{Attributes: map[string]*sdkTerraform.ResourceAttrDiff{}}
	for k, v := range changes {
		diff.Attributes[k] = &sdkTerraform.ResourceAttrDiff{Old: formatUnitValue(attrs[k]), New: formatUnitValue(v)}
	}
	d, err := sdkSchema.InternalMap(resourceAliCloudMongodbAccount().Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("failed to build resource data: %v", err)
	}
	return d
}

func formatUnitValue(v interface{}) string {
	return fmt.Sprint(v)
}

func TestUnitMongodbAccountSchemaContract(t *testing.T) {
	r := resourceAliCloudMongodbAccount()
	assert.True(t, r.Schema["db_name"].ForceNew, "db_name is part of the account identity")
	assert.True(t, r.Schema["db_name"].Computed)
	assert.True(t, r.Schema["roles"].Computed, "roles are returned by DescribeAccounts")
	assert.True(t, r.Schema["roles"].ForceNew)
	assert.True(t, r.Schema["character_type"].ForceNew)
}

func TestUnitMongodbAccountDescribeRejectsInvalidId(t *testing.T) {
	service := &MongodbServiceV2{}
	if _, err := service.DescribeMongodbAccount("dds-1"); err == nil {
		t.Fatal("expected an error for an id without an account name")
	}
	if _, err := service.DescribeMongodbAccount("dds-1:db:acc:extra"); err == nil {
		t.Fatal("expected an error for an id with too many parts")
	}
}

func TestUnitMongodbAccountStateRefreshChecksFieldError(t *testing.T) {
	service := &MongodbServiceV2{}
	refresh := service.MongodbAccountStateRefreshFuncWithApi("dds-1:admin:acc", "AccountStatus", nil, func(id string) (map[string]interface{}, error) {
		return map[string]interface{}{"DBName": "admin"}, nil
	})
	if _, _, err := refresh(); err == nil {
		t.Fatal("expected an error when the watched field is missing from the response")
	}
}

func TestUnitMongodbAccountReadNotFoundClearsId(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
	}, nil)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *MongodbServiceV2, id string) (map[string]interface{}, error) {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Account", id), errs.NotFoundMsg, "mock")
	})
	defer patches.Reset()

	assert.Nil(t, resourceAliCloudMongodbAccount().Read(d, &connectivity.AliyunClient{}))
	assert.Equal(t, "", d.Id())
}

func TestUnitMongodbAccountReadSetsDatabaseScopedId(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc1", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc1",
	}, nil)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *MongodbServiceV2, id string) (map[string]interface{}, error) {
		return map[string]interface{}{
			"DBInstanceId":  "dds-1",
			"AccountName":   "acc1",
			"DBName":        "admin",
			"CharacterType": "db",
			"AccountStatus": "Available",
			"Roles": map[string]interface{}{
				"Role": []interface{}{map[string]interface{}{"DBName": "admin", "Role": "read"}},
			},
		}, nil
	})
	defer patches.Reset()

	assert.Nil(t, resourceAliCloudMongodbAccount().Read(d, &connectivity.AliyunClient{}))
	assert.Equal(t, "dds-1:admin:acc1", d.Id())
	assert.Equal(t, "admin", d.Get("db_name"))
	assert.Equal(t, 1, d.Get("roles").(*sdkSchema.Set).Len())
}

func TestUnitMongodbAccountReadPropagatesError(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
	}, nil)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *MongodbServiceV2, id string) (map[string]interface{}, error) {
		return nil, errors.New("Forbidden.RAM")
	})
	defer patches.Reset()

	err := resourceAliCloudMongodbAccount().Read(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected the DescribeAccounts error to be propagated")
	}
	assert.Equal(t, "dds-1:acc", d.Id())
}

func TestUnitMongodbAccountUpdateDescriptionDescribeError(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:admin:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
		"db_name":      "admin",
	}, map[string]interface{}{
		"account_description": "updated",
	})
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *MongodbServiceV2, id string) (map[string]interface{}, error) {
		return nil, errors.New("Forbidden.RAM")
	})
	defer patches.Reset()

	if err := resourceAliCloudMongodbAccount().Update(d, &connectivity.AliyunClient{}); err == nil {
		t.Fatal("expected the pre-update DescribeAccounts error to be propagated")
	}
}

func TestUnitMongodbAccountUpdateShardDescriptionRejected(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:admin:acc", map[string]interface{}{
		"instance_id":  "dds-1",
		"account_name": "acc",
		"db_name":      "admin",
	}, map[string]interface{}{
		"account_description": "updated",
	})
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&MongodbServiceV2{}), "DescribeMongodbAccount", func(_ *MongodbServiceV2, id string) (map[string]interface{}, error) {
		return map[string]interface{}{"DBInstanceId": "dds-1", "AccountName": "acc", "DBName": "admin", "CharacterType": "db", "AccountStatus": "Available"}, nil
	})
	defer patches.Reset()

	err := resourceAliCloudMongodbAccount().Update(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected a shard account description update to be rejected")
	}
	assert.Contains(t, err.Error(), "replica set")
}

func TestUnitMongodbAccountCreateRejectsRolesForReplicaSet(t *testing.T) {
	d := sdkSchema.TestResourceDataRaw(t, resourceAliCloudMongodbAccount().Schema, map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "root",
		"account_password": "YourPassword+12345",
		"character_type":   "normal",
		"db_name":          "admin",
		"roles": []interface{}{
			map[string]interface{}{"role": "read", "db_name": "admin"},
		},
	})
	err := resourceAliCloudMongodbAccount().Create(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected roles on a replica set account to be rejected")
	}
	assert.Contains(t, err.Error(), "config server account")
}

func TestUnitMongodbAccountCreateRejectsRolesForShardAccount(t *testing.T) {
	d := sdkSchema.TestResourceDataRaw(t, resourceAliCloudMongodbAccount().Schema, map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "shard_acc",
		"account_password": "YourPassword+12345",
		"character_type":   "db",
		"db_name":          "admin",
		"roles": []interface{}{
			map[string]interface{}{"role": "read", "db_name": "admin"},
		},
	})
	err := resourceAliCloudMongodbAccount().Create(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected roles on a shard direct-connection account to be rejected")
	}
	assert.Contains(t, err.Error(), "InvalidRoles.NotSupported")
}

func TestUnitMongodbAccountCreateRequiresRolesForCsAccount(t *testing.T) {
	d := sdkSchema.TestResourceDataRaw(t, resourceAliCloudMongodbAccount().Schema, map[string]interface{}{
		"instance_id":      "dds-1",
		"account_name":     "cs_acc",
		"account_password": "YourPassword+12345",
		"character_type":   "cs",
		"db_name":          "admin",
	})
	err := resourceAliCloudMongodbAccount().Create(d, &connectivity.AliyunClient{})
	if err == nil {
		t.Fatal("expected a config server account without roles to be rejected")
	}
	assert.Contains(t, err.Error(), "required")
}

func TestUnitMongodbAccountDeleteNoop(t *testing.T) {
	d := unitMongodbAccountResourceData(t, "dds-1:acc", map[string]interface{}{}, nil)
	assert.Nil(t, resourceAliCloudMongodbAccount().Delete(d, &connectivity.AliyunClient{}))
}

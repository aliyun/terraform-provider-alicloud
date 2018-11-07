package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type DrdsService struct {
	client *connectivity.AliyunClient
}

// crate Drdsinstance
func (s *DrdsService) CreateDrdsInstance(req *drds.CreateDrdsInstanceRequest) (response *drds.CreateDrdsInstanceResponse, err error) {

	//resp, err := s.client.WithDrdsClient(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.CreateDrdsInstance(req)
	})
	resp, _ := raw.(*drds.CreateDrdsInstanceResponse)

	if err != nil {
		return resp, fmt.Errorf("createDrdsInstance got an error: %#v", err)
	}

	return resp, nil
}

func (s *DrdsService) DescribeDrdsInstance(drdsInstanceId string) (response *drds.DescribeDrdsInstanceResponse, err error) {
	req := drds.CreateDescribeDrdsInstanceRequest()
	req.DrdsInstanceId = drdsInstanceId
	//resp, err := client.drdsconn.DescribeDrdsInstance(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstance(req)
	})
	if err != nil {
		return nil, fmt.Errorf("describe drds instance error")
	}
	resp, _ := raw.(*drds.DescribeDrdsInstanceResponse)

	if resp == nil {
		return resp, GetNotFoundErrorFromString(GetNotFoundMessage("Instance", drdsInstanceId))

	}
	return resp, nil
}

func (s *DrdsService) DescribeDrdsInstances(regionId string) (response *drds.DescribeDrdsInstancesResponse, err error) {
	req := drds.CreateDescribeDrdsInstancesRequest()
	req.Type = string(Private)
	//resp, err := client.drdsconn.DescribeDrdsInstances(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstances(req)
	})
	resp, _ := raw.(*drds.DescribeDrdsInstancesResponse)

	return resp, err

}

func (s *DrdsService) DescribeRegions() (response *drds.DescribeRegionsResponse, err error) {
	req := drds.CreateDescribeRegionsRequest()
	//resp, err := client.drdsconn.DescribeRegions(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeRegions(req)
	})
	resp, _ := raw.(*drds.DescribeRegionsResponse)

	return resp, err
}

func (s *DrdsService) ModifyDrdsInstanceDescription(request *drds.ModifyDrdsInstanceDescriptionRequest) (response *drds.ModifyDrdsInstanceDescriptionResponse, err error) {
	req := drds.CreateModifyDrdsInstanceDescriptionRequest()
	req.DrdsInstanceId = request.DrdsInstanceId
	req.Description = request.Description
	//resp, err := client.drdsconn.ModifyDrdsInstanceDescription(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.ModifyDrdsInstanceDescription(req)
	})
	resp, _ := raw.(*drds.ModifyDrdsInstanceDescriptionResponse)

	return resp, err

}

func (s *DrdsService) RemoveDrdsInstance(drdsInstanceId string) (response *drds.RemoveDrdsInstanceResponse, err error) {
	req := drds.CreateRemoveDrdsInstanceRequest()
	req.DrdsInstanceId = drdsInstanceId
	//resp, err := client.drdsconn.RemoveDrdsInstance(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.RemoveDrdsInstance(req)
	})
	resp, _ := raw.(*drds.RemoveDrdsInstanceResponse)

	return resp, err
}

func (s *DrdsService) CreateDrdsDB(drdsInstanceId string, dbName string, encode string, password string) (response *drds.CreateDrdsDBResponse, err error) {
	req := drds.CreateCreateDrdsDBRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.Encode = encode
	req.Password = password
	//resp, err := client.drdsconn.CreateDrdsDB(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.CreateDrdsDB(req)
	})
	resp, _ := raw.(*drds.CreateDrdsDBResponse)

	return resp, err
}

func (s *DrdsService) DescribeDrdsDB(dbName string, drdsInstanceId string) (response *drds.DescribeDrdsDBResponse, err error) {
	req := drds.CreateDescribeDrdsDBRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	//resp, err := client.drdsconn.DescribeDrdsDB(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsDB(req)
	})
	resp, _ := raw.(*drds.DescribeDrdsDBResponse)

	return resp, err
}

func (s *DrdsService) DeleteDrdsDB(dbName string, drdsInstanceId string) (response *drds.DeleteDrdsDBResponse, err error) {
	req := drds.CreateDeleteDrdsDBRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	//resp, err := client.drdsconn.DeleteDrdsDB(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DeleteDrdsDB(req)
	})
	resp, _ := raw.(*drds.DeleteDrdsDBResponse)

	return resp, err
}

func (s *DrdsService) ModifyDrdsDBPasswd(dbName string, drdsInstanceId string, newPasswd string) (response *drds.ModifyDrdsDBPasswdResponse, err error) {
	req := drds.CreateModifyDrdsDBPasswdRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	//resp, err := client.drdsconn.ModifyDrdsDBPasswd(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.ModifyDrdsDBPasswd(req)
	})
	resp, _ := raw.(*drds.ModifyDrdsDBPasswdResponse)

	return resp, err
}

func (s *DrdsService) DescribeDrdsDBs(drdsInstanceId string) (response *drds.DescribeDrdsDBsResponse, err error) {
	req := drds.CreateDescribeDrdsDBsRequest()
	req.DrdsInstanceId = drdsInstanceId
	//resp, err := client.drdsconn.DescribeDrdsDBs(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsDBs(req)
	})
	resp, _ := raw.(*drds.DescribeDrdsDBsResponse)

	return resp, err
}

func (s *DrdsService) DescribeDrdsDBIpWhiteList(drdsInstanceId string, dbName string, groupName string) (response *drds.DescribeDrdsDBIpWhiteListResponse, err error) {
	req := drds.CreateDescribeDrdsDBIpWhiteListRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	if groupName == "" {
		req.GroupName = "DEFAULT_GROUP"
	} else {
		req.GroupName = groupName
	}
	//resp, err := client.drdsconn.DescribeDrdsDBIpWhiteList(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsDBIpWhiteList(req)
	})
	resp, _ := raw.(*drds.DescribeDrdsDBIpWhiteListResponse)

	return resp, err
}

func (s *DrdsService) ModifyDrdsIpWhiteList(drdsInstanceId string, dbName string, ipWhiteList string, mode bool, groupName string) (response *drds.ModifyDrdsIpWhiteListResponse, err error) {
	req := drds.CreateModifyDrdsIpWhiteListRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.Mode = requests.NewBoolean(mode)
	req.IpWhiteList = ipWhiteList
	if groupName == "" {
		req.GroupName = "DEFAULT_GROUP"
	} else {
		req.GroupName = groupName
	}
	//resp, err := client.drdsconn.ModifyDrdsIpWhiteList(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.ModifyDrdsIpWhiteList(req)
	})
	resp, _ := raw.(*drds.ModifyDrdsIpWhiteListResponse)

	return resp, err
}

func (s *DrdsService) DescribeRdsList(drdsInstanceId string, dbName string) (response *drds.DescribeRdsListResponse, err error) {
	req := drds.CreateDescribeRdsListRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	//resp, err := client.drdsconn.DescribeRdsList(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeRdsList(req)
	})
	resp, _ := raw.(*drds.DescribeRdsListResponse)

	return resp, err
}

func (s *DrdsService) ModifyFullTableScan(drdsInstanceId string, dbName string, tableNames string, fulltableScan bool) (response *drds.ModifyFullTableScanResponse, err error) {
	req := drds.CreateModifyFullTableScanRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.TableNames = tableNames
	req.FullTableScan = requests.NewBoolean(fulltableScan)
	//resp, err := client.drdsconn.ModifyFullTableScan(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.ModifyFullTableScan(req)
	})
	resp, _ := raw.(*drds.ModifyFullTableScanResponse)

	return resp, err
}

func (s *DrdsService) DescribeShardDBs(drdsInstanceId string, dbName string) (response *drds.DescribeShardDBsResponse, err error) {
	req := drds.CreateDescribeShardDBsRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	//resp, err := client.drdsconn.DescribeShardDBs(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeShardDBs(req)
	})
	resp, _ := raw.(*drds.DescribeShardDBsResponse)

	return resp, err
}

func (s *DrdsService) DeleteFailedDrdsDB(drdsInstanceId string, dbName string) (response *drds.DeleteFailedDrdsDBResponse, err error) {
	req := drds.CreateDeleteFailedDrdsDBRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	//resp, err := client.drdsconn.DeleteFailedDrdsDB(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DeleteFailedDrdsDB(req)
	})
	resp, _ := raw.(*drds.DeleteFailedDrdsDBResponse)

	return resp, err
}

func (s *DrdsService) ModifyRdsReadWeight(drdsInstanceId string, dbName string, instanceNames string, weight string) (response *drds.ModifyRdsReadWeightResponse, err error) {
	req := drds.CreateModifyRdsReadWeightRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	//resp, err := client.drdsconn.ModifyRdsReadWeight(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.ModifyRdsReadWeight(req)
	})
	resp, _ := raw.(*drds.ModifyRdsReadWeightResponse)

	return resp, err
}

func (s *DrdsService) CreateReadOnlyAccount(drdsInstanceId string, dbName string, password string) (response *drds.CreateReadOnlyAccountResponse, err error) {
	req := drds.CreateCreateReadOnlyAccountRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.Password = password
	//resp, err := client.drdsconn.CreateReadOnlyAccount(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.CreateReadOnlyAccount(req)
	})
	resp, _ := raw.(*drds.CreateReadOnlyAccountResponse)

	return resp, err
}

func (s *DrdsService) DescribeReadOnlyAccount(drdsInstanceId string, dbName string) (response *drds.DescribeReadOnlyAccountResponse, err error) {
	req := drds.CreateDescribeReadOnlyAccountRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	//resp, err := client.drdsconn.DescribeReadOnlyAccount(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeReadOnlyAccount(req)
	})
	resp, _ := raw.(*drds.DescribeReadOnlyAccountResponse)

	return resp, err
}

func (s *DrdsService) ModifyReadOnlyAccountPassword(drdsInstanceId string, dbName string, accountName string, originPassword string, newPassword string) (response *drds.ModifyReadOnlyAccountPasswordResponse, err error) {
	req := drds.CreateModifyReadOnlyAccountPasswordRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	req.AccountName = accountName
	req.OriginPassword = originPassword
	req.NewPasswd = newPassword
	//resp, err := client.drdsconn.ModifyReadOnlyAccountPassword(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.ModifyReadOnlyAccountPassword(req)
	})
	resp, _ := raw.(*drds.ModifyReadOnlyAccountPasswordResponse)

	return resp, err
}

func (s *DrdsService) RemoveReadOnlyAccount(drdsInstanceId string, dbName string, accountName string) (response *drds.RemoveReadOnlyAccountResponse, err error) {
	req := drds.CreateRemoveReadOnlyAccountRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.AccountName = accountName
	//resp, err := client.drdsconn.RemoveReadOnlyAccount(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.RemoveReadOnlyAccount(req)
	})
	resp, _ := raw.(*drds.RemoveReadOnlyAccountResponse)

	return resp, err
}

func (s *DrdsService) CreateDrdsAccount(drdsInstanceId string, dbName string, password string) (response *drds.CreateDrdsAccountResponse, err error) {
	req := drds.CreateCreateDrdsAccountRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.Password = password
	//resp, err := client.drdsconn.CreateDrdsAccount(req)
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.CreateDrdsAccount(req)
	})
	resp, _ := raw.(*drds.CreateDrdsAccountResponse)

	return resp, err
}

func convertTypeValue(returnedType string, rawType string) InstanceType {
	var i InstanceType
	returnedInstanceType := InstanceType(returnedType)
	switch returnedInstanceType {
	case PrivateType_:
		i = PrivateType
	}
	return i
}

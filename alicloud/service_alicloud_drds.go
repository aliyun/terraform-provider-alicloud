package alicloud

import (

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

// crate Drdsinstance
func (client *AliyunClient) CreateDrdsInstance(zoneID string,regionId string,description string, quantity int, specification string, payType string,vpcId string, vswitchId string, instanceSeries string) (response *drds.CreateDrdsInstanceResponse, err error) {
	req := drds.CreateCreateDrdsInstanceRequest()
	req.ZoneId = zoneID
	req.PayType = payType
	req.Description = description
	req.Quantity = requests.NewInteger(quantity)
	req.Specification = specification
	req.Type = string(PrivateType_)
	req.VpcId = vpcId
	req.VswitchId = vswitchId
	req.IsHa = requests.NewBoolean(true)
	req.InstanceSeries = instanceSeries
	resp, err := client.drdsconn.CreateDrdsInstance(req)
	if err != nil {
		return nil, fmt.Errorf("create drds error")
	}
	if resp == nil {
		return resp, fmt.Errorf("create drds error")
	}

	return resp, nil
}



func (client *AliyunClient)DescribeDrdsInstance(drdsInstanceId string)(response *drds.DescribeDrdsInstanceResponse, err error)  {
	req := drds.CreateDescribeDrdsInstanceRequest()
	req.DrdsInstanceId = drdsInstanceId
	resp, err := client.drdsconn.DescribeDrdsInstance(req)
	if  err != nil {
		return nil, fmt.Errorf("describe drds instance error")
	}
	if resp == nil {
		return resp, fmt.Errorf("describe drds instance error")
	}
	return resp, nil
}

func (client *AliyunClient)DescribeDrdsInstances(regionId string) (response *drds.DescribeDrdsInstancesResponse, err error) {
	req := drds.CreateDescribeDrdsInstancesRequest()
	req.Type = string(Private)
	resp, err := client.drdsconn.DescribeDrdsInstances(req)
	return resp, err

}


func (client *AliyunClient)DescribeRegions() (response *drds.DescribeRegionsResponse, err error) {
	req := drds.CreateDescribeRegionsRequest()
	resp, err := client.drdsconn.DescribeRegions(req)
	return resp, err
}

func (client *AliyunClient)ModifyDrdsInstanceDescription(request *drds.ModifyDrdsInstanceDescriptionRequest)(response *drds.ModifyDrdsInstanceDescriptionResponse, err error)  {
	req := drds.CreateModifyDrdsInstanceDescriptionRequest()
	req.DrdsInstanceId = request.DrdsInstanceId
	req.Description = request.Description
	resp, err := client.drdsconn.ModifyDrdsInstanceDescription(req)
	return resp, err

}

func (client *AliyunClient)RemoveDrdsInstance(drdsInstanceId string)(response *drds.RemoveDrdsInstanceResponse, err error)  {
	req := drds.CreateRemoveDrdsInstanceRequest()
	req.DrdsInstanceId = drdsInstanceId
	resp, err := client.drdsconn.RemoveDrdsInstance(req)
	return resp,err
}

func (client *AliyunClient)CreateDrdsDB(drdsInstanceId string, dbName string, encode string, password string)(response *drds.CreateDrdsDBResponse, err error) {
	req := drds.CreateCreateDrdsDBRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.Encode = encode
	req.Password = password
	resp, err := client.drdsconn.CreateDrdsDB(req)
	return resp,err
}

func (client *AliyunClient)DescribeDrdsDB(dbName string, drdsInstanceId string) (response *drds.DescribeDrdsDBResponse, err error)  {
	req := drds.CreateDescribeDrdsDBRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	resp,err := client.drdsconn.DescribeDrdsDB(req)
	return resp, err
}

func (client *AliyunClient)DeleteDrdsDB(dbName string, drdsInstanceId string) (response *drds.DeleteDrdsDBResponse, err error){
	req := 	drds.CreateDeleteDrdsDBRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	resp, err := client.drdsconn.DeleteDrdsDB(req)
	return resp,err
}

func (client *AliyunClient)ModifyDrdsDBPasswd(dbName string, drdsInstanceId string, newPasswd string)(response *drds.ModifyDrdsDBPasswdResponse, err error)  {
	req := drds.CreateModifyDrdsDBPasswdRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	resp, err := client.drdsconn.ModifyDrdsDBPasswd(req)
	return resp, err
}

func (client *AliyunClient)DescribeDrdsDBs(drdsInstanceId string) (response *drds.DescribeDrdsDBsResponse, err error) {
	req := drds.CreateDescribeDrdsDBsRequest()
	req.DrdsInstanceId = drdsInstanceId
	resp, err := client.drdsconn.DescribeDrdsDBs(req)
	return resp, err
}

func (client *AliyunClient)DescribeDrdsDBIpWhiteList(drdsInstanceId string, dbName string, groupName string) (response *drds.DescribeDrdsDBIpWhiteListResponse, err error)  {
	req := drds.CreateDescribeDrdsDBIpWhiteListRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	if groupName == "" {
		req.GroupName = "DEFAULT_GROUP"
	} else {
		req.GroupName = groupName
	}
	resp, err := client.drdsconn.DescribeDrdsDBIpWhiteList(req)
	return resp, err
}

func (client *AliyunClient)ModifyDrdsIpWhiteList(drdsInstanceId string, dbName string, ipWhiteList string, mode bool, groupName string) (response *drds.ModifyDrdsIpWhiteListResponse, err error)  {
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
	resp,err := client.drdsconn.ModifyDrdsIpWhiteList(req)
	return resp,err
}

func (client *AliyunClient)DescribeRdsList(drdsInstanceId string, dbName string)(response *drds.DescribeRdsListResponse, err error) {
	req := drds.CreateDescribeRdsListRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	resp,err :=	client.drdsconn.DescribeRdsList(req)
	return resp, err
}

func (client *AliyunClient)ModifyFullTableScan(drdsInstanceId string, dbName string, tableNames string, fulltableScan bool)(response *drds.ModifyFullTableScanResponse, err error) {
	req := drds.CreateModifyFullTableScanRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.TableNames = tableNames
	req.FullTableScan = requests.NewBoolean(fulltableScan)
	resp, err := client.drdsconn.ModifyFullTableScan(req)
	return resp, err
}


func (client *AliyunClient)DescribeShardDBs(drdsInstanceId string, dbName string) (response *drds.DescribeShardDBsResponse, err error) {
	req := drds.CreateDescribeShardDBsRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	resp, err := client.drdsconn.DescribeShardDBs(req)
	return resp, err
}

func (client *AliyunClient)DeleteFailedDrdsDB(drdsInstanceId string, dbName string)(response *drds.DeleteFailedDrdsDBResponse, err error) {
	req := drds.CreateDeleteFailedDrdsDBRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	resp, err := client.drdsconn.DeleteFailedDrdsDB(req)
	return resp, err
}

func (client *AliyunClient)ModifyRdsReadWeight(drdsInstanceId string, dbName string, instanceNames string, weight string)(response *drds.ModifyRdsReadWeightResponse, err error) {
	req := drds.CreateModifyRdsReadWeightRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	resp, err := client.drdsconn.ModifyRdsReadWeight(req)
	return resp, err
}

func (client *AliyunClient)CreateReadOnlyAccount(drdsInstanceId string, dbName string, password string) (response *drds.CreateReadOnlyAccountResponse, err error) {
	req := drds.CreateCreateReadOnlyAccountRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.Password = password
	resp,err := client.drdsconn.CreateReadOnlyAccount(req)
	return resp,err
}

func (client *AliyunClient)DescribeReadOnlyAccount(drdsInstanceId string, dbName string)(response *drds.DescribeReadOnlyAccountResponse, err error) {
	req := drds.CreateDescribeReadOnlyAccountRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	resp, err := client.drdsconn.DescribeReadOnlyAccount(req)
	return resp,err
}

func (client *AliyunClient)ModifyReadOnlyAccountPassword(drdsInstanceId string, dbName string, accountName string, originPassword string, newPassword string)(response *drds.ModifyReadOnlyAccountPasswordResponse, err error) {
	req := drds.CreateModifyReadOnlyAccountPasswordRequest()
	req.DbName = dbName
	req.DrdsInstanceId = drdsInstanceId
	req.AccountName = accountName
	req.OriginPassword = originPassword
	req.NewPasswd = newPassword
	resp, err := client.drdsconn.ModifyReadOnlyAccountPassword(req)
	return resp, err
}

func (client *AliyunClient)RemoveReadOnlyAccount(drdsInstanceId string, dbName string, accountName string)(response *drds.RemoveReadOnlyAccountResponse, err error) {
	req := drds.CreateRemoveReadOnlyAccountRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.AccountName = accountName
	resp,err := client.drdsconn.RemoveReadOnlyAccount(req)
	return resp, err
}

func (client *AliyunClient)CreateDrdsAccount(drdsInstanceId string, dbName string, password string)(response *drds.CreateDrdsAccountResponse, err error)  {
	req := drds.CreateCreateDrdsAccountRequest()
	req.DrdsInstanceId = drdsInstanceId
	req.DbName = dbName
	req.Password = password
	resp, err := client.drdsconn.CreateDrdsAccount(req)
	return resp,err
}

func convertTypeValue(returnedType string, rawType string) InstanceType {
	var i InstanceType
	returnedInstanceType := InstanceType(returnedType)
	switch returnedInstanceType {
	case PrivateType_:
		i = PrivateType
	case PublicType_:
		i = PublicType
	}
	return i
}



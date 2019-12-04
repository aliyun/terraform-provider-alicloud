package market

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// ProductItem is a nested struct in market response
type ProductItem struct {
	Code             string `json:"Code" xml:"Code"`
	Name             string `json:"Name" xml:"Name"`
	CategoryId       int64  `json:"CategoryId" xml:"CategoryId"`
	SupplierId       int64  `json:"SupplierId" xml:"SupplierId"`
	SupplierName     string `json:"SupplierName" xml:"SupplierName"`
	ShortDescription string `json:"ShortDescription" xml:"ShortDescription"`
	Tags             string `json:"Tags" xml:"Tags"`
	SuggestedPrice   string `json:"SuggestedPrice" xml:"SuggestedPrice"`
	TargetUrl        string `json:"TargetUrl" xml:"TargetUrl"`
	ImageUrl         string `json:"ImageUrl" xml:"ImageUrl"`
	Score            string `json:"Score" xml:"Score"`
	OperationSystem  string `json:"OperationSystem" xml:"OperationSystem"`
	WarrantyDate     string `json:"WarrantyDate" xml:"WarrantyDate"`
	DeliveryDate     string `json:"DeliveryDate" xml:"DeliveryDate"`
	DeliveryWay      string `json:"DeliveryWay" xml:"DeliveryWay"`
}

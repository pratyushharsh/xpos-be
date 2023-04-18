package src

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"model"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SyncRepository struct{}

var (
	transactionType = "TRN"
	productType     = "PROD"
	customerType    = "CUS"
	configType      = "CFG"

	taxConfigType      = "TAX"
	invoiceConfigType  = "INV"
	sequenceConfigType = "SEQ"

	taxService = TaxRepository{}
)

func buildWriteRequestForProducts(storeId *string, utcTimeStamp *int64, products *[]*model.ProductEntity, writeRequest []*dynamodb.WriteRequest) ([]*dynamodb.WriteRequest, *ServiceError) {
	for _, product := range *products {
		pk := "STORE#" + *storeId
		sk := "PROD#" + *product.ProductId
		product.LastSyncAt = utcTimeStamp
		// Convert the item to DynamoDB AttributeValues
		dao := model.ProductEntityDao{
			PK:            &pk,
			SK:            &sk,
			GPK1:          storeId,
			GSK1:          utcTimeStamp,
			Type:          &productType,
			ProductEntity: product,
		}

		av, err := dynamodbattribute.MarshalMap(dao)
		if err != nil {
			return nil, &ServiceError{
				ErrorCode:    "500",
				ErrorMessage: "Error Unparsing input request",
			}
		}
		writeRequest = append(writeRequest, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		})
	}
	return writeRequest, nil
}

func buildWriteRequestForCustomers(storeId *string, utcTimeStamp *int64, customers *[]*model.CustomerEntity, writeRequest []*dynamodb.WriteRequest) ([]*dynamodb.WriteRequest, *ServiceError) {
	for _, customer := range *customers {

		pk := "STORE#" + *storeId
		sk := "CUS#" + *customer.ContactId

		customer.LastSyncAt = utcTimeStamp
		// Convert the item to DynamoDB AttributeValues
		dao := model.CustomerEntityDao{
			PK:             &pk,
			SK:             &sk,
			GPK1:           storeId,
			GSK1:           utcTimeStamp,
			Type:           &customerType,
			CustomerEntity: customer,
		}

		av, err := dynamodbattribute.MarshalMap(dao)
		if err != nil {
			return nil, &ServiceError{
				ErrorCode:    "500",
				ErrorMessage: "Error Unparsing input request",
			}
		}
		writeRequest = append(writeRequest, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		})
	}
	return writeRequest, nil
}

func buildWriteRequestForTaxConfig(storeId *string, utcTimeStamp *int64, clientTg *[]*model.TaxGroupEntity, writeRequest []*dynamodb.WriteRequest) ([]*dynamodb.WriteRequest, *ServiceError) {
	pk := "STORE#" + *storeId
	sk := configType + "#" + taxConfigType
	// Convert the item to DynamoDB AttributeValues

	taxGroup := make(map[string]*model.TaxGroupEntity)

	var serverTaxGroup []*model.TaxGroupEntity
	existingTaxGroup := taxService.getTaxGroupsForStore(*storeId)

	if existingTaxGroup != nil && existingTaxGroup.TaxGroups != nil {
		for _, taxGroup := range *existingTaxGroup.TaxGroups {
			serverTaxGroup = append(serverTaxGroup, taxGroup)
		}
	}

	sort.Sort(model.TaxGroupEntitySorter(*clientTg))
	sort.Sort(model.TaxGroupEntitySorter(serverTaxGroup))

	clientTaxGroup := *clientTg
	// Merge the tax rules
	var merged []*model.TaxGroupEntity
	var i, j int
	for i < len(serverTaxGroup) && j < len(clientTaxGroup) {
		if strings.Compare(*serverTaxGroup[i].GroupId, *clientTaxGroup[j].GroupId) < 0 {
			merged = append(merged, serverTaxGroup[i])
			i++
		} else if strings.Compare(*serverTaxGroup[i].GroupId, *clientTaxGroup[j].GroupId) > 0 {
			merged = append(merged, clientTaxGroup[j])
			j++
		} else {
			mergedRule := serverTaxGroup[i].Merge(*clientTaxGroup[j])
			merged = append(merged, &mergedRule)
			i++
			j++
		}
	}

	for i < len(serverTaxGroup) {
		merged = append(merged, serverTaxGroup[i])
		i++
	}

	for j < len(clientTaxGroup) {
		merged = append(merged, clientTaxGroup[j])
		j++
	}

	// Merge the two taxgroups
	for _, taxConfig := range merged {
		taxGroup[*taxConfig.GroupId] = taxConfig
	}

	dao := model.TaxGroupDao{
		PK:        &pk,
		SK:        &sk,
		GPK1:      storeId,
		GSK1:      utcTimeStamp,
		Type:      &taxConfigType,
		TaxGroups: &taxGroup,
	}

	av, err := dynamodbattribute.MarshalMap(dao)
	if err != nil {
		return nil, &ServiceError{
			ErrorCode:    "500",
			ErrorMessage: "Error Unparsing input request",
		}
	}
	writeRequest = append(writeRequest, &dynamodb.WriteRequest{
		PutRequest: &dynamodb.PutRequest{
			Item: av,
		},
	})
	return writeRequest, nil
}

func buildWriteRequestForInvoiceConfig(storeId *string, utcTimeStamp *int64, invoiceConfig *[]*model.ReportConfigEntity, writeRequest []*dynamodb.WriteRequest) ([]*dynamodb.WriteRequest, *ServiceError) {

	for _, invoiceCfg := range *invoiceConfig {
		pk := "STORE#" + *storeId
		sk := configType + "#" + invoiceConfigType
		// Convert the item to DynamoDB AttributeValues

		dao := model.ReportConfigDao{
			PK:                 &pk,
			SK:                 &sk,
			GPK1:               storeId,
			GSK1:               utcTimeStamp,
			Type:               &invoiceConfigType,
			ReportConfigEntity: invoiceCfg,
		}

		av, err := dynamodbattribute.MarshalMap(dao)
		if err != nil {
			return nil, &ServiceError{
				ErrorCode:    "500",
				ErrorMessage: "Error Unparsing input request",
			}
		}
		writeRequest = append(writeRequest, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		})
	}
	return writeRequest, nil
}

func buildWriteRequestForSequenceConfig(storeId *string, utcTimeStamp *int64, sequenceConfigs *[]*model.SequenceEntity, writeRequest []*dynamodb.WriteRequest) ([]*dynamodb.WriteRequest, *ServiceError) {
	pk := "STORE#" + *storeId
	sk := configType + "#" + sequenceConfigType
	// Convert the item to DynamoDB AttributeValues

	sequences := make(map[string]*model.SequenceEntity)

	for _, seqConfig := range *sequenceConfigs {
		sequences[*seqConfig.Name] = seqConfig
	}

	dao := model.SequenceEntityDao{
		PK:              &pk,
		SK:              &sk,
		GPK1:            storeId,
		GSK1:            utcTimeStamp,
		Type:            &sequenceConfigType,
		SequenceEntitys: &sequences,
	}

	av, err := dynamodbattribute.MarshalMap(dao)
	if err != nil {
		return nil, &ServiceError{
			ErrorCode:    "500",
			ErrorMessage: "Error Unparsing input request",
		}
	}
	writeRequest = append(writeRequest, &dynamodb.WriteRequest{
		PutRequest: &dynamodb.PutRequest{
			Item: av,
		},
	})
	return writeRequest, nil
}

func (sr *SyncRepository) UpdateSyncData(storeId string, req *UpdateSyncRequest) (*UpdateSyncResponse, *ServiceError) {
	db := GetDynamoDbClient()

	now := time.Now()
	utcTimeStamp := now.UnixMicro()
	//Convert the item to DynamoDB AttributeValues
	var writeRequest []*dynamodb.WriteRequest

	response := &UpdateSyncResponse{
		LastSyncedAt: &utcTimeStamp,
	}

	// For each transaction request create a write request
	tran, tranErr := BuildWriteRequestForTransaction(&storeId, &utcTimeStamp, req.Transactions)
	response.Transactions = &SyncResponse{
		Data:         tran,
		Error:        tranErr,
		LastSyncedAt: &utcTimeStamp,
	}

	// For each product request create a write request
	writeRequest, err := buildWriteRequestForProducts(&storeId, &utcTimeStamp, req.Products, writeRequest)
	if err != nil {
		return nil, err
	}

	// For each customer request create a write request
	writeRequest, err = buildWriteRequestForCustomers(&storeId, &utcTimeStamp, req.Customers, writeRequest)
	if err != nil {
		return nil, err
	}

	// Create Request for config element.
	if req.Config != nil {
		if req.Config.TaxConfig != nil && len(*req.Config.TaxConfig) > 0 {
			writeRequest, err = buildWriteRequestForTaxConfig(&storeId, &utcTimeStamp, req.Config.TaxConfig, writeRequest)
		}
	}

	if req.Config != nil {
		if req.Config.InvoiceConfig != nil {
			writeRequest, err = buildWriteRequestForInvoiceConfig(&storeId, &utcTimeStamp, req.Config.InvoiceConfig, writeRequest)
		}
	}

	if req.Config != nil {
		if req.Config.SequenceConfig != nil {
			writeRequest, err = buildWriteRequestForSequenceConfig(&storeId, &utcTimeStamp, req.Config.SequenceConfig, writeRequest)
		}
	}

	// Writing Request To The DynamoDB
	for i := 0; i < len(writeRequest); i += 25 {
		end := i + 25
		if end > len(writeRequest) {
			end = len(writeRequest)
		}
		batchWriteRequest := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				DataTable: writeRequest[i:end],
			},
		}
		_, err := db.BatchWriteItem(batchWriteRequest)
		if err != nil {
			log.Printf("Error writing to dynamo: %v", err)
			return nil, &ServiceError{
				ErrorCode:    "500",
				ErrorMessage: "Error Updating Sync Data",
			}
		}
	}

	return response, nil
}

func (sr *SyncRepository) GetSyncData(storeId string, from *int64, to *int64) (*GetSyncResponse, *ServiceError) {
	db := GetDynamoDbClient()

	// @TODO CHECK IF STORE NOT EXIST

	inpReq := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":storeId": {
				S: aws.String(storeId),
			},
		},
		IndexName:              aws.String("GPK1-GSK1-index"),
		KeyConditionExpression: aws.String("GPK1 = :storeId"),
		TableName:              aws.String(DataTable),
	}

	if from != nil && to != nil {
		inpReq.ExpressionAttributeValues[":from"] = &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(*from, 10)),
		}
		inpReq.ExpressionAttributeValues[":to"] = &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(*to, 10)),
		}
		inpReq.KeyConditionExpression = aws.String("GPK1 = :storeId AND GSK1 BETWEEN :from AND :to")
	} else if from != nil {
		inpReq.ExpressionAttributeValues[":from"] = &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(*from, 10))}
		inpReq.KeyConditionExpression = aws.String(*inpReq.KeyConditionExpression + " AND GSK1 >= :from")
	} else if to != nil {
		inpReq.ExpressionAttributeValues[":to"] = &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(*to, 10))}
		inpReq.KeyConditionExpression = aws.String(*inpReq.KeyConditionExpression + " AND GSK1 < :to")
	}

	res, err := db.Query(inpReq)

	// Error in input request
	if err != nil {
		return nil, &ServiceError{
			ErrorCode:    "500",
			ErrorMessage: err.Error(),
		}
	}

	// For each item parse check the type and add to the response
	var products []*model.ProductEntity
	var customers []*model.CustomerEntity
	var transactions []*model.TransactionHeaderEntity
	var taxGroups []*model.TaxGroupEntity
	var invoiceConfig []*model.ReportConfigEntity
	var sequenceConfig []*model.SequenceEntity

	for _, item := range res.Items {
		typ := item["Type"].S
		if *typ == productType {
			var tp model.ProductEntity
			_ = dynamodbattribute.UnmarshalMap(item, &tp)
			products = append(products, &tp)
		} else if *typ == customerType {
			var tp model.CustomerEntity
			_ = dynamodbattribute.UnmarshalMap(item, &tp)
			customers = append(customers, &tp)
		} else if *typ == transactionType {
			var tp model.TransactionHeaderEntity
			_ = dynamodbattribute.UnmarshalMap(item, &tp)
			transactions = append(transactions, &tp)
		} else if *typ == taxConfigType {
			var tp *model.TaxGroupDao
			_ = dynamodbattribute.UnmarshalMap(item, &tp)

			if tp != nil && tp.TaxGroups != nil {
				for _, taxGroup := range *tp.TaxGroups {
					taxGroups = append(taxGroups, taxGroup)
				}
			}
		} else if *typ == invoiceConfigType {
			var tp *model.ReportConfigEntity
			_ = dynamodbattribute.UnmarshalMap(item, &tp)
			invoiceConfig = append(invoiceConfig, tp)
		} else if *typ == sequenceConfigType {
			var tp *model.SequenceEntityDao
			_ = dynamodbattribute.UnmarshalMap(item, &tp)

			if tp != nil && tp.SequenceEntitys != nil {
				for _, sequence := range *tp.SequenceEntitys {
					sequenceConfig = append(sequenceConfig, sequence)
				}
			}
		}
	}

	return &GetSyncResponse{
		Transactions: SyncData{
			Data: &transactions,
			From: from,
			To:   to,
		},
		Customers: SyncData{
			Data: &customers,
			From: from,
			To:   to,
		},
		Products: SyncData{
			Data: &products,
			From: from,
			To:   to,
		},
		Config: ConfigOutput{
			TaxConfig: SyncData{
				Data: &taxGroups,
				From: from,
				To:   to,
			},
			InvoiceConfig: SyncData{
				Data: &invoiceConfig,
				From: from,
				To:   to,
			},
			SequenceConfig: SyncData{
				Data: &sequenceConfig,
				From: from,
				To:   to,
			},
		},
	}, nil
}

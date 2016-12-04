package main

import (
	"encoding/base64"
	"strconv"
	"errors"
	"fmt"
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
//	"github.com/op/go-logging"
)

type AssetManagementChaincode struct {
}

// Init method will be called during deployment.
// The deploy transaction metadata is supposed to contain the administrator cert
func (t *AssetManagementChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Init Chaincode...")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	// Create invoice table
	err := stub.CreateTable("Invoice", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Number", Type: shim.ColumnDefinition_INT32, Key: true},
		&shim.ColumnDefinition{Name: "Price", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "Status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "DeliveryDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "RequestDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "PaymentDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "SupplierId", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "BuyerId", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "SupplierCert", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "BuyerCert", Type: shim.ColumnDefinition_BYTES, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating Invoice table.")
	}

	// Create PaymentRequest table
    err = stub.CreateTable("PaymentRequest", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Id", Type: shim.ColumnDefinition_INT32, Key: true},
		&shim.ColumnDefinition{Name: "Invoice", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "DiscountRate", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "PayerId", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "PayerCert", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "Status", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating PaymentRequest table.")
	}


	supplierRole, err := stub.GetCallerMetadata()
	fmt.Printf("Assiger role is %v\n", string(supplierRole))

	if err != nil {
		return nil, fmt.Errorf("Failed getting metadata, [%v]", err)
	}

	if len(supplierRole) == 0 {
		return nil, errors.New("Invalid supplier role. Empty.")
	}

	stub.PutState("supplierRole", supplierRole)

	fmt.Println("Init Chaincode...done")

	return nil, nil
}

func errorJson(function string, errorMsg error) []byte {
	error := `{"error": { "code":-1,"message":"` + function + `","data":"` + errorMsg.Error() + `"}}`;
	return []byte(error);
}

func (t *AssetManagementChaincode) createInvoice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Create invoice...")

	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 7")
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		throwError := errors.New("Expecting integer value for invoice number")
		return errorJson("createInvoice", throwError), throwError
	}

	price, err := strconv.Atoi(args[1])
	if err != nil {
		throwError := errors.New("Expecting integer value for invoice price")
		return errorJson("createInvoice", throwError), throwError
	}

	deliveryDate := args[2]

	supplierId, err := strconv.Atoi(args[3])
	if err != nil {
		throwError := errors.New("Expecting integer value for invoice supplierId")
		return errorJson("createInvoice", throwError), throwError
	}

	buyerId, err := strconv.Atoi(args[4])
	if err != nil {
		throwError := errors.New("Expecting integer value for invoice buyerId")
		return errorJson("createInvoice", throwError), throwError
	}
    fmt.Println("Invoice number = ", number)	

	supplier, err := base64.StdEncoding.DecodeString(args[5])
	if err != nil {
		return nil, errors.New("Failed decodinf supplier")
	}
	fmt.Println("Supplier cert bytes = ", supplier)	

	buyer, err := base64.StdEncoding.DecodeString(args[6])
	if err != nil {
		return nil, errors.New("Failed decodinf buyer")
	}
	fmt.Println("Buyer cert bytes = ", buyer)	

//************************************************************************
	//Enable this when membersrvc conf available

	// Verify the identity of the caller
	// Only a supplier can create an invoice
	// supplierRole, err := stub.GetState("supplierRole")

	// if err != nil {
	// 	return nil, errors.New("Failed fetching supplier identity")
	// }


	// callerRole, err := stub.ReadCertAttribute("role")
	// if err != nil {
	// 	fmt.Printf("Error reading attribute 'role' [%v] \n", err)
	// 	return nil, fmt.Errorf("Failed fetching caller role. Error was [%v]", err)
	// }

	// caller := string(callerRole[:])
	// supplierStr := string(supplierRole[:])

	// if caller != supplierStr {
	// 	fmt.Printf("Caller is not supplier - caller %v assigner %v\n", caller, supplierStr)
	// 	return nil, fmt.Errorf("The caller does not have the rights to invoke create invoice. Expected role [%v], caller role [%v]", supplierStr, caller)
	// }
//***********************************************************************8	

	// Create an invoice
	fmt.Println("Creating new invoice, number: [%s] ,price: [%s], deliveryDate: [%s], supplier is [% x], buyer is [% x]",number,price,deliveryDate, supplier,buyer)

	ok, err := stub.InsertRow("Invoice", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(number)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(price)}},
			&shim.Column{Value: &shim.Column_String_{String_: "Pending"}},
			&shim.Column{Value: &shim.Column_String_{String_: deliveryDate}},
			&shim.Column{Value: &shim.Column_String_{String_: deliveryDate}},
			&shim.Column{Value: &shim.Column_String_{String_: deliveryDate}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(supplierId)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(buyerId)}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: supplier}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: buyer}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Invoice with this number was already created.")
	}

	fmt.Println("Create invoice...done!")

	return nil, err
}

func (t *AssetManagementChaincode) approveInvoice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Approve invoice...")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		throwError := errors.New("Expecting integer value for invoice number")
		return errorJson("approveInvoice", throwError), throwError
	}
	
	buyer, err := base64.StdEncoding.DecodeString(args[1])
	if err != nil {
		return nil, errors.New("Failed decodinf buyer")
	}
	fmt.Println("Buyer cert bytes = ", buyer)	



	// Verify the identity of the caller
	// Only the owner can transfer one of his assets
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_Int32{Int32: int32(number)}}
	columns = append(columns, col1)

	row, err := stub.GetRow("Invoice", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving invoice [%s]: [%s]", number, err)
	}

	realBuyer := row.Columns[9].GetBytes()
	fmt.Println("Real buyer of [%s] is [% x]", number, realBuyer)
	if len(realBuyer) == 0 {
		return nil, fmt.Errorf("Invalid real buyer. Nil")
	}

	ok := bytes.Equal(buyer, realBuyer)

	if ok != true {
		return nil, fmt.Errorf("Caller is not allowed to do this operation")	
	}

	// Approve an invoice
	fmt.Println("Approving the invoice, number: [%s] , buyer is [% x]",number,buyer)

	supplier := row.Columns[8].GetBytes()
	buyerId := row.Columns[7].GetInt32()
	supplierId := row.Columns[6].GetInt32()
	deliveryDate := row.Columns[3].GetString_()
	requestDate := row.Columns[4].GetString_()
	paymentDate := row.Columns[5].GetString_()
	price := row.Columns[1].GetInt32()

	// update from balance
	_, err = stub.ReplaceRow(
		"Invoice",
		shim.Row{
			Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(number)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(price)}},
			&shim.Column{Value: &shim.Column_String_{String_: "Approved"}},
			&shim.Column{Value: &shim.Column_String_{String_: deliveryDate}},
			&shim.Column{Value: &shim.Column_String_{String_: requestDate}},
			&shim.Column{Value: &shim.Column_String_{String_: paymentDate}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(supplierId)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(buyerId)}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: supplier}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: realBuyer}}},
			},
		)
	if err != nil {
		throwError := errors.New("Failed update status row.")
		return errorJson("transfer", throwError), throwError
	}

	
	fmt.Println("Approve invoice...done!")

	return nil, err
}

func (t *AssetManagementChaincode) createPaymentRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Create a payment request...")

	// &shim.ColumnDefinition{Name: "Id", Type: shim.ColumnDefinition_INT32, Key: true},
	// 	&shim.ColumnDefinition{Name: "Invoice", Type: shim.ColumnDefinition_INT32, Key: false},
	// 	&shim.ColumnDefinition{Name: "DiscountRate", Type: shim.ColumnDefinition_INT32, Key: false},
	// 	&shim.ColumnDefinition{Name: "PayerId", Type: shim.ColumnDefinition_INT32, Key: false},
	// 	&shim.ColumnDefinition{Name: "PayerCert", Type: shim.ColumnDefinition_BYTES, Key: false},


	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	payment, err := strconv.Atoi(args[0])
	if err != nil {
		throwError := errors.New("Expecting integer value for payment request id")
		return errorJson("createPaymentRequest", throwError), throwError
	}
	number, err := strconv.Atoi(args[1])
	if err != nil {
		throwError := errors.New("Expecting integer value for invoice number")
		return errorJson("createPaymentRequest", throwError), throwError
	}
	discountRate, err := strconv.Atoi(args[2])
	if err != nil {
		throwError := errors.New("Expecting integer value for discountRate")
		return errorJson("createPaymentRequest", throwError), throwError
	}

	requestDate := args[3]
	
	fmt.Println("Payment request id = ", number)
    fmt.Println("Invoice number = ", number)	

	buyer, err := base64.StdEncoding.DecodeString(args[4])
	if err != nil {
		return nil, errors.New("Failed decodinf buyer")
	}
	fmt.Println("Buyer cert bytes = ", buyer)	

	// Verify the identity of the caller
	// Only the owner can transfer one of his assets
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_Int32{Int32: int32(number)}}
	columns = append(columns, col1)

	row, err := stub.GetRow("Invoice", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving invoice [%s]: [%s]", number, err)
	}

	realBuyer := row.Columns[9].GetBytes()
	fmt.Println("Real buyer of [%s] is [% x]", number, realBuyer)
	if len(realBuyer) == 0 {
		return nil, fmt.Errorf("Invalid real buyer. Nil")
	}

	ok := bytes.Equal(buyer, realBuyer)

	if ok != true {
		return nil, fmt.Errorf("Caller is not allowed to do this operation")	
	}

	// Create a payment request
	fmt.Println("Creating new payment request, number: [%s] ,paymentID: [%s], discountRate: [%s], buyer is [% x]",number,payment,discountRate, buyer)

	ok, err = stub.InsertRow("PaymentRequest", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(payment)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(number)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(discountRate)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: -1}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: buyer}},
			&shim.Column{Value: &shim.Column_String_{String_: "Pending"}},
	}})

	if !ok && err == nil {
		return nil, errors.New("payment request with this id was already created.")
	}

	//Update invoice request date

	supplier := row.Columns[8].GetBytes()
	buyerId := row.Columns[7].GetInt32()
	supplierId := row.Columns[6].GetInt32()
	deliveryDate := row.Columns[3].GetString_()
	paymentDate := row.Columns[5].GetString_()
	price := row.Columns[1].GetInt32()
	status := row.Columns[2].GetString_()

	// update from balance
	_, err = stub.ReplaceRow(
		"Invoice",
		shim.Row{
			Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(number)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(price)}},
			&shim.Column{Value: &shim.Column_String_{String_: status}},
			&shim.Column{Value: &shim.Column_String_{String_: deliveryDate}},
			&shim.Column{Value: &shim.Column_String_{String_: requestDate}},
			&shim.Column{Value: &shim.Column_String_{String_: paymentDate}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(supplierId)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(buyerId)}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: supplier}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: realBuyer}}},
			},
		)
	if err != nil {
		throwError := errors.New("Failed update status row.")
		return errorJson("transfer", throwError), throwError
	}

	fmt.Println("Create payment request...done!")

	return nil, err
}

func (t *AssetManagementChaincode) assignPaymentRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Assign a payment request...")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	payment, err := strconv.Atoi(args[0])
	if err != nil {
		throwError := errors.New("Expecting integer value for payment request id")
		return errorJson("assignPaymentRequest", throwError), throwError
	}
	payerId, err := strconv.Atoi(args[1])
	if err != nil {
		throwError := errors.New("Expecting integer value for payer id")
		return errorJson("assignPaymentRequest", throwError), throwError
	}
	
	
	fmt.Println("Payment request id = ", payment)
    fmt.Println("Payer id = ", payerId)	

	payer, err := base64.StdEncoding.DecodeString(args[2])
	if err != nil {
		return nil, errors.New("Failed decodinf payer")
	}
	fmt.Println("Payer cert bytes = ", payer)	

	// Verify the identity of the caller
	// Only the owner can transfer one of his assets
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_Int32{Int32: int32(payment)}}
	columns = append(columns, col1)

	row, err := stub.GetRow("PaymentRequest", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving payment request [%s]: [%s]", payment, err)
	}

	oldPayerId := row.Columns[3].GetInt32()
	fmt.Println("Real payer of [%s] is [% x]", payment, oldPayerId)
	if oldPayerId != -1 {
		return nil, fmt.Errorf("Payment request [%s] already has payer with id = [%s]",payment,oldPayerId)
	}

	invoice := row.Columns[1].GetInt32()
	discountRate := row.Columns[2].GetInt32()

	// Assign a payment request
	fmt.Println("Assigning a payment request, paymentID: [%s], payerId: [%s], payer is [% x]",payment,payerId, payer)

	// update from balance
	_, err = stub.ReplaceRow(
		"PaymentRequest",
		shim.Row{
			Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(payment)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(invoice)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(discountRate)}},
			&shim.Column{Value: &shim.Column_Int32{Int32: int32(payerId)}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: payer}},
			&shim.Column{Value: &shim.Column_String_{String_: "Assigned"}},
			},
		})
	if err != nil {
		throwError := errors.New("Failed update payer row.")
		return errorJson("transfer", throwError), throwError
	}

	fmt.Println("Create payment request...done!")

	return nil, err
}


func (t *AssetManagementChaincode) isCaller(stub shim.ChaincodeStubInterface, certificate []byte) (bool, error) {
	fmt.Println("Check caller...")

	// In order to enforce access control, we require that the
	// metadata contains the signature under the signing key corresponding
	// to the verification key inside certificate of
	// the payload of the transaction (namely, function name and args) and
	// the transaction binding (to avoid copying attacks)

	// Verify \sigma=Sign(certificate.sk, tx.Payload||tx.Binding) against certificate.vk
	// \sigma is in the metadata

	sigma, err := stub.GetCallerMetadata()
	if err != nil {
		return false, errors.New("Failed getting metadata")
	}


	ok := bytes.Equal(sigma, certificate)

	if ok != true {
		return false, errors.New("Caller is not allowed to do this operation")	
	}
//************************************************************************
	//Enable this when signature is available

	// payload, err := stub.GetPayload()
	// if err != nil {
	// 	return false, errors.New("Failed getting payload")
	// }
	// binding, err := stub.GetBinding()
	// if err != nil {
	// 	return false, errors.New("Failed getting binding")
	// }

	// fmt.Println("passed certificate [% x]", certificate)
	// fmt.Println("passed sigma [% x]", sigma)
	// fmt.Println("passed payload [% x]", payload)
	// fmt.Println("passed binding [% x]", binding)

// 	ok, err := stub.VerifySignature(
// 		certificate,
// 		sigma,
// 		append(payload, binding...),
// 	)
// 	if err != nil {
// //		myLogger.Errorf("Failed checking signature [%s]", err)
// 		return ok, err
// 	}
// 	if !ok {
// //		myLogger.Error("Invalid signature")
// 	}

	fmt.Println("Check caller...Verified!")

	return ok, err
//************************************************************************
}

// Invoke will be called for every transaction.
// Supported functions are the following:
// "assign(asset, owner)": to assign ownership of assets. An asset can be owned by a single entity.
// Only an administrator can call this function.
// "transfer(asset, newOwner)": to transfer the ownership of an asset. Only the owner of the specific
// asset can call this function.
// An asset is any string to identify it. An owner is representated by one of his ECert/TCert.
func (t *AssetManagementChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "createInvoice" {
		// Assign ownership
		return t.createInvoice(stub, args)
	} else if function == "approveInvoice" {
		// Transfer ownership
		return t.approveInvoice(stub, args)
	} else if function == "createPaymentRequest" {
		// Transfer ownership
		return t.createPaymentRequest(stub, args)
	} else if function == "assignPaymentRequest" {
		// Transfer ownership
		return t.assignPaymentRequest(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}



func (t *AssetManagementChaincode) invoice_info(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Query invoice...")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		throwError := errors.New("Expecting integer value for invoice number")
		return errorJson("invoice_info", throwError), throwError
	}
	
	buyer, err := base64.StdEncoding.DecodeString(args[1])
	if err != nil {
		return nil, errors.New("Failed decodinf buyer")
	}
	fmt.Println("Buyer cert bytes = ", buyer)	



	// Verify the identity of the caller
	// Only the owner can transfer one of his assets
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_Int32{Int32: int32(number)}}
	columns = append(columns, col1)

	row, err := stub.GetRow("Invoice", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving invoice [%s]: [%s]", number, err)
	}

	realBuyer := row.Columns[9].GetBytes()
	fmt.Println("Real buyer of [%s] is [% x]", number, realBuyer)
	if len(realBuyer) == 0 {
		return nil, fmt.Errorf("Invalid real buyer. Nil")
	}

	ok := bytes.Equal(buyer, realBuyer)

	if ok != true {
		return nil, fmt.Errorf("Caller is not allowed to do this operation")	
	}

	
	deliveryDate := row.Columns[3].GetString_()
	requestDate := row.Columns[4].GetString_()
	paymentDate := row.Columns[5].GetString_()
	status := row.Columns[2].GetString_()
	price := row.Columns[1].GetInt32()

	jsonResp := `{"invoice":"` + strconv.Itoa(int(number)) + `","price":"` + strconv.Itoa(int(price)) + `",` +
		`"delivery_date":"` + deliveryDate + `","request_date":"` + requestDate +
		`","payment_date":"` + paymentDate + `","status":"` + status +`"}`
	
	fmt.Println(jsonResp)
	fmt.Println("Query invoice...done!")

	return []byte(jsonResp), nil
}

func (t *AssetManagementChaincode) payment_info(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Query a payment request...")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	payment, err := strconv.Atoi(args[0])
	if err != nil {
		throwError := errors.New("Expecting integer value for payment id")
		return errorJson("payment_info", throwError), throwError
	}
	payerId, err :=strconv.Atoi(args[1])
	if err != nil {
		throwError := errors.New("Expecting integer value for payer id")
		return errorJson("payment_info", throwError), throwError
	}
	
	
	fmt.Println("Payment request id = ", payment)
    fmt.Println("Payer id = ", payerId)	

	payer, err := base64.StdEncoding.DecodeString(args[2])
	if err != nil {
		return nil, errors.New("Failed decodinf payer")
	}
	fmt.Println("Payer cert bytes = ", payer)	

	// Verify the identity of the caller
	// Only the owner can transfer one of his assets
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_Int32{Int32: int32(payment)}}
	columns = append(columns, col1)

	row, err := stub.GetRow("PaymentRequest", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving payment request [%s]: [%s]", payment, err)
	}

	oldPayerId := row.Columns[3].GetInt32()

	if int32(oldPayerId) != int32(payerId) {
				realPayer := row.Columns[4].GetBytes()
				fmt.Println("Real payer of [%s] is [% x]", payment, realPayer)
				if len(realPayer) == 0 {
					return nil, fmt.Errorf("Invalid real payer. Nil")
				}

				ok := bytes.Equal(payer, realPayer)

				if ok != true {
					return nil, fmt.Errorf("Payment request already has payer")	
				}
	}

	invoice := row.Columns[1].GetInt32()
	discountRate := row.Columns[2].GetInt32()
	status := row.Columns[5].GetString_()


	jsonResp := `{"invoice":"` + strconv.Itoa(int(invoice)) + `","paymentId":"` + strconv.Itoa(int(payment)) + `",` +
		`"discountRate":"` + strconv.Itoa(int(discountRate)) + `"status":"` + status +`"}`
	
	fmt.Println(jsonResp)
	fmt.Println("Query payment request...done!")

	return []byte(jsonResp), nil
}

// Query callback representing the query of a chaincode
// Supported functions are the following:
// "query(asset)": returns the owner of the asset.
// Anyone can invoke this function.
func (t *AssetManagementChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Query [%s]", function)

	if function == "invoice_info" {
		// Get invoice_info
		return t.invoice_info(stub, args)
	} else if function == "payment_info" {
		// Get payment_info
		return t.payment_info(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(AssetManagementChaincode))
	if err != nil {
		fmt.Printf("Error starting AssetManagementChaincode: %s", err)
	}
}

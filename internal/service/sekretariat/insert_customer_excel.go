package sekretariat

import (
	"context"
	"log"
	"sekretariat/internal/entity/sekretariat"
	"sekretariat/pkg/errors"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// ImportCustomersFromExcel reads customers from an Excel file and inserts them into the database.
func (s Service) ImportCustomersFromExcel(ctx context.Context) error {

	filePath := `C:\Users\user\go\src\sekretariat\internal\service\sekretariat\CustomerPBM.xlsx`
	companyID := 2
	// Open  the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][ImportCustomersFromExcel][OPEN_FILE]")
	}
	defer f.Close()

	// Get all sheet names
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return errors.New("[SERVICE][ImportCustomersFromExcel][NO_SHEETS_FOUND]")
	}

	// Read the first sheet
	rows, err := f.GetRows(sheetNames[1])
	if err != nil {
		return errors.Wrap(err, "[SERVICE][ImportCustomersFromExcel][GET_ROWS]")
	}

	log.Println(len(rows))

	// Skip header row
	if len(rows) == 0 {
		return errors.New("[SERVICE][ImportCustomersFromExcel][NO_ROWS_FOUND]")
	}

	for i, row := range rows[1:] { // Assuming the first row is the header

		log.Println("data number : ", i)
		if len(row) < 2 {
			continue // Skip invalid rows
		}

		// Assuming columns: ID, Name, Email, etc.
		customer := sekretariat.Customer{
			CustomerID:    row[0],
			CustomerName:  row[1],
			Address:       row[12],
			PIC:           row[13],
			PenandaTangan: row[10],
			Jabatan:       row[11],
			NoTelp:        row[14],
			UpdatedBy:     "Admin",
			CompanyID:     companyID,
		}

		// Insert the customer into the database
		err := s.CreateCustomer(ctx, customer)
		if err != nil {
			return errors.Wrap(err, "[SERVICE][ImportCustomersFromExcel][CREATE_CUSTOMER]")
		}
	}

	return nil
}

// ImportCustomersFromExcel reads customers from an Excel file and inserts them into the database.
func (s Service) ImportContractsFromExcel(ctx context.Context) (interface{}, error) {
	filePath := `C:\Users\user\go\src\sekretariat\internal\service\sekretariat\CustomerPBM.xlsx`
	companyID := 2

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return filePath, errors.Wrap(err, "[SERVICE][ImportContractsFromExcel][OPEN_FILE]")
	}
	defer f.Close()

	// Get all sheet names
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return filePath, errors.New("[SERVICE][ImportContractsFromExcel][NO_SHEETS_FOUND]")
	}

	// Read the first sheet (adjust index as needed)
	rows, err := f.GetRows(sheetNames[1])
	if err != nil {
		return filePath, errors.Wrap(err, "[SERVICE][ImportContractsFromExcel][GET_ROWS]")
	}

	log.Println("Total rows:", len(rows))

	// Define range

	groupedHeaders := make(map[string]sekretariat.KontrakHeader)
	groupedDetails := make(map[string][]sekretariat.KontrakDetail)

	// Iterate through rows to collect details
	for _, row := range rows[1:] {
		if len(row) < 19 {
			log.Printf("[SERVICE][ImportContractsFromExcel][ROW_%d][INSUFFICIENT_DATA]", len(row))
			continue
		}

		contractNumber := row[2]
		qty, err := strconv.Atoi(row[1])
		if err != nil {
			return filePath, errors.Wrap(err, "[SERVICE][ImportContractsFromExcel][CONVERT_QUANTITY]")
		}
		price, err := strconv.ParseFloat(row[6], 64) // Ensure column index is correct
		if err != nil {
			return filePath, errors.Wrap(err, "[SERVICE][ImportContractsFromExcel][CONVERT_PRICE]")
		}

		detail := sekretariat.KontrakDetail{
			NoKontrak:          contractNumber,
			Quantity:           qty,
			Tipe:               row[5],
			Speed:              row[17],
			Harga:              price,
			FreeCopy:           row[7] + " copy",
			OverCopy:           row[8] + "/copy",
			FreeCopyColor:      "",
			OverCopyColor:      "",
			PeriodeAwalString:  row[3],
			PeriodeAkhirString: row[4],
			Penempatan:         row[11],
			ActiveYN:           "Y",
			UpdatedBy:          "Admin",
		}

		groupedDetails[contractNumber] = append(groupedDetails[contractNumber], detail)
	}

	// Iterate through rows again to collect headers
	for _, row := range rows[1:] {
		if len(row) < 19 {
			continue
		}

		contractNumber := row[2]

		// Check if header already exists
		if _, exists := groupedHeaders[contractNumber]; exists {
			continue
		}

		customerID, err := s.data.GetCustomerIDByNameAndAddress(ctx, row[0], row[11])
		if err != nil {
			return filePath, errors.Wrap(err, "[SERVICE][ImportContractsFromExcel][GET_CUSTOMER_ID]")
		}

		bankID, err := strconv.Atoi(row[14])
		if err != nil {
			return filePath, errors.Wrap(err, "[SERVICE][ImportContractsFromExcel][CONVERT_BANK_ID]")
		}
		paymentID, err := strconv.Atoi(row[18])
		if err != nil {
			return filePath, errors.Wrap(err, "[SERVICE][ImportContractsFromExcel][CONVERT_PAYMENT_ID]")
		}
		deposit, err := strconv.ParseFloat(row[15], 64)
		if err != nil {
			return filePath, errors.Wrap(err, "[SERVICE][ImportContractsFromExcel][CONVERT_DEPOSIT]")
		}

		header := sekretariat.KontrakHeader{
			NoKontrak:   contractNumber,
			TanggalBuat: row[16],
			CompanyID:   companyID,
			Customer: sekretariat.Customer{
				CustomerID: customerID,
			},
			Bank: sekretariat.Bank{
				ID: bankID,
			},
			Pembayaran: sekretariat.Pembayaran{
				ID: paymentID,
			},
			Deposit:           deposit,
			DendaSatuPersenYN: "Y",
			ActiveYN:          "Y",
			UpdatedBy:         "Admin",
			Details:           groupedDetails[contractNumber], // Associate details with the contract header
		}

		groupedHeaders[contractNumber] = header
	}

	// Insert headers and details into the database
	for _, header := range groupedHeaders {
		// Check if contract already exists
		existingContract, err := s.data.GetContractByNumber(ctx, header.NoKontrak)
		log.Println("existingContract", existingContract)
		if err != nil {
			return filePath, errors.Wrap(err, "[SERVICE][ImportContractsFromExcel][CHECK_CONTRACT_EXISTS]")
		}

		if existingContract != false {
			log.Printf("Contract %s already exists, skipping insert.", header.NoKontrak)
			continue
		}

		// err = s.CreateContract(ctx, header)
		// if err != nil {
		// 	return filePath, errors.Wrap(err, "[SERVICE][ImportContractsFromExcel][CREATE_CONTRACT]")
		// }
	}

	return groupedHeaders, nil
}

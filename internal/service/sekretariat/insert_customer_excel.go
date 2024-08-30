package sekretariat

import (
	"context"
)

// ImportCustomersFromExcel reads customers from an Excel file and inserts them into the database.
func (s Service) ImportCustomersFromExcel(ctx context.Context) error {

	// filePath := `C:\Users\user\go\src\sekretariat\internal\service\sekretariat\CustomerPBM.xlsx`
	// companyID := 2
	// // Open  the Excel file
	// f, err := excelize.OpenFile(filePath)
	// if err != nil {
	// 	return errors.Wrap(err, "[SERVICE][ImportCustomersFromExcel][OPEN_FILE]")
	// }
	// defer f.Close()

	// // Get all sheet names
	// sheetNames := f.GetSheetList()
	// if len(sheetNames) == 0 {
	// 	return errors.New("[SERVICE][ImportCustomersFromExcel][NO_SHEETS_FOUND]")
	// }

	// // Read the first sheet
	// rows, err := f.GetRows(sheetNames[0])
	// if err != nil {
	// 	return errors.Wrap(err, "[SERVICE][ImportCustomersFromExcel][GET_ROWS]")
	// }

	// log.Println(len(rows))

	// // Skip header row
	// if len(rows) == 0 {
	// 	return errors.New("[SERVICE][ImportCustomersFromExcel][NO_ROWS_FOUND]")
	// }

	// for i, row := range rows[292:] { // Assuming the first row is the header

	// 	log.Println("data number : ", i)
	// 	if len(row) < 2 {
	// 		continue // Skip invalid rows
	// 	}

	// 	// Assuming columns: ID, Name, Email, etc.
	// 	customer := sekretariat.Customer{
	// 		CustomerName:  row[0],
	// 		Address:       row[11],
	// 		PIC:           row[12],
	// 		PenandaTangan: row[9],
	// 		Jabatan:       row[10],
	// 		NoTelp:        row[13],
	// 		UpdatedBy:     "Admin",
	// 		CompanyID:     companyID,
	// 	}

	// 	// Insert the customer into the database
	// 	err := s.CreateCustomer(ctx, customer)
	// 	if err != nil {
	// 		return errors.Wrap(err, "[SERVICE][ImportCustomersFromExcel][CREATE_CUSTOMER]")
	// 	}
	// }

	return nil
}

package sekretariat

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type (
	// Data ...
	Data struct {
		db   *sqlx.DB
		stmt map[string]*sqlx.Stmt
	}

	// Key - value map for query selector
	statement struct {
		key   string
		query string
	}
)

// Query constants
const (
	// Contracts
	getAllContractsHeader  = "GetAllContractsHeader"
	qGetAllContractsHeader = `SELECT no_kontrak, tanggal_buat, kh.company_id, kh.id_customer,
							cu.nama_customer, cu.alamat, cu.pic, cu.penandatangan, cu.jabatan, cu.no_telp,
							ba.bank_id, ba.bank_name, ba.nomor_rekening, ba.atas_nama,
							active_yn, kh.updated_by, kh.updated_at  FROM kontrak_header kh
							LEFT JOIN customer cu ON kh.id_customer = cu.id_customer
							LEFT JOIN bank ba ON kh.bank_id = ba.bank_id
							WHERE kh.company_id = ?
							AND
	 							((no_kontrak LIKE ? OR ? = '')
	 						OR
	 							(cu.nama_customer LIKE ? OR ? = ''))`

	getAllContractsHeaderPage  = "GetAllContractsHeaderPage"
	qGetAllContractsHeaderPage = `SELECT no_kontrak, tanggal_buat, kh.company_id, kh.id_customer,
								cu.nama_customer, cu.alamat, cu.pic, cu.penandatangan, cu.jabatan, cu.no_telp,
								ba.bank_id, ba.bank_name, ba.nomor_rekening, ba.atas_nama,
								active_yn, kh.updated_by, kh.updated_at  FROM kontrak_header kh
								LEFT JOIN customer cu ON kh.id_customer = cu.id_customer
								LEFT JOIN bank ba ON kh.bank_id = ba.bank_id
								WHERE kh.company_id = ?
								AND
									((no_kontrak LIKE ? OR ? = '')
								OR
									(cu.nama_customer LIKE ? OR ? = '')) LIMIT ?,?`

	getAllContractsHeaderCount  = "GetAllContractsHeaderCount"
	qGetAllContractsHeaderCount = `SELECT count(*) FROM kontrak_header kh
								LEFT JOIN customer cu ON kh.id_customer = cu.id_customer
								LEFT JOIN bank ba ON kh.bank_id = ba.bank_id WHERE kh.company_id = ?`

	getContractHeaderByContractNumber  = "GetContractHeaderByContractNumber"
	qGetContractHeaderByContractNumber = `SELECT no_kontrak, tanggal_buat, kh.company_id, kh.id_customer,
											cu.nama_customer, cu.alamat, cu.pic, cu.penandatangan, cu.jabatan, cu.no_telp,
											ba.bank_id, ba.bank_name, ba.nomor_rekening, ba.atas_nama,
											active_yn, kh.updated_by, kh.updated_at  FROM sys.kontrak_header kh
											LEFT JOIN sys.customer cu ON kh.id_customer = cu.id_customer
											LEFT JOIN sys.bank ba ON kh.bank_id = ba.bank_id
											WHERE kh.company_id = ? AND no_kontrak LIKE ?`

	getContractDetailsByContractNumber  = `GetContractDetailsByContractNumber`
	qGetContractDetailsByContractNumber = `SELECT no_kontrak, quantity, tipe_mesin, speed, harga_sewa, 
										free_copy, over_copy, free_copy_color, over_copy_color, periode_awal, periode_akhir, penempatan,
										active_yn, updated_by, updated_at FROM kontrak_detail WHERE no_kontrak LIKE ?`

	createContractHeader  = "CreateContractHeader"
	qCreateContractHeader = `INSERT INTO sys.kontrak_header(no_kontrak, tanggal_buat, company_id, id_customer, bank_id, active_yn, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?)`

	createContractDetail  = "CreateContractDetail"
	qCreateContractDetail = `INSERT INTO sys.kontrak_detail(no_kontrak, quantity, tipe_mesin, speed,
							harga_sewa, free_copy, over_copy, free_copy_color, over_copy_color, periode_awal, periode_akhir, penempatan, active_yn, updated_by)
							VALUES (?, ?, ?, ?,
							?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Customers
	getAllCustomers  = "GetAllCustomers"
	qGetAllCustomers = `SELECT * FROM customer WHERE company_id = ?`

	getCustomerFiltered  = "GetCustomerFiltered"
	qGetCustomerFiltered = `SELECT * FROM customer WHERE company_id = ? 
						AND
							((id_customer LIKE ? OR ? = '')
	 					OR
	 						(nama_customer LIKE ? OR ? = '')) LIMIT ?,?`

	getCustomerFilteredCount  = "GetCustomerFilteredCount"
	qGetCustomerFilteredCount = `SELECT COUNT(*) FROM customer WHERE company_id = ? 
						AND
							((id_customer LIKE ? OR ? = '')
	 					OR
	 						(nama_customer LIKE ? OR ? = ''))`

	getCustomerByID  = "GetCustomerByID"
	qGetCustomerByID = `SELECT * FROM customer WHERE id_customer = ?`

	createCustomer  = "CreateCustomer"
	qCreateCustomer = `INSERT INTO sys.customer(id_customer, company_id, nama_customer, alamat, pic, penandatangan, jabatan, no_telp, updated_by)
						VALUES(?,?,?,?,?,?,?,?,?)`

	updateCustomer  = "UpdateCustomer"
	qUpdateCustomer = `UPDATE sys.customer
						SET 
							nama_customer = COALESCE(NULLIF(?,''), nama_customer), 
							alamat = COALESCE(NULLIF(?,''), alamat), 
							pic = COALESCE(NULLIF(?,''), pic), 
							penandatangan = COALESCE(NULLIF(?,''), penandatangan), 
							jabatan = COALESCE(NULLIF(?,''), jabatan), 
							no_telp = COALESCE(NULLIF(?,''), no_telp), 
							updated_by = COALESCE(NULLIF(?,''), updated_by)
						WHERE 
							id_customer = ?`

	// Company
	getAllCompanies  = "GetAllCompanies"
	qGetAllCompanies = `SELECT * FROM company`

	getCompanyByID  = "GetCompanyByID"
	qGetCompanyByID = `SELECT * FROM company WHERE company_id = ?`

	// Bank
	getAllBanks  = "GetAllBanks"
	qGetAllBanks = `SELECT * FROM bank`

	getBankByID  = "GetBankByID"
	qGetBankByID = `SELECT * FROM bank WHERE bank_id = ?`
)

// Add queries to key value order to be initialized as prepared statements
var (
	readStmt = []statement{
		//Contracts
		{getAllContractsHeader, qGetAllContractsHeader},
		{getAllContractsHeaderPage, qGetAllContractsHeaderPage},
		{getAllContractsHeaderCount, qGetAllContractsHeaderCount},
		{getContractHeaderByContractNumber, qGetContractHeaderByContractNumber},
		{getContractDetailsByContractNumber, qGetContractDetailsByContractNumber},

		//Customers
		{getAllCustomers, qGetAllCustomers},
		{getCustomerFiltered, qGetCustomerFiltered},
		{getCustomerFilteredCount, qGetCustomerFilteredCount},
		{getCustomerByID, qGetCustomerByID},

		//Company
		{getAllCompanies, qGetAllCompanies},
		{getCompanyByID, qGetCompanyByID},

		//Bank
		{getAllBanks, qGetAllBanks},
		{getBankByID, qGetBankByID},
	}
	insertStmt = []statement{
		{createContractHeader, qCreateContractHeader},
		{createContractDetail, qCreateContractDetail},

		{createCustomer, qCreateCustomer},
	}
	updateStmt = []statement{
		{updateCustomer, qUpdateCustomer},
	}
	deleteStmt = []statement{}
)

// New returns data instance
func New(db *sqlx.DB) Data {
	d := Data{
		db: db,
	}

	d.initStmt()
	return d
}

// Initialize queries as prepared statements
func (d *Data) initStmt() {
	var (
		err   error
		stmts = make(map[string]*sqlx.Stmt)
	)

	for _, v := range readStmt {
		stmts[v.key], err = d.db.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize statement key %v, err : %v", v.key, err)
		}
	}
	for _, v := range insertStmt {
		stmts[v.key], err = d.db.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize insert statement key %v, err : %v", v.key, err)
		}
	}

	for _, v := range updateStmt {
		stmts[v.key], err = d.db.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize update statement key %v, err : %v", v.key, err)
		}
	}

	for _, v := range deleteStmt {
		stmts[v.key], err = d.db.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize delete statement key %v, err : %v", v.key, err)
		}
	}

	d.stmt = stmts
}

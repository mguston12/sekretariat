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
							p.payment_id, p.paling_lambat, p.melunasi, p.tertunda,
							deposit, denda_satupersenyn,
							active_yn, kh.updated_by, kh.updated_at  FROM kontrak_header kh
							LEFT JOIN customer cu ON kh.id_customer = cu.id_customer
							LEFT JOIN bank ba ON kh.bank_id = ba.bank_id
							LEFT JOIN payment_method p ON p.payment_id = kh.payment_id
							WHERE kh.company_id = ? AND kh.active_yn = 'Y'
							AND
	 							((no_kontrak LIKE ? OR ? = '')
	 						OR
	 							(cu.nama_customer LIKE ? OR ? = ''))`

	getAllContractsHeaderPage  = "GetAllContractsHeaderPage"
	qGetAllContractsHeaderPage = `SELECT no_kontrak, tanggal_buat, kh.company_id, kh.id_customer,
								cu.nama_customer, cu.alamat, cu.pic, cu.penandatangan, cu.jabatan, cu.no_telp,
								ba.bank_id, ba.bank_name, ba.nomor_rekening, ba.atas_nama, 
								p.payment_id, p.paling_lambat, p.melunasi, p.tertunda,
								deposit, denda_satupersenyn,
								active_yn, kh.updated_by, kh.updated_at  FROM kontrak_header kh
								LEFT JOIN customer cu ON kh.id_customer = cu.id_customer
								LEFT JOIN bank ba ON kh.bank_id = ba.bank_id
								LEFT JOIN payment_method p ON p.payment_id = kh.payment_id
								WHERE kh.company_id = ? AND kh.active_yn = 'Y'
								AND
									((no_kontrak LIKE ? OR ? = '')
								OR
									(cu.nama_customer LIKE ? OR ? = '')) LIMIT ?,?`

	getAllContractsHeaderCount  = "GetAllContractsHeaderCount"
	qGetAllContractsHeaderCount = `SELECT count(*) FROM kontrak_header kh
								LEFT JOIN customer cu ON kh.id_customer = cu.id_customer
								LEFT JOIN bank ba ON kh.bank_id = ba.bank_id 
								LEFT JOIN payment_method p ON p.payment_id = kh.payment_id
								WHERE kh.company_id = ? AND kh.active_yn = 'Y'
								AND
									((no_kontrak LIKE ? OR ? = '')
								OR
									(cu.nama_customer LIKE ? OR ? = ''))`

	getContractHeaderByContractNumber  = "GetContractHeaderByContractNumber"
	qGetContractHeaderByContractNumber = `SELECT id_header, no_kontrak, tanggal_buat, kh.company_id, kh.id_customer,
											cu.nama_customer, cu.alamat, cu.pic, cu.penandatangan, cu.jabatan, cu.no_telp,
											ba.bank_id, ba.bank_name, ba.nomor_rekening, ba.atas_nama, 
											p.payment_id, p.paling_lambat, p.melunasi, p.tertunda,
											deposit, denda_satupersenyn,
											active_yn, kh.updated_by, kh.updated_at  FROM kontrak_header kh
											LEFT JOIN customer cu ON kh.id_customer = cu.id_customer
											LEFT JOIN bank ba ON kh.bank_id = ba.bank_id
											LEFT JOIN payment_method p ON p.payment_id = kh.payment_id
											WHERE kh.company_id = ? AND no_kontrak LIKE ? AND kh.active_yn = 'Y'`

	getContractDetailsByContractNumber  = `GetContractDetailsByContractNumber`
	qGetContractDetailsByContractNumber = `SELECT id_detail, no_kontrak, quantity, tipe_mesin, speed, harga_sewa, 
										free_copy, over_copy, free_copy_color, over_copy_color, periode_awal, periode_akhir, penempatan,
										active_yn, updated_by, updated_at FROM kontrak_detail WHERE no_kontrak LIKE ?`

	getContractByNumber  = `GetContractByNumber`
	qGetContractByNumber = `SELECT no_kontrak FROM kontrak_header WHERE no_kontrak = ?`

	getContractExp30Days  = "GetContractExp30Days"
	qGetContractExp30Days = `SELECT d.*
							FROM kontrak_detail d
							JOIN kontrak_header h ON h.no_kontrak = d.no_kontrak
							WHERE h.company_id = ?
							AND d.periode_akhir >= DATE_FORMAT(CURDATE() + INTERVAL 1 MONTH, '%Y-%m-01')
							AND d.periode_akhir < DATE_FORMAT(CURDATE() + INTERVAL 2 MONTH, '%Y-%m-01');`

	createContractHeader  = "CreateContractHeader"
	qCreateContractHeader = `INSERT INTO kontrak_header(no_kontrak, tanggal_buat, company_id, id_customer, bank_id, payment_id, deposit, denda_satupersenyn, active_yn, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	createContractDetail  = "CreateContractDetail"
	qCreateContractDetail = `INSERT INTO kontrak_detail(no_kontrak, quantity, tipe_mesin, speed,
							harga_sewa, free_copy, over_copy, free_copy_color, over_copy_color, periode_awal, periode_akhir, penempatan, active_yn, updated_by)
							VALUES (?, ?, ?, ?,
							?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	updateContractHeader  = "UpdateContractHeader"
	qUpdateContractHeader = `UPDATE kontrak_header
					SET 
						no_kontrak = COALESCE(NULLIF(?,''), no_kontrak),
						tanggal_buat = COALESCE(NULLIF(?,''), tanggal_buat), 
						id_customer = COALESCE(NULLIF(?,''), id_customer), 
						bank_id = COALESCE(NULLIF(?,''), bank_id), 
						payment_id = COALESCE(NULLIF(?,''), payment_id), 
						deposit = COALESCE(NULLIF(?,''), 0),
						denda_satupersenyn = COALESCE(NULLIF(?,''), denda_satupersenyn),
						active_yn = COALESCE(NULLIF(?,''), active_yn),
						updated_by = COALESCE(NULLIF(?,''), updated_by)
					WHERE 
						id_header = ?`

	updateContractDetail  = "UpdateContractDetail"
	qUpdateContractDetail = `UPDATE kontrak_detail
						SET 
							active_yn = COALESCE(NULLIF(?,''), active_yn),
							updated_by = COALESCE(NULLIF(?,''), updated_by)
						WHERE 
							no_kontrak = ?`

	deleteContractDetail  = "DeleteContractDetail"
	qDeleteContractDetail = `DELETE FROM kontrak_detail WHERE no_kontrak = ?`

	// Customers
	getCustomer  = "GetCustomer"
	qGetCustomer = `SELECT * FROM customer WHERE (nama_customer LIKE ? OR ? = '')`

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

	getCustomerByNameAndAddress    = "GetCustomerIDByNameAndAddress"
	qGetCustomerIDByNameAndAddress = `SELECT id_customer FROM customer WHERE nama_customer = ? AND alamat = ?`

	createCustomer  = "CreateCustomer"
	qCreateCustomer = `INSERT INTO customer(id_customer, company_id, nama_customer, alamat, pic, penandatangan, jabatan, no_telp, updated_by)
						VALUES(?,?,?,?,?,?,?,?,?)`

	updateCustomer  = "UpdateCustomer"
	qUpdateCustomer = `UPDATE customer
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
	qGetBankByID = `SELECT * FROM bank WHERE bank_id = ? ORDER BY company_id ASC`

	getBankByCompanyID  = "GetBankByCompanyID"
	qGetBankByCompanyID = `SELECT * FROM bank WHERE company_id = ?`

	createBank  = "CreateBank"
	qCreateBank = `INSERT INTO bank(bank_name, nomor_rekening, atas_nama, company_id)
						VALUES(?,?,?,?)`

	updateBank  = "UpdateBank"
	qUpdateBank = `UPDATE bank
						SET
							bank_name = COALESCE(NULLIF(?,''), bank_name),
							nomor_rekening = COALESCE(NULLIF(?,''), nomor_rekening),
							atas_nama = COALESCE(NULLIF(?,''), atas_nama)
						WHERE
							bank_id = ?`

	deleteBankByID  = "DeleteBankByID"
	qDeleteBankByID = `DELETE FROM bank WHERE bank_id = ?`

	// PaymentMethod
	getPaymentMethod  = "GetPaymentMethod"
	qGetPaymentMethod = `SELECT * FROM payment_method`
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
		{getContractExp30Days, qGetContractExp30Days},
		{getContractByNumber, qGetContractByNumber},

		//Customers
		{getAllCustomers, qGetAllCustomers},
		{getCustomerFiltered, qGetCustomerFiltered},
		{getCustomerFilteredCount, qGetCustomerFilteredCount},
		{getCustomerByID, qGetCustomerByID},
		{getCustomerByNameAndAddress, qGetCustomerIDByNameAndAddress},
		{getCustomer, qGetCustomer},
		//Company
		{getAllCompanies, qGetAllCompanies},
		{getCompanyByID, qGetCompanyByID},

		//Bank
		{getAllBanks, qGetAllBanks},
		{getBankByID, qGetBankByID},
		{getBankByCompanyID, qGetBankByCompanyID},

		//Payment Method
		{getPaymentMethod, qGetPaymentMethod},
	}
	insertStmt = []statement{
		{createContractHeader, qCreateContractHeader},
		{createContractDetail, qCreateContractDetail},

		{createCustomer, qCreateCustomer},
		{createBank, qCreateBank},
	}
	updateStmt = []statement{
		{updateContractHeader, qUpdateContractHeader},
		{updateContractDetail, qUpdateContractDetail},

		{updateCustomer, qUpdateCustomer},
		{updateBank, qUpdateBank},
	}
	deleteStmt = []statement{
		{deleteContractDetail, qDeleteContractDetail},

		{deleteBankByID, qDeleteBankByID},
	}
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

package sekretariat

import "time"

type KontrakHeader struct {
	NoKontrak   string `db:"no_kontrak" json:"no_kontrak"`
	TanggalBuat string `db:"tanggal_buat" json:"tanggal_buat"`
	CompanyID   int    `db:"company_id" json:"company_id"`
	Customer
	Bank
	Pembayaran
	Deposit           float64         `db:"deposit" json:"deposit"`
	DendaSatuPersenYN string          `db:"denda_satupersenyn" json:"denda_satupersenyn"`
	ActiveYN          string          `db:"active_yn" json:"active_yn"`
	UpdatedBy         string          `db:"updated_by" json:"updated_by"`
	UpdatedAt         time.Time       `db:"updated_at" json:"updated_at"`
	Details           []KontrakDetail `json:"details"`
}

type KontrakDetail struct {
	ID                 int       `db:"id_detail" json:"id_detail"`
	NoKontrak          string    `db:"no_kontrak" json:"no_kontrak"`
	Quantity           int       `db:"quantity" json:"quantity"`
	Tipe               string    `db:"tipe_mesin" json:"tipe_mesin"`
	Speed              string    `db:"speed" json:"speed"`
	Harga              float64   `db:"harga_sewa" json:"harga_sewa"`
	FreeCopy           string    `db:"free_copy" json:"free_copy"`
	OverCopy           string    `db:"over_copy" json:"over_copy"`
	FreeCopyColor      string    `db:"free_copy_color" json:"free_copy_color"`
	OverCopyColor      string    `db:"over_copy_color" json:"over_copy_color"`
	PeriodeAwal        time.Time `db:"periode_awal" json:"periode_awal"`
	PeriodeAwalString  string    `json:"periode_awal_string"`
	PeriodeAkhir       time.Time `db:"periode_akhir" json:"periode_akhir"`
	PeriodeAkhirString string    `json:"periode_akhir_string"`
	Penempatan         string    `db:"penempatan" json:"penempatan"`
	ActiveYN           string    `db:"active_yn" json:"active_yn"`
	UpdatedBy          string    `db:"updated_by" json:"updated_by"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

type Customer struct {
	CustomerID    string    `db:"id_customer" json:"id_customer"`
	CompanyID     int       `db:"company_id" json:"company_id"`
	CustomerName  string    `db:"nama_customer" json:"nama_customer"`
	Address       string    `db:"alamat" json:"alamat_customer"`
	PIC           string    `db:"pic" json:"pic"`
	PenandaTangan string    `db:"penandatangan" json:"penandatangan"`
	Jabatan       string    `db:"jabatan" json:"jabatan"`
	NoTelp        string    `db:"no_telp" json:"no_telp"`
	UpdatedBy     string    `db:"updated_by" json:"updated_by"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type Company struct {
	ID      int    `db:"company_id" json:"company_id"`
	Name    string `db:"company_name" json:"company_name"`
	Address string `db:"company_address" json:"company_address"`
	Phone   string `db:"company_phonenumber" json:"company_phonenumber"`
	PIC     string `db:"company_pic" json:"company_pic"`
}

type Bank struct {
	ID        int    `db:"bank_id" json:"bank_id"`
	Name      string `db:"bank_name" json:"bank_name"`
	Norek     string `db:"nomor_rekening" json:"nomor_rekening"`
	AtasNama  string `db:"atas_nama" json:"atas_nama"`
	CompanyID int    `db:"company_id" json:"company_id"`
}

type Pembayaran struct {
	ID           int `db:"payment_id" json:"payment_id"`
	PalingLambat int `db:"paling_lambat" json:"paling_lambat"`
	Melunasi     int `db:"melunasi" json:"melunasi"`
	Tertunda     int `db:"tertunda" json:"tertunda"`
}

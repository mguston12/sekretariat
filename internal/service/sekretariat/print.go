package sekretariat

import (
	"bytes"
	"context"
	"math"
	"strings"

	"sekretariat/pkg/errors"
	"strconv"

	"github.com/jung-kurt/gofpdf"
	"github.com/yudapc/go-rupiah"
)

func (s Service) PrintKontrak(ctx context.Context, company int, no_kontrak string) (bytes.Buffer, error) {
	b := bytes.Buffer{}

	kontrak, err := s.GetDataContractByContractNumber(ctx, company, no_kontrak)
	if err != nil {
		return b, errors.Wrap(err, "[SERVICE][PrintKontrak]")
	}

	customer, err := s.data.GetCustomerByID(ctx, kontrak.CustomerID)
	if err != nil {
		return b, errors.Wrap(err, "[SERVICE][PrintKontrak]")
	}

	dataCompany, err := s.data.GetCompanyByID(ctx, kontrak.CompanyID)
	if err != nil {
		return b, errors.Wrap(err, "[SERVICE][PrintKontrak]")
	}

	bank, err := s.data.GetBankByID(ctx, kontrak.Bank.ID)
	if err != nil {
		return b, errors.Wrap(err, "[SERVICE][PrintKontrak]")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	initialY := 0.0

	if kontrak.CompanyID == 3 || kontrak.CompanyID == 2 {
		pdf.SetXY(25, 33)
	} else {
		pdf.SetXY(25, 25)
	}

	// KONTRAK BARU (PER 9 AGUSTUS 2024)
	pdf.SetFont("times", "B", 10.5)
	initialY = pdf.GetY()

	pdf.Text(74, initialY+5, "PERJANJIAN SEWA MENYEWA")
	pdf.Text(70, initialY+10, "MESIN PRINTER MULTIFUNCTION")

	pdf.SetXY(53, initialY+11)
	paragraphWidth := 100.0
	text := "Nomor: SPK. " + kontrak.NoKontrak
	pdf.MultiCell(paragraphWidth, 5, text, "", "C", false)

	height := initialY + 10

	// Add title

	paragraphWidth = 167.0

	pdf.SetXY(20, height+10)
	pdf.SetFont("times", "", 10.5)

	text = "Pada hari ini " + kontrak.TanggalBuat + ", dibuat dan ditandatangani Perjanjian Sewa Menyewa Mesin Printer Multifunction (untuk selanjutnya disebut sebagai \"Perjanjian\") oleh dan antara :"

	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)

	// PIHAK SATU 32.8
	pdf.Text(25, height+28.8, "1.")
	pdf.SetXY(30, height+25)
	pdf.SetFont("times", "", 10.5)
	paragraphWidth = 157.0
	textPihakPertama := dataCompany.PIC + ", Direktur, sah mewakili " + strings.ToUpper(dataCompany.Name) + " suatu perseroan terbatas yang berkedudukan di " + dataCompany.Address + " (\"Pihak Pertama\"); dan"
	pdf.MultiCell(paragraphWidth, 4.5, textPihakPertama, "", "J", false)

	// PIHAK DUA
	pdf.Text(25, height+40.5, "2.")
	pdf.SetXY(30, height+37)
	pdf.SetFont("times", "", 10.5)
	paragraphWidth = 157.0

	initialY = pdf.GetY()

	textPihakKedua := customer.PenandaTangan + ", " + customer.Jabatan + ", mewakili " + strings.ToUpper(customer.CustomerName) + ", beralamatkan di " + customer.Address + " (\"Pihak Kedua\")."
	pdf.MultiCell(paragraphWidth, 4.5, textPihakKedua, "", "J", false)

	finalY := pdf.GetY()
	lines := (finalY - initialY) / 5

	height = height + 35 + (lines+1)*5

	paragraphWidth = 168.0

	pdf.SetXY(20, height)
	pdf.SetFont("times", "", 10.5)
	text = "Pihak Pertama dan Pihak Kedua secara bersama-sama disebut sebagai \"Para Pihak\" dan masing-masing sebagai \"Pihak\"."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)

	pdf.SetXY(20, height+12)
	pdf.SetFont("times", "", 10.5)
	text = "Kontrak Perjanjian ini menjadi satu kesatuan dari form pemesanan (Order Confirmation) dan tidak terpisahkan dari konfirmasi dan permohonan sewa."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)

	pdf.SetXY(20, height+24)
	pdf.SetFont("times", "", 10.5)
	text = "Menyatakan bahwa Para Pihak setuju dan sepakat untuk menandatangani perjanjian  untuk  menyewa Mesin milik Pihak Pertama tersebut untuk dipergunakan sebagai  sarana  operasional Pihak Kedua dengan ketentuan dan syarat-syarat sebagai berikut:"
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)

	height = height + 43

	// PASAL 1
	pdf.SetFont("times", "B", 10.5)
	pdf.Text(95, height, "PASAL 1")
	pdf.Text(83, height+4, "HAK DAN KEWAJIBAN")

	pdf.Text(21, height+9, "I.")
	pdf.Text(26, height+9, "Pihak Pertama")

	height = height + 11.5

	if len(kontrak.Details) == 1 {
		pdf.SetFont("times", "", 10.5)
		pdf.Text(26, height+2.5, "1.")

		initialY = pdf.GetY()

		pdf.SetXY(32, height-1)
		paragraphWidth = 159.0

		terbilangQuantity := numberToWords(kontrak.Details[0].Quantity)
		text = "Menyediakan " + strconv.Itoa(kontrak.Details[0].Quantity) + " (" + terbilangQuantity + ")" + " unit Mesin Printer Multifunction \"" + kontrak.Details[0].Tipe + "\", dalam kondisi baik untuk disewakan kepada Pihak Kedua di " + customer.CustomerName + ", penempatan Mesin di " + kontrak.Details[0].Penempatan + "."
		pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
		finalY = pdf.GetY()

		lines = (finalY - initialY) / 9.5
	}

	if lines < 3.5 {
		height = height + (math.Floor(lines))*4.5 + 1.5
	} else {
		height = height + (math.Ceil(lines))*4.5 + 1.5

	}

	pdf.Text(26, height+1, "2.")
	pdf.SetXY(32, height-2.5)
	paragraphWidth = 159.0
	text = "Jika Mesin rusak dan tidak bisa diperbaiki, Pihak Pertama akan menggantinya dengan Mesin serupa dalam jangka waktu 2 (dua) hari kerja."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)

	pdf.Text(26, height+10.2, "3.")
	pdf.SetXY(32, height+6.5)
	paragraphWidth = 159.0
	text = "Sisa bahan pakai dan suku cadang bekas tetap milik Pihak Pertama."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)

	height = height + 16

	pdf.SetFont("times", "B", 10.5)
	pdf.Text(20, height, "II.")
	pdf.Text(26, height, "Pihak Kedua")

	height = height + 5

	count := 0.0
	height_number := height
	height_text := height - 3.3

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number, "1.")
	pdf.SetXY(32, height_text)
	paragraphWidth = 159.0
	text = "Menyediakan kertas dan memastikan aliran listrik sesuai spesifikasi Mesin."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count++

	pdf.Text(26, height_number+count*4.5, "2.")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	text = "Bertanggung jawab atas fasilitas listrik dan penggunaan Mesin, menjaga Mesin dari kerusakan, termasuk akibat kebakaran."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 2

	pdf.Text(26, height_number+count*4.5, "3.")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	text = "Tidak boleh menyalahgunakan Mesin, merusak perlengkapan, memindahkan Mesin tanpa izin, menyewakan Mesin kepada pihak ketiga, atau menjadikannya jaminan."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 2

	pdf.Text(26, height_number+count*4.5, "4.")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	text = "Bertanggung jawab atas kerugian jika melanggar ketentuan."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 2

	height = height_number - 2 + count*4.5

	pdf.SetFont("times", "B", 10.5)
	pdf.Text(95, height, "PASAL 2")
	pdf.Text(90, height+4, "HARGA SEWA")

	height = height + 9

	count = 0.0
	number := 1
	height_number = height
	height_text = height - 3.3

	formatRupiah := rupiah.FormatRupiah(kontrak.Details[0].Harga)

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text)
	paragraphWidth = 159.0
	freeCopyBW, overCopyBW, freeCopyColor, overCopyColor := "", "", "", ""

	if kontrak.Details[0].FreeCopy != "" && kontrak.Details[0].FreeCopy != "0" {
		freeCopyBW = ", Free Copy BW " + kontrak.Details[0].FreeCopy
	}
	if kontrak.Details[0].OverCopy != "" && kontrak.Details[0].OverCopy != "0" {
		overCopyBW = ", Over Copy BW " + kontrak.Details[0].OverCopy
	}
	if kontrak.Details[0].FreeCopyColor != "" && kontrak.Details[0].FreeCopyColor != "0" {
		freeCopyColor = ", Free Copy Color " + kontrak.Details[0].FreeCopyColor
	}
	if kontrak.Details[0].OverCopyColor != "" && kontrak.Details[0].OverCopyColor != "0" {
		overCopyColor = ", Over Copy Color " + kontrak.Details[0].OverCopyColor
	}
	initialY = pdf.GetY()

	text = "Biaya sewa: " + formatRupiah + "/bulan" + freeCopyBW + overCopyBW + freeCopyColor + overCopyColor + "."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)

	finalY = pdf.GetY()

	lines = (finalY - initialY) / 4.5
	count = count + lines
	number++

	if dataCompany.ID == 3 {
		pdf.SetFont("times", "", 10.5)
		pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
		pdf.SetXY(32, height_text+count*4.5)
		paragraphWidth = 159.0
		text = "Harga belum termasuk PPN."
		pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
		count++
		number++
	}
	if kontrak.Deposit != 0.00 {
		formatDeposit := rupiah.FormatRupiah(kontrak.Deposit)
		pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
		pdf.SetXY(32, height_text+count*4.5)
		paragraphWidth = 159.0
		text = "Pihak Kedua membayar deposit sebesar " + formatDeposit + " yang dikembalikan oleh Pihak Pertama kepada Pihak Kedua saat perjanjian kontrak sewa Mesin berakhir dan Pihak Kedua sudah menyelesaikan semua kewajiban."
		pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
		count = count + 3
		number++
	}

	pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	text = "Pembayaran tagihan dapat dilakukan secara Cash atau ditransfer ke Rekening dengan data sebagai berikut :"
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 3

	pdf.Text(33, height_text+count*4.5, "-    Nama Rekening")
	pdf.Text(80, height_text+count*4.5, ":")
	pdf.Text(83, height_text+count*4.5, bank.AtasNama)
	count++

	pdf.Text(33, height_text+count*4.5, "-    Bank")
	pdf.Text(80, height_text+count*4.5, ":")
	pdf.Text(83, height_text+count*4.5, bank.Name)
	count++

	pdf.Text(33, height_text+count*4.5, "-    Nomor Rekening")
	pdf.Text(80, height_text+count*4.5, ":")
	pdf.Text(83, height_text+count*4.5, bank.Norek)

	pdf.AddPage()
	// Page 2

	pdf.SetXY(25, 25)
	height = 25

	pdf.SetFont("times", "B", 10.5)
	pdf.Text(95, height, "PASAL 3")
	pdf.Text(82, height+4, "CARA PEMBAYARAN")

	height = height + 9

	count = 0.0
	number = 1
	height_number = height
	height_text = height - 3.3

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text)
	paragraphWidth = 159.0
	text = "Pembayaran dilakukan paling lambat " + strconv.Itoa(kontrak.PalingLambat) + " hari kerja setelah menerima kwitansi."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count++
	number++

	if kontrak.DendaSatuPersenYN == "Y" {
		pdf.SetFont("times", "", 10.5)
		pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
		pdf.SetXY(32, height_text+count*4.5)
		paragraphWidth = 159.0
		text = "Apabila Pihak Kedua tidak melunasi dalam waktu " + strconv.Itoa(kontrak.Melunasi) + " hari, Pihak Kedua akan dikenakan denda keterlambatan pembayaran sebesar 1%."
		pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
		count = count + 2
		number++
	}

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	text = "Biaya copy dihitung berdasarkan counter Mesin."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count++
	number++

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	text = "Jika pembayaran tertunda lebih dari " + strconv.Itoa(kontrak.Tertunda) + " hari, Mesin bisa dikunci oleh Pihak Pertama. Setelah 60 hari, Mesin dapat ditarik oleh Pihak Pertama dan Pihak Kedua wajib membayar semua tagihan dan denda."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 2
	number++

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	text = "Pihak Kedua harus membantu dalam penarikan Mesin jika diperlukan."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 2

	height = height_number + count*4.5

	pdf.SetFont("times", "B", 10.5)
	pdf.Text(95, height, "PASAL 4")
	pdf.Text(85, height+4, "PERUBAHAN HARGA")

	height = height + 9

	count = 0.0
	height_number = height
	height_text = height - 3.3

	pdf.SetFont("times", "", 10.5)
	pdf.SetXY(32, height_text)
	paragraphWidth = 159.0
	text = "Jika terjadi perubahan harga, Pihak Pertama akan memberi tahu 15 hari sebelumnya. Pihak Kedua dapat menolak dalam 7 hari setelah pemberitahuan jika tidak setuju."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 3

	height = height_number + count*4.5

	pdf.SetFont("times", "B", 10.5)
	pdf.Text(95, height, "PASAL 5")
	pdf.Text(83, height+4, "JANGKA WAKTU SEWA")

	height = height + 9

	count = 0.0
	number = 1
	height_number = height
	height_text = height - 3.3

	indonesianDateFormat := "02 January 2006"
	periodeAwalString := kontrak.Details[0].PeriodeAwal.Format(indonesianDateFormat)
	periodeAkhirString := kontrak.Details[0].PeriodeAkhir.Format(indonesianDateFormat)

	monthReplacements := map[string]string{
		"January":   "Januari",
		"February":  "Februari",
		"March":     "Maret",
		"April":     "April",
		"May":       "Mei",
		"June":      "Juni",
		"July":      "Juli",
		"August":    "Agustus",
		"September": "September",
		"October":   "Oktober",
		"November":  "November",
		"December":  "Desember",
	}

	// Replace month names in the string
	for enMonth, idMonth := range monthReplacements {
		periodeAwalString = strings.Replace(periodeAwalString, enMonth, idMonth, -1)
		periodeAkhirString = strings.Replace(periodeAkhirString, enMonth, idMonth, -1)
	}

	startYear, startMonth, startDay := kontrak.Details[0].PeriodeAwal.Date()
	endYear, endMonth, endDay := kontrak.Details[0].PeriodeAkhir.Date()

	// Hitung selisih tahun dan bulan
	yearDiff := endYear - startYear
	monthDiff := int(endMonth) - int(startMonth)

	// Total selisih bulan
	totalMonthDiff := yearDiff*12 + monthDiff

	// Jika Anda ingin memastikan selisih ini tetap 12
	if startDay < endDay {
		totalMonthDiff += 1
	}

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text)
	paragraphWidth = 159.0
	terbilangMonthDiff := numberToWords(totalMonthDiff)
	text = "Pihak Kedua menyewa Mesin Printer Multifunction kepada Pihak Pertama dengan jangka waktu sewa adalah " + strconv.Itoa(totalMonthDiff) + " (" + terbilangMonthDiff + ")" + " bulan, terhitung sejak " + periodeAwalString + " hingga " + periodeAkhirString + "."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 2
	number++

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	text = "Setelah jangka waktu sewa berakhir, Para Pihak setuju untuk memperpanjang secara otomatis(terus menerus) selama satu tahun berikutnya, kecuali salah satu pihak tidak berniat memperpanjang maka dapat dituangkan dalam surat konfirmasi penghentian kontrak."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 4

	height = height_number + count*4.5

	pdf.SetFont("times", "B", 10.5)
	pdf.Text(95, height, "PASAL 6")
	pdf.Text(85, height+4, "PEMBATALAN SEWA")

	height = height + 9

	count = 0.0
	number = 1
	height_number = height
	height_text = height - 3.3

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text)
	paragraphWidth = 159.0
	text = "Pihak Pertama dapat menghentikan sewa jika Pihak Kedua tidak membayar 2 bulan berturut-turut."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count++
	number++

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	text = "Pihak Kedua dapat menghentikan sewa jika Pihak Pertama tidak memenuhi kewajiban."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count++
	number++

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	brandNewYN := 0

	if len(kontrak.Details) > 0 && strings.Contains(kontrak.Details[0].Tipe, "Brand New") {
		brandNewYN = 2
	} else {
		brandNewYN = 1
	}

	if kontrak.Deposit != 0.00 {
		text = "Jika Pihak Kedua menghentikan sewa, harus memberi pemberitahuan 1 bulan sebelumnya dan membayar denda atau penalty sebesar " + strconv.Itoa(brandNewYN) + " bulan sewa. Deposit yang sudah dibayarkan hangus tidak dapat dikembalikan."
		count += 1
	} else {
		text = "Jika Pihak Kedua menghentikan sewa, harus memberi pemberitahuan 1 bulan sebelumnya dan membayar denda atau penalty sebesar " + strconv.Itoa(brandNewYN) + " bulan sewa."
	}

	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 2
	number++

	pdf.SetFont("times", "", 10.5)
	pdf.Text(26, height_number+count*4.5, strconv.Itoa(number)+".")
	pdf.SetXY(32, height_text+count*4.5)
	paragraphWidth = 159.0
	text = "Jika perjanjian dibatalkan, Mesin harus ditarik dalam 5 hari kerja, dan Pihak Kedua tetap wajib membayar semua tagihan hingga bulan berjalan."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 3

	height = height_number + count*4.5

	pdf.SetFont("times", "B", 10.5)
	pdf.Text(95, height, "PASAL 7")
	pdf.Text(91.5, height+4, "ADDENDUM")

	height = height + 9

	count = 0.0
	height_number = height
	height_text = height - 3.3

	pdf.SetFont("times", "", 10.5)
	pdf.SetXY(32, height_text)
	paragraphWidth = 159.0
	text = "Perubahan atau tambahan perjanjian akan diatur dalam addendum yang merupakan kesatuan dan bagian yang tidak dapat dipisahkan dari Perjanjian ini."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 3

	height = height_number + count*4.5

	pdf.SetFont("times", "B", 10.5)
	pdf.Text(95, height, "PASAL 8")
	pdf.Text(94, height+4, "PENUTUP")

	height = height + 9

	count = 0.0
	height_number = height
	height_text = height - 3.3

	pdf.SetFont("times", "", 10.5)
	pdf.SetXY(32, height_text)
	paragraphWidth = 159.0
	text = "Perjanjian ini ditandatangani di Jakarta dalam rangkap 2 (dua), masing-masing bermaterai dan memiliki kekuatan hukum yang sama."
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "J", false)
	count = count + 3

	height = height + count*4.5

	pdf.SetFont("times", "B", 10)
	pdf.Text(45, height, "Pihak Pertama")
	pdf.Text(140, height, "Pihak Kedua")

	pdf.Text(38, height+5, strings.ToUpper(dataCompany.Name))
	pdf.SetXY(100, height+1.2)
	paragraphWidth = 100.0
	text = strings.ToUpper(customer.CustomerName)
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "C", false)

	pdf.SetFont("times", "BU", 10)
	pdf.Text(45, height+40, dataCompany.PIC)
	pdf.SetXY(115, height+36.2)
	paragraphWidth = 70.0
	text = customer.PenandaTangan
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "C", false)

	pdf.SetFont("times", "B", 10)
	pdf.Text(50, height+44, "Direktur")

	pdf.SetXY(125, height+40)
	paragraphWidth = 50.0
	text = customer.Jabatan
	pdf.MultiCell(paragraphWidth, 4.5, text, "", "C", false)

	// KONTRAK LAMA
	// pdf.SetFont("times", "B", 11)
	// initialY := pdf.GetY()
	// pdf.Text(74, initialY+5, "PERJANJIAN SEWA MENYEWA")
	// pdf.Text(70, initialY+9, "MESIN PRINTER MULTIFUNCTION")

	// pdf.SetXY(53, initialY+9)
	// paragraphWidth := 100.0
	// text := "Nomor: SPK. " + kontrak.NoKontrak
	// pdf.MultiCell(paragraphWidth, 5, text, "", "C", false)

	// height := initialY + 9

	// // Add title

	// paragraphWidth = 167.0

	// pdf.SetXY(20, 45)
	// pdf.SetFont("times", "", 11)

	// text = "Pada hari ini " + kontrak.TanggalBuat + ", dibuat dan ditandatangani Perjanjian Sewa Menyewa Mesin Printer Multifunction (untuk selanjutnya disebut sebagai \"Perjanjian\") oleh dan antara :"

	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// // PIHAK SATU 32.8
	// pdf.Text(25, height+32.8, "1.")
	// pdf.SetXY(30, height+29)
	// pdf.SetFont("times", "", 11)
	// paragraphWidth = 157.0
	// textPihakPertama := dataCompany.PIC + " selaku Direktur, bertindak untuk dan atas nama serta sah mewakili " + strings.ToUpper(dataCompany.Name) + " suatu perseroan terbatas yang berkedudukan di " + dataCompany.Address + " (selanjutnya disebut sebagai \"Pihak Pertama\"); dan"
	// pdf.MultiCell(paragraphWidth, 5, textPihakPertama, "", "J", false)

	// // PIHAK DUA
	// pdf.Text(25, height+51.5, "2.")
	// pdf.SetXY(30, height+48)
	// pdf.SetFont("times", "", 11)
	// paragraphWidth = 157.0

	// initialY = pdf.GetY()

	// textPihakKedua := customer.PenandaTangan + " selaku " + customer.Jabatan + ", bertindak untuk dan atas nama serta sah mewakili " + strings.ToUpper(customer.CustomerName) + ", beralamatkan di " + customer.Address + " (selanjutnya disebut sebagai \"Pihak Kedua\")."
	// pdf.MultiCell(paragraphWidth, 5, textPihakKedua, "", "J", false)

	// finalY := pdf.GetY()
	// lines := (finalY - initialY) / 5

	// height = height + 48 + (lines+1)*5

	// paragraphWidth = 168.0

	// pdf.SetXY(20, height)
	// pdf.SetFont("times", "", 11)
	// text = "Pihak Pertama dan Pihak Kedua secara bersama-sama disebut sebagai \"Para Pihak\" dan masing-masing sebagai \"Pihak\"."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.SetXY(20, height+15)
	// pdf.SetFont("times", "", 11)
	// text = "Objek dari Perjanjian ini adalah mesin Printer Multifunction, milik Pihak Pertama untuk  selanjutnya disebut \"Mesin\"."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.SetXY(20, height+30)
	// pdf.SetFont("times", "", 11)
	// text = "Menyatakan bahwa Para Pihak setuju dan sepakat untuk menandatangani perjanjian  untuk  menyewa Mesin milik Pihak Pertama tersebut untuk dipergunakan sebagai  sarana  operasional Pihak Kedua dengan ketentuan dan syarat-syarat sebagai berikut:"
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// height = height + 30

	// // PASAL 1
	// pdf.SetFont("times", "B", 11)
	// pdf.Text(95, height+25, "PASAL 1")
	// pdf.Text(83, height+29, "HAK DAN KEWAJIBAN")

	// pdf.Text(21, height+38, "I.")
	// pdf.Text(26, height+38, "Hak dan kewajiban Pihak Pertama")

	// height = height + 40.5

	// pdf.SetFont("times", "", 11)
	// pdf.Text(26, height+3.5, "1.")

	// initialY = pdf.GetY()
	// pdf.SetXY(32, height)
	// paragraphWidth = 159.0
	// text = "Pihak Pertama menyediakan 1 (Satu) unit mesin Printer Multifunction merk \"" + kontrak.Details[0].Tipe + "\", dengan keadaan dapat dipergunakan dengan baik untuk disewakan kepada Pihak Kedua yang ditempatkan di " + customer.CustomerName + ", beralamatkan di " + kontrak.Details[0].Penempatan + ". Penyediaan yang di maksud ayat I.1 termasuk pelaksanaan pemasanganya, bahan pakai (kecuali kertas)."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)
	// finalY = pdf.GetY()

	// lines = (finalY - initialY) / 10

	// height = height + (lines+1)*5
	// pdf.Text(25.5, height, "2.")
	// pdf.SetXY(32, height-4)
	// paragraphWidth = 159.0
	// text = "Apabila Mesin tersebut tidak dapat berfungsi dengan baik setelah diadakan perbaikan dan evaluasi secara teknis maka Pihak Pertama akan  menyediakan  dan  menukar  Mesin dengan  mesin lain yang sejenis yang kondisinya lebih baik paling lambat dalam 2 (dua) hari kerja."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(25.5, height+16.2, "3.")
	// pdf.SetXY(32, height+12.5)
	// paragraphWidth = 159.0
	// text = "Dalam hubungan dengan ayat I butir 2 Pasal ini, sisa  bahan  pakai  dan suku cadang bekas dikembalikan dan menjadi milik Pihak Pertama."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(25.5, height+27.7, "4.")
	// pdf.SetXY(32, height+24)
	// paragraphWidth = 159.0
	// text = "Pihak Pertama wajib memberikan Spesifikasi atas Mesin kepada Pihak Kedua."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(190, 274, "1")

	// pdf.SetFont("times", "I", 9)
	// pdf.Text(20, 274, "Perjanjian Sewa Menyewa Mesin Printer Multifunction")
	// pdf.Rect(142, 277, 50, 12, "D")
	// pdf.Line(167, 277, 167, 289)

	// pdf.SetFont("times", "UB", 7)
	// pdf.Text(143, 279.5, "Pihak Pertama")
	// pdf.Text(168, 279.5, "Pihak Kedua")
	// pdf.AddPage()

	// // Page 2
	// pdf.SetFont("times", "B", 11)
	// pdf.Text(18, 25, "II.")
	// pdf.Text(26, 25, "Hak dan kewajiban Pihak Kedua")

	// height = 30
	// pdf.SetFont("times", "", 11)
	// pdf.Text(26, height, "1.")
	// pdf.SetXY(32, height-3.8)
	// paragraphWidth = 159.0
	// text = "Pihak Kedua, untuk terlaksananya perjanjian   sewa  menyewa  dengan   baik,   menyediakan   kertas sesuai dengan spesifikasi Mesin, termasuk aliran listrik yang cukup, baik dalam jumlah maupun kadar tegangan sesuai dengan kebutuhan spesifikasi Mesin."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+16, "2.")
	// pdf.SetXY(32, height+12.2)
	// paragraphWidth = 159.0
	// text = "Setiap pengadaan fasilitas yang diperlukan  untuk  aliran  listrik, misal : Stop kontak; Transformator; Stabilizer; Switch hub; Kabel; atau alat konektor ke internet dan peralatan lainnya menjadi tanggung jawab  Pihak Kedua. Pihak Kedua bertanggung jawab dalam penggunaan Mesin dengan baik dan menjaga segala kemungkinan yang dapat mengakibatkan kerusakan Mesin."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+37, "3.")
	// pdf.SetXY(32, height+33.2)
	// paragraphWidth = 159.0
	// text = "Pihak Kedua tidak diperkenankan untuk:"
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(33, height+42, "-")
	// pdf.SetXY(38.5, height+38.2)
	// paragraphWidth = 152.0
	// text = "Menyalahgunakan pemakaian dan atau fungsi  mesin yang bertentangan dengan UU dan atau peraturan yang berlaku di Negara Republik Indonesia"
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(33, height+52, "-")
	// pdf.SetXY(38, height+48.2)
	// paragraphWidth = 152.0
	// text = "Mengubah; merusak; menghilangkan  perlengkapan  mesin, tanda pengenalan atau meteran Mesin (counter)"
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(33, height+62, "-")
	// pdf.SetXY(38, height+58.2)
	// paragraphWidth = 152.0
	// text = "Membongkar; memasang dan memindahkan mesin tanpa persetujuan Pihak Pertama. "
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(33, height+67, "-")
	// pdf.SetXY(38, height+63.2)
	// paragraphWidth = 152.0
	// text = "Memindah sewakan Mesin ke Pihak Ketiga. "
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(33, height+72, "-")
	// pdf.SetXY(38, height+68.2)
	// paragraphWidth = 152.0
	// text = "Menjadikan mesin sebagai jaminan apabila ada sengketa antara Pihak  Kedua dengan pihak lainnya."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+83, "4.")
	// pdf.SetXY(32, height+79.2)
	// paragraphWidth = 160.0
	// text = "Pihak  Kedua  bertanggung  jawab  serta  mengganti kerugian  kepada  Pihak  Pertama apabila Pihak Kedua melanggar ketentutan-ketentuan yang tercantum dalam pasal 1 ayat II."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// // Get Curr Height
	// initialY = pdf.GetY()

	// // PASAL 2

	// pdf.SetFont("times", "B", 11)
	// pdf.Text(95, initialY+10, "PASAL 2")
	// pdf.Text(83, initialY+14, "PEMINDAHAN LOKASI")

	// height = initialY + 22

	// pdf.SetFont("times", "", 11)

	// pdf.Text(26, height, "1.")
	// pdf.SetXY(32, height-3.8)
	// paragraphWidth = 159.0
	// text = "Apabila Pihak Kedua akan memindahkan Mesin ke lokasi lainnya, maka Pihak Kedua harus memberitahukan secara tertulis kepada Pihak  Pertama  sekurang-kurangnya 1 (satu) minggu sebelumnya dan pemindahan tersebut atas persetujuan Pihak Pertama."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+15, "2.")
	// pdf.SetXY(32, height-3.8+15)
	// paragraphWidth = 159.0
	// text = "Untuk pemindahan Mesin tersebut, Pihak Pertama akan mengenakan biaya transportasi kepada pihak Kedua dan besar biaya tersebut disesuaikan dengan kesepakatan Para Pihak."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+25, "3.")
	// pdf.SetXY(32, height-3.8+25)
	// paragraphWidth = 159.0
	// text = "Segala akibat yang timbul dari pemindahan lokasi Mesin tersebut merupakan tanggung jawab pihak yang melaksanakan."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// // Get Curr Height
	// initialY = pdf.GetY()

	// // PASAL 3

	// pdf.SetFont("times", "B", 11)
	// pdf.Text(95, initialY+10, "PASAL 3")
	// pdf.Text(91, initialY+14, "HARGA SEWA")

	// height = initialY + 22

	// pdf.SetFont("times", "", 11)
	// pdf.Text(26, height, "Pihak Pertama menyewakan Mesin kepada Pihak Kedua dengan ketentuan Biaya sebagai berikut : ")
	// pdf.Text(26, height+5, "1.")
	// pdf.SetXY(32, height-3.8+5)
	// paragraphWidth = 159.0
	// formatRupiah := rupiah.FormatRupiah(kontrak.Details[0].Harga)
	// formatOverCopy := rupiah.FormatRupiah(float64(kontrak.Details[0].OverCopy))
	// text = "Biaya sewa Mesin untuk per unit " + kontrak.Details[0].Tipe + ", Speed " + kontrak.Details[0].Speed + " copy/menit adalah " + formatRupiah + "/bulan, dengan Free Copy " + strconv.Itoa(kontrak.Details[0].FreeCopy) + " copy dan pemakaian lebih dari " + strconv.Itoa(kontrak.Details[0].FreeCopy) + " copy akan dikenakan biaya pemakaian " + formatOverCopy + "/copy."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+20, "2.")
	// pdf.SetXY(32, height-3.8+20)
	// paragraphWidth = 159.0
	// text = "Pembayaran tagihan dapat dilakukan secara Cash atau ditransfer ke Rekening dengan data sebagai berikut :"
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(33, height+30, "-    Nama Rekening")
	// pdf.Text(33, height+35, "-    Bank")
	// pdf.Text(33, height+40, "-    Nomor Rekening")

	// pdf.Text(80, height+30, ":")
	// pdf.Text(80, height+35, ":")
	// pdf.Text(80, height+40, ":")
	// pdf.Text(83, height+30, bank.AtasNama)
	// pdf.Text(83, height+35, bank.Name)
	// pdf.Text(83, height+40, bank.Norek)

	// // Get Curr Height
	// initialY = pdf.GetY()

	// // PASAL 4

	// pdf.SetFont("times", "B", 11)
	// pdf.Text(95, initialY+27, "PASAL 4")
	// pdf.Text(60, initialY+31, "SERVICE MESIN DAN PENGGANTIAN SPAREPART")

	// height = initialY + 31

	// pdf.SetFont("times", "", 11)
	// pdf.Text(26, height+5, "Dalam hal ini service Mesin dan penggantian sparepart hanya  boleh  dilakukan  oleh  teknisi atau  petugas")
	// pdf.Text(26, height+10, "Pihak Pertama.")

	// pdf.Text(190, 274, "2")

	// pdf.SetFont("times", "I", 9)
	// pdf.Text(20, 274, "Perjanjian Sewa Menyewa Mesin Printer Multifunction")
	// pdf.Rect(142, 277, 50, 12, "D")
	// pdf.Line(167, 277, 167, 289)

	// pdf.SetFont("times", "UB", 7)
	// pdf.Text(143, 279.5, "Pihak Pertama")
	// pdf.Text(168, 279.5, "Pihak Kedua")

	// pdf.AddPage()

	// // Page 3
	// pdf.SetFont("times", "B", 11)
	// initialY = pdf.GetY()

	// // Pasal 5
	// pdf.Text(95, initialY+15, "PASAL 5")
	// pdf.Text(83, initialY+19, "CARA PEMBAYARAN")

	// height = initialY + 26

	// pdf.SetFont("times", "", 11)

	// pdf.Text(26, height, "1.")
	// pdf.SetXY(32, height-3.8)
	// paragraphWidth = 159.0
	// text = "Pihak Kedua akan membayar kepada Pihak Pertama setiap bulannya paling lambat 15 (Lima belas) hari kerja setelah kwitansi tagihan diterima berdasarkan biaya sewa mesin dan jumlah copy yang dibuat Pihak Kedua dalam periode 1 (satu) bulan."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+15, "2.")
	// pdf.SetXY(32, height-3.8+15)
	// paragraphWidth = 159.0
	// text = "Apabila Pihak Kedua belum melunasi tagihan setelah 30 (Tiga Puluh) hari kerja sejak kwitansi diterima sedangkan persyaratan Administrasi telah terpenuhi dengan benar maka akan dikenakan biaya keterlambatan pembayaran sebesar 1% (satu persen) dari total tagihan."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+31, "3.")
	// pdf.SetXY(32, height-3.8+31)
	// paragraphWidth = 159.0
	// text = "Dalam menghitung biaya penggandaan copy dalam satu bulan, jumlah copy dihitung berdasarkan counter pada saat pencatatan dan dikurangkan dengan counter pada catatan bulan sebelumnya."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+42, "4.")
	// pdf.SetXY(32, height-3.8+42)
	// paragraphWidth = 159.0
	// text = "Apabila Pihak Kedua belum melakukan pembayaran atas tagihan yang sudah jatuh tempo melebihi 30 (tiga puluh) hari sejak kwitansi diterima, maka Pihak Pertama berhak untuk memasang password pada Mesin sehingga Mesin dalam keadaan tidak aktif. Dalam hal ini Pihak Kedua tetap berkewajiban untuk membayar semua biaya sewa atas Mesin tersebut. Password akan dibuka oleh Pihak Pertama apabila Pihak Kedua sudah melunasi pembayaran yang sudah jatuh tempo."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+68, "5.")
	// pdf.SetXY(32, height-3.8+68)
	// paragraphWidth = 159.0
	// text = "Apabila Pihak Kedua belum melakukan pembayaran atas tagihan yang sudah jatuh tempo melebihi 60 (enam puluh) hari sejak kwitansi diterima, maka Pihak Pertama berhak untuk menarik Mesin  dan Pihak Kedua menyatakan tidak berkeberatan atas penarikan tersebut. Dalam hal ini maka semua tagihan yang belum dilunasi termasuk denda keterlambatan pembayaran tetap harus diselesaikan pembayarannya oleh Pihak Kedua selambat-lambatnya 7 (tujuh) hari kerja sejak mesin ditarik oleh Pihak Pertama."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+99, "6.")
	// pdf.SetXY(32, height-3.8+99)
	// paragraphWidth = 159.0
	// text = "Dalam hal terjadi penarikan Mesin oleh Pihak Pertama dari lokasi pihak Kedua maka tidak diperlukan lagi persetujuan lisan maupun tertulis dari Pihak Kedua, tetapi Pihak Kedua tetap berkewajiban untuk membantu kelancaran penarikan mesin apabila diperlukan surat pengantar dan atau surat keterangan pengeluaran Mesin dari lokasi Pihak Kedua."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.SetFont("times", "B", 11)
	// initialY = pdf.GetY()

	// // Pasal 6
	// pdf.Text(95, initialY+10, "PASAL 6")
	// pdf.Text(83, initialY+14, "PERUBAHAN HARGA")

	// height = initialY + 22

	// pdf.SetFont("times", "", 11)

	// pdf.Text(26, height, "1.")
	// pdf.SetXY(32, height-3.8)
	// paragraphWidth = 159.0
	// text = "Apabila terdapat perubahan nilai tukar rupiah (Rp) terhadap US Dollar atau perubahan kebijakan ekonomi dan atau politik yang mengakibatkan lonjakan harga suku cadang dan operasional maka, Pihak Pertama akan mengadakan perubahan atas syarat harga sewa dan biaya penggandaan dan Pihak Pertama akan memberitahukan kepada Pihak Kedua selambat-lambatnya 15 (lima belas) hari sebelum perubahan dilaksanakan. Apabila perubahan tersebut tidak disetujui, maka Pihak Kedua dapat memberitahukan secara tertulis kepada Pihak Pertama dengan disertakan rincian alasan penolakannya dalam jangka waktu 7 (tujuh) hari setelah menerima pemberitahuan perubahan harga."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+36, "2.")
	// pdf.SetXY(32, height-3.8+36)
	// paragraphWidth = 159.0
	// text = "Apabila Pihak Kedua tidak mengajukan keberatan secara tertulis seperti yang tertuang dalam Pasal 6 butir 1 maka Para Pihak menganggap bahwa perubahan harga tersebut telah disetujui."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.SetFont("times", "B", 11)
	// initialY = pdf.GetY()

	// // Pasal 7
	// pdf.Text(95, initialY+10, "PASAL 7")
	// pdf.Text(83, initialY+14, "JANGKA WAKTU SEWA")

	// height = initialY + 22

	// pdf.SetFont("times", "", 11)

	// pdf.Text(26, height, "1.")
	// pdf.SetXY(32, height-3.8)
	// paragraphWidth = 159.0

	// indonesianDateFormat := "02 January 2006"
	// periodeAwalString := kontrak.Details[0].PeriodeAwal.Format(indonesianDateFormat)
	// periodeAkhirString := kontrak.Details[0].PeriodeAkhir.Format(indonesianDateFormat)

	// monthReplacements := map[string]string{
	// 	"January":   "Januari",
	// 	"February":  "Februari",
	// 	"March":     "Maret",
	// 	"April":     "April",
	// 	"May":       "Mei",
	// 	"June":      "Juni",
	// 	"July":      "Juli",
	// 	"August":    "Agustus",
	// 	"September": "September",
	// 	"October":   "Oktober",
	// 	"November":  "November",
	// 	"December":  "Desember",
	// }

	// // Replace month names in the string
	// for enMonth, idMonth := range monthReplacements {
	// 	periodeAwalString = strings.Replace(periodeAwalString, enMonth, idMonth, -1)
	// 	periodeAkhirString = strings.Replace(periodeAkhirString, enMonth, idMonth, -1)
	// }

	// text = "Para Pihak setuju dan sepakat bahwa jangka waktu sewa kerjasama yang diatur dalam Perjanjian ini adalah 12 (dua belas) bulan terhitung sejak tanggal " + periodeAwalString + " sampai dengan " + periodeAkhirString + "."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+11, "2.")
	// pdf.SetXY(32, height-3.8+11)
	// paragraphWidth = 159.0
	// text = "Setelah jangka waktu sewa berakhir, Para Pihak setuju untuk memperpanjang secara otomatis (terus menerus) selama satu tahun berikutnya, kecuali salah satu pihak tidak berniat memperpanjang maka dapat dituangkan dalam surat konfirmasi penghentian kontrak."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+27, "3.")
	// pdf.SetXY(32, height-3.8+27)
	// paragraphWidth = 159.0
	// text = "Selama berlaku Perpanjangan Otomatis maka semua hak dan kewajiban dalam pasal-pasal surat perjanjian Nomor: SPK. " + kontrak.NoKontrak + " ini tetap berlaku dan mengikat para pihak."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(190, 274, "3")

	// pdf.SetFont("times", "I", 9)
	// pdf.Text(20, 274, "Perjanjian Sewa Menyewa Mesin Printer Multifunction")
	// pdf.Rect(142, 277, 50, 12, "D")
	// pdf.Line(167, 277, 167, 289)

	// pdf.SetFont("times", "UB", 7)
	// pdf.Text(143, 279.5, "Pihak Pertama")
	// pdf.Text(168, 279.5, "Pihak Kedua")
	// pdf.AddPage()

	// // Page 4
	// pdf.SetFont("times", "B", 11)
	// initialY = pdf.GetY()

	// // Pasal 8
	// pdf.Text(95, initialY+15, "PASAL 8")
	// pdf.Text(83, initialY+19, "PEMBATALAN SEWA")

	// height = initialY + 26

	// pdf.SetFont("times", "", 11)

	// pdf.Text(26, height, "1.")
	// pdf.SetXY(32, height-3.8)
	// paragraphWidth = 159.0
	// text = "Pihak Pertama dapat menghentikan  hubungan  sewa menyewa  bila  Pihak  Kedua  tidak  memenuhi kewajiban  dalam   membayar tagihan  selama 2 (dua) bulan berturut-turut dan Pihak Kedua berkewajiban untuk segera melunasi tagihan yang belum terbayar ditambah dengan denda."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+16, "2.")
	// pdf.SetXY(32, height-3.8+16)
	// paragraphWidth = 159.0
	// text = "Pihak Kedua dapat menghentikan hubungan sewa menyewa  apabila  Pihak  Pertama tidak dapat memenuhi kewajibannya."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+27, "3.")
	// pdf.SetXY(32, height-3.8+27)
	// paragraphWidth = 159.0
	// text = "Apabila Pihak Kedua menghentikan hubungan sewa menyewa di luar sebab yang tercantum dalam pasal 8 ayat 2, Pihak Kedua wajib memberitahukan secara tertulis kepada Pihak Pertama paling lambat 1 (satu) bulan sebelumnya dan kepada Pihak Kedua dikenakan biaya sebesar 1 (satu) bulan tagihan minimum. Untuk mesin dikembalikan kepada Pihak Pertama dan menjadi hak milik Pihak Pertama."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+53, "4.")
	// pdf.SetXY(32, height-3.8+53)
	// paragraphWidth = 159.0
	// text = "Apabila Pihak Pertama dan atau Pihak  Kedua  membatalkan  Perjanjian  ini,  maka Perjanjian ini  batal dengan  sendirinya  pada  saat  mesin  telah  ditarik kembali  dari   tempat   Pihak  Kedua.  Walaupun  mesin  tersebut telah  ditarik kembali, Pihak Kedua tetap terikat untuk memenuhi semua  kewajiban  membayar hutang-hutangnya dan perhitungan pembayaran sampai dengan  pemakaian  bulan berjalan. Mesin diambil maksimal 5 (lima) hari kerja setelah pemberitahuan  oleh Pihak Pertama kepada Pihak Kedua."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.SetFont("times", "B", 11)
	// initialY = pdf.GetY()

	// // Pasal 9
	// pdf.Text(95, initialY+10, "PASAL 9")
	// pdf.Text(83, initialY+14, "FORCE MAJEURE")

	// height = initialY + 22

	// pdf.SetFont("times", "", 11)

	// pdf.Text(26, height, "1.")
	// pdf.SetXY(32, height-3.8)
	// paragraphWidth = 159.0
	// text = "Yang dimaksud  dengan    force  majeure   dalam   Perjanjian  Sewa  Menyewa  ini adalah  peristiwa-peristiwa  yang  secara langsung mempengaruhi pelaksanaan Perjanjian  ini  dan  terjadi di luar kekuasaan dan kemampuan salah satu Pihak untuk mengatasinya dan dibenarkan  oleh  Pihak  lainnya  yaitu: bencana alam; wabah penyakit; huru-hara; perang; Kebakaran; dan peraturan pemerintah  mengenai keadaan bahaya, sehingga Para Pihak terpaksa tidak dapat memenuhi kewajibannya."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+31, "2.")
	// pdf.SetXY(32, height-3.8+31)
	// paragraphWidth = 159.0
	// text = "Diluar ketentuan tersebut di atas tidak dapat dikategorikan sebagai force majeure."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+37, "3.")
	// pdf.SetXY(32, height-3.8+37)
	// paragraphWidth = 159.0
	// text = "Dalam hal terjadi force majeure dimaksud Pasal 9 ayat 1., pihak yang tidak dapat melaksanakan kewajibannya, wajib memberitahukan secara tertulis kepada pihak lainnya dengan surat pernyataan kebenaran dari pejabat berwenang setempat dalam waktu 10 (sepuluh) hari sejak saat terjadinya   force majeure."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+58, "4.")
	// pdf.SetXY(32, height-3.8+58)
	// paragraphWidth = 159.0
	// text = "Kelalaian dan keterlambatan dalam memenuhi kewajiban pemberitahuan dimaksud Pasal 9 ayat 3, mengakibatkan tidak diakuinya peristiwa-peritiwa dimaksud Pasal 9 ayat I sebagai force majeure."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.SetFont("times", "B", 11)
	// initialY = pdf.GetY()

	// // Pasal 10
	// pdf.Text(95, initialY+10, "PASAL 10")
	// pdf.Text(75, initialY+14, "PENYELESAIAN PERSELISIHAN")

	// height = initialY + 22

	// pdf.SetFont("times", "", 11)

	// pdf.Text(26, height, "1.")
	// pdf.SetXY(32, height-3.8)
	// paragraphWidth = 159.0
	// text = "Jika terjadi selisih paham berkenaan dengan perjanjian ini antara Para  Pihak, maka selisih paham tersebut akan diselesaikan dengan cara musyawarah."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(26, height+11, "2.")
	// pdf.SetXY(32, height-3.8+11)
	// paragraphWidth = 159.0
	// text = "Apabila selisih paham tersebut dalam Pasal 10 ayat 1 di atas tidak dapat diselesaikan oleh Para Pihak, maka Para Pihak sepakat untuk memilih domisili hukum tetap di Kantor Kepaniteraan Pengadilan Negeri Jakarta Selatan."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.Text(190, 274, "4")

	// pdf.SetFont("times", "I", 9)
	// pdf.Text(20, 274, "Perjanjian Sewa Menyewa Mesin Printer Multifunction")
	// pdf.Rect(142, 277, 50, 12, "D")
	// pdf.Line(167, 277, 167, 289)

	// pdf.SetFont("times", "UB", 7)
	// pdf.Text(143, 279.5, "Pihak Pertama")
	// pdf.Text(168, 279.5, "Pihak Kedua")
	// pdf.AddPage()

	// // Page 5
	// pdf.SetFont("times", "B", 11)
	// initialY = pdf.GetY()

	// // Pasal 11
	// pdf.Text(95, initialY+15, "PASAL 11")
	// pdf.Text(91, initialY+19, "ADDENDUM")

	// height = initialY + 26

	// pdf.SetFont("times", "", 11)

	// pdf.SetXY(26, height-3.8)
	// paragraphWidth = 159.0
	// text = "Hal-hal yang belum cukup diatur dan perubahan-perubahan dalam perjanjian ini akan diatur kemudian atas dasar pemufakatan  Para Pihak yang akan dituangkan dalam bentuk Addendum yang merupakan kesatuan dan bagian yang tidak dapat dipisahkan dari Perjanjian ini."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// pdf.SetFont("times", "B", 11)
	// initialY = pdf.GetY()

	// // Pasal 12
	// pdf.Text(95, initialY+10, "PASAL 12")
	// pdf.Text(95, initialY+14, "PENUTUP")

	// height = initialY + 22

	// pdf.SetFont("times", "", 11)

	// pdf.SetXY(26, height-3.8)
	// paragraphWidth = 159.0
	// text = "Perjanjian ini dibuat dan ditandatangani oleh Para Pihak di Jakarta dalam rangkap 2 (dua) masing-masing bermaterai yang cukup dan memiliki kekuatan hukum yang sama serta segalanya dibuat dengan tujuan dan itikad baik."
	// pdf.MultiCell(paragraphWidth, 5, text, "", "J", false)

	// initialY = pdf.GetY()

	// height = initialY + 15

	// pdf.SetFont("times", "B", 10)
	// pdf.Text(45, height, "Pihak Pertama")
	// pdf.Text(130, height, "Pihak Kedua")

	// pdf.SetXY(30, height+1.2)
	// paragraphWidth = 50.0
	// text = strings.ToUpper(dataCompany.Name)
	// pdf.MultiCell(paragraphWidth, 5, text, "", "C", false)
	// // pdf.Text(38, height+5, strings.ToUpper(dataCompany.Name))

	// pdf.SetXY(100, height+1.2)
	// paragraphWidth = 80
	// text = strings.ToUpper(customer.CustomerName)
	// pdf.MultiCell(paragraphWidth, 5, text, "", "C", false)

	// pdf.SetFont("times", "BU", 10)
	// // pdf.Text(45, height+40, dataCompany.PIC)

	// pdf.SetXY(16, height+36.2)
	// paragraphWidth = 80
	// text = dataCompany.PIC
	// pdf.MultiCell(paragraphWidth, 5, text, "", "C", false)

	// pdf.SetXY(100, height+36.2)
	// paragraphWidth = 80
	// text = customer.PenandaTangan
	// pdf.MultiCell(paragraphWidth, 5, text, "", "C", false)

	// pdf.SetFont("times", "B", 10)
	// pdf.Text(50, height+44, "Direktur")

	// pdf.SetXY(100, height+40.2)
	// paragraphWidth = 80
	// text = customer.Jabatan
	// pdf.MultiCell(paragraphWidth, 5, text, "", "C", false)

	// pdf.Text(190, 274, "5")

	// pdf.SetFont("times", "I", 9)
	// pdf.Text(20, 274, "Perjanjian Sewa Menyewa Mesin Printer Multifunction")
	// pdf.Rect(142, 277, 50, 12, "D")
	// pdf.Line(167, 277, 167, 289)

	// pdf.SetFont("times", "UB", 7)
	// pdf.Text(143, 279.5, "Pihak Pertama")
	// pdf.Text(168, 279.5, "Pihak Kedua")

	err = pdf.Output(&b)
	if err != nil {
		return b, errors.Wrap(err, "[SERVICE][PrintKontrak]")
	}

	return b, nil
}

func numberToWords(num int) string {
	if num < 1 || num > 100 {
		return "Angka di luar jangkauan"
	}

	// Define words for numbers
	units := []string{"", "satu", "dua", "tiga", "empat", "lima", "enam", "tujuh", "delapan", "sembilan"}
	teens := []string{"sepuluh", "sebelas", "dua belas", "tiga belas", "empat belas", "lima belas", "enam belas", "tujuh belas", "delapan belas", "sembilan belas"}
	tens := []string{"", "", "dua puluh", "tiga puluh", "empat puluh", "lima puluh", "enam puluh", "tujuh puluh", "delapan puluh", "sembilan puluh"}

	if num == 100 {
		return "seratus"
	}

	var words string
	if num >= 20 {
		words = tens[num/10] + " "
		num %= 10
	}

	if num >= 10 {
		words = teens[num-10] + " "
		num = 0
	}

	if num > 0 {
		words += units[num] + " "
	}

	return strings.TrimSpace(words)
}

// func formatNumberWithDotSeparator(number int) string {
// 	// Convert the number to a string
// 	numStr := strconv.Itoa(number)

// 	// Reverse the string to make it easier to insert separators
// 	reversed := reverseString(numStr)

// 	var result strings.Builder
// 	for i, c := range reversed {
// 		if i > 0 && i%3 == 0 {
// 			result.WriteRune('.')
// 		}
// 		result.WriteRune(c)
// 	}

// 	// Reverse the string back to its original order
// 	return reverseString(result.String())
// }

// // Helper function to reverse a string
// func reverseString(s string) string {
// 	runes := []rune(s)
// 	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
// 		runes[i], runes[j] = runes[j], runes[i]
// 	}
// 	return string(runes)
// }

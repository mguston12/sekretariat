package sekretariat

import (
	"context"
	"fmt"
	"math"
	"sekretariat/internal/entity/sekretariat"
	"sekretariat/pkg/errors"
	"time"
)

func (s Service) GetAllContractsHeader(ctx context.Context, company int, keyword string, page, length int) ([]sekretariat.KontrakHeader, int, error) {
	limit := length
	offset := (page - 1) * length
	var lastPage int

	if page != 0 && length != 0 {
		headers, count, err := s.data.GetAllContractsHeaderCount(ctx, company)
		if err != nil {
			return headers, lastPage, errors.Wrap(err, "[SERVICE][GetAllContractsHeader][COUNT]")

		}
		lastPage = int(math.Ceil(float64(count) / float64(length)))

		headers, err = s.data.GetAllContractsHeaderPage(ctx, company, keyword, offset, limit)
		if err != nil {
			return headers, lastPage, errors.Wrap(err, "[SERVICE][GetAllContractsHeader][PAGE]")
		}
		return headers, lastPage, nil
	}

	headers, err := s.data.GetAllContractsHeader(ctx, company, keyword)
	if err != nil {
		return headers, lastPage, errors.Wrap(err, "[SERVICE][GetAllContractsHeader]")
	}
	return headers, lastPage, nil
}

func (s Service) GetDataContractByContractNumber(ctx context.Context, company int, no_kontrak string) (sekretariat.KontrakHeader, error) {
	header, err := s.data.GetContractsHeaderByContractNumber(ctx, company, no_kontrak)
	if err != nil {
		return header, errors.Wrap(err, "[SERVICE][GetDataContractByContractNumber]")
	}

	details, err := s.data.GetContractDetailsByContractNumber(ctx, no_kontrak)
	if err != nil {
		return header, errors.Wrap(err, "[SERVICE][GetDataContractByContractNumber]")
	}

	detailMap := make(map[string][]sekretariat.KontrakDetail)
	for _, detail := range details {
		detailMap[detail.NoKontrak] = append(detailMap[detail.NoKontrak], detail)
	}

	header.Details = append(header.Details, detailMap[header.NoKontrak]...)

	return header, nil
}

// func (s Service) GetContractFiltered(ctx context.Context, company int, id, name string) ([]sekretariat.Kontrak, error) {
// 	contract, err := s.data.GetContractFiltered(ctx, company, id, name)
// 	if err != nil {
// 		return contract, errors.Wrap(err, "[SERVICE][GetContractFiltered]")
// 	}
// 	return contract, nil
// }

func (s Service) CreateContract(ctx context.Context, header sekretariat.KontrakHeader) error {
	err := s.data.CreateContractHeader(ctx, header)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateContract][Header]")
	}

	for _, detail := range header.Details {
		detail.NoKontrak = header.NoKontrak
		detail.UpdatedBy = header.UpdatedBy

		layoutFormat := "2006-01-02"
		_periodeAwal, _ := time.Parse(layoutFormat, detail.PeriodeAwalString)
		_periodeAkhir, _ := time.Parse(layoutFormat, detail.PeriodeAkhirString)

		detail.PeriodeAwal = _periodeAwal
		detail.PeriodeAkhir = _periodeAkhir

		err := s.data.CreateContractDetail(ctx, detail)
		if err != nil {
			return errors.Wrap(err, "[SERVICE][CreateContract][Detail]")
		}
	}

	err = s.data.IncreaseCounterContract(ctx, header.CompanyID)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateContract][IncreaseCounter]")
	}

	return nil
}

func (s Service) UpdateContract(ctx context.Context, header sekretariat.KontrakHeader) error {
	err := s.data.UpdateContractHeader(ctx, header)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][UpdateContract][1]")
	}

	err = s.data.DeleteContractDetail(ctx, header.NoKontrak)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][UpdateContract][2]")
	}

	for _, detail := range header.Details {
		detail.NoKontrak = header.NoKontrak
		detail.UpdatedBy = header.UpdatedBy

		if detail.PeriodeAwal.IsZero() {
			layoutFormat := "2006-01-02"

			_periodeAwal, _ := time.Parse(layoutFormat, detail.PeriodeAwalString)
			_periodeAkhir, _ := time.Parse(layoutFormat, detail.PeriodeAkhirString)

			detail.PeriodeAwal = _periodeAwal
			detail.PeriodeAkhir = _periodeAkhir
		}

		err = s.data.CreateContractDetail(ctx, detail)
		if err != nil {
			return errors.Wrap(err, "[SERVICE][UpdateContract][3]")
		}
	}

	return nil
}

func (s Service) GetCounterContract(ctx context.Context, company int) (string, error) {
	var com string
	var month string

	switch company {
	case 1:
		com = "PB"
	case 2:
		com = "PBM"
	case 3:
		com = "MMU"
	}

	switch int(time.Now().Month()) {
	case 1:
		month = "I"
	case 2:
		month = "II"
	case 3:
		month = "III"
	case 4:
		month = "IV"
	case 5:
		month = "V"
	case 6:
		month = "VI"
	case 7:
		month = "VII"
	case 8:
		month = "VIII"
	case 9:
		month = "IX"
	case 10:
		month = "X"
	case 11:
		month = "XI"
	case 12:
		month = "XII"
	}

	id, err := s.data.GetCounterContract(ctx, company)
	if err != nil {
		return "", errors.Wrap(err, "[SERVICE][GetCounterContract]")
	}

	return fmt.Sprintf("%04d/%s/%s/%d", id, com, month, time.Now().Year()), nil
}

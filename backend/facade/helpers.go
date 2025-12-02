package facade

import (
	"github.com/Station-Manager/errors"
	"github.com/Station-Manager/types"
)

func mergeCountryIntoContactedStation(station *types.ContactedStation, country types.Country) error {
	const op errors.Op = "facade.mergeCountryIntoContactedStation"
	if station == nil {
		return errors.New(op).Msg("Contacted station parameter is nil")
	}
	station.Country = country.Name
	station.Cont = country.Continent

	return nil
}

func mergeIntoQso(qso *types.Qso, country types.Country) error {
	const op errors.Op = "facade.mergeCountryIntoQso"
	if qso == nil {
		return errors.New(op).Msg("QSO parameter is nil")
	}
	qso.CountryDetails = country

	return nil
}

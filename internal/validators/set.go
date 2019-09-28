package validators

import (
	"github.com/google/uuid"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/dispatcher/common"
	"github.com/ttacon/libphonenumber"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

type ValidatorSet struct {
}

var (
	availablePositions = map[string]bool{
		common.UserProfilePositionCEO:               true,
		common.UserProfilePositionCTO:               true,
		common.UserProfilePositionCMO:               true,
		common.UserProfilePositionCFO:               true,
		common.UserProfilePositionProjectManagement: true,
		common.UserProfilePositionGenericManagement: true,
		common.UserProfilePositionSoftwareDeveloper: true,
		common.UserProfilePositionMarketing:         true,
		common.UserProfilePositionSupport:           true,
	}

	availableAnnualIncome = []*billing.RangeInt{
		{From: 0, To: 1000},
		{From: 1000, To: 10000},
		{From: 10000, To: 100000},
		{From: 100000, To: 1000000},
		{From: 1000000, To: 0},
	}

	availableNumberOfEmployees = []*billing.RangeInt{
		{From: 1, To: 10},
		{From: 11, To: 50},
		{From: 51, To: 100},
		{From: 100, To: 0},
	}

	availableTariffPaymentAmountRange = []*billing.PriceTableCurrency{
		{From: 0.75, To: 5},
	}

	zipUsaRegexp      = regexp.MustCompile("^[0-9]{5}(?:-[0-9]{4})?$")
	nameRegexp        = regexp.MustCompile("^[\\p{L}\\p{M} \\-\\']+$")
	companyNameRegexp = regexp.MustCompile("^[\\p{L}\\p{M} \\-\\.0-9]+$")
	zipGeneralRegexp  = regexp.MustCompile("^\\d{0,30}$")
	swiftRegexp       = regexp.MustCompile("^[A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?$")
	cityRegexp        = regexp.MustCompile("^[\\p{L}\\p{M} \\-\\.]+$")
)

// PhoneValidator
func (v *ValidatorSet) PhoneValidator(fl validator.FieldLevel) bool {
	_, err := libphonenumber.Parse(fl.Field().String(), "US")
	return err == nil
}

// UuidValidator
func (v *ValidatorSet) UuidValidator(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())
	return err == nil
}

// ZipUsaValidator
func (v *ValidatorSet) ZipUsaValidator(fl validator.FieldLevel) bool {
	return zipUsaRegexp.MatchString(fl.Field().String())
}

// NameValidator
func (v *ValidatorSet) NameValidator(fl validator.FieldLevel) bool {
	return nameRegexp.MatchString(fl.Field().String())
}

// PositionValidator
func (v *ValidatorSet) PositionValidator(fl validator.FieldLevel) bool {
	_, ok := availablePositions[fl.Field().String()]
	return ok
}

// CompanyValidator
func (v *ValidatorSet) CompanyValidator(sl validator.StructLevel) {
	company := sl.Current().Interface().(grpc.UserProfileCompany)
	res := v.RangeIntValidator(company.AnnualIncome, availableAnnualIncome)

	if res == false {
		sl.ReportError(company.AnnualIncome, "AnnualIncome", "annual_income", "annual_income", "")
	}

	res = v.RangeIntValidator(company.NumberOfEmployees, availableNumberOfEmployees)

	if res == false {
		sl.ReportError(company.NumberOfEmployees, "NumberOfEmployees", "number_of_employees", "number_of_employees", "")
	}
}

// MerchantTariffRatesValidator
func (v *ValidatorSet) MerchantTariffRatesValidator(sl validator.StructLevel) {
	tariff := sl.Current().Interface().(grpc.GetMerchantTariffRatesRequest)

	if tariff.AmountFrom <= 0 && tariff.AmountTo <= 0 {
		return
	}

	res := v.RangeFloatValidator(
		&billing.PriceTableCurrency{
			From: tariff.AmountFrom,
			To:   tariff.AmountTo,
		},
		availableTariffPaymentAmountRange,
	)

	if res == false {
		sl.ReportError(tariff.AmountFrom, "AmountFrom", "amount_from", "amount_from", "")
	}
}

// RangeIntValidator
func (v *ValidatorSet) RangeIntValidator(in *billing.RangeInt, rng []*billing.RangeInt) bool {
	for _, v := range rng {
		if in.From == v.From && in.To == v.To {
			return true
		}
	}

	return false
}

// RangeFloatValidator
func (v *ValidatorSet) RangeFloatValidator(in *billing.PriceTableCurrency, rng []*billing.PriceTableCurrency) bool {
	for _, v := range rng {
		if in.From == v.From && in.To == v.To {
			return true
		}
	}

	return false
}

// CompanyNameValidator
func (v *ValidatorSet) CompanyNameValidator(fl validator.FieldLevel) bool {
	return companyNameRegexp.MatchString(fl.Field().String())
}

// MerchantCompanyValidator
func (v *ValidatorSet) MerchantCompanyValidator(sl validator.StructLevel) {
	company := sl.Current().Interface().(billing.MerchantCompanyInfo)

	reg := zipGeneralRegexp

	if v, ok := common.ZipRegexp[company.Country]; ok {
		reg = v
	}

	match := reg.MatchString(company.Zip)

	if !match {
		sl.ReportError(company.Zip, "Zip", "zip", "zip", "")
	}
}

// SwiftValidator
func (v *ValidatorSet) SwiftValidator(fl validator.FieldLevel) bool {
	return swiftRegexp.MatchString(fl.Field().String())
}

// CityValidator
func (v *ValidatorSet) CityValidator(fl validator.FieldLevel) bool {
	return cityRegexp.MatchString(fl.Field().String())
}

// WorldRegionValidator
func (v *ValidatorSet) WorldRegionValidator(fl validator.FieldLevel) bool {
	_, ok := common.TariffRegions[fl.Field().String()]
	return ok
}

// New
func New() *ValidatorSet {
	return &ValidatorSet{}
}

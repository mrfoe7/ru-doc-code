package kpp

import (
	ru_doc_code "github.com/mrfoe7/ru-doc-code"
)

const kppLength = 9

type KPP struct {
	Code         ru_doc_code.TaxRegionCode          `code:"length=4" valid:"runelength4_custom"`
	Reason       ru_doc_code.RegistrationReasonCode `code:"length=2" valid:"runelength2_custom"`
	SerialNumber ru_doc_code.SerialNumber           `code:"length=3" valid:"runelength3_custom"`
}

// Validate check to valid KPP format
// example: input format is 773643301
func Validate(kpp string) (bool, error) {
	if len(kpp) != kppLength {
		name, err := ru_doc_code.GetModuleName()
		if err != nil {
			return false, err
		}
		return false, &ru_doc_code.CommonError{
			Method: name,
			Err:    ru_doc_code.ErrInvalidLength,
		}
	}

	_, err := ru_doc_code.StrToArr(kpp)
	if err != nil {
		return false, err
	}

	var value KPP
	err = ru_doc_code.Unmarshal([]byte(kpp), &value)
	if err != nil {
		return false, err
	}

	// todo: validate tax region/office ru_doc_code.TaxRegionCode(kpp[:4])

	_, ok := ru_doc_code.SupportedRegistrationReasonSet[value.Reason]
	if !ok {
		return false, ru_doc_code.ErrRegistrationReasonCode
	}

	return true, nil
}

// DOES NOT IMPLEMENTED. Generate generate random type kpp string value.
func Generate() string {
	return ""
}

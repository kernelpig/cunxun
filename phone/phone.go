package phone

import (
	"fmt"
	"strconv"
	"strings"

	libphone "github.com/ttacon/libphonenumber"

	e "wangqingang/cunxun/error"
)

type PhoneInfo struct {
	CountryCode    string
	NationalNumber string
}

// var v string = "86 13681454478"
func ValidPhone(str string) (*PhoneInfo, error) {
	if str == "" {
		return nil, e.S(e.MPhoneErr, e.PhoneEmpty)
	}

	pieces := strings.Split(str, " ")
	if len(pieces) < 2 {
		detail := fmt.Sprintf("phone format error, expect {CC}{Space}{Number}")
		return nil, e.SD(e.MPhoneErr, e.PhoneFormatErr, detail)
	}

	countryCodeStr := pieces[0]
	countryCode, err := strconv.Atoi(strings.TrimLeft(countryCodeStr, "+"))
	if err != nil {
		return nil, e.SP(e.MPhoneErr, e.PhoneInvalidCountryCode, err)
	}

	region := libphone.GetRegionCodeForCountryCode(countryCode)
	if region == libphone.UNKNOWN_REGION {
		return nil, e.S(e.MPhoneErr, e.PhoneUnknownRegion)
	}

	phoneNumber, err := libphone.Parse(str, region)
	if err != nil {
		return nil, e.SP(e.MPhoneErr, e.PhoneParseNumberErr, err)
	}

	if !libphone.IsValidNumberForRegion(phoneNumber, region) {
		return nil, e.SP(e.MPhoneErr, e.PhoneRegionMismatch, err)
	}

	phoneInfo := &PhoneInfo{
		CountryCode:    fmt.Sprintf("%d", phoneNumber.GetCountryCode()),
		NationalNumber: fmt.Sprintf("%d", phoneNumber.GetNationalNumber()),
	}
	return phoneInfo, nil
}

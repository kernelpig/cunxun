package phone

import (
	"strconv"
	"strings"

	libphone "github.com/ttacon/libphonenumber"

	"fmt"
	e "wangqingang/cunxun/error"
)

// var v string = "86 13681454478"
func ValidPhone(str string) error {
	if str == "" {
		return e.SP(e.MPhoneErr, e.PhoneEmpty, nil)
	}

	pieces := strings.Split(str, " ")
	if len(pieces) < 2 {
		return e.SP(e.MPhoneErr, e.PhoneFormatErr, fmt.Errorf("phone format error, expect {CC}{Space}{Number}"))
	}

	countryCodeStr := pieces[0]
	countryCode, err := strconv.Atoi(strings.TrimLeft(countryCodeStr, "+"))
	if err != nil {
		return e.SP(e.MPhoneErr, e.PhoneInvalidCountryCode, err)
	}

	region := libphone.GetRegionCodeForCountryCode(countryCode)
	if region == libphone.UNKNOWN_REGION {
		return e.SP(e.MPhoneErr, e.PhoneUnknownRegion, err)
	}

	phoneNumber, err := libphone.Parse(str, region)
	if err != nil {
		return e.SP(e.MPhoneErr, e.PhoneParseNumberErr, err)
	}

	if !libphone.IsValidNumberForRegion(phoneNumber, region) {
		return e.SP(e.MPhoneErr, e.PhoneRegionMismatch, err)
	}

	return nil
}

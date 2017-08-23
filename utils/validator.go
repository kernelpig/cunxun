package utils

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"sync"

	libphone "github.com/ttacon/libphonenumber"
	validator "gopkg.in/go-playground/validator.v8"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var structValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}
	return nil
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		config := &validator.Config{TagName: "binding"}
		v.validate = validator.New(config)
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

// var v string = "+86 13681454478"
func ValidPhone(str string) error {
	if str == "" {
		return errors.New("empty string")
	}

	pieces := strings.Split(str, " ")
	if len(pieces) < 2 {
		return errors.New("phone format error, expect {+CC}{Space}{Number}")
	}

	countryCodeStr := pieces[0]
	countryCode, err := strconv.Atoi(strings.TrimLeft(countryCodeStr, "+"))
	if err != nil {
		return errors.New("invalid country code")
	}

	region := libphone.GetRegionCodeForCountryCode(countryCode)
	if region == libphone.UNKNOWN_REGION {
		return errors.New("unknown country code region")
	}

	phoneNumber, err := libphone.Parse(str, region)
	if err != nil {
		return errors.New("parse phone number failed")
	}

	if !libphone.IsValidNumberForRegion(phoneNumber, region) {
		return errors.New("mismatch phone region")
	}

	return nil
}

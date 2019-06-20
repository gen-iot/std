package std

import (
	"bytes"
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	zhTranslations "gopkg.in/go-playground/validator.v9/translations/zh"
)

/**
 * Created by xuchao on 2019-03-06 .
 */

// use a single instance , it caches struct info
var (
	uni       *ut.UniversalTranslator
	gValidate *validator.Validate
)

type LANG int

const (
	LANG_ZH LANG = iota
	LANG_EN
)

func init() {
	zh := zh.New()
	en := en.New()
	uni = ut.New(zh, zh, en)
	gValidate = validator.New()
	trans, found := uni.GetTranslator("zh")
	if found {
		if err := zhTranslations.RegisterDefaultTranslations(gValidate, trans); err != nil {
			panic(err)
		}
	}
}
func Str2Lang(lang string) LANG {
	switch lang {
	case "zh":
		return LANG_ZH
	case "cn":
		return LANG_EN
	}
	return LANG_ZH
}

func Lang2Str(lang LANG) string {

	switch lang {
	case LANG_ZH:
		return "zh"
	case LANG_EN:
		return "cn"
	}
	return "zh"
}

func ValidateStructWithLanguage(lang LANG, i interface{}) error {
	e := gValidate.Struct(i)
	if e != nil {
		if _, ok := e.(*validator.InvalidValidationError); ok {
			return e
		}
		// translate all error at once
		var buffer bytes.Buffer
		rawErrs := e.(validator.ValidationErrors)
		trans, found := uni.GetTranslator(Lang2Str(lang))
		if found {
			tansErrs := rawErrs.Translate(trans)
			for _, err := range tansErrs {
				buffer.WriteString(fmt.Sprintf("%s\n", err))
			}
		} else {
			for _, err := range rawErrs {
				buffer.WriteString(fmt.Sprintf("param:'%s' type:'%s' miss match with check:'%s'\n", err.Field(), err.Kind(), err.Tag()))
			}
		}
		return errors.New(buffer.String())
	}
	return nil
}

func ValidateStruct(i interface{}) error {
	return ValidateStructWithLanguage(LANG_ZH, i)
}

func Verify(i interface{}) error {
	return ValidateStructWithLanguage(LANG_ZH, i)
}

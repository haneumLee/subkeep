package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// validate is the singleton validator instance.
var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Register custom validations.
	_ = validate.RegisterValidation("billing_cycle", validateBillingCycle)
	_ = validate.RegisterValidation("currency_krw", validateCurrencyKRW)
}

// GetValidator returns the singleton validator instance.
func GetValidator() *validator.Validate {
	return validate
}

// ValidateStruct validates a struct and returns a formatted error message.
func ValidateStruct(s interface{}) *AppError {
	if err := validate.Struct(s); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return ErrValidation(err.Error())
		}

		messages := make([]string, 0, len(validationErrors))
		for _, fe := range validationErrors {
			messages = append(messages, formatFieldError(fe))
		}

		return ErrValidation(strings.Join(messages, "; "))
	}
	return nil
}

// formatFieldError converts a FieldError into a human-readable message.
func formatFieldError(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s 필드는 필수입니다", field)
	case "email":
		return fmt.Sprintf("%s 필드는 유효한 이메일이어야 합니다", field)
	case "min":
		return fmt.Sprintf("%s 필드는 최소 %s 이상이어야 합니다", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s 필드는 최대 %s 이하여야 합니다", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s 필드는 %s 이상이어야 합니다", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s 필드는 %s 이하여야 합니다", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s 필드는 [%s] 중 하나여야 합니다", field, fe.Param())
	case "url":
		return fmt.Sprintf("%s 필드는 유효한 URL이어야 합니다", field)
	case "billing_cycle":
		return fmt.Sprintf("%s 필드는 monthly, yearly, weekly 중 하나여야 합니다", field)
	case "currency_krw":
		return fmt.Sprintf("%s 필드는 KRW이어야 합니다", field)
	default:
		return fmt.Sprintf("%s 필드의 유효성 검증에 실패했습니다 (%s)", field, fe.Tag())
	}
}

// validateBillingCycle checks that the billing cycle is one of the allowed values.
func validateBillingCycle(fl validator.FieldLevel) bool {
	cycle := fl.Field().String()
	switch cycle {
	case "monthly", "yearly", "weekly":
		return true
	default:
		return false
	}
}

// validateCurrencyKRW checks that the currency is KRW.
func validateCurrencyKRW(fl validator.FieldLevel) bool {
	return fl.Field().String() == "KRW"
}

// MonthlyAmount converts an amount to its monthly equivalent based on the billing cycle.
// 금액 환산 규칙 (FRS F-03):
//   - monthly: 그대로
//   - yearly:  amount / 12 (반올림)
//   - weekly:  amount × 52 / 12 (반올림)
func MonthlyAmount(amount int64, billingCycle string) int64 {
	switch billingCycle {
	case "yearly":
		return roundDiv(amount, 12)
	case "weekly":
		return roundDiv(amount*52, 12)
	default: // monthly
		return amount
	}
}

// roundDiv performs integer division with rounding to nearest.
func roundDiv(numerator, denominator int64) int64 {
	if denominator == 0 {
		return 0
	}
	return (numerator + denominator/2) / denominator
}

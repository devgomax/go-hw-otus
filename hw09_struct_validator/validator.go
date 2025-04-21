package hw09structvalidator

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// ValidationError структура, содержащая имя поля и ошибку его валидации.
type ValidationError struct {
	Field string
	Err   error
}

// ValidationErrors слайс структур ValidationError для хранения всех ошибок валидации.
type ValidationErrors []ValidationError

// Error реализует интерфейс error. Возвращает текст всех ошибок валидации.
func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	slices.SortFunc(v, func(a, b ValidationError) int {
		return cmp.Compare(a.Field, b.Field)
	})

	messages := make([]string, 0, len(v))
	for _, ve := range v {
		messages = append(messages, fmt.Sprintf("field %q: %v", ve.Field, ve.Err))
	}

	return strings.Join(messages, "\n")
}

var (
	ErrInvalidValueType = errors.New("invalid tag value type")
	ErrInvalidRegexp    = errors.New("can't compile regexp")
	ErrInvalidFormat    = errors.New("validator doesn't match the format {identifier:value}")
	ErrUnknownValidator = errors.New("unknown validator")
)

// Validate валидирует поля структуры на основе тега "validate".
func Validate(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	var validationErrors ValidationErrors

	t := val.Type()

	for i := range t.NumField() {
		field := val.Field(i)

		tag := t.Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}

		if err := validateField(field, t.Field(i).Name, tag); err != nil {
			var ves ValidationErrors
			if errors.As(err, &ves) {
				validationErrors = append(validationErrors, ves...)
				continue
			}

			return fmt.Errorf("field validation failed at %q: %w", t.Field(i).Name, err)
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateField(fieldValue reflect.Value, fieldName, tag string) error {
	switch kind := fieldValue.Kind(); {
	case kind == reflect.Pointer:
		return validateField(fieldValue.Elem(), fieldName, tag)
	case reflect.Int <= kind && kind <= reflect.Int64:
		return validateNumber(fieldValue.Int(), tag, fieldName)
	case reflect.Uint <= kind && kind <= reflect.Uintptr:
		return validateNumber(fieldValue.Uint(), tag, fieldName)
	case kind == reflect.Float32, kind == reflect.Float64:
		return validateNumber(fieldValue.Float(), tag, fieldName)
	case kind == reflect.String:
		return validateString(fieldValue.String(), tag, fieldName)
	case kind == reflect.Slice:
		return validateSlice(fieldValue, tag, fieldName)
	case kind == reflect.Struct:
		if strings.Contains(tag, "nested") {
			return Validate(fieldValue.Interface())
		}
		return nil
	default:
		return nil
	}
}

// Numbers ограничение на все целочисленные типы и типы с плавающей точкой.
type Numbers interface {
	constraints.Integer | constraints.Float
}

func validateSlice(value reflect.Value, tag, fieldName string) error {
	var validationErrors ValidationErrors

	for i := range value.Len() {
		element := value.Index(i)
		elemName := fmt.Sprintf("%s[%d]", fieldName, i)

		if err := validateField(element, elemName, tag); err != nil {
			var ves ValidationErrors
			if errors.As(err, &ves) {
				validationErrors = append(validationErrors, ves...)
				continue
			}

			return err
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateNumber[T Numbers](value T, tag string, fieldName string) error {
	validators, err := getValidators(tag)
	if err != nil {
		return err
	}

	var validationErrors ValidationErrors

	for k, v := range validators {
		switch k {
		case "min":
			limit, inErr := strconv.ParseFloat(v, 64)
			if inErr != nil {
				return errors.Join(ErrInvalidValueType, fmt.Errorf(`value of %q must be a number: %w`, k, inErr))
			}

			if cmp.Less(float64(value), limit) {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value must be >= %v", limit),
				})
			}
		case "max":
			limit, inErr := strconv.ParseFloat(v, 64)
			if inErr != nil {
				return errors.Join(ErrInvalidValueType, fmt.Errorf(`value of %q must be a number: %w`, k, inErr))
			}

			if cmp.Less(limit, float64(value)) {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value must be <= %v", limit),
				})
			}
		case "in":
			strnums := strings.Split(v, ",")
			if len(strnums) == 0 {
				return nil
			}

			set := make([]float64, 0, len(strnums))
			for _, n := range strnums {
				num, inErr := strconv.ParseFloat(n, 64)
				if inErr != nil {
					return errors.Join(ErrInvalidValueType, fmt.Errorf(`all values of %q must be numbers: %w`, k, inErr))
				}

				set = append(set, num)
			}

			if !slices.Contains(set, float64(value)) {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value must be one of [%s]", v),
				})
			}
		default:
			return ErrUnknownValidator
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateString[T interface{ ~string }](value T, tag string, fieldName string) error {
	validators, err := getValidators(tag)
	if err != nil {
		return err
	}

	var validationErrors ValidationErrors

	for k, v := range validators {
		switch k {
		case "len":
			limit, err := strconv.Atoi(v)
			if err != nil {
				return errors.Join(ErrInvalidValueType, fmt.Errorf(`value of %q must be a number: %w`, k, err))
			}

			if len(value) != limit {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value must be a string of length %v", limit),
				})
			}
		case "regexp":
			pattern, err := regexp.Compile(v)
			if err != nil {
				return errors.Join(ErrInvalidRegexp, fmt.Errorf(`value of %q must be a valid regexp string: %w`, k, err))
			}

			if !pattern.MatchString(string(value)) {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value must match a regex %q", pattern.String()),
				})
			}
		case "in":
			strs := strings.Split(v, ",")
			if len(strs) == 0 {
				return nil
			}

			if !slices.Contains(strs, string(value)) {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf(`value must be one of [%s]`, v),
				})
			}
		default:
			return ErrUnknownValidator
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func getValidators(tag string) (map[string]string, error) {
	validators := strings.Split(tag, "|")
	m := make(map[string]string, len(validators))

	for _, validator := range validators {
		parts := strings.Split(validator, ":")
		if len(parts) != 2 {
			return nil, ErrInvalidFormat
		}

		m[parts[0]] = parts[1]
	}

	return m, nil
}

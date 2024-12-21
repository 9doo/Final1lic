
package calc_test

import (
	"testing"

	"github.com/9doo/Final1lic/pkg/calculation"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calc.Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error", testCase.expression) // ошибок не должно быть
			}
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult) // неверный ответ
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:        "invalid_operator",
			expression:  "1+1*",
			expectedErr: calc.ErrInvalidExpression,
		},
		{
			name:        "invalid_operator",
			expression:  "2+2**2",
			expectedErr: calc.ErrInvalidExpression,
		},
		{
			name:        "invalid_parentheses",
			expression:  "((2+2-*(2",
			expectedErr: calc.ErrInvalidExpression,
		},
		{
			name:        "empty_expression",
			expression:  "",
			expectedErr: calc.ErrEmptyExpression,
		},
		{
			name:        "division_by_zero",
			expression:  "1/0",
			expectedErr: calc.ErrDivisionByZero,
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calc.Calc(testCase.expression)
			if err == nil {
				t.Fatalf("expression %s is invalid but result %f was obtained", testCase.expression, val) // ошибка отсутствует
			}

			if err != testCase.expectedErr {
				t.Fatalf("expected error %v, got %v", testCase.expectedErr, err) // проверка, что ошибка соответствует ожидаемой
			}
		})
	}
}

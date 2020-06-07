package lab1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"errors"
)

func TestCalculatePostfix(t *testing.T) {

	tests := []struct {
		input string
		exp int
		expErr error
	}{
		{"1 2 + 3 * 4 - 5 /", 1, nil},
		{"2 3 ^ 1 +", 9, nil},
		{"7 3 + 5 / 2 4 ^ * 8 - 6 + 1 +", 31, nil},
		{"4 2 - 3 * 5 +", 11, nil},
		{"6 2 ^ 9 / 1 * 2 3 * - 3 + 8 +", 9, nil},
		{"15&*", 0, errors.New("Invalid element of expression.")},
		{"", 0, errors.New("Expression is not complete.")},
		{"+ 5 6 7", 0, errors.New("Expression is wrong.")},

	}

	for _, test := range tests{
		res, err := CalculatePostfix(test.input)
		assert.Equal(t, res, test.exp)
		assert.Equal(t, err, test.expErr)
	}
}

func ExampleCalculatePostfix() {
	res, _ := CalculatePostfix("2 2 + 2 *")
	fmt.Println(res)

	// Output:
	// 8
}

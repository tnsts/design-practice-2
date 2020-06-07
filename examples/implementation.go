package lab1

import (
		"math"
		"strings"
		"strconv"
		"errors"
	)

func CalculatePostfix(input string) (int, error) {
	expression := strings.Fields(input)
	var stack []int

	for _, element := range expression {
		if element == "+" || element == "-" || element == "/" || element == "*" || element == "^" {
			if len(stack) < 2 {
				return 0, errors.New("Expression is wrong.")
			}

			var values []int
			stack, values = stack[:len(stack) - 2], stack[len(stack) - 2:]
			a, b := values[0], values[1]

			switch element {
			case "+":
				stack = append(stack, a + b)
			case "-":
				stack = append(stack, a - b)
			case "*":
				stack = append(stack, a * b)
			case "/":
				stack = append(stack, a / b)
			case "^":
				stack = append(stack, int(math.Pow(float64(a), float64(b))))
			}
		} else {
			value, err := strconv.Atoi(element)
			if err != nil {
				return 0, errors.New("Invalid element of expression.")
			}

			stack = append(stack, value)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("Expression is not complete.")
	}

	return stack[0], nil
}

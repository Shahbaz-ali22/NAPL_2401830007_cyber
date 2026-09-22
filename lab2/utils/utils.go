// Package utils provides small string and mathematical helper functions.
package utils

import "unicode"

// Reverse returns s with its Unicode characters in reverse order.
func Reverse(s string) string {
	runes := []rune(s)
	for left, right := 0, len(runes)-1; left < right; left, right = left+1, right-1 {
		runes[left], runes[right] = runes[right], runes[left]
	}
	return string(runes)
}

// CountVowels returns the number of English vowels in s.
func CountVowels(s string) int {
	count := 0
	for _, r := range s {
		switch unicode.ToLower(r) {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return count
}

// Factorial returns n! for non-negative n. It returns 0 for negative n.
func Factorial(n int) int {
	if n < 0 {
		return 0
	}

	result := 1
	for value := 2; value <= n; value++ {
		result *= value
	}
	return result
}

// Power returns base raised to a non-negative exponent.
func Power(base, exponent int) int {
	if exponent < 0 {
		return 0
	}

	result := 1
	for exponent > 0 {
		if exponent%2 == 1 {
			result *= base
		}
		base *= base
		exponent /= 2
	}
	return result
}

package utility

import "strconv"

func IsLeap(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

func CheckLeap(acc int, password []rune) bool {
	if len(password) == 0 {
		return false
	}

	c := password[0]

	if c < '0' || c > '9' {
		return false
	}

	n := int(c - '0')
	if IsLeap(acc + n) {
		return true
	}

	return CheckLeap(n*10, password[1:])
}

func GetDigit(i int) int {
	return len(strconv.Itoa(i))
}

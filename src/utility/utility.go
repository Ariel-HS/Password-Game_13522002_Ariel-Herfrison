package utility

import (
	"math"
	"regexp"
	"strconv"
)

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

func CheckPrime(x int) bool {
	sqRoot := int(math.Sqrt(float64(x)))
	for i := 2; i <= sqRoot; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func HasNumber(i int, x int) bool {
	// check if x is in i
	copy := i

	for copy > 0 {
		digit := copy % 10
		copy = copy / 10

		if digit == x {
			return true
		}
	}

	return false
}

func GetNumberCount(str []rune) int {
	ctr := 0
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= '0' && c <= '9' {
			ctr++
		}
	}

	return ctr
}

func GetUppercaseCount(str []rune) int {
	ctr := 0
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 'A' && c <= 'Z' {
			ctr++
		}
	}

	return ctr
}

func GetSCharCount(str string) int {
	r := regexp.MustCompile("[!@#$%^&*()\\-_=+\\\\|\\[\\]{};:\\/?.<>'\"]")
	ctr := len(r.FindAllString(str, -1))

	return ctr
}

func GetMonthCount(str []rune) int {
	months := []string{`january`, `february`, `march`, `april`, `may`, `june`, `july`, `august`, `september`, `october`, `november`, `december`}
	ctr := 0
	for _, month := range months {
		match, _ := regexp.MatchString(`(?i)`+month, string(str))

		if match {
			ctr++
		}
	}

	return ctr
}

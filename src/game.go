package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/Ariel-HS/Password-Game_13522002_Ariel-Herfrison/utility"
	_ "github.com/mattn/go-sqlite3"
)

// import "./utility/utility.go"

type Rule struct {
	Emoji string
	Text  string
	Extra template.HTML
}

type Country struct {
	name string
	flag string
}

type Captcha struct {
	answer string
	image  string
}

func main() {
	var pauled bool
	var superPauled bool
	var highScore int
	var combustible bool
	var Rules []Rule
	var country1 Country
	var country2 Country
	var country3 Country
	var captcha Captcha
	bestTime := "-"
	bestTimeInt := -1
	var difficulty string

	//--> Start of Database Setup <--//
	db, err := sql.Open("sqlite3", "password_game.db")

	if err != nil {
		fmt.Println("Database not found")
		return
	}
	defer db.Close()

	rowsCountry, _ := db.Query(`
		SELECT * FROM country
	`)

	var countries []Country
	var name string
	var flag string
	for rowsCountry.Next() {
		rowsCountry.Scan(&name, &flag)

		countries = append(countries, Country{name: name, flag: flag})
	}

	rowsCaptcha, _ := db.Query(`
		SELECT * FROM captcha
	`)

	var captchas []Captcha
	var answer string
	var image string
	for rowsCaptcha.Next() {
		rowsCaptcha.Scan(&answer, &image)

		captchas = append(captchas, Captcha{answer: answer, image: image})
	}
	//--> End of Database Setup <--//

	handler := func(w http.ResponseWriter, r *http.Request) {
		difficulty = r.PostFormValue("difficulty")
		if difficulty == "" {
			tmpl := template.Must(template.ParseFiles("difficulty.html"))

			tmpl.Execute(w, nil)

			return
		}

		pauled = false
		superPauled = false
		combustible = false
		highScore = 1

		country1 = countries[rand.Intn(len(countries))]
		country2 = countries[rand.Intn(len(countries))]
		for country1.name == country2.name {
			country2 = countries[rand.Intn(len(countries))]
		}
		country3 = countries[rand.Intn(len(countries))]
		for country1.name == country3.name || country2.name == country3.name {
			country3 = countries[rand.Intn(len(countries))]
		}
		captcha = captchas[rand.Intn(len(captchas))]

		if difficulty == "easy" {
			Rules = []Rule{
				{Emoji: "❌", Text: "Your password must include the current time"},
				{Emoji: "❌", Text: "The length of your password must be a prime number"},
				{Emoji: "❌", Text: "Your password must include the length of your password"},
				{Emoji: "❌", Text: "At least 30% of your password must be in digits"},
				{Emoji: "❌", Text: "Your password must contain one of the following words: I want IRK | I need IRK | I love IRK"},
				{Emoji: "❌", Text: "A sacrifice must be made. Pick 2 letters that you will no longer be able to use"},
				{Emoji: "❌", Text: "🐔 Paul has hatched ! Please don't forget to feed him. He eats 1 🐛 every 20 second"},
				{Emoji: "❌", Text: "Your password must include a leap year"},
				{Emoji: "❌", Text: "Your password must include this CAPTCHA"},
				{Emoji: "❌", Text: "🥚 This is my chicken Paul. He hasn't hatched yet. Please put him in your password and keep him safe"},
				{Emoji: "❌", Text: "Oh no! Your password is on fire 🔥. Quick, put it out!"},
				{Emoji: "❌", Text: "The Roman numerals in your password should multiply to 35"},
				{Emoji: "❌", Text: "Your password must include one of these countries"},
				{Emoji: "❌", Text: "Your password must include a Roman numeral"},
				{Emoji: "❌", Text: "Your password must include a month of the year"},
				{Emoji: "❌", Text: "The digits in your password must add up to 45"},
				{Emoji: "❌", Text: "Your password must include a special character (! @ # $ % ^ & * ( ) - _ = + \\ | [ ] { } ; : / ? . < > ' \")"},
				{Emoji: "❌", Text: "Your password must include an uppercase letter"},
				{Emoji: "❌", Text: "Your password must include a number"},
				{Emoji: "❌", Text: "Your password must be at least 3 characters"},
			}
		} else if difficulty == "normal" {
			Rules = []Rule{
				{Emoji: "❌", Text: "Your password must include 2 hours after the current time"},
				{Emoji: "❌", Text: "The length of your password must be a prime number that contains the number '3'"},
				{Emoji: "❌", Text: "Your password must include the length of your password + 3"},
				{Emoji: "❌", Text: "At least 40% of your password must be in digits"},
				{Emoji: "❌", Text: "Your password must contain two of the following words: I want IRK | I need IRK | I love IRK"},
				{Emoji: "❌", Text: "A sacrifice must be made. Pick 2 letters that you will no longer be able to use"},
				{Emoji: "❌", Text: "🐔 Paul has hatched ! Please don't forget to feed him. He eats 2 🐛 every 20 second"},
				{Emoji: "❌", Text: "Your password must include a leap year"},
				{Emoji: "❌", Text: "Your password must include this CAPTCHA"},
				{Emoji: "❌", Text: "🥚 This is my chicken Paul. He hasn't hatched yet. Please put him in your password and keep him safe"},
				{Emoji: "❌", Text: "Oh no! Your password is on fire 🔥. Quick, put it out!"},
				{Emoji: "❌", Text: "The Roman numerals in your password should multiply to 35"},
				{Emoji: "❌", Text: "Your password must include 2 of these countries"},
				{Emoji: "❌", Text: "Your password must include 2 Roman numerals"},
				{Emoji: "❌", Text: "Your password must include 2 different months of the year"},
				{Emoji: "❌", Text: "The digits in your password must add up to 35"},
				{Emoji: "❌", Text: "Your password must include 2 special characters (! @ # $ % ^ & * ( ) - _ = + \\ | [ ] { } ; : / ? . < > ' \")"},
				{Emoji: "❌", Text: "Your password must include 3 uppercase letters"},
				{Emoji: "❌", Text: "Your password must include 2 number"},
				{Emoji: "❌", Text: "Your password must be at least 5 characters"},
			}
		} else if difficulty == "hard" {
			Rules = []Rule{
				{Emoji: "❌", Text: "Your password must include 2 hours and 30 minutes after the current time"},
				{Emoji: "❌", Text: "The length of your password must be a prime number that contains the number '3' and '7'"},
				{Emoji: "❌", Text: "Your password must include the length of your password + 13"},
				{Emoji: "❌", Text: "At least 50% of your password must be in digits"},
				{Emoji: "❌", Text: "Your password must contain all of the following words: I want IRK | I need IRK | I love IRK"},
				{Emoji: "❌", Text: "A sacrifice must be made. Pick 2 letters that you will no longer be able to use"},
				{Emoji: "❌", Text: "🐔 Paul has hatched ! Please don't forget to feed him. He eats 3 🐛 every 20 second"},
				{Emoji: "❌", Text: "Your password must include a leap year"},
				{Emoji: "❌", Text: "Your password must include this CAPTCHA"},
				{Emoji: "❌", Text: "🥚 This is my chicken Paul. He hasn't hatched yet. Please put him in your password and keep him safe"},
				{Emoji: "❌", Text: "Oh no! Your password is on fire 🔥. Quick, put it out!"},
				{Emoji: "❌", Text: "The Roman numerals in your password should multiply to 35"},
				{Emoji: "❌", Text: "Your password must include all of these countries"},
				{Emoji: "❌", Text: "Your password must include 4 Roman numeral"},
				{Emoji: "❌", Text: "Your password must include 3 different months of the year"},
				{Emoji: "❌", Text: "The digits in your password must add up to 30"},
				{Emoji: "❌", Text: "Your password must include 3 special characters (! @ # $ % ^ & * ( ) - _ = + \\ | [ ] { } ; : / ? . < > ' \")"},
				{Emoji: "❌", Text: "Your password must include 5 uppercase letters"},
				{Emoji: "❌", Text: "Your password must include 3 numbers"},
				{Emoji: "❌", Text: "Your password must be at least 10 characters"},
			}
		}
		Rules[5].Extra = `<div class="row justify-content-center m-1" id="keyboard">
						<div class="row justify-content-center m-1">
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">A</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">B</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">C</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">D</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">E</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">F</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">G</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">H</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">I</button>
						</div>
						<div class="row justify-content-center m-1">
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()" autocomplete="off">J</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">K</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">L</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">M</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">N</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">O</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">P</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">Q</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">R</button>
						</div>
						<div class="row justify-content-center m-1 w-75">
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">S</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">T</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">U</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">V</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">W</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">X</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">Y</button>
							<button data-bs-toggle="button" data-bs-toggle="button" class="col btn btn-light mx-1" onclick="buttonPress()">Z</button>
						</div>
						<div class="row justify-content-center m-1 w-25">
							<button class="col btn btn-secondary" type="submit" 
							style="border-color:black; min-height: 25px;"
							hx-post="/sacrifice/" 
							hx-vals='js:{password: getPassword(), sacrifice: sacrifice(), extra: getExtra()}'
							hx-swap="multi:#rule-list:innerHTML,#inputBackground:innerHTML,#inputEntry:innerHTML,#Length:innerHTML,#timer-fire:outerHTML,#timer-paul:outerHTML,#game-over:outerHTML" 
							required="required"
							hx-ext="multi-swap"
							>🔥 Sacrifice</button>
						</div>
					</div>`

		tmpl := template.Must(template.ParseFiles("index.html"))

		rules := map[string][]Rule{
			"Rules": Rules[19:],
		}

		tmpl.Execute(w, rules)
	}

	selectHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		tmpl := template.Must(template.ParseFiles("difficulty.html"))

		tmpl.Execute(w, nil)
	}

	checkRules := func(password []rune) {
		Rule20 := func(password []rune) {
			now := ""
			if difficulty == "easy" {
				now = time.Now().Format("15:04")
			} else if difficulty == "normal" {
				now = time.Now().Add(time.Hour * 2).Format("15:04")
			} else if difficulty == "hard" {
				now = time.Now().Add(time.Hour*2 + time.Minute*30).Format("15:04")
			}
			// fmt.Println(now)

			match, _ := regexp.MatchString(now, string(password))

			if match {
				Rules[0].Emoji = "✅"
			} else {
				Rules[0].Emoji = "❌"
			}
		}

		Rule19 := func(password []rune) {
			num := len(password)

			isTrue := true

			if difficulty == "normal" && !utility.HasNumber(num, 3) {
				isTrue = false
			} else if difficulty == "hard" && (!utility.HasNumber(num, 3) || !utility.HasNumber(num, 7)) {
				isTrue = false
			}

			if isTrue && utility.CheckPrime(num) { // if isPrime and contains the needed numbers
				Rules[1].Emoji = "✅"
			} else {
				Rules[1].Emoji = "❌"
			}

			Rule20(password)
		}

		Rule18 := func(password []rune) {
			num := ""
			if difficulty == "easy" {
				num = strconv.Itoa(len(password))
			} else if difficulty == "normal" {
				num = strconv.Itoa(len(password) + 3)
			} else if difficulty == "hard" {
				num = strconv.Itoa(len(password) + 13)
			}
			match, _ := regexp.MatchString(num, string(password))

			if match {
				Rules[2].Emoji = "✅"
			} else {
				Rules[2].Emoji = "❌"
			}

			Rule19(password)
		}

		Rule17 := func(password []rune) {
			num := 0
			if difficulty == "easy" {
				num = len(password) * 3 / 10
			} else if difficulty == "normal" {
				num = len(password) * 4 / 10
			} else if difficulty == "hard" {
				num = len(password) * 5 / 10
			}
			ctr := utility.GetNumberCount(password)

			if ctr >= num {
				Rules[3].Emoji = "✅"
			} else {
				Rules[3].Emoji = "❌"
			}

			Rule18(password)
		}

		Rule16 := func(password []rune) {
			var match bool
			if difficulty == "easy" {
				match, _ = regexp.MatchString("(I want IRK)|(I need IRK)|(I love IRK)", string(password))
			} else if difficulty == "normal" {
				reg := regexp.MustCompile(`(I want IRK).*((I need IRK)|(I love IRK))|(I need IRK).*((I want IRK)|(I love IRK))|(I love IRK).*((I want IRK)|(I need IRK))`)
				match = reg.MatchString(string(password))
			} else if difficulty == "hard" {
				reg := regexp.MustCompile(`(I want IRK).*((I need IRK).*(I love IRK)|(I love IRK).*(I need IRK))|(I need IRK).*((I want IRK).*(I love IRK)|(I love IRK).*(I want IRK))|(I love IRK).*((I want IRK).*(I need IRK)|(I need IRK).*(I want IRK))`)
				match = reg.MatchString(string(password))
			}

			if match {
				Rules[4].Emoji = "✅"
			} else {
				Rules[4].Emoji = "❌"
			}

			Rule17(password)
		}

		// Rule15 skipped

		Rule14 := func(password []rune) {
			match, _ := regexp.MatchString("🐔", string(password))
			if superPauled && !match {
				Rules[6].Emoji = "❌"
				Rules[9].Emoji = "❌"
			}

			Rule16(password)
		}

		Rule13 := func(password []rune) {
			hasLeap := false

			for i := 0; i < len(password); i++ {
				c := password[i]

				if c >= '0' && c <= '9' {
					if utility.CheckLeap(0, password[i:]) {
						hasLeap = true
						break
					}
				}
			}

			if hasLeap {
				Rules[7].Emoji = "✅"
			} else {
				Rules[7].Emoji = "❌"
			}

			Rule14(password)
		}

		Rule12 := func(password []rune) {
			// fmt.Println("Captcha", captcha.answer)
			match, _ := regexp.MatchString(captcha.answer, string(password))
			if match {
				Rules[8].Emoji = "✅"
			} else {
				Rules[8].Emoji = "❌"
			}

			Rule13(password)
		}

		Rule11 := func(password []rune) {
			match, _ := regexp.MatchString("🥚", string(password))
			if !pauled && match {
				Rules[9].Emoji = "✅"
			} else if pauled && !match {
				Rules[9].Emoji = "❌"
			}

			Rule12(password)
		}

		Rule10 := func(password []rune) {
			match, _ := regexp.MatchString("🔥", string(password))
			if !match {
				// fmt.Println("check this")
				Rules[10].Emoji = "✅"
			} else {
				// fmt.Println("onfire")
				Rules[10].Emoji = "❌"
			}

			Rule11(password)
		}

		Rule9 := func(password []rune) {
			// (\s*I*\s+)*
			r := regexp.MustCompile(`^((I?[^IVXLCDM]+)*XXXV([^IVXLCDM]+I?)*)$|^((I?[^IVXLCDM]+)*V[^IVXLCDM]+(I?[^IVXLCDM]+)*VII([^IVXLCDM]+I?)*)$|^((I?[^IVXLCDM]+)*VII[^IVXLCDM]+(I?[^IVXLCDM]+)*V([^IVXLCDM]+I?)*)$`)
			match := r.MatchString(string(password))
			if match {
				Rules[11].Emoji = "✅"
			} else {
				Rules[11].Emoji = "❌"
			}

			Rule10(password)
		}

		Rule8 := func(password []rune) {
			// fmt.Println("Country", country1.name, country2.name, country3.name)
			str := ""
			if difficulty == "easy" {
				str = `(?i)((` + country1.name + `)|(` + country2.name + `)|(` + country3.name + `))`
			} else if difficulty == "normal" {
				str = `(?i)((` + country1.name + `)((` + country2.name + `)|(` + country3.name + `))|(` + country2.name + `)((` + country1.name + `)|(` + country3.name + `))|(` + country3.name + `)((` + country1.name + `)|(` + country2.name + `)))`
			} else if difficulty == "hard" {
				str = `(?i)((` + country1.name + `)((` + country2.name + `)(` + country3.name + `)|(` + country3.name + `)(` + country2.name + `))|(` + country2.name + `)((` + country1.name + `)(` + country3.name + `)|(` + country3.name + `)(` + country1.name + `))|(` + country3.name + `)((` + country1.name + `)(` + country2.name + `)|(` + country2.name + `)(` + country1.name + `)))`
			}
			r := regexp.MustCompile(str)
			match := r.MatchString(string(password))
			if match {
				Rules[12].Emoji = "✅"
			} else {
				Rules[12].Emoji = "❌"
			}

			Rule9(password)
		}

		Rule7 := func(password []rune) {
			var reg *regexp.Regexp
			if difficulty == "easy" {
				reg = regexp.MustCompile(`[IVXLCDM]`)
			} else if difficulty == "normal" {
				reg = regexp.MustCompile(`[IVXLCDM].*[IVXLCDM]`)
			} else if difficulty == "hard" {
				reg = regexp.MustCompile(`[IVXLCDM].*[IVXLCDM].*[IVXLCDM].*[IVXLCDM]`)
			}
			match := reg.MatchString(string(password))
			if match {
				Rules[13].Emoji = "✅"
			} else {
				Rules[13].Emoji = "❌"
			}

			Rule8(password)
		}

		Rule6 := func(password []rune) {
			n := utility.GetMonthCount(password)

			if (difficulty == "easy" && n >= 1) || (difficulty == "normal" && n >= 2) || (difficulty == "hard" && n >= 3) {
				Rules[14].Emoji = "✅"
			} else {
				Rules[14].Emoji = "❌"
			}

			Rule7(password)
		}

		Rule5 := func(password []rune) {
			sum := func() int {
				acc := 0
				for i := 0; i < len(password); i++ {
					if password[i] >= '0' && password[i] <= '9' {
						acc += int(password[i] - '0')
					}
				}

				return acc
			}

			total := sum()

			if (difficulty == "easy" && total == 45) || (difficulty == "normal" && total == 35) || (difficulty == "hard" && total == 30) {
				Rules[15].Emoji = "✅"

			} else {
				Rules[15].Emoji = "❌"
			}

			Rule6(password)
		}

		Rule4 := func(password []rune) {
			match := false
			if difficulty == "easy" {
				match, _ = regexp.MatchString("[!@#$%^&*()\\-_=+\\\\|\\[\\]{};:\\/?.<>'\"]", string(password))
			} else if difficulty == "normal" {
				match, _ = regexp.MatchString("[!@#$%^&*()\\-_=+\\\\|\\[\\]{};:\\/?.<>'\"].*[!@#$%^&*()\\-_=+\\\\|\\[\\]{};:\\/?.<>'\"]", string(password))
			} else if difficulty == "hard" {
				match, _ = regexp.MatchString("[!@#$%^&*()\\-_=+\\\\|\\[\\]{};:\\/?.<>'\"].*[!@#$%^&*()\\-_=+\\\\|\\[\\]{};:\\/?.<>'\"].*[!@#$%^&*()\\-_=+\\\\|\\[\\]{};:\\/?.<>'\"]", string(password))
			}

			if match {
				Rules[16].Emoji = "✅"
			} else {
				Rules[16].Emoji = "❌"
			}

			Rule5(password)
		}

		Rule3 := func(password []rune) {
			match := false
			if difficulty == "easy" {
				match, _ = regexp.MatchString("[A-Z]", string(password))
			} else if difficulty == "normal" {
				match, _ = regexp.MatchString("[A-Z].*[A-Z].*[A-Z]", string(password))
			} else if difficulty == "hard" {
				match, _ = regexp.MatchString("[A-Z].*[A-Z].*[A-Z].*[A-Z].*[A-Z]", string(password))
			}

			if match {
				Rules[17].Emoji = "✅"
			} else {
				Rules[17].Emoji = "❌"
			}

			Rule4(password)
		}

		Rule2 := func(password []rune) {
			match := false
			if difficulty == "easy" {
				match, _ = regexp.MatchString("[0-9]", string(password))
			} else if difficulty == "normal" {
				match, _ = regexp.MatchString("[0-9].*[0-9]", string(password))
			} else if difficulty == "hard" {
				match, _ = regexp.MatchString("[0-9].*[0-9].*[0-9]", string(password))
			}

			if match {
				Rules[18].Emoji = "✅"
			} else {
				Rules[18].Emoji = "❌"
			}

			Rule3(password)
		}

		Rule1 := func(password []rune) {
			if (difficulty == "easy" && len(password) >= 3) || (difficulty == "normal" && len(password) >= 5) || (difficulty == "hard" && len(password) >= 10) {
				Rules[19].Emoji = "✅"
			} else {
				Rules[19].Emoji = "❌"
			}

			Rule2(password)
		}

		Rule1(password)
	}

	apply := func(w http.ResponseWriter, r *http.Request, password []rune) {
		background := ""

		allCorrect := func() bool {
			for i := 1; i <= highScore; i++ {
				if Rules[20-i].Emoji != "✅" {
					return false
				}
			}

			return true
		}

		Rule20 := func(password []rune) {
			if highScore < 20 {
				return
			}

			if allCorrect() {
				fmt.Println("Success")

				timeStr := r.PostFormValue("time")

				timeInt, _ := strconv.Atoi(r.PostFormValue("timeInt"))
				if timeInt < bestTimeInt || bestTimeInt == -1 {
					bestTimeInt = timeInt
					bestTime = timeStr
				}

				str := `<div id="game-over" class="flex-column" style="position: absolute; width: 100%; height: 100%; display: flex;">
							<div class="row flex-grow-1 align-items-center justify-content-center" style="display: flex;font-size: 48;">
								<div style="position: absolute; width: 100%; height: 100%; background-color: grey; opacity: 0.5;"></div>
								<div class="align-self-center align-items-center justify-content-center" style="display: flex; height: 200px; background-color: black; color: gold; position: absolute; width: 100%;">
									<div class="col align-self-center align-items-center justify-content-center">
										<div class="row row position align-self-center align-items-center justify-content-center"
										style="text-align: center; font-size: 320%;">
											GREAT TRIAL CONQUERED
										</div>
										<div class="row position align-self-center align-items-center justify-content-center"
										style="font-size: 120%;">
											Your Time: ` + timeStr + `
										</div>
										<div class="row position align-self-center align-items-center justify-content-center"
										style="font-size: 120%;">
											Best Time: ` + bestTime + `
										</div>
									</div>
								</div>

								<script>
									gameOver()
								</script>
							</div>
						</div>`
				tmpl, _ := template.New("t").Parse(str)
				tmpl.Execute(w, str)
				superPauled = false
				combustible = false
			}
		}

		Rule19 := func(password []rune) {
			if highScore < 19 {
				return
			}

			if highScore < 20 && allCorrect() {
				highScore = 20
			}

			Rule20(password)
		}

		Rule18 := func(password []rune) {
			if highScore < 18 {
				return
			}

			if highScore < 19 && allCorrect() {
				highScore = 19
			}

			Rule19(password)
		}

		Rule17 := func(password []rune) {
			if highScore < 17 {
				return
			}

			if highScore < 18 && allCorrect() {
				highScore = 18
			}

			Rule18(password)
		}

		Rule16 := func(password []rune) {
			if highScore < 16 {
				return
			}

			if highScore < 17 && allCorrect() {
				highScore = 17
			}

			Rule17(password)
		}

		Rule15 := func(password []rune) {
			if highScore < 15 {
				return
			}

			if highScore < 16 && allCorrect() {
				highScore = 16
			}

			Rule16(password)
		}

		Rule14 := func(password []rune) {
			if highScore < 14 {
				return
			}

			if Rules[6].Emoji == "❌" {
				fmt.Println("failure")

				time := r.PostFormValue("time")

				str := `<div id="game-over" class="flex-column" style="position: absolute; width: 100%; height: 100%; display: flex;">
							<div class="row flex-grow-1 align-items-center justify-content-center" style="display: flex;font-size: 48;">
								<div style="position: absolute; width: 100%; height: 100%; background-color: grey; opacity: 0.5;"></div>
								<div class="align-self-center align-items-center justify-content-center" style="display: flex; height: 200px; background-color: black; color: red; position: absolute; width: 100%;">
									<div class="col align-self-center align-items-center justify-content-center">
										<div class="row row position align-self-center align-items-center justify-content-center"
										style="text-align: center; font-size: 320%;">
											HATCH PAUL IS KILL
										</div>
										<div class="row position align-self-center align-items-center justify-content-center"
										style="font-size: 120%;">
											Your Time: ` + time + `
										</div>
										<div class="row position align-self-center align-items-center justify-content-center"
										style="font-size: 120%;">
											Best Time: ` + bestTime + `
										</div>
									</div>
								</div>

								<script>
									gameOver()
								</script>
							</div>
						</div>`
				tmpl, _ := template.New("t").Parse(str)
				tmpl.Execute(w, str)
				highScore = 1
				superPauled = false
				combustible = false

				return
			}

			if highScore < 15 {
				highScore = 15

				str := `<div class="row justify-content-center m-1">
								<div class="col-3"></div>
								<div class="col-3 align-self-center align-items-center justify-content-center" 
								style="display: flex;" id="captcha">
									<img src="` + captcha.image + `" width="96" height="64">
								</div>
								<div class="col-3 align-self-center align-items-center justify-content-start"
								style="display: flex;">
									<button class="btn"
									type="submit" hx-post="/reCaptcha/"
									hx-swap="multi:#rule-list:innerHTML,#inputBackground:innerHTML,#inputEntry:innerHTML,#Length:innerHTML,#timer-fire:outerHTML,#timer-paul:outerHTML,#game-over:outerHTML" 
									hx-vals='js:{password: getPassword(), extra: getExtra()}'
									hx-ext="multi-swap"
									>
										<font size="5">🔄</font>
									</button>
								</div>
							</div>`

				Rules[8].Extra = template.HTML(str)

				str = `<div id = "timer-fire">
								<div
								hx-post="/timerFire/"
								hx-trigger="every 1.5s"
								hx-swap="multi:#rule-list:innerHTML,#inputBackground:innerHTML,#inputEntry:innerHTML,#Length:innerHTML,#timer-fire:outerHTML,#timer-paul:outerHTML,#game-over:outerHTML" 
								hx-vals='js:{password: getPassword(), extra: getExtra()}'
								hx-ext="multi-swap"
								></div>
							</div>`
				tmpl, _ := template.New("t").Parse(str)
				tmpl.Execute(w, str)
			}

			Rule15(password)
		}

		Rule13 := func(password []rune) {
			if highScore < 13 {
				return
			}

			if highScore < 14 && allCorrect() {
				r := regexp.MustCompile(`^(.*)🥚(.*)$`)
				newPassword := r.ReplaceAllString(string(password), "${1}🐔$2")

				str := `<div id="inputEntry" class="form-control" style="position: relative; max-width: 650px; background: transparent; word-wrap: break-word; display: inline-block;" contenteditable="true"
                        >` +
					newPassword + `</div>`

				tmpl, _ := template.New("t").Parse(str)
				tmpl.Execute(w, str)

				password = []rune(newPassword)
				background = newPassword

				pauled = false
				superPauled = true

				str = `<div id = "timer-paul">
							<div
							hx-post="/timerPaul/"
							hx-trigger="every 20s"
							hx-swap="multi:#rule-list:innerHTML,#inputBackground:innerHTML,#inputEntry:innerHTML,#Length:innerHTML,#timer-fire:outerHTML,#timer-paul:outerHTML,#game-over:outerHTML" 
							hx-vals='js:{password: getPassword(), extra: getExtra()}'
							hx-ext="multi-swap"
							></div>
						</div>`
				tmpl, _ = template.New("t").Parse(str)
				tmpl.Execute(w, str)

				Rules[6].Emoji = "✅"

				highScore = 14
			}

			Rule14(password)
		}

		Rule12 := func(password []rune) {
			if highScore < 12 {
				return
			}

			if highScore < 13 && allCorrect() {
				highScore = 13
			}

			Rule13(password)
		}

		Rule11 := func(password []rune) {
			if highScore < 11 {
				return
			}

			if pauled && Rules[9].Emoji == "❌" {
				fmt.Println("failure")

				time := r.PostFormValue("time")

				str := `<div id="game-over" class="flex-column" style="position: absolute; width: 100%; height: 100%; display: flex;">
							<div class="row flex-grow-1 align-items-center justify-content-center" style="display: flex;font-size: 48;">
								<div style="position: absolute; width: 100%; height: 100%; background-color: grey; opacity: 0.5;"></div>
								<div class="align-self-center align-items-center justify-content-center" style="display: flex; height: 200px; background-color: black; color: red; position: absolute; width: 100%;">
									<div class="col align-self-center align-items-center justify-content-center">
										<div class="row row position align-self-center align-items-center justify-content-center"
										style="text-align: center; font-size: 320%;">
											PAUL IS KILL
										</div>
										<div class="row position align-self-center align-items-center justify-content-center"
										style="font-size: 120%;">
											Your Time: ` + time + `
										</div>
										<div class="row position align-self-center align-items-center justify-content-center"
										style="font-size: 120%;">
											Best Time: ` + bestTime + `
										</div>
									</div>
								</div>

								<script>
									gameOver()
								</script>
							</div>
						</div>`
				tmpl, _ := template.New("t").Parse(str)
				tmpl.Execute(w, str)
				highScore = 1
				pauled = false
				combustible = false

				return
			}

			if highScore < 12 && allCorrect() {
				highScore = 12

				// insert captcha
				str := `<div class="row justify-content-center m-1">
								<div class="col-3"></div>
								<div class="col-3 align-self-center align-items-center justify-content-center" 
								style="display: flex;" id="captcha">
									<img src="` + captcha.image + `" width="96" height="64">
								</div>
								<div class="col-3 align-self-center align-items-center justify-content-start"
								style="display: flex;">
									<button class="btn"
									type="submit" hx-post="/reCaptcha/"
									hx-swap="multi:#rule-list:innerHTML,#inputBackground:innerHTML,#inputEntry:innerHTML,#Length:innerHTML,#timer-fire:outerHTML,#timer-paul:outerHTML,#game-over:outerHTML" 
									hx-vals='js:{password: getPassword()}'
									hx-ext="multi-swap"
									>
										<font size="5">🔄</font>
									</button>
								</div>
							</div>`

				Rules[8].Extra = template.HTML(str)
				pauled = true
			}

			Rule12(password)
		}

		Rule10 := func(password []rune) {
			if highScore < 10 {
				return
			}

			// fmt.Println("check this")
			if highScore < 11 && allCorrect() {
				highScore = 11
			}

			Rule11(password)
		}

		Rule9 := func(password []rune) {
			if highScore < 9 {
				return
			}

			if Rules[11].Emoji == "✅" {
				if highScore < 10 && allCorrect() {
					newPassword := string(password[:len(password)-1]) + "🔥"
					str := `<div id="inputEntry" class="form-control" style="position: relative; max-width: 650px; background: transparent; word-wrap: break-word; display: inline-block;" contenteditable="true"
                        >` +
						newPassword + `</div>`

					tmpl, _ := template.New("t").Parse(str)
					tmpl.Execute(w, str)

					highScore = 10

					fire := []rune("🔥")[0]
					password[len(password)-1] = fire
					background = string(password)

					// fmt.Println("Called")
					str = `<div id = "timer-fire">
								<div
								hx-post="/timerFire/"
								hx-trigger="every 1.5s"
								hx-swap="multi:#rule-list:innerHTML,#inputBackground:innerHTML,#inputEntry:innerHTML,#Length:innerHTML,#timer-fire:outerHTML,#timer-paul:outerHTML,#game-over:outerHTML" 
								hx-vals='js:{password: getPassword()}'
								hx-ext="multi-swap"
								></div>
							</div>`
					tmpl, _ = template.New("t").Parse(str)
					tmpl.Execute(w, str)

					combustible = true
				}
			} else {
				m := regexp.MustCompile(`[IVXLCDM]`)
				background = m.ReplaceAllString(string(password), `<span style="background-color: firebrick;">${0}</span>`)
			}

			Rule10(password)
		}

		Rule8 := func(password []rune) {
			if highScore < 8 {
				return
			}

			if highScore < 9 && allCorrect() {
				highScore = 9
			}

			Rule9(password)
		}

		Rule7 := func(password []rune) {
			if highScore < 7 {
				return
			}

			if highScore < 8 && allCorrect() {
				highScore = 8

				str := `<div class="row justify-content-center m1-3">
							<div class="col-3"><img src="` + country1.flag + `" width="96" height="64"></div>
							<div class="col-3"><img src="` + country2.flag + `" width="96" height="64"></div>
							<div class="col-3"><img src="` + country3.flag + `" width="96" height="64"></div>
						</div>`

				Rules[12].Extra = template.HTML(str)
			}

			Rule8(password)
		}

		Rule6 := func(password []rune) {
			if highScore < 6 {
				return
			}

			if highScore < 7 && allCorrect() {
				highScore = 7
			}

			Rule7(password)
		}

		Rule5 := func(password []rune) {
			if highScore < 5 {
				return
			}

			if Rules[15].Emoji == "✅" {
				if highScore < 6 && allCorrect() {
					highScore = 6
				}
			} else {
				m := regexp.MustCompile("([0-9])+")
				background = m.ReplaceAllString(string(password), `<span style="background-color: firebrick;">${0}</span>`)
			}

			Rule6(password)
		}

		Rule4 := func(password []rune) {
			if highScore < 4 {
				return
			}

			if highScore < 5 && allCorrect() {
				highScore = 5
			}

			Rule5(password)
		}

		Rule3 := func(password []rune) {
			if highScore < 3 {
				return
			}

			if highScore < 4 && allCorrect() {
				highScore = 4
			}

			Rule4(password)
		}

		Rule2 := func(password []rune) {
			if highScore < 2 {
				return
			}

			if highScore < 3 && allCorrect() {
				highScore = 3
			}

			Rule3(password)
		}

		Rule1 := func(password []rune) {
			if highScore < 2 && allCorrect() {
				highScore = 2
			}

			Rule2(password)
		}

		background = string(password)
		passLength := len(password)
		log.Print(passLength)

		Rule1(password)

		tmpl := template.Must(template.ParseFiles("index.html"))

		fmt.Println(highScore)
		// fmt.Println(Rules[(20 - highScore):])
		rules := map[string][]Rule{
			"Rules": Rules[(20 - highScore):],
		}

		tmpl.ExecuteTemplate(w, "rule-list-element", rules)

		str := `<label id="Length" for="Length" class="form-label">` + strconv.Itoa(passLength) + `</label>`
		tmpl, _ = template.New("t").Parse(str)
		tmpl.Execute(w, passLength)

		fmt.Println("background:", background)
		fmt.Println("password", string(password))

		str = `<div id="inputBackground" class="form-control" style="position: absolute; border: none; max-width: 649px; word-wrap: break-word; display: inline-block;">` +
			background + `</div>`
		tmpl, _ = template.New("t").Parse(str)
		tmpl.Execute(w, background)
	}

	checkAndApply := func(w http.ResponseWriter, r *http.Request, password []rune) {
		checkRules(password)
		apply(w, r, password)
	}

	reCaptcha := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ReCaptcha")
		newCaptcha := captchas[rand.Intn(len(captchas))]
		fmt.Println("before", captcha.answer, newCaptcha.answer)
		for captcha.answer == newCaptcha.answer {
			newCaptcha = captchas[rand.Intn(len(captchas))]
		}
		fmt.Println("after", captcha.answer, newCaptcha.answer)
		captcha = newCaptcha

		str := `<div class="row justify-content-center m-1">
				<div class="col-3"></div>
				<div class="col-3 align-self-center align-items-center justify-content-center" 
				style="display: flex;" id="captcha">
					<img src="` + captcha.image + `" width="96" height="64">
				</div>
				<div class="col-3 align-self-center align-items-center justify-content-start"
				style="display: flex;">
					<button class="btn"
					type="submit" hx-post="/reCaptcha/"
					hx-swap="multi:#rule-list:innerHTML,#inputBackground:innerHTML,#inputEntry:innerHTML,#Length:innerHTML,#timer-fire:outerHTML,#timer-paul:outerHTML,#game-over:outerHTML"  
                    hx-vals='js:{password: getPassword(), extra: getExtra()}'
                    hx-ext="multi-swap"
					>
						<font size="5">🔄</font>
					</button>
				</div>
			</div>`

		Rules[8].Extra = template.HTML(str)

		if highScore > 14 {
			newExtra := r.PostFormValue("extra")
			Rules[5].Extra = template.HTML(newExtra)
		}

		password := []rune(r.PostFormValue("password"))
		checkAndApply(w, r, password)
	}

	cheat := func(w http.ResponseWriter, r *http.Request, password []rune) {
		reg := regexp.MustCompile("cheat")

		password = []rune(reg.ReplaceAllString(string(password), ""))
		// fmt.Println("PASSWORD: 2", string(password))

		checkRules(password)

		now := ""
		length := ""
		retry := false
		Cheat5 := func(password []rune) []rune {
			sum := func() int {
				acc := 0
				for i := 0; i < len(password); i++ {
					if password[i] >= '0' && password[i] <= '9' {
						acc += int(password[i] - '0')
					}
				}

				return acc
			}

			total := sum()

			difference := 0
			if difficulty == "easy" {
				difference = 45 - total
			} else if difficulty == "normal" {
				difference = 35 - total
			} else if difficulty == "hard" {
				difference = 30 - total
			}

			if difference != 0 && highScore >= 5 {
				if difference > 0 {
					temp := regexp.MustCompile("(.)$")
					nums := ""

					for i := 9; i > 0 && difference > 0; i-- {
						n := difference / i
						difference -= i * n
						for j := 0; j < n; j++ {
							nums += strconv.Itoa(i)
						}
						// print("Hey", i, n, nums)
					}

					password = []rune(temp.ReplaceAllString(string(password), "${0}"+nums))
				} else {
					difference = -difference
					tempPassword := string(password)
					regCaptcha := regexp.MustCompile(captcha.answer)
					locCaptcha := regCaptcha.FindStringIndex(tempPassword)
					regTime := regexp.MustCompile(now)
					locTime := regTime.FindStringIndex(tempPassword)
					regLength := regexp.MustCompile(length)
					locLength := regLength.FindStringIndex(tempPassword)
					fmt.Println("LOC", locCaptcha, locLength, locTime)
					fmt.Println(tempPassword[locCaptcha[0]:locCaptcha[1]],
						tempPassword[locTime[0]:locTime[1]],
						tempPassword[locLength[0]:locLength[1]])

					for i := 0; i < len(tempPassword) && difference > 0; i++ { // avoid captcha, time, and length
						if locCaptcha != nil && i >= locCaptcha[0] && i < locCaptcha[1] { // if captcha exists and i is in captcha, skip
							continue
						}
						if now != "" && locTime != nil && i >= locTime[0] && i < locTime[1] { // if time exists and i is in time, skip
							continue
						}
						if length != "" && locLength != nil && i >= locLength[0] && i < locLength[1] { // if length exists and i is in length, skip
							continue
						}

						c := tempPassword[i]
						// if char is number
						if c >= '0' && c <= '9' {
							x := int(c - '0')
							if x <= difference {
								difference -= x
								x = 0
							} else {
								x -= difference
								difference = 0
							}
							// fmt.Println(tempPassword[:i], "|", x, "|", tempPassword[i+1:])
							tempPassword = tempPassword[:i] + strconv.Itoa(x) + tempPassword[i+1:]
						}
					}

					if difference > 0 {
						// reCaptcha
						fmt.Println("ReCaptcha try")
						newCaptcha := captchas[rand.Intn(len(captchas))]
						// fmt.Println("before", captcha.answer, newCaptcha.answer)
						for captcha.answer == newCaptcha.answer {
							newCaptcha = captchas[rand.Intn(len(captchas))]
						}
						// fmt.Println("after", captcha.answer, newCaptcha.answer)
						captcha = newCaptcha
						fmt.Println("after", captcha.answer)

						str := `<div class="row justify-content-center m-1">
								<div class="col-3"></div>
								<div class="col-3 align-self-center align-items-center justify-content-center" 
								style="display: flex;" id="captcha">
									<img src="` + captcha.image + `" width="96" height="64">
								</div>
								<div class="col-3 align-self-center align-items-center justify-content-start"
								style="display: flex;">
									<button class="btn"
									type="submit" hx-post="/reCaptcha/"
									hx-swap="multi:#rule-list:innerHTML,#inputBackground:innerHTML,#inputEntry:innerHTML,#Length:innerHTML,#timer-fire:outerHTML,#timer-paul:outerHTML,#game-over:outerHTML"  
									hx-vals='js:{password: getPassword(), extra: getExtra()}'
									hx-ext="multi-swap"
									>
										<font size="5">🔄</font>
									</button>
								</div>
							</div>`

						Rules[8].Extra = template.HTML(str)

						retry = true

						tempPassword = tempPassword[:locCaptcha[0]] + tempPassword[locCaptcha[1]:]
					}

					password = []rune(tempPassword)
				}
			}

			return password
		}

		Cheat18_19 := func(password []rune) []rune {
			if highScore < 18 {
				return Cheat5(password)
			}

			oldLength := len(password)

			// check rule 18
			checkStr := ""
			if difficulty == "easy" {
				checkStr = strconv.Itoa(oldLength)
			} else if difficulty == "normal" {
				checkStr = strconv.Itoa(oldLength + 3)
			} else if difficulty == "hard" {
				checkStr = strconv.Itoa(oldLength + 13)
			}
			match, _ := regexp.MatchString(checkStr, string(password))

			targetAppend := oldLength
			if !match {
				nDigit := utility.GetDigit(oldLength) // get digit of current password length
				targetAppend += nDigit                // new length = old length + digit of old length
				if difficulty == "normal" {
					targetAppend += 3
				} else if difficulty == "hard" {
					targetAppend += 13
				}

				if utility.GetDigit(targetAppend) != nDigit { // if digit will change (e.g 99 -> 101)
					targetAppend += 1 // modify the string
				}

				length = strconv.Itoa(targetAppend)
			}

			tempPassword := []rune(string(password) + length)
			// fmt.Println("PASSWORD: 3", string(password))
			newLength := len(tempPassword)

			if highScore < 19 { // if rule 19 not apply, skip
				return Cheat5(tempPassword)
			}

			// check rule 19
			// if new length is not prime and not contain necessary number
			isPrime := utility.CheckPrime(newLength)
			if !isPrime || ((difficulty == "normal" || difficulty == "hard") && !utility.HasNumber(newLength, 3)) || (difficulty == "hard" && !utility.HasNumber(newLength, 7)) {
				// search the next length of password that is prime and contains the necessary number
				targetLength := targetAppend + 1
				for !utility.CheckPrime(targetLength) || targetLength < targetAppend+utility.GetDigit(targetAppend) ||
					((difficulty == "normal" || difficulty == "hard") && !utility.HasNumber(targetLength, 3)) ||
					(difficulty == "hard" && !utility.HasNumber(targetLength, 7)) {
					targetLength++
				}

				appendStr := ""                                                   // string to be appended to password
				diff := targetLength - oldLength - utility.GetDigit(targetLength) // number of extra characters needed
				for j := 0; j < diff; j++ {
					appendStr += "0"
				}

				if difficulty == "easy" {
					length = strconv.Itoa(targetLength)
					appendStr += length
				} else if difficulty == "normal" {
					appendStr += strconv.Itoa(targetLength + 3)
					length = strconv.Itoa(targetLength + 3)
				} else if difficulty == "hard" {
					appendStr += strconv.Itoa(targetLength + 13)
					length = strconv.Itoa(targetLength + 13)
				}

				password = []rune(string(password) + appendStr)
				// fmt.Println("PASSWORD: 4", string(password))
			}

			// fmt.Println("PASSWORD: 5", string(password))
			return Cheat5(password)
		}

		Cheat17 := func(password []rune) []rune {
			n := 0
			if difficulty == "easy" {
				n = len(password) * 3 / 10
			} else if difficulty == "normal" {
				n = len(password) * 4 / 10
			} else if difficulty == "hard" {
				n = len(password) * 5 / 10
			}
			ctr := utility.GetNumberCount(password)

			if ctr < n && highScore >= 17 {
				// fmt.Println("Insufficient, ", ctr, n, string(password))
				num := 0
				if difficulty == "easy" {
					num = len(password) * 3 / 7
				} else if difficulty == "normal" {
					num = len(password) * 2 / 3
				} else if difficulty == "hard" {
					num = len(password)
				}

				difference := num - ctr
				nums := ""
				for i := 0; i < difference; i++ {
					nums += "0"
				}

				password = []rune(string(password) + nums)
			}

			return Cheat18_19(password)
		}

		Cheat20 := func(password []rune) []rune {
			if Rules[0].Emoji != "✅" && highScore >= 20 {
				if difficulty == "easy" {
					now = time.Now().Format("15:04")
				} else if difficulty == "normal" {
					now = time.Now().Add(time.Hour * 2).Format("15:04")
				} else if difficulty == "hard" {
					now = time.Now().Add(time.Hour*2 + time.Minute*30).Format("15:04")
				}
				password = []rune(string(password) + now)
			}

			return Cheat17(password)
		}

		Cheat16 := func(password []rune) []rune {
			if Rules[4].Emoji != "✅" && highScore >= 16 {
				if difficulty == "easy" {
					password = []rune(string(password) + "I want IRK")
				} else if difficulty == "normal" {
					match1, _ := regexp.MatchString(string(password), "I want IRK")
					match2, _ := regexp.MatchString(string(password), "I need IRK")
					match3, _ := regexp.MatchString(string(password), "I love IRK")

					if !match1 && !match2 && match3 {
						password = []rune(string(password) + "I want IRKI need IRK")
					} else {
						if !match1 {
							password = []rune(string(password) + "I want IRK")
						} else if !match2 {
							password = []rune(string(password) + "I need IRK")
						}
					}
				} else if difficulty == "hard" {
					str := string(password)
					match1, _ := regexp.MatchString(string(password), "I want IRK")
					match2, _ := regexp.MatchString(string(password), "I need IRK")
					match3, _ := regexp.MatchString(string(password), "I love IRK")

					if !match1 {
						str += "I want IRK"
					}
					if !match2 {
						str += "I need IRK"
					}
					if !match3 {
						str += "I love IRK"
					}
					password = []rune(str)
				}
			}

			return Cheat20(password)
		}

		Cheat14 := func(password []rune) []rune {
			temp := regexp.MustCompile("🐔")
			match, _ := regexp.MatchString(`🐛`, string(password))

			if !match {
				if difficulty == "easy" {
					password = []rune(temp.ReplaceAllString(string(password), "🐔🐛"))
				} else if difficulty == "normal" {
					password = []rune(temp.ReplaceAllString(string(password), "🐔🐛🐛"))
				} else if difficulty == "hard" {
					password = []rune(temp.ReplaceAllString(string(password), "🐔🐛🐛🐛"))
				}
			}

			return Cheat16(password)
		}

		Cheat12 := func(password []rune) []rune {
			if Rules[8].Emoji != "✅" && highScore >= 12 {
				temp := regexp.MustCompile("(.)$")

				password = []rune(temp.ReplaceAllString(string(password), "${0}"+captcha.answer))
			}

			return Cheat14(password)
		}

		Cheat11 := func(password []rune) []rune {
			if Rules[9].Emoji != "✅" && highScore >= 11 {
				temp := regexp.MustCompile("(.)$")

				password = []rune(temp.ReplaceAllString(string(password), "${0}🥚"))
			}

			return Cheat12(password)
		}

		Cheat10 := func(password []rune) []rune {
			if Rules[10].Emoji != "✅" && highScore >= 10 {
				temp := regexp.MustCompile("🔥")

				password = []rune(temp.ReplaceAllString(string(password), ""))
			}

			return Cheat11(password)
		}

		Cheat9 := func(password []rune) []rune {

			cheat9Helper := func(password string) []rune {
				temp := regexp.MustCompile("V([^IV][^V]*)VI+") // check for V VII possibilities, get last occurence
				loc := temp.FindStringIndex(password)
				// fmt.Println(loc)

				if loc != nil {
					passLeft := password[:loc[0]]
					passRight := password[loc[1]:]
					passMid := password[loc[0]:loc[1]]

					// fmt.Println(passLeft, passMid, passRight)

					// if contain other roman numerals (fail), clear it
					clearLeft := regexp.MustCompile("[VXLCDM]|I{2,}|I$")
					clearRight := regexp.MustCompile("[VXLCDM]|I{2,}|^I")
					passLeft = clearLeft.ReplaceAllString(passLeft, "")
					passRight = clearRight.ReplaceAllString(passRight, "")

					clearMid := regexp.MustCompile("I+V")
					clearMid2 := regexp.MustCompile("[XLCDM]|[^V]I{2,}")
					passMid = temp.ReplaceAllString(passMid, "V${1}VII")
					passMid = clearMid2.ReplaceAllString(clearMid.ReplaceAllString(passMid, "V"), "")

					// fmt.Println(passLeft, passMid, passRight)
					password = passLeft + passMid + passRight
					// fmt.Println(password)

					return []rune(password)
				}

				temp = regexp.MustCompile("VI*([^V]+)VI*") // else check for V(II) V(II) possibilities
				loc = temp.FindStringIndex(password)
				// fmt.Println(loc)

				if loc != nil {
					passLeft := password[:loc[0]]
					passRight := password[loc[1]:]
					passMid := password[loc[0]:loc[1]]

					// fmt.Println(passLeft, passMid, passRight)

					// if contain other roman numerals (fail), clear it
					clearLeft := regexp.MustCompile("[VXLCDM]|I{2,}|I$")
					clearRight := regexp.MustCompile("[VXLCDM]|I{2,}|^I")
					passLeft = clearLeft.ReplaceAllString(passLeft, "")
					passRight = clearRight.ReplaceAllString(passRight, "")

					clearMid := regexp.MustCompile("I+V")
					clearMid2 := regexp.MustCompile("[XLCDM]|[^V]I{2,}")
					passMid = temp.ReplaceAllString(passMid, "VII${1}V")
					passMid = clearMid2.ReplaceAllString(clearMid.ReplaceAllString(passMid, "V"), "")

					// fmt.Println(passLeft, passMid, passRight)
					password = passLeft + passMid + passRight
					// fmt.Println(password)

					return []rune(password)
				}

				temp = regexp.MustCompile("X+V*|X*V+") // else modify to XXXV
				loc = temp.FindStringIndex(password)
				// fmt.Println(loc)

				if loc != nil {
					passLeft := password[:loc[0]]
					passRight := password[loc[1]:]
					passMid := password[loc[0]:loc[1]]

					// fmt.Println(passLeft, passMid, passRight)

					// if contain other roman numerals (fail), clear it
					clearLeft := regexp.MustCompile("[VXLCDM]|I{2,}|I$")
					clearRight := regexp.MustCompile("[VXLCDM]|I{2,}|^I")
					passLeft = clearLeft.ReplaceAllString(passLeft, "")
					passRight = clearRight.ReplaceAllString(passRight, "")

					passMid = temp.ReplaceAllString(passMid, "XXXV")

					// fmt.Println(passLeft, passMid, passRight)
					password = passLeft + passMid + passRight
					// fmt.Println(password)

					return []rune(password)
				}

				// else add XXXV
				clear := regexp.MustCompile("[VXLCDM]|I{2,}")
				password = clear.ReplaceAllString(password, "")
				clear = regexp.MustCompile("I$") // if pass ends with I, add a space
				password = clear.ReplaceAllString(password, "I ")

				return []rune(password + "XXXV")
			}

			tempPassword := string(password)

			r := regexp.MustCompile(`^((I?[^IVXLCDM]+)*XXXV([^IVXLCDM]+I?)*)$|^((I?[^IVXLCDM]+)*V[^IVXLCDM]+(I?[^IVXLCDM]+)*VII([^IVXLCDM]+I?)*)$|^((I?[^IVXLCDM]+)*VII[^IVXLCDM]+(I?[^IVXLCDM]+)*V([^IVXLCDM]+I?)*)$`)
			match := r.MatchString(tempPassword)
			if !match && highScore >= 9 {
				password = cheat9Helper(tempPassword)
			}

			return Cheat10(password)
		}

		Cheat8 := func(password []rune) []rune {
			if Rules[12].Emoji != "✅" && highScore >= 8 {
				if difficulty == "easy" {
					password = []rune(string(password) + country1.name)
				} else if difficulty == "normal" {
					match1, _ := regexp.MatchString(string(password), country1.name)
					match2, _ := regexp.MatchString(string(password), country2.name)
					match3, _ := regexp.MatchString(string(password), country3.name)

					if !match1 && !match2 && match3 {
						password = []rune(string(password) + country1.name + country2.name)
					} else {
						if !match1 {
							password = []rune(string(password) + country1.name)
						} else if !match2 {
							password = []rune(string(password) + country2.name)
						}
					}
				} else if difficulty == "hard" {
					str := string(password)
					match1, _ := regexp.MatchString(string(password), country1.name)
					match2, _ := regexp.MatchString(string(password), country2.name)
					match3, _ := regexp.MatchString(string(password), country3.name)

					if !match1 {
						str += country1.name
					}
					if !match2 {
						str += country2.name
					}
					if !match3 {
						str += country3.name
					}
					password = []rune(str)
				}
			}

			return Cheat9(password)
		}

		Cheat7 := func(password []rune) []rune {
			if Rules[13].Emoji != "✅" && highScore >= 7 {
				password = []rune(string(password) + "XXXV")
			}

			return Cheat8(password)
		}

		Cheat6 := func(password []rune) []rune {
			if highScore >= 6 {
				ctr := 0
				if difficulty == "easy" {
					ctr = 1 - utility.GetMonthCount(password)
				} else if difficulty == "normal" {
					ctr = 2 - utility.GetMonthCount(password)
				} else if difficulty == "hard" {
					ctr = 3 - utility.GetMonthCount(password)
				}

				str := string(password)
				months := []string{`january`, `february`, `march`, `april`, `may`, `june`, `july`, `august`, `september`, `october`, `november`, `december`}
				for i := 0; i < 12 && ctr > 0; i++ {
					match, _ := regexp.MatchString(`(?i)`+months[i], str)

					if !match {
						ctr--
						str += months[i]
					}
				}
				password = []rune(str)
			}

			return Cheat7(password)
		}

		Cheat4 := func(password []rune) []rune {
			if Rules[16].Emoji != "✅" && highScore >= 4 {
				str := "@"
				n := 0
				if difficulty == "normal" {
					n = 1 - utility.GetSCharCount(string(password))
				} else if difficulty == "hard" {
					n = 2 - utility.GetSCharCount(string(password))
				}

				for i := 0; i < n; i++ {
					str += "@"
				}

				password = []rune(string(password) + str)
			}

			return Cheat6(password)
		}

		Cheat3 := func(password []rune) []rune {
			if Rules[17].Emoji != "✅" && highScore >= 3 {
				str := "A"
				n := 0
				if difficulty == "normal" {
					n = 2 - utility.GetUppercaseCount(password)
				} else if difficulty == "hard" {
					n = 4 - utility.GetUppercaseCount(password)
				}

				for i := 0; i < n; i++ {
					str += "A"
				}

				password = []rune(string(password) + str)
			}

			return Cheat4(password)
		}

		Cheat2_13 := func(password []rune) []rune {

			if (Rules[18].Emoji != "✅" && highScore >= 2) || (Rules[7].Emoji != "✅" && highScore > 12) {
				str := "0"
				n := 0
				if difficulty == "normal" {
					n = 1 - utility.GetNumberCount(password)
				} else if difficulty == "hard" {
					n = 2 - utility.GetNumberCount(password)
				}

				for i := 0; i < n; i++ {
					str += "0"
				}

				password = []rune(string(password) + str)
			}

			return Cheat3(password)
		}

		Cheat1 := func(password []rune) []rune {
			if Rules[19].Emoji != "✅" {
				temp := regexp.MustCompile("(.?)$")
				a := ""

				n := 0
				if difficulty == "easy" {
					n = 3 - len(password)
				} else if difficulty == "normal" {
					n = 5 - len(password)
				} else if difficulty == "hard" {
					n = 10 - len(password)
				}

				for i := 0; i < n; i++ {
					a = temp.ReplaceAllString(a, "${0}a")
				}

				password = []rune(temp.ReplaceAllString(string(password), "${0}"+a))

				fmt.Println("Cheat 1:", string(password))
			}

			return Cheat2_13(password)
			// return password
		}

		// fmt.Println("Before cheat:", string(password))
		password = Cheat1(password)
		for retry {
			fmt.Println("Retry")
			retry = false
			checkRules(password)
			password = Cheat1(password)
		}

		checkAndApply(w, r, password)
		newPassword := string(password)
		fmt.Println("Cheat:", newPassword, highScore)
		str := `<div id="inputEntry" class="form-control" style="position: relative; max-width: 650px; background: transparent; word-wrap: break-word; display: inline-block;" contenteditable="true"
                        >` +
			newPassword + `</div>`

		tmpl, _ := template.New("t").Parse(str)
		tmpl.Execute(w, str)
	}

	checkHelper := func(w http.ResponseWriter, r *http.Request) {
		log.Print("Request received")
		passwordStr := r.PostFormValue("password")
		// fmt.Println("PASSWORD: ", passwordStr)
		password := []rune(passwordStr)

		match, _ := regexp.MatchString("cheat", passwordStr)
		if match {
			fmt.Println("Cheating")
			cheat(w, r, password)

			return
		}

		log.Print(string(password))

		checkAndApply(w, r, password)
	}

	timerFire := func(w http.ResponseWriter, r *http.Request) {
		if !combustible {
			return
		}

		password := []rune(r.PostFormValue("password"))
		match, _ := regexp.MatchString("🔥", string(password))
		if !match {
			// 1/30 chance to occur again
			// if rand.Intn(30) == 1 {
			// 	newPassword := string(password[:len(password)-1]) + "🔥"
			// 	str := `<div id="inputEntry" class="form-control" style="position: absolute; width: 700px; background: transparent; display: flex;" contenteditable="true">` +
			// 		newPassword + `</div>`

			// 	tmpl, _ := template.New("t").Parse(str)
			// 	tmpl.Execute(w, str)

			// 	fire := []rune("🔥")[0]
			// 	password[len(password)-1] = fire

			// 	if highScore > 14 {
			// 		newExtra := r.PostFormValue("extra")
			// 		Rules[5].Extra = template.HTML(newExtra)
			// 	}

			// 	checkAndApply(w, r, password)
			// }

			return
		}

		// fmt.Println("Timer activated")

		m := regexp.MustCompile("[^🔥]🔥|🔥[^🔥]")
		newPassword := m.ReplaceAllString(string(password), "${1}🔥🔥$2")
		password = []rune(newPassword)

		str := `<div id="inputEntry" class="form-control" style="position: relative; max-width: 650px; background: transparent; word-wrap: break-word; display: inline-block;" contenteditable="true"
                        >` +
			newPassword + `</div>`

		tmpl, _ := template.New("t").Parse(str)
		tmpl.Execute(w, newPassword)

		if highScore > 14 {
			newExtra := r.PostFormValue("extra")
			Rules[5].Extra = template.HTML(newExtra)
		}

		checkAndApply(w, r, password)
	}

	timerPaul := func(w http.ResponseWriter, r *http.Request) {
		if !superPauled {
			return
		}

		password := []rune(r.PostFormValue("password"))
		var m *regexp.Regexp
		if difficulty == "easy" {
			m = regexp.MustCompile(`^(.*?)🐛(.*)`)
		} else if difficulty == "normal" {
			m = regexp.MustCompile(`^(.*?)🐛(.*?)🐛(.*?)`)
		} else { // difficulty == "hard"
			m = regexp.MustCompile(`^(.*?)🐛(.*?)🐛(.*?)🐛(.*?)`)
		}
		match := m.MatchString(string(password))
		if superPauled && !match {
			fmt.Println("failure")

			time := r.PostFormValue("time")

			str := `<div id="game-over" class="flex-column" style="position: absolute; width: 100%; height: 100%; display: flex;">
						<div class="row flex-grow-1 align-items-center justify-content-center" style="display: flex;font-size: 48;">
							<div style="position: absolute; width: 100%; height: 100%; background-color: grey; opacity: 0.5;"></div>
							<div class="align-self-center align-items-center justify-content-center" style="display: flex; height: 200px; background-color: black; color: red; position: absolute; width: 100%;">
								<div class="col align-self-center align-items-center justify-content-center">
									<div class="row row position align-self-center align-items-center justify-content-center"
									style="text-align: center; font-size: 320%;">
										PAUL IS STARVE
									</div>
									<div class="row position align-self-center align-items-center justify-content-center"
									style="font-size: 120%;">
										Your Time: ` + time + `
									</div>
									<div class="row position align-self-center align-items-center justify-content-center"
									style="font-size: 120%;">
										Best Time: ` + bestTime + `
									</div>
								</div>
							</div>

							<script>
								gameOver()
    						</script>
						</div>
					</div>`
			tmpl, _ := template.New("t").Parse(str)
			tmpl.Execute(w, str)
			highScore = 1
			superPauled = false
			combustible = false
			Rules[6].Emoji = "❌"
			Rules[9].Emoji = "❌"
			return
		}

		newPassword := m.ReplaceAllString(string(password), "${1}$2")
		str := `<div id="inputEntry" class="form-control" style="position: relative; max-width: 650px; background: transparent; word-wrap: break-word; display: inline-block;" contenteditable="true"
                        >` +
			newPassword + `</div>`

		tmpl, _ := template.New("t").Parse(str)
		tmpl.Execute(w, str)

		if highScore > 14 {
			newExtra := r.PostFormValue("extra")
			Rules[5].Extra = template.HTML(newExtra)
		}

		password = []rune(newPassword)
		checkAndApply(w, r, password)
	}

	sacrifice := func(w http.ResponseWriter, r *http.Request) {
		success := r.PostFormValue("sacrifice")

		if success != "true" {
			return
		}

		newExtra := r.PostFormValue("extra")
		// fmt.Println("New Extra")
		// fmt.Println(newExtra)
		Rules[5].Emoji = "✅"
		Rules[5].Extra = template.HTML(newExtra)

		password := []rune(r.PostFormValue("password"))
		checkAndApply(w, r, password)
	}

	// save := func(w http.ResponseWriter, r *http.Request) {

	// }

	// load := func(w http.ResponseWriter, r *http.Request) {
	// 	// saveFiles, _ := db.Query(`
	// 	// 	SELECT * FROM save
	// 	// `)
	// 	query, _ := db.Prepare(`INSERT INTO save (date, difficulty, time, entry, background, highscore, country1, country2, country3, captcha, sacrifice) VALUES
	// 	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	// 	`)
	// 	//date, diff, time, entr, bg, hs, c1, c2, c3, capt, sac
	// 	time := time.Now().Format("2006-01-02")
	// 	query.Exec(time, difficulty, "", "", "", highScore, country1.name, country2.name, country3.name, captcha.answer, "")
	// }

	http.HandleFunc("/", selectHandler)
	http.HandleFunc("/handle/", handler)
	http.HandleFunc("/check/", checkHelper)
	http.HandleFunc("/timerFire/", timerFire)
	http.HandleFunc("/reCaptcha/", reCaptcha)
	http.HandleFunc("/timerPaul/", timerPaul)
	http.HandleFunc("/sacrifice/", sacrifice)
	// http.HandleFunc("/save/", save)
	// http.HandleFunc("/load/", load)
	http.ListenAndServe(":1334", nil)
}

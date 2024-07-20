package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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
	Rules := []Rule{
		{Emoji: "❌", Text: "Your password must include the current time"},
		{Emoji: "❌", Text: "The length of your password must be a prime number"},
		{Emoji: "❌", Text: "Your password must include the length of your password"},
		{Emoji: "❌", Text: "At least 30% of your password must be in digits"},
		{Emoji: "❌", Text: "Your password must contain one of the following words: I want IRK | I need IRK | I love IRK"},
		{Emoji: "❌", Text: "A sacrifice must be made. Pick 2 letters that you will no longer be able to use"},
		{Emoji: "✅", Text: "🐔 Paul has hatched ! Please don't forget to feed him. He eats 1 🐛 every 20 second"},
		{Emoji: "❌", Text: "Your password must include a leap year"},
		{Emoji: "❌", Text: "Your password must include this CAPTCHA"},
		{Emoji: "❌", Text: "🥚 This is my chicken Paul. He hasn't hatched yet. Please put him in your password and keep him safe"},
		{Emoji: "❌", Text: "Oh no! Your password is on fire 🔥. Quick, put it out!"},
		{Emoji: "❌", Text: "The Roman numerals in your password should multiply to 35"},
		{Emoji: "❌", Text: "Your password must include one of this country"},
		{Emoji: "❌", Text: "Your password must include a Roman numeral"},
		{Emoji: "❌", Text: "Your password must include a month of the year"},
		{Emoji: "❌", Text: "The digits in your password must add up to 45"},
		{Emoji: "❌", Text: "Your password must include a special character (! @ # $ % ^ & * ( ) - _ = + \\ | [ ] { } ; : / ? . < > ' \")"},
		{Emoji: "❌", Text: "Your password must include an uppercase letter"},
		{Emoji: "❌", Text: "Your password must include a number"},
		{Emoji: "❌", Text: "Your password must be at least 5 characters"},
	}
	var country1 Country
	var country2 Country
	var country3 Country
	var captcha Captcha
	bestTime := "-"
	bestTimeInt := -1

	//--> Start of Database Setup <--//
	db, _ := sql.Open("sqlite3", "password_game.db")

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

	db.Close()
	//--> End of Database Setup <--//

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		pauled = false
		superPauled = false
		combustible = false
		highScore = 1
		for i := 0; i < len(Rules); i++ {
			Rules[i].Emoji = "❌"
			Rules[i].Extra = ""
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

	check := func(w http.ResponseWriter, r *http.Request, password []rune) {
		background := ""

		allCorrect := func() bool {
			for i := 1; i < highScore; i++ {
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

			now := time.Now().Format("15:04")
			// fmt.Println(now)

			match, _ := regexp.MatchString(now, string(password))

			if match {
				Rules[0].Emoji = "✅"

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
			} else {
				Rules[0].Emoji = "❌"
			}

		}

		Rule19 := func(password []rune) {
			if highScore < 19 {
				return
			}

			num := len(password)
			sqRoot := int(math.Sqrt(float64(num)))

			isPrime := true
			for i := 2; i <= sqRoot; i++ {
				if num%i == 0 {
					isPrime = false
				}
			}

			if isPrime {
				if highScore < 20 && allCorrect() {
					highScore = 20
				}
				Rules[1].Emoji = "✅"
			} else {
				Rules[1].Emoji = "❌"
			}

			Rule20(password)
		}

		Rule18 := func(password []rune) {
			if highScore < 18 {
				return
			}

			num := strconv.Itoa(len(password))
			match, _ := regexp.MatchString(num, string(password))

			if match {
				if highScore < 19 && allCorrect() {
					highScore = 19
				}
				Rules[2].Emoji = "✅"
			} else {
				Rules[2].Emoji = "❌"
			}

			Rule19(password)
		}

		Rule17 := func(password []rune) {
			if highScore < 17 {
				return
			}

			num := len(password) * 3 / 10
			ctr := 0

			for i := 0; i < len(password); i++ {
				c := password[i]
				if c >= '0' && c <= '9' {
					ctr++
				}
			}

			if ctr >= num {
				if highScore < 18 && allCorrect() {
					highScore = 18
				}
				Rules[3].Emoji = "✅"
			} else {
				Rules[3].Emoji = "❌"
			}

			Rule18(password)
		}

		Rule16 := func(password []rune) {
			if highScore < 16 {
				return
			}

			match, _ := regexp.MatchString("(I want IRK)|(I need IRK)|(I love IRK)", string(password))

			if match {
				if highScore < 17 && allCorrect() {
					highScore = 17
				}
				Rules[4].Emoji = "✅"
			} else {
				Rules[4].Emoji = "❌"
			}

			Rule17(password)
		}

		Rule15 := func(password []rune) {
			if highScore < 15 {
				return
			}

			if highScore < 16 && Rules[5].Emoji == "✅" {
				highScore = 16
			}

			Rule16(password)
		}

		Rule14 := func(password []rune) {
			if highScore < 14 {
				return
			}

			match, _ := regexp.MatchString("🐔", string(password))
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

				Rules[6].Emoji = "❌"
				Rules[9].Emoji = "❌"
				return
			}

			if highScore < 15 {
				highScore = 15
			}

			Rule15(password)
		}

		Rule13 := func(password []rune) {
			if highScore < 13 {
				return
			}

			hasLeap := false

			for i := 0; i < len(password); i++ {
				c := password[i]

				if c >= '0' && c <= '9' {
					if CheckLeap(0, password[i:]) {
						hasLeap = true
						break
					}
				}
			}

			if hasLeap {
				if highScore < 14 {
					r := regexp.MustCompile(`^(.*)🥚(.*)$`)
					newPassword := r.ReplaceAllString(string(password), "${1}🐔$2")

					str := `<div id="inputEntry" class="form-control" style="position: absolute; width: 700px; background: transparent; display: flex;" contenteditable="true">` +
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
								hx-vals='js:{password: getPassword()}'
								hx-ext="multi-swap">
								></div>
							</div>`
					tmpl, _ = template.New("t").Parse(str)
					tmpl.Execute(w, str)

					Rules[6].Emoji = "✅"

					highScore = 14
				}
				Rules[7].Emoji = "✅"
			} else {
				Rules[7].Emoji = "❌"
			}

			Rule14(password)
		}

		Rule12 := func(password []rune) {
			if highScore < 12 {
				return
			}

			fmt.Println("Captcha", captcha.answer)
			match, _ := regexp.MatchString(captcha.answer, string(password))
			if match {
				if highScore < 13 && allCorrect() {
					highScore = 13
				}
				Rules[8].Emoji = "✅"
			} else {
				Rules[8].Emoji = "❌"
			}

			Rule13(password)
		}

		Rule11 := func(password []rune) {
			if highScore < 11 {
				return
			}

			match, _ := regexp.MatchString("🥚", string(password))
			if pauled && !match {
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

				Rules[9].Emoji = "❌"
				return
			}

			if !pauled && match {
				// fmt.Println("pauled")
				if highScore < 12 && allCorrect() {
					highScore = 12

					captcha = captchas[rand.Intn(len(captchas))]

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
				}
				pauled = true
				Rules[9].Emoji = "✅"
			}
			// else {
			// 	Rules[9].Emoji = "❌"
			// }

			Rule12(password)
		}

		Rule10 := func(password []rune) {
			if highScore < 10 {
				return
			}

			match, _ := regexp.MatchString("🔥", string(password))
			if !match {
				// fmt.Println("check this")
				if highScore < 11 && allCorrect() {
					highScore = 11
				}
				Rules[10].Emoji = "✅"
			} else {
				// fmt.Println("onfire")
				Rules[10].Emoji = "❌"
			}

			Rule11(password)
		}

		Rule9 := func(password []rune) {
			if highScore < 9 {
				return
			}

			// (\s*I*\s+)*
			r := regexp.MustCompile(`^((I?[^IVXLCDM]+)*XXXV([^IVXLCDM]+I?)*)$|^((I?[^IVXLCDM]+)*V[^IVXLCDM]+(I?[^IVXLCDM]+)*VII([^IVXLCDM]+I?)*)$|^((I?[^IVXLCDM]+)*VII[^IVXLCDM]+(I?[^IVXLCDM]+)*V([^IVXLCDM]+I?)*)$`)
			match := r.MatchString(string(password))
			if match {
				if highScore < 10 && allCorrect() {
					newPassword := string(password[:len(password)-1]) + "🔥"
					str := `<div id="inputEntry" class="form-control" style="position: absolute; width: 700px; background: transparent; display: flex;" contenteditable="true">` +
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
								hx-ext="multi-swap">
								></div>
							</div>`
					tmpl, _ = template.New("t").Parse(str)
					tmpl.Execute(w, str)

					combustible = true
				}
				Rules[11].Emoji = "✅"
			} else {
				Rules[11].Emoji = "❌"

				m := regexp.MustCompile(`[IVXLCDM]`)
				background = m.ReplaceAllString(string(password), `<span style="background-color: firebrick;">${0}</span>`)
			}

			Rule10(password)
		}

		Rule8 := func(password []rune) {
			if highScore < 8 {
				return
			}

			fmt.Println("Country", country1.name, country2.name, country3.name)
			str := `(?i)(` + country1.name + `)|(` + country2.name + `)|(` + country3.name + `)`
			r := regexp.MustCompile(str)
			match := r.MatchString(string(password))
			if match {
				if highScore < 9 && allCorrect() {
					highScore = 9
				}
				Rules[12].Emoji = "✅"
			} else {
				Rules[12].Emoji = "❌"
			}

			Rule9(password)
		}

		Rule7 := func(password []rune) {
			if highScore < 7 {
				return
			}

			r := regexp.MustCompile(`[IVXLCDM]`)
			match := r.MatchString(string(password))
			if match {
				if highScore < 8 && allCorrect() {
					highScore = 8

					country1 = countries[rand.Intn(len(countries))]
					country2 = countries[rand.Intn(len(countries))]
					for country1.name == country2.name {
						country2 = countries[rand.Intn(len(countries))]
					}
					country3 = countries[rand.Intn(len(countries))]
					for country1.name == country3.name || country2.name == country3.name {
						country3 = countries[rand.Intn(len(countries))]
					}

					str := `<div class="row justify-content-center m1-3">
								<div class="col-3"><img src="` + country1.flag + `" width="96" height="64"></div>
								<div class="col-3"><img src="` + country2.flag + `" width="96" height="64"></div>
								<div class="col-3"><img src="` + country3.flag + `" width="96" height="64"></div>
							</div>`

					Rules[12].Extra = template.HTML(str)
				}
				Rules[13].Emoji = "✅"
			} else {
				Rules[13].Emoji = "❌"
			}

			Rule8(password)
		}

		Rule6 := func(password []rune) {
			if highScore < 6 {
				return
			}

			r := regexp.MustCompile(`(?i)(january)|(february)|(march)|(april)|(may)|(june)|(july)|(august)|(september)|(october)|(november)|(december)`)
			match := r.MatchString(string(password))
			if match {
				if highScore < 7 && allCorrect() {
					highScore = 7
				}
				Rules[14].Emoji = "✅"
			} else {
				Rules[14].Emoji = "❌"
			}

			Rule7(password)
		}

		Rule5 := func(password []rune) {
			if highScore < 5 {
				return
			}

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

			if total == 45 {
				if highScore < 6 && allCorrect() {
					highScore = 6
				}
				Rules[15].Emoji = "✅"

			} else {
				Rules[15].Emoji = "❌"

				m := regexp.MustCompile("([0-9])+")
				background = m.ReplaceAllString(string(password), `<span style="background-color: firebrick;">${0}</span>`)
			}

			Rule6(password)
		}

		Rule4 := func(password []rune) {
			if highScore < 4 {
				return
			}

			match, _ := regexp.MatchString("[!@#$%^&*()\\-_=+\\\\|\\[\\]{};:\\/?.<>'\"]", string(password))

			if match {
				if highScore < 5 && allCorrect() {
					highScore = 5
				}
				Rules[16].Emoji = "✅"
			} else {
				Rules[16].Emoji = "❌"
			}

			Rule5(password)
		}

		Rule3 := func(password []rune) {
			if highScore < 3 {
				return
			}

			match, _ := regexp.MatchString("[A-Z]", string(password))

			if match {
				if highScore < 4 && allCorrect() {
					highScore = 4
				}
				Rules[17].Emoji = "✅"
			} else {
				Rules[17].Emoji = "❌"
			}

			Rule4(password)
		}

		Rule2 := func(password []rune) {
			if highScore < 2 {
				return
			}

			match, _ := regexp.MatchString("[0-9]", string(password))

			if match {
				if highScore < 3 && allCorrect() {
					highScore = 3
				}
				Rules[18].Emoji = "✅"
			} else {
				Rules[18].Emoji = "❌"
			}

			Rule3(password)
		}

		Rule1 := func(password []rune) {
			if len(password) >= 5 {
				if highScore < 2 {
					highScore = 2
				}
				Rules[19].Emoji = "✅"
			} else {
				Rules[19].Emoji = "❌"
			}

			Rule2(password)
		}

		background = string(password)
		passLength := len(password)
		log.Print(passLength)

		Rule1(password)

		tmpl := template.Must(template.ParseFiles("index.html"))

		rules := map[string][]Rule{
			"Rules": Rules[(20 - highScore):],
		}

		tmpl.ExecuteTemplate(w, "rule-list-element", rules)

		str := `<label id="Length" for="Length" class="form-label">` + strconv.Itoa(passLength) + `</label>`
		tmpl, _ = template.New("t").Parse(str)
		tmpl.Execute(w, passLength)

		fmt.Println("background:", background)
		fmt.Println("password", string(password))

		str = `<div id="inputBackground">` +
			background + `</div>`
		tmpl, _ = template.New("t").Parse(str)
		tmpl.Execute(w, background)
	}

	timerFire := func(w http.ResponseWriter, r *http.Request) {
		if !combustible {
			return
		}

		password := []rune(r.PostFormValue("password"))
		match, _ := regexp.MatchString("🔥", string(password))
		if !match {
			// 1/30 chance to occur again
			if rand.Intn(30) == 1 {
				newPassword := string(password[:len(password)-1]) + "🔥"
				str := `<div id="inputEntry" class="form-control" style="position: absolute; width: 700px; background: transparent; display: flex;" contenteditable="true">` +
					newPassword + `</div>`

				tmpl, _ := template.New("t").Parse(str)
				tmpl.Execute(w, str)

				fire := []rune("🔥")[0]
				password[len(password)-1] = fire
			}
			check(w, r, password)
			return
		}

		// fmt.Println("Timer activated")

		m := regexp.MustCompile("[^🔥]🔥|🔥[^🔥]")
		newPassword := m.ReplaceAllString(string(password), "${1}🔥🔥$2")
		password = []rune(newPassword)

		str := `<div id="inputEntry" class="form-control" style="position: absolute; width: 700px; background: transparent; display: flex;" contenteditable="true">` +
			newPassword + `</div>`

		tmpl, _ := template.New("t").Parse(str)
		tmpl.Execute(w, newPassword)

		check(w, r, password)
	}

	checkHelper := func(w http.ResponseWriter, r *http.Request) {
		log.Print("Request received")
		passwordStr := r.PostFormValue("password")

		match, _ := regexp.MatchString("^cheat$", passwordStr)
		if match {
			// cheat(w, r)

			return
		}

		password := []rune(passwordStr)
		log.Print(string(password))

		check(w, r, password)
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
                    hx-vals='js:{password: getPassword()}'
                    hx-ext="multi-swap"
					>
						<font size="5">🔄</font>
					</button>
				</div>
			</div>`

		Rules[8].Extra = template.HTML(str)

		password := []rune(r.PostFormValue("password"))
		check(w, r, password)
	}

	timerPaul := func(w http.ResponseWriter, r *http.Request) {
		if !superPauled {
			return
		}

		password := []rune(r.PostFormValue("password"))
		m := regexp.MustCompile(`^(.*?)🐛(.*)`)
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
			superPauled = false
			Rules[6].Emoji = "❌"
			Rules[9].Emoji = "❌"
			return
		}

		newPassword := m.ReplaceAllString(string(password), "${1}$2")
		str := `<div id="inputEntry" class="form-control" style="position: absolute; width: 700px; background: transparent; display: flex;" contenteditable="true">` +
			newPassword + `</div>`

		tmpl, _ := template.New("t").Parse(str)
		tmpl.Execute(w, str)

		password = []rune(newPassword)
		check(w, r, password)
	}

	sacrifice := func(w http.ResponseWriter, r *http.Request) {
		success := r.PostFormValue("sacrifice")

		if success != "true" {
			return
		}

		newExtra := r.PostFormValue("extra")
		fmt.Println("New Extra")
		fmt.Println(newExtra)
		Rules[5].Emoji = "✅"
		Rules[5].Extra = template.HTML(newExtra)

		password := []rune(r.PostFormValue("password"))
		check(w, r, password)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/check/", checkHelper)
	http.HandleFunc("/timerFire/", timerFire)
	http.HandleFunc("/reCaptcha/", reCaptcha)
	http.HandleFunc("/timerPaul/", timerPaul)
	http.HandleFunc("/sacrifice/", sacrifice)
	http.ListenAndServe(":1334", nil)
}

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

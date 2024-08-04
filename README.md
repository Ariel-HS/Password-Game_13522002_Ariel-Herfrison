# PASSWORD GAME
## Description
A redux of the neal.fun/password-game with slight changes to the rules. Make a password with up to 20 rules. Go through all 20 rules to win the game. Some rules will cause you to lose instantly. Type "cheat" into the entrybox to see an answer.

## Getting Started
### Dependencies
* Go Language
* Bootstrap 5
* HTMX
* Docker
* sqlite3

### Executing program
* Run in localhost
  1. Open the terminal from the root directory
  ```
  cd Password-Game_13522002_Ariel-Herfrison
  ```
  2. Build the docker image
  ```
  docker build . -t [image-name]
  ```
  3. Run the docker container
  ```
  docker run -p 1334:1334 [image-name]
  ```
  *make sure docker desktop is running
* Run from the internet
  1. Open the following link: https://passwordgamelite-nvvtv47y.b4a.run/

## Screenshots
* Difficulty screen
* Main screen
* Game over screen
* Win screen
  
## Authors
Ariel Herfrison

## Acknowledgments
* [Password-Game](https://neal.fun/password-game/)

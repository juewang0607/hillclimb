package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/emedvedev/enigma"
)

// construct of the enigma
var enigmaSettings = struct {
	Reflector string
	Rings     []int
	Positions []byte
	Rotors    []string
}{
	Reflector: "C-thin",
	Rings:     []int{1, 1, 1, 16},
	Positions: []byte{'A', 'A', 'B', 'Q'},
	Rotors:    []string{"Beta", "II", "IV", "III"},
}

var character = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var rotor1_set = []string{"Beta", "Gamma", "I", "II", "V", "VI"}
var rotor2_set = []string{"I", "II", "V", "VI"}
var rotor_config = make([]enigma.RotorConfig, 4)
var rotor1 = enigmaSettings.Rotors[0]
var rotor2 = enigmaSettings.Rotors[1]
var position1 = enigmaSettings.Positions[0]
var position2 = enigmaSettings.Positions[1]

// function to caculate ICscore
func IocScore(ciphertext string) float64 {
	ioc_score := 0.0
	for i := 0; i < len(character); i++ {
		times_occur := float64(strings.Count(ciphertext, string(character[i])))
		ioc_score += (times_occur * (times_occur - 1))
	}
	ioc_score = float64(ioc_score / (float64(len(ciphertext)) * (float64(len(ciphertext) - 1))))
	return ioc_score
}

// function to swap character
func swap_character(old string, new string, tempString string) string {
	var temp string
	temp = old
	tempString = strings.ReplaceAll(tempString, new, "-")
	tempString = strings.ReplaceAll(tempString, old, new)
	tempString = strings.ReplaceAll(tempString, "-", temp)
	return tempString
}
func lexicographical(currentPlugboard string) []string {
	var newPlugboard []string
	for i := 0; i < len(currentPlugboard); i++ {
		// if the character has connected with others, it will be replaced as "*"
		if character[i] != currentPlugboard[i] && string(character[i]) != "*" && string(currentPlugboard[i]) != "*" {
			newPlugboard = append(newPlugboard, string(character[i])+string(currentPlugboard[i]))
			character = strings.Replace(character, string(character[i]), "*", -1)
			character = strings.Replace(character, string(currentPlugboard[i]), "*", -1)
			currentPlugboard = strings.Replace(currentPlugboard, string(character[i]), "*", -1)
			//fmt.Println(currentPlugboard)
			currentPlugboard = strings.Replace(currentPlugboard, string(currentPlugboard[i]), "*", -1)
			//fmt.Println(currentPlugboard)
		}
		//fmt.Println(newPlugboard)
	}
	return newPlugboard
}

// function to generate a plugboard
func generate_plugboard(currentPlugboard string) []string {

	var i int
	var enigmaPlugboardSetting []string
	var tempDefaultAlphabets string
	tempDefaultAlphabets = character
	for i = 0; i < len(currentPlugboard); i++ {
		if currentPlugboard[i] != tempDefaultAlphabets[i] && string(currentPlugboard[i]) != "-" && string(tempDefaultAlphabets[i]) != "-" {
			enigmaPlugboardSetting = append(enigmaPlugboardSetting, string(currentPlugboard[i])+string(tempDefaultAlphabets[i]))
			var x = string(currentPlugboard[i])
			var y = string(tempDefaultAlphabets[i])
			currentPlugboard = strings.ReplaceAll(currentPlugboard, x, "-")
			currentPlugboard = strings.ReplaceAll(currentPlugboard, y, "-")
			tempDefaultAlphabets = strings.ReplaceAll(tempDefaultAlphabets, x, "-")
			tempDefaultAlphabets = strings.ReplaceAll(tempDefaultAlphabets, y, "-")
		}
	}
	//fmt.Println(enigmaPlugboardSetting)
	return enigmaPlugboardSetting
}

// we need to find the plugboard with highest score
func hillclimbing(data string) string {

	var i int
	var j int
	var tempPlugboard string
	var tempPlugboard1 string
	var tempPlugboard2 string
	var tempPlugboard3 string
	var tempPlugboard4 string
	var currentPlugboard string
	var bestPlugboard string
	var currentTempPlugboard string
	var IOC float64
	var maxIOC float64
	var enigmaPlugboard []string
	//var decrypted string
	bestPlugboard = character
	currentPlugboard = character
	maxIOC = 0.0
	for i = 0; i < 26; i++ {

		currentPlugboard = bestPlugboard
		IOC = 0.0
		for j = i + 1; j < 26; j++ {
			// fmt.Println(currentPlugboard[j], character[j])
			if string(currentPlugboard[j]) != string(character[j]) {
				tempPlugboard = swap_character(string(character[j]), string(currentPlugboard[j]), currentPlugboard)
				tempPlugboard = swap_character(string(character[i]), string(currentPlugboard[i]), tempPlugboard)
				tempPlugboard1 = swap_character(string(character[i]), string(currentPlugboard[i]), tempPlugboard)
				tempPlugboard2 = swap_character(string(character[j]), string(currentPlugboard[j]), tempPlugboard)
				tempPlugboard3 = swap_character(string(character[i]), string(currentPlugboard[i]), tempPlugboard)
				tempPlugboard4 = swap_character(string(character[j]), string(currentPlugboard[j]), tempPlugboard)
				var IOC1 = IocScore(strings.ReplaceAll(data, string(character[i]), string(currentPlugboard[i])))
				var IOC2 = IocScore(strings.ReplaceAll(data, string(character[j]), string(currentPlugboard[j])))
				var IOC3 = IocScore(strings.ReplaceAll(data, string(character[i]), string(currentPlugboard[i])))
				var IOC4 = IocScore(strings.ReplaceAll(data, string(character[j]), string(currentPlugboard[j])))

				if math.Max(math.Max(IOC, IOC2), math.Max(IOC3, IOC4)) == IOC1 {
					IOC = IOC1
					tempPlugboard = tempPlugboard1
				} else if math.Max(math.Max(IOC1, IOC2), math.Max(IOC3, IOC4)) == IOC2 {
					IOC = IOC2
					tempPlugboard = tempPlugboard2
				} else if math.Max(math.Max(IOC1, IOC2), math.Max(IOC3, IOC4)) == IOC3 {
					IOC = IOC3
					tempPlugboard = tempPlugboard3
				} else {
					IOC = IOC4
					tempPlugboard = tempPlugboard4
				}
			} else {
				tempPlugboard = swap_character(string(character[i]), string(currentPlugboard[j]), currentPlugboard)
				enigmaSettings.Rotors[0] = rotor1
				enigmaSettings.Rotors[1] = rotor2
				enigmaSettings.Positions[0] = position1
				enigmaSettings.Positions[1] = position2
				enigmaPlugboard = generate_plugboard(tempPlugboard)
				rotor_config[0] = enigma.RotorConfig{enigmaSettings.Rotors[0], enigmaSettings.Positions[0], 1}
				rotor_config[1] = enigma.RotorConfig{enigmaSettings.Rotors[1], enigmaSettings.Positions[1], 1}
				rotor_config[2] = enigma.RotorConfig{"IV", 'B', 1}
				rotor_config[3] = enigma.RotorConfig{"III", 'Q', 16}
				Enigma := enigma.NewEnigma(rotor_config, "C-thin", enigmaPlugboard)
				PT := Enigma.EncodeString(enigma.SanitizePlaintext(data))
				//decrypted = strings.ReplaceAll(data, string(character[i]), string(currentPlugboard[j]))
				IOC = IocScore(PT)
			}
			if IOC > maxIOC {
				maxIOC = IOC
				currentTempPlugboard = tempPlugboard
				// fmt.Println(IOC, tempPlugboard)
			}
		}
		bestPlugboard = currentTempPlugboard
	}
	return bestPlugboard
}

var trigramScores = make(map[string]float64)

// function to caculate trigram
func caculate_trigram_score(text, plugboard string) float64 {
	var score = 0.0
	var i int
	var enigmaPlugboard = generate_plugboard(plugboard)
	var trigramScores = make(map[string]float64)
	var trigram_pair []string
	var total_score float64
	//initialize the setting of enigma
	rotor_config[0] = enigma.RotorConfig{rotor1, position1, 1}
	rotor_config[1] = enigma.RotorConfig{rotor2, position2, 1}
	rotor_config[2] = enigma.RotorConfig{"IV", 'B', 1}
	rotor_config[3] = enigma.RotorConfig{"III", 'Q', 16}
	Enigma := enigma.NewEnigma(rotor_config, "C-thin", enigmaPlugboard)
	PT := Enigma.EncodeString(enigma.SanitizePlaintext(text))
	//fmt.Println(PT)
	// read the trigram file
	result, err := ioutil.ReadFile("english_trigrams.txt")
	if err != nil {
		log.Fatal(err)
	}
	trigram_pair = strings.Split(string(result), "\n")
	for i := 0; i < len(trigram_pair)-1; i++ {
		//fmt.Println(trigram_pair[i])
		frequency := strings.Split(trigram_pair[i], " ")
		//fmt.Println(frequency)
		// convert the frenquency to an int
		var value, err = strconv.Atoi(frequency[1])
		if err != nil {
			log.Fatal(err)
		}
		trigramScores[frequency[0]] = float64(value)
		total_score += float64(value)
	}
	for i := 0; i < len(trigram_pair)-1; i++ {
		var k = strings.Split(trigram_pair[i], " ")
		trigramScores[k[0]] = math.Log(trigramScores[k[0]] / total_score)
	}
	for i = 0; i < len(PT)-3; i++ {
		score += float64(trigramScores[PT[i:i+3]])
	}
	return score
}
func main() {
	//fmt.Println(rotor1)
	//fmt.Println(rotor2)
	//fmt.Println(position1)
	//fmt.Println(position2)
	//read files
	var plaintext string

	fileIO, err := os.OpenFile(os.Args[1], os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	defer fileIO.Close()
	rawBytes, err := ioutil.ReadAll(fileIO)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(rawBytes), "\n")
	for i, line := range lines {
		if i == 0 {
			plaintext = string(line)
		}
	}
	// declare plugboard
	var best_plugboard string
	var current_plugboard string
	// declare trigram score
	var trigram_score float64
	// declare the original plaintext

	best_score := -100000000.0
	// There are 8 rotors in total
	// Sincere the last two will not be overwriten
	// The first two rotors have 30 ways wo combine
	// For example: "Beta" "I" "IV" "III"
	// "Beta" "II" "IV" "III"
	// We will have 6 cases
	// dummy value will be overwritten

	var best_rotor1 string
	var best_rotor2 string
	var best_position1 byte
	var best_position2 byte
	// Then we declare the start_position
	// Same as the rotor, the start+position of the third and fourth rotor are

	// 6 rotors except "IV" and "III" which are placed into the rotor3 and totor4
	var rotor1_set = []string{"Beta", "Gamma", "I", "II", "V", "VI"}
	var rotor2_set = []string{"I", "II", "V", "VI"}

	// 26 characters stored into a string
	var start_pos = []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I',
		'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	// rotor1
	for i := 0; i < len(rotor1_set); i++ {

		rotor1 = rotor1_set[i]
		// rotor1
		for j := 0; j < len(rotor2_set); j++ {
			//fmt.Println(rotor2)
			rotor2 = rotor2_set[j]
			//fmt.Println(rotor2)
			fmt.Println("check rotor: "+rotor1, rotor2)
			// position1
			for m := 0; m < len(start_pos); m++ {

				position1 = start_pos[m]
				// position2
				for n := 0; n < len(start_pos); n++ {

					position2 = start_pos[n]
					//hillclimb the plaintext
					current_plugboard = hillclimbing(plaintext)
					//fmt.Println(plaintext)
					//fmt.Println(current_plugboard)
					// caculate the trigram score
					trigram_score = caculate_trigram_score(plaintext, current_plugboard)
					//fmt.Println(trigram_score)
					if trigram_score > best_score {
						best_score = trigram_score
						best_plugboard = current_plugboard
						best_rotor1 = enigmaSettings.Rotors[0]
						best_rotor2 = enigmaSettings.Rotors[1]
						best_position1 = enigmaSettings.Positions[0]
						best_position2 = enigmaSettings.Positions[1]
					}
				}
			}
		}
	}

	//var bestEnigmaPlugboard = generate_plugboard(best_plugboard)
	//rotor_config[0] = enigma.RotorConfig{best_rotor1, best_position1, 1}
	//rotor_config[1] = enigma.RotorConfig{best_rotor2, best_position2, 1}
	//rotor_config[2] = enigma.RotorConfig{"IV", 'B', 1}
	//rotor_config[3] = enigma.RotorConfig{"III", 'Q', 16}
	//Enigma := enigma.NewEnigma(rotor_config, "C-thin", bestEnigmaPlugboard)
	//PT := Enigma.EncodeString(enigma.SanitizePlaintext(plaintext))
	fmt.Println(best_rotor1, best_rotor2, "IV", "III")
	fmt.Println(string(best_position1), string(best_position2), string('B'), string('Q'))
	fmt.Println(strings.Join(lexicographical(best_plugboard), " "))
	//fmt.Println(PT)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type request struct {
	Player int `json:"player"`
	Dice   int `json:"dice"`
}

type response struct {
	Message string `json:"message"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/dicegame", startGame).Methods("POST")

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}

func startGame(w http.ResponseWriter, r *http.Request) {
	var (
		err                      error
		reqData                  *request
		statusCode, winningPoint int
		winnerArr                []int
		winnerStr, respMsg       string
	)

	// parse json from request body
	err = json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		statusCode = http.StatusBadRequest
		respMsg = "error when reading data"
	}

	fmt.Println()
	fmt.Println("Start The Game")
	winnerArr, winningPoint = processGame(reqData.Player, reqData.Dice)

	for _, v := range winnerArr {
		winnerStr += strconv.Itoa(v) + ","
	}
	winnerStr = strings.TrimSuffix(winnerStr, ",")

	respMsg = "Player " + winnerStr + " wins the game with " + strconv.Itoa(winningPoint) + " point."

	response := &response{
		Message: respMsg,
	}
	statusCode = http.StatusOK

	// convert data into json and send as response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func processGame(playerCount, diceCount int) (winners []int, winningPoint int) {
	var (
		playersDiceCount []int
		playersScore     []int
	)

	additionalDices := make(map[int]int)
	for i := 0; i < playerCount; i++ {
		playersScore = append(playersScore, 0)
		playersDiceCount = append(playersDiceCount, diceCount)
	}

	iteration := 0
	for playerCount > 1 {
		for ind, v := range playersDiceCount {
			fmt.Println("Player " + strconv.Itoa(ind))
			if v == 0 {
				fmt.Println("-")
				continue
			}

			for itr := 1; itr <= v; itr++ {
				dice := rand.Intn(6) + 1

				if dice == 1 {
					// if dice shows 1 then player's dice count would be decreased, adding additional dices to next player
					addInd := ind + 1
					playersDiceCount[ind]--
					if addInd == len(playersDiceCount) {
						addInd = 0
					}
					additionalDices[addInd]++
				} else if dice == 6 {
					// if dice shows 6 then player's score would increased and it's dice count would be decreased
					playersDiceCount[ind]--
					playersScore[ind]++
				}
				fmt.Print(strconv.Itoa(dice) + ";")
			}
			fmt.Println()
		}

		playerCount = 0
		// recount players that is still playing in the next iteration
		for ind, v := range playersDiceCount {
			playersDiceCount[ind] = v + additionalDices[ind]
			additionalDices[ind] = 0

			if playersDiceCount[ind] > 0 {
				playerCount++
			}
		}
		fmt.Println("players score :", playersScore, "players dice count :", playersDiceCount, "player count :", playerCount, "iteration :", iteration)
		iteration++
	}

	for i, v := range playersScore {
		if v > winningPoint {
			winningPoint = v
			winners = nil
			winners = append(winners, i+1)
		} else if v == winningPoint {
			winners = append(winners, i+1)
		}
	}

	return winners, winningPoint
}

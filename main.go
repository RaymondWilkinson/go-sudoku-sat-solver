package main

import (
	"encoding/json"
	"github.com/akrylysov/algnhsa"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sudoku-sat-solver/sudoku"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/generate", generateHandler)
	mux.HandleFunc("/solve", solveHandler)

	handler := cors.AllowAll().Handler(mux)

	if region := os.Getenv("AWS_REGION"); len(region) > 0 {
		algnhsa.ListenAndServe(handler, nil)
	} else {
		if err := http.ListenAndServe(":9990", handler); err != nil {
			log.Fatal(err)
		}
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	difficulty := r.URL.Query().Get("difficulty")

	if len(difficulty) == 0 {
		difficulty = "easy"
	}

	result, err := sudoku.Generate(difficulty)
	if err != nil {
		JSONError(w, err.Error(), 422)
		return
	}

	JSONSuccess(w, map[string]string{
		"grid": result.String(),
	})
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	gridString := r.URL.Query().Get("grid")

	if len(gridString) == 0 {
		JSONError(w, "A Sudoku grid is required.", 422)
		return
	}

	result, err := sudoku.SolveFromString(gridString)
	if err != nil {
		JSONError(w, err.Error(), 422)
		return
	}

	JSONSuccess(w, map[string]string{
		"grid": result.String(),
	})
}

func JSONSuccess(w http.ResponseWriter, data map[string]string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(data)
}

func JSONError(w http.ResponseWriter, errMsg string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"message": errMsg,
	})
}

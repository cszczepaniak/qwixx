package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/cszczepaniak/qwixx/internal/game"
	"github.com/starfederation/datastar-go/datastar"
)

// Signals matches the client-side store (game + backend-derived state + dialog state).
type Signals struct {
	Game                  game.State `json:"game"`
	DisabledIndices       [4]int     `json:"disabledIndices"`
	UnlockedLocks         [4]bool    `json:"unlockedLocks"`
	Score                 int        `json:"score"`
	ConfirmRow            int        `json:"confirmRow"`
	ConfirmCol            int        `json:"confirmCol"`
	ShowConfirmClearCross bool       `json:"showConfirmClearCross"`
	ConfirmClearAll       bool       `json:"confirmClearAll"`
	ShowScore             bool       `json:"showScore"`
}

// PopulateDerived fills DisabledIndices, UnlockedLocks, and Score from the current Game state.
func (s *Signals) PopulateDerived() {
	for i := range 4 {
		s.DisabledIndices[i] = s.Game.DisabledIndex(i)
		s.UnlockedLocks[i] = s.Game.UnlockedLock(i)
	}
	s.Score = s.Game.Score()
}

func SetCross(w http.ResponseWriter, r *http.Request) {
	var sigs Signals
	if err := datastar.ReadSignals(r, &sigs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	row, _ := strconv.Atoi(r.URL.Query().Get("row"))
	col, _ := strconv.Atoi(r.URL.Query().Get("col"))
	sigs.Game.SetCross(row, col)
	sigs.PopulateDerived()
	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(sigs); err != nil {
		log.Printf("patch signals: %v", err)
		return
	}
}

func UnsetCross(w http.ResponseWriter, r *http.Request) {
	var sigs Signals
	if err := datastar.ReadSignals(r, &sigs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	row, _ := strconv.Atoi(r.URL.Query().Get("row"))
	col, _ := strconv.Atoi(r.URL.Query().Get("col"))
	sigs.Game.UnsetCross(row, col)
	sigs.ConfirmRow = -1
	sigs.ConfirmCol = -1
	sigs.ShowConfirmClearCross = false
	sigs.PopulateDerived()
	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(sigs); err != nil {
		log.Printf("patch signals: %v", err)
		return
	}
}

func SetMissed(w http.ResponseWriter, r *http.Request) {
	var sigs Signals
	if err := datastar.ReadSignals(r, &sigs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	i, _ := strconv.Atoi(r.URL.Query().Get("i"))
	sigs.Game.SetMissed(i)
	sigs.PopulateDerived()
	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(sigs); err != nil {
		log.Printf("patch signals: %v", err)
		return
	}
}

func UnsetMissed(w http.ResponseWriter, r *http.Request) {
	var sigs Signals
	if err := datastar.ReadSignals(r, &sigs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	i, _ := strconv.Atoi(r.URL.Query().Get("i"))
	sigs.Game.UnsetMissed(i)
	sigs.PopulateDerived()
	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(sigs); err != nil {
		log.Printf("patch signals: %v", err)
		return
	}
}

func ClearAll(w http.ResponseWriter, r *http.Request) {
	var sigs Signals
	if err := datastar.ReadSignals(r, &sigs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sigs.Game.ClearAll()
	sigs.ConfirmClearAll = false
	sigs.PopulateDerived()
	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(sigs); err != nil {
		log.Printf("patch signals: %v", err)
		return
	}
}

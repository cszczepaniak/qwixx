package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cszczepaniak/qwixx/internal/game"
	"github.com/starfederation/datastar-go/datastar"
)

// HTTPError carries an HTTP status code for the adapter to send.
type HTTPError struct {
	Code int
	Err  error
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return http.StatusText(e.Code)
}

func (e *HTTPError) Unwrap() error { return e.Err }

// Error returns an error that the adapter will send as the given HTTP status code.
func Error(code int, err error) error {
	return &HTTPError{Code: code, Err: err}
}

// HandlerFunc is a handler that returns an error. The router adapts it to http.HandlerFunc.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

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

func parseRowCol(r *http.Request) (row, col int, err error) {
	row, err = strconv.Atoi(r.URL.Query().Get("row"))
	if err != nil {
		return 0, 0, Error(http.StatusBadRequest, fmt.Errorf("invalid row: %w", err))
	}
	col, err = strconv.Atoi(r.URL.Query().Get("col"))
	if err != nil {
		return 0, 0, Error(http.StatusBadRequest, fmt.Errorf("invalid col: %w", err))
	}
	if row < 0 || row >= 4 || col < 0 || col >= 12 {
		return 0, 0, Error(http.StatusBadRequest, errors.New("row must be 0-3, col must be 0-11"))
	}
	return row, col, nil
}

func parseMissedIndex(r *http.Request) (int, error) {
	i, err := strconv.Atoi(r.URL.Query().Get("i"))
	if err != nil {
		return 0, Error(http.StatusBadRequest, fmt.Errorf("invalid index: %w", err))
	}
	if i < 0 || i >= 4 {
		return 0, Error(http.StatusBadRequest, errors.New("index must be 0-3"))
	}
	return i, nil
}

func SetCross(w http.ResponseWriter, r *http.Request) error {
	var sigs Signals
	if err := datastar.ReadSignals(r, &sigs); err != nil {
		return Error(http.StatusBadRequest, err)
	}
	row, col, err := parseRowCol(r)
	if err != nil {
		return err
	}
	sigs.Game.SetCross(row, col)
	sigs.PopulateDerived()
	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(sigs); err != nil {
		return fmt.Errorf("patch signals: %w", err)
	}
	return nil
}

func UnsetCross(w http.ResponseWriter, r *http.Request) error {
	var sigs Signals
	if err := datastar.ReadSignals(r, &sigs); err != nil {
		return Error(http.StatusBadRequest, err)
	}
	row, col, err := parseRowCol(r)
	if err != nil {
		return err
	}
	sigs.Game.UnsetCross(row, col)
	sigs.ConfirmRow = -1
	sigs.ConfirmCol = -1
	sigs.ShowConfirmClearCross = false
	sigs.PopulateDerived()
	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(sigs); err != nil {
		return fmt.Errorf("patch signals: %w", err)
	}
	return nil
}

func SetMissed(w http.ResponseWriter, r *http.Request) error {
	var sigs Signals
	if err := datastar.ReadSignals(r, &sigs); err != nil {
		return Error(http.StatusBadRequest, err)
	}
	i, err := parseMissedIndex(r)
	if err != nil {
		return err
	}
	sigs.Game.SetMissed(i)
	sigs.PopulateDerived()
	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(sigs); err != nil {
		return fmt.Errorf("patch signals: %w", err)
	}
	return nil
}

func UnsetMissed(w http.ResponseWriter, r *http.Request) error {
	var sigs Signals
	if err := datastar.ReadSignals(r, &sigs); err != nil {
		return Error(http.StatusBadRequest, err)
	}
	i, err := parseMissedIndex(r)
	if err != nil {
		return err
	}
	sigs.Game.UnsetMissed(i)
	sigs.PopulateDerived()
	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(sigs); err != nil {
		return fmt.Errorf("patch signals: %w", err)
	}
	return nil
}

func ClearAll(w http.ResponseWriter, r *http.Request) error {
	var sigs Signals
	if err := datastar.ReadSignals(r, &sigs); err != nil {
		return Error(http.StatusBadRequest, err)
	}
	sigs.Game.ClearAll()
	sigs.ConfirmClearAll = false
	sigs.PopulateDerived()
	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(sigs); err != nil {
		return fmt.Errorf("patch signals: %w", err)
	}
	return nil
}

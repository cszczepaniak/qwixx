package views

import "fmt"

func LockStripeShowExpr(rowIdx int) string {
	return fmt.Sprintf("!$unlockedLocks[%d]", rowIdx)
}

func CrossShowExpr(rowIdx, colIdx int) string {
	return fmt.Sprintf("$game.crosses[%d][%d]", rowIdx, colIdx)
}

func ConfirmCrossClickExpr(rowIdx, colIdx int) string {
	return fmt.Sprintf("$showConfirmClearCross = true; $confirmRow = %d; $confirmCol = %d", rowIdx, colIdx)
}

func DisabledExpr(rowIdx, colIdx int) string {
	return fmt.Sprintf("%d <= $disabledIndices[%d]", colIdx, rowIdx)
}

func SetCrossURL(rowIdx, colIdx int) string {
	return fmt.Sprintf("@post('/action/set-cross?row=%d&col=%d')", rowIdx, colIdx)
}

func MissedShowExpr(i int) string {
	return fmt.Sprintf("$game.missedRolls[%d]", i)
}

func SetMissedURL(i int) string {
	return fmt.Sprintf("@post('/action/set-missed?i=%d')", i)
}

func UnsetMissedURL(i int) string {
	return fmt.Sprintf("@post('/action/unset-missed?i=%d')", i)
}

func UnsetCrossConfirmURL() string {
	return "@post('/action/unset-cross?row=' + $confirmRow + '&col=' + $confirmCol)"
}

func ClearAllConfirmURL() string {
	return "@post('/action/clear-all'); $confirmClearAll = false"
}

var rowColors = []struct {
	mg, fg, bg, text string
}{
	{"bg-red-500", "bg-red-800", "bg-red-100", "text-red-500"},
	{"bg-yellow-300", "bg-yellow-600", "bg-yellow-100", "text-yellow-300"},
	{"bg-green-500", "bg-green-700", "bg-green-100", "text-green-500"},
	{"bg-blue-700", "bg-blue-900", "bg-blue-100", "text-blue-700"},
}

func RowMarginClass(rowIdx int) string {
	if rowIdx < 0 || rowIdx >= len(rowColors) {
		return "mx-auto w-max px-3 py-2"
	}
	return "mx-auto w-max px-3 py-2 " + rowColors[rowIdx].mg
}

func RowInnerClass(rowIdx int) string {
	if rowIdx < 0 || rowIdx >= len(rowColors) {
		return "flex flex-col items-center space-y-0.5 rounded-md p-0.5 md:flex-row md:space-x-0.5 md:space-y-0"
	}
	return "flex flex-col items-center space-y-0.5 rounded-md p-0.5 md:flex-row md:space-x-0.5 md:space-y-0 " + rowColors[rowIdx].fg
}

func CellClass(rowIdx, colIdx int) string {
	if rowIdx < 0 || rowIdx >= len(rowColors) {
		return "flex h-12 w-12 items-center justify-around rounded-md p-2 text-center text-2xl font-bold"
	}
	base := "flex h-12 w-12 items-center justify-around rounded-md p-2 text-center text-2xl font-bold " + rowColors[rowIdx].bg + " " + rowColors[rowIdx].text
	if colIdx == 10 {
		return base + " mt-2 md:ml-2 md:mt-0"
	}
	return base
}

// CellLabel returns the display label for a cell (e.g. "2", "12", "LOCK").
func CellLabel(rowIdx, colIdx int) string {
	if colIdx == 11 {
		return "LOCK"
	}
	if rowIdx <= 1 {
		if colIdx < 10 {
			return fmt.Sprintf("%d", colIdx+2)
		}
		return "12"
	}
	if colIdx < 10 {
		return fmt.Sprintf("%d", 12-colIdx)
	}
	return "2"
}

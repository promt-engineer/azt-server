package engine

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine/utils"
	baseUtils "bitbucket.org/play-workspace/base-slot-server/utils"
)

type Window struct {
	config *baseUtils.Chooser[int, int]

	matrix        [][]int
	baseReels     []map[int]int
	topReels      []top
	columnHeights []int
	scatterCount  int

	changeWindowSize bool
}

type top struct {
	symbol   int
	index    int
	absIndex int
}

func (w *Window) GetWidth() int {
	return windowWidth
}

func (w *Window) GetHeight(col int) int {
	height := w.columnHeights[col]
	if reelWithTop(col) {
		height++
	}

	return height
}

func (w *Window) GetSymbol(colIndex, rowIndex int) int {
	return w.matrix[colIndex][rowIndex]
}

func (w *Window) Compute(stops []int, reels []map[int]int, maxIndexes []int, deletedIndexes []map[int]struct{}) error {
	w.matrix = make([][]int, windowWidth)
	w.baseReels = make([]map[int]int, windowWidth)

	if w.changeWindowSize {
		w.columnHeights = make([]int, windowWidth)
	}

	for index, stop := range stops {
		reel := reels[index]
		reelMaxIndex := maxIndexes[index] + 1

		scatterCount := 0
		appended := 0

		if index < len(stops)-1 {
			if w.changeWindowSize {
				height, err := w.config.Pick()
				if err != nil {
					return err
				}

				w.columnHeights[index] = height
			}

			height := w.columnHeights[index]

			windowLine := make(map[int]int, height)

			for i := 0; appended < height; i++ {
				symbolIndex := mod(stop-i, reelMaxIndex)

				if symbol, exist := reel[symbolIndex]; exist {
					if _, isDeleted := deletedIndexes[index][symbolIndex]; !isDeleted {
						if symbol == scatterSymbol {
							scatterCount++
						}

						// we cannot have two scatter symbols in the same reel, need to replace with symbol 2
						if scatterCount > 1 {
							symbol = 2
						}

						windowLine[symbolIndex] = symbol
						w.matrix[index] = append(w.matrix[index], symbol)

						appended++
					}
				}
			}

			w.baseReels[index] = windowLine

			continue
		}

		// top range is the index from 1 to 4
		topRange := 1

		w.topReels = make([]top, topReelHeight)

		for i := 0; appended < topReelHeight; i++ {
			absIndex := mod(stop-i, reelMaxIndex)

			if symbol, exist := reel[absIndex]; exist {
				if _, isDeleted := deletedIndexes[index][absIndex]; !isDeleted {
					w.matrix[topRange] = append(w.matrix[topRange], symbol)
					w.topReels[appended] = top{symbol: symbol, index: topRange, absIndex: absIndex}

					appended++
					topRange++
				}
			}
		}
	}

	w.changeWindowSize = false

	return nil
}

func (w *Window) CheckWin() []utils.Win[int] {
	return utils.CheckWindow[int](w, wildSymbol, &scatterSymbol)
}

func (w *Window) GetIndexesBySymbol(s int) [][]int {
	results := make([][]int, windowWidth)

	for reelIndex, reel := range w.matrix {
		var find bool

		for symbolIndex, symbol := range reel {
			if symbol == s || (symbol == wildSymbol && s != scatterSymbol) {
				results[reelIndex] = append(results[reelIndex], symbolIndex)

				find = true

				break
			}
		}

		if !find && s != scatterSymbol {
			break
		}
	}

	return results
}

func (w *Window) GetAbsIndexesBySymbol(s int) [][]int {
	results := make([][]int, topReelIndex+1)

	for reelIndex, reel := range w.baseReels {
		var find bool

		var scatterCount int

		for symbolIndex, symbol := range reel {
			if s == symbol {
				results[reelIndex] = append(results[reelIndex], symbolIndex)
				find = true
			}

			if symbol == scatterSymbol {
				scatterCount++
			}

			if scatterCount > 1 && symbol == scatterSymbol && s == 2 {
				results[reelIndex] = append(results[reelIndex], symbolIndex)
				find = true
			}
		}

		if reelIndex >= 1 && reelIndex <= 4 {
			topItem := w.topReels[reelIndex-1]

			if topItem.symbol == s || topItem.symbol == wildSymbol {
				results[topReelIndex] = append(results[topReelIndex], topItem.absIndex)
				find = true
			}
		}

		if !find {
			break
		}
	}

	return results
}

func (w *Window) Matrix() [][]int {
	return w.matrix
}

func (w *Window) GetScatterSymbol() *int {
	return &scatterSymbol
}

func (w *Window) SetScatterQty(qty int) {
	w.scatterCount = qty
}

func (w *Window) GetScatterQty() int {
	return w.scatterCount
}

func (w *Window) GetMultiplierQty() int {
	return 0
}

func (w *Window) SetMultiplierQty(_ int) {}

func (w *Window) GetMultiplierSymbol() *int {
	return nil
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func reelWithTop(index int) bool {
	return index >= 1 && index <= 4
}

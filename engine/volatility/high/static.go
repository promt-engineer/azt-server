package high

import "aztec-pyramids/engine/models"

const (
	mainRTPPercentage = 0.71

	lowRTPBase  = 61.7425
	highRTPBase = 145.9775

	lowRTPAnte  = 40.15373017
	highRTPAnte = 94.96774038

	lowBuyBonus  = 72.92713848
	highBuyBonus = 109.14804601

	lowBuyBonusCorrectionCoef  = 1.32
	highBuyBonusCorrectionCoef = 0.82
)

var (
	baseLowBonusRate  = 336.3379524
	baseHighBonusRate = 137.055768

	baseLow4ScatterCount = 28072.0
	baseLow5ScatterCount = 1623.0
	baseLow6ScatterCount = 37.0

	baseHigh4ScatterCount = 67322.0
	baseHigh5ScatterCount = 5431.0
	baseHigh6ScatterCount = 210.0
	// ------------------------------------------------
	anteLowBonusRate  = 211.26587
	anteHighBonusRate = 83.61259

	anteLow4ScatterCount = 22105289.0
	anteLow5ScatterCount = 1517287.0
	anteLow6ScatterCount = 44285.0

	anteHigh4ScatterCount = 54498750.0
	anteHigh5ScatterCount = 5094300.0
	anteHigh6ScatterCount = 206560.0
)

var configBaseLow = map[int]int{
	2: 20,
	3: 14,
	4: 5,
	5: 2,
	6: 2,
	7: 1,
}

var configBaseHigh = map[int]int{
	2: 5,
	3: 10,
	4: 7,
	5: 2,
	6: 1,
	7: 2,
}

var configAnteLow = map[int]int{
	2: 20,
	3: 14,
	4: 5,
	5: 1,
	6: 1,
	7: 1,
}

var configAnteHigh = map[int]int{
	2: 5,
	3: 10,
	4: 7,
	5: 2,
	6: 1,
	7: 1,
}

var configBonusLow = map[int]int{
	2: 20,
	3: 14,
	4: 5,
	5: 1,
	6: 1,
	7: 1,
}

var configBonusHigh = map[int]int{
	2: 5,
	3: 10,
	4: 7,
	5: 2,
	6: 1,
	7: 1,
}

var rtpTable = map[models.BonusChoice]models.BonusRtp{
	{5, 10, false}: {32.48972245, 85.05909033},
	{10, 5, false}: {43.36345892, 124.57216918},
	{15, 1, false}: {41.22794560, 139.86753553},

	{9, 10, false}: {63.47289469, 171.32703918},
	{14, 5, false}: {68.48061880, 202.68683621},
	{19, 1, false}: {62.83811559, 215.54698777},

	{13, 10, false}: {98.9367, 273.6813},
	{18, 5, false}:  {98.1782, 297.2685},
	{23, 1, false}:  {88.9567, 307.3743},
}

var randomRtpTable = map[int]models.BonusRtp{
	4: {47.55363133, 137.77638482},
	5: {73.55507977, 217.87710674},
	6: {103.9803, 314.0011},
}

// Ordinary base game spin

var availableTopReels = [][]int{topReelSet1, topReelSet2, topReelSet3}
var availableReels = [][][]int{reelSet1, reelSet2, reelSet3}

var topReelSet1 = []int{4, 3, 2, 8, 5, 10, 7, 8, 7, 3, 7, 4, 5, 3, 8, 1, 2, 8, 10, 10, 6, 6, 2, 8, 8, 10, 7, 6, 9, 7, 5, 9, 6, 10, 8, 10, 4, 5, 10, 10, 8, 6, 9, 9, 8, 5, 9, 9, 8, 4, 5, 9, 9, 9, 6, 6, 5, 10, 9, 9, 4, 9, 8, 10, 9, 10, 5, 9, 10, 4, 11, 8, 3, 10, 10, 10, 9, 10, 9, 8, 8, 10, 6, 8, 3, 7, 11, 7, 5, 9}
var reelSet1 = [][]int{
	{7, 8, 4, 8, 12, 8, 8, 9, 6, 8, 2, 3, 6, 6, 8, 4, 6, 8, 9, 9, 9, 6, 10, 9, 9, 12, 6, 9, 9, 9, 10, 10, 6, 7, 1, 10, 6, 7, 10, 6, 10, 9, 10, 10, 9, 12, 7, 9, 5, 7, 10, 9, 2, 2, 9, 8, 10, 8, 7, 10, 5, 10, 8, 10, 4, 10, 10, 7, 8, 9, 2, 3, 8, 10, 3, 10, 12, 3, 4, 8, 5, 5, 4, 5, 10, 10, 5, 8, 5, 9, 8, 8, 3, 5, 8, 8, 7, 7},
	{9, 3, 9, 10, 10, 12, 8, 10, 10, 10, 2, 9, 10, 10, 10, 5, 8, 6, 9, 4, 7, 8, 12, 8, 8, 8, 9, 9, 5, 7, 5, 8, 5, 8, 9, 8, 3, 7, 9, 3, 6, 7, 6, 9, 7, 9, 6, 1, 6, 7, 8, 3, 4, 9, 8, 8, 12, 9, 7, 4, 2, 2, 4, 4, 9, 7, 3, 6, 5, 8, 6, 10, 10, 3, 8, 8, 5, 8, 5, 10, 10, 10, 10, 5, 10, 6, 10},
	{9, 8, 2, 8, 5, 12, 9, 3, 9, 9, 6, 10, 6, 5, 5, 4, 8, 5, 8, 7, 5, 7, 7, 7, 9, 12, 3, 5, 2, 9, 8, 10, 4, 7, 9, 5, 9, 9, 10, 10, 4, 12, 7, 8, 4, 10, 9, 6, 9, 8, 10, 5, 3, 8, 5, 10, 10, 10, 3, 10, 8, 10, 10, 2, 10, 4, 9, 9, 8, 8, 3, 4, 12, 8, 6, 8, 10, 6, 10, 3, 7, 1, 10, 9, 6, 9, 6, 7, 8},
	{7, 9, 7, 5, 6, 2, 8, 5, 3, 9, 10, 10, 6, 10, 5, 2, 10, 4, 8, 8, 6, 3, 10, 8, 10, 6, 8, 1, 10, 5, 6, 3, 6, 5, 4, 12, 5, 2, 8, 7, 10, 8, 8, 5, 10, 10, 9, 9, 9, 8, 9, 3, 6, 7, 8, 9, 12, 8, 6, 7, 9, 9, 5, 8, 4, 4, 2, 10, 4, 9, 8, 9, 8, 10, 9, 4, 7, 10, 3, 9, 3, 8, 9, 5, 12, 10, 8, 8, 9, 7, 9, 8, 10, 7},
	{6, 12, 4, 5, 6, 3, 10, 8, 1, 2, 5, 5, 3, 9, 9, 9, 3, 5, 7, 9, 9, 9, 8, 7, 10, 4, 10, 9, 2, 9, 6, 8, 10, 10, 6, 6, 3, 8, 7, 9, 4, 8, 7, 8, 10, 3, 7, 12, 9, 8, 6, 10, 8, 8, 4, 8, 8, 6, 2, 8, 7, 10, 8, 9, 7, 9, 5, 4, 2, 5, 9, 12, 10, 10, 9, 5, 3, 7, 10, 10, 9, 10, 8, 7, 4, 10, 8, 9},
	{8, 9, 12, 10, 3, 10, 9, 9, 8, 1, 8, 9, 10, 8, 5, 4, 8, 9, 10, 10, 2, 12, 10, 8, 8, 7, 7, 4, 5, 4, 8, 8, 7, 10, 10, 7, 9, 4, 6, 5, 8, 9, 5, 3, 10, 4, 10, 10, 9, 8, 9, 8, 5, 2, 5, 9, 6, 9, 9, 10, 9, 8, 9, 8, 3, 9, 6, 7, 3, 6, 10, 7, 5, 9, 10, 9, 9, 12, 10, 10, 7, 2, 3, 3, 8, 10, 8, 4, 8, 10, 6, 9, 6, 6},
}

// Base game spin with a power bet (the cost of the spin is x1.25 with the same pay table)

var topReelSet2 = []int{4, 3, 2, 8, 5, 10, 7, 8, 7, 3, 7, 4, 5, 3, 8, 1, 2, 8, 10, 10, 6, 6, 2, 8, 8, 10, 7, 6, 9, 7, 5, 9, 6, 10, 8, 10, 4, 5, 10, 10, 8, 6, 9, 9, 8, 5, 9, 9, 8, 4, 5, 9, 9, 9, 6, 6, 5, 10, 9, 9, 4, 9, 8, 10, 9, 10, 5, 9, 10, 4, 11, 8, 3, 10, 10, 10, 9, 10, 9, 8, 8, 10, 6, 8, 3, 7, 11, 7, 5, 9}
var reelSet2 = [][]int{
	{7, 8, 4, 8, 12, 8, 8, 9, 6, 8, 2, 3, 6, 6, 8, 4, 6, 8, 9, 9, 9, 6, 10, 9, 9, 12, 6, 9, 9, 9, 10, 10, 6, 7, 1, 10, 6, 7, 10, 6, 10, 9, 10, 10, 9, 12, 7, 9, 5, 7, 10, 9, 2, 2, 9, 8, 10, 8, 7, 10, 5, 10, 8, 10, 4, 10, 10, 7, 8, 9, 2, 3, 8, 10, 3, 10, 12, 3, 4, 8, 5, 5, 4, 5, 10, 10, 5, 8, 5, 9, 8, 8, 3, 5, 8, 8, 7, 7},
	{9, 3, 9, 10, 10, 12, 8, 10, 10, 10, 2, 9, 10, 10, 10, 5, 8, 6, 9, 4, 7, 8, 12, 8, 8, 8, 9, 9, 5, 7, 5, 8, 5, 8, 9, 8, 3, 7, 9, 3, 6, 7, 6, 9, 7, 9, 6, 1, 6, 7, 8, 3, 4, 9, 8, 8, 12, 9, 7, 4, 2, 2, 4, 4, 9, 7, 3, 6, 5, 8, 6, 10, 10, 3, 8, 8, 5, 12, 8, 5, 10, 10, 10, 10, 5, 10, 6, 10},
	{9, 8, 2, 8, 5, 12, 9, 3, 9, 9, 6, 10, 6, 5, 5, 4, 8, 5, 8, 7, 5, 7, 7, 7, 9, 12, 3, 5, 2, 9, 8, 10, 4, 7, 9, 5, 9, 9, 10, 10, 4, 12, 7, 8, 4, 10, 9, 6, 9, 8, 10, 5, 3, 8, 5, 10, 10, 10, 3, 10, 8, 10, 10, 2, 10, 4, 9, 9, 8, 8, 3, 4, 12, 8, 6, 8, 10, 6, 10, 3, 7, 1, 10, 9, 6, 9, 6, 7, 8},
	{7, 9, 7, 5, 6, 2, 8, 12, 5, 3, 9, 10, 10, 6, 10, 5, 2, 10, 4, 8, 8, 6, 3, 10, 8, 10, 6, 8, 1, 10, 5, 6, 3, 6, 5, 4, 12, 5, 2, 8, 7, 10, 8, 8, 5, 10, 10, 9, 9, 9, 8, 9, 3, 6, 7, 8, 9, 12, 8, 6, 7, 9, 9, 5, 8, 4, 4, 2, 10, 4, 9, 8, 9, 8, 10, 9, 4, 7, 10, 3, 9, 3, 8, 9, 5, 12, 10, 8, 8, 9, 7, 9, 8, 10, 7},
	{6, 12, 4, 5, 6, 3, 10, 8, 1, 2, 5, 5, 3, 9, 9, 9, 3, 5, 7, 9, 9, 9, 8, 7, 10, 4, 12, 10, 9, 2, 9, 6, 8, 10, 10, 6, 6, 3, 8, 7, 9, 4, 8, 7, 8, 10, 3, 7, 12, 9, 8, 6, 10, 8, 8, 4, 8, 8, 6, 2, 8, 7, 10, 8, 9, 7, 9, 5, 4, 2, 5, 9, 12, 10, 10, 9, 5, 3, 7, 10, 10, 9, 10, 8, 7, 4, 10, 8, 9},
	{8, 9, 12, 10, 3, 10, 9, 9, 8, 1, 8, 9, 10, 8, 5, 4, 8, 9, 10, 10, 2, 12, 10, 8, 8, 7, 7, 4, 5, 4, 8, 8, 7, 10, 10, 7, 9, 4, 6, 5, 8, 9, 5, 3, 10, 4, 10, 10, 9, 12, 8, 9, 8, 5, 2, 5, 9, 6, 9, 9, 10, 9, 8, 9, 8, 3, 9, 6, 7, 3, 6, 10, 7, 5, 9, 10, 9, 9, 12, 10, 10, 7, 2, 3, 3, 8, 10, 8, 4, 8, 10, 6, 9, 6, 6},
}

// Free game

var topReelSet3 = []int{4, 3, 2, 8, 5, 10, 7, 8, 7, 3, 7, 4, 5, 3, 8, 1, 2, 8, 10, 10, 6, 6, 2, 8, 8, 10, 7, 6, 9, 7, 5, 9, 6, 10, 8, 10, 4, 5, 10, 10, 8, 6, 9, 9, 8, 5, 9, 9, 8, 4, 5, 9, 9, 9, 6, 6, 5, 10, 9, 9, 4, 9, 8, 10, 9, 10, 5, 9, 10, 4, 11, 8, 3, 10, 10, 10, 9, 10, 9, 8, 8, 10, 6, 8, 3, 7, 11, 7, 5, 9}
var reelSet3 = [][]int{
	{7, 8, 4, 8, 12, 8, 8, 9, 6, 8, 2, 3, 6, 6, 8, 4, 6, 8, 9, 9, 9, 6, 10, 9, 9, 6, 9, 9, 9, 10, 10, 6, 7, 1, 10, 6, 7, 10, 6, 10, 9, 10, 10, 9, 12, 7, 9, 5, 7, 10, 9, 2, 2, 9, 8, 10, 8, 7, 10, 5, 10, 8, 10, 4, 10, 10, 7, 8, 9, 2, 3, 8, 10, 3, 10, 3, 4, 8, 5, 5, 4, 5, 10, 10, 5, 8, 5, 9, 8, 8, 3, 5, 8, 8, 7, 7},
	{9, 3, 9, 10, 10, 12, 8, 10, 10, 10, 2, 9, 10, 10, 10, 5, 8, 6, 9, 4, 7, 8, 8, 8, 8, 9, 9, 5, 7, 5, 8, 5, 8, 9, 8, 3, 7, 9, 3, 6, 7, 6, 9, 7, 9, 6, 1, 6, 7, 8, 3, 4, 9, 8, 8, 12, 9, 7, 4, 2, 2, 4, 4, 9, 7, 3, 6, 5, 8, 6, 10, 10, 3, 8, 8, 5, 8, 5, 10, 10, 10, 10, 5, 10, 6, 10},
	{9, 8, 2, 8, 5, 12, 9, 3, 9, 9, 6, 10, 6, 5, 5, 4, 8, 5, 8, 7, 5, 7, 7, 7, 9, 3, 5, 2, 9, 8, 10, 4, 7, 9, 5, 9, 9, 10, 10, 4, 12, 7, 8, 4, 10, 9, 6, 9, 8, 10, 5, 3, 8, 5, 10, 10, 10, 3, 10, 8, 10, 10, 2, 10, 4, 9, 9, 8, 8, 3, 4, 8, 6, 8, 10, 6, 10, 3, 7, 1, 10, 9, 6, 9, 6, 7, 8},
	{7, 9, 7, 5, 6, 2, 8, 5, 3, 9, 10, 10, 6, 10, 5, 2, 10, 4, 8, 8, 6, 3, 10, 8, 10, 6, 8, 1, 10, 5, 6, 3, 6, 5, 4, 12, 5, 2, 8, 7, 10, 8, 8, 5, 10, 10, 9, 9, 9, 8, 9, 3, 6, 7, 8, 9, 8, 6, 7, 9, 9, 5, 8, 4, 4, 2, 10, 4, 9, 8, 9, 8, 10, 9, 4, 7, 10, 3, 9, 3, 8, 9, 5, 12, 10, 8, 8, 9, 7, 9, 8, 10, 7},
	{6, 12, 4, 5, 6, 3, 10, 8, 1, 2, 5, 5, 3, 9, 9, 9, 3, 5, 7, 9, 9, 9, 8, 7, 10, 4, 10, 9, 2, 9, 6, 8, 10, 10, 6, 6, 3, 8, 7, 9, 4, 8, 7, 8, 10, 3, 7, 12, 9, 8, 6, 10, 8, 8, 4, 8, 8, 6, 2, 8, 7, 10, 8, 9, 7, 9, 5, 4, 2, 5, 9, 10, 10, 9, 5, 3, 7, 10, 10, 9, 10, 8, 7, 4, 10, 8, 9},
	{8, 9, 10, 3, 10, 9, 9, 8, 1, 8, 9, 10, 8, 5, 4, 8, 9, 10, 10, 2, 12, 10, 8, 8, 7, 7, 4, 5, 4, 8, 8, 7, 10, 10, 7, 9, 4, 6, 5, 8, 9, 5, 3, 10, 4, 10, 10, 9, 8, 9, 8, 5, 2, 5, 9, 6, 9, 9, 10, 9, 8, 9, 8, 3, 9, 6, 7, 3, 6, 10, 7, 5, 9, 10, 9, 9, 12, 10, 10, 7, 2, 3, 3, 8, 10, 8, 4, 8, 10, 6, 9, 6, 6},
}

var buyBonusStops = [][]int{
	{4, 5, 5, 35, 10, 10, 0},   // *, *, *, *, -, -, -
	{5, 15, 5, 34, 47, 20, 0},  // *, -, *, -, *, *, -
	{39, 56, 6, 28, 47, 20, 0}, // -, *, *, -, *, *, -
	{4, 56, 85, 28, 47, 20, 0}, // *, *, -, -, *, *, -
	{4, 56, 85, 36, 35, 20, 0}, // *, *, -, *, -, *, -
}

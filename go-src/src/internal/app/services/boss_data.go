package services

import "nro/src/internal/core/domain"

// BossData holds the static data for all bosses.
var BossData = map[int]*domain.BossTemplate{
	// Broly
	1: {
		ID:          1,
		Name:        "Broly",
		Gender:      2,
		Outfit:      []int16{294, 295, 296, -1, -1, -1},
		Damage:      1000,
		HP:          []int64{100000, 200000, 500000},
		MapJoin:     []int{5},
		SkillTemp:   [][]int{{0, 1, 1000}, {1, 1, 2000}},
		TextS:       []string{"Ta là Broly", "Ta sẽ tiêu diệt tất cả"},
		TextM:       []string{"Đỡ đòn này", "Yaaaa!"},
		TextE:       []string{"Ta sẽ quay lại", "Không thể nào..."},
		SecondsRest: 10,
		TypeAppear:  0,
		AIType:      "Broly",
	},
	// Android 13
	13: {
		ID:          13,
		Name:        "Android 13",
		Gender:      1,
		Outfit:      []int16{300, 301, 302, -1, -1, -1},
		Damage:      1200,
		HP:          []int64{200000},
		MapJoin:     []int{84},
		SkillTemp:   [][]int{{2, 1, 3000}},
		TextS:       []string{"Ta là Android 13"},
		TextM:       []string{"Chết đi!"},
		TextE:       []string{"Hự..."},
		SecondsRest: 15,
		TypeAppear:  0,
		AIType:      "Android",
		AndroidID:   13,
	},
}

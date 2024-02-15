package models

type HeroModel struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	LocalizedName string `json:"localized_name"`
	Image         string `json:"img"`
	HeraldPicks   int    `json:"1_pick"`
	HeraldWins    int    `json:"1_win"`
	GuardianPicks int    `json:"2_pick"`
	GuardianWins  int    `json:"2_win"`
	CrusaderPicks int    `json:"3_pick"`
	CrusaderWins  int    `json:"3_win"`
	ArchonPicks   int    `json:"4_pick"`
	ArchonWins    int    `json:"4_win"`
	LegendPicks   int    `json:"5_pick"`
	LegendWins    int    `json:"5_win"`
	AncientPicks  int    `json:"6_pick"`
	AncientWins   int    `json:"6_win"`
	DivinePicks   int    `json:"7_pick"`
	DivineWins    int    `json:"7_win"`
	ImmortalPicks int    `json:"8_pick"`
	ImmortalWins  int    `json:"8_win"`
}

func (H *HeroModel) GetWinrates() []float64 {
	ans := make([]float64, 0)
	ans = append(ans, float64(H.HeraldWins)/float64(H.HeraldPicks))
	ans = append(ans, float64(H.GuardianWins)/float64(H.GuardianPicks))
	ans = append(ans, float64(H.CrusaderWins)/float64(H.CrusaderPicks))
	ans = append(ans, float64(H.ArchonWins)/float64(H.ArchonPicks))
	ans = append(ans, float64(H.LegendWins)/float64(H.LegendPicks))
	ans = append(ans, float64(H.AncientWins)/float64(H.AncientPicks))
	ans = append(ans, float64(H.DivineWins)/float64(H.DivinePicks))
	ans = append(ans, float64(H.ImmortalWins)/float64(H.ImmortalPicks))

	return ans
}

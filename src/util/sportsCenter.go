package util

type MatchInfoQuery struct {
	Cp_id                 string `json:"cp_id"`
	Match_name            string `json:"match_name"`
	Match_level           string `json:"match_level"`
	Match_start_time_from string `json:"match_start_time_from"`
	Match_start_time_to   string `json:"match_start_time_to"`
	Match_end_time_from   string `json:"match_end_time_from"`
	Match_end_time_to     string `json:"match_end_time_to"`
	Entry_fee             string `json:"entry_fee"`
	Bonus_per_match       string `json:"bonus_per_match"`
	Page                  int    `json:"page"`
	Page_size             int    `json:"page_size"`
}

type SonMatchQuery struct {
	Cp_id     string `json:"cp_id"`
	Match_id  string `json:"match_id"`
	Status    string `json:"status"`
	Page      int    `json:"page"`
	Page_size int    `json:"page_size"`
}

type PlayerCashoutReq struct {
	Cp_id            string `json:"cp_id"`
	Player_id        string `json:"player_id"`
	Player_id_number string `json:"player_id_number"`
}

type AwardResultReq struct {
	Cp_id     string `json:"cp_id"`
	Match_id  string `json:"match_id"`
	Page      int    `json:"page"`
	Page_size int    `json:"page_size"`
}

type PlayerWalletInfoQuery struct {
	Cp_id     string `json:"cp_id"`
	Player_id string `json:"player_id"`
	Page      int    `json:"page"`
	Page_size int    `json:"page_size"`
}

type PlayerWalletBalanceQuery struct {
	Cp_id     string `json:"cp_id"`
	Player_id string `json:"player_id"`
}

type PlayerWalletListQuery struct {
	Cp_id     string `json:"cp_id"`
	Player_id string `json:"player_id"`
	Page      int    `json:"page"`
	Page_size int    `json:"page_size"`
}

type PlayerMasterScoreQuery struct {
	Cp_id            string `json:"cp_id"`
	Player_id_number string `json:"player_id_number"`
}

type PlayerWalletTransaction struct {
	Cp_id      string  `json:"cp_id"`
	Player_id  string  `json:"player_id"`
	Order_id   string  `json:"order_id"`
	Amount     float64 `json:"amount"`
	Notes      string  `json:"notes"`
	Notify_url string  `json:"notify_url"`
}

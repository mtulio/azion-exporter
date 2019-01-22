package collector

type mIndexing struct {
	Products struct {
		Num1441740010 struct {
			Requests struct {
				Total [][]interface{} `json:"total"`
			} `json:"requests"`
		} `json:"1441740010"`
	} `json:"products"`
}

type mIndexing5xx struct {
	Products struct {
		Num1441740010 struct {
			StatusCode struct {
				Code [][]interface{} `json:"5xx"`
			} `json:"status_code"`
		} `json:"1441740010"`
	} `json:"products"`
}

type mIndexing500 struct {
	Products struct {
		Num1441740010 struct {
			StatusCode struct {
				Code [][]interface{} `json:"500"`
			} `json:"status_code"`
		} `json:"1441740010"`
	} `json:"products"`
}

type mIndexing502 struct {
	Products struct {
		Num1441740010 struct {
			StatusCode struct {
				Code [][]interface{} `json:"502"`
			} `json:"status_code"`
		} `json:"1441740010"`
	} `json:"products"`
}

type mIndexing503 struct {
	Products struct {
		Num1441740010 struct {
			StatusCode struct {
				Code [][]interface{} `json:"503"`
			} `json:"status_code"`
		} `json:"1441740010"`
	} `json:"products"`
}

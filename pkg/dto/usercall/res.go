package usercall

type OnInputFuzzyResp struct {
	Values []string `json:"values"`
}

type OnInputValidateResp struct {
	Msg string `json:"msg"`
}

type OnTableDeleteRowsResp struct {
}

type OnTableUpdateRowResp struct {
}

type OnTableSearchResp struct {
}

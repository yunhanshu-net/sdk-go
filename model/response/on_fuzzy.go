package response

type OnInputFuzzy struct {
	Values []string `json:"values"`
}

type OnInputValidate struct {
	Msg string `json:"msg"`
}

type OnTableDeleteRows struct {
}

type OnTableUpdateRow struct {
}

type OnTableSearch struct {
}

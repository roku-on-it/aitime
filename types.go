package main

type (
	ChatReq struct {
		Model    string       `json:"model"`
		Stream   bool         `json:"stream"`
		Messages []RequestMsg `json:"messages"`
	}

	RequestMsg struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)

type (
	ChatResponse struct {
		Message ResponseMsg `json:"message"`
	}

	ResponseMsg struct {
		Content string `json:"content"`
	}
)

type (
	WatchNextResponse struct {
		Episode int  `json:"number"`
		Season  int  `json:"season_number"`
		Show    Show `json:"show"`
	}

	Show struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}
)

type HandlePromptReqBody struct {
	Prompt string `json:"prompt"`
}

type FavoriteShowsResponse struct {
	Shows []Show `json:"shows"`
}

type (
	EpisodesResponse struct {
		Data []Data `json:"data"`
	}

	Data struct {
		ID     int `json:"id"`
		Number int `json:"number"`
		Season `json:"season"`
	}

	Season struct {
		Number int `json:"number"`
	}
)

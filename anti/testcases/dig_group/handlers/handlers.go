package handlers

import dig "github.com/3timeslazy/anti-dig"

type Handler struct{}

type ListV1Out struct {
	dig.Out
	Handlers []Handler `group:"handlers_v1,flatten"`
}

func NewListV1() ListV1Out {
	return ListV1Out{}
}

type ListV2Out struct {
	dig.Out
	Handlers []Handler `group:"handlers_v2,flatten"`
}

func NewListV2() ListV2Out {
	return ListV2Out{}
}

type HandlerV1 struct {
	dig.Out
	Handler Handler `group:"handlers_v1"`
}

func NewV1() HandlerV1 {
	return HandlerV1{}
}

func NewV2() Handler {
	return Handler{}
}

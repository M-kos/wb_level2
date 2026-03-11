package dto

type EventsResponse struct {
	Result []*EventDto `json:"result"`
}

type EventResponse struct {
	Result *EventDto `json:"result"`
}

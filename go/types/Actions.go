package types

type Action string

const (
	Post   Action = "POST"
	Put    Action = "PUT"
	Patch  Action = "PATCH"
	Delete Action = "DELETE"
	Get    Action = "GET"
)

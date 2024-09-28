package internal

const (
	HealthCommand = "health"
	SetCommand    = "set"
	GetCommand    = "get"
	RemoveCommand = "remove"
)

type Command interface {
}

type CommandHealth struct {
}

type CommandSet struct {
	key   string
	value string
}

type CommandGet struct {
	key string
}

type CommandRemove struct {
	key string
}

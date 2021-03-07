package command

import "github.com/HotPotatoC/kvstore/pkg/hashtable"

type CommandOp int

const (
	SET CommandOp = iota
	GET
	DEL
	LIST
	KEYS
)

func (c CommandOp) String() string {
	return [...]string{"set", "get", "del", "list", "keys"}[c]
}

type Command interface {
	String() string
	Execute(args []string) []byte
}

func New(db *hashtable.HashTable, cmd CommandOp) Command {
	var command Command
	switch cmd {
	case SET:
		command = makeSetCommand(db)
		break
	case GET:
		command = makeGetCommand(db)
		break
	case DEL:
		command = makeDelCommand(db)
		break
	case LIST:
		command = makeListCommand(db)
		break
	case KEYS:
		command = makeKeysCommand(db)
		break
	default:
		command = nil
		break
	}

	return command
}

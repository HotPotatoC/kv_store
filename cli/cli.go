package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"

	"github.com/HotPotatoC/kvstore/command"
	"github.com/HotPotatoC/kvstore/packet"
	"github.com/HotPotatoC/kvstore/pkg/comm"
	"github.com/HotPotatoC/kvstore/pkg/utils"
)

// CLI represents the cli client
type CLI struct {
	comm   *comm.Comm
	reader *bufio.Reader
}

// New creates a new CLI client
func New(addr string) *CLI {
	comm, err := comm.New(addr)
	if err != nil {
		log.Fatal(err)
	}

	return &CLI{
		comm:   comm,
		reader: bufio.NewReader(os.Stdin),
	}
}

// Start runs the CLI client
func (c *CLI) Start() {
	go func() {
	start:
		for {
			fmt.Printf("%s> ", c.comm.Connection().RemoteAddr().String())

			input, err := c.reader.ReadBytes('\n')
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			preprocessed, err := c.preprocess(input)
			if err != nil {
				log.Println(err)
				continue start
			}

			err = c.comm.Send(preprocessed.Bytes())
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			msg, _, err := c.comm.Read()
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			fmt.Print(string(msg))
		}
	}()

	<-utils.WaitForSignals(os.Interrupt, syscall.SIGTERM)
	c.comm.Connection().Close()
	os.Exit(0)
}

func (c *CLI) preprocess(data []byte) (*bytes.Buffer, error) {
	var packet *packet.Packet
	var err error

	rawCmd := bytes.Split(data, []byte(" "))[0]
	cmd := bytes.ToLower(
		bytes.TrimSpace(rawCmd))
	args := bytes.TrimSpace(
		bytes.TrimPrefix(data, rawCmd))

	switch string(cmd) {
	case command.SET.String():
		if packet, err = c.set(args); err != nil {
			return nil, err
		}
	case command.GET.String():
		if packet, err = c.get(args); err != nil {
			return nil, err
		}
	case command.DEL.String():
		if packet, err = c.del(args); err != nil {
			return nil, err
		}
	case command.LIST.String():
		if packet, err = c.list(); err != nil {
			return nil, err
		}
	case command.KEYS.String():
		if packet, err = c.keys(); err != nil {
			return nil, err
		}
	case command.FLUSH.String():
		if packet, err = c.flush(); err != nil {
			return nil, err
		}
	case command.INFO.String():
		if packet, err = c.info(); err != nil {
			return nil, err
		}
	case "exit":
		c.comm.Connection().Close()
		os.Exit(0)
	default:
		return nil, command.ErrCommandDoesNotExist
	}

	buffer, err := packet.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed processing input: %v", err)
	}

	return buffer, nil
}

func (c *CLI) set(args []byte) (*packet.Packet, error) {
	if len(bytes.Split(args, []byte(" "))) < 2 {
		return nil, command.ErrMissingKeyValueArg
	}
	return packet.NewPacket(command.SET, args), nil
}

func (c *CLI) get(args []byte) (*packet.Packet, error) {
	if bytes.Equal(args, []byte("")) {
		return nil, command.ErrMissingKeyArg
	}
	return packet.NewPacket(command.GET, args), nil
}

func (c *CLI) del(args []byte) (*packet.Packet, error) {
	if bytes.Equal(args, []byte("")) {
		return nil, command.ErrMissingKeyArg
	}
	return packet.NewPacket(command.DEL, args), nil
}

func (c *CLI) list() (*packet.Packet, error) {
	return packet.NewPacket(command.LIST, []byte("")), nil
}

func (c *CLI) keys() (*packet.Packet, error) {
	return packet.NewPacket(command.KEYS, []byte("")), nil
}

func (c *CLI) flush() (*packet.Packet, error) {
	return packet.NewPacket(command.FLUSH, []byte("")), nil
}

func (c *CLI) info() (*packet.Packet, error) {
	return packet.NewPacket(command.INFO, []byte("")), nil
}
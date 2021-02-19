package client

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/HotPotatoC/kvstore/internal/server"
	"github.com/HotPotatoC/kvstore/pkg/comm"
	"github.com/HotPotatoC/kvstore/pkg/logger"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

type client struct {
	comm   *comm.Comm
	reader *bufio.Reader
}

func init() {
	log = logger.NewLogger()
}

func New(addr string) *client {
	comm, err := comm.New(addr)
	if err != nil {
		log.Fatal(server.ErrConnectionRefused)
	}

	return &client{
		comm:   comm,
		reader: bufio.NewReader(os.Stdin),
	}
}

func (c *client) StartCLI() {
	for {
		fmt.Printf("%s> ", c.comm.Connection().RemoteAddr().String())

		input, err := c.reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		err = c.comm.Send(input)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		msg, err := c.comm.Read()
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		fmt.Println(string(msg))
	}
}

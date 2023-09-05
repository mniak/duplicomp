package duplicomp

import (
	"io"

	"google.golang.org/protobuf/proto"
)

type Pipe struct {
	In  Stream
	Out Stream
}

func (p *Pipe) Run() error {
	var err error
	for {
		var msg proto.Message
		msg, err = p.In.Receive()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		err = p.Out.Send(msg)
		if err != nil {
			break
		}
	}
	return err
}

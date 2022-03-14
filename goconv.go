package goconv

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type request struct {
	filename string
	w        io.Writer
	errChan  chan error
}

// Service .
type Service struct {
	queue chan request
}

// NewService .
func NewService() *Service {
	s := &Service{
		queue: make(chan request),
	}
	go s.run()
	return s
}

func (s *Service) run() {
	// run a loop to consume the queue
	for {
		select {
		case data := <-s.queue:
			cmd := exec.Command("unoconv", "-f", "pdf", "--stdout", data.filename)
			cmd.Stdout = data.w
			err := cmd.Run()
			if err != nil {
				data.errChan <- err
			} else {
				data.errChan <- nil
			}
		}
	}
}

// Convert submit a conversion request to the service
func (s *Service) Convert(r io.Reader, w io.Writer) error {
	tempfile, err := ioutil.TempFile(os.TempDir(), "goconv*")
	if err != nil {
		return err
	}
	defer os.Remove(tempfile.Name())
	io.Copy(tempfile, r)
	tempfile.Close()
	errchan := make(chan error)
	req := request{
		tempfile.Name(),
		w,
		errchan,
	}
	s.queue <- req
	return <-errchan
}

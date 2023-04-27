package socketwaiter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

type socketwaiter struct {
	wd uint32
	fd uint32
}

func WaitForSocketFile(path string, filename string, timeout time.Duration) <-chan error {
	outchan := make(chan error, 1)
	internalchan := make(chan error, 1)
	closechan := make(chan int, 1)
	fd, err := unix.InotifyInit1(unix.IN_CLOEXEC | unix.IN_NONBLOCK)
	if fd == -1 {
		outchan <- err
		return outchan
	}
	inotifyfile := os.NewFile(uintptr(fd), "")
	go func() {
		<-closechan
		err := inotifyfile.Close()
		if err != nil {
			fmt.Println("Error closing inotify fd", err)
		}
	}()
	// add watch
	wd, err := unix.InotifyAddWatch(fd, path, unix.IN_CREATE)
	if wd == -1 {
		outchan <- err
		closechan <- 1
		return outchan
	}
	// read output
	go readEvents(filename, inotifyfile, internalchan)
	go func() {
		select {
		case e := <-internalchan:
			outchan <- e
			closechan <- 1
		case <-time.After(timeout):
			outchan <- fmt.Errorf("Socket file wait timed out")
			closechan <- 1
		}
	}()
	return outchan
}

func readEvents(filename string, file *os.File, outchan chan<- error) {
	var buf [unix.SizeofInotifyEvent * 4096]byte
	for {
		n, err := file.Read(buf[:])
		// handle read error
		switch {
		case errors.Unwrap(err) == os.ErrClosed:
			outchan <- err
		case err != nil:
			fmt.Println("Error reading inotifyfile", err)
			continue
		}
		// handle too short output
		if n < unix.SizeofInotifyEvent {
			var err error
			if n == 0 {
				// If EOF is received. This should really never happen.
				err = io.EOF
			} else {
				// Read was too short.
				err = errors.New("notify: short read in readEvents()")
			}
			fmt.Println("Error occurred socket wait", err)
			continue
		}

		var offset uint32
		for offset <= uint32(n-unix.SizeofInotifyEvent) {
			var (
				// Point "raw" to the event in the buffer
				raw     = (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))
				nameLen = uint32(raw.Len)
			)
			if nameLen <= 0 {
				fmt.Println("name too short")
				offset += unix.SizeofInotifyEvent
				continue
			}
			// Point "bytes" at the first byte of the filename
			bytes := (*[unix.PathMax]byte)(unsafe.Pointer(&buf[offset+unix.SizeofInotifyEvent]))[:nameLen:nameLen]
			// The filename is padded with NULL bytes. TrimRight() gets rid of those.

			name := strings.TrimRight(string(bytes[0:nameLen]), "\000")
			if name == filename {
				outchan <- nil
				return
			}
			offset += unix.SizeofInotifyEvent + nameLen
		}

	}
}

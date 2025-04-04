package times

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/zimeg/emporia-time/internal/errors"
)

// bufferWriter contains two writers to write to and a bounds for toggles
type bufferWriter struct {
	bounds string
	buff   *bytes.Buffer
	std    io.Writer
	stored bool
}

// Write writes to the stream until bounds is reached then writes to buffer
func (bw *bufferWriter) Write(p []byte) (int, error) {
	output := string(p)
	if strings.HasPrefix(output, bw.bounds) {
		bw.stored = true
		output = strings.TrimPrefix(output, bw.bounds)
	}
	if bw.stored {
		n, err := bw.buff.Write([]byte(output))
		if err != nil {
			return n, errors.Wrap(errors.ErrWriteBuffer, err)
		}
		return len(p), nil
	} else {
		n, err := bw.std.Write(p)
		if err != nil {
			return n, errors.Wrap(errors.ErrWriteOutput, err)
		}
		return n, nil
	}
}

// timerCommand forms the command struct for a timer to be parsed
func timerCommand(command []string, stderr bufferWriter) *exec.Cmd {
	timer, err := exec.LookPath("time")
	if err != nil {
		timer = "time"
	} else {
		timer, err = filepath.Abs(timer)
		if err != nil {
			timer = "time"
		}
	}
	if strings.HasPrefix(command[0], "./") {
		command = append([]string{"source"}, command...)
	}
	timeShell := []string{
		"(" + strings.Join(command, " ") + ")",
		";",
		"EMPORIA_TIME_EXIT_CODE_STATUS=$?",
		";",
		"1>&2",
		"echo",
		stderr.bounds,
		"1>&2",
		"echo",
		"code $EMPORIA_TIME_EXIT_CODE_STATUS",
	}
	timeArgs := []string{
		"-p",
		"sh",
		"-c",
		strings.Join(timeShell, " "),
	}
	cmd := exec.Command(timer, timeArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	return cmd
}

// makeBounds creates a random string to denote the end of command output
func makeBounds() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const size = 64
	var bounds strings.Builder
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for range size {
		bounds.WriteByte(charset[random.Intn(len(charset))])
	}
	bounds.WriteByte('\n')
	return bounds.String()
}

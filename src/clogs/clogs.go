package clogs

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/xiejw/lunar/base/errors"
)

type FileItem struct {
	Path     string // relative to domain
	Size     int64
	Checksum string
}

type CmdKind int

const (
	CmdNew CmdKind = iota
	CmdDel
)

type CmdLog struct {
	Kind      CmdKind
	FileItem  FileItem
	Timestamp int64
}

type CmdLogs struct {
	Cmds      []*CmdLog
	VersionID int // same as len(Cmds)
}

// Serialize a CmdLog to a string line (no newline).
func (clog *CmdLog) ToOneLine() string {
	// re: file size. 12v is larger than 300G, which should be enough, which should be enough
	// re: epoch. in 200 years later, epoch will need one more char.
	switch clog.Kind {
	case CmdNew:
		return fmt.Sprintf(
			"+\t%12v\t%v\t%v\t%v",
			clog.FileItem.Size,
			clog.FileItem.Checksum,
			clog.Timestamp,
			clog.FileItem.Path)
	case CmdDel:
		return fmt.Sprintf(
			"-\t%12v\t%v\t%v\t%v",
			clog.FileItem.Size,
			clog.FileItem.Checksum,
			clog.Timestamp,
			clog.FileItem.Path)
	default:
		panic("unsupported CmdKind.")
	}
}

// Deserialize lines from r to CmdLogs.
func FromLines(r io.Reader) (*CmdLogs, error) {
	buf_r := bufio.NewReader(r)

	cmds := make([]*CmdLog, 0)

	stop := false
	for !stop {
		line, err := buf_r.ReadString('\n')

		// phase 1: break lines and handle EOF.
		switch err {
		case nil:
			// continue.
		case io.EOF:
			stop = true
			// still need to process line.
		default:
			return nil, errors.From(err).EmitNote("unexpected error during read one line from reader.")
		}

		// phase 2: remove while spaces.
		line = strings.Trim(line, " \n\t")
		if len(line) == 0 {
			continue
		}

		// phase 3: parse.
		parts := strings.SplitN(line, "\t", 5)
		if len(parts) != 5 {
			return nil, errors.New("failed to parse leading char (incorrect number parts): %v", parts)
		}

		cl := &CmdLog{}
		// CmdKind Size Checksum Timestamp Path

		switch parts[0] {
		case "+":
			cl.Kind = CmdNew
		case "-":
			cl.Kind = CmdDel
		default:
			return nil, errors.New("failed to parse leading char: %v", line)
		}

		cl.FileItem.Size, err = strconv.ParseInt(strings.Trim(parts[1], " "), 10, 64)
		if err != nil {
			return nil, errors.WrapNote(err, "failed to parse size: %v", line)
		}

		cl.FileItem.Checksum = parts[2]

		cl.Timestamp, err = strconv.ParseInt(strings.Trim(parts[3], " "), 10, 64)
		if err != nil {
			return nil, errors.WrapNote(err, "failed to parse timestamp: %v", line)
		}

		cl.FileItem.Path = parts[4]

		cmds = append(cmds, cl)
	}

	clgs := &CmdLogs{
		Cmds:      cmds,
		VersionID: len(cmds),
	}
	return clgs, nil
}

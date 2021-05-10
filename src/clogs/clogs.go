package clogs

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/xiejw/fsx/src/errors"
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

// serialization
func ToOneLine(clog *CmdLog) string {
	// re: file size. 12v is larger than 300G, which should be enough, which should be enough
	// re: epoch. in 200 years later, epoch will need one more char.
	switch clog.Kind {
	case CmdNew:
		return fmt.Sprintf(
			"+ %12v %v %v %v",
			clog.FileItem.Size,
			clog.FileItem.Checksum,
			clog.Timestamp,
			clog.FileItem.Path)
	case CmdDel:
		return fmt.Sprintf(
			"- %12v %v %v %v",
			clog.FileItem.Size,
			clog.FileItem.Checksum,
			clog.Timestamp,
			clog.FileItem.Path)
	default:
		panic("unsupported CmdKind.")
	}
}

func FromLines(r io.Reader) (*CmdLogs, error) {
	buf_r := bufio.NewReader(r)

	items := make([]*CmdLog, 0)

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
		switch line[0] {
		case '+':
			items = append(items, &CmdLog{
				Kind: CmdNew,
			})
		case '-':
			items = append(items, &CmdLog{
				Kind: CmdDel,
			})
		default:
			return nil, errors.New("failed to parse leading char: %v", line)
		}
	}

	clgs := &CmdLogs{
		Cmds:      items,
		VersionID: len(items),
	}
	return clgs, nil
}

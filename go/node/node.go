package node

import (
	"fmt"

	"github.com/xiejw/fsx/go/cmdlog"
)

type Node struct {
	Name        string // Must be uinque in the cluster.
	NextVersion uint64 // Points to next verison.
	IsMaster    bool
	CmdLogs     []cmdlog.CmdLog // Ordered CmdLog
}

// Performs a sanity check on the state.
func (state *Node) Check() error {
	if state.Name == "" {
		return fmt.Errorf("name cannot be empty.")
	}
	if state.NextVersion != uint64(len(state.CmdLogs)) {
		return fmt.Errorf(
			"NextVersion field check error. got: %v, expected: %v.",
			state.NextVersion, len(state.CmdLogs))
	}
	for i, cmdLog := range state.CmdLogs {
		if cmdLog.Version != uint64(i) {
			return fmt.Errorf("At %v, the CmdLog version is not right: %v.",
				i, cmdLog.Version)
		}
	}

	return nil
}

// What's the value to marshal?
//
// func (state *Node) Marshal(w io.Writer) error {
// 	c, err := json.MarshalIndent(state, "", "  ")
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = w.Write(c)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func Unmarshal(data []byte) (*Node, error) {
// 	state := &Node{}
// 	err := json.Unmarshal(data, state)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return state, nil
// }
//
// func (state *Node) String() string {
// 	var buf bytes.Buffer
// 	err := state.Marshal(&buf)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return buf.String()
// }

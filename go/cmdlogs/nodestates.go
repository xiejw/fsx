package cmdlogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type NodeState struct {
	Name        string   // Must be uinque in the cluster.
	NextVersion uint64   // Points to next verison.
	CmdLogs     []CmdLog // Ordered CmdLog
}

// Performs a sanity check on the state.
func (state *NodeState) Check() error {
	if state.Name == "" {
		return fmt.Errorf("name cannot be empty.")
	}
	if state.NextVersion != uint64(len(state.CmdLogs)) {
		return fmt.Errorf(
			"NextVersion field check error. got: %v, expected: %v",
			state.NextVersion, len(state.CmdLogs))
	}

	return nil
}

func (state *NodeState) Marshal(w io.Writer) error {
	c, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(c)
	if err != nil {
		return err
	}
	return nil
}

func Unmarshal(data []byte) (*NodeState, error) {
	state := &NodeState{}
	err := json.Unmarshal(data, state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (state *NodeState) String() string {
	var buf bytes.Buffer
	err := state.Marshal(&buf)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

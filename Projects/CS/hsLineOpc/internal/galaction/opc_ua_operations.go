package galaction

import (
	"github.com/gopcua/opcua/ua"
)

func MustParseNodeID(s string) *ua.NodeID {
	return ua.MustParseNodeID(s)
}

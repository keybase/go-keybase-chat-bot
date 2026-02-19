// Code generated to Go types using avdl-compiler v1.4.10 (https://github.com/keybase/node-avdl-compiler). DO NOT EDIT.
//   Input file: ../client/protocol/avdl/keybase1/os.avdl

package keybase1

import (
	"fmt"
)

type RuntimeGroup int

const (
	RuntimeGroup_UNKNOWN     RuntimeGroup = 0
	RuntimeGroup_LINUXLIKE   RuntimeGroup = 1
	RuntimeGroup_DARWINLIKE  RuntimeGroup = 2
	RuntimeGroup_WINDOWSLIKE RuntimeGroup = 3
)

func (o RuntimeGroup) DeepCopy() RuntimeGroup { return o }

var RuntimeGroupMap = map[string]RuntimeGroup{
	"UNKNOWN":     0,
	"LINUXLIKE":   1,
	"DARWINLIKE":  2,
	"WINDOWSLIKE": 3,
}

var RuntimeGroupRevMap = map[RuntimeGroup]string{
	0: "UNKNOWN",
	1: "LINUXLIKE",
	2: "DARWINLIKE",
	3: "WINDOWSLIKE",
}

func (o RuntimeGroup) String() string {
	if v, ok := RuntimeGroupRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

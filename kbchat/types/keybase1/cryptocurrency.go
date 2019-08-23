// Auto-generated types using avdl-compiler v1.4.1 (https://github.com/keybase/node-avdl-compiler)
//   Input file: ../client/protocol/avdl/keybase1/cryptocurrency.avdl

package keybase1

import (
	context "golang.org/x/net/context"
)

type RegisterAddressRes struct {
	Type   string `codec:"type" json:"type"`
	Family string `codec:"family" json:"family"`
}

func (o RegisterAddressRes) DeepCopy() RegisterAddressRes {
	return RegisterAddressRes{
		Type:   o.Type,
		Family: o.Family,
	}
}

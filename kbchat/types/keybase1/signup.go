// Auto-generated types using avdl-compiler v1.4.1 (https://github.com/keybase/node-avdl-compiler)
//   Input file: ../client/protocol/avdl/keybase1/signup.avdl

package keybase1

import (
	context "golang.org/x/net/context"
)

type SignupRes struct {
	PassphraseOk bool `codec:"passphraseOk" json:"passphraseOk"`
	PostOk       bool `codec:"postOk" json:"postOk"`
	WriteOk      bool `codec:"writeOk" json:"writeOk"`
}

func (o SignupRes) DeepCopy() SignupRes {
	return SignupRes{
		PassphraseOk: o.PassphraseOk,
		PostOk:       o.PostOk,
		WriteOk:      o.WriteOk,
	}
}

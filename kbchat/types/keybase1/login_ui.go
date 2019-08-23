// Auto-generated types using avdl-compiler v1.4.1 (https://github.com/keybase/node-avdl-compiler)
//   Input file: ../client/protocol/avdl/keybase1/login_ui.avdl

package keybase1

import (
	context "golang.org/x/net/context"
)

type ResetPromptType int

const (
	ResetPromptType_COMPLETE         ResetPromptType = 0
	ResetPromptType_ENTER_NO_DEVICES ResetPromptType = 1
	ResetPromptType_ENTER_FORGOT_PW  ResetPromptType = 2
)

func (o ResetPromptType) DeepCopy() ResetPromptType { return o }

var ResetPromptTypeMap = map[string]ResetPromptType{
	"COMPLETE":         0,
	"ENTER_NO_DEVICES": 1,
	"ENTER_FORGOT_PW":  2,
}

var ResetPromptTypeRevMap = map[ResetPromptType]string{
	0: "COMPLETE",
	1: "ENTER_NO_DEVICES",
	2: "ENTER_FORGOT_PW",
}

func (e ResetPromptType) String() string {
	if v, ok := ResetPromptTypeRevMap[e]; ok {
		return v
	}
	return ""
}

type PassphraseRecoveryPromptType int

const (
	PassphraseRecoveryPromptType_ENCRYPTED_PGP_KEYS PassphraseRecoveryPromptType = 0
)

func (o PassphraseRecoveryPromptType) DeepCopy() PassphraseRecoveryPromptType { return o }

var PassphraseRecoveryPromptTypeMap = map[string]PassphraseRecoveryPromptType{
	"ENCRYPTED_PGP_KEYS": 0,
}

var PassphraseRecoveryPromptTypeRevMap = map[PassphraseRecoveryPromptType]string{
	0: "ENCRYPTED_PGP_KEYS",
}

func (e PassphraseRecoveryPromptType) String() string {
	if v, ok := PassphraseRecoveryPromptTypeRevMap[e]; ok {
		return v
	}
	return ""
}

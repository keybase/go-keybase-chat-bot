// Code generated to Go types using avdl-compiler v1.4.10 (https://github.com/keybase/node-avdl-compiler). DO NOT EDIT.
//   Input file: ../client/protocol/avdl/keybase1/debugging.avdl

package keybase1

type FirstStepResult struct {
	ValPlusTwo int `codec:"valPlusTwo" json:"valPlusTwo"`
}

func (o FirstStepResult) DeepCopy() FirstStepResult {
	return FirstStepResult{
		ValPlusTwo: o.ValPlusTwo,
	}
}

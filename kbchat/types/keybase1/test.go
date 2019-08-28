// Auto-generated types using avdl-compiler v1.4.1 (https://github.com/keybase/node-avdl-compiler)
//   Input file: ../client/protocol/avdl/keybase1/test.avdl

package keybase1

// Result from calling test(..).
type Test struct {
	Reply string `codec:"reply" json:"reply"`
}

func (o Test) DeepCopy() Test {
	return Test{
		Reply: o.Reply,
	}
}

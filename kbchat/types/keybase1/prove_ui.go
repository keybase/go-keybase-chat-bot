// Auto-generated types using avdl-compiler v1.4.1 (https://github.com/keybase/node-avdl-compiler)
//   Input file: ../client/protocol/avdl/keybase1/prove_ui.avdl

package keybase1

type PromptOverwriteType int

const (
	PromptOverwriteType_SOCIAL PromptOverwriteType = 0
	PromptOverwriteType_SITE   PromptOverwriteType = 1
)

func (o PromptOverwriteType) DeepCopy() PromptOverwriteType { return o }

var PromptOverwriteTypeMap = map[string]PromptOverwriteType{
	"SOCIAL": 0,
	"SITE":   1,
}

var PromptOverwriteTypeRevMap = map[PromptOverwriteType]string{
	0: "SOCIAL",
	1: "SITE",
}

func (e PromptOverwriteType) String() string {
	if v, ok := PromptOverwriteTypeRevMap[e]; ok {
		return v
	}
	return ""
}

type ProveParameters struct {
	LogoFull    []SizedImage `codec:"logoFull" json:"logoFull"`
	LogoBlack   []SizedImage `codec:"logoBlack" json:"logoBlack"`
	Title       string       `codec:"title" json:"title"`
	Subtext     string       `codec:"subtext" json:"subtext"`
	Suffix      string       `codec:"suffix" json:"suffix"`
	ButtonLabel string       `codec:"buttonLabel" json:"buttonLabel"`
}

func (o ProveParameters) DeepCopy() ProveParameters {
	return ProveParameters{
		LogoFull: (func(x []SizedImage) []SizedImage {
			if x == nil {
				return nil
			}
			ret := make([]SizedImage, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.LogoFull),
		LogoBlack: (func(x []SizedImage) []SizedImage {
			if x == nil {
				return nil
			}
			ret := make([]SizedImage, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.LogoBlack),
		Title:       o.Title,
		Subtext:     o.Subtext,
		Suffix:      o.Suffix,
		ButtonLabel: o.ButtonLabel,
	}
}

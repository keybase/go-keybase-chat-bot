// Auto-generated by avdl-compiler v1.4.0 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/keybase1/metadata_update.avdl

package keybase1

type RekeyRequest struct {
	FolderID string `codec:"folderID" json:"folderID"`
	Revision int64  `codec:"revision" json:"revision"`
}

func (o RekeyRequest) DeepCopy() RekeyRequest {
	return RekeyRequest{
		FolderID: o.FolderID,
		Revision: o.Revision,
	}
}

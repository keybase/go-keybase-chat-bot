// Code generated to Go types using avdl-compiler v1.4.10 (https://github.com/keybase/node-avdl-compiler). DO NOT EDIT.
//   Input file: ../client/protocol/avdl/keybase1/simple_fs.avdl

package keybase1

import (
	"errors"
	"fmt"
)

type OpID [16]byte

func (o OpID) DeepCopy() OpID {
	var ret OpID
	copy(ret[:], o[:])
	return ret
}

type KBFSRevision int64

func (o KBFSRevision) DeepCopy() KBFSRevision {
	return o
}

type KBFSArchivedType int

const (
	KBFSArchivedType_REVISION        KBFSArchivedType = 0
	KBFSArchivedType_TIME            KBFSArchivedType = 1
	KBFSArchivedType_TIME_STRING     KBFSArchivedType = 2
	KBFSArchivedType_REL_TIME_STRING KBFSArchivedType = 3
)

func (o KBFSArchivedType) DeepCopy() KBFSArchivedType { return o }

var KBFSArchivedTypeMap = map[string]KBFSArchivedType{
	"REVISION":        0,
	"TIME":            1,
	"TIME_STRING":     2,
	"REL_TIME_STRING": 3,
}

var KBFSArchivedTypeRevMap = map[KBFSArchivedType]string{
	0: "REVISION",
	1: "TIME",
	2: "TIME_STRING",
	3: "REL_TIME_STRING",
}

func (o KBFSArchivedType) String() string {
	if v, ok := KBFSArchivedTypeRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type KBFSArchivedParam struct {
	KBFSArchivedType__ KBFSArchivedType `codec:"KBFSArchivedType" json:"KBFSArchivedType"`
	Revision__         *KBFSRevision    `codec:"revision,omitempty" json:"revision,omitempty"`
	Time__             *Time            `codec:"time,omitempty" json:"time,omitempty"`
	TimeString__       *string          `codec:"timeString,omitempty" json:"timeString,omitempty"`
	RelTimeString__    *string          `codec:"relTimeString,omitempty" json:"relTimeString,omitempty"`
}

func (o *KBFSArchivedParam) KBFSArchivedType() (ret KBFSArchivedType, err error) {
	switch o.KBFSArchivedType__ {
	case KBFSArchivedType_REVISION:
		if o.Revision__ == nil {
			err = errors.New("unexpected nil value for Revision__")
			return ret, err
		}
	case KBFSArchivedType_TIME:
		if o.Time__ == nil {
			err = errors.New("unexpected nil value for Time__")
			return ret, err
		}
	case KBFSArchivedType_TIME_STRING:
		if o.TimeString__ == nil {
			err = errors.New("unexpected nil value for TimeString__")
			return ret, err
		}
	case KBFSArchivedType_REL_TIME_STRING:
		if o.RelTimeString__ == nil {
			err = errors.New("unexpected nil value for RelTimeString__")
			return ret, err
		}
	}
	return o.KBFSArchivedType__, nil
}

func (o KBFSArchivedParam) Revision() (res KBFSRevision) {
	if o.KBFSArchivedType__ != KBFSArchivedType_REVISION {
		panic("wrong case accessed")
	}
	if o.Revision__ == nil {
		return
	}
	return *o.Revision__
}

func (o KBFSArchivedParam) Time() (res Time) {
	if o.KBFSArchivedType__ != KBFSArchivedType_TIME {
		panic("wrong case accessed")
	}
	if o.Time__ == nil {
		return
	}
	return *o.Time__
}

func (o KBFSArchivedParam) TimeString() (res string) {
	if o.KBFSArchivedType__ != KBFSArchivedType_TIME_STRING {
		panic("wrong case accessed")
	}
	if o.TimeString__ == nil {
		return
	}
	return *o.TimeString__
}

func (o KBFSArchivedParam) RelTimeString() (res string) {
	if o.KBFSArchivedType__ != KBFSArchivedType_REL_TIME_STRING {
		panic("wrong case accessed")
	}
	if o.RelTimeString__ == nil {
		return
	}
	return *o.RelTimeString__
}

func NewKBFSArchivedParamWithRevision(v KBFSRevision) KBFSArchivedParam {
	return KBFSArchivedParam{
		KBFSArchivedType__: KBFSArchivedType_REVISION,
		Revision__:         &v,
	}
}

func NewKBFSArchivedParamWithTime(v Time) KBFSArchivedParam {
	return KBFSArchivedParam{
		KBFSArchivedType__: KBFSArchivedType_TIME,
		Time__:             &v,
	}
}

func NewKBFSArchivedParamWithTimeString(v string) KBFSArchivedParam {
	return KBFSArchivedParam{
		KBFSArchivedType__: KBFSArchivedType_TIME_STRING,
		TimeString__:       &v,
	}
}

func NewKBFSArchivedParamWithRelTimeString(v string) KBFSArchivedParam {
	return KBFSArchivedParam{
		KBFSArchivedType__: KBFSArchivedType_REL_TIME_STRING,
		RelTimeString__:    &v,
	}
}

func (o KBFSArchivedParam) DeepCopy() KBFSArchivedParam {
	return KBFSArchivedParam{
		KBFSArchivedType__: o.KBFSArchivedType__.DeepCopy(),
		Revision__: (func(x *KBFSRevision) *KBFSRevision {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.Revision__),
		Time__: (func(x *Time) *Time {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.Time__),
		TimeString__: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.TimeString__),
		RelTimeString__: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.RelTimeString__),
	}
}

type KBFSArchivedPath struct {
	Path             string               `codec:"path" json:"path"`
	ArchivedParam    KBFSArchivedParam    `codec:"archivedParam" json:"archivedParam"`
	IdentifyBehavior *TLFIdentifyBehavior `codec:"identifyBehavior,omitempty" json:"identifyBehavior,omitempty"`
}

func (o KBFSArchivedPath) DeepCopy() KBFSArchivedPath {
	return KBFSArchivedPath{
		Path:          o.Path,
		ArchivedParam: o.ArchivedParam.DeepCopy(),
		IdentifyBehavior: (func(x *TLFIdentifyBehavior) *TLFIdentifyBehavior {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.IdentifyBehavior),
	}
}

type KBFSPath struct {
	Path             string               `codec:"path" json:"path"`
	IdentifyBehavior *TLFIdentifyBehavior `codec:"identifyBehavior,omitempty" json:"identifyBehavior,omitempty"`
}

func (o KBFSPath) DeepCopy() KBFSPath {
	return KBFSPath{
		Path: o.Path,
		IdentifyBehavior: (func(x *TLFIdentifyBehavior) *TLFIdentifyBehavior {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.IdentifyBehavior),
	}
}

type PathType int

const (
	PathType_LOCAL         PathType = 0
	PathType_KBFS          PathType = 1
	PathType_KBFS_ARCHIVED PathType = 2
)

func (o PathType) DeepCopy() PathType { return o }

var PathTypeMap = map[string]PathType{
	"LOCAL":         0,
	"KBFS":          1,
	"KBFS_ARCHIVED": 2,
}

var PathTypeRevMap = map[PathType]string{
	0: "LOCAL",
	1: "KBFS",
	2: "KBFS_ARCHIVED",
}

func (o PathType) String() string {
	if v, ok := PathTypeRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type Path struct {
	PathType__     PathType          `codec:"PathType" json:"PathType"`
	Local__        *string           `codec:"local,omitempty" json:"local,omitempty"`
	Kbfs__         *KBFSPath         `codec:"kbfs,omitempty" json:"kbfs,omitempty"`
	KbfsArchived__ *KBFSArchivedPath `codec:"kbfsArchived,omitempty" json:"kbfsArchived,omitempty"`
}

func (o *Path) PathType() (ret PathType, err error) {
	switch o.PathType__ {
	case PathType_LOCAL:
		if o.Local__ == nil {
			err = errors.New("unexpected nil value for Local__")
			return ret, err
		}
	case PathType_KBFS:
		if o.Kbfs__ == nil {
			err = errors.New("unexpected nil value for Kbfs__")
			return ret, err
		}
	case PathType_KBFS_ARCHIVED:
		if o.KbfsArchived__ == nil {
			err = errors.New("unexpected nil value for KbfsArchived__")
			return ret, err
		}
	}
	return o.PathType__, nil
}

func (o Path) Local() (res string) {
	if o.PathType__ != PathType_LOCAL {
		panic("wrong case accessed")
	}
	if o.Local__ == nil {
		return
	}
	return *o.Local__
}

func (o Path) Kbfs() (res KBFSPath) {
	if o.PathType__ != PathType_KBFS {
		panic("wrong case accessed")
	}
	if o.Kbfs__ == nil {
		return
	}
	return *o.Kbfs__
}

func (o Path) KbfsArchived() (res KBFSArchivedPath) {
	if o.PathType__ != PathType_KBFS_ARCHIVED {
		panic("wrong case accessed")
	}
	if o.KbfsArchived__ == nil {
		return
	}
	return *o.KbfsArchived__
}

func NewPathWithLocal(v string) Path {
	return Path{
		PathType__: PathType_LOCAL,
		Local__:    &v,
	}
}

func NewPathWithKbfs(v KBFSPath) Path {
	return Path{
		PathType__: PathType_KBFS,
		Kbfs__:     &v,
	}
}

func NewPathWithKbfsArchived(v KBFSArchivedPath) Path {
	return Path{
		PathType__:     PathType_KBFS_ARCHIVED,
		KbfsArchived__: &v,
	}
}

func (o Path) DeepCopy() Path {
	return Path{
		PathType__: o.PathType__.DeepCopy(),
		Local__: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.Local__),
		Kbfs__: (func(x *KBFSPath) *KBFSPath {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.Kbfs__),
		KbfsArchived__: (func(x *KBFSArchivedPath) *KBFSArchivedPath {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.KbfsArchived__),
	}
}

type DirentType int

const (
	DirentType_FILE DirentType = 0
	DirentType_DIR  DirentType = 1
	DirentType_SYM  DirentType = 2
	DirentType_EXEC DirentType = 3
)

func (o DirentType) DeepCopy() DirentType { return o }

var DirentTypeMap = map[string]DirentType{
	"FILE": 0,
	"DIR":  1,
	"SYM":  2,
	"EXEC": 3,
}

var DirentTypeRevMap = map[DirentType]string{
	0: "FILE",
	1: "DIR",
	2: "SYM",
	3: "EXEC",
}

func (o DirentType) String() string {
	if v, ok := DirentTypeRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type PrefetchStatus int

const (
	PrefetchStatus_NOT_STARTED PrefetchStatus = 0
	PrefetchStatus_IN_PROGRESS PrefetchStatus = 1
	PrefetchStatus_COMPLETE    PrefetchStatus = 2
)

func (o PrefetchStatus) DeepCopy() PrefetchStatus { return o }

var PrefetchStatusMap = map[string]PrefetchStatus{
	"NOT_STARTED": 0,
	"IN_PROGRESS": 1,
	"COMPLETE":    2,
}

var PrefetchStatusRevMap = map[PrefetchStatus]string{
	0: "NOT_STARTED",
	1: "IN_PROGRESS",
	2: "COMPLETE",
}

func (o PrefetchStatus) String() string {
	if v, ok := PrefetchStatusRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type PrefetchProgress struct {
	Start        Time  `codec:"start" json:"start"`
	EndEstimate  Time  `codec:"endEstimate" json:"endEstimate"`
	BytesTotal   int64 `codec:"bytesTotal" json:"bytesTotal"`
	BytesFetched int64 `codec:"bytesFetched" json:"bytesFetched"`
}

func (o PrefetchProgress) DeepCopy() PrefetchProgress {
	return PrefetchProgress{
		Start:        o.Start.DeepCopy(),
		EndEstimate:  o.EndEstimate.DeepCopy(),
		BytesTotal:   o.BytesTotal,
		BytesFetched: o.BytesFetched,
	}
}

type Dirent struct {
	Time                 Time             `codec:"time" json:"time"`
	Size                 int              `codec:"size" json:"size"`
	Name                 string           `codec:"name" json:"name"`
	DirentType           DirentType       `codec:"direntType" json:"direntType"`
	LastWriterUnverified User             `codec:"lastWriterUnverified" json:"lastWriterUnverified"`
	Writable             bool             `codec:"writable" json:"writable"`
	PrefetchStatus       PrefetchStatus   `codec:"prefetchStatus" json:"prefetchStatus"`
	PrefetchProgress     PrefetchProgress `codec:"prefetchProgress" json:"prefetchProgress"`
	SymlinkTarget        string           `codec:"symlinkTarget" json:"symlinkTarget"`
}

func (o Dirent) DeepCopy() Dirent {
	return Dirent{
		Time:                 o.Time.DeepCopy(),
		Size:                 o.Size,
		Name:                 o.Name,
		DirentType:           o.DirentType.DeepCopy(),
		LastWriterUnverified: o.LastWriterUnverified.DeepCopy(),
		Writable:             o.Writable,
		PrefetchStatus:       o.PrefetchStatus.DeepCopy(),
		PrefetchProgress:     o.PrefetchProgress.DeepCopy(),
		SymlinkTarget:        o.SymlinkTarget,
	}
}

type DirentWithRevision struct {
	Entry    Dirent       `codec:"entry" json:"entry"`
	Revision KBFSRevision `codec:"revision" json:"revision"`
}

func (o DirentWithRevision) DeepCopy() DirentWithRevision {
	return DirentWithRevision{
		Entry:    o.Entry.DeepCopy(),
		Revision: o.Revision.DeepCopy(),
	}
}

type RevisionSpanType int

const (
	RevisionSpanType_DEFAULT   RevisionSpanType = 0
	RevisionSpanType_LAST_FIVE RevisionSpanType = 1
)

func (o RevisionSpanType) DeepCopy() RevisionSpanType { return o }

var RevisionSpanTypeMap = map[string]RevisionSpanType{
	"DEFAULT":   0,
	"LAST_FIVE": 1,
}

var RevisionSpanTypeRevMap = map[RevisionSpanType]string{
	0: "DEFAULT",
	1: "LAST_FIVE",
}

func (o RevisionSpanType) String() string {
	if v, ok := RevisionSpanTypeRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type ErrorNum int

func (o ErrorNum) DeepCopy() ErrorNum {
	return o
}

type OpenFlags int

const (
	OpenFlags_READ      OpenFlags = 0
	OpenFlags_REPLACE   OpenFlags = 1
	OpenFlags_EXISTING  OpenFlags = 2
	OpenFlags_WRITE     OpenFlags = 4
	OpenFlags_APPEND    OpenFlags = 8
	OpenFlags_DIRECTORY OpenFlags = 16
)

func (o OpenFlags) DeepCopy() OpenFlags { return o }

var OpenFlagsMap = map[string]OpenFlags{
	"READ":      0,
	"REPLACE":   1,
	"EXISTING":  2,
	"WRITE":     4,
	"APPEND":    8,
	"DIRECTORY": 16,
}

var OpenFlagsRevMap = map[OpenFlags]string{
	0:  "READ",
	1:  "REPLACE",
	2:  "EXISTING",
	4:  "WRITE",
	8:  "APPEND",
	16: "DIRECTORY",
}

func (o OpenFlags) String() string {
	if v, ok := OpenFlagsRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type Progress int

func (o Progress) DeepCopy() Progress {
	return o
}

type SimpleFSListResult struct {
	Entries  []Dirent `codec:"entries" json:"entries"`
	Progress Progress `codec:"progress" json:"progress"`
}

func (o SimpleFSListResult) DeepCopy() SimpleFSListResult {
	return SimpleFSListResult{
		Entries: (func(x []Dirent) []Dirent {
			if x == nil {
				return nil
			}
			ret := make([]Dirent, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.Entries),
		Progress: o.Progress.DeepCopy(),
	}
}

type FileContent struct {
	Data     []byte   `codec:"data" json:"data"`
	Progress Progress `codec:"progress" json:"progress"`
}

func (o FileContent) DeepCopy() FileContent {
	return FileContent{
		Data: (func(x []byte) []byte {
			if x == nil {
				return nil
			}
			return append([]byte{}, x...)
		})(o.Data),
		Progress: o.Progress.DeepCopy(),
	}
}

type AsyncOps int

const (
	AsyncOps_LIST                    AsyncOps = 0
	AsyncOps_LIST_RECURSIVE          AsyncOps = 1
	AsyncOps_READ                    AsyncOps = 2
	AsyncOps_WRITE                   AsyncOps = 3
	AsyncOps_COPY                    AsyncOps = 4
	AsyncOps_MOVE                    AsyncOps = 5
	AsyncOps_REMOVE                  AsyncOps = 6
	AsyncOps_LIST_RECURSIVE_TO_DEPTH AsyncOps = 7
	AsyncOps_GET_REVISIONS           AsyncOps = 8
)

func (o AsyncOps) DeepCopy() AsyncOps { return o }

var AsyncOpsMap = map[string]AsyncOps{
	"LIST":                    0,
	"LIST_RECURSIVE":          1,
	"READ":                    2,
	"WRITE":                   3,
	"COPY":                    4,
	"MOVE":                    5,
	"REMOVE":                  6,
	"LIST_RECURSIVE_TO_DEPTH": 7,
	"GET_REVISIONS":           8,
}

var AsyncOpsRevMap = map[AsyncOps]string{
	0: "LIST",
	1: "LIST_RECURSIVE",
	2: "READ",
	3: "WRITE",
	4: "COPY",
	5: "MOVE",
	6: "REMOVE",
	7: "LIST_RECURSIVE_TO_DEPTH",
	8: "GET_REVISIONS",
}

func (o AsyncOps) String() string {
	if v, ok := AsyncOpsRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type ListFilter int

const (
	ListFilter_NO_FILTER            ListFilter = 0
	ListFilter_FILTER_ALL_HIDDEN    ListFilter = 1
	ListFilter_FILTER_SYSTEM_HIDDEN ListFilter = 2
)

func (o ListFilter) DeepCopy() ListFilter { return o }

var ListFilterMap = map[string]ListFilter{
	"NO_FILTER":            0,
	"FILTER_ALL_HIDDEN":    1,
	"FILTER_SYSTEM_HIDDEN": 2,
}

var ListFilterRevMap = map[ListFilter]string{
	0: "NO_FILTER",
	1: "FILTER_ALL_HIDDEN",
	2: "FILTER_SYSTEM_HIDDEN",
}

func (o ListFilter) String() string {
	if v, ok := ListFilterRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type ListArgs struct {
	OpID   OpID       `codec:"opID" json:"opID"`
	Path   Path       `codec:"path" json:"path"`
	Filter ListFilter `codec:"filter" json:"filter"`
}

func (o ListArgs) DeepCopy() ListArgs {
	return ListArgs{
		OpID:   o.OpID.DeepCopy(),
		Path:   o.Path.DeepCopy(),
		Filter: o.Filter.DeepCopy(),
	}
}

type ListToDepthArgs struct {
	OpID   OpID       `codec:"opID" json:"opID"`
	Path   Path       `codec:"path" json:"path"`
	Filter ListFilter `codec:"filter" json:"filter"`
	Depth  int        `codec:"depth" json:"depth"`
}

func (o ListToDepthArgs) DeepCopy() ListToDepthArgs {
	return ListToDepthArgs{
		OpID:   o.OpID.DeepCopy(),
		Path:   o.Path.DeepCopy(),
		Filter: o.Filter.DeepCopy(),
		Depth:  o.Depth,
	}
}

type RemoveArgs struct {
	OpID      OpID `codec:"opID" json:"opID"`
	Path      Path `codec:"path" json:"path"`
	Recursive bool `codec:"recursive" json:"recursive"`
}

func (o RemoveArgs) DeepCopy() RemoveArgs {
	return RemoveArgs{
		OpID:      o.OpID.DeepCopy(),
		Path:      o.Path.DeepCopy(),
		Recursive: o.Recursive,
	}
}

type ReadArgs struct {
	OpID   OpID  `codec:"opID" json:"opID"`
	Path   Path  `codec:"path" json:"path"`
	Offset int64 `codec:"offset" json:"offset"`
	Size   int   `codec:"size" json:"size"`
}

func (o ReadArgs) DeepCopy() ReadArgs {
	return ReadArgs{
		OpID:   o.OpID.DeepCopy(),
		Path:   o.Path.DeepCopy(),
		Offset: o.Offset,
		Size:   o.Size,
	}
}

type WriteArgs struct {
	OpID   OpID  `codec:"opID" json:"opID"`
	Path   Path  `codec:"path" json:"path"`
	Offset int64 `codec:"offset" json:"offset"`
}

func (o WriteArgs) DeepCopy() WriteArgs {
	return WriteArgs{
		OpID:   o.OpID.DeepCopy(),
		Path:   o.Path.DeepCopy(),
		Offset: o.Offset,
	}
}

type CopyArgs struct {
	OpID                   OpID `codec:"opID" json:"opID"`
	Src                    Path `codec:"src" json:"src"`
	Dest                   Path `codec:"dest" json:"dest"`
	OverwriteExistingFiles bool `codec:"overwriteExistingFiles" json:"overwriteExistingFiles"`
}

func (o CopyArgs) DeepCopy() CopyArgs {
	return CopyArgs{
		OpID:                   o.OpID.DeepCopy(),
		Src:                    o.Src.DeepCopy(),
		Dest:                   o.Dest.DeepCopy(),
		OverwriteExistingFiles: o.OverwriteExistingFiles,
	}
}

type MoveArgs struct {
	OpID                   OpID `codec:"opID" json:"opID"`
	Src                    Path `codec:"src" json:"src"`
	Dest                   Path `codec:"dest" json:"dest"`
	OverwriteExistingFiles bool `codec:"overwriteExistingFiles" json:"overwriteExistingFiles"`
}

func (o MoveArgs) DeepCopy() MoveArgs {
	return MoveArgs{
		OpID:                   o.OpID.DeepCopy(),
		Src:                    o.Src.DeepCopy(),
		Dest:                   o.Dest.DeepCopy(),
		OverwriteExistingFiles: o.OverwriteExistingFiles,
	}
}

type GetRevisionsArgs struct {
	OpID     OpID             `codec:"opID" json:"opID"`
	Path     Path             `codec:"path" json:"path"`
	SpanType RevisionSpanType `codec:"spanType" json:"spanType"`
}

func (o GetRevisionsArgs) DeepCopy() GetRevisionsArgs {
	return GetRevisionsArgs{
		OpID:     o.OpID.DeepCopy(),
		Path:     o.Path.DeepCopy(),
		SpanType: o.SpanType.DeepCopy(),
	}
}

type OpDescription struct {
	AsyncOp__              AsyncOps          `codec:"asyncOp" json:"asyncOp"`
	List__                 *ListArgs         `codec:"list,omitempty" json:"list,omitempty"`
	ListRecursive__        *ListArgs         `codec:"listRecursive,omitempty" json:"listRecursive,omitempty"`
	ListRecursiveToDepth__ *ListToDepthArgs  `codec:"listRecursiveToDepth,omitempty" json:"listRecursiveToDepth,omitempty"`
	Read__                 *ReadArgs         `codec:"read,omitempty" json:"read,omitempty"`
	Write__                *WriteArgs        `codec:"write,omitempty" json:"write,omitempty"`
	Copy__                 *CopyArgs         `codec:"copy,omitempty" json:"copy,omitempty"`
	Move__                 *MoveArgs         `codec:"move,omitempty" json:"move,omitempty"`
	Remove__               *RemoveArgs       `codec:"remove,omitempty" json:"remove,omitempty"`
	GetRevisions__         *GetRevisionsArgs `codec:"getRevisions,omitempty" json:"getRevisions,omitempty"`
}

func (o *OpDescription) AsyncOp() (ret AsyncOps, err error) {
	switch o.AsyncOp__ {
	case AsyncOps_LIST:
		if o.List__ == nil {
			err = errors.New("unexpected nil value for List__")
			return ret, err
		}
	case AsyncOps_LIST_RECURSIVE:
		if o.ListRecursive__ == nil {
			err = errors.New("unexpected nil value for ListRecursive__")
			return ret, err
		}
	case AsyncOps_LIST_RECURSIVE_TO_DEPTH:
		if o.ListRecursiveToDepth__ == nil {
			err = errors.New("unexpected nil value for ListRecursiveToDepth__")
			return ret, err
		}
	case AsyncOps_READ:
		if o.Read__ == nil {
			err = errors.New("unexpected nil value for Read__")
			return ret, err
		}
	case AsyncOps_WRITE:
		if o.Write__ == nil {
			err = errors.New("unexpected nil value for Write__")
			return ret, err
		}
	case AsyncOps_COPY:
		if o.Copy__ == nil {
			err = errors.New("unexpected nil value for Copy__")
			return ret, err
		}
	case AsyncOps_MOVE:
		if o.Move__ == nil {
			err = errors.New("unexpected nil value for Move__")
			return ret, err
		}
	case AsyncOps_REMOVE:
		if o.Remove__ == nil {
			err = errors.New("unexpected nil value for Remove__")
			return ret, err
		}
	case AsyncOps_GET_REVISIONS:
		if o.GetRevisions__ == nil {
			err = errors.New("unexpected nil value for GetRevisions__")
			return ret, err
		}
	}
	return o.AsyncOp__, nil
}

func (o OpDescription) List() (res ListArgs) {
	if o.AsyncOp__ != AsyncOps_LIST {
		panic("wrong case accessed")
	}
	if o.List__ == nil {
		return
	}
	return *o.List__
}

func (o OpDescription) ListRecursive() (res ListArgs) {
	if o.AsyncOp__ != AsyncOps_LIST_RECURSIVE {
		panic("wrong case accessed")
	}
	if o.ListRecursive__ == nil {
		return
	}
	return *o.ListRecursive__
}

func (o OpDescription) ListRecursiveToDepth() (res ListToDepthArgs) {
	if o.AsyncOp__ != AsyncOps_LIST_RECURSIVE_TO_DEPTH {
		panic("wrong case accessed")
	}
	if o.ListRecursiveToDepth__ == nil {
		return
	}
	return *o.ListRecursiveToDepth__
}

func (o OpDescription) Read() (res ReadArgs) {
	if o.AsyncOp__ != AsyncOps_READ {
		panic("wrong case accessed")
	}
	if o.Read__ == nil {
		return
	}
	return *o.Read__
}

func (o OpDescription) Write() (res WriteArgs) {
	if o.AsyncOp__ != AsyncOps_WRITE {
		panic("wrong case accessed")
	}
	if o.Write__ == nil {
		return
	}
	return *o.Write__
}

func (o OpDescription) Copy() (res CopyArgs) {
	if o.AsyncOp__ != AsyncOps_COPY {
		panic("wrong case accessed")
	}
	if o.Copy__ == nil {
		return
	}
	return *o.Copy__
}

func (o OpDescription) Move() (res MoveArgs) {
	if o.AsyncOp__ != AsyncOps_MOVE {
		panic("wrong case accessed")
	}
	if o.Move__ == nil {
		return
	}
	return *o.Move__
}

func (o OpDescription) Remove() (res RemoveArgs) {
	if o.AsyncOp__ != AsyncOps_REMOVE {
		panic("wrong case accessed")
	}
	if o.Remove__ == nil {
		return
	}
	return *o.Remove__
}

func (o OpDescription) GetRevisions() (res GetRevisionsArgs) {
	if o.AsyncOp__ != AsyncOps_GET_REVISIONS {
		panic("wrong case accessed")
	}
	if o.GetRevisions__ == nil {
		return
	}
	return *o.GetRevisions__
}

func NewOpDescriptionWithList(v ListArgs) OpDescription {
	return OpDescription{
		AsyncOp__: AsyncOps_LIST,
		List__:    &v,
	}
}

func NewOpDescriptionWithListRecursive(v ListArgs) OpDescription {
	return OpDescription{
		AsyncOp__:       AsyncOps_LIST_RECURSIVE,
		ListRecursive__: &v,
	}
}

func NewOpDescriptionWithListRecursiveToDepth(v ListToDepthArgs) OpDescription {
	return OpDescription{
		AsyncOp__:              AsyncOps_LIST_RECURSIVE_TO_DEPTH,
		ListRecursiveToDepth__: &v,
	}
}

func NewOpDescriptionWithRead(v ReadArgs) OpDescription {
	return OpDescription{
		AsyncOp__: AsyncOps_READ,
		Read__:    &v,
	}
}

func NewOpDescriptionWithWrite(v WriteArgs) OpDescription {
	return OpDescription{
		AsyncOp__: AsyncOps_WRITE,
		Write__:   &v,
	}
}

func NewOpDescriptionWithCopy(v CopyArgs) OpDescription {
	return OpDescription{
		AsyncOp__: AsyncOps_COPY,
		Copy__:    &v,
	}
}

func NewOpDescriptionWithMove(v MoveArgs) OpDescription {
	return OpDescription{
		AsyncOp__: AsyncOps_MOVE,
		Move__:    &v,
	}
}

func NewOpDescriptionWithRemove(v RemoveArgs) OpDescription {
	return OpDescription{
		AsyncOp__: AsyncOps_REMOVE,
		Remove__:  &v,
	}
}

func NewOpDescriptionWithGetRevisions(v GetRevisionsArgs) OpDescription {
	return OpDescription{
		AsyncOp__:      AsyncOps_GET_REVISIONS,
		GetRevisions__: &v,
	}
}

func (o OpDescription) DeepCopy() OpDescription {
	return OpDescription{
		AsyncOp__: o.AsyncOp__.DeepCopy(),
		List__: (func(x *ListArgs) *ListArgs {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.List__),
		ListRecursive__: (func(x *ListArgs) *ListArgs {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.ListRecursive__),
		ListRecursiveToDepth__: (func(x *ListToDepthArgs) *ListToDepthArgs {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.ListRecursiveToDepth__),
		Read__: (func(x *ReadArgs) *ReadArgs {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.Read__),
		Write__: (func(x *WriteArgs) *WriteArgs {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.Write__),
		Copy__: (func(x *CopyArgs) *CopyArgs {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.Copy__),
		Move__: (func(x *MoveArgs) *MoveArgs {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.Move__),
		Remove__: (func(x *RemoveArgs) *RemoveArgs {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.Remove__),
		GetRevisions__: (func(x *GetRevisionsArgs) *GetRevisionsArgs {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.GetRevisions__),
	}
}

type GetRevisionsResult struct {
	Revisions []DirentWithRevision `codec:"revisions" json:"revisions"`
	Progress  Progress             `codec:"progress" json:"progress"`
}

func (o GetRevisionsResult) DeepCopy() GetRevisionsResult {
	return GetRevisionsResult{
		Revisions: (func(x []DirentWithRevision) []DirentWithRevision {
			if x == nil {
				return nil
			}
			ret := make([]DirentWithRevision, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.Revisions),
		Progress: o.Progress.DeepCopy(),
	}
}

type OpProgress struct {
	Start        Time     `codec:"start" json:"start"`
	EndEstimate  Time     `codec:"endEstimate" json:"endEstimate"`
	OpType       AsyncOps `codec:"opType" json:"opType"`
	BytesTotal   int64    `codec:"bytesTotal" json:"bytesTotal"`
	BytesRead    int64    `codec:"bytesRead" json:"bytesRead"`
	BytesWritten int64    `codec:"bytesWritten" json:"bytesWritten"`
	FilesTotal   int64    `codec:"filesTotal" json:"filesTotal"`
	FilesRead    int64    `codec:"filesRead" json:"filesRead"`
	FilesWritten int64    `codec:"filesWritten" json:"filesWritten"`
}

func (o OpProgress) DeepCopy() OpProgress {
	return OpProgress{
		Start:        o.Start.DeepCopy(),
		EndEstimate:  o.EndEstimate.DeepCopy(),
		OpType:       o.OpType.DeepCopy(),
		BytesTotal:   o.BytesTotal,
		BytesRead:    o.BytesRead,
		BytesWritten: o.BytesWritten,
		FilesTotal:   o.FilesTotal,
		FilesRead:    o.FilesRead,
		FilesWritten: o.FilesWritten,
	}
}

type SimpleFSQuotaUsage struct {
	UsageBytes      int64 `codec:"usageBytes" json:"usageBytes"`
	ArchiveBytes    int64 `codec:"archiveBytes" json:"archiveBytes"`
	LimitBytes      int64 `codec:"limitBytes" json:"limitBytes"`
	GitUsageBytes   int64 `codec:"gitUsageBytes" json:"gitUsageBytes"`
	GitArchiveBytes int64 `codec:"gitArchiveBytes" json:"gitArchiveBytes"`
	GitLimitBytes   int64 `codec:"gitLimitBytes" json:"gitLimitBytes"`
}

func (o SimpleFSQuotaUsage) DeepCopy() SimpleFSQuotaUsage {
	return SimpleFSQuotaUsage{
		UsageBytes:      o.UsageBytes,
		ArchiveBytes:    o.ArchiveBytes,
		LimitBytes:      o.LimitBytes,
		GitUsageBytes:   o.GitUsageBytes,
		GitArchiveBytes: o.GitArchiveBytes,
		GitLimitBytes:   o.GitLimitBytes,
	}
}

type FolderSyncMode int

const (
	FolderSyncMode_DISABLED FolderSyncMode = 0
	FolderSyncMode_ENABLED  FolderSyncMode = 1
	FolderSyncMode_PARTIAL  FolderSyncMode = 2
)

func (o FolderSyncMode) DeepCopy() FolderSyncMode { return o }

var FolderSyncModeMap = map[string]FolderSyncMode{
	"DISABLED": 0,
	"ENABLED":  1,
	"PARTIAL":  2,
}

var FolderSyncModeRevMap = map[FolderSyncMode]string{
	0: "DISABLED",
	1: "ENABLED",
	2: "PARTIAL",
}

func (o FolderSyncMode) String() string {
	if v, ok := FolderSyncModeRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type FolderSyncConfig struct {
	Mode  FolderSyncMode `codec:"mode" json:"mode"`
	Paths []string       `codec:"paths" json:"paths"`
}

func (o FolderSyncConfig) DeepCopy() FolderSyncConfig {
	return FolderSyncConfig{
		Mode: o.Mode.DeepCopy(),
		Paths: (func(x []string) []string {
			if x == nil {
				return nil
			}
			ret := make([]string, len(x))
			for i, v := range x {
				vCopy := v
				ret[i] = vCopy
			}
			return ret
		})(o.Paths),
	}
}

type FolderSyncConfigAndStatus struct {
	Config FolderSyncConfig `codec:"config" json:"config"`
	Status FolderSyncStatus `codec:"status" json:"status"`
}

func (o FolderSyncConfigAndStatus) DeepCopy() FolderSyncConfigAndStatus {
	return FolderSyncConfigAndStatus{
		Config: o.Config.DeepCopy(),
		Status: o.Status.DeepCopy(),
	}
}

type FolderSyncConfigAndStatusWithFolder struct {
	Folder Folder           `codec:"folder" json:"folder"`
	Config FolderSyncConfig `codec:"config" json:"config"`
	Status FolderSyncStatus `codec:"status" json:"status"`
}

func (o FolderSyncConfigAndStatusWithFolder) DeepCopy() FolderSyncConfigAndStatusWithFolder {
	return FolderSyncConfigAndStatusWithFolder{
		Folder: o.Folder.DeepCopy(),
		Config: o.Config.DeepCopy(),
		Status: o.Status.DeepCopy(),
	}
}

type SyncConfigAndStatusRes struct {
	Folders       []FolderSyncConfigAndStatusWithFolder `codec:"folders" json:"folders"`
	OverallStatus FolderSyncStatus                      `codec:"overallStatus" json:"overallStatus"`
}

func (o SyncConfigAndStatusRes) DeepCopy() SyncConfigAndStatusRes {
	return SyncConfigAndStatusRes{
		Folders: (func(x []FolderSyncConfigAndStatusWithFolder) []FolderSyncConfigAndStatusWithFolder {
			if x == nil {
				return nil
			}
			ret := make([]FolderSyncConfigAndStatusWithFolder, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.Folders),
		OverallStatus: o.OverallStatus.DeepCopy(),
	}
}

type FolderWithFavFlags struct {
	Folder     Folder `codec:"folder" json:"folder"`
	IsFavorite bool   `codec:"isFavorite" json:"isFavorite"`
	IsIgnored  bool   `codec:"isIgnored" json:"isIgnored"`
	IsNew      bool   `codec:"isNew" json:"isNew"`
}

func (o FolderWithFavFlags) DeepCopy() FolderWithFavFlags {
	return FolderWithFavFlags{
		Folder:     o.Folder.DeepCopy(),
		IsFavorite: o.IsFavorite,
		IsIgnored:  o.IsIgnored,
		IsNew:      o.IsNew,
	}
}

type KbfsOnlineStatus int

const (
	KbfsOnlineStatus_OFFLINE KbfsOnlineStatus = 0
	KbfsOnlineStatus_TRYING  KbfsOnlineStatus = 1
	KbfsOnlineStatus_ONLINE  KbfsOnlineStatus = 2
)

func (o KbfsOnlineStatus) DeepCopy() KbfsOnlineStatus { return o }

var KbfsOnlineStatusMap = map[string]KbfsOnlineStatus{
	"OFFLINE": 0,
	"TRYING":  1,
	"ONLINE":  2,
}

var KbfsOnlineStatusRevMap = map[KbfsOnlineStatus]string{
	0: "OFFLINE",
	1: "TRYING",
	2: "ONLINE",
}

func (o KbfsOnlineStatus) String() string {
	if v, ok := KbfsOnlineStatusRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type FSSettings struct {
	SpaceAvailableNotificationThreshold int64 `codec:"spaceAvailableNotificationThreshold" json:"spaceAvailableNotificationThreshold"`
	SfmiBannerDismissed                 bool  `codec:"sfmiBannerDismissed" json:"sfmiBannerDismissed"`
	SyncOnCellular                      bool  `codec:"syncOnCellular" json:"syncOnCellular"`
}

func (o FSSettings) DeepCopy() FSSettings {
	return FSSettings{
		SpaceAvailableNotificationThreshold: o.SpaceAvailableNotificationThreshold,
		SfmiBannerDismissed:                 o.SfmiBannerDismissed,
		SyncOnCellular:                      o.SyncOnCellular,
	}
}

type SimpleFSStats struct {
	ProcessStats      ProcessRuntimeStats `codec:"processStats" json:"processStats"`
	BlockCacheDbStats []string            `codec:"blockCacheDbStats" json:"blockCacheDbStats"`
	SyncCacheDbStats  []string            `codec:"syncCacheDbStats" json:"syncCacheDbStats"`
	RuntimeDbStats    []DbStats           `codec:"runtimeDbStats" json:"runtimeDbStats"`
}

func (o SimpleFSStats) DeepCopy() SimpleFSStats {
	return SimpleFSStats{
		ProcessStats: o.ProcessStats.DeepCopy(),
		BlockCacheDbStats: (func(x []string) []string {
			if x == nil {
				return nil
			}
			ret := make([]string, len(x))
			for i, v := range x {
				vCopy := v
				ret[i] = vCopy
			}
			return ret
		})(o.BlockCacheDbStats),
		SyncCacheDbStats: (func(x []string) []string {
			if x == nil {
				return nil
			}
			ret := make([]string, len(x))
			for i, v := range x {
				vCopy := v
				ret[i] = vCopy
			}
			return ret
		})(o.SyncCacheDbStats),
		RuntimeDbStats: (func(x []DbStats) []DbStats {
			if x == nil {
				return nil
			}
			ret := make([]DbStats, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.RuntimeDbStats),
	}
}

type SubscriptionTopic int

const (
	SubscriptionTopic_FAVORITES           SubscriptionTopic = 0
	SubscriptionTopic_JOURNAL_STATUS      SubscriptionTopic = 1
	SubscriptionTopic_ONLINE_STATUS       SubscriptionTopic = 2
	SubscriptionTopic_DOWNLOAD_STATUS     SubscriptionTopic = 3
	SubscriptionTopic_FILES_TAB_BADGE     SubscriptionTopic = 4
	SubscriptionTopic_OVERALL_SYNC_STATUS SubscriptionTopic = 5
	SubscriptionTopic_SETTINGS            SubscriptionTopic = 6
	SubscriptionTopic_UPLOAD_STATUS       SubscriptionTopic = 7
)

func (o SubscriptionTopic) DeepCopy() SubscriptionTopic { return o }

var SubscriptionTopicMap = map[string]SubscriptionTopic{
	"FAVORITES":           0,
	"JOURNAL_STATUS":      1,
	"ONLINE_STATUS":       2,
	"DOWNLOAD_STATUS":     3,
	"FILES_TAB_BADGE":     4,
	"OVERALL_SYNC_STATUS": 5,
	"SETTINGS":            6,
	"UPLOAD_STATUS":       7,
}

var SubscriptionTopicRevMap = map[SubscriptionTopic]string{
	0: "FAVORITES",
	1: "JOURNAL_STATUS",
	2: "ONLINE_STATUS",
	3: "DOWNLOAD_STATUS",
	4: "FILES_TAB_BADGE",
	5: "OVERALL_SYNC_STATUS",
	6: "SETTINGS",
	7: "UPLOAD_STATUS",
}

func (o SubscriptionTopic) String() string {
	if v, ok := SubscriptionTopicRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type PathSubscriptionTopic int

const (
	PathSubscriptionTopic_CHILDREN PathSubscriptionTopic = 0
	PathSubscriptionTopic_STAT     PathSubscriptionTopic = 1
)

func (o PathSubscriptionTopic) DeepCopy() PathSubscriptionTopic { return o }

var PathSubscriptionTopicMap = map[string]PathSubscriptionTopic{
	"CHILDREN": 0,
	"STAT":     1,
}

var PathSubscriptionTopicRevMap = map[PathSubscriptionTopic]string{
	0: "CHILDREN",
	1: "STAT",
}

func (o PathSubscriptionTopic) String() string {
	if v, ok := PathSubscriptionTopicRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type DownloadInfo struct {
	DownloadID        string   `codec:"downloadID" json:"downloadID"`
	Path              KBFSPath `codec:"path" json:"path"`
	Filename          string   `codec:"filename" json:"filename"`
	StartTime         Time     `codec:"startTime" json:"startTime"`
	IsRegularDownload bool     `codec:"isRegularDownload" json:"isRegularDownload"`
}

func (o DownloadInfo) DeepCopy() DownloadInfo {
	return DownloadInfo{
		DownloadID:        o.DownloadID,
		Path:              o.Path.DeepCopy(),
		Filename:          o.Filename,
		StartTime:         o.StartTime.DeepCopy(),
		IsRegularDownload: o.IsRegularDownload,
	}
}

type DownloadState struct {
	DownloadID  string  `codec:"downloadID" json:"downloadID"`
	Progress    float64 `codec:"progress" json:"progress"`
	EndEstimate Time    `codec:"endEstimate" json:"endEstimate"`
	LocalPath   string  `codec:"localPath" json:"localPath"`
	Error       string  `codec:"error" json:"error"`
	Done        bool    `codec:"done" json:"done"`
	Canceled    bool    `codec:"canceled" json:"canceled"`
}

func (o DownloadState) DeepCopy() DownloadState {
	return DownloadState{
		DownloadID:  o.DownloadID,
		Progress:    o.Progress,
		EndEstimate: o.EndEstimate.DeepCopy(),
		LocalPath:   o.LocalPath,
		Error:       o.Error,
		Done:        o.Done,
		Canceled:    o.Canceled,
	}
}

type DownloadStatus struct {
	RegularDownloadIDs []string        `codec:"regularDownloadIDs" json:"regularDownloadIDs"`
	States             []DownloadState `codec:"states" json:"states"`
}

func (o DownloadStatus) DeepCopy() DownloadStatus {
	return DownloadStatus{
		RegularDownloadIDs: (func(x []string) []string {
			if x == nil {
				return nil
			}
			ret := make([]string, len(x))
			for i, v := range x {
				vCopy := v
				ret[i] = vCopy
			}
			return ret
		})(o.RegularDownloadIDs),
		States: (func(x []DownloadState) []DownloadState {
			if x == nil {
				return nil
			}
			ret := make([]DownloadState, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.States),
	}
}

type UploadState struct {
	UploadID   string   `codec:"uploadID" json:"uploadID"`
	TargetPath KBFSPath `codec:"targetPath" json:"targetPath"`
	Error      *string  `codec:"error,omitempty" json:"error,omitempty"`
	Canceled   bool     `codec:"canceled" json:"canceled"`
}

func (o UploadState) DeepCopy() UploadState {
	return UploadState{
		UploadID:   o.UploadID,
		TargetPath: o.TargetPath.DeepCopy(),
		Error: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.Error),
		Canceled: o.Canceled,
	}
}

type FilesTabBadge int

const (
	FilesTabBadge_NONE            FilesTabBadge = 0
	FilesTabBadge_UPLOADING_STUCK FilesTabBadge = 1
	FilesTabBadge_AWAITING_UPLOAD FilesTabBadge = 2
	FilesTabBadge_UPLOADING       FilesTabBadge = 3
)

func (o FilesTabBadge) DeepCopy() FilesTabBadge { return o }

var FilesTabBadgeMap = map[string]FilesTabBadge{
	"NONE":            0,
	"UPLOADING_STUCK": 1,
	"AWAITING_UPLOAD": 2,
	"UPLOADING":       3,
}

var FilesTabBadgeRevMap = map[FilesTabBadge]string{
	0: "NONE",
	1: "UPLOADING_STUCK",
	2: "AWAITING_UPLOAD",
	3: "UPLOADING",
}

func (o FilesTabBadge) String() string {
	if v, ok := FilesTabBadgeRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type GUIViewType int

const (
	GUIViewType_DEFAULT GUIViewType = 0
	GUIViewType_TEXT    GUIViewType = 1
	GUIViewType_IMAGE   GUIViewType = 2
	GUIViewType_AUDIO   GUIViewType = 3
	GUIViewType_VIDEO   GUIViewType = 4
	GUIViewType_PDF     GUIViewType = 5
)

func (o GUIViewType) DeepCopy() GUIViewType { return o }

var GUIViewTypeMap = map[string]GUIViewType{
	"DEFAULT": 0,
	"TEXT":    1,
	"IMAGE":   2,
	"AUDIO":   3,
	"VIDEO":   4,
	"PDF":     5,
}

var GUIViewTypeRevMap = map[GUIViewType]string{
	0: "DEFAULT",
	1: "TEXT",
	2: "IMAGE",
	3: "AUDIO",
	4: "VIDEO",
	5: "PDF",
}

func (o GUIViewType) String() string {
	if v, ok := GUIViewTypeRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type GUIFileContext struct {
	ViewType    GUIViewType `codec:"viewType" json:"viewType"`
	ContentType string      `codec:"contentType" json:"contentType"`
	Url         string      `codec:"url" json:"url"`
}

func (o GUIFileContext) DeepCopy() GUIFileContext {
	return GUIFileContext{
		ViewType:    o.ViewType.DeepCopy(),
		ContentType: o.ContentType,
		Url:         o.Url,
	}
}

type SimpleFSSearchHit struct {
	Path string `codec:"path" json:"path"`
}

func (o SimpleFSSearchHit) DeepCopy() SimpleFSSearchHit {
	return SimpleFSSearchHit{
		Path: o.Path,
	}
}

type SimpleFSSearchResults struct {
	Hits       []SimpleFSSearchHit `codec:"hits" json:"hits"`
	NextResult int                 `codec:"nextResult" json:"nextResult"`
}

func (o SimpleFSSearchResults) DeepCopy() SimpleFSSearchResults {
	return SimpleFSSearchResults{
		Hits: (func(x []SimpleFSSearchHit) []SimpleFSSearchHit {
			if x == nil {
				return nil
			}
			ret := make([]SimpleFSSearchHit, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.Hits),
		NextResult: o.NextResult,
	}
}

type IndexProgressRecord struct {
	EndEstimate Time  `codec:"endEstimate" json:"endEstimate"`
	BytesTotal  int64 `codec:"bytesTotal" json:"bytesTotal"`
	BytesSoFar  int64 `codec:"bytesSoFar" json:"bytesSoFar"`
}

func (o IndexProgressRecord) DeepCopy() IndexProgressRecord {
	return IndexProgressRecord{
		EndEstimate: o.EndEstimate.DeepCopy(),
		BytesTotal:  o.BytesTotal,
		BytesSoFar:  o.BytesSoFar,
	}
}

type SimpleFSIndexProgress struct {
	OverallProgress IndexProgressRecord `codec:"overallProgress" json:"overallProgress"`
	CurrFolder      Folder              `codec:"currFolder" json:"currFolder"`
	CurrProgress    IndexProgressRecord `codec:"currProgress" json:"currProgress"`
	FoldersLeft     []Folder            `codec:"foldersLeft" json:"foldersLeft"`
}

func (o SimpleFSIndexProgress) DeepCopy() SimpleFSIndexProgress {
	return SimpleFSIndexProgress{
		OverallProgress: o.OverallProgress.DeepCopy(),
		CurrFolder:      o.CurrFolder.DeepCopy(),
		CurrProgress:    o.CurrProgress.DeepCopy(),
		FoldersLeft: (func(x []Folder) []Folder {
			if x == nil {
				return nil
			}
			ret := make([]Folder, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.FoldersLeft),
	}
}

type ArchiveJobStartPathType int

const (
	ArchiveJobStartPathType_KBFS ArchiveJobStartPathType = 0
	ArchiveJobStartPathType_GIT  ArchiveJobStartPathType = 1
)

func (o ArchiveJobStartPathType) DeepCopy() ArchiveJobStartPathType { return o }

var ArchiveJobStartPathTypeMap = map[string]ArchiveJobStartPathType{
	"KBFS": 0,
	"GIT":  1,
}

var ArchiveJobStartPathTypeRevMap = map[ArchiveJobStartPathType]string{
	0: "KBFS",
	1: "GIT",
}

func (o ArchiveJobStartPathType) String() string {
	if v, ok := ArchiveJobStartPathTypeRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type ArchiveJobStartPath struct {
	ArchiveJobStartPathType__ ArchiveJobStartPathType `codec:"archiveJobStartPathType" json:"archiveJobStartPathType"`
	Kbfs__                    *KBFSPath               `codec:"kbfs,omitempty" json:"kbfs,omitempty"`
	Git__                     *string                 `codec:"git,omitempty" json:"git,omitempty"`
}

func (o *ArchiveJobStartPath) ArchiveJobStartPathType() (ret ArchiveJobStartPathType, err error) {
	switch o.ArchiveJobStartPathType__ {
	case ArchiveJobStartPathType_KBFS:
		if o.Kbfs__ == nil {
			err = errors.New("unexpected nil value for Kbfs__")
			return ret, err
		}
	case ArchiveJobStartPathType_GIT:
		if o.Git__ == nil {
			err = errors.New("unexpected nil value for Git__")
			return ret, err
		}
	}
	return o.ArchiveJobStartPathType__, nil
}

func (o ArchiveJobStartPath) Kbfs() (res KBFSPath) {
	if o.ArchiveJobStartPathType__ != ArchiveJobStartPathType_KBFS {
		panic("wrong case accessed")
	}
	if o.Kbfs__ == nil {
		return
	}
	return *o.Kbfs__
}

func (o ArchiveJobStartPath) Git() (res string) {
	if o.ArchiveJobStartPathType__ != ArchiveJobStartPathType_GIT {
		panic("wrong case accessed")
	}
	if o.Git__ == nil {
		return
	}
	return *o.Git__
}

func NewArchiveJobStartPathWithKbfs(v KBFSPath) ArchiveJobStartPath {
	return ArchiveJobStartPath{
		ArchiveJobStartPathType__: ArchiveJobStartPathType_KBFS,
		Kbfs__:                    &v,
	}
}

func NewArchiveJobStartPathWithGit(v string) ArchiveJobStartPath {
	return ArchiveJobStartPath{
		ArchiveJobStartPathType__: ArchiveJobStartPathType_GIT,
		Git__:                     &v,
	}
}

func (o ArchiveJobStartPath) DeepCopy() ArchiveJobStartPath {
	return ArchiveJobStartPath{
		ArchiveJobStartPathType__: o.ArchiveJobStartPathType__.DeepCopy(),
		Kbfs__: (func(x *KBFSPath) *KBFSPath {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.Kbfs__),
		Git__: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.Git__),
	}
}

type SimpleFSArchiveJobDesc struct {
	JobID                string           `codec:"jobID" json:"jobID"`
	KbfsPathWithRevision KBFSArchivedPath `codec:"kbfsPathWithRevision" json:"kbfsPathWithRevision"`
	GitRepo              *string          `codec:"gitRepo,omitempty" json:"gitRepo,omitempty"`
	OverwriteZip         bool             `codec:"overwriteZip" json:"overwriteZip"`
	StartTime            Time             `codec:"startTime" json:"startTime"`
	StagingPath          string           `codec:"stagingPath" json:"stagingPath"`
	TargetName           string           `codec:"targetName" json:"targetName"`
	ZipFilePath          string           `codec:"zipFilePath" json:"zipFilePath"`
}

func (o SimpleFSArchiveJobDesc) DeepCopy() SimpleFSArchiveJobDesc {
	return SimpleFSArchiveJobDesc{
		JobID:                o.JobID,
		KbfsPathWithRevision: o.KbfsPathWithRevision.DeepCopy(),
		GitRepo: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.GitRepo),
		OverwriteZip: o.OverwriteZip,
		StartTime:    o.StartTime.DeepCopy(),
		StagingPath:  o.StagingPath,
		TargetName:   o.TargetName,
		ZipFilePath:  o.ZipFilePath,
	}
}

type SimpleFSFileArchiveState int

const (
	SimpleFSFileArchiveState_ToDo       SimpleFSFileArchiveState = 0
	SimpleFSFileArchiveState_InProgress SimpleFSFileArchiveState = 1
	SimpleFSFileArchiveState_Complete   SimpleFSFileArchiveState = 2
	SimpleFSFileArchiveState_Skipped    SimpleFSFileArchiveState = 3
)

func (o SimpleFSFileArchiveState) DeepCopy() SimpleFSFileArchiveState { return o }

var SimpleFSFileArchiveStateMap = map[string]SimpleFSFileArchiveState{
	"ToDo":       0,
	"InProgress": 1,
	"Complete":   2,
	"Skipped":    3,
}

var SimpleFSFileArchiveStateRevMap = map[SimpleFSFileArchiveState]string{
	0: "ToDo",
	1: "InProgress",
	2: "Complete",
	3: "Skipped",
}

func (o SimpleFSFileArchiveState) String() string {
	if v, ok := SimpleFSFileArchiveStateRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type SimpleFSArchiveFile struct {
	State        SimpleFSFileArchiveState `codec:"state" json:"state"`
	DirentType   DirentType               `codec:"direntType" json:"direntType"`
	Sha256SumHex string                   `codec:"sha256SumHex" json:"sha256SumHex"`
}

func (o SimpleFSArchiveFile) DeepCopy() SimpleFSArchiveFile {
	return SimpleFSArchiveFile{
		State:        o.State.DeepCopy(),
		DirentType:   o.DirentType.DeepCopy(),
		Sha256SumHex: o.Sha256SumHex,
	}
}

type SimpleFSArchiveJobState struct {
	Desc        SimpleFSArchiveJobDesc         `codec:"desc" json:"desc"`
	Manifest    map[string]SimpleFSArchiveFile `codec:"manifest" json:"manifest"`
	Phase       SimpleFSArchiveJobPhase        `codec:"phase" json:"phase"`
	BytesTotal  int64                          `codec:"bytesTotal" json:"bytesTotal"`
	BytesCopied int64                          `codec:"bytesCopied" json:"bytesCopied"`
	BytesZipped int64                          `codec:"bytesZipped" json:"bytesZipped"`
}

func (o SimpleFSArchiveJobState) DeepCopy() SimpleFSArchiveJobState {
	return SimpleFSArchiveJobState{
		Desc: o.Desc.DeepCopy(),
		Manifest: (func(x map[string]SimpleFSArchiveFile) map[string]SimpleFSArchiveFile {
			if x == nil {
				return nil
			}
			ret := make(map[string]SimpleFSArchiveFile, len(x))
			for k, v := range x {
				kCopy := k
				vCopy := v.DeepCopy()
				ret[kCopy] = vCopy
			}
			return ret
		})(o.Manifest),
		Phase:       o.Phase.DeepCopy(),
		BytesTotal:  o.BytesTotal,
		BytesCopied: o.BytesCopied,
		BytesZipped: o.BytesZipped,
	}
}

type SimpleFSArchiveJobPhase int

const (
	SimpleFSArchiveJobPhase_Queued   SimpleFSArchiveJobPhase = 0
	SimpleFSArchiveJobPhase_Indexing SimpleFSArchiveJobPhase = 1
	SimpleFSArchiveJobPhase_Indexed  SimpleFSArchiveJobPhase = 2
	SimpleFSArchiveJobPhase_Copying  SimpleFSArchiveJobPhase = 3
	SimpleFSArchiveJobPhase_Copied   SimpleFSArchiveJobPhase = 4
	SimpleFSArchiveJobPhase_Zipping  SimpleFSArchiveJobPhase = 5
	SimpleFSArchiveJobPhase_Done     SimpleFSArchiveJobPhase = 6
)

func (o SimpleFSArchiveJobPhase) DeepCopy() SimpleFSArchiveJobPhase { return o }

var SimpleFSArchiveJobPhaseMap = map[string]SimpleFSArchiveJobPhase{
	"Queued":   0,
	"Indexing": 1,
	"Indexed":  2,
	"Copying":  3,
	"Copied":   4,
	"Zipping":  5,
	"Done":     6,
}

var SimpleFSArchiveJobPhaseRevMap = map[SimpleFSArchiveJobPhase]string{
	0: "Queued",
	1: "Indexing",
	2: "Indexed",
	3: "Copying",
	4: "Copied",
	5: "Zipping",
	6: "Done",
}

func (o SimpleFSArchiveJobPhase) String() string {
	if v, ok := SimpleFSArchiveJobPhaseRevMap[o]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(o))
}

type SimpleFSArchiveState struct {
	Jobs        map[string]SimpleFSArchiveJobState `codec:"jobs" json:"jobs"`
	LastUpdated Time                               `codec:"lastUpdated" json:"lastUpdated"`
}

func (o SimpleFSArchiveState) DeepCopy() SimpleFSArchiveState {
	return SimpleFSArchiveState{
		Jobs: (func(x map[string]SimpleFSArchiveJobState) map[string]SimpleFSArchiveJobState {
			if x == nil {
				return nil
			}
			ret := make(map[string]SimpleFSArchiveJobState, len(x))
			for k, v := range x {
				kCopy := k
				vCopy := v.DeepCopy()
				ret[kCopy] = vCopy
			}
			return ret
		})(o.Jobs),
		LastUpdated: o.LastUpdated.DeepCopy(),
	}
}

type SimpleFSArchiveJobErrorState struct {
	Error     string `codec:"error" json:"error"`
	NextRetry Time   `codec:"nextRetry" json:"nextRetry"`
}

func (o SimpleFSArchiveJobErrorState) DeepCopy() SimpleFSArchiveJobErrorState {
	return SimpleFSArchiveJobErrorState{
		Error:     o.Error,
		NextRetry: o.NextRetry.DeepCopy(),
	}
}

type SimpleFSArchiveJobStatus struct {
	Desc            SimpleFSArchiveJobDesc        `codec:"desc" json:"desc"`
	Phase           SimpleFSArchiveJobPhase       `codec:"phase" json:"phase"`
	TodoCount       int                           `codec:"todoCount" json:"todoCount"`
	InProgressCount int                           `codec:"inProgressCount" json:"inProgressCount"`
	CompleteCount   int                           `codec:"completeCount" json:"completeCount"`
	SkippedCount    int                           `codec:"skippedCount" json:"skippedCount"`
	TotalCount      int                           `codec:"totalCount" json:"totalCount"`
	BytesTotal      int64                         `codec:"bytesTotal" json:"bytesTotal"`
	BytesCopied     int64                         `codec:"bytesCopied" json:"bytesCopied"`
	BytesZipped     int64                         `codec:"bytesZipped" json:"bytesZipped"`
	Error           *SimpleFSArchiveJobErrorState `codec:"error,omitempty" json:"error,omitempty"`
}

func (o SimpleFSArchiveJobStatus) DeepCopy() SimpleFSArchiveJobStatus {
	return SimpleFSArchiveJobStatus{
		Desc:            o.Desc.DeepCopy(),
		Phase:           o.Phase.DeepCopy(),
		TodoCount:       o.TodoCount,
		InProgressCount: o.InProgressCount,
		CompleteCount:   o.CompleteCount,
		SkippedCount:    o.SkippedCount,
		TotalCount:      o.TotalCount,
		BytesTotal:      o.BytesTotal,
		BytesCopied:     o.BytesCopied,
		BytesZipped:     o.BytesZipped,
		Error: (func(x *SimpleFSArchiveJobErrorState) *SimpleFSArchiveJobErrorState {
			if x == nil {
				return nil
			}
			tmp := x.DeepCopy()
			return &tmp
		})(o.Error),
	}
}

type SimpleFSArchiveStatus struct {
	Jobs        []SimpleFSArchiveJobStatus `codec:"jobs" json:"jobs"`
	LastUpdated Time                       `codec:"lastUpdated" json:"lastUpdated"`
}

func (o SimpleFSArchiveStatus) DeepCopy() SimpleFSArchiveStatus {
	return SimpleFSArchiveStatus{
		Jobs: (func(x []SimpleFSArchiveJobStatus) []SimpleFSArchiveJobStatus {
			if x == nil {
				return nil
			}
			ret := make([]SimpleFSArchiveJobStatus, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.Jobs),
		LastUpdated: o.LastUpdated.DeepCopy(),
	}
}

type SimpleFSArchiveJobFreshness struct {
	CurrentTLFRevision KBFSRevision `codec:"currentTLFRevision" json:"currentTLFRevision"`
}

func (o SimpleFSArchiveJobFreshness) DeepCopy() SimpleFSArchiveJobFreshness {
	return SimpleFSArchiveJobFreshness{
		CurrentTLFRevision: o.CurrentTLFRevision.DeepCopy(),
	}
}

type SimpleFSArchiveCheckArchiveResult struct {
	Desc               SimpleFSArchiveJobDesc `codec:"desc" json:"desc"`
	CurrentTLFRevision KBFSRevision           `codec:"currentTLFRevision" json:"currentTLFRevision"`
	PathsWithIssues    map[string]string      `codec:"pathsWithIssues" json:"pathsWithIssues"`
}

func (o SimpleFSArchiveCheckArchiveResult) DeepCopy() SimpleFSArchiveCheckArchiveResult {
	return SimpleFSArchiveCheckArchiveResult{
		Desc:               o.Desc.DeepCopy(),
		CurrentTLFRevision: o.CurrentTLFRevision.DeepCopy(),
		PathsWithIssues: (func(x map[string]string) map[string]string {
			if x == nil {
				return nil
			}
			ret := make(map[string]string, len(x))
			for k, v := range x {
				kCopy := k
				vCopy := v
				ret[kCopy] = vCopy
			}
			return ret
		})(o.PathsWithIssues),
	}
}

type SimpleFSArchiveAllFilesResult struct {
	TlfPathToJobDesc map[string]SimpleFSArchiveJobDesc `codec:"tlfPathToJobDesc" json:"tlfPathToJobDesc"`
	TlfPathToError   map[string]string                 `codec:"tlfPathToError" json:"tlfPathToError"`
	SkippedTLFPaths  []string                          `codec:"skippedTLFPaths" json:"skippedTLFPaths"`
}

func (o SimpleFSArchiveAllFilesResult) DeepCopy() SimpleFSArchiveAllFilesResult {
	return SimpleFSArchiveAllFilesResult{
		TlfPathToJobDesc: (func(x map[string]SimpleFSArchiveJobDesc) map[string]SimpleFSArchiveJobDesc {
			if x == nil {
				return nil
			}
			ret := make(map[string]SimpleFSArchiveJobDesc, len(x))
			for k, v := range x {
				kCopy := k
				vCopy := v.DeepCopy()
				ret[kCopy] = vCopy
			}
			return ret
		})(o.TlfPathToJobDesc),
		TlfPathToError: (func(x map[string]string) map[string]string {
			if x == nil {
				return nil
			}
			ret := make(map[string]string, len(x))
			for k, v := range x {
				kCopy := k
				vCopy := v
				ret[kCopy] = vCopy
			}
			return ret
		})(o.TlfPathToError),
		SkippedTLFPaths: (func(x []string) []string {
			if x == nil {
				return nil
			}
			ret := make([]string, len(x))
			for i, v := range x {
				vCopy := v
				ret[i] = vCopy
			}
			return ret
		})(o.SkippedTLFPaths),
	}
}

type SimpleFSArchiveAllGitReposResult struct {
	GitRepoToJobDesc map[string]SimpleFSArchiveJobDesc `codec:"gitRepoToJobDesc" json:"gitRepoToJobDesc"`
	GitRepoToError   map[string]string                 `codec:"gitRepoToError" json:"gitRepoToError"`
}

func (o SimpleFSArchiveAllGitReposResult) DeepCopy() SimpleFSArchiveAllGitReposResult {
	return SimpleFSArchiveAllGitReposResult{
		GitRepoToJobDesc: (func(x map[string]SimpleFSArchiveJobDesc) map[string]SimpleFSArchiveJobDesc {
			if x == nil {
				return nil
			}
			ret := make(map[string]SimpleFSArchiveJobDesc, len(x))
			for k, v := range x {
				kCopy := k
				vCopy := v.DeepCopy()
				ret[kCopy] = vCopy
			}
			return ret
		})(o.GitRepoToJobDesc),
		GitRepoToError: (func(x map[string]string) map[string]string {
			if x == nil {
				return nil
			}
			ret := make(map[string]string, len(x))
			for k, v := range x {
				kCopy := k
				vCopy := v
				ret[kCopy] = vCopy
			}
			return ret
		})(o.GitRepoToError),
	}
}

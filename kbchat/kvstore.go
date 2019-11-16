package kbchat

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat/types/keybase1"
)

type GetEntry struct {
	Result keybase1.KVGetResult `json:"result"`
	Error  Error                `json:"error,omitempty"`
}

type PutEntry struct {
	Result keybase1.KVPutResult `json:"result"`
	Error  Error                `json:"error,omitempty"`
}

type DeleteEntry struct {
	Result keybase1.KVDeleteEntryResult `json:"result"`
	Error  Error                        `json:"error,omitempty"`
}

type ListNamespaces struct {
	Result keybase1.KVListNamespaceResult `json:"result"`
	Error  Error                          `json:"error,omitempty"`
}

type ListEntryKeys struct {
	Result keybase1.KVListEntryResult `json:"result"`
	Error  Error                      `json:"error,omitempty"`
}

type KVStoreAPI interface {
	PutEntry(teamName string, namespace string, entryKey string, entryValue string, revision int) (keybase1.KVPutResult, error)
	DeleteEntry(teamName string, namespace string, entryKey string, revision int) (keybase1.KVDeleteEntryResult, error)
	GetEntry(teamName string, namespace string, entryKey string) (keybase1.KVGetResult, error)
	ListNamespaces(teamName string) (keybase1.KVListNamespaceResult, error)
	ListEntryKeys(teamName string, namespace string) (keybase1.KVListEntryResult, error)
}

func (a *API) PutEntry(teamName string, namespace string, entryKey string, entryValue string, revision int) (result keybase1.KVPutResult, err error) {

	type PutArgs struct {
		Method string `json:"method"`
		Params struct {
			Options struct {
				Team       string `json:"team"`
				Namespace  string `json:"namespace"`
				EntryKey   string `json:"entryKey"`
				EntryValue string `json:"entryValue"`
				Revision   int    `json:"revision,omitempty"`
			} `json:"options"`
		} `json:"params"`
	}

	args := PutArgs{Method: "put"}
	args.Params.Options.Team = teamName
	args.Params.Options.Namespace = namespace
	args.Params.Options.EntryKey = entryKey
	args.Params.Options.EntryValue = entryValue

	if revision != 0 {
		args.Params.Options.Revision = revision
	}

	apiInput, err := json.Marshal(args)
	if err != nil {
		return result, err
	}

	cmd := a.runOpts.Command("kvstore", "api")
	cmd.Stdin = strings.NewReader(string(apiInput))
	bytes, err := cmd.Output()
	if err != nil {
		return result, fmt.Errorf("failed to call keybase kvstore api: %v", err)
	}

	entry := PutEntry{}
	err = json.Unmarshal(bytes, &entry)
	if err != nil {
		return result, fmt.Errorf("failed to parse output from keybase kvstore api: %v", err)
	}
	if entry.Error.Message != "" {
		return result, fmt.Errorf("received error from keybase kvstore api: %s", entry.Error.Message)
	}
	return entry.Result, nil
}

func (a *API) DeleteEntry(teamName string, namespace string, entryKey string, revision int) (result keybase1.KVDeleteEntryResult, err error) {

	type DeleteArgs struct {
		Method string `json:"method"`
		Params struct {
			Options struct {
				Team      string `json:"team"`
				Namespace string `json:"namespace"`
				EntryKey  string `json:"entryKey"`
				Revision  int    `json:"revision,omitempty"`
			} `json:"options"`
		} `json:"params"`
	}

	args := DeleteArgs{Method: "del"}
	args.Params.Options.Team = teamName
	args.Params.Options.Namespace = namespace
	args.Params.Options.EntryKey = entryKey

	if revision != 0 {
		args.Params.Options.Revision = revision
	}

	apiInput, err := json.Marshal(args)
	if err != nil {
		return result, err
	}

	cmd := a.runOpts.Command("kvstore", "api")
	cmd.Stdin = strings.NewReader(string(apiInput))
	bytes, err := cmd.Output()
	if err != nil {
		return result, fmt.Errorf("failed to call keybase kvstore api: %v", err)
	}

	entry := DeleteEntry{}
	err = json.Unmarshal(bytes, &entry)
	if err != nil {
		return result, fmt.Errorf("failed to parse output from keybase kvstore api: %v", err)
	}
	if entry.Error.Message != "" {
		return result, fmt.Errorf("received error from keybase kvstore api: %s", entry.Error.Message)
	}
	return entry.Result, nil
}

func (a *API) GetEntry(teamName string, namespace string, entryKey string) (result keybase1.KVGetResult, err error) {

	type GetArgs struct {
		Method string `json:"method"`
		Params struct {
			Options struct {
				Team      string `json:"team"`
				Namespace string `json:"namespace"`
				EntryKey  string `json:"entryKey"`
			} `json:"options"`
		} `json:"params"`
	}

	args := GetArgs{Method: "get"}
	args.Params.Options.Team = teamName
	args.Params.Options.Namespace = namespace
	args.Params.Options.EntryKey = entryKey

	apiInput, err := json.Marshal(args)
	if err != nil {
		return result, err
	}

	cmd := a.runOpts.Command("kvstore", "api")
	cmd.Stdin = strings.NewReader(string(apiInput))
	bytes, err := cmd.Output()
	if err != nil {
		return result, fmt.Errorf("failed to call keybase kvstore api: %v", err)
	}

	entry := GetEntry{}
	err = json.Unmarshal(bytes, &entry)

	if err != nil {
		return result, fmt.Errorf("failed to parse output from keybase kvstore api: %v", err)
	}
	if entry.Error.Message != "" {
		return result, fmt.Errorf("received error from keybase kvstore api: %s", entry.Error.Message)
	}
	return entry.Result, nil
}

func (a *API) ListNamespaces(teamName string) (result keybase1.KVListNamespaceResult, err error) {
	apiInput := fmt.Sprintf(`{"method": "list", "params": {"options": {"team": "%s"}}}`, teamName)
	cmd := a.runOpts.Command("kvstore", "api")
	cmd.Stdin = strings.NewReader(apiInput)
	bytes, err := cmd.Output()
	if err != nil {
		return result, fmt.Errorf("failed to call keybase kvstore api: %v", err)
	}

	var namespaces ListNamespaces
	err = json.Unmarshal(bytes, &namespaces)
	if err != nil {
		return result, fmt.Errorf("failed to parse output from keybase kvstore api: %v", err)
	}
	if namespaces.Error.Message != "" {
		return result, fmt.Errorf("received error from keybase kvstore api: %s", namespaces.Error.Message)
	}
	return namespaces.Result, nil
}

func (a *API) ListEntryKeys(teamName string, namespace string) (result keybase1.KVListEntryResult, err error) {
	apiInput := fmt.Sprintf(`{"method": "list", "params": {"options": {"team": "%s", "namespace": "%s"}}}`, teamName, namespace)
	cmd := a.runOpts.Command("kvstore", "api")
	cmd.Stdin = strings.NewReader(apiInput)
	bytes, err := cmd.Output()
	if err != nil {
		return result, fmt.Errorf("failed to call keybase kvstore api: %v", err)
	}

	entryKeys := ListEntryKeys{}
	err = json.Unmarshal(bytes, &entryKeys)
	if err != nil {
		return result, fmt.Errorf("failed to parse output from keybase kvstore api: %v", err)
	}
	if entryKeys.Error.Message != "" {
		return result, fmt.Errorf("received error from keybase kvstore api: %s", entryKeys.Error.Message)
	}
	return entryKeys.Result, nil
}

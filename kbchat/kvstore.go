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

func (a *API) PutEntry(teamName string, namespace string, entryKey string, entryValue string, revision int) (keybase1.KVPutResult, error) {
	empty := keybase1.KVPutResult{}

	apiInput := fmt.Sprintf(`{"method": "put", "params": {"options": {"team": "%s", "namespace": "%s", "entryKey": "%s", "entryValue": "%s"}}}`, teamName, namespace, entryKey, entryValue)
	if revision != 0 {
		apiInput = fmt.Sprintf(`{"method": "put", "params": {"options": {"team": "%s", "namespace": "%s", "entryKey": "%s", "entryValue": "%s", "revision": %d}}}`, teamName, namespace, entryKey, entryValue, revision)
	}

	cmd := a.runOpts.Command("kvstore", "api")
	cmd.Stdin = strings.NewReader(apiInput)
	bytes, err := cmd.Output()
	if err != nil {
		return empty, fmt.Errorf("failed to call keybase kvstore api: %v", err)
	}

	entry := PutEntry{}
	err = json.Unmarshal(bytes, &entry)
	if err != nil {
		return empty, fmt.Errorf("failed to parse output from keybase kvstore api: %v", err)
	}
	if entry.Error.Message != "" {
		return empty, fmt.Errorf("received error from keybase kvstore api: %s", entry.Error.Message)
	}
	return entry.Result, nil
}

func (a *API) DeleteEntry(teamName string, namespace string, entryKey string, revision int) (keybase1.KVDeleteEntryResult, error) {
	empty := keybase1.KVDeleteEntryResult{}

	apiInput := fmt.Sprintf(`{"method": "delete", "params": {"options": {"team": "%s", "namespace": "%s", "entryKey": "%s"}}}`, teamName, namespace, entryKey)
	if revision != 0 {
		apiInput = fmt.Sprintf(`{"method": "delete", "params": {"options": {"team": "%s", "namespace": "%s", "entryKey": "%s", "revision": "%d"}}}`, teamName, namespace, entryKey, revision)
	}

	cmd := a.runOpts.Command("kvstore", "api")
	cmd.Stdin = strings.NewReader(apiInput)
	bytes, err := cmd.Output()
	if err != nil {
		return empty, fmt.Errorf("failed to call keybase kvstore api: %v", err)
	}

	entry := DeleteEntry{}
	err = json.Unmarshal(bytes, &entry)
	if err != nil {
		return empty, fmt.Errorf("failed to parse output from keybase kvstore api: %v", err)
	}
	if entry.Error.Message != "" {
		return empty, fmt.Errorf("received error from keybase kvstore api: %s", entry.Error.Message)
	}
	return entry.Result, nil
}

func (a *API) GetEntry(teamName string, namespace string, entryKey string) (keybase1.KVGetResult, error) {
	empty := keybase1.KVGetResult{}

	apiInput := fmt.Sprintf(`{"method": "get", "params": {"options": {"team": "%s", "namespace": "%s", "entryKey": "%s"}}}`, teamName, namespace, entryKey)
	cmd := a.runOpts.Command("kvstore", "api")
	cmd.Stdin = strings.NewReader(apiInput)
	bytes, err := cmd.Output()
	if err != nil {
		return empty, fmt.Errorf("failed to call keybase kvstore api: %v", err)
	}

	entry := GetEntry{}
	err = json.Unmarshal(bytes, &entry)
	if err != nil {
		return empty, fmt.Errorf("failed to parse output from keybase kvstore api: %v", err)
	}
	if entry.Error.Message != "" {
		return empty, fmt.Errorf("received error from keybase kvstore api: %s", entry.Error.Message)
	}
	return entry.Result, nil
}

func (a *API) ListNamespaces(teamName string) ([]string, error) {
	empty := []string{}

	apiInput := fmt.Sprintf(`{"method": "list", "params": {"options": {"team": "%s"}}}`, teamName)
	cmd := a.runOpts.Command("kvstore", "api")
	cmd.Stdin = strings.NewReader(apiInput)
	bytes, err := cmd.Output()
	if err != nil {
		return empty, fmt.Errorf("failed to call keybase kvstore api: %v", err)
	}

	var namespaces ListNamespaces
	err = json.Unmarshal(bytes, &namespaces)
	if err != nil {
		return empty, fmt.Errorf("failed to parse output from keybase kvstore api: %v", err)
	}
	if namespaces.Error.Message != "" {
		return empty, fmt.Errorf("received error from keybase kvstore api: %s", namespaces.Error.Message)
	}
	return namespaces.Result.Namespaces, nil
}

func (a *API) ListEntryKeys(teamName string, namespace string) ([]keybase1.KVListEntryKey, error) {
	empty := []keybase1.KVListEntryKey{}

	apiInput := fmt.Sprintf(`{"method": "list", "params": {"options": {"team": "%s", "namespace": "%s"}}}`, teamName, namespace)
	cmd := a.runOpts.Command("kvstore", "api")
	cmd.Stdin = strings.NewReader(apiInput)
	bytes, err := cmd.Output()
	if err != nil {
		return empty, fmt.Errorf("failed to call keybase kvstore api: %v", err)
	}

	entryKeys := ListEntryKeys{}
	err = json.Unmarshal(bytes, &entryKeys)
	if err != nil {
		return empty, fmt.Errorf("failed to parse output from keybase kvstore api: %v", err)
	}
	if entryKeys.Error.Message != "" {
		return empty, fmt.Errorf("received error from keybase kvstore api: %s", entryKeys.Error.Message)
	}
	return entryKeys.Result.EntryKeys, nil
}

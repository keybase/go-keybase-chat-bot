package kbchat

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"testing"
)

func randomTempDir() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	dir := path.Join(os.TempDir(), "keybase_bot_"+hex.EncodeToString(bytes))
	return dir, nil
}

func TestListChannels(t *testing.T) {
	dir, _ := randomTempDir()
	fmt.Printf("randomTempDir() = %+v\n", dir)
}

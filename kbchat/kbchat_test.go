package kbchat

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type Bots struct {
	Alice1   *OneshotOptions
	Alice2   *OneshotOptions
	Bob1     *OneshotOptions
	Charlie1 *OneshotOptions
}

type Config struct {
	Bots
}

func readAndParseConfig() (Config, error) {
	var config Config
	data, err := ioutil.ReadFile("test_config.yaml")
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func randomTempDir() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	dir := path.Join(os.TempDir(), "keybase_bot_"+hex.EncodeToString(bytes))
	return dir, nil
}

func whichKeybase() (string, error) {
	cmd := exec.Command("which", "keybase")

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	location := strings.TrimSpace(string(out))
	return location, nil
}

func copyFile(source, dest string) error {
	sourceData, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, sourceData, 0777)
	if err != nil {
		return err
	}

	return nil
}

func prepWorkingDir(workingDir string) (string, error) {
	kbLocation, err := whichKeybase()
	if err != nil {
		return "", err
	}

	err = os.Mkdir(workingDir, 0777)
	if err != nil {
		return "", err
	}
	kbDestination := path.Join(workingDir, "keybase")

	err = copyFile(kbLocation, kbDestination)
	if err != nil {
		return "", err
	}

	return kbDestination, nil
}

func TestSendMessageByTlfName(t *testing.T) {
	config, err := readAndParseConfig()
	require.NoError(t, err)
	dir, err := randomTempDir()
	require.NoError(t, err)
	kbLocation, err := prepWorkingDir(dir)
	require.NoError(t, err)

	kbc, err := Start(RunOptions{KeybaseLocation: kbLocation, HomeDir: dir, Oneshot: config.Bots.Alice1, StartService: true})
	require.NoError(t, err, "error %s")

	tlfName := fmt.Sprintf("%s,%s", kbc.Username(), "kb_monbot")
	res, err := kbc.SendMessageByTlfName(tlfName, "test")
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	err = kbc.Shutdown()
	require.NoError(t, err)
}

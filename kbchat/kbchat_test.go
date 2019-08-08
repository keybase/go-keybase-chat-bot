package kbchat

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
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

// Creates the working directory and copies over the keybase binary in PATH.
// We do this to avoid any version mismatch issues.
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

func deleteWorkingDir(workingDir string) error {
	return os.RemoveAll(workingDir)
}

var kbc *API
var config Config

func TestMain(m *testing.M) {
	var err error
	config, err = readAndParseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in reading config: %v\n", err)
	}
	dir, err := randomTempDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in generating directory: %v\n", err)
	}
	kbLocation, err := prepWorkingDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in preparing working directory: %v\n", err)
	}
	kbc, err = Start(RunOptions{KeybaseLocation: kbLocation, HomeDir: dir, Oneshot: config.Bots.Alice1, StartService: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in starting service: %v\n", err)
	}

	flag.Parse()
	code := m.Run()

	err = kbc.Shutdown()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on service shutdown: %v\n", err)
	}
	err = deleteWorkingDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error working directory teardown: %v\n", err)
	}
	os.Exit(code)
}

func TestGetUsername(t *testing.T) {
	require.Equal(t, kbc.GetUsername(), config.Bots.Alice1.Username)
}

func TestGetConversations(t *testing.T) {
	conversations, err := kbc.GetConversations(false)
	require.NoError(t, err)
	require.Greater(t, len(conversations), 0)
}

func TestGetTextMessages(t *testing.T) {
	channel := Channel{
		Name: fmt.Sprintf("%s,%s", config.Bots.Alice1.Username, config.Bots.Charlie1.Username),
	}
	messages, err := kbc.GetTextMessages(channel, false)
	require.NoError(t, err)
	require.Greater(t, len(messages), 0)
}

func TestSendMessage(t *testing.T) {}

func TestSendMessageByConvID(t *testing.T) {}

func TestSendMessageByTlfName(t *testing.T) {
	tlfName := fmt.Sprintf("%s,%s", kbc.Username(), "kb_monbot")
	res, err := kbc.SendMessageByTlfName(tlfName, "test")
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)
}

func TestSendMessageByTeamName(t *testing.T) {}

func TestSendAttachmentByTeam(t *testing.T) {}

func TestReactByChannel(t *testing.T) {}

func TestReactByConvID(t *testing.T) {}

func TestAdvertiseCommands(t *testing.T) {}

func TestListChannels(t *testing.T) {}

func TestJoinChannel(t *testing.T) {}

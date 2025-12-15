package swiftx

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/hiscaler/swiftx-go/config"
)

var client *Client

func TestMain(m *testing.M) {
	b, err := os.ReadFile("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var cfg config.Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	client = NewClient(cfg)
	m.Run()
}

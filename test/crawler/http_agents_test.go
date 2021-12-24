package crawler

import (
	"fmt"
	"testing"
	"github.com/laconlab/lacon-go-tiny-scrapy/src/crawler"
	"gopkg.in/yaml.v2"
)

func TestHttpAgentYaml(t *testing.T) {
    config := `
    User-Agents:
        - Test1
        - Test-2
        - Test_3
    `

    agents := &crawler.HttpAgents{}

    if err := yaml.Unmarshal([]byte(config), &agents); err != nil {
        t.Error(err)
    }

    fmt.Println(agents)

    if agents.Next() != "Test1" {
        t.Error("Failed Test1")
    }

    if agents.Next() != "Test-2" {
        t.Error("Failed Test1")
    }

    if agents.Next() != "Test_3" {
        t.Error("Failed Test1")
    }

    if agents.Next() != "Test1" {
        t.Error("Failed Test1")
    }
}

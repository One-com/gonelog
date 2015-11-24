package gonelog_test

import (
	"github.com/One-com/gonelog/log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Hello world")
	os.Exit(m.Run())
}

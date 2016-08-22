package utilities

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	// arrange
	configUtil := NewConfigUtil()
	// act
	dbName := configUtil.GetConfig("dbName")
	// assert
	if dbName != "test" {
		t.Error("Error db name")
	}
}
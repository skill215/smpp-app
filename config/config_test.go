package config_test

import (
	"testing"

	"github.com/skill215/smpp-app/config"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	conf, err := config.GetSmppConf("smpp-app.yaml")
	assert.Nil(t, err)
	assert.True(t, len(conf.App.SmppConn) > 0)
}

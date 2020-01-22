package config

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"knative.dev/test-infra/tools/wathola/internal/ensure"
	"os"
	"path"
	"testing"
)

var id = uuid.New()

func TestReadIfPresent(t *testing.T) {
	// given
	expanded := ensureConfigFileNotPresent()
	data := []byte(`[sender]
address = 'http://default-broker.event-example.svc.cluster.local/'
`)
	ensure.NoError(ioutil.WriteFile(expanded, data, 0644))
	defer func() { ensure.NoError(os.Remove(expanded)) }()

	// when
	ReadIfPresent()

	// then
	assert.Equal(t,
		"http://default-broker.event-example.svc.cluster.local/",
		Instance.Sender.Address)
}

func TestReadIfPresentAndInvalid(t *testing.T) {
	// given
	origLogFatal := logFatal
	defer func() { logFatal = origLogFatal } ()
	expanded := ensureConfigFileNotPresent()
	data := []byte(`[sender]
address = 'http://default-broker.event-example.svc.cluster.local/
`)
	ensure.NoError(ioutil.WriteFile(expanded, data, 0644))
	defer func() { ensure.NoError(os.Remove(expanded)) }()
	var errors []string
	logFatal = func(args ...interface{}) {
		errors = append(errors, fmt.Sprint(args))
	}

	// when
	ReadIfPresent()

	// then
	assert.Contains(t, errors, "[(2, 12): unclosed string]")
}

func TestReadIfNotPresent(t *testing.T) {
	// given
	ensureConfigFileNotPresent()

	// when
	ReadIfPresent()

	// then
	assert.Equal(t,
		"http://localhost:22110/",
		Instance.Sender.Address)
}

func ensureConfigFileNotPresent() string {
	Instance = defaultValues()
	location = fmt.Sprintf("~/tmp/wathola-%v/config.toml", id.String())
	expanded, err := homedir.Expand(location)
	ensure.NoError(err)
	dir := path.Dir(expanded)
	ensure.NoError(os.MkdirAll(dir, os.ModePerm))
	if _, err := os.Stat(expanded); err == nil {
		ensure.NoError(os.Remove(expanded))
	}

	return expanded
}

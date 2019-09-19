// Package util supports various needs for running tests
package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	// ArtifactsDir is the dir containing artifacts
	ArtifactsDir = "artifacts"
	Filename     = "meta.json"
)

// client holds metadata as a string:string map, as well as path for storing
// metadata
type client struct {
	MetaData map[string]string
	path     string
}

// getLocalArtifactsDir gets the artifacts directory where prow looks for artifacts.
// By default, it will look at the env var ARTIFACTS.
// This is copy pasted from test-infra/shared/prow
// TODO(chaodaiG): refactor this so that they can be shared
func getLocalArtifactsDir() string {
	dir := os.Getenv("ARTIFACTS")
	if dir == "" {
		log.Printf("Env variable ARTIFACTS not set. Using %s instead.", ArtifactsDir)
		dir = ArtifactsDir
	}
	return dir
}

// NewClient creates a client, takes custom directory for storing `metadata.json`.
// It reads existing `meta.json` file if it exists, otherwise creates it.
// Errors out if there is any file i/o problem other than file not exist error.
func NewClient(dir string) (*client, error) {
	c := &client{
		MetaData: make(map[string]string),
	}
	if dir == "" {
		dir = getLocalArtifactsDir()
	}
	c.path = path.Join(dir, Filename)
	_, err := os.Stat(dir)
	if err == nil || os.IsNotExist(err) {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0777)
		}
		if err == nil {
			err = c.sync()
		}
	}

	if err != nil {
		return nil, err
	}
	return c, nil
}

// sync reads from meta file and convert it to Meta, returns empty
// Meta if file doesn't exist yet, and returns error if there is any i/o or
// unmarshall error
func (c *client) sync() error {
	_, err := os.Stat(c.path)
	if os.IsNotExist(err) {
		log.Println("write file")
		body, _ := json.Marshal(&c.MetaData)
		err = ioutil.WriteFile(c.path, body, 0777)
	} else {
		var body []byte
		body, err = ioutil.ReadFile(c.path)
		if err == nil {
			err = json.Unmarshal(body, &c.MetaData)
		}
	}

	return err
}

// Set sets key:val pair, and overrides if it exists
func (c *client) Set(key, val string) error {
	err := c.sync()
	if err != nil {
		return err
	}
	if oldVal, ok := c.MetaData[key]; ok {
		log.Printf("Overriding meta %q:%q with new value %q", key, oldVal, val)
	}
	c.MetaData[key] = val
	body, _ := json.Marshal(c.MetaData)
	return ioutil.WriteFile(c.path, body, 0777)
}

// Get gets val for key
func (c *client) Get(key string) (string, error) {
	var res string
	err := c.sync()
	if err == nil {
		if val, ok := c.MetaData[key]; ok {
			res = val
		} else {
			err = fmt.Errorf("key %q doesn't exist", key)
		}
	}
	return res, err
}

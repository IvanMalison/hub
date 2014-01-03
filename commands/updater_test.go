package commands

import (
	"fmt"
	"github.com/bmizerany/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUpdater_downloadFile(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/gh.zip", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, "1234")
	})

	path, err := downloadFile(fmt.Sprintf("%s/gh.zip", server.URL))
	assert.Equal(t, nil, err)

	content, err := ioutil.ReadFile(path)
	assert.Equal(t, nil, err)
	assert.Equal(t, "1234", string(content))
	assert.Equal(t, "gh.zip", filepath.Base(path))
}

func TestUpdater_unzipExecutable(t *testing.T) {
	target, _ := ioutil.TempFile("", "unzip-test")
	defer target.Close()

	source, _ := os.Open(filepath.Join("..", "fixtures", "gh.zip"))
	defer source.Close()

	_, err := io.Copy(target, source)
	assert.Equal(t, nil, err)

	exec, err := unzipExecutable(target.Name())
	assert.Equal(t, nil, err)
	assert.Equal(t, "gh", filepath.Base(exec))
}

func TestUpdater_timeToUpdate(t *testing.T) {
	// file doesn't exist
	timestampDir, _ := ioutil.TempDir("", "timestampDir-test")
	timestampPath := filepath.Join(timestampDir, "gh-update")
	updater := Updater{timestampPath: timestampPath}

	assert.T(t, updater.timeToUpdate())
	timestamp, err := ioutil.ReadFile(timestampPath)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", string(timestamp))

	// invalid timestamp format
	timestampFile, _ := ioutil.TempFile("", "timestampFile-test")
	updater = Updater{timestampPath: timestampFile.Name()}
	assert.T(t, updater.timeToUpdate())
	timestamp, err = ioutil.ReadFile(timestampFile.Name())
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", string(timestamp))

	// dev version
	updater = Updater{CurrentVersion: "dev"}
	assert.T(t, !updater.timeToUpdate())
}

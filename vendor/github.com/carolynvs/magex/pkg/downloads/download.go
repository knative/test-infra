package downloads

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/carolynvs/magex/pkg/gopath"
	"github.com/carolynvs/magex/shx"
	"github.com/carolynvs/magex/xplat"
)

// PostDownloadHook is the handler called after downloading a file, which returns the absolute path to the binary.
type PostDownloadHook func(archivePath string) (string, error)

// DownloadOptions is the configuration settings used to download a file.
type DownloadOptions struct {
	// UrlTemplate is the Go template for the URL to download. Required.
	// Available Template Variables:
	//   - {{.GOOS}}
	//   - {{.GOARCH}}
	//   - {{.EXT}}
	//   - {{.VERSION}}
	UrlTemplate string

	// Name of the binary, excluding OS specific file extension. Required.
	Name string

	// Version to replace {{.VERSION}} in the URL template. Optional depending on whether or not the version is in the UrlTemplate.
	Version string

	// Ext to replace {{.EXT}} in the URL template. Optional, defaults to xplat.FileExt().
	Ext string

	// OsReplacement maps from a GOOS to the os keyword used for the download. Optional, defaults to empty.
	OsReplacement map[string]string

	// ArchReplacement maps from a GOARCH to the arch keyword used for the download. Optional, defaults to empty.
	ArchReplacement map[string]string

	// Hook to call after downloading the file.
	Hook PostDownloadHook
}

// DownloadToGopathBin takes a Go templated URL and expands template variables
// - srcTemplate is the URL
// - version is the version to substitute into the template
// - ext is the file extension to substitute into the template
//
// Template Variables:
// - {{.GOOS}}
// - {{.GOARCH}}
// - {{.EXT}}
// - {{.VERSION}}
func DownloadToGopathBin(opts DownloadOptions) error {

	if err := gopath.EnsureGopathBin(); err != nil {
		return err
	}
	bin := gopath.GetGopathBin()
	return Download(bin, opts)
}

func Download(destDir string, opts DownloadOptions) error {
	src, err := RenderTemplate(opts.UrlTemplate, opts)
	if err != nil {
		return err
	}
	log.Printf("Downloading %s...", src)

	// Download to a temp file
	tmpDir, err := ioutil.TempDir("", "magex")
	if err != nil {
		return fmt.Errorf("could not create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	tmpFile := filepath.Join(tmpDir, filepath.Base(src))

	r, err := http.Get(src)
	if err != nil {
		return fmt.Errorf("could not resolve %s: %w", src, err)
	}
	defer r.Body.Close()
	if r.StatusCode >= 400 {
		return fmt.Errorf("error downloading %s (%d): %s", src, r.StatusCode, r.Status)
	}

	f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", tmpFile, err)
	}
	defer f.Close()

	// Download to the temp file
	_, err = io.Copy(f, r.Body)
	if err != nil {
		return fmt.Errorf("error downloading %s: %w", src, err)
	}
	f.Close()

	// Call a hook to allow for extracting or modifying the downloaded file
	var tmpBin = tmpFile
	if opts.Hook != nil {
		tmpBin, err = opts.Hook(tmpFile)
		if err != nil {
			return err
		}
	}

	// Make the binary executable
	err = os.Chmod(tmpBin, 0755)
	if err != nil {
		return fmt.Errorf("could not make %s executable: %w", tmpBin, err)
	}

	// Move it to the destination
	destPath := filepath.Join(destDir, opts.Name+xplat.FileExt())
	if err := shx.Copy(tmpBin, destPath); err != nil {
		return fmt.Errorf("error copying %s to %s: %w", tmpBin, destPath, err)
	}
	return nil
}

// RenderTemplate takes a Go templated string and expands template variables
// Available Template Variables:
// - {{.GOOS}}
// - {{.GOARCH}}
// - {{.EXT}}
// - {{.VERSION}}
func RenderTemplate(tmplContents string, opts DownloadOptions) (string, error) {
	tmpl, err := template.New("url").Parse(tmplContents)
	if err != nil {
		return "", fmt.Errorf("error parsing %s as a Go template: %w", opts.UrlTemplate, err)
	}

	srcData := struct {
		GOOS    string
		GOARCH  string
		EXT     string
		VERSION string
	}{
		GOOS:    runtime.GOOS,
		GOARCH:  runtime.GOARCH,
		EXT:     opts.Ext,
		VERSION: opts.Version,
	}

	if overrideGoos, ok := opts.OsReplacement[runtime.GOOS]; ok {
		srcData.GOOS = overrideGoos
	}

	if overrideGoarch, ok := opts.ArchReplacement[runtime.GOARCH]; ok {
		srcData.GOARCH = overrideGoarch
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, srcData)
	if err != nil {
		return "", fmt.Errorf("error rendering %s as a Go template with data: %#v: %w", opts.UrlTemplate, srcData, err)
	}

	return buf.String(), nil
}

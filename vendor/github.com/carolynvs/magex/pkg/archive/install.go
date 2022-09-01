package archive

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/carolynvs/magex/pkg/downloads"
	"github.com/carolynvs/magex/xplat"
	"github.com/mholt/archiver/v3"
	_ "github.com/mholt/archiver/v3"
)

// DownloadArchiveOptions are the set of options available for DownloadToGopathBin.
type DownloadArchiveOptions struct {
	downloads.DownloadOptions

	// ArchiveExtensions maps from the GOOS to the expected extension. Required.
	// For example, windows may use .zip while darwin/linux uses .tgz.
	ArchiveExtensions map[string]string

	// TargetFileTemplate specifies the path to the target binary in the archive. Required.
	// Supports the same templating as downloads.DownloadOptions.UrlTemplate.
	TargetFileTemplate string
}

// DownloadToGopathBin downloads an archived file to GOPATH/bin.
func DownloadToGopathBin(opts DownloadArchiveOptions) error {
	// determine the appropriate file extension based on the OS, e.g. windows gets .zip, otherwise .tgz
	opts.Ext = opts.ArchiveExtensions[runtime.GOOS]
	if opts.Ext == "" {
		return fmt.Errorf("no archive file extension was specified for the current GOOS (%s)", runtime.GOOS)
	}

	if opts.Hook == nil {
		opts.Hook = ExtractBinaryFromArchiveHook(opts)
	}

	return downloads.DownloadToGopathBin(opts.DownloadOptions)
}

// ExtractBinaryFromArchiveHook is the default hook for DownloadToGopathBin.
func ExtractBinaryFromArchiveHook(opts DownloadArchiveOptions) downloads.PostDownloadHook {
	return func(archiveFile string) (binPath string, err error) {
		// Save the binary next to the archive file in the temp directory
		outDir := filepath.Dir(archiveFile)

		// Render the name of the file in the archive
		opts.Ext = xplat.FileExt()
		targetFile, err := downloads.RenderTemplate(opts.TargetFileTemplate, opts.DownloadOptions)
		if err != nil {
			return "", fmt.Errorf("error rendering TargetFileTemplate %q with data %#v: %w", opts.TargetFileTemplate, opts.DownloadOptions, err)
		}

		log.Printf("extracting %s from %s...\n", targetFile, archiveFile)

		// Extract the binary
		err = archiver.Extract(archiveFile, targetFile, outDir)
		if err != nil {
			return "", fmt.Errorf("unable to unpack %s: %w", archiveFile, err)
		}

		// The extracted file may be nested depending on its position in the archive
		binFile := filepath.Join(outDir, targetFile)

		// Check that file was extracted, Extract doesn't error out if you give it a missing targetFile
		if _, err := os.Stat(binFile); os.IsNotExist(err) {
			return "", fmt.Errorf("could not find %s in the archive", targetFile)
		}

		return binFile, nil
	}
}

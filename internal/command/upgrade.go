package command

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strings"

	"github.com/axetroy/dvs/internal/dir"
	"github.com/axetroy/dvs/internal/util"
	Version "github.com/axetroy/dvs/internal/version"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// upgrade dvs
func Upgrade(version string, force bool) error {
	var cacheDir = dir.CacheDir

	var (
		err         error
		tarFilename = fmt.Sprintf("dvs_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
		tarFilepath = path.Join(cacheDir, tarFilename)
	)

	downloadURL := fmt.Sprintf("https://github.com/axetroy/dvs/releases/download/%s/%s", version, tarFilename)

	defer func() {
		if err != nil {
			fmt.Printf("If the upgrade fails, download from the `%s` and upgrade manually.\n", downloadURL)
		}
	}()

	// get current dvs version
	dvsExecutablePath, err := os.Executable()

	if err != nil {
		return err
	}

	currentVersion := Version.GetCurrentUsingVersion()

	if !force && version == Version.GetCurrentUsingVersion() {
		fmt.Printf("You are using the latest version `%s`\n", color.GreenString(version))
		return nil
	}

	fmt.Printf("Upgrade dvs `%s` to `%s`\n", currentVersion, version)

	defer os.RemoveAll(cacheDir)

	quit := make(chan os.Signal)
	signal.Notify(quit, util.GetAbortSignals()...)

	go func() {
		<-quit
		fmt.Printf("What made you cancel the download? you can download the file via `%s` and update manually.\n", downloadURL)
		fmt.Println("Good Luck :)")
		_ = os.RemoveAll(cacheDir)
		os.Exit(255)
	}()

	if err = util.DownloadFile(tarFilepath, downloadURL); err != nil {
		return errors.Wrap(err, "download fail")
	}

	// decompress the tag
	if err := decompress(tarFilepath, cacheDir); err != nil {
		return errors.Wrap(err, "unzip fail")
	}

	downloadeddvsFilepath := path.Join(cacheDir, "dvs")

	if runtime.GOOS == "windows" && !strings.HasSuffix(downloadeddvsFilepath, ".exe") {
		// Ensure to add '.exe' to given path on Windows
		downloadeddvsFilepath += ".exe"
	}

	if err := util.UpgradeCommand(downloadeddvsFilepath, dvsExecutablePath); err != nil {
		return errors.Wrap(err, "upgrade fail")
	}

	ps := exec.Command(dvsExecutablePath, "--help")

	ps.Stderr = os.Stderr
	ps.Stdout = os.Stdout

	if err := ps.Run(); err != nil {
		return errors.Wrap(err, "upgrade fail")
	}

	fmt.Printf("dvs upgrade success at `%s`\n", dvsExecutablePath)

	return nil
}

// decompress tar.gz
func decompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)

	if err != nil {
		return errors.Wrapf(err, "open file `%s` fail", tarFile)
	}

	defer srcFile.Close()

	gr, err := gzip.NewReader(srcFile)

	if err != nil {
		return errors.Wrapf(err, "read zip file fail")
	}

	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Wrap(err, "read from zip file fail")
		}

		filename := path.Join(dest, hdr.Name)

		file, err := os.Create(filename)

		if err != nil {
			return errors.Wrapf(err, "unzip and create file `%s` fail\n", filename)
		}

		if runtime.GOOS != "windows" {
			if err := file.Chmod(os.FileMode(hdr.Mode)); err != nil {
				_ = file.Close()
				return errors.Wrap(err, "change file mode fail")
			}
		}

		if _, err := io.Copy(file, tr); err != nil {
			_ = file.Close()
			return errors.Wrap(err, "copy file from zip fail")
		}

		_ = file.Close()
	}

	return nil
}

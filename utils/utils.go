package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// the root asset directory where all the external files used are stored
const ASSET_ROOT_DIR = "assets"

// AssetPath will return the relative path to the file represented by
// assetName otherwise it will return an error if the file doesn't exist.
func AssetPath(assetName string) (string, error) {

	relativePath := path.Join("..", ASSET_ROOT_DIR, assetName)

	if _, err := os.Stat(relativePath); os.IsNotExist(err) {
		return "", err
	}

	return relativePath, nil
}

// SysAssetPath will return the relative path to the file represented by
// assetName but also add in the GOOS after the filename and before the
// extension. This allows loading system-specific files with one command instead
// of a complicated switch statement every time. For instance if your GOOS ==
// "linux" then you can open "config_linux.json" by just asking for
// "config.json". This is extremely powerful and a core component of loading
// system-specific shell commands based on the current GOOS.
func SysAssetPath(assetName string) (string, error) {

	var relativeName bytes.Buffer
	var extIndex int

	fileExt := filepath.Ext(assetName)

	if fileExt != "" {
		// if there is an extension, insert right before it
		extIndex = strings.Index(assetName, fileExt)
	} else {
		// if there is no extension, insert at the end of the name
		extIndex = len(assetName)
	}

	relativeName.WriteString(path.Join("..", ASSET_ROOT_DIR, assetName[0:extIndex]))

	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		relativeName.WriteString("_")
		relativeName.WriteString(runtime.GOOS)
	default:
		return "", fmt.Errorf("Invalid GOOS for asset string: %v", runtime.GOOS)
	}

	relativeName.WriteString(assetName[extIndex:])
	path := relativeName.String()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("Relative file does not exist: %v", path)
	}

	return relativeName.String(), nil
}

// FullDateString will return the current time formatted as a string.
func FullDateString() string {
	return time.Now().String()
}

// FullDateStringSafe returns the current time as a string with only file-name
// safe characters. Used to quickly and easily generate unique file names based
// off of the current system time.
func FullDateStringSafe() string {
	t := time.Now()
	return fmt.Sprintf("[%v-%02d-%02d][%02d_%02d_%02d.%02d]",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

// TimeStampFileName will generate a string to be used to uniquely name and
// identify a new file. The characters used in the file name are guaranteed to
// be safe for a file. The given fileBaseName will be used to append to the
// beginning of the file name to give some control over the file name.
func TimeStampFileName(fileBaseName string, fileExtension string) string {
	dts := FullDateStringSafe()

	var nameBuffer bytes.Buffer
	nameBuffer.WriteString(fileBaseName)
	nameBuffer.WriteString("_")
	nameBuffer.WriteString(dts)
	nameBuffer.WriteString(fileExtension)

	return nameBuffer.String()
}

// ReadLines reads in a file by path and returns a slice of strings
// credit to: https://stackoverflow.com/a/18479916/584947
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// ExternalIPAddress will get this current computer's external IP address.
// credit to: https://gist.github.com/jniltinho/9788121
func ExternalIPAddress() (string, error) {

	var ipBuffer bytes.Buffer

	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	_, copyErr := io.Copy(&ipBuffer, resp.Body)
	if copyErr != nil {
		return "", copyErr
	}

	return strings.Trim(ipBuffer.String(), "\n"), nil
}

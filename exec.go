package you

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var cmd *exec.Cmd
var cmdPort = ""

func Exec(port string) {
	cmdPort = port
	app := appPath()

	if !fileExists(app) {
		logrus.Fatalf("executable file not exists: %s", app)
		return
	}

	cmd = exec.Command(app, "--port", port)
	go func() {
		if err := cmd.Run(); err != nil {
			logrus.Fatalf("executable file error: %v", err)
			return
		}
	}()

	time.Sleep(5 * time.Second)
	logrus.Info("helper exec running ...")
}

func appPath() string {
	app := "bin/"

	if runtime.GOARCH == "arm" || runtime.GOARCH == "arm64" {
		app += "-arm64"
	}

	if runtime.GOOS == "windows" {
		app += ".exe"
	}

	switch runtime.GOOS {
	case "linux":
		if runtime.GOARCH == "arm" || runtime.GOARCH == "arm64" {
			app += "linux/helper-arm64"
		} else {
			app += "linux/helper"
		}
	case "darwin":
		app += "osx/helper"
	case "windows":
		app += "windows/helper.exe"
	default:
		logrus.Fatalf("Unsupported platform: %s", runtime.GOOS)
	}
	return app
}

func Exit() {
	if cmd == nil {
		return
	}
	_ = cmd.Process.Kill()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

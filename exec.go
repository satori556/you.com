package you

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var cmd *exec.Cmd
var cmdPort = ""

// tips: arm64 平台 特征过于明显，无法过盾
func Exec(port string, stdout io.Writer, stderr io.Writer) {
	cmdPort = port
	app := appPath()

	if !fileExists(app) {
		logrus.Fatalf("executable file not exists: %s", app)
		return
	}

	cmd = exec.Command(app, "--port", port)
	if stdout == nil {
		cmd.Stdout = os.Stdout
	}
	if stderr == nil {
		cmd.Stderr = os.Stderr
	}

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
		// 可惜了，arm过不了验证
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

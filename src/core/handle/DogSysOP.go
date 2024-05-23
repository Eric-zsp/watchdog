package handle

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/Eric-zsp/watchdog/src/core/global"
	"github.com/Eric-zsp/watchdog/src/core/utils/files"

	gologs "github.com/cn-joyconn/gologs"
	"github.com/cn-joyconn/goutils/filetool"
	"github.com/cn-joyconn/goutils/strtool"
)

type DogUpgradeOP struct{}

func (op *DogUpgradeOP) getSaveUpgradeDir(appName string) string {
	return path.Join(filetool.SelfDir(), global.UpgradeFileDir, appName)
}
func (op *DogUpgradeOP) getSaveUpgradeFile(appName string) (string, error) {
	dir := op.getSaveUpgradeDir(appName)
	result := path.Join(dir, "upgrade.zip")
	if !filetool.IsExist(dir) {
		err := os.MkdirAll(dir, 0755)
		return result, err
	}
	return result, nil
}
func (op *DogUpgradeOP) getUpgradeUnzipDir(appName string) string {
	return path.Join(op.getSaveUpgradeDir(appName), "update")
}
func (op *DogUpgradeOP) getUpgradeBackupDir(appName string) string {
	return path.Join(op.getSaveUpgradeDir(appName), "backup")
}

func (op *DogUpgradeOP) unzipDir(appName string) error {
	unzipDir := op.getUpgradeUnzipDir(appName)
	if filetool.IsExist(unzipDir) {
		gologs.GetLogger("default").Sugar().Info("unzipDir delete temp dir ", appName)
		os.RemoveAll(unzipDir)
	}
	err := os.MkdirAll(unzipDir, 0755)
	if err != nil {
		return err
	}
	zipFile, err1 := op.getSaveUpgradeFile(appName)
	if err1 != nil {
		return err1
	}
	gologs.GetLogger("default").Sugar().Info("unzipDir begin DeCompress ", appName)
	err3 := files.DeCompress(zipFile, unzipDir)
	if err3 != nil {
		gologs.GetLogger("default").Sugar().Info("unzipDir DeCompress error ", err3.Error(), appName)
	}
	return err3
}
func (op *DogUpgradeOP) moveFiles(appName string, baseDir string, listPaths string) error {
	if strtool.IsBlank(listPaths) {
		return nil
	}
	fpaths := strings.Split(listPaths, ",")
	if fpaths != nil {
		var oldPath string
		var backupPath string
		var upgradePath string
		backupRoot := op.getUpgradeBackupDir(appName)
		if filetool.IsExist(backupRoot) {
			gologs.GetLogger("default").Sugar().Info("backup delete temp dir ", backupRoot)
			os.RemoveAll(backupRoot)
		}
		upgradeRoot := op.getUpgradeUnzipDir(appName)
		for _, fp := range fpaths {
			if !strtool.IsBlank(fp) {

				gologs.GetLogger("default").Sugar().Info("replace file ", fp, " ", appName)
				oldPath = path.Join(baseDir, fp)
				backupPath = path.Join(backupRoot, fp)
				upgradePath = path.Join(upgradeRoot, fp)
				if files.IsDirectory(oldPath) {
					gologs.GetLogger("default").Sugar().Info("CopyDir ", oldPath, " to ", backupPath)
					files.CopyDir(oldPath, backupPath)
					gologs.GetLogger("default").Sugar().Info("Remove dir ", oldPath)
					os.RemoveAll(oldPath)
					gologs.GetLogger("default").Sugar().Info("CopyDir ", upgradePath, " to ", oldPath)
					files.CopyDir(upgradePath, oldPath)
				} else {
					gologs.GetLogger("default").Sugar().Info("CopyFile ", oldPath, " to ", backupPath)
					files.CopyFile(oldPath, backupPath)
					gologs.GetLogger("default").Sugar().Info("Remove file", oldPath)
					os.Remove(oldPath)
					gologs.GetLogger("default").Sugar().Info("CopyFile ", upgradePath, " to ", oldPath)
					files.CopyFile(upgradePath, oldPath)
				}
			}
		}
	}
	return nil
}

func (op *DogUpgradeOP) DoUpgrade(svrID string, baseDir string, durl string, listPath string) error {

	gologs.GetLogger("default").Sugar().Info("DoUpgrade ", svrID)
	savePath, _ := op.getSaveUpgradeFile(svrID)
	gologs.GetLogger("default").Sugar().Info("begin download ", durl)
	downloadFile(durl, savePath)

	e1 := op.unzipDir(svrID)
	if e1 != nil {
		return e1
	}
	gologs.GetLogger("default").Sugar().Info("stop service  ", svrID)
	var cmdStr string
	if runtime.GOOS == "windows" {
		cmdStr, e1 := files.SaveFileByes("stopApp.bat", []byte(global.WindowsCMDAdminAuth+"net stop "+svrID))
		if e1 != nil {
			gologs.GetLogger("default").Sugar().Error("the agent  restarting app happen error,the error msg:" + e1.Error())
		}
		gologs.GetLogger("default").Sugar().Info(cmdStr)
		cmd := exec.Command("cmd.exe", "/C", cmdStr)
		e1 = cmd.Run()
		if e1 != nil {
			gologs.GetLogger("default").Sugar().Info("the agent start app happen error,the error msg:" + e1.Error())
		}
	} else {
		cmdStr = "systemctl stop " + svrID
		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		e1 = cmd.Run()
		if e1 != nil {
			gologs.GetLogger("default").Sugar().Info("the agent stop app happen error,the error msg:" + e1.Error())
		}
	}

	gologs.GetLogger("default").Sugar().Info("begin move files   ", svrID)
	op.moveFiles(svrID, baseDir, listPath)
	gologs.GetLogger("default").Sugar().Info("begin start service   ", svrID)
	if runtime.GOOS == "windows" {
		cmdStr, e1 := files.SaveFileByes("startApp.bat", []byte(global.WindowsCMDAdminAuth+"net start "+svrID))
		if e1 != nil {
			gologs.GetLogger("default").Sugar().Error("the agent  start app happen error,the error msg:" + e1.Error())
		}
		gologs.GetLogger("default").Sugar().Info(cmdStr)
		cmd := exec.Command("cmd.exe", "/C", cmdStr)
		e1 = cmd.Run()
		if e1 != nil {
			gologs.GetLogger("default").Sugar().Info("the agent start app happen error,the error msg:" + e1.Error())
		}
	} else {
		cmdStr = "systemctl start " + svrID
		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		e1 = cmd.Run()
		if e1 != nil {
			gologs.GetLogger("default").Sugar().Info("the agent start app happen error,the error msg:" + e1.Error())
		}
	}

	return nil
}

func downloadFile(url, filename string) {
	// 创建一个新的HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		gologs.GetLogger("default").Sugar().Errorf("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码是否为200 OK
	if resp.StatusCode != http.StatusOK {
		gologs.GetLogger("default").Sugar().Error("Error: Non-200 status code returned")
		return
	}

	if filetool.IsExist(filename) {
		os.Remove(filename)
	}

	// 创建一个文件用于写入
	out, err := os.Create(filename)
	if err != nil {
		gologs.GetLogger("default").Sugar().Errorf("Error creating file:", err)
		return
	}
	defer out.Close()

	// 将响应体复制到文件中
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		gologs.GetLogger("default").Sugar().Errorf("Error writing to file:", err)
		return
	}

	gologs.GetLogger("default").Sugar().Errorf("Successfully downloaded and saved the file:", filename)
}

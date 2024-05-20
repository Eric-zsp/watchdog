package handle

import (
	"os"
	"os/exec"
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
	return filetool.SelfDir() + "/" + global.UpgradeFileDir + "/" + appName
}
func (op *DogUpgradeOP) getSaveUpgradeFile(appName string) (string, error) {
	dir := op.getSaveUpgradeDir(appName)
	result := dir + "/upgrade.zip"
	if !filetool.IsExist(dir) {
		err := os.MkdirAll(dir, 0755)
		return result, err
	}
	return result, nil
}
func (op *DogUpgradeOP) getUpgradeUnzipDir(appName string) string {
	return op.getSaveUpgradeDir(appName) + "/update/"
}
func (op *DogUpgradeOP) getUpgradeBackupDir(appName string) string {
	return op.getSaveUpgradeDir(appName) + "/backup/"
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
		upgradeRoot := op.getUpgradeUnzipDir(appName)
		for _, fp := range fpaths {
			if !strtool.IsBlank(fp) {

				gologs.GetLogger("default").Sugar().Info("replace file ", fp, " ", appName)
				oldPath = baseDir + fp
				backupPath = backupRoot + fp
				upgradePath = upgradeRoot + fp
				if files.IsDirectory(oldPath) {
					files.CopyDir(oldPath, backupPath)
					os.RemoveAll(oldPath)
					files.CopyFile(upgradePath, oldPath)
				} else {
					files.CopyFile(oldPath, backupPath)
					os.Remove(oldPath)
					files.CopyFile(upgradePath, oldPath)
				}
			}
		}
	}
	return nil
}

func (op *DogUpgradeOP) DoUpgrade(svrID string, baseDir string, durl string, listPath string) error {

	gologs.GetLogger("default").Sugar().Info("DoUpgrade ", svrID)
	savePath := op.getUpgradeUnzipDir(svrID)
	gologs.GetLogger("default").Sugar().Info("begin download ", durl)
	files.DownLoadFile(durl, savePath)

	e1 := op.unzipDir(svrID)
	if e1 != nil {
		return e1
	}
	gologs.GetLogger("default").Sugar().Info("stop service  ", svrID)
	if runtime.GOOS == "windows" {
		cmd := exec.Command(global.UpgradeFileDir + "net stop " + svrID)
		cmd.Run()
		// cmd = exec.Command("copy", op.getCurrentCfgDir()+"*", dir)
	} else {
		cmd := exec.Command("systemctl stop " + svrID)
		cmd.Run()
	}
	gologs.GetLogger("default").Sugar().Info("begin move files   ", svrID)
	op.moveFiles(svrID, baseDir, listPath)
	gologs.GetLogger("default").Sugar().Info("begin start service   ", svrID)
	if runtime.GOOS == "windows" {
		cmd := exec.Command(global.UpgradeFileDir + "net start " + svrID)
		cmd.Run()
		// cmd = exec.Command("copy", op.getCurrentCfgDir()+"*", dir)
	} else {
		cmd := exec.Command("systemctl start " + svrID)
		cmd.Run()
	}

	return nil
}
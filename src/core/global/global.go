package global

//升级文件存储位置
var UpgradeFileDir = "upgrade_work"

//程序版本号
var Version = "1.0.0"

// var WindowsCMDAdminAuth = ""
var WindowsCMDAdminAuth = `@echo off 
%1 mshta vbscript:CreateObject("Shell.Application").ShellExecute("cmd.exe","/c %~s0 ::","","runas",1)(window.close)&&exit 
cd /d "%~dp0" 
`

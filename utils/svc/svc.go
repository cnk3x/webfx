package svc

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/template"
)

// 将自己安装成服务
func ServiceInstall(name, workingDirectory, execStart string, dryRun bool) (err error) {
	data := map[string]string{
		"Name":             name,
		"WorkingDirectory": workingDirectory,
		"ExecStart":        execStart,
	}

	buf := bytes.NewBuffer(nil)
	if err = template.Must(template.New("unit").Parse(unitTemplate)).Execute(buf, data); err != nil {
		return
	}

	sFile := fmt.Sprintf("/etc/systemd/system/%s.service", name)

	if runtime.GOOS != "linux" {
		if _, e := exec.LookPath("systemctl"); e != nil {
			return fmt.Errorf("systemctl is not available on this system")
		}
	}

	fmt.Println("will install service:")
	fmt.Println(sFile)
	fmt.Println("----------------------------------------")
	fmt.Println(buf.String())

	if dryRun {
		return
	}

	ServiceControl(name, "stop")
	if err = os.WriteFile(sFile, buf.Bytes(), 0o644); err != nil {
		return
	}

	fmt.Printf("systemd unit file written to %s\n", sFile)
	return ServiceControl("daemon-reload")
}

// 控制服务
func ServiceControl(name string, actions ...string) (err error) {
	execRun := func(name string, arg ...string) error {
		cmd := exec.Command(name, arg...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	if name == "daemon-reload" {
		return execRun("systemctl", name)
	}
	for _, action := range actions {
		if err = execRun("systemctl", action, name); err != nil {
			return
		}
	}
	return
}

const unitTemplate = `[Unit]
Description=hxweb for {{.Name}}

[Service]
Type=simple
WorkingDirectory={{.WorkingDirectory}}
ExecStart={{.ExecStart}}
KillSignal=SIGINT
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
`

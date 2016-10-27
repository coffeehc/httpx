/**
 * Created by coffee on 15/11/23.
 */
package utils
import (
	"html/template"
	"os"
	"path"
	"fmt"
	"flag"
)


func InstallService(serviceName string,args string){
	if !flag.Parsed(){
		flag.Parse()
	}
	_args:=flag.Args()
	needInstall := false
	for _,arg:=range _args{
		if arg == "install"{
			needInstall = true
			break
		}
	}
	if needInstall {
		service := struct {
			ServiceName string
			ServicePath string
			ServiceAegs string
			WorkDir string
		}{serviceName,GetAppPath(),args,GetAppDir()}
		t:= template.Must(template.New("").Parse(sysvScript))
		file, err := os.OpenFile(path.Join("/etc/init.d", serviceName), os.O_CREATE | os.O_RDWR, 0766)
		if err != nil {
			fmt.Printf("打开文件[%s]失败:%s", path.Join("/etc/init.d", serviceName), err)
			os.Exit(-1)
		}
		t.Execute(file, service)
		os.Exit(0)
	}
}


var sysvScript = `#!/bin/sh
# For RedHat and cousins:
# chkconfig: - 99 01
# description: {{.ServiceName}}
# processname: {{.ServicePath}}

### BEGIN INIT INFO
# Provides:          {{.ServicePath}}
# Required-Start:
# Required-Stop:
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
### END INIT INFO

cmd="{{.ServicePath}} {{.ServiceAegs}}"

name=$(basename $0)
pid_file="/var/run/$name.pid"
stdout_log="/var/log/$name.log"
stderr_log="/var/log/$name.err"

get_pid() {
    cat "$pid_file"
}

is_running() {
    [ -f "$pid_file" ] && ps $(get_pid) > /dev/null 2>&1
}

case "$1" in
    start)
        if is_running; then
            echo "Already started"
        else
            echo "Starting $name"
            cd '{{.WorkDir}}'
            $cmd >> "$stdout_log" 2>> "$stderr_log" &
            echo $! > "$pid_file"
            if ! is_running; then
                echo "Unable to start, see $stdout_log and $stderr_log"
                exit 1
            fi
        fi
    ;;
    stop)
        if is_running; then
            echo -n "Stopping $name.."
            kill $(get_pid)
            for i in {1..10}
            do
                if ! is_running; then
                    break
                fi
                echo -n "."
                sleep 1
            done
            echo
            if is_running; then
                echo "Not stopped; may still be shutting down or shutdown may have failed"
                exit 1
            else
                echo "Stopped"
                if [ -f "$pid_file" ]; then
                    rm "$pid_file"
                fi
            fi
        else
            echo "Not running"
        fi
    ;;
    restart)
        $0 stop
        if is_running; then
            echo "Unable to stop, will not attempt to start"
            exit 1
        fi
        $0 start
    ;;
    status)
        if is_running; then
            echo "Running"
        else
            echo "Stopped"
            exit 1
        fi
    ;;
    *)
    echo "Usage: $0 {start|stop|restart|status}"
    exit 1
    ;;
esac
exit 0

`
package ps

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"syscall"
)

func RunLinuxShell(command string) (stdoutRsStr, stderrRsStr string, pid, errNum int, err error) {
	stdoutRsStr = ""
	stderrRsStr = ""
	pid = 0
	errNum = 0

	//command := fmt.Sprintf("ls -l /proc/%d/exe", 123)
	cmd := exec.Command("/bin/bash", "-c", command) //不加第一个第二个参数会报错

	//cmd.Stdout = os.Stdout // cmd.Stdout -> stdout  重定向到标准输出，逐行实时打印
	//cmd.Stderr = os.Stderr // cmd.Stderr -> stderr
	//也可以重定向文件 cmd.Stderr= fd (文件打开的描述符即可)

	stdout, err := cmd.StdoutPipe() //创建标准输出管道
	if err != nil {
		//fmt.Println(err)
		return stdoutRsStr, stderrRsStr,pid, errNum, err
	}

	stderr, err := cmd.StderrPipe() //创建错误输出管道
	if err != nil {
		//fmt.Println(err)
		return stdoutRsStr, stderrRsStr,pid, errNum, err
	}

	defer func() {
		stdout.Close()
		stderr.Close()
	}()

	if err := cmd.Start(); err != nil {
		return stdoutRsStr, stderrRsStr,pid, errNum, err
	}

	//fmt.Println(cmd.Args) //查看当前执行命令

	pid = cmd.Process.Pid //查看命令pid
	//fmt.Println(cmdPid)

	stdoutRsByte, err := ioutil.ReadAll(stdout) // 读取输出结果
	if err != nil {
		return stdoutRsStr, stderrRsStr,pid, errNum, err
	}
	stdoutRsStr = string(stdoutRsByte)
	//fmt.Println("stdout", stdoutRsStr)

	stderrRsByte, err := ioutil.ReadAll(stderr) // 读取输出结果
	if err != nil {
		return stdoutRsStr, stderrRsStr,pid, errNum, err
	}

	stderrRsStr = string(stderrRsByte)
	//fmt.Println("stderror", stderrRsStr)

	if err := cmd.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			//fmt.Println("cmd exit status")
			errNum = ex.Sys().(syscall.WaitStatus).ExitStatus() //获取命令执行返回状态，相当于shell: echo $?
		}
	}

	//fmt.Println("errno", errNum)
	return stdoutRsStr, stderrRsStr,pid, errNum, err
}

func main() {

	command := fmt.Sprintf("ls -l /proc/%d/exe", 123)
	stdout, stderr, errNo,pid, err := RunLinuxShell(command)
	fmt.Println("stdout:", stdout, "stderr:", stderr, "errNo:", "pid:",pid,errNo, "err:", err)
}

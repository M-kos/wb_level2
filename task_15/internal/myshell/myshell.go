package myshell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

const (
	LogicOr         = "||"
	LogicAnd        = "&&"
	Pipe            = "|"
	RedirectBack    = "<"
	RedirectForward = ">"
)

var Commands = map[string]struct{}{
	"cd":   {},
	"pwd":  {},
	"echo": {},
	"kill": {},
	"ps":   {},
}

type LogicPart struct {
	cmd      string
	operator string
}

type MyShell struct {
}

func NewMyShell() *MyShell {
	return &MyShell{}
}

func (sh *MyShell) Run() {
	reader := bufio.NewReader(os.Stdin)

	sh.start(reader)
}

func (sh *MyShell) start(reader *bufio.Reader) {
	for {
		fmt.Print(">>> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			fmt.Println("\nexit")
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		err = sh.processLine(line)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (sh *MyShell) processLine(line string) error {
	logicParts := sh.splitByLogicOperators(line)
	var err error

	for _, part := range logicParts {
		if part.operator == LogicAnd && err != nil {
			continue
		}
		if part.operator == LogicOr && err == nil {
			continue
		}

		err = sh.executePipeline(part.cmd)
	}

	return err
}

func (sh *MyShell) splitByLogicOperators(line string) []LogicPart {
	result := make([]LogicPart, 0)
	j := 0

	for i := 0; i < len(line)-1; i++ {
		if line[i:i+2] == LogicOr {
			result = append(result, LogicPart{
				cmd:      strings.TrimSpace(line[j:i]),
				operator: LogicOr,
			})

			j = i + 2
			i = i + 2
			continue
		}
		if line[i:i+2] == LogicAnd {
			result = append(result, LogicPart{
				cmd:      strings.TrimSpace(line[j:i]),
				operator: LogicAnd,
			})

			j = i + 2
			i = i + 2
			continue
		}
	}

	if j < len(line) {
		result = append(result, LogicPart{
			cmd:      strings.TrimSpace(line[j:]),
			operator: "",
		})
	}

	return result
}

func (sh *MyShell) executePipeline(line string) error {
	parts := strings.Split(line, Pipe)

	var prevPipe io.Reader
	var pipes []*os.File // Для отслеживания открытых pipe'ов
	var procs []*exec.Cmd

	defer func() {
		for _, p := range pipes {
			_ = p.Close()
		}
	}()

	for i, p := range parts {
		p = strings.TrimSpace(p)
		args := sh.parseArgs(p)

		if len(args) == 0 {
			continue
		}

		if _, ok := Commands[args[0]]; ok && len(parts) == 1 {
			return sh.runBuiltin(args)
		}

		cmd := exec.Command(args[0])
		args = args[1:]

		cmd.Args = append([]string{cmd.Path}, args...)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		args = sh.processRedirects(args, cmd)

		if prevPipe != nil {
			cmd.Stdin = prevPipe
		}

		if i < len(parts)-1 {
			r, w, err := os.Pipe()
			if err != nil {
				return err
			}
			cmd.Stdout = w
			pipes = append(pipes, w)
			prevPipe = r
		} else {
			cmd.Stdout = os.Stdout
		}

		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return err
		}

		procs = append(procs, cmd)
	}

	var lastErr error
	for _, p := range procs {
		if err := p.Wait(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (sh *MyShell) parseArgs(line string) []string {
	parts := strings.Fields(line)
	return sh.expandEnvVars(parts)
}

func (sh *MyShell) expandEnvVars(args []string) []string {
	for i, a := range args {
		if strings.Contains(a, "$") {
			args[i] = os.ExpandEnv(a)
		}
	}
	return args
}

func (sh *MyShell) runBuiltin(args []string) error {
	if len(args) == 0 {
		return nil
	}

	switch args[0] {
	case "cd":
		target := "."
		if len(args) > 1 {
			target = args[1]
		}

		if err := os.Chdir(target); err != nil {
			return err
		}

		return nil
	case "pwd":
		d, err := os.Getwd()
		if err != nil {
			return err
		}

		fmt.Println(d)

		return nil
	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
		return nil
	case "kill":
		if len(args) < 2 {
			return fmt.Errorf("usage: kill <pid>")
		}

		pid, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("kill: invalid pid")
		}

		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			return err
		}

		return nil
	case "ps":
		cmd := exec.Command("ps")
		out, err := cmd.Output()
		if err != nil {
			return err
		}

		fmt.Println(string(out))

		return nil
	default:
		return fmt.Errorf("unknown builtin command: %s", args[0])
	}
}

func (sh *MyShell) processRedirects(args []string, cmd *exec.Cmd) []string {
	var res []string

	for i := 0; i < len(args); i++ {
		if args[i] == RedirectForward && i < len(args)-1 {
			f, err := os.Create(args[i+1])
			if err != nil {
				fmt.Println("redirect error:", err)
				return res
			}
			cmd.Stdout = f
			i++
			continue
		}

		if args[i] == RedirectBack && i+1 < len(args) {
			f, err := os.Open(args[i+1])
			if err != nil {
				fmt.Println("redirect error:", err)
				return res
			}
			cmd.Stdin = f
			i++
			continue
		}

		res = append(res, args[i])
	}

	return res
}

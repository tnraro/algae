package alga

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sync"
	"tnraro/algae/internal/util"
)

func CreateAlga(name string, compose string, env string) (string, *AlgaError) {
	if !CheckName(name) {
		return "", Errorf(400, `The alga name "%s" does not match the pattern`, name)
	}
	if hasAlga(name) {
		return "", Errorf(409, "The alga \"%s\" already exists", name)
	}
	if err := createDir(name); err != nil {
		clear(name)
		return "", err
	}
	if err := write(name, "compose.yml", compose); err != nil {
		clear(name)
		return "", err
	}
	if err := write(name, ".env", env); err != nil {
		clear(name)
		return "", err
	}
	log0, err0 := config(name)
	if err0 != nil {
		clear(name)
		return "", err0
	}
	log1, err1 := upAlgaWithPull(name)
	if err1 != nil {
		clear(name)
		return "", err1
	}
	return log0 + log1, nil
}

type getAlga struct {
	Name    string
	Compose string
	Env     string
}

func GetAlga(name string) (*getAlga, *AlgaError) {
	if !CheckName(name) {
		return nil, Errorf(400, `The alga name "%s" does not match the pattern`, name)
	}
	if !hasAlga(name) {
		return nil, Errorf(404, "The alga \"%s\" not exists", name)
	}
	compose, err0 := read(name, "compose.yml")
	if err0 != nil {
		return nil, Errorf(500, `compose.yml not exists`)
	}
	env, err1 := read(name, ".env")
	if err1 != nil {
		return nil, Errorf(500, ".env not exists")
	}
	return &getAlga{
		Name:    name,
		Compose: compose,
		Env:     env,
	}, nil
}
func DeleteAlga(name string) (string, *AlgaError) {
	if !CheckName(name) {
		return "", Errorf(400, `The alga name "%s" does not match the pattern`, name)
	}
	if !hasAlga(name) {
		return "", Errorf(404, "The alga \"%s\" not exists", name)
	}
	logs, err := downAlga(name)
	if err != nil {
		return "", err
	}
	if err := clear(name); err != nil {
		return "", err
	}
	return logs, nil
}
func UpdateAlga(name string, compose string, env string) (string, *AlgaError) {
	if !CheckName(name) {
		return "", Errorf(400, `The alga name "%s" does not match the pattern`, name)
	}
	if !hasAlga(name) {
		return "", Errorf(404, "The alga \"%s\" not exists", name)
	}
	log0, err0 := downAlga(name)
	if err0 != nil {
		return "", err0
	}
	composeBackup := ""
	if compose != "" {
		before, err := read(name, "compose.yml")
		if err != nil {
			return "", err
		}
		composeBackup = before
		if err := write(name, "compose.yml", compose); err != nil {
			composeBackup = ""
			return "", err
		}
	}
	envBackup := ""
	if env != "" {
		before, err := read(name, ".env")
		if err != nil {
			return "", err
		}
		envBackup = before
		if err := write(name, ".env", env); err != nil {
			envBackup = ""
			if composeBackup != "" {
				if err := write(name, "compose.yml", composeBackup); err != nil {
					fmt.Println("failed to restore compose.yml", err)
				}
			}
			return "", err
		}
	}
	log1, err1 := config(name)
	if err1 != nil {
		if composeBackup != "" {
			if err := write(name, "compose.yml", composeBackup); err != nil {
				fmt.Println("failed to restore compose.yml", err)
			}
		}
		if envBackup != "" {
			if err := write(name, ".env", envBackup); err != nil {
				fmt.Println("failed to restore .env", err)
			}
		}
		return "", err1
	}
	log2, err2 := upAlgaWithPull(name)
	if err2 != nil {
		if composeBackup != "" {
			if err := write(name, "compose.yml", composeBackup); err != nil {
				fmt.Println("failed to restore compose.yml", err)
			}
		}
		if envBackup != "" {
			if err := write(name, ".env", envBackup); err != nil {
				fmt.Println("failed to restore .env", err)
			}
		}
		return "", err1
	}
	return log0 + log1 + log2, nil
}
func GetAlgae() ([]string, *AlgaError) {
	files, err := os.ReadDir(util.DataDir("algae"))
	if err != nil {
		return nil, Error(500, err.Error())
	}

	result := make([]string, len(files))
	for i, f := range files {
		result[i] = f.Name()
	}
	return result, nil
}

func run(name string, command string, args ...string) (string, *AlgaError) {
	cmd := exec.Command(command, args...)
	cmd.Dir = AlgaDir(name)
	result, err := cmd.Output()

	if err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("failed executing:", args, err)
			return "", Error(500, err.Error())
		case *exec.ExitError:
			fmt.Println("exit:", e.ExitCode(), string(e.Stderr))
			return "", Error(500, string(e.Stderr))
		default:
			fmt.Println("unexpected error:", err)
			return "", Error(500, err.Error())
		}
	}
	return string(result), nil
}

func AlgaDir(name string, v ...string) string {
	return util.DataDir("algae", name, path.Join(v...))
}

func hasAlga(name string) bool {
	if _, err := os.Stat(AlgaDir(name)); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		fmt.Println(err)
		return false
	}
}
func config(name string) (string, *AlgaError) {
	return run(name, "docker", "compose", "config")
}
func upAlgaWithPull(name string) (string, *AlgaError) {
	return run(name, "docker", "compose", "up", "-d", "--pull=always")
}
func downAlga(name string) (string, *AlgaError) {
	return run(name, "docker", "compose", "down", "--remove-orphans")
}

func write(name string, filename string, content string) *AlgaError {
	err := os.WriteFile(AlgaDir(name, filename), []byte(content), 0755)
	if err != nil {
		return Errorf(500, "failed to write %s/%s", name, filename)
	}
	return nil
}

func read(name string, filename string) (string, *AlgaError) {
	result, err := os.ReadFile(AlgaDir(name, filename))
	if err != nil {
		return "", Errorf(500, "failed to read %s/%s", name, filename)
	}
	return string(result), nil
}

func createDir(name string) *AlgaError {
	err := os.MkdirAll(AlgaDir(name), 0755)
	if err != nil {
		return Errorf(500, "failed to make %s", name)
	}
	return nil
}

func clear(name string) *AlgaError {
	err := os.RemoveAll(AlgaDir(name))
	if err != nil {
		return Errorf(500, "failed to clear %s", name)
	}
	return nil
}

func Login(registry string, username string, password string) (string, *AlgaError) {
	cmd := exec.Command("docker", "login", registry, "-u", username, "-p", password)
	result, err := cmd.Output()

	if err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("failed executing:", err)
			return "", Error(500, err.Error())
		case *exec.ExitError:
			fmt.Println("exit:", e.ExitCode(), string(e.Stderr))
			return "", Error(500, string(e.Stderr))
		default:
			fmt.Println("unexpected error:", err)
			return "", Error(500, err.Error())
		}
	}
	return string(result), nil
}

var (
	nameRe      *regexp.Regexp
	nameReError error
	nameReOnce  sync.Once
)

func initNameRe() {
	nameRe, nameReError = regexp.Compile(`^[\w_-]{2,}$`)
}

func CheckName(name string) bool {
	nameReOnce.Do(initNameRe)
	if nameReError != nil {
		return false
	}
	return nameRe.MatchString(name)
}

type AlgaError struct {
	Code int
	msg  string
}

func (e *AlgaError) Error() string { return e.msg }
func Errorf(code int, format string, a ...any) *AlgaError {
	return &AlgaError{
		Code: code,
		msg:  fmt.Sprintf(format, a...),
	}
}
func Error(code int, msg string) *AlgaError {
	return &AlgaError{
		Code: code,
		msg:  msg,
	}
}

package builder

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/k8s-community/cicd"
	ghIntegr "github.com/k8s-community/github-integration/client"
)

// Process do CICD work: go get of repo, git checkout to given commit, make test and make deploy
func Process(log logrus.FieldLogger, task Task) {
	logger := log.WithFields(logrus.Fields{"source": task.prefix, "user": task.user, "repo": task.repo, "commit": task.commit})

	// TODO: it's good to use something like build.Default.GOPATH, but it doesn't work with daemon
	gopath := os.Getenv("GOPATH")

	url := fmt.Sprintf("%s/%s/%s", task.prefix, task.user, task.repo)
	dir := fmt.Sprintf("%s/src/%s", gopath, url)

	logger.Infof("Remove dir %s", dir)
	err := os.RemoveAll(dir)
	processCommandResult(task.callback, "", err)
	if err != nil {
		logger.Errorf("Couldn't remove directory %s: %s", dir, err)
		return
	}

	var output string

	out, err := runCommand(logger, []string{}, gopath, "go", "get", "-u", url)
	output += out
	processCommandResult(task.callback, output, err)
	if err != nil {
		return
	}

	out, err = runCommand(logger, []string{}, dir, "git", "checkout", task.commit)
	output += out
	processCommandResult(task.callback, output, err)
	if err != nil {
		return
	}

	// Prepare typical Makefile by template from k8s-community/myapp
	out, err = runCommand(
		logger, []string{}, dir, "cp",
		os.Getenv("GOPATH")+"/src/github.com/k8s-community/myapp/Makefile", ".",
	)
	output += out
	processCommandResult(task.callback, output, err)
	if err != nil {
		return
	}

	userEnv := []string{
		"USERSPACE=" + task.user,
		"NAMESPACE=" + task.user,
		"APP=" + task.repo,
		"RELEASE=" + task.version,
	}

	out, err = runCommand(logger, userEnv, dir, "make", "test")
	output += out
	processCommandResult(task.callback, output, err)
	if err != nil {
		return
	}

	if task.task == cicd.TaskDeploy {
		out, err = runCommand(logger, userEnv, dir, "make", "deploy")
		output += out
		processCommandResult(task.callback, output, err)
		if err != nil {
			return
		}
	}

	task.callback(ghIntegr.StateSuccess, output)
}

func runCommand(logger logrus.FieldLogger, env []string, dir, name string, arg ...string) (string, error) {
	logger = logger.WithFields(logrus.Fields{
		"command":        name + " " + strings.Join(arg, " "),
		"additional_env": strings.Join(env, " "),
	})

	logger.Infof("Execute command...")
	command := exec.Command(name, arg...)

	osEnv := append(os.Environ(), env...)
	command.Env = osEnv
	command.Dir = dir

	out, err := command.CombinedOutput()
	commandOut := string(out)

	if len(out) > 0 {
		logger.Info(commandOut)
	}

	if err != nil {
		logger.Errorf("Command failed: %s", err)
		return commandOut, err
	}

	logger.Infof("Done")
	return commandOut, nil
}

func processCommandResult(callback Callback, output string, err error) {
	if err != nil {
		callback(ghIntegr.StateError, output+" \n\nError: "+err.Error())
	} else {
		callback(ghIntegr.StatePending, output)
	}
}

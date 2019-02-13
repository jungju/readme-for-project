package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"text/template"
)

const targetComment = "<!-- rfp:projects:simple1 -->"

func main() {
	token := flag.String("token", "", "token")
	user := flag.String("user", "", "user or org")
	repo := flag.String("repo", "", "repo name")
	targetFile := flag.String("file", "", "target file")
	flag.Parse()

	if token == nil || *token == "" {
		log.Fatal("Require token")
	} else if user == nil || *user == "" {
		log.Fatal("Require user")
	} else if repo == nil || *repo == "" {
		log.Fatal("Require repo")
	} else if targetFile == nil || *targetFile == "" {
		log.Fatal("Require file")
	}

	ret, err := Executer("ok_sh", "./ok.sh", []string{"-j", "list_projects", *user, *repo}, true, "./", []string{"OK_SH_ACCEPT=application/vnd.github.inertia-preview+json", "OK_SH_TOKEN=" + *token})
	if err != nil {
		log.Fatal(err.Error())
	}

	projects := []*Project{}
	if err := json.Unmarshal([]byte(ret), &projects); err != nil {
		log.Fatal(err.Error(), "response:", string(ret))
	}

	final, err := ExecuteTemplateSource(templateSample1, projects, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	replaceContent(*targetFile, final)

	fmt.Println("done.")
}

func Executer(cmdType, cmdName string, cmdArgs []string, cmdWait bool, executeDir string, addEnv []string) (string, error) {
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, addEnv...)
	if executeDir != "" {
		cmd.Dir = executeDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func ExecuteTemplateSource(templateString string, data interface{}, funcMaps template.FuncMap, delims ...[]string) (string, error) {
	t := template.New(reflect.TypeOf(data).Name()).Funcs(funcMaps)
	if len(delims) > 0 {
		t.Delims(delims[0][0], delims[0][1]).Parse(templateString)
	}

	t, err := t.Parse(templateString)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func replaceContent(path, newContents string) {
	goFile := mustRead(path)
	updateGoFile := strings.Replace(goFile, targetComment, targetComment+"\n"+newContents, -1)
	mustWrite(path, updateGoFile)
}
func mustRead(path string) string {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	return string(read)
}

func mustWrite(path, out string) {
	if err := ioutil.WriteFile(path, []byte(out), 0644); err != nil {
		panic(err.Error())
	}
	exec.Command("gofmt", "-w", path).Output()
}

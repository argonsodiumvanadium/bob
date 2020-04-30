package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"gopkg.in/yaml.v2"
	"github.com/tatsushid/go-prettytable"
)

type (

	Attributes struct {
		Name string                `yaml:name`
		History map[string]string   `yaml:projects`
	}

	DirAttr struct {
		Authors  []string
		Name string
		Version string
		CreationDate string
		Description string
		ExecPath string

		Dependencies []string
	}

)

const (
	PROJ_YAML_FILE = "attributes.yaml"
	SRC_YAML_FILE = "data.yaml"
)

func currTime() string {
	return time.Now().Format("Mon Jan _2 15:04:05 2006")
}

func main () {
	args := os.Args
	decidePath(args[1:])
}

func handleErr (e error) {
	fmt.Print("\u001B[91m")
	if e != nil {
		defer fmt.Print("\u001B[0m")
		panic(e)
	}
}

func printErr (args ... interface{}) {
	fmt.Print("\u001B[91m")
	for _,val := range args {
		fmt.Print(val)
	}
	fmt.Println("\u001B[0m")
}

func decidePath (args []string) {
	switch args[0] {
	case "build":
		if len(args) > 1 {
			buildProject(args[1])
		} else {
			printErr("You have not given any directory name to bob, Thus it failed to successfully create a project")
		}
	case "list":
		if len(args) > 1 {
			listProject(args[1:])
		} else {
			listProjects(getAttributes())
		}
	}
}

func buildProject (projName string) {
	handleErr(os.Mkdir(projName,0777))
	handleErr(os.Mkdir(projName+"/src",0777))

	attr := getAttributes()
	file,err := os.OpenFile(projName+"/"+PROJ_YAML_FILE,os.O_CREATE,0777)
	handleErr(err)
	file.Close()

	buildAttributes := &DirAttr{}
	buildAttributes.Name = projName
	buildAttributes.Authors = append(buildAttributes.Authors,attr.Name)
	buildAttributes.Version = "v 0.1.0"
	buildAttributes.CreationDate = currTime()
	buildAttributes.Description = "this is a project built by "+arrToString(buildAttributes.Authors)+" and is called "+buildAttributes.Name

	file,err = os.OpenFile(projName+"/"+PROJ_YAML_FILE,os.O_WRONLY,0444)
	defer file.Close()
	handleErr(err)
	data,err := yaml.Marshal(buildAttributes)
	handleErr(err)
	file.Write(data)
	dir, err := os.Getwd()
	handleErr(err)
	attr.History[projName] = dir+"/"+projName
	attr.save()
	fmt.Println("\u001B[96mThe Project was \u001B[92mSUCCESSFULLY\u001B[96m created \u001b[0m")
}

func arrToString (args []string) (ret string) {
	for _,val := range args[:len(args)-1] {
		ret += val+", "
	}
	ret += args[len(args)-1]
	return
}

func getAttributes () (*Attributes) {
	attr := &Attributes{}
	data,err := ioutil.ReadFile(SRC_YAML_FILE)
	handleErr(err)
	handleErr(yaml.Unmarshal(data,attr))
	
	for attr.Name == "" {
		fmt.Print("\u001b[92mNo name was specified in my yaml file would you care to provide your name \u001B[96m: ")
		
		input := ""
		
		fmt.Scanln(&input)
		fmt.Print("\u001B[0m")

		strings.Replace(input,"\r","",-1)
		strings.Replace(input,"\n","",-1)
		
		attr.Name = input
	}

	return attr
}

func (self *Attributes) save () {
	file,err := os.OpenFile(SRC_YAML_FILE,os.O_CREATE,0777)
	handleErr(err)
	file.Close()
	file,err = os.OpenFile(SRC_YAML_FILE,os.O_WRONLY,0777)
	defer file.Close()
	handleErr(err)
	data,err := yaml.Marshal(self)
	handleErr(err)
	file.Write(data)
}

func listProjects (attr *Attributes) {
	table, err := prettytable.NewTable(
		[]prettytable.Column{
			{Header: "\u001B[92m   Project Name   \u001B[0m"},
			{Header: "\u001B[96m   Project Path   \u001B[0m"},
	}...)

	handleErr(err)

	table.Separator = " \u001B[35m|\u001B[0m "

	table.AddRow("\u001B[35m------------------\u001B[0m","\u001B[35m------------------\u001B[0m")    //for the sake of seperation

	for key,val := range attr.History {
		table.AddRow("\u001B[92m"+key+"\u001B[0m","\u001B[96m"+val+"\u001B[0m")
	}

	fmt.Print("\n\n\n")
	table.Print()
	fmt.Print("\n\n\n")

	attr.save()
}

func listProject (args []string) {
	attr := getAttributes()
	fmt.Println("\n\n\n")

	for _,val := range args {
		path := attr.History[val]
		path = path+"/"+PROJ_YAML_FILE

		data,err := ioutil.ReadFile(path)
		handleErr(err)

		fmt.Println("\u001b[92m---- "+val+" ----\u001b[0m")
		fmt.Println(string(data),"\n\n\n")
	}
}
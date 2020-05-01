package main

import (
	"os"
	"os/exec"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"gopkg.in/yaml.v2"
	"github.com/tatsushid/go-prettytable"
)

type (

	Attributes struct {
		Name string                  `yaml:name`
		Projects map[string]string   `yaml:projects`
	}

	DirAttr struct {
		Authors  []string
		Name string
		Version string
		CreationDate string
		Description string
		ExecPath string

		Dependencies map[string][]string
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
	if len(args) > 0 {
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
		case "add":
			if len(args) > 2 {
				addProject(args[1],args[2])
			} else {
				addProjectAfterGettingName(args[1])
			}
		case "rm":
			rmProject(args[1])
		case "del":
			deleteProject(args[1])
		case "init":
			initializeDependencies(args[1])
		}
	}
}

func buildProject (projName string) {
	handleErr(os.Mkdir(projName,0777))
	handleErr(os.Mkdir(projName+"/src",0777))

	attr := getAttributes()

	attr.makeAttributesFile(projName)

	dir, err := os.Getwd()
	handleErr(err)
	attr.Projects[projName] = dir+"/"+projName
	attr.save()
	fmt.Println("\u001B[96mThe Project was \u001B[92mSUCCESSFULLY\u001B[96m created \u001b[0m")
}

func (attr *Attributes) makeAttributesFile (path string) {
	file,err := os.OpenFile(path+"/"+PROJ_YAML_FILE,os.O_CREATE,0777)
	handleErr(err)
	file.Close()

	elems := strings.Split(path,"/")
	name := elems[len(elems)-1]

	BuildAttributes := &DirAttr{}
	BuildAttributes.Name = name
	BuildAttributes.Authors = append(BuildAttributes.Authors,attr.Name)
	BuildAttributes.Version = "v 0.1.0"
	BuildAttributes.CreationDate = currTime()
	BuildAttributes.Description = "this is a project built by "+arrToString(BuildAttributes.Authors)+" and is called "+BuildAttributes.Name
	BuildAttributes.Dependencies = make(map[string][]string)

	file,err = os.OpenFile(path+"/"+PROJ_YAML_FILE,os.O_WRONLY,0444)
	defer file.Close()
	handleErr(err)
	data,err := yaml.Marshal(BuildAttributes)
	handleErr(err)
	file.Write(data)
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

	if attr.Projects == nil {
		attr.Projects = make(map[string]string)
	}

	return attr
}

func (self *Attributes) getBuildAttr (path string) (*DirAttr) {
	buildAttr := &DirAttr{}

	path = path+"/"+PROJ_YAML_FILE

	data,err := ioutil.ReadFile(path)
	handleErr(err)

	handleErr(yaml.Unmarshal(data,buildAttr))

	return buildAttr
}

func (self *Attributes) save () {
	var err = os.Remove(SRC_YAML_FILE)
	handleErr(err)
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

	for key,val := range attr.Projects {
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
		path := attr.Projects[val]
		path = path+"/"+PROJ_YAML_FILE

		data,err := ioutil.ReadFile(path)
		handleErr(err)

		fmt.Println("\u001b[92m---- "+val+" ----\u001b[0m")
		fmt.Println(string(data),"\n\n\n")
	}
}

func addProjectAfterGettingName (path string) {
	elems := strings.Split(path,"/")
	name := elems[len(elems)-1]
	addProject(path,name)
}

func addProject (path,name string) {
	attr := getAttributes()
	attr.makeAttributesFile(path)

	dir, err := os.Getwd()
	handleErr(err)
	attr.Projects[name] = dir+"/"+path
	attr.save()
	fmt.Println("\u001B[96mThe Project was \u001B[92mSUCCESSFULLY\u001B[96m added\u001b[0m")
}

func rmProject (name string) {
	attr := getAttributes()

	if _,ok := attr.Projects[name];ok {
		delete(attr.Projects,name)
		fmt.Println("\u001B[93mThe Project was \u001B[91mSUCCESSFULLY\u001B[93m removed\u001b[0m")
		attr.save()
	} else {
		printErr("The project name provided does not exist (remove trailing '/' if any)")
	}

}

func deleteProject (name string) {
	rmProject(name)
	attr := getAttributes()
	err := exec.Command("rm","-rf",attr.Projects[name]).Run()
	handleErr(err)
}

func initializeDependencies (name string) {
	attr := getAttributes()
	if _,ok := attr.Projects[name];!ok {
		printErr("The Project does not exist")
		return
	}
	buildAttr := attr.getBuildAttr(attr.Projects[name])
	fmt.Print("\u001b[0m")

	for tool,v := range buildAttr.Dependencies {
		for _,substance := range v {
			command := strings.Split(tool," ")
			command = append(command,strings.Split(substance," ") ...)
			fmt.Println("\u001b[96mRunning ... \u001b[92m",tool,substance,"\u001b[0m")
			cmd := exec.Command(command[0],command[1:] ... )
			err := cmd.Run()
			if err != nil {
				printErr("some error occured")
				fmt.Println(err)
			}
			out,_ := cmd.Output()
			fmt.Println(string(out))
		}
	}

	fmt.Println("\u001B[96mdependencies are \u001B[92mINSTALLED\u001B[0m")
}
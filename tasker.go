package tasker

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strings"
)

//Task common task definition
type Task struct {
	name, datetime, status string
}

//TaskCreate used in creating tasks
type TaskCreate struct {
	schedule    string
	modifier    string
	days        string
	months      string
	idletime    string
	taskname    string
	taskrun     string
	starttime   string
	interval    string
	endtime     string
	duration    string
	terminate   string
	startdate   string
	enddate     string
	channelName string
	noPassword  string
	markDelete  string
	force       string
	level       string
	delaytime   string
}

const (
	taskerFile = "SCHTASKS"
)

var (
	//Schedules list of available scheduling scheme
	Schedules = struct {
		MINUTE, HOURLY, DAILY, WEEKLY, MONTHLY  string
		ONCE, ONSTART, ONLOGON, ONIDLE, ONEVENT string
	}{
		MINUTE: "MINUTE", HOURLY: "HOURLY", DAILY: "DAILY", WEEKLY: "WEEKLY",
		MONTHLY: "MONTHLY", ONCE: "ONCE", ONSTART: "ONSTART", ONLOGON: "ONLOGON",
		ONIDLE: "ONIDLE", ONEVENT: "ONEVENT",
	}

	//Days list of days
	Days = struct {
		MON, TUE, WED, THU, FRI, SAT, SUN, ALL string
	}{
		MON: "MON", TUE: "TUE", WED: "WED",
		THU: "THU", FRI: "FRI", SAT: "SAT", SUN: "SUN",
		ALL: "*",
	}

	//Months list of months
	Months = struct {
		JAN, FEB, MAR, APR, MAY, JUN string
		JUL, AUG, SEP, OCT, NOV, DEC string
		ALL                          string
	}{
		JAN: "JAN", FEB: "FEB", MAR: "MAR", APR: "APR", MAY: "MAY", JUN: "JUN",
		JUL: "JUL", AUG: "AUG", SEP: "SEP", OCT: "OCT", NOV: "NOV", DEC: "DEC",
		ALL: "*",
	}
	//Level Run Levels
	Level = struct {
		LIMITED, HIGHEST string
	}{
		LIMITED: "LIMITED", HIGHEST: "HIGHEST",
	}

	//Commands
	_Create = struct {
		command     string
		schedule    string
		modifier    string
		days        string
		months      string
		idletime    string
		taskname    string
		taskrun     string
		starttime   string
		interval    string
		endtime     string
		duration    string
		terminate   string
		startdate   string
		enddate     string
		channelName string
		noPassword  string
		markDelete  string
		force       string
		level       string
		delaytime   string
	}{
		command:     "/CREATE",
		schedule:    "/SC",
		modifier:    "/MO",
		days:        "/D",
		months:      "/M",
		idletime:    "/I",
		taskname:    "/TN",
		taskrun:     "/TR",
		starttime:   "/ST",
		interval:    "/RI",
		endtime:     "/ET",
		duration:    "/DU",
		terminate:   "/K",
		startdate:   "/SD",
		enddate:     "/ED",
		channelName: "/EC",
		noPassword:  "/NP",
		markDelete:  "/Z",
		force:       "/F",
		level:       "/RL",
		delaytime:   "/DELAY",
	}
	_Delete = struct {
		command  string
		taskname string
		force    string
	}{
		command:  "/DELETE",
		taskname: "/TN",
		force:    "/F",
	}
	_Query = struct {
		command     string
		format      string
		formatCSV   string
		formatLIST  string
		formatTABLE string
		noHeader    string
	}{
		command:     "/QUERY",
		format:      "/FO",
		formatCSV:   "CSV",
		formatLIST:  "LIST",
		formatTABLE: "TABLE",
		noHeader:    "/NH",
	}
)

//SchTask definitions
type SchTask struct {
	bin    string
	prefix string
}

//New creates a new tasker object
func New() SchTask {
	return SchTask{
		bin:    taskerFile,
		prefix: "go-wintask-",
	}
}

func catch(out []byte, e error) {
	if e != nil {
		log.Fatal(string(out))
	}
}

//Create  Enables an administrator to create scheduled tasks on a local or
//remote system.
func (task SchTask) Create(taskcreate TaskCreate) {
	cmds := []string{}

	//make commands
	cmds = append(cmds, _Create.taskname)
	cmds = append(cmds, taskcreate.taskname)

	cmd := exec.Command(task.bin, cmds...)

	output, err := cmd.CombinedOutput()
	catch(output, err)
}

//Delete Deletes one or more scheduled tasks.
func (task SchTask) Delete(taskname string, own, force bool) {
	cmd := &exec.Cmd{}
	if !force {
		cmd = exec.Command(task.bin, _Delete.command, _Delete.taskname, taskname)
	} else {
		cmd = exec.Command(task.bin, _Delete.command, _Delete.taskname, taskname, _Delete.force)
	}

	output, err := cmd.CombinedOutput()
	catch(output, err)
}

//Query Enables an administrator to display the scheduled tasks on the
//local or remote system.
func (task SchTask) Query(name string, own bool) []Task {
	taskList := make([]Task, 0)

	cmd := exec.Command(task.bin, _Query.command, _Query.format, _Query.formatCSV, _Query.noHeader)
	output, err := cmd.CombinedOutput()
	catch(output, err)

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		tx := strings.Replace(scanner.Text(), "\"", "", -1)
		ts := strings.Split(tx, ",")

		tname := strings.TrimSpace(ts[0])

		if own {
			name = task.prefix + name
		}

		if name == "*" || name == "" || strings.Contains(strings.ToLower(tname), strings.ToLower(name)) {
			dtime := strings.TrimSpace(ts[1])
			stat := strings.TrimSpace(ts[2])
			taskList = append(taskList, Task{tname, dtime, stat})
		}
	}

	return taskList
}

//Change updates a task
func (task SchTask) Change(name string, own bool) {

}

//Run

//End

//ShowSid

package tasker

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

//Task common task definition
type Task struct {
	name, datetime, status string
}

//TaskCreate used in creating tasks
//Examples
//==> Creates a scheduled task "doc" on the remote machine "ABC"
//	which runs notepad.exe every hour under user "runasuser".
//	SCHTASKS /Create /S ABC /U user /P password /RU runasuser
//		 /RP runaspassword /SC HOURLY /TN doc /TR notepad
//
//==> Creates a scheduled task "accountant" on the remote machine
//	"ABC" to run calc.exe every five minutes from the specified
//	start time to end time between the start date and end date.
//	SCHTASKS /Create /S ABC /U domain\user /P password /SC MINUTE
//		 /MO 5 /TN accountant /TR calc.exe /ST 12:00 /ET 14:00
//		 /SD 06/06/2006 /ED 06/06/2006 /RU runasuser /RP userpassword
//
//==> Creates a scheduled task "gametime" to run freecell on the
//	first Sunday of every month.
//	SCHTASKS /Create /SC MONTHLY /MO first /D SUN /TN gametime
//		 /TR c:\windows\system32\freecell
//
//==> Creates a scheduled task "report" on remote machine "ABC"
//	to run notepad.exe every week.
//	SCHTASKS /Create /S ABC /U user /P password /RU runasuser
//		 /RP runaspassword /SC WEEKLY /TN report /TR notepad.exe
//
//==> Creates a scheduled task "logtracker" on remote machine "ABC"
//	to run notepad.exe every five minutes starting from the
//	specified start time with no end time. The /RP password will be
//	prompted for.
//	SCHTASKS /Create /S ABC /U domain\user /P password /SC MINUTE
//		 /MO 5 /TN logtracker
//		 /TR c:\windows\system32\notepad.exe /ST 18:30
//		 /RU runasuser /RP
//
//==> Creates a scheduled task "gaming" to run freecell.exe starting
//	at 12:00 and automatically terminating at 14:00 hours every day
//	SCHTASKS /Create /SC DAILY /TN gaming /TR c:\freecell /ST 12:00
//		 /ET 14:00 /K
//
//==> Creates a scheduled task "EventLog" to run wevtvwr.msc starting
//	whenever event 101 is published in the System channel
//	SCHTASKS /Create /TN EventLog /TR wevtvwr.msc /SC ONEVENT
//		 /EC System /MO *[System/EventID=101]
//
//==> Spaces in file paths can be used by using two sets of quotes, one
//	set for CMD.EXE and one for SchTasks.exe.  The outer quotes for CMD
//	need to be double quotes; the inner quotes can be single quotes or
//	escaped double quotes:
//	SCHTASKS /Create
//   /tr "'c:\program files\internet explorer\iexplorer.exe'
//   \"c:\log data\today.xml\"" ...
type TaskCreate struct {
	// /RU  username      Specifies the "run as" user account (user context)
	// 					  under which the task runs. For the system account,
	//					  valid values are "", "NT AUTHORITY\SYSTEM"
	//					  or "SYSTEM".
	//					  For v2 tasks, "NT AUTHORITY\LOCALSERVICE" and
	//					  "NT AUTHORITY\NETWORKSERVICE" are also available as well
	//					  as the well known SIDs for all three.
	Username string

	// /RP  [password]    Specifies the password for the "run as" user.
	//					  To prompt for the password, the value must be either
	//					  "*" or none. This password is ignored for the
	//					  system account. Must be combined with either /RU or
	//					  /XML switch.
	Password string

	// /TN   taskname     Specifies the string in the form of path\name
	//                    which uniquely identifies this scheduled task.
	Taskname string

	// /TR   taskrun      Specifies the path and file name of the program to be
	//					  run at the scheduled time.
	//					  Example: C:\windows\system32\calc.exe
	//								/create /tn "Run my script" /tr "C:\test 2\myscript.cmd \"one of the arguments to pass\"
	//								anotherargument \"This is the third argument\" lastargument"
	Taskrun   string
	Arguments []string

	// /SC   schedule     Specifies the schedule frequency.
	//                    Valid schedule types: MINUTE, HOURLY, DAILY, WEEKLY,
	//                    MONTHLY, ONCE, ONSTART, ONLOGON, ONIDLE, ONEVENT.
	Schedule string

	// /MO   modifier     Refines the schedule type to allow finer control over
	//				      schedule recurrence. Valid values are listed in the
	//				      "Modifiers" section below.
	// MINUTE:  1 - 1439 minutes.
	// HOURLY:  1 - 23 hours.
	// DAILY:   1 - 365 days.
	// WEEKLY:  weeks 1 - 52.
	// ONCE:    No modifiers.
	// ONSTART: No modifiers.
	// ONLOGON: No modifiers.
	// ONIDLE:  No modifiers.
	// MONTHLY: 1 - 12, or FIRST, SECOND, THIRD, FOURTH, LAST, LASTDAY.
	Modifier string

	// /D    days         Specifies the day of the week to run the task. Valid
	//                    values: MON, TUE, WED, THU, FRI, SAT, SUN and for
	//                    MONTHLY schedules 1 - 31 (days of the month).
	//                    Wildcard "*" specifies all days.
	Days []string

	// /M    months       Specifies month(s) of the year. Defaults to the first
	//                    day of the month. Valid values: JAN, FEB, MAR, APR,
	//                    MAY, JUN, JUL, AUG, SEP, OCT, NOV, DEC. Wildcard "*"
	//                    specifies all months.
	Months []string

	// /I    idletime     Specifies the amount of idle time to wait before
	//                    running a scheduled ONIDLE task.
	//                    Valid range: 1 - 999 minutes.
	Idletime string

	// /ST   starttime    Specifies the start time to run the task. The time
	//                    format is HH:mm (24 hour time) for example, 14:30 for
	//                    2:30 PM. Defaults to current time if /ST is not
	//                    specified.  This option is required with /SC ONCE.
	Starttime string

	// /RI   interval     Specifies the repetition interval in minutes. This is
	//                    not applicable for schedule types: MINUTE, HOURLY,
	//                    ONSTART, ONLOGON, ONIDLE, ONEVENT.
	//                    Valid range: 1 - 599940 minutes.
	//                    If either /ET or /DU is specified, then it defaults to
	//                    10 minutes.
	Interval string

	// /ET   endtime      Specifies the end time to run the task. The time format
	//                    is HH:mm (24 hour time) for example, 14:50 for 2:50 PM.
	//                    This is not applicable for schedule types: ONSTART,
	//                    ONLOGON, ONIDLE, ONEVENT.
	Endtime string

	// /DU   duration     Specifies the duration to run the task. The time
	//                    format is HH:mm. This is not applicable with /ET and
	//                    for schedule types: ONSTART, ONLOGON, ONIDLE, ONEVENT.
	//                    For /V1 tasks, if /RI is specified, duration defaults
	//                    to 1 hour.
	Duration string

	// /K     terminate   Terminates the task at the endtime or duration time.
	//                    This is not applicable for schedule types: ONSTART,
	//                    ONLOGON, ONIDLE, ONEVENT. Either /ET or /DU must be
	//                    specified.
	Terminate bool

	// /SD   startdate    Specifies the first date on which the task runs. The
	//                    format is mm/dd/yyyy. Defaults to the current
	//                    date. This is not applicable for schedule types: ONCE,
	//                    ONSTART, ONLOGON, ONIDLE, ONEVENT.
	Startdate string

	// /ED   enddate      Specifies the last date when the task should run. The
	//                    format is mm/dd/yyyy. This is not applicable for
	//                    schedule types: ONCE, ONSTART, ONLOGON, ONIDLE, ONEVENT.
	Enddate string

	// /EC   channelName  Specifies the event channel for OnEvent triggers.
	ChannelName string

	// /NP    noPassword  No password is stored.  The task runs non-interactively
	//                    as the given user.  Only local resources are available.
	NoPassword bool

	// /Z     markDelete  Marks the task for deletion after its final run.
	MarkDelete bool

	// /F                 Forcefully creates the task and suppresses warnings if
	//                    the specified task already exists.
	Force bool

	// /RL   level        Sets the Run Level for the job. Valid values are
	//                    LIMITED and HIGHEST. The default is LIMITED.
	Level string

	// /DELAY delaytime   Specifies the wait time to delay the running of the
	//                    task after the trigger is fired.  The time format is
	//                    mmmm:ss.  This option is only valid for schedule types
	//                    ONSTART, ONLOGON, ONEVENT.
	Delaytime string
}

const (
	taskerFile = "SCHTASKS"
)

var (
	//Debug Enables debugging, commands won't be performed just logged.
	Debug      = false
	dbgMessage = "You are currently in debug mode."

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
	/*************Create**************/
	_Create = struct {
		Command     string
		username    string
		password    string
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
		preVista    string
		level       string
		delaytime   string
	}{
		Command:     "/CREATE",
		username:    "/RU",
		password:    "/RP",
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
		preVista:    "/V1",
		force:       "/F",
		level:       "/RL",
		delaytime:   "/DELAY",
	}
	/*************Delete**************/
	_Delete = struct {
		Command  string
		taskname string
		force    string
	}{
		Command:  "/DELETE",
		taskname: "/TN",
		force:    "/F",
	}
	/*************Query**************/
	_Query = struct {
		Command     string
		format      string
		formatCSV   string
		formatLIST  string
		formatTABLE string
		noHeader    string
	}{
		Command:     "/QUERY",
		format:      "/FO",
		formatCSV:   "CSV",
		formatLIST:  "LIST",
		formatTABLE: "TABLE",
		noHeader:    "/NH",
	}
	/*************Change**************/
	_Change = struct {
		Command string
	}{
		Command: "/CHANGE",
	}
	/*************Run**************/
	_Run = struct {
		Command   string
		immediate string
		taskname  string
	}{
		Command:   "/RUN",
		taskname:  "/TN",
		immediate: "/I",
	}
	/*************End**************/
	_End = struct {
		Command  string
		taskname string
	}{
		Command:  "/END",
		taskname: "/TN",
	}
	/*************ShowSid**************/
	_ShowSid = struct {
		Command  string
		taskname string
	}{
		Command:  "/SHOWSID",
		taskname: "/TN",
	}
)

//SchTask definitions
type SchTask struct {
	bin           string
	prefix        string
	compatibility bool
}

//New creates a new tasker object
func New(com bool) SchTask {
	return SchTask{
		bin:           taskerFile,
		prefix:        "go-wintask-",
		compatibility: com,
	}
}

func catch(out []byte, e error) {
	if e != nil {
		log.Fatal(string(out))
	}
}

func getCurrDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func getCurrExe() string {
	_, file := path.Split(os.Args[0])
	return file
}

//TaskMake for generating tasks
func (task SchTask) TaskMake(taskcreate TaskCreate, command string, own bool) []string {
	cmds := []string{}
	/****make commands****/
	//Append the command
	cmds = append(cmds, command)
	if taskcreate.Username != "" {
		cmds = append(cmds, _Create.username)
		cmds = append(cmds, taskcreate.Username)
	}
	//password string
	if taskcreate.Password != "" {
		cmds = append(cmds, _Create.password)
		cmds = append(cmds, taskcreate.Password)
	}
	//Schedule
	if taskcreate.Schedule != "" {
		cmds = append(cmds, _Create.schedule)
		cmds = append(cmds, taskcreate.Schedule)
	}
	//Modifier
	if taskcreate.Modifier != "" {
		cmds = append(cmds, _Create.modifier)
		cmds = append(cmds, taskcreate.Modifier)
	}
	//days []string
	if len(taskcreate.Days) > 0 {
		f := ""
		for _, d := range taskcreate.Days {
			f += d + ","
		}
		f = strings.TrimSuffix(f, ",")
		if f != "" {
			cmds = append(cmds, _Create.days)
			cmds = append(cmds, f)
		}
	}
	//months []string
	if len(taskcreate.Months) > 0 {
		f := ""
		for _, d := range taskcreate.Months {
			f += d + ","
		}
		f = strings.TrimSuffix(f, ",")
		if f != "" {
			cmds = append(cmds, _Create.months)
			cmds = append(cmds, f)
		}
	}
	//idletime string
	if taskcreate.Idletime != "" {
		cmds = append(cmds, _Create.idletime)
		cmds = append(cmds, taskcreate.Idletime)
	}
	//starttime string
	if taskcreate.Starttime != "" {
		cmds = append(cmds, _Create.starttime)
		cmds = append(cmds, taskcreate.Starttime)
	}
	//interval string
	if taskcreate.Interval != "" {
		cmds = append(cmds, _Create.interval)
		cmds = append(cmds, taskcreate.Interval)
	}
	//endtime string
	if taskcreate.Endtime != "" {
		cmds = append(cmds, _Create.endtime)
		cmds = append(cmds, taskcreate.Endtime)
	}
	//duration string
	if taskcreate.Duration != "" {
		cmds = append(cmds, _Create.duration)
		cmds = append(cmds, taskcreate.Duration)
	}
	//terminate string
	if taskcreate.Terminate {
		cmds = append(cmds, _Create.terminate)
	}
	//startdate string
	if taskcreate.Startdate != "" {
		cmds = append(cmds, _Create.startdate)
		cmds = append(cmds, taskcreate.Startdate)
	}
	//enddate string
	if taskcreate.Enddate != "" {
		cmds = append(cmds, _Create.enddate)
		cmds = append(cmds, taskcreate.Enddate)
	}
	//channelName string
	if taskcreate.ChannelName != "" {
		cmds = append(cmds, _Create.channelName)
		cmds = append(cmds, taskcreate.ChannelName)
	}
	//No Password
	if taskcreate.NoPassword {
		cmds = append(cmds, _Create.noPassword)
	}
	//Force
	if taskcreate.Force {
		cmds = append(cmds, _Create.force)
	}
	//level string
	if taskcreate.Level != "" {
		cmds = append(cmds, _Create.level)
		cmds = append(cmds, taskcreate.Level)
	}
	//delaytime string
	if taskcreate.Delaytime != "" {
		cmds = append(cmds, _Create.delaytime)
		cmds = append(cmds, taskcreate.Delaytime)
	}
	//Add taskname
	cmds = append(cmds, _Create.taskname)
	name := taskcreate.Taskname
	if own {
		name = task.prefix + name
	}
	cmds = append(cmds, name)
	//Add taskrun
	cmds = append(cmds, _Create.taskrun)
	run := taskcreate.Taskrun
	if run == "" {
		run = path.Join(getCurrDir(), getCurrExe())
	}
	args := ""
	//append the args
	for _, arg := range taskcreate.Arguments {
		if strings.IndexRune(arg, ' ') >= 0 {
			args += "\"" + arg + "\" "
		} else {
			args += arg + " "
		}
	}
	run = "\"" + run + "\" " + strings.TrimSpace(args)
	run = strings.TrimSpace(run)
	//run = "\"" + run + "\""
	cmds = append(cmds, run)
	//markDelete bool
	if taskcreate.MarkDelete {
		cmds = append(cmds, _Create.preVista)
		cmds = append(cmds, _Create.markDelete)
	}

	if Debug {
		fmt.Println("Commands:", cmds)
	}
	return cmds
}

//Create  Enables an administrator to create scheduled tasks on a local or
//remote system.
func (task SchTask) Create(taskcreate TaskCreate) string {
	cmds := task.TaskMake(taskcreate, _Create.Command, true)

	if Debug {
		return dbgMessage
	}

	cmd := exec.Command(task.bin, cmds...)

	output, err := cmd.CombinedOutput()
	catch(output, err)

	return string(output)
}

//Delete Deletes one or more scheduled tasks.
func (task SchTask) Delete(taskname string, own, force bool) string {
	cmd := &exec.Cmd{}

	if Debug {
		return dbgMessage
	}

	if own {
		taskname = task.prefix + taskname
	}

	if !force {
		cmd = exec.Command(task.bin, _Delete.Command, _Delete.taskname, taskname)
	} else {
		cmd = exec.Command(task.bin, _Delete.Command, _Delete.taskname, taskname, _Delete.force)
	}

	output, err := cmd.CombinedOutput()
	catch(output, err)

	return string(output)
}

//Query Enables an administrator to display the scheduled tasks on the
//local or remote system.
func (task SchTask) Query(name string, own bool) []Task {
	taskList := make([]Task, 0)

	cmd := &exec.Cmd{}
	if task.compatibility {
		cmd = exec.Command(task.bin, _Query.Command, _Query.format, _Query.formatCSV)
	} else {
		cmd = exec.Command(task.bin, _Query.Command, _Query.format, _Query.formatCSV, _Query.noHeader)
	}

	output, err := cmd.CombinedOutput()
	catch(output, err)

	if own {
		tmp := name
		if name == "*" {
			tmp = ""
		}
		name = task.prefix + tmp
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		tx := strings.Replace(scanner.Text(), "\"", "", -1)

		//skip
		if task.compatibility && strings.HasPrefix(tx, "TaskName") {
			continue
		}

		ts := strings.Split(tx, ",")

		tname := strings.TrimSpace(ts[0])

		if name == "*" || name == "" || strings.Contains(strings.ToLower(tname), strings.ToLower(name)) {
			dtime := strings.TrimSpace(ts[1])
			stat := strings.TrimSpace(ts[2])
			taskList = append(taskList, Task{tname, dtime, stat})
		}
	}

	return taskList
}

//Change Changes the program to run, or user account and password used
//by a scheduled task.
func (task SchTask) Change(taskcreate TaskCreate, own bool) string {
	cmds := task.TaskMake(taskcreate, _Change.Command, own)

	if Debug {
		return dbgMessage
	}

	cmd := exec.Command(task.bin, cmds...)

	output, err := cmd.CombinedOutput()
	catch(output, err)

	return string(output)
}

//Run Runs a scheduled task on demand.
func (task SchTask) Run(taskName string, own bool) string {

	if Debug {
		return dbgMessage
	}

	if own {
		taskName = task.prefix + taskName
	}
	cmd := exec.Command(task.bin, _Run.Command, _Run.taskname, taskName, _Run.immediate)

	output, err := cmd.CombinedOutput()
	catch(output, err)

	return string(output)
}

//End Stops a running scheduled task.
func (task SchTask) End(taskName string, own bool) string {

	if Debug {
		return dbgMessage
	}

	if own {
		taskName = task.prefix + taskName
	}
	cmd := exec.Command(task.bin, _End.Command, _End.taskname, taskName)

	output, err := cmd.CombinedOutput()
	catch(output, err)

	return string(output)
}

//ShowSid Shows the SID for the task's dedicated user.
func (task SchTask) ShowSid(taskName string, own bool) string {

	if Debug {
		return dbgMessage
	}

	if own {
		taskName = task.prefix + taskName
	}
	taskName = "\\" + taskName
	cmd := exec.Command(task.bin, _ShowSid.Command, _ShowSid.taskname, taskName)

	output, err := cmd.CombinedOutput()
	catch(output, err)

	return string(output)
}

//ShowHelp displays help for the command
func (task SchTask) ShowHelp(command string) string {
	cmd := exec.Command(task.bin, command, "/?")

	output, err := cmd.CombinedOutput()
	catch(output, err)

	return string(output)
}

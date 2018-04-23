package tasker

import (
	"fmt"
	"testing"
	"time"
)

var (
	tasker     = New()
	taskName   = "Test"
	executable = "notepad.exe"
)

func TestQuery(t *testing.T) {
	output := tasker.Query("TEST", false)
	fmt.Printf("%+v\n", output)
}

func TestCreate(t *testing.T) {
	timeNow := time.Now().Local()
	timeStr := timeNow.Add(time.Minute).Format("15:04")
	timeStrPlus := timeNow.Add(time.Minute * 2).Format("15:04")

	output := tasker.Create(TaskCreate{
		Taskname:  taskName,
		Taskrun:   executable,
		Starttime: timeStr,
		Terminate: true,
		Endtime:   timeStrPlus,
		Schedule:  Schedules.DAILY,
		Interval:  "0",
	})
	fmt.Printf("%+v\n", output)
}

func TestDelete(t *testing.T) {
	output := tasker.Delete(taskName, true, true)
	fmt.Printf("%+v\n", output)
}

func TestChange(t *testing.T) {
	timeNow := time.Now().Local()
	timeStr := timeNow.Add(time.Minute).Format("15:04")
	timeStrPlus := timeNow.Add(time.Minute * 2).Format("15:04")

	output := tasker.Change(TaskCreate{
		Taskname:  taskName,
		Taskrun:   executable,
		Starttime: timeStr,
		Terminate: true,
		Endtime:   timeStrPlus,
		Interval:  "0",
	}, true)
	fmt.Printf("%+v\n", output)
}

func TestRun(t *testing.T) {
	output := tasker.Run(taskName, true)
	fmt.Printf("%+v\n", output)
}

func TestEnd(t *testing.T) {
	output := tasker.End(taskName, true)
	fmt.Printf("%+v\n", output)
}

func TestShowSid(t *testing.T) {
	output := tasker.ShowSid(taskName, true)
	fmt.Printf("%+v\n", output)
}

func TestShowHelp(t *testing.T) {
	output := tasker.ShowHelp(_Create.Command)
	fmt.Printf("%+v\n", output)
}

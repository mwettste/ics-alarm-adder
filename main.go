package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ics "github.com/arran4/golang-ical"
)

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func addReminderToAllEvents(calendar *ics.Calendar) {
	for _, event := range calendar.Events() {
		alarm := event.AddAlarm()
		alarm.SetAction(ics.ActionDisplay)
		alarm.SetTrigger("-PT30M")
	}
}

func hasExistingAlarmsOnEvents(calendar *ics.Calendar) bool {
	for _, v := range calendar.Events() {
		if len(v.Alarms()) != 0 {
			return true
		}
	}

	return false
}

func userWantsToAddAdditionalAlarms() bool {
	var input string
	for {
		fmt.Println("Found existing alarms, do you want to add additional ones? (y/n):")
		fmt.Scanln(&input)

		if input == "y" {
			return true
		}

		if input == "n" {
			return false
		}
	}
}

func main() {
	var filename string
	flag.StringVar(&filename, "f", "", "The ICS file to add alarms to.")

	flag.Parse()

	if filename == "" {
		fmt.Printf("Please provide the filename with '-f thefile.ics'\n")
		return
	}

	file, err := os.OpenFile(filename, os.O_CREATE, os.ModeAppend)

	if err != nil {
		fmt.Printf("Error opening file '%v': %v\n", filename, err)
		return
	}

	defer file.Close()
	calendar, err := ics.ParseCalendar(file)

	if err != nil {
		fmt.Printf("Could not parse calendar: %v\n", err)
		return
	}

	if hasExistingAlarmsOnEvents(calendar) {
		if !userWantsToAddAdditionalAlarms() {
			fmt.Println("Aborting program without adding events...")
			return
		}
	}

	addReminderToAllEvents(calendar)

	newfile, _ := os.Create(fmt.Sprintf("%v-with-notifications.ics", fileNameWithoutExtension(file.Name())))
	newfile.WriteString(calendar.Serialize())
	newfile.Close()

	fmt.Printf("Added alarms to your events, you can find the new ics file here: %v\n", newfile.Name())
}

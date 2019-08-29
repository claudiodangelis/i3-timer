package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
)

var alarmCommandFlag = flag.String("alarm-command", "",
	"command to be executed when the alarm fires")
var debugFlag = flag.Bool("debug", false, "print debug messages")
var colorsFlag = flag.Bool("colors", false, "colorized timer")
var autostartFlag = flag.Bool("autostart", false, "start timer automatically")
var recurrentFlag = flag.Bool("recurrent", false, "re-start timer after alarm")
var durationFlag = flag.Int("duration", 5, "default duration in minutes")

func debug(args ...interface{}) {
	if *debugFlag {
		log.Println(args...)
	}
}

// Timer holds data of the current timer
type Timer struct {
	Duration    time.Duration `json:"duration"`
	StartTime   time.Time     `json:"startTime"`
	ShowElapsed bool          `json:"showElapsed"`
}

// Remaining returns the remaining duration
func (t *Timer) Remaining() time.Duration {
	var elapsed time.Duration
	if t.IsRunning() {
		elapsed = time.Since(t.StartTime)
	}
	return (t.Duration - elapsed).Truncate(time.Duration(time.Second))
}

// IsLessThanOneMinute returns true if more than one minute is still missing
func (t *Timer) IsLessThanOneMinute() bool {
	return time.Since(t.StartTime)-t.Duration < time.Minute
}

// IsRunning check if the timer is running
func (t *Timer) IsRunning() bool {
	return !time.Time.IsZero(t.StartTime)
}

// IsNotRunning check if the timer is not running
func (t *Timer) IsNotRunning() bool {
	return !t.IsRunning()
}

// AddMinute adds one minute to the duration of the timer
func (t *Timer) AddMinute() {
	t.Duration += time.Duration(time.Minute)
	t.Save()
}

// RemoveMinute removes one minute from the duration of the timer
func (t *Timer) RemoveMinute() {
	if t.Duration > time.Duration(0) {
		t.Duration -= time.Duration(time.Minute)
		if t.Duration < 0 {
			t.Duration = 0
		}
		t.Save()
	}
}

// Alarm executes what has been passed as `-alarm-command`
func (t *Timer) Alarm() {
	if *alarmCommandFlag != "" {
		cmd := exec.Command("sh", "-c", *alarmCommandFlag)
		cmd.Start()
	}
}

// Reset the timer
func (t *Timer) Reset() {
	t.StartTime = time.Time{}
	t.Save()
}

// Start the timer
func (t *Timer) Start() {
	t.StartTime = time.Now()
	t.Save()
}

// ToggleView will toggle between showing elapsed and remaining time.
func (t *Timer) ToggleView() {
	if t.ShowElapsed {
		t.ShowElapsed = false
	} else {
		t.ShowElapsed = true
	}

	t.Save()
}

// String formats the timer
func (t *Timer) String() string {
	var elapsed time.Duration
	if t.IsRunning() {
		elapsed = time.Since(t.StartTime)
	}
	markupStart := ""
	markupEnd := ""
	if *colorsFlag && t.IsRunning() {
		// Get elapsed
		r := t.Remaining()
		if r < t.Duration/4 {
			// If remaining is < 25% of duration, print it red
			markupStart = "<span color='red'>"
		} else if r < t.Duration/2 {
			// If remaining is < 50% of duration, print it yellow
			markupStart = "<span color='yellow'>"
		} else if r > t.Duration/2 {
			// If remaining is > 50% of duration, print it green
			markupStart = "<span color='#00ff00'>"
		}
		markupEnd = "</span>"
	}

	// Default to showing remaining time.
	timerValue := (t.Duration - elapsed).Truncate(time.Duration(time.Second))

	// Show elapsed if timer have been toggled.
	if t.ShowElapsed {
		timerValue = elapsed.Truncate(time.Duration(time.Second))
	}

	return fmt.Sprintf("%sTimer: %s%s",
		markupStart,
		timerValue,
		markupEnd)
}

// Save the content of the timer
func (t *Timer) Save() error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	configFile := filepath.Join(currentUser.HomeDir, ".i3-timer.json")
	j, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, j, 0644)
}

// Button represents the mouse button click
type Button string

// LoadTimer reads timer data from file
func LoadTimer() (Timer, error) {
	var t Timer
	// Check if timer exists
	currentUser, err := user.Current()
	if err != nil {
		return t, err
	}
	configFile := filepath.Join(currentUser.HomeDir, ".i3-timer.json")
	// If config file does not exist yet, create a default one
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Duration = time.Duration(*durationFlag) * time.Minute
		j, err := json.Marshal(t)
		if err != nil {
			return t, err
		}
		if err := ioutil.WriteFile(configFile, j, 0644); err != nil {
			return t, err
		}
	} else {
		// TODO: Check invalid timer
		b, err := ioutil.ReadFile(configFile)
		if err != nil {
			return t, err
		}
		if err := json.Unmarshal(b, &t); err != nil {
			return t, err
		}
	}
	return t, nil
}

const (
	// LeftButton toggles the timer display (elapsed/remaining)
	LeftButton Button = "1"
	// MiddleButton starts the timer
	MiddleButton Button = "2"
	// RightButton stops the timer
	RightButton Button = "3"
	// ScrollUpButton adds 1 minute to duration
	ScrollUpButton Button = "4"
	// ScrollDownButton removes 1 minute to duration
	ScrollDownButton Button = "5"
)

func main() {
	flag.Parse()
	timer, err := LoadTimer()
	if err != nil {
		panic(err)
	}
	switch Button(os.Getenv("BLOCK_BUTTON")) {
	case LeftButton:
		// Toggle elapsed/remaining.
		if timer.IsRunning() {
			timer.ToggleView()
		}
	case MiddleButton:
		// Start the timer if not started yet
		if time.Time.IsZero(timer.StartTime) {
			timer.Start()
		}
	case RightButton:
		// Stop the timer if it's started
		if timer.IsRunning() {
			timer.Reset()
		}
	case ScrollUpButton:
		if timer.IsNotRunning() {
			timer.AddMinute()
		}
	case ScrollDownButton:
		if timer.IsNotRunning() {
			timer.RemoveMinute()
		}
	}
	if timer.IsRunning() && timer.Remaining() <= 0 {
		timer.Alarm()
		timer.Reset()
		if *recurrentFlag {
			timer.Start()
		}
	}
	if !timer.IsRunning() && *autostartFlag {
		timer.Start()
	}

	fmt.Println(timer.String())
}

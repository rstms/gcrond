package cmd

import (
	"bufio"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"syscall"
)

var PATTERN = regexp.MustCompile(`(\S+\s+\S+\s+\S+\s+\S+\s+\S+)\s+(\S.*)$`)
var SECONDS_PATTERN = regexp.MustCompile(`(\S+\s+\S+\s+\S+\s+\S+\s+\S+\s+\S+)\s+(\S.*)$`)

func cronExec(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Printf("gcrond exec failed: '%v': %v\n", args, err)
	}
}

type Cron struct {
	scheduler gocron.Scheduler
}

func NewCron(crontab, exec, shell, flag string, seconds bool) (*Cron, error) {
	if runtime.GOOS == "windows" {
		if shell == "" {
			shell = os.Getenv("COMSPEC")
		}
		if flag == "" {
			flag = "/C"
		}
	} else {
		if shell == "" {
			shell = "/bin/sh"
		}
		if flag == "" {
			flag = "-c"
		}
	}
	cron, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed creating scheduler: %v", err)
	}
	if exec != "" {
		err := addJob(cron, shell, flag, exec, seconds)
		if err != nil {
			return nil, err
		}
	} else {
		file, err := os.Open(crontab)
		if err != nil {
			return nil, fmt.Errorf("failed opening crontab: %v", err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			err := addJob(cron, shell, flag, line, seconds)
			if err != nil {
				return nil, err
			}

		}
		err = scanner.Err()
		if err != nil {
			return nil, fmt.Errorf("failed reading crontab %s: %v", crontab, err)
		}
	}
	return &Cron{scheduler: cron}, nil
}

func addJob(cron gocron.Scheduler, shell, flag, line string, seconds bool) error {
	pattern := PATTERN
	if seconds {
		pattern = SECONDS_PATTERN
	}
	fields := pattern.FindStringSubmatch(line)
	if len(fields) != 3 {
		return fmt.Errorf("crontab syntax error: '%s'", line)
	}
	args := []string{shell, flag, fields[2]}
	_, err := cron.NewJob(gocron.CronJob(fields[1], seconds), gocron.NewTask(cronExec, args))
	if err != nil {
		return fmt.Errorf("failed creating job: '%s': %v", line, err)
	}
	return nil

}

func (c *Cron) Run() error {
	c.scheduler.Start()
	defer c.scheduler.Shutdown()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	<-sigint
	return nil
}

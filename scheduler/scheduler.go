package scheduler

import (
	"github.com/lvkeliang/corn/parser"
	"sort"
	"sync"
	"time"
)

type Task struct {
	Expression *parser.CronExpression
	Job        func()
}

type Cron struct {
	tasks []*Task
	stop  chan bool
	wg    sync.WaitGroup
}

func NewCron() *Cron {
	return &Cron{
		stop: make(chan bool),
	}
}

func (c *Cron) AddTask(expression *parser.CronExpression, job func()) {
	c.tasks = append(c.tasks, &Task{Expression: expression, Job: job})
}

func (c *Cron) Start() {
	for _, task := range c.tasks {
		c.wg.Add(1)
		go func(task *Task) {
			defer c.wg.Done()
			for {
				nextTime := c.getNextTime(task.Expression)

				select {
				case <-time.After(nextTime.Sub(time.Now())):
					task.Job()
				case <-c.stop:
					return
				}
			}
		}(task)
	}
}

func (c *Cron) Stop() {
	close(c.stop)
	c.wg.Wait()
}

func (c *Cron) getNextTime(expression *parser.CronExpression) time.Time {
	now := time.Now()
	currentYear := now.Year()
	for year := currentYear; year < currentYear+5; year++ {
		for _, month := range expression.Month.Value {
			for _, day := range expression.DayOfMonth.Value {
				if day > 28 && month == 2 && !isLeapYear(year) {
					continue
				}
				if day > 30 && (month == 4 || month == 6 || month == 9 || month == 11) {
					continue
				}
				date := time.Date(year, time.Month(month), day+1, 0, 0, 0, 0, time.Local)
				if date.Before(now) {
					continue
				}
				if !contains(expression.DayOfWeek.Value, int(date.Weekday())) {
					continue
				}
				for _, hour := range expression.Hour.Value {
					for _, minute := range expression.Minute.Value {
						for _, second := range expression.Second.Value {
							date := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)

							if date.Before(now) {
								continue
							}

							return date
						}
					}
				}
			}
		}
	}
	return time.Now().AddDate(10, 0, 0)
}

func isLeapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}

func contains(slice []int, value int) bool {
	sort.Ints(slice)
	index := sort.SearchInts(slice, value)
	return index < len(slice) && slice[index] == value
}

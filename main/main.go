package main

import (
	"fmt"
	"github.com/lvkeliang/corn/parser"
	"github.com/lvkeliang/corn/scheduler"
	"time"
)

//func main() {
//	// expression := "1-5,10-15/2 * 8-22/3 ? * MON-FRI"
//	// expression := "*/10 6-8 * * * *"
//	expression := "1 1 8-22/3 ? * MON-FRI"
//	ce := parser.NewCronExpression()
//	err := ce.Parse(expression)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Printf("%+v\n", ce)
//}

func main() {
	expression1 := parser.NewCronExpression()
	err := expression1.Parse("*/5 * * * * *")
	if err != nil {
		fmt.Println(err)
		return
	}

	expression2 := parser.NewCronExpression()
	err = expression2.Parse("*/10 * * * * *")
	if err != nil {
		fmt.Println(err)
		return
	}

	cron := scheduler.NewCron()

	cron.AddTask(expression1, func() {
		fmt.Println("Task 1")
	})

	cron.AddTask(expression2, func() {
		fmt.Println("Task 2")
	})

	cron.Start()

	time.Sleep(time.Minute)

	cron.Stop()
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/KMConner/kyodai-go/internal/auth"
	"github.com/KMConner/kyodai-go/internal/kulasis"
	"github.com/jessevdk/go-flags"
	"os"
	"strconv"
	"strings"
	"time"
)

type timeslotOptions struct {
	defaultOptions
	Term string `short:"s" long:"semester" choice:"first" choice:"second" required:"true"`
}

func (opt *timeslotOptions) Execute(_ []string) error {
	authInfo := opt.GetInfo()
	timeSlot, err := kulasis.RetrieveTimeSlot(authInfo)
	if err != nil {
		return err
	}
	var semester kulasis.Semester
	if opt.Term == "first" {
		semester = kulasis.First
	} else {
		semester = kulasis.Second
	}
	for d := 1; d <= 5; d++ {
		for p := 1; p <= 5; p++ {
			dp := kulasis.DayPeriod{
				Semester: semester,
				Day:      time.Weekday(d),
				Period:   p,
			}
			lecture := timeSlot.GetLecture(dp)
			if lecture != nil {
				fmt.Printf("[%s %d] %s\n", dp.Day.String(), dp.Period, lecture.LectureName)
			}
		}
	}
	return nil
}

type getMailOptions struct {
	defaultOptions
	GetNew bool `short:"n" long:"long"`
}

func (opt *getMailOptions) Execute(_ []string) error {
	authInfo := opt.GetInfo()
	timeSlot, err := kulasis.RetrieveTimeSlot(authInfo)
	if err != nil {
		return err
	}
	var lectures []*kulasis.Lecture
	if opt.GetNew {
		lectures = timeSlot.GetNewLecture()
	} else {
		lectures = timeSlot.GetAllLectures()
	}

	for i, l := range lectures {
		fmt.Printf("%d: %s\n", i+1, l.LectureName)
	}
	println("Select lectures to read course mail.")

	reader := bufio.NewReader(os.Stdin)
	numStr, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	num, err := strconv.Atoi(strings.TrimSpace(numStr))
	if err != nil {
		return err
	}
	if num < 1 || num > len(lectures) {
		return errors.New("INVALID SELECTION")
	}
	lecture := lectures[num-1]

	titles, err := lecture.GetCourseMailTitles()
	if err != nil {
		return err
	}

	for _, t := range *titles {
		mail, err := t.GetContent()
		if err != nil {
			return err
		}

		fmt.Printf("[%s] - %s\n%s\n##########\n\n", mail.Title, mail.Date, mail.TextBody)
	}
	return nil
}

type defaultOptions struct {
	AccountId string `short:"a" long:"account" required:"true" env:"ACCOUNT_ID"`
	Token     string `short:"t" long:"token" required:"true" env:"ACCESS_TOKEN"`
}

func (opt *defaultOptions) GetInfo() auth.Info {
	authInfo := auth.Info{
		AccessToken: opt.Token,
		Account:     opt.AccountId,
	}
	return authInfo
}

func main() {
	defaults := defaultOptions{}
	parser := flags.NewParser(&defaults, flags.Default)
	timeslot := timeslotOptions{}
	mail := getMailOptions{}
	_, e := parser.AddCommand("timeslot", "Show timeslot",
		"Print time slot to console", &timeslot)
	if e != nil {
		println(e.Error())
		return
	}

	_, e = parser.AddCommand("mail", "Get mails",
		"Get mail", &mail)
	if e != nil {
		println(e.Error())
		return
	}

	_, e = parser.Parse()
	if e != nil {
		return
	}
}

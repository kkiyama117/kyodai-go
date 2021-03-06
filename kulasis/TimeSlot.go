package kulasis

import (
	"fmt"
	"net/url"
	"time"
)

type Semester int

const (
	First Semester = iota
	FirstIntensive
	Second
	SecondIntensive
	FullYear
	FullYearIntensive
)

type lectureId struct {
	DepartmentNo int
	LectureNo    int
}

type Lecture struct {
	DepartmentName string
	DepartmentNo   int
	IsNew          bool
	LectureName    string
	LectureNo      int
	RoomName       string
	TeacherName    string
	info           *Info
}

type DayPeriod struct {
	Semester Semester
	Day      time.Weekday
	Period   int
}

type TimeSlot struct {
	times    map[DayPeriod]lectureId
	lectures map[lectureId]*Lecture
}

func RetrieveTimeSlot(info Info) (*TimeSlot, error) {
	var timeSlotRaw timeSlotRaw
	timeslotUrl, err := url.Parse("https://www.k.kyoto-u.ac.jp/api/app/v1/timetable/get_table")
	if err != nil {
		return nil, err
	}

	err = accessWithToken(*timeslotUrl, &info, &timeSlotRaw)
	if err != nil {
		return nil, err
	}
	timeSlot := timeSlotRaw.toTimeSlot()
	for _, v := range timeSlot.lectures {
		v.info = &info
	}
	return timeSlot, nil
}

func (slot *TimeSlot) GetLecture(dp DayPeriod) *Lecture {
	if v, ok := slot.times[dp]; ok {
		if l, ok := slot.lectures[v]; ok {
			return l
		}
	}
	return nil
}

func (slot *TimeSlot) GetAllLectures() []*Lecture {
	var ret []*Lecture
	for _, v := range slot.lectures {
		ret = append(ret, v)
	}
	return ret
}

func (slot *TimeSlot) GetNewLecture() []*Lecture {
	var ret []*Lecture
	for _, v := range slot.lectures {
		if v.IsNew {
			ret = append(ret, v)
		}
	}
	return ret
}

func (lec *Lecture) GetCourseMailTitles() (*[]CourseMailTitle, error) {
	var mails courseMailCollectionRaw
	mailListUrl, err := url.Parse(fmt.Sprintf(
		"https://www.k.kyoto-u.ac.jp/api/app/v1/support/course_mail_list?departmentNo=%d&lectureNo=%d",
		lec.DepartmentNo, lec.LectureNo))

	if err != nil {
		return nil, err
	}

	err = accessWithToken(*mailListUrl, lec.info, &mails)
	if err != nil {
		return nil, err
	}

	ret := make([]CourseMailTitle, 0)
	for _, m := range mails.CourseMails {
		m.info = lec.info
		ret = append(ret, m)
	}
	return &ret, nil
}

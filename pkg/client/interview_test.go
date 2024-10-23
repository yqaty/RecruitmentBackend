package client

import (
	"UniqueRecruitmentBackend/pkg"
	"testing"
	"time"
)

const localAddr = "http://localhost:3333"
const devAddr = "https://dev.back.recruitment2023.hustunique.com"

const recruitmentID = "320772a9-a362-437b-a689-3eafa46625e7"

func TestCreateInterviews(t *testing.T) {
	cli, _ := NewClient(&Opts{Addr: localAddr})
	err := cli.CreateInterview([]pkg.CreateInterviewOpts{{
		Date:   time.Date(2024, 5, 1, 0, 0, 0, 0, time.Local),
		Period: pkg.Morning,
		Start:  time.Date(2024, 5, 1, 8, 0, 0, 0, time.Local),
		End:    time.Date(2024, 5, 1, 12, 0, 0, 0, time.Local),
	}, {
		Date:   time.Date(2024, 5, 1, 0, 0, 0, 0, time.Local),
		Period: pkg.Afternoon,
		Start:  time.Date(2024, 5, 1, 14, 0, 0, 0, time.Local),
		End:    time.Date(2024, 5, 1, 18, 0, 0, 0, time.Local),
	}, {
		Date:   time.Date(2024, 5, 2, 0, 0, 0, 0, time.Local),
		Period: pkg.Morning,
		Start:  time.Date(2024, 5, 2, 8, 0, 0, 0, time.Local),
		End:    time.Date(2024, 5, 2, 12, 0, 0, 0, time.Local),
	}, {
		Date:   time.Date(2024, 5, 2, 0, 0, 0, 0, time.Local),
		Period: pkg.Afternoon,
		Start:  time.Date(2024, 5, 2, 14, 0, 0, 0, time.Local),
		End:    time.Date(2024, 5, 2, 18, 0, 0, 0, time.Local),
	},
	}, recruitmentID, pkg.Web)

	if err != nil {
		t.Fatal(err)
	}
	t.Log("success create interviews")
}

func TestDeleteInterviews(t *testing.T) {
	cli, _ := NewClient(&Opts{Addr: localAddr})
	err := cli.DeleteInterviews([]pkg.DeleteInterviewOpts{{"70ad3db8-b229-401e-8a7d-9edd8e0cec99"}, {"6ae32ed8-e61a-418d-b5a3-4a7abca26bfc"}}, recruitmentID, pkg.Web)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success delete interviews")
}

func TestListnterviews(t *testing.T) {
	cli, _ := NewClient(&Opts{Addr: localAddr})
	resp, err := cli.ListInterviews(recruitmentID, pkg.Web)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", resp)
}

func TestSetApplicationInterviewTime(t *testing.T) {

}

func TestSelectInterviewSlots(t *testing.T) {

}

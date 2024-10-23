package client

import (
	"testing"
	"time"

	"UniqueRecruitmentBackend/pkg"
)

func TestCreateRecruitment(t *testing.T) {
	cli, _ := NewClient(&Opts{Addr: localAddr})
	resp, err := cli.CreateRecruitment(&pkg.CreateRecOpts{
		Name:      "2024夏",
		Beginning: time.Date(2024, 5, 1, 0, 0, 0, 0, time.Local),
		Deadline:  time.Date(2024, 6, 1, 0, 0, 0, 0, time.Local),
		End:       time.Date(2024, 7, 1, 0, 0, 0, 0, time.Local),
	})

	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestListRecruitment(t *testing.T) {
	cli, _ := NewClient(&Opts{Addr: localAddr})
	resp, err := cli.ListRecruitment("320772a9-a362-437b-a689-3eafa46625e7")

	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestGetLatestRecruitment(t *testing.T) {
	cli, _ := NewClient(&Opts{Addr: localAddr})
	resp, err := cli.GetLastestRecruitment()

	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

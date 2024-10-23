package client

import (
	"UniqueRecruitmentBackend/pkg"
	"testing"
)

func TestCreateApplication(t *testing.T) {
	opts := []pkg.CreateAppOpts{
		{
			Grade:         "dayi",
			Institute:     "计算机科学与技术学院",
			Major:         "计算机科学与技术",
			Rank:          "1",
			Group:         "web",
			Intro:         "我是一个大一的学生",
			RecruitmentID: recruitmentID,
			Referrer:      "张三",
			IsQuick:       true,
			//Resume:
		},
		{
			Grade:         "daer",
			Institute:     "人工智能学院",
			Major:         "人工智能",
			Rank:          "2",
			Group:         "lab",
			Intro:         "我是一个大二的学生",
			RecruitmentID: recruitmentID,
			Referrer:      "李四",
			IsQuick:       false,
		},
		{
			Grade:         "dasan",
			Institute:     "软件学院",
			Major:         "软件工程",
			Rank:          "3",
			Group:         "ai",
			Intro:         "我是一个大三的学生",
			RecruitmentID: recruitmentID,
			Referrer:      "王五",
			IsQuick:       true,
		},
		{
			Grade:         "dasi",
			Institute:     "新闻传播学院",
			Major:         "广播电视新闻",
			Rank:          "4",
			Group:         "web",
			Intro:         "我是一个大四的学生",
			RecruitmentID: recruitmentID,
			Referrer:      "赵六",
			IsQuick:       false,
		},
		{
			Grade:         "dawu",
			Institute:     "计算机科学与技术学院",
			Major:         "计算机科学与技术",
			Rank:          "5",
			Group:         "lab",
			Intro:         "我是一个大五的学生",
			RecruitmentID: recruitmentID,
			Referrer:      "孙七",
			IsQuick:       true,
		},
	}

	cli, _ := NewClient(&Opts{Addr: localAddr})
	for _, opt := range opts {
		_, err := cli.CreateApplication(&opt)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestListApplications(t *testing.T) {
	cli, _ := NewClient(&Opts{Addr: localAddr})
	resp, err := cli.ListApplication("c549b215-3b27-4978-b537-2228b81848fb")

	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

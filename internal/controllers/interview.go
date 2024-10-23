package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"

	"UniqueRecruitmentBackend/internal/common"
	"UniqueRecruitmentBackend/internal/models"
	"UniqueRecruitmentBackend/internal/utils"
	"UniqueRecruitmentBackend/pkg"
	"UniqueRecruitmentBackend/pkg/grpc"
)

// GetRecruitmentInterviews get recruitment interviews
// @Id get_recruitment_interviews.
// @Summary get recruitment interviews.
// @Description get recruitment interviews, candidate can't see slotNumber and selectNumber of Interviews (will get interviews of groups or unique), guarantee slotNumber > selectNumber
// @Tags interviews
// @Accept  json
// @Produce  json
// @Param	rid path string true "recruitment id"
// @Param 	name path pkg.Group true "pkg.Group"
// @Success 200 {object} common.JSONResult{data=[]pkg.Interview} ""
// @Failure 400 {object} common.JSONResult{} "code is not 0 and msg not empty"
// @Router /recruitments/{rid}/interviews/{name} [get]
func GetRecruitmentInterviews(c *gin.Context) {
	var (
		interviews []pkg.Interview
		err        error
	)

	defer func() { common.Resp(c, interviews, err) }()

	opts := &pkg.GetInterviewsOpts{}
	if err = c.ShouldBindUri(opts); err != nil {
		return
	}
	if err = opts.Validate(); err != nil {
		return
	}

	interviews, err = models.GetInterviewsByRidAndNameWithoutApp(opts.Rid, opts.Name)
	return
}

// CreateRecruitmentInterviews create recruitment interviews
// @Id create_recruitment_interviews.
// @Summary create recruitment interviews.
// @Description create recruitment interviews.
// @Tags interviews
// @Accept  json
// @Produce  json
// @Param	rid path string true "recruitment id"
// @Param 	name path pkg.Group true "pkg.Group"
// @Param	[]pkg.CreateInterviewOpts body []pkg.CreateInterviewOpts true "create interview opts"
// @Success 200 {object} common.JSONResult{} ""
// @Failure 400 {object} common.JSONResult{} "code is not 0 and msg not empty"
// @Router /recruitments/{rid}/interviews/{name} [post]
func CreateRecruitmentInterviews(c *gin.Context) {
	var (
		opts []pkg.CreateInterviewOpts
		user *pkg.UserDetail
		err  error
	)

	defer func() { common.Resp(c, nil, err) }()

	if err = c.ShouldBind(&opts); err != nil {
		return
	}

	rid := c.Param("rid")
	name := pkg.Group(c.Param("name"))
	name = pkg.Group(strings.ToLower(string(name)))

	if rid == "" {
		err = fmt.Errorf("request param wrong, you should set rid")
		return
	}
	if _, ok := pkg.GroupMap[name]; !ok {
		err = fmt.Errorf("request param wrong, name set wrong")
		return
	}

	user, err = grpc.GetUserInfoByUID(common.GetUID(c))
	if err != nil {
		return
	}

	// member can only update his group's interview or team interview (组面/群面
	if err = checkInterviewGroupName(user, name); err != nil {
		return
	}

	err = models.CreateInterviews(opts, name, rid)
	return
}

// DeleteRecruitmentInterviews delete recruitment interviews
// @Id delete_recruitment_interviews.
// @Summary delete recruitment interviews.
// @Description delete recruitment interviews.
// @Tags interviews
// @Accept  json
// @Produce  json
// @Param	rid path string true "recruitment id"
// @Param 	name path pkg.Group true "pkg.Group or unique"
// @Param	[]pkg.DeleteInterviewOpts body []pkg.DeleteInterviewOpts true "delete interview opts"
// @Success 200 {object} common.JSONResult{} ""
// @Failure 400 {object} common.JSONResult{} "code is not 0 and msg not empty"
// @Router /recruitments/{rid}/interviews/{name} [delete]
func DeleteRecruitmentInterviews(c *gin.Context) {
	var (
		opts []pkg.DeleteInterviewOpts
		user *pkg.UserDetail
		err  error
	)

	defer func() { common.Resp(c, nil, err) }()

	if err = c.ShouldBind(&opts); err != nil {
		return
	}

	rid := c.Param("rid")
	name := pkg.Group(c.Param("name"))
	if rid == "" {
		err = fmt.Errorf("request param wrong, you should set rid")
		return
	}
	if _, ok := pkg.GroupMap[name]; !ok {
		err = fmt.Errorf("request param wrong, name set wrong")
		return
	}

	user, err = grpc.GetUserInfoByUID(common.GetUID(c))
	if err != nil {
		return
	}

	// member can only update his group's interview or team interview (组面/群面
	if err = checkInterviewGroupName(user, name); err != nil {
		return
	}

	err = models.DeleteInterviews(opts, name, rid)
	return
}

// SetRecruitmentInterviews set recruitment interviews
// @Id set_recruitment_interviews.
// @Summary set recruitment interviews.
// @Description get recruitment interviews, use PUt method to prevent resource are duplicated
// @Tags interviews
// @Accept  json
// @Produce  json
// @Param	rid path string true "recruitment id"
// @Param 	name path pkg.Group true "pkg.Group or unique"
// @Param	[]pkg.UpdateInterviewOpts body []pkg.UpdateInterviewOpts true "update interview info"
// @Success 200 {object} common.JSONResult{} ""
// @Failure 400 {object} common.JSONResult{} "code is not 0 and msg not empty"
// @Router /recruitments/{rid}/interviews/{name} [put]
func SetRecruitmentInterviews(c *gin.Context) {
	var (
		r                *pkg.Recruitment
		user             *pkg.UserDetail
		originInterviews []pkg.Interview
		err              error
	)

	defer func() { common.Resp(c, nil, err) }()

	rid := c.Param("rid")
	name := pkg.Group(c.Param("name"))
	uid := common.GetUID(c)
	if rid == "" {
		err = fmt.Errorf("request param wrong, you should set rid")
		return
	}
	if _, ok := pkg.GroupMap[name]; !ok {
		err = fmt.Errorf("request param wrong, name set wrong")
		return
	}

	var interviews []pkg.UpdateInterviewOpts
	if err = c.ShouldBind(&interviews); err != nil {
		return
	}

	// judge whether the recruitment has expired
	r, err = models.GetRecruitmentById(rid)
	if err != nil {
		return
	}
	if r.End.Before(time.Now()) {
		err = fmt.Errorf("recruitment %s has already ended", r.Name)
		return
	}

	user, err = grpc.GetUserInfoByUID(uid)
	if err != nil {
		return
	}

	// member can only update his group's interview or team interview (组面/群面
	if err = checkInterviewGroupName(user, name); err != nil {
		return
	}

	var interviewsToAdd []pkg.Interview
	var interviewIdsToDel []string
	interviewsToUpdate := make(map[string]pkg.Interview)
	for _, interview := range interviews {
		if interview.Uid != "" {
			// update
			interviewsToUpdate[interview.Uid] = pkg.Interview{
				Common: pkg.Common{
					Uid: interview.Uid,
				},
				Name:          name,
				RecruitmentID: rid,
				Date:          interview.Date,
				Period:        interview.Period,
				Start:         interview.Start,
				End:           interview.End,
			}
		} else {
			// add
			interviewsToAdd = append(interviewsToAdd, pkg.Interview{
				Name:          name,
				RecruitmentID: rid,
				Date:          interview.Date,
				Period:        interview.Period,
				Start:         interview.Start,
				End:           interview.End,
			})
		}
	}

	originInterviews, err = models.GetInterviewsByRidAndName(rid, name)
	if err != nil {
		return
	}

	var iids []string
	for _, origin := range originInterviews {
		iids = append(iids, origin.Uid)
	}
	interviewsCannotBeUpdate, err := models.GetInterviewsCannotBeUpdate(iids)
	zapx.Infof("interviews can not be update: %v", interviewsCannotBeUpdate)

	var errors []string

	for _, origin := range originInterviews {
		interview, toUpdate := interviewsToUpdate[origin.Uid]
		_, cannotBeUpdate := interviewsCannotBeUpdate[origin.Uid]
		if toUpdate {
			// update
			// check whether the interview time has been selected by candidates
			if cannotBeUpdate && !checkUpdateInterview(&origin, &interview) {
				errors = append(errors, fmt.Sprintf("[interview %s have been selected]", origin.Uid))
			} else {
				if errdb := models.UpdateInterview(&interview); errdb != nil {
					errors = append(errors, fmt.Sprintf("[update interviews db failed, err: %s]", errdb.Error()))
				}
			}
		} else {
			// delete
			// check whether the interview time has been selected by candidates
			if cannotBeUpdate {
				errors = append(errors, fmt.Sprintf("[interview %s have been selected]", origin.Uid))
			} else {
				interviewIdsToDel = append(interviewIdsToDel, origin.Uid)
			}
		}
	}

	if errdb := models.AddAndDeleteInterviews(interviewsToAdd, interviewIdsToDel); errdb != nil {
		errors = append(errors, fmt.Sprintf("[add and delete interviews db failed, err: %s]", errdb.Error()))
	}

	if len(errors) != 0 {
		err = fmt.Errorf("%v", errors)
		return
	}
	return
}

// check user's group == name
func checkInterviewGroupName(user *pkg.UserDetail, name pkg.Group) error {
	if name != pkg.Unique {
		if !utils.CheckInGroups(user.Groups, name) {
			err := fmt.Errorf("you can't set other group's interview time")
			return err
		}
	}
	return nil
}

// check if interview times are equal
func checkUpdateInterview(origin *pkg.Interview, interview *pkg.Interview) bool {
	if !origin.Date.Equal(interview.Date) {
		return false
	}
	if !origin.Start.Equal(interview.Date) {
		return false
	}
	if !origin.End.Equal(interview.Date) {
		return false
	}
	if origin.Period != interview.Period {
		return false
	}
	return true
}

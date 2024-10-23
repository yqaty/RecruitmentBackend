package models

import (
	"UniqueRecruitmentBackend/global"
	"UniqueRecruitmentBackend/pkg"
	"fmt"
)

func GetInterviewById(iid string) (*pkg.Interview, error) {
	db := global.GetDB()
	var interview pkg.Interview
	if err := db.Model(&pkg.Interview{}).
		Where("uid = ?", iid).
		First(&interview).Error; err != nil {
		return nil, err
	}
	return &interview, nil
}

func GetInterviewsByRidAndNameWithoutApp(rid string, name pkg.Group) ([]pkg.Interview, error) {
	db := global.GetDB()
	var res []pkg.Interview
	if err := db.Model(&pkg.Interview{}).
		//Omit("\"selectNumber\", \"slotNumber\"").
		//Where("\"selectNumber\" < \"slotNumber\"").
		Where("\"recruitmentId\" = ? AND name = ?", rid, name).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetInterviewsByRidAndNameWithoutAppByMember(rid string, name pkg.Group) ([]pkg.Interview, error) {
	db := global.GetDB()
	var res []pkg.Interview
	if err := db.Model(&pkg.Interview{}).
		Where("\"recruitmentId\" = ? AND name = ?", rid, name).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetInterviewsByRidAndName(rid string, name pkg.Group) ([]pkg.Interview, error) {
	db := global.GetDB()
	var res []pkg.Interview
	if err := db.Model(&pkg.Interview{}).
		Preload("Applications").
		Where("\"recruitmentId\" = ? AND name = ?", rid, name).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetInterviewsByIdsAndName(ids []string, name pkg.Group) ([]pkg.Interview, error) {
	db := global.GetDB()
	var interviews []pkg.Interview
	if err := db.Where("uid in ?", ids).
		Where("name = ?", name).
		Find(&interviews).Error; err != nil {
		return nil, err
	}
	return interviews, nil
}

func UpdateInterview(interview *pkg.Interview) error {
	db := global.GetDB()
	if err := db.Model(&pkg.Interview{}).
		Where("\"uid\" = ?", interview.Uid).
		Updates(map[string]interface{}{
			"date":   interview.Date,
			"period": interview.Period,
			"start":  interview.Start,
			"end":    interview.End,
			"name":   interview.Name,
		}).Error; err != nil {
		return err
	}
	return nil
}

func AddAndDeleteInterviews(interviewsToAdd []pkg.Interview, interviewIdsToDel []string) (err error) {
	db := global.GetDB()
	if len(interviewsToAdd) != 0 {
		if errCreate := db.Create(interviewsToAdd).Error; errCreate != nil {
			return errCreate
		}
	}
	if len(interviewIdsToDel) != 0 {
		if errDelete := db.Delete(&pkg.Interview{}, "uid in ?", interviewIdsToDel).Error; errDelete != nil {
			return errDelete
		}
	}
	return
}

func GetInterviewsCannotBeUpdate(iids []string) (map[string]struct{}, error) {
	db := global.GetDB()
	interviewsCannotBeUpdate := make(map[string]struct{})
	res := []string{}

	if len(iids) == 0 {
		return interviewsCannotBeUpdate, nil
	}

	// get the interview uid that has been selected by the application
	if err := db.Table("interview_selections").
		Select("DISTINCT interview_uid").
		Where("interview_uid IN ?", iids).
		Find(&res).Error; err != nil {
		return nil, err
	}
	for _, val := range res {
		interviewsCannotBeUpdate[val] = struct{}{}
	}

	// get the interview uid that has been allocated by the application
	if err := db.Model(&pkg.Application{}).
		Select("DISTINCT \"interviewAllocationsGroupId\"").
		Where("\"interviewAllocationsGroupId\" IN ?", iids).
		Find(&res).Error; err != nil {
		return nil, err
	}
	for _, val := range res {
		interviewsCannotBeUpdate[val] = struct{}{}
	}

	if err := db.Model(&pkg.Application{}).
		Select("DISTINCT \"interviewAllocationsTeamId\"").
		Where("\"interviewAllocationsTeamId\" IN ?", iids).
		Find(&res).Error; err != nil {
		return nil, err
	}
	for _, val := range res {
		interviewsCannotBeUpdate[val] = struct{}{}
	}
	return interviewsCannotBeUpdate, nil
}

func CreateInterviews(opts []pkg.CreateInterviewOpts, name pkg.Group, rid string) (err error) {
	db := global.GetDB()
	var errs []error
	for _, opt := range opts {
		dbErr := db.Model(&pkg.Interview{}).Create(&pkg.Interview{
			RecruitmentID: rid,
			Name:          name,
			Date:          opt.Date,
			Period:        opt.Period,
			Start:         opt.Start,
			End:           opt.End,
		}).Error
		if dbErr != nil {
			errs = append(errs, dbErr)
		}
	}
	if len(errs) != 0 {
		err = fmt.Errorf("%v", errs)
	}
	return
}

func DeleteInterviews(opts []pkg.DeleteInterviewOpts, name pkg.Group, rid string) (err error) {
	db := global.GetDB()
	var errs []error
	for _, opt := range opts {
		var dbErr error
		var res []string
		// get the interview uid that has been selected by the application
		dbErr = db.Table("interview_selections").
			Select("DISTINCT interview_uid").
			Where("interview_uid = ?", opt.Iid).
			Find(&res).Error
		if dbErr != nil {
			errs = append(errs, dbErr)
			continue
		} else if len(res) != 0 {
			errs = append(errs, fmt.Errorf("interview %s have been selected", opt.Iid))
			continue
		}

		// get the interview uid that has been allocated by the application
		dbErr = db.Model(&pkg.Application{}).
			Select("uid").
			Where("\"interviewAllocationsGroupId\" = ? or \"interviewAllocationsTeamId\" = ?", opt.Iid, opt.Iid).
			Find(&res).Error
		if dbErr != nil {
			errs = append(errs, dbErr)
			continue
		} else if len(res) != 0 {
			errs = append(errs, fmt.Errorf("interview %s have been allocated", opt.Iid))
			continue
		}

		dbErr = db.Model(&pkg.Interview{}).
			Where("\"recruitmentId\" = ? AND name = ? AND uid = ?", rid, name, opt.Iid).
			Delete(&pkg.Interview{}).Error
		if dbErr != nil {
			errs = append(errs, dbErr)
		}
	}

	if len(errs) != 0 {
		err = fmt.Errorf("%v", errs)
	}
	return
}

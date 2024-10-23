package models

import (
	"UniqueRecruitmentBackend/global"
	"UniqueRecruitmentBackend/pkg"
	"UniqueRecruitmentBackend/pkg/grpc"
	"encoding/json"
	"errors"
)

func CreateRecruitment(opts *pkg.CreateRecOpts) (r *pkg.Recruitment, err error) {
	db := global.GetDB()
	if db.Model(&pkg.Recruitment{}).
		Where("name = ?", opts.Name).
		Find(r).RowsAffected > 0 {
		return nil, errors.New("recruitment with the same name cannot be created")
	}

	r = &pkg.Recruitment{
		Name:      opts.Name,
		Beginning: opts.Beginning,
		Deadline:  opts.Deadline,
		End:       opts.End,
	}
	err = db.Model(&pkg.Recruitment{}).Create(r).Error
	return
}

func UpdateRecruitment(opts *pkg.UpdateRecOpts) error {
	bytes, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	var r pkg.Recruitment
	if err := json.Unmarshal(bytes, &r); err != nil {
		return err
	}
	r.Uid = opts.Rid

	db := global.GetDB()
	return db.Updates(&r).Error
}

func GetRecruitmentById(rid string) (*pkg.Recruitment, error) {
	db := global.GetDB()
	var r pkg.Recruitment
	if err := db.Model(&pkg.Recruitment{}).
		Where("uid = ?", rid).
		Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func GetFullRecruitmentById(rid string) (*pkg.Recruitment, error) {
	db := global.GetDB()
	var r pkg.Recruitment
	//remember preload need the struct filed name
	var err error
	if err = db.Model(&pkg.Recruitment{}).
		Preload("Applications").
		Preload("Interviews").
		Preload("Applications.InterviewSelections").
		Preload("Applications.Comments").
		Preload("Applications.InterviewAllocationsGroup").
		Preload("Applications.InterviewAllocationsTeam").
		Where("uid = ?", rid).Find(&r).Error; err != nil {
		return nil, err
	}

	for i := range r.Applications {
		r.Applications[i].UserDetail, err = grpc.GetUserInfoByUID(r.Applications[i].CandidateID)
		if err != nil {
			return nil, err
		}
	}
	return &r, err
}

func GetAllRecruitment() ([]pkg.Recruitment, error) {
	db := global.GetDB()
	var r []pkg.Recruitment
	err := db.Model(&pkg.Recruitment{}).
		Order("beginning DESC").
		Find(&r).Error
	return r, err
}

// GetPendingRecruitment get the latest recruitment
func GetPendingRecruitment() (*pkg.Recruitment, error) {
	db := global.GetDB()
	var r pkg.Recruitment
	if err := db.Model(&pkg.Recruitment{}).
		Select("uid").
		Order("beginning DESC").
		Limit(1).
		Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func GetRecruitmentStatistics(rid string) (map[string]int, error) {
	var results []struct {
		Group string
		Count int
	}
	statistics := make(map[string]int)
	db := global.GetDB()
	if err := db.Model(&pkg.Recruitment{}).
		Select("applications.group, count(*) as count").
		Joins("JOIN applications on recruitments.uid = applications.\"recruitmentId\"").
		Where("recruitments.uid = ?", rid).
		Group("applications.group").
		Scan(&results).Error; err != nil {
		return map[string]int{}, err
	}

	for _, result := range results {
		statistics[result.Group] = result.Count
	}

	return statistics, nil
}

func UpdateStressTestTime(opts *pkg.SetStressTestTimeOpts) error {
	db := global.GetDB()
	if err := db.Model(&pkg.Recruitment{}).
		Where("uid = ?", opts.Rid).
		Updates(map[string]interface{}{
			"\"stressTestStart\"": opts.Start,
			"\"stressTestEnd\"":   opts.End,
		}).Error; err != nil {
		return err
	}
	return nil
}

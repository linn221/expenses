package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
	"gorm.io/gorm"
)

type costPerspectiveService struct {
	*app.Container
}

func (s *costPerspectiveService) db(session *app.Session) *gorm.DB {
	if session.Tx != nil {
		return session.Tx
	}
	return s.Container.DB.WithContext(session.Ctx)
}

func (s *costPerspectiveService) Create(session *app.Session, model models.CostPerspective) (*models.CostPerspective, error) {
	return createModel(s.Container, session, model)
}

func (s *costPerspectiveService) Update(session *app.Session, existing *models.CostPerspective, input *models.CostPerspective) (*models.CostPerspective, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()

	result, err := updateModel(s.Container, session, existing, input)
	if err != nil {
		return nil, err
	}

	applications, err := s.replaceApplications(session, existing.Id, input.Applications)
	if err != nil {
		return nil, err
	}
	result.Applications = applications
	return result, session.Commit()
}

func (s *costPerspectiveService) Delete(session *app.Session, record *models.CostPerspective) (*models.CostPerspective, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()

	err := s.db(session).Where("cost_perspective_id = ?", record.Id).Delete(&models.CostPerspectiveApplication{}).Error
	if err != nil {
		return nil, err
	}

	result, err := deleteModel(s.Container, session, record)
	if err != nil {
		return nil, err
	}
	return result, session.Commit()
}

func (s *costPerspectiveService) Get(session *app.Session) (*models.CostPerspective, error) {
	return selectResource[models.CostPerspective](s.Container, session, "Applications")
}

func (s *costPerspectiveService) List(session *app.Session, filter *dtos.CostPerspectiveFilter) ([]models.CostPerspective, error) {
	db := s.Container.DB.WithContext(session.Ctx).Model(&models.CostPerspective{}).Preload("Applications").Where("user_id = ?", session.UserId)
	if filter.Name != "" {
		db.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Type != "" {
		db.Where("`type` = ?", filter.Type)
	}
	if filter.IsActive != nil {
		db.Where("is_active = ?", filter.IsActive)
	}

	var results []models.CostPerspective
	err := db.Find(&results).Error
	return results, err
}

func (s *costPerspectiveService) replaceApplications(session *app.Session, costPerspectiveId int, applications []models.CostPerspectiveApplication) ([]models.CostPerspectiveApplication, error) {
	err := s.db(session).Where("cost_perspective_id = ?", costPerspectiveId).Delete(&models.CostPerspectiveApplication{}).Error
	if err != nil {
		return nil, err
	}
	if len(applications) == 0 {
		return []models.CostPerspectiveApplication{}, nil
	}

	for i := range applications {
		applications[i].Id = 0
		applications[i].CostPerspectiveId = costPerspectiveId
	}
	err = s.db(session).Create(&applications).Error
	return applications, err
}

func (s *costPerspectiveService) HttpResponse() dtos.HttpResponseFromModel[models.CostPerspective] {
	return dtos.RespondCostPerspective
}

func (s *costPerspectiveService) HttpResponseCollection() dtos.HttpResponseFromModelCollection[models.CostPerspective] {
	return dtos.RespondCostPerspectiveCollection
}

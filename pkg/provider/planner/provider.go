package plan

import "github.com/fuzzingbits/hub/pkg/entity"

type Provider interface {
	Get(userUUID string, date string) (entity.Planner, error)
	Save(plan entity.Planner) error
	Delete(plan entity.Planner) error
	GetAll(userUUID string) ([]entity.Planner, error)
}

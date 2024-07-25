package scheduler

import (
	"gpu-cloudsim/models"
)

type Scheduler interface {
	Schedule(containers []*models.Container, hosts []*models.Host) error
}

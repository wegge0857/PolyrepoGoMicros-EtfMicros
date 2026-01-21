package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Etf 是业务领域对象
type Etf struct {
	Id       int64
	Name     string
	TenYCagr float64
	Star     int64
}

// Data 层是 Biz 层接口的具体实现者 （依赖倒置，而非通常的业务依赖数据库）
// 实现逻辑：在data层 返回格式 必须为biz.EtfRepo
type EtfRepo interface {
	FindByID(ctx context.Context, id int64) (*Etf, error)
	UpdateStar(ctx context.Context, id int64, kind string) (int64, error)
}

// EtfUseCase 是业务逻辑主体
type EtfUseCase struct {
	repo EtfRepo
	log  *log.Helper
}

func NewEtfUseCase(repo EtfRepo, logger log.Logger) *EtfUseCase {
	return &EtfUseCase{repo: repo, log: log.NewHelper(logger)}
}

// Get 获取用户业务逻辑
func (uc *EtfUseCase) Get(ctx context.Context, id int64) (*Etf, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *EtfUseCase) UpdateStar(ctx context.Context, id int64, kind string) (int64, error) {
	return uc.repo.UpdateStar(ctx, id, kind)
}
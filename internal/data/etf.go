package data

import (
	"context"
	"etfMicros/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type EtfRepo struct {
	data *Data //  数据源句柄: 指向 internal/data/data.go 中定义的 Data 结构体的指针
	log  *log.Helper
}

// 注入阶段（启动时）： wire 运行 NewEtfRepo。
// 生成的对象被塞进biz层的接口里，最终通过 NewEtfUseCase 注入给 uc.repo。
// Wire 把 Data 层的实例“注入”到了 Biz 层 （依赖注入）
func NewEtfRepo(data *Data, logger log.Logger) biz.EtfRepo { // 这里的参数由Wire去传入
	return &EtfRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 调用阶段（请求时）： 当用户访问接口，Service 调用 uc.Get(ctx, id)。
// 在 uc内部，uc.repo.FindByID(ctx, id) ---> uc.repo 实际上是一个接口 触发 已实现的 FindByID方法
func (r *EtfRepo) FindByID(ctx context.Context, id int64) (*biz.Etf, error) {
	var u Etf
	db := r.data.db.WithContext(ctx) // WithContext(ctx) 确保遵循链路追踪和超时控制
	// 使用 GORM 的 First 方法按主键查询
	if err := db.First(&u, id).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	// 将 Data 层的模型转换为 Biz 层的业务实体
	return &biz.Etf{
		Id:    u.ID,
		Name:  u.EtfName,
		TenYCagr: u.TenYCagr,
	}, nil
}

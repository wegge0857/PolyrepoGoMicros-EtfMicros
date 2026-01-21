package data

import (
	"context"
	"etfMicros/internal/biz"

	"database/sql" // 添加这一行

	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
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

	var e Etf
	db := r.data.db.WithContext(ctx) // WithContext(ctx) 确保遵循链路追踪和超时控制
	// 使用 GORM 的 First 方法按主键查询
	if err := db.First(&e, id).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	// 将 Data 层的模型转换为 Biz 层的业务实体
	return &biz.Etf{
		Id:       e.ID,
		Name:     e.EtfName,
		TenYCagr: e.TenYCagr,
		Star:     e.Star,
	}, nil
}

// 加收藏数
func (r *EtfRepo) UpdateStar(ctx context.Context, id int64, kind string) (int64, error) {
	log.Info("--------->UpdateStar begin")
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return 0, err
	}

	var latestStar int64

	// 特殊情况 判断是否是回退逻辑
	if barrier.Op == dtmimp.OpCompensate || barrier.Op == dtmimp.OpRollback {
		db := r.data.db.WithContext(ctx).Model(&Etf{})
		expr := gorm.Expr("CASE WHEN star > 0 THEN star - 1 ELSE 0 END")

		err := db.Where("id = ?", id).Update("star", expr).Error
		if err != nil {
			return 0, err
		}

		// 注意：这里继续复用 db 实例，但要小心之前的 Where 条件残留，重新指定结果接收变量
		err = r.data.db.WithContext(ctx).Model(&Etf{}).
			Where("id = ?", id).
			Select("star").
			Scan(&latestStar).Error

		return latestStar, err
	}

	// 正常加的逻辑
	sqlDB, err := r.data.db.DB()
	if err != nil {
		return 0, err
	}

	err = barrier.CallWithDB(sqlDB, func(sTx *sql.Tx) error {
		// 1. 【关键】将原生 sTx 包装进 GORM
		// 这样 gdb 就是一个已经开启了事务、且使用了 DTM 屏障连接的 GORM 对象
		gdb, err := gorm.Open(
			mysql.New(mysql.Config{Conn: sTx}),
			&gorm.Config{})
		if err != nil {
			return err
		}

		// 使用 ctx 保证链路追踪
		tx := gdb.WithContext(ctx)

		// 2. 准备 SQL 表达式
		expr := "star + 1"

		// 3. 执行业务逻辑 (完全使用 tx 对象)
		// 更新操作
		if err := tx.Table("etfs").Where("id = ?", id).
			Update("star", gorm.Expr(expr)).Error; err != nil {
			return err
		}

		// 查询操作
		return tx.Table("etfs").Where("id = ?", id).
			Select("star").Scan(&latestStar).Error
	})

	if err != nil {
		return 0, err
	}

	return latestStar, err
}

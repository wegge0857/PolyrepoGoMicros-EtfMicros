package data

import (
	"etfMicros/internal/conf"
	"time"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewEtfRepo)

type Data struct {
	db *gorm.DB // 添加 GORM 句柄
}

type Etf struct {
	ID        int64     `gorm:"primarykey;autoIncrement" json:"id"`
	EtfName   string    `gorm:"column:etf_name;uniqueIndex;type:varchar(100);comment:etf名字" json:"etf_name"`
	TenYCagr  float64   `gorm:"column:ten_y_cagr;default:0;comment:近十年年化收益率" json:"ten_y_cagr"`
	Star      int64     `gorm:"column:star;default:0;comment:用户星级或点赞数" json:"star"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}

// NewData 负责初始化数据库连接
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "etfMicros/data"))
	// 这里使用 gorm.Open 连接数据库
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		l.Error("etfMicros 数据库连接失败...")
		return nil, nil, err
	}

	sqlDB, _ := db.DB()
	// --- 连接池核心配置 ---
	// 1. 设置连接可复用的最大时间 必须小于数据库服务器的 wait_timeout
	sqlDB.SetConnMaxLifetime(time.Minute * 30)
	// 2. 设置闲置连接的最大存活时间，防止拿到数据库已关掉的“死连接”
	sqlDB.SetConnMaxIdleTime(time.Minute * 2)
	// 3. 设置最大闲置连接数
	sqlDB.SetMaxIdleConns(10)
	// 4. 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	l.Info("mysql 连接池 initialized")

	// 自动迁移表结构（可选，方便测试）
	db.AutoMigrate(&Etf{})
	l.Info("etfMicros 数据库连接成功...")
	cleanup := func() {
		l.Info("etfMicros 关闭mysql数据资源...") //Ctrl + C 或 Kill
		// 如果需要显式关闭 DB，可以在这里获取 sql.DB 并 Close
	}

	// 如果你的表就在当前连接的库里，直接写 "barrier" 即可
	dtmcli.SetBarrierTableName("barrier")

	return &Data{db: db}, cleanup, nil
}

package service

import (
	"context"
	"etfMicros/internal/biz"

	etfV1 "github.com/wegge0857/PolyrepoGoMicros-ApiLink/etf/v1"
)

type EtfService struct {
	// 当前你在 proto 里定义的接口，如果你没写，就会自动走 UnimplementedEtfServer 的逻辑，就不会报错
	etfV1.UnimplementedEtfServer
	// 类型嵌入（Type Embedding） 或 匿名组合
	// 即EtfService 自动获得了 UnimplementedEtfServer 的所有方法。EtfService实现了EtfServer接口的所有方法

	// 添加这一行，以便在方法中使用 s.uc
	uc *biz.EtfUseCase
}

// 在参数中加入 uc *biz.EtfUseCase
func NewEtfService(uc *biz.EtfUseCase) *EtfService {
	return &EtfService{
		uc: uc,
	}
}

func (s *EtfService) GetEtf(ctx context.Context, req *etfV1.GetEtfRequest) (*etfV1.GetEtfReply, error) {
	// 调用 biz 层
	u, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &etfV1.GetEtfReply{
		Id:       u.Id,
		Name:     u.Name,
		TenYCagr: u.TenYCagr,
	}, nil
}

func (s *EtfService) UpdateStar(ctx context.Context, req *etfV1.UpdateStarRequest) (*etfV1.UpdateStarReply, error) {
	num, err := s.uc.UpdateStar(ctx, req.Id, req.Kind)
	if err != nil {
		return nil, err
	}
	return &etfV1.UpdateStarReply{
		CurrentStar: num,
	}, nil
}

func (s *EtfService) CreateEtf(ctx context.Context, req *etfV1.CreateEtfRequest) (*etfV1.CreateEtfReply, error) {
	return &etfV1.CreateEtfReply{}, nil
}
func (s *EtfService) UpdateEtf(ctx context.Context, req *etfV1.UpdateEtfRequest) (*etfV1.UpdateEtfReply, error) {
	return &etfV1.UpdateEtfReply{}, nil
}
func (s *EtfService) DeleteEtf(ctx context.Context, req *etfV1.DeleteEtfRequest) (*etfV1.DeleteEtfReply, error) {
	return &etfV1.DeleteEtfReply{}, nil
}
func (s *EtfService) ListEtf(ctx context.Context, req *etfV1.ListEtfRequest) (*etfV1.ListEtfReply, error) {
	return &etfV1.ListEtfReply{}, nil
}

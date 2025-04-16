package biz

import (
	"context"
	"seckil_service/dao/redis"
	"seckil_service/proto"
	"seckil_service/rpc"
	"strconv"
	"time"

	"go.uber.org/zap"
)

//���ö�ʱ�����������ά��Ա�ڻ��ʼǰ�����Ԥ��
//ɸѡ�ȵ����ݣ���ǰ����Ʒ��Ϣ����redis
// ˢ��ָ����Ʒ�Ŀ�浽 Redis

func refreshStockForGoods(ctx context.Context, goodsID int64,userID int64) error{
	stock,err := rpc.StockCli.GetStock(ctx,&proto.GetStockReq{
		GoodsId: goodsID,
	})
	if err != nil {
        zap.L().Error("��ȡ���ʧ��", zap.Error(err), zap.Int64("goodsID", goodsID))
        return err
    }
	// �������Ϣͬ���� Redis
    // �������Ϣͬ���� Redis
    stockKey := "seckill:stock:" + strconv.FormatInt(goodsID, 10) // �޸�����
    err = redis.GetClient().Set(ctx, stockKey, stock.Stock, time.Minute*5).Err()
    if err != nil {
        zap.L().Error("ͬ����浽Redisʧ��", zap.Error(err), zap.Int64("goodsID", goodsID))
        return err
    }

    zap.L().Info("���Ԥ�����", zap.Int64("goodsID", goodsID), zap.Int64("stock", stock.Stock))
    return nil
}
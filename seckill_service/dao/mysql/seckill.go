package mysql

import (
	"context"
	"seckil_service/errno"
	"seckil_service/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DecrementStock �����ݿ��пۼ����
func DecrementStock(ctx context.Context, goodsId int64, quantity int) (int64, error) {
	// ʹ�� gorm �� WithContext �������������Ĵ��ݸ����ݿ����
	result := db.WithContext(ctx).
		Model(&model.Goods{}).
		Where("goods_id = ? AND stock >= ?", goodsId, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	// �������Ƿ�ɹ�
	if result.Error != nil {
		zap.L().Error("���ݿ�ۼ����ʧ��", zap.Error(result.Error))
		return 0, result.Error
	}

	// ���û���б����£����ش���
	if result.RowsAffected == 0 {
		zap.L().Error("��治�����Ʒ������", zap.Int64("goodsId", goodsId))
		return 0, errno.ErrStockInsufficient
	}

	// ����Ӱ�������
	return result.RowsAffected, nil
}

// GetStockByGoodsId ������Ʒ ID ��ѯ�����Ϣ��
func GetStockByGoodsId(ctx context.Context, goodsId int64) (*model.Stock, error) {
	// ��ʼ��һ�� Stock �ṹ�����ڴ洢��ѯ�����
	var data model.Stock

	// ʹ�� GORM ��ѯ����¼��
	err := db.WithContext(ctx).
		Model(&model.Stock{}).          // ָ��������ģ��Ϊ Stock ��
		Where("goods_id = ?", goodsId). // ������Ʒ ID ��ѯ��
		First(&data).                   // ��ȡ��һ����¼��
		Error                           // ��ȡ��ѯ����Ĵ�����Ϣ��

	// �����ѯʧ���Ҵ����Ǽ�¼δ�ҵ����򷵻� ErrQueryFailed��
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errno.ErrQueryFailed
	}

	// ��¼��ѯ�������־��
	zap.L().Info("��ѯ���Ŀ����Ϣ", zap.Any("data", data))

	return &data, nil // ���ز�ѯ�����
}
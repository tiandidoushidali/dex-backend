package authlogic

import (
	"dex/app/api/internal/svc"
	"dex/app/api/internal/types"
	"dex/app/common/constants"
	"dex/app/data/model/comer"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"golang.org/x/net/context"

	"github.com/ethereum/go-ethereum/common"
)

type LoginByWalletAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWalletAddressLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByWalletAddressLogic {
	return &LoginByWalletAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginByWalletAddressLogic) WalletAddressLogin(req *types.WalletAddressLoginReq) (resp *types.JwtAuthResp, err error) {
	// 1.获取nonce
	nonceCacheKey := fmt.Sprintf(constants.RedisAuthWalletLoginNonce, req.WalletAddress)
	nonce, err := l.svcCtx.RedisClient.Get(l.ctx, nonceCacheKey).Result()
	if err != nil {
		if err.Error() == string(redis.Nil) {
			return nil, errors.New(-1, "Nonce not exist")
		}
		logx.WithContext(l.ctx).Errorf("get nonce from redis err: %+v", err)
		return nil, errors.New(-1, "Nonce not exist")
	}
	// 2.校验签名
	err = validateEthWallet(req.WalletAddress, nonce, req.Signature)
	if err != nil {
		return nil, err
	}
	// 查询用户
	comerInfo, err := comer.FindComerByAddress(l.svcCtx.MysqlDB, req.WalletAddress)
	if err != nil {
		logx.Errorf("find comer err: %+v", err)
		return nil, err
	}
	if comerInfo == nil {
		// 4. 创建新用户
		newComer := comer.Comer{
			Address: req.WalletAddress,
		}
		err = comer.InsertComer(l.svcCtx.MysqlDB, &newComer)
		if err != nil {
			return nil, err
		}
		comerInfo = &newComer
	}
	// 6. 删除nonce
	_, err = l.svcCtx.RedisClient.Del(l.ctx, nonceCacheKey).Result()
	if err != nil {
		logx.Errorf("redis remove nonce failed: %+v", err)
	}
	// 7. 生成jwt
	claim := make(jwt.MapClaims)
	claim["comerID"] = comerInfo.ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(l.svcCtx.Config.JwtAuth.AccessSecret)) // 对称加密必须是[]byte
	if err != nil {
		logx.Errorf("sign token err: %+v", err)
		return nil, err
	}

	resp = &types.JwtAuthResp{
		Token: tokenStr,
	}
	return resp, err
}

func validateEthWallet(address, nonce, signature string) error {
	logx.Infof("validateEthWallet address:%s nonce:%s signature:%s", address, nonce, signature)
	addrKey := common.HexToAddress(address)
	//  解析签名（注意：ethers.js 返回的是 65 字节签名，含 r, s, v）
	sig := hexutil.MustDecode(signature)
	// Ethers.js 会返回 V 为 27/28，我们需要减去 27 变成 0/1
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}
	// 加入 Ethereum 签名前缀（EIP-191）
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(nonce), nonce)
	msg256 := crypto.Keccak256([]byte(msg))
	pubKey, err := crypto.SigToPub(msg256, sig)
	if err != nil {
		return err
	}
	recoverAddr := crypto.PubkeyToAddress(*pubKey)
	if recoverAddr != addrKey {
		logx.Errorf("recover addr %s is not match with pub addr: %s", recoverAddr.Hex(), addrKey.Hex())
		return errors.New(-1, "Address mismatch")
	}

	return nil
}

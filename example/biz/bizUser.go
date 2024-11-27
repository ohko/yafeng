package biz

import (
	"errors"
	"example/config"

	"github.com/ohko/yafeng"
)

func UserLogin(ctx *yafeng.Context, account, password string) (*config.TableUser, error) {
	// 验证输入数据
	if account == "" || password == "" {
		return nil, errors.New("Account/Password is empty")
	}

	// 查询用户
	info := &config.TableUser{Account: account}
	if err := ctx.Tx.Where(&info).First(&info).Error; err != nil {
		return nil, ErrLoginFailed
	}

	// 验证密码
	if info.Password != yafeng.Hash(password) {
		return nil, ErrLoginFailed
	}

	// 返回用户信息
	return info, nil
}

func UserChange(ctx *yafeng.Context, id int, password string) error {
	// 验证输入数据
	if id == 0 {
		return errors.New("ID is empty")
	}

	// 查询用户
	info := &config.TableUser{ID: id}
	if err := ctx.Tx.Where(&info).First(&info).Error; err != nil {
		return err
	}

	// 更新
	if err := ctx.Tx.Where(&config.TableUser{ID: id}).Updates(&config.TableUser{Password: yafeng.Hash(password)}).Error; err != nil {
		return err
	}

	return nil
}

func UserInfo(ctx *yafeng.Context, id int) (*config.TableUser, error) {
	// 验证输入数据
	if id == 0 {
		return nil, errors.New("ID is empty")
	}

	// 查询用户
	info := &config.TableUser{ID: id}
	if err := ctx.Tx.Where(&info).First(&info).Error; err != nil {
		return nil, err
	}

	// 返回用户信息
	return info, nil
}

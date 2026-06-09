package handler

import (
	"encoding/json"
	"errors"

	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type SysConfigHandler struct {
	base BaseHandler
}

func NewSysConfigHandler(injector do.Injector) *SysConfigHandler {
	return &SysConfigHandler{do.MustInvoke[BaseHandler](injector)}
}

// GetConfig godoc
//
//	@Tags			SysConfig
//	@Summary		获取系统设置(部分不敏感的)
//	@Description	敏感信息不返回,包括各种key密钥
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	vo.SysConfigVO
//	@Router			/sysConfig/get [post]
func (s SysConfigHandler) GetConfig(c echo.Context) error {
	var (
		config db.SysConfig
		result vo.SysConfigVO
	)

	if err := s.base.db.First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return SuccessResp(c, h{})
		}
		return FailRespWithMsg(c, Fail, "读取系统配置异常")
	}
	err := json.Unmarshal([]byte(config.Content), &result)
	if err != nil {
		return FailRespWithMsg(c, Fail, "读取系统配置异常")
	}
	result.Version = s.base.cfg.Version
	result.CommitId = s.base.cfg.CommitId

	suffix := result.S3.ThumbnailSuffix
	result.S3 = vo.S3VO{
		ThumbnailSuffix: suffix,
	}
	return SuccessResp(c, result)
}

// GetFullConfig godoc
//
//	@Tags		SysConfig
//	@Summary	获取系统设置(完整的)
//	@Accept		json
//	@Produce	json
//	@Param		x-api-token	header		string	true	"登录TOKEN"
//	@Success	200			{object}	vo.FullSysConfigVO
//	@Success	200
//	@Router		/api/sysConfig/get [post]
func (s SysConfigHandler) GetFullConfig(c echo.Context) error {
	var (
		config db.SysConfig
		result vo.FullSysConfigVO
	)

	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil || currentUser.Id != 1 {
		return FailRespWithMsg(c, Fail, "需要先登录")
	}
	if err := s.base.db.First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return SuccessResp(c, h{})
		}
		return FailRespWithMsg(c, Fail, "读取系统配置异常")
	}
	err := json.Unmarshal([]byte(config.Content), &result)
	if err != nil {
		return FailRespWithMsg(c, Fail, "读取系统配置异常")
	}
	result.Version = s.base.cfg.Version
	result.CommitId = s.base.cfg.CommitId
	return SuccessResp(c, result)
}

// SaveConfig godoc
//
//	@Tags		SysConfig
//	@Summary	保存系统设置
//	@Accept		json
//	@Produce	json
//	@Param		object		body	vo.FullSysConfigVO	true	"保存系统设置"
//	@Param		x-api-token	header	string				true	"登录TOKEN"
//	@Success	200
//	@Router		/api/sysConfig/save [post]
func (s SysConfigHandler) SaveConfig(c echo.Context) error {
	var (
		config db.SysConfig
		result vo.FullSysConfigVO
	)
	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil || currentUser.Id != 1 {
		return FailRespWithMsg(c, Fail, "需要先登录")
	}

	if err := c.Bind(&result); err != nil {
		s.base.log.Info().Msgf("保存配置错误,%s", err)
		return FailResp(c, ParamError)
	}

	data, err := json.Marshal(result)
	if err != nil {
		return FailRespWithMsg(c, Fail, "读取系统配置异常")
	}

	// 读取已有配置:NotFound 表示首次保存(新建),其它错误视为真正的读取失败
	existing := s.base.db.First(&config).Error
	if existing != nil && !errors.Is(existing, gorm.ErrRecordNotFound) {
		return FailRespWithMsg(c, Fail, "读取系统配置异常")
	}
	config.Content = string(data)

	// 配置保存与管理员用户名更新需要原子化,避免中途失败导致不一致
	if err := s.base.db.Transaction(func(tx *gorm.DB) error {
		if errors.Is(existing, gorm.ErrRecordNotFound) {
			if err := tx.Save(&config).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Updates(&config).Error; err != nil {
				return err
			}
		}
		return tx.Table("User").Where("id=?", 1).Update("username", result.AdminUserName).Error
	}); err != nil {
		return FailRespWithMsg(c, Fail, "保存系统配置异常")
	}
	return SuccessResp(c, h{})
}

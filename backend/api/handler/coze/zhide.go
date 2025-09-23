package coze

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/coze-dev/coze-studio/backend/api/model/zhide"
	"github.com/coze-dev/coze-studio/backend/application/user"
	"github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/hertzutil/domain"
	"github.com/coze-dev/coze-studio/backend/pkg/i18n"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

// @router /passport/web/zhide/login [POST]
func ZhideLogin(ctx context.Context, c *app.RequestContext) {
	req := &zhide.ZhideLoginRequest{}
	if err := c.BindAndValidate(req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// 1️⃣ 简单白名单校验，这里先不校验，后边再完善

	// 2️⃣ 请求平台验签
	cli := &http.Client{Timeout: 5 * time.Second}
	b, _ := json.Marshal(map[string]string{"token": req.Token})
	resp, err := cli.Post(req.Url, "application/json", bytes.NewReader(b))
	if err != nil || resp.StatusCode != http.StatusOK {
		c.String(http.StatusBadRequest, "token 验签失败")
		return
	}

	var pinfo zhide.ZhideUser
	if err = json.NewDecoder(resp.Body).Decode(&pinfo); err != nil {
		c.String(http.StatusBadRequest, "返回格式异常")
		return
	}

	if pinfo.Email == "" {
		c.String(http.StatusBadRequest, "用户缺少 Email")
		return
	}

	locale := string(i18n.GetLocale(ctx))

	loginResp, sessionKey, err := user.UserApplicationSVC.PassportWebZhideLogin(ctx, locale, &pinfo)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.SetCookie(entity.SessionKey,
		sessionKey,
		consts.SessionMaxAgeSecond,
		"/", domain.GetOriginHost(c),
		protocol.CookieSameSiteDefaultMode,
		false, true)

	c.JSON(http.StatusOK, loginResp)
}

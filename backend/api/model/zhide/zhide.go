package zhide

type ZhideLoginRequest struct {
	Token string `form:"token,required" json:"token" query:"token,required"`
	Url   string `form:"url,required" json:"url" query:"url,required"`
}

type ZhideLoginResponse struct {
	Msg string `form:"msg,required" json:"msg" query:"msg,required"`
}

// 平台返回的用户信息
type ZhideUser struct {
	Id        uint64 `json:"id"`
	Nickname  string `json:"nickname"`
	Gender    int8   `json:"gender"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	Signature string `json:"signature"`
}

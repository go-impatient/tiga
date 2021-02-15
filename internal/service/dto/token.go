package dto

// AccessDetails ...
type AccessDetails struct {
	TokenUUID string
	AppID     string
	AppSecret string
}

// TokenDetails ...
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUUID    string
	RefreshUUID  string
	AtExpires    int64 // Token的过期时间
	RtExpires    int64 // 刷新Token的过期时间
}

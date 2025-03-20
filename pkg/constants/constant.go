package constants

const (
	CAPTCHA_ID = "captcha_id:"
	CAPTCHA_TIMEOUT = 60*5
	LIMITER_TIMES = 20
	LIMITER_TIME = 30
	JWT_CONTEXT_KEY = "token"
	FIELD_ERROR_INFO = "field_error_info"
	PASSWORD_REGEX   = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[A-Za-z\d]{8,32}$`
	VALIDATION_ERROR = 444
	TOKEN_EXPIRATION = 30
	TOKEN_CHECK = 3
)
const (
	USER_ONLINE = "online"
	USER_OFFLINE = "offline"
	FRIEND_STATUS_WAITING = "waiting"
	FRIEND_STATUS_FRIENDLY  = "friendly"
	FRIEND_STATUS_BLACKLIST = "blacklist"
)

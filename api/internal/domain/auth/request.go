package auth

type RequestSetPassword struct {
	Body struct {
		Username string `json:"username" minLength:"4" maxLength:"32"`
		Password string `json:"password" minLength:"8" maxLength:"32"`
	}
}

type RequestLoginPlayer struct {
	Body struct {
		Username string `json:"username" minLength:"4" maxLength:"32"`
		Password string `json:"password" minLength:"8" maxLength:"32"`
	}
}

type ResponseLoginPlayer struct {
	Body struct {
		Token string `json:"token" readOnly:"true"`
	}
}

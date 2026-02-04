package auth

type RequestRegisterCreatePlayer struct {
	Body struct {
		Username string  `json:"username" minLength:"4" maxLength:"32"`
		Img      *string `json:"img" format:"uri" required:"false"`
		Email    *string `json:"email" format:"email" required:"false"`
		Password string  `json:"password" minLength:"8" maxLength:"32"`
	}
}

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

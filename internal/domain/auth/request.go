package auth

type RequestRegisterCreatePlayer struct {
	Body struct {
		Username string  `json:"username"`
		Img      *string `json:"img" required:"false"`
		Email    *string `json:"email" format:"email" required:"false"`
		Password string  `json:"password"`
	}
}

type RequestSetPassword struct {
	Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
}

type RequestLoginPlayer struct {
	Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
}

type ResponseLoginPlayer struct {
	Body struct {
		Token string `json:"token"`
	}
}

package forms

type PassWordLoginForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式要用正则限制，自定义validator，注意binding要加mobile验证
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=10"`
}

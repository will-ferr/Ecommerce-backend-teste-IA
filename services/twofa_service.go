package services

import (
	"bytes"
	"encoding/base64"
	"image/png"

	"smart-choice/models"
	"smart-choice/repository"

	"github.com/pquerna/otp/totp"
)

func Generate2FA(user *models.User) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SmartChoice",
		AccountName: user.Email,
	})
	if err != nil {
		return "", "", err
	}

	user.TwoFASecret = key.Secret()
	user.TwoFA = true
	if err := repository.UpdateUser(user); err != nil {
		return "", "", err
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return "", "", err
	}
	png.Encode(&buf, img)
	qrCode := base64.StdEncoding.EncodeToString(buf.Bytes())

	return qrCode, key.Secret(), nil
}

func Validate2FA(user *models.User, code string) bool {
	return totp.Validate(code, user.TwoFASecret)
}

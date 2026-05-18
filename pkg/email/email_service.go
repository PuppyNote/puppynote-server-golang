package email

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/PuppyNote/puppynote-server-golang/config"
	"gopkg.in/gomail.v2"
)

func SendVerificationCode(to string) (string, error) {
	code := generateCode()
	store.save(to, code)

	if err := sendMail(to, code); err != nil {
		return "", err
	}
	return code, nil
}

func VerifyCode(email, code string) bool {
	return store.verify(email, code)
}

func generateCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

func sendMail(to, code string) error {
	cfg := config.AppConfig.Email

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("PuppyNote <%s>", cfg.Username))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "[PuppyNote] 이메일 인증번호를 확인해주세요")
	m.SetBody("text/html", buildEmailBody(code))

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return d.DialAndSend(m)
}

func buildEmailBody(code string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; text-align: center;">
  <div style="max-width: 500px; margin: 0 auto; padding: 40px 20px;">
    <div style="background: linear-gradient(135deg, #eebd2b, #f0d060); padding: 20px; border-radius: 12px 12px 0 0;">
      <h1 style="color: white; margin: 0;">🐾 PuppyNote</h1>
    </div>
    <div style="border: 1px solid #eee; padding: 40px; border-radius: 0 0 12px 12px;">
      <p style="color: #333; font-size: 16px;">이메일 인증번호입니다.</p>
      <div style="font-size: 36px; font-weight: bold; letter-spacing: 8px; color: #eebd2b; padding: 20px;">%s</div>
      <p style="color: #888; font-size: 13px;">인증번호는 5분간 유효합니다.</p>
    </div>
  </div>
</body>
</html>`, code)
}

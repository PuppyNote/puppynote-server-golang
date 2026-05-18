package errors

import "net/http"

type PuppyNoteError struct {
	Code    int
	Message string
}

func (e *PuppyNoteError) Error() string {
	return e.Message
}

func New(message string) *PuppyNoteError {
	return &PuppyNoteError{Code: http.StatusBadRequest, Message: message}
}

func NewWithCode(code int, message string) *PuppyNoteError {
	return &PuppyNoteError{Code: code, Message: message}
}

var (
	ErrUnauthorized    = NewWithCode(http.StatusUnauthorized, "인증이 필요합니다.")
	ErrTokenExpired    = NewWithCode(http.StatusUnauthorized, "토큰이 만료되었습니다.")
	ErrTokenInvalid    = NewWithCode(http.StatusUnauthorized, "유효하지 않은 토큰입니다.")
	ErrUserNotFound    = New("존재하지 않는 유저입니다.")
	ErrPetNotFound     = New("존재하지 않는 반려동물입니다.")
	ErrDuplicateEmail  = New("이미 사용중인 이메일입니다.")
	ErrInvalidPassword = New("비밀번호가 일치하지 않습니다.")
	ErrEmailNotVerified = New("이메일 인증이 필요합니다.")
	ErrInvalidCode     = New("유효하지 않은 인증 코드입니다.")
	ErrExpiredCode     = New("만료된 인증 코드입니다.")
	ErrFoodNotRelated  = New("음식에 관한 질문만 해주세요.")
)

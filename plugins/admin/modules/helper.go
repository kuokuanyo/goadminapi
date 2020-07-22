package modules

import uuid "github.com/satori/go.uuid"

func Uuid() string {
	uid, _ := uuid.NewV4()
	rst := uid.String()
	return rst
}

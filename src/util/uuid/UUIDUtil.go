package uuid

import "github.com/satori/go.uuid"

func GetRandomUUID() string {
	uid,_:= uuid.NewV4()
	return uid.String()
}

package app

import (
	"github.com/Andrey-Kachow/goauth-backdev/pkg/db"
	"github.com/Andrey-Kachow/goauth-backdev/pkg/msg"
)

type ApplicationContext struct {
	NotificationService msg.NotificationService
	TokenDB             db.TokenDB
}

var appContext ApplicationContext

func Init() {
	appContext = ApplicationContext{
		NotificationService: msg.ProvideNotificationService(),
		TokenDB:             db.ProvideApplicationTokenDB(),
	}
}

func Context() ApplicationContext {
	return appContext
}

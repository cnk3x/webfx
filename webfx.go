package webfx

import (
	"github.com/cnk3x/webfx/db"
	"github.com/cnk3x/webfx/utils/fxs"
	"github.com/cnk3x/webfx/web"

	"go.uber.org/fx"
)

func Run(fxSet ...fx.Option) {
	fx.New(
		fx.WithLogger(fxs.Logger),
		fx.Provide(db.GormOpen, web.NewJwt),
		fx.Options(fxSet...),
		fx.Invoke(web.Run),
	).Run()
}

package client

import (
	"github.com/artisanhe/tools/courier"
)

type IRequest interface {
	Do() courier.Result
}

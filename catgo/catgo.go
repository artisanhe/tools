package catgo

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/artisanhe/tools/catgo/cat-go/cat"
	"github.com/artisanhe/tools/conf"
)

func init() {
	cat.ResetLogger(&innerLogger{l: logrus.WithField("Component", "Catgo")})
}

type Catgo struct {
	Domain      string   `conf:"env"`
	Debug       bool     `conf:"env"`
	HttpServers []string `conf:"env"`
}

func (c Catgo) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"HttpServers": []string{"cat:2280"},
	}
}

func (c Catgo) MarshalDefaults(v interface{}) {
	if cat, ok := v.(*Catgo); ok {
		if cat.Domain == "" {
			cat.Domain = "g7pay.local"
		}
		cat.Debug = true
		if len(cat.HttpServers) == 0 {
			cat.HttpServers = []string{"cat:2280"}
		}
	}
}

func (c *Catgo) Init() {
	if c.Debug {
		cat.DebugOn()
	}
	services := []cat.HostConfig{}

	for _, s := range c.HttpServers {
		h, p, e := parseHost(s)
		if e != nil {
			panic(e)
		}
		services = append(services, cat.HostConfig{
			Host: h,
			Port: p,
		})
	}
	if err := cat.InitWithServer(c.Domain, services); err != nil {
		panic(err)
	}

}

type innerLogger struct {
	l *logrus.Entry
}

func (i *innerLogger) SetOutput(io.Writer) {

}

func (i *innerLogger) Printf(format string, args ...interface{}) {
	i.l.Printf(format, args...)
}

func parseHost(raw string) (string, int, error) {
	ss := strings.Split(raw, ":")
	if len(ss) != 2 {
		return "", 0, errors.New("parse host error")
	}
	port, err := strconv.Atoi(ss[1])
	if err != nil {
		return "", 0, fmt.Errorf("parse host error : %v", err)
	}
	if port <= 0 || port > 65535 {
		return "", 0, errors.New("parse host error : port illegal")
	}

	return ss[0], port, nil

}

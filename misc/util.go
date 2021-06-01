package misc

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	PersonType string = "1"
	CompayType string = "2"
)

func Float64ToInt64(i float64) (int64, error) {
	strTemp := fmt.Sprintf("%.f", i*100)
	num, err := strconv.ParseInt(strTemp, 10, 64)
	if err != nil {
		logrus.Warningf("bal to int64 fail[err:%v]", err)
		return 0, err
	}
	return num, nil
}

func Int64ToFloat64(i int64) (float64, error) {
	strTemp := fmt.Sprintf("%.2f", float64(i)/float64(100))
	num, err := strconv.ParseFloat(strTemp, 64)
	if err != nil {
		logrus.Errorf("bal to float64 fail[err:%v]", err)
		return 0, err
	}
	return num, nil
}

// Convert int64 to float64 and round to two decimals
// Return string type.
func Int64ToFloat64String(i int64) string {
	return fmt.Sprintf("%.2f", float64(i)/float64(100))
}

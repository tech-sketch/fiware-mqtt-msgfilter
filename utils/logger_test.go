/*
Package utils : utilities.

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package utils

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/tech-sketch/fiware-mqtt-msgfilter/mock"
)

func TestLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWriter := mock.NewMockiWriter(ctrl)

	logger := &Logger{
		Name:   "test",
		writer: mockWriter,
		now: func() time.Time {
			return time.Date(2018, 12, 31, 18, 0, 17, 987, time.Local)
		},
	}

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |debug| [test] foo")
	logger.Debugf("foo")

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |debug| [test] foo %s", "HOGE")
	logger.Debugf("foo %s", "HOGE")

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |debug| [test] foo %s:%d", "HOGE", 1)
	logger.Debugf("foo %s:%d", "HOGE", 1)

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |info | [test] foo")
	logger.Infof("foo")

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |info | [test] foo %s", "HOGE")
	logger.Infof("foo %s", "HOGE")

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |info | [test] foo %s:%d", "HOGE", 1)
	logger.Infof("foo %s:%d", "HOGE", 1)

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |warn | [test] foo")
	logger.Warnf("foo")

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |warn | [test] foo %s", "HOGE")
	logger.Warnf("foo %s", "HOGE")

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |warn | [test] foo %s:%d", "HOGE", 1)
	logger.Warnf("foo %s:%d", "HOGE", 1)

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |error| [test] foo")
	logger.Errorf("foo")

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |error| [test] foo %s", "HOGE")
	logger.Errorf("foo %s", "HOGE")

	mockWriter.EXPECT().Printf("[APP] 2018/12/31 - 18:00:17 |error| [test] foo %s:%d", "HOGE", 1)
	logger.Errorf("foo %s:%d", "HOGE", 1)
}

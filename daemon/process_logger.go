package daemon

import (
	"github.com/ihaiker/gokit/commons"
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/gokit/logs/appenders"
	"io"
)

const LOGGER_LINES_NUM = 100

type TailLogger func(id, line string)

type ProcessLogger struct {
	file    io.Writer
	lines   []string
	lineIdx *atomic.AtomicInt
	tails   map[string]TailLogger
}

func NewLogger(loggerFile string) (logger *ProcessLogger, err error) {
	logger = new(ProcessLogger)
	writers := make([]io.Writer, 0)
	if loggerFile != "" {
		if _, match := appenders.MatchDailyRollingFile(loggerFile); match {
			if logger.file, err = appenders.NewDailyRollingFileOut(loggerFile); err != nil {
				return nil, err
			} else {
				writers = append(writers, logger.file)
			}
		} else {
			if logger.file, err = files.New(loggerFile).GetWriter(true); err != nil {
				return
			} else {
				writers = append(writers, logger.file)
			}
		}
	}
	logger.lineIdx = atomic.NewAtomicInt(0)
	logger.lines = make([]string, LOGGER_LINES_NUM)
	logger.tails = map[string]TailLogger{}
	return
}

func (self *ProcessLogger) Reg(id string, tail TailLogger, lines int) {
	logger.Debug("reg logger ", id)
	self.tails[id] = tail
	if lines > LOGGER_LINES_NUM {
		lines = LOGGER_LINES_NUM
	}
	idx := self.lineIdx.Get()
	if idx < lines {
		idx = 0
	} else if idx < LOGGER_LINES_NUM {
		idx -= lines
	}
	for i := idx; i < idx+lines; i++ {
		line := self.lines[(idx+i)%LOGGER_LINES_NUM]
		if line != "" {
			tail(id, line)
		}
	}
}

func (self *ProcessLogger) UnReg(id string) {
	delete(self.tails, id)
}

func (self *ProcessLogger) Write(p []byte) (int, error) {
	line := string(p)
	self.lines[self.lineIdx.GetAndIncrement(1)%LOGGER_LINES_NUM] = line
	for id, tail := range self.tails {
		commons.SafeExec(func() {
			tail(id, line)
		})
	}
	if self.file != nil {
		return self.file.Write(p)
	}
	return len(p), nil
}

func (self *ProcessLogger) Close() error {
	if self.file != nil {
		return (self.file.(io.Closer)).Close()
	}
	return nil
}

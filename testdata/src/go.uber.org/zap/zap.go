package zap

import "errors"

type Field struct{}

type Logger struct{}

type SugaredLogger struct{}

type Option interface{}

type Core interface{}

func NewProduction(...Option) (*Logger, error) { return &Logger{}, nil }
func NewDevelopment(...Option) (*Logger, error) { return &Logger{}, nil }
func New(_ Core, _ ...Option) *Logger           { return &Logger{} }

func (l *Logger) Info(msg string, fields ...Field)   {}
func (l *Logger) Debug(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)   {}
func (l *Logger) Error(msg string, fields ...Field)  {}
func (l *Logger) DPanic(msg string, fields ...Field) {}
func (l *Logger) Panic(msg string, fields ...Field)  {}
func (l *Logger) Fatal(msg string, fields ...Field)  {}
func (l *Logger) Sync() error                        { return errors.New("stub") }
func (l *Logger) Sugar() *SugaredLogger              { return &SugaredLogger{} }

func (s *SugaredLogger) Info(args ...interface{})                       {}
func (s *SugaredLogger) Infof(template string, args ...interface{})     {}
func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{}) {}
func (s *SugaredLogger) Debug(args ...interface{})                      {}
func (s *SugaredLogger) Debugf(template string, args ...interface{})    {}
func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{}) {}
func (s *SugaredLogger) Warn(args ...interface{})                       {}
func (s *SugaredLogger) Warnf(template string, args ...interface{})     {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{}) {}
func (s *SugaredLogger) Error(args ...interface{})                      {}
func (s *SugaredLogger) Errorf(template string, args ...interface{})    {}
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {}
func (s *SugaredLogger) Sync() error                                    { return errors.New("stub") }

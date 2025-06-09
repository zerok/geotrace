package logging

import "log/slog"

func Addr(name string) slog.Attr {
	return slog.String("addr", name)
}

func Err(err error) slog.Attr {
	return slog.Any("err", err)
}

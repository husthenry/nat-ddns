package constants

import "errors"

var ErrShortWrite = errors.New("short write")
var ErrMarshal = errors.New("data marshal")
var EOF = errors.New("EOF")

package websocket

import "github.com/gorilla/websocket"

// IsCloseError checks if the error is a close error with the given codes.
// Will return true if is a CloseError but no codes are provided.
func IsCloseError(err error, codes ...int) bool {
	if err == nil {
		return false
	}

	e, ok := err.(*websocket.CloseError)
	if !ok {
		return false
	}

	if len(codes) == 0 {
		return true
	}

	for _, code := range codes {
		if e.Code == code {
			return true
		}
	}

	return false
}

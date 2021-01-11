//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {

	var requestControl chan bool = make(chan bool)

	if !ValidateRequest(u) {
		return false
	}

	go func(ctrl chan<- bool) {
		process()
		ctrl <- true
	}(requestControl)

	t := time.Tick(time.Second * 1)

	for {

		<-t

		select {

		case done := <-requestControl:

			return done

		default:

			atomic.AddInt64(&u.TimeUsed, 1)

			if !ValidateRequest(u) {
				return false
			}
		}

	}
}

// ValidateRequest checks if a request is still valid
func ValidateRequest(u *User) bool {
	return u.IsPremium || atomic.LoadInt64(&u.TimeUsed) < 10
}

func main() {
	RunMockServer()
}

// Package notify has functions and types used for sending notifications on different communication channel
package notify

// Notifier can be any type that can SendMessage
type Notifier interface {
	SendMessage(string) error
}

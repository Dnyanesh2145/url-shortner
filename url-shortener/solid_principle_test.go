package main

import (
	"testing"
)

// Interface pattern
type Notification interface {
	Send(to, message string)
}

type EmailNotification struct{}
func (EmailNotification) Send(to, message string) {}

type SMSNotification struct{}
func (SMSNotification) Send(to, message string) {}

// IF/SWITCH pattern
func notifyIf(notificationType, to, msg string) {
	if notificationType == "email" {
		// send email
	} else if notificationType == "sms" {
		// send sms
	}
}

// Interface-based
func notifyInterface(n Notification, to, msg string) {
	n.Send(to, msg)
}

// BENCHMARKS
func BenchmarkNotifyIf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		notifyIf("email", "x", "y")
		notifyIf("sms", "x", "y")
	}
}

func BenchmarkNotifyInterface(b *testing.B) {
	email := EmailNotification{}
	sms := SMSNotification{}
	for i := 0; i < b.N; i++ {
		notifyInterface(email, "x", "y")
		notifyInterface(sms, "x", "y")
	}
}
 
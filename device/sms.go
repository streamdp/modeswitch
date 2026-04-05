package device

import (
	"bytes"
	"fmt"

	"github.com/streamdp/modeswitch/config"
)

const (
	SmsUnread    = "unread"
	SmsNoContent = "no-content"
)

// SmsList - Show a list of received SMS messages.
// Possible arguments:
//
//	    unread - Show a list of unread SMS messages only.
//		id - Show message with the given identifier.
//		no-content - Disable message text output.
func SmsList(sender Sender, c *config.UserConfig, args ...string) ([]byte, error) {
	var cmd = fmt.Sprintf("sms %s list", c.InterfaceId)

	for _, arg := range args {
		cmd += " " + arg
	}

	out, err := sender.Send(cmd)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// SmsDelete - Delete SMS message.
func SmsDelete(sender Sender, c *config.UserConfig, id string) (bool, error) {
	out, err := sender.Send(
		fmt.Sprintf("sms %s delete %s", c.InterfaceId, id),
	)
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte("message deleted")), nil
}

// SmsRead - Mark SMS as read.
func SmsRead(sender Sender, c *config.UserConfig, id string) (bool, error) {
	out, err := sender.Send(
		fmt.Sprintf("sms %s read %s", c.InterfaceId, id),
	)
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte("message deleted")), nil
}

// SmsSend - Send SMS to specified number. The maximum value of saved incoming SMS
// messages in the router's memory is 128. If the memory is full, the oldest SMS
// from the memory will be automatically deleted when a new SMS is received.
func SmsSend(sender Sender, c *config.UserConfig, to, msg string) (bool, error) {
	out, err := sender.Send(
		fmt.Sprintf("sms %s send %s %s", c.InterfaceId, to, msg),
	)
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte("message sent")), nil
}

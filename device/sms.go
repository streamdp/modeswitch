package device

import (
	"bytes"
	"fmt"
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
func (d *Device) SmsList(args ...string) ([]byte, error) {
	var cmd = fmt.Sprintf("sms %s list", d.cfg.InterfaceId)

	for _, arg := range args {
		cmd += " " + arg
	}

	out, err := d.sender.Send(cmd)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// SmsDelete - Delete SMS message.
func (d *Device) SmsDelete(id string) (bool, error) {
	out, err := d.sender.Send(
		fmt.Sprintf("sms %s delete %s", d.cfg.InterfaceId, id),
	)
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte("message deleted")), nil
}

// SmsRead - Mark SMS as read.
func (d *Device) SmsRead(id string) (bool, error) {
	out, err := d.sender.Send(
		fmt.Sprintf("sms %s read %s", d.cfg.InterfaceId, id),
	)
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte("message deleted")), nil
}

// SmsSend - Send SMS to specified number. The maximum value of saved incoming SMS
// messages in the router's memory is 128. If the memory is full, the oldest SMS
// from the memory will be automatically deleted when a new SMS is received.
func (d *Device) SmsSend(to, msg string) (bool, error) {
	out, err := d.sender.Send(
		fmt.Sprintf("sms %s send %s %s", d.cfg.InterfaceId, to, msg),
	)
	if err != nil {
		return false, err
	}

	return bytes.Contains(out, []byte("message sent")), nil
}

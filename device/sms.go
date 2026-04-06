package device

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	SmsUnread    = "unread"
	SmsNoContent = "no-content"
)

type SMSData struct {
	NVFreeSlots   int        `json:"nv_free_slots"`
	NVTotalSlots  int        `json:"nv_total_slots"`
	SimFreeSlots  int        `json:"sim_free_slots"`
	SimTotalSlots int        `json:"sim_total_slots"`
	MessagesCount int        `json:"messages_count"`
	Messages      []*Message `json:"messages"`
}

type Message struct {
	ID         string `json:"id"`
	Read       bool   `json:"read"`
	From       string `json:"from"`
	Timestamp  string `json:"timestamp"`
	Parts      int    `json:"parts"`
	TotalParts int    `json:"total_parts"`
	Text       string `json:"text"`
}

func (m *Message) String() string {
	if m.From != "" && m.Text != "" {
		return fmt.Sprintf("From: %s\nText: %s", m.From, m.Text)
	}

	return fmt.Sprintf("ID: %s\nRead: %v", m.ID, m.Read)
}

// SmsList - Show a list of received SMS messages.
// Possible arguments:
//
//	    unread - Show a list of unread SMS messages only.
//		id - Show message with the given identifier.
//		no-content - Disable message text output.
func (d *Device) SmsList(args ...string) (*Response, error) {
	var cmd = fmt.Sprintf("sms %s list", d.cfg.InterfaceId)

	for _, arg := range args {
		cmd += " " + arg
	}

	out, err := d.sender.Send(cmd)
	if err != nil {
		return nil, err
	}

	if err = errorLookup(out); err != nil {
		return nil, err
	}

	return &Response{SMS: parseSMSBytes(out)}, nil
}

// SmsDelete - Delete SMS message.
func (d *Device) SmsDelete(id string) (*Response, error) {
	out, err := d.sender.Send(
		fmt.Sprintf("sms %s delete %s", d.cfg.InterfaceId, id),
	)
	if err != nil {
		return nil, err
	}

	return &Response{Text: string(out)}, nil
}

// SmsRead - Mark SMS as read.
func (d *Device) SmsRead(id string) (*Response, error) {
	out, err := d.sender.Send(
		fmt.Sprintf("sms %s read %s", d.cfg.InterfaceId, id),
	)
	if err != nil {
		return nil, err
	}

	return &Response{Text: string(out)}, nil
}

// SmsSend - Send SMS to specified number. The maximum value of saved incoming SMS
// messages in the router's memory is 128. If the memory is full, the oldest SMS
// from the memory will be automatically deleted when a new SMS is received.
func (d *Device) SmsSend(to, msg string) (*Response, error) {
	out, err := d.sender.Send(
		fmt.Sprintf("sms %s send %s %s", d.cfg.InterfaceId, to, msg),
	)
	if err != nil {
		return nil, err
	}

	return &Response{Text: string(out)}, nil
}

func errorLookup(raw []byte) error {
	if i := bytes.Index(raw, []byte("Command::Base error")); i != -1 {
		return errors.New(string(raw[i:]))
	}
	return nil
}

func parseSMSBytes(raw []byte) *SMSData {
	var data = new(SMSData)
	var currentMsg *Message

	scanner := bufio.NewScanner(bytes.NewReader(raw))

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		if strings.Contains(line, "slots") ||
			strings.Contains(line, "messages-count") {
			parts := strings.SplitN(trimmed, ":", 2)
			key := strings.TrimSpace(parts[0])
			val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))

			switch key {
			case "nv-free-slots":
				data.NVFreeSlots = val
			case "nv-total-slots":
				data.NVTotalSlots = val
			case "sim-free-slots":
				data.SimFreeSlots = val
			case "sim-total-slots":
				data.SimTotalSlots = val
			case "messages-count":
				data.MessagesCount = val
			}

			continue
		}

		if strings.Contains(line, "messages, id =") {
			if currentMsg != nil {
				data.Messages = append(data.Messages, currentMsg)
			}

			idPart := strings.Split(line, "=")
			if len(idPart) > 1 {
				id := strings.TrimSpace(idPart[1])
				currentMsg = &Message{ID: strings.TrimSuffix(id, ":")}
			}

			continue
		}

		if currentMsg != nil {
			if currentMsg.Text != "" {
				currentMsg.Text += "\n" + trimmed

				continue
			}

			parts := strings.SplitN(trimmed, ": ", 2)
			if len(parts) != 2 {
				continue
			}

			key, val := parts[0], parts[1]

			switch key {
			case "read":
				currentMsg.Read = val == "yes"
			case "from":
				currentMsg.From = val
			case "timestamp":
				currentMsg.Timestamp = val
			case "parts":
				currentMsg.Parts, _ = strconv.Atoi(val)
			case "total-parts":
				currentMsg.TotalParts, _ = strconv.Atoi(val)
			case "text":
				currentMsg.Text = val
			}
		}
	}

	if currentMsg != nil {
		data.Messages = append(data.Messages, currentMsg)
	}

	return data
}

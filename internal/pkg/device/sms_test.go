package device

import (
	"reflect"
	"testing"
)

func Test_parseSMSBytes(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name string
		args args
		want *SMSData
	}{
		{
			name: "parse sms list no-content output bytes",
			args: args{
				raw: []byte(`

nv-free-slots: 12
nv-total-slots: 23
sim-free-slots: 10
sim-total-slots: 10
messages-count: 5

	messages, id = nv-3:
	read: yes
	
	messages, id = nv-7:
	read: yes
	
	messages, id = nv-2:
	read: yes
	
	messages, id = nv-0:
	read: yes
	
	messages, id = nv-1:
	read: yes

`),
			},
			want: &SMSData{
				NVFreeSlots:   12,
				NVTotalSlots:  23,
				SimFreeSlots:  10,
				SimTotalSlots: 10,
				MessagesCount: 5,
				Messages: []*Message{
					{
						ID:   "nv-3",
						Read: true,
					},
					{
						ID:   "nv-7",
						Read: true,
					},
					{
						ID:   "nv-2",
						Read: true,
					},
					{
						ID:   "nv-0",
						Read: true,
					},
					{
						ID:   "nv-1",
						Read: true,
					},
				},
			},
		},
		{
			name: "parse sms list output bytes",
			args: args{
				raw: []byte(`

nv-free-slots: 23
nv-total-slots: 23
sim-free-slots: 0
sim-total-slots: 15

	messages, id = sim-0:
	read: yes
	from: +79658283425
	timestamp: Thu Aug 20 14:39:57 2020
	parts: 1
	total-parts: 1
	text: Accepted

	messages, id = sim-1:
	read: yes
	from: MegaFon
	timestamp: Wed Sep 9 13:57:21 2020
	parts: 2
	total-parts: 2
	text: 636-269 — your personal login code.
	Do not share this code with anyone.

	messages, id = sim-3:
	read: yes
	from: +79658283425
	timestamp: Wed Sep 9 16:32:26 2020
	parts: 1
	total-parts: 1
	text: Our time to your time to yes to

	messages, id = sim-4:
	read: yes
	from: +79658283425
	timestamp: Mon Sep 14 17:14:11 2020
	parts: 1
	total-parts: 1
	text: Ok

	messages, id = sim-5:
	read: yes
	from: MegaFon
	timestamp: Wed Sep 16 10:24:46 2020
	parts: 7
	total-parts: 7
	text: Listen to audiobooks on management, leadership,
	personal efficiency and self-development 2 weeks free!
	Just subscribe to the MegaFon AudioBooks and listen to them without advertising on any convenient device.
	The cost after the trial period - 1 euro / day.
	Payment from the phone account without card binding. Cancel subscriptions at any time: pay only for days
	of usage. Learn more:
	http://i.megafon.com/Q2XadzRp9xusLwS1

	messages, id = sim-12:
	read: no
	from: +79252384670
	timestamp: Fri Sep 18 19:02:27 2020
	parts: 3
	total-parts: 4
	text: This subscriber left you 18.09.2020 at 18:35
	voice message. You can listen to it for free by
	number 0525. / Listen to podcasts and book parodies in
	convenient application without advertising for 5 e/d.

`),
			},
			want: &SMSData{
				NVFreeSlots:   23,
				NVTotalSlots:  23,
				SimFreeSlots:  0,
				SimTotalSlots: 15,
				Messages: []*Message{
					{
						ID:         "sim-0",
						Read:       true,
						From:       "+79658283425",
						Timestamp:  "Thu Aug 20 14:39:57 2020",
						Parts:      1,
						TotalParts: 1,
						Text:       "Accepted",
					},
					{
						ID:         "sim-1",
						Read:       true,
						From:       "MegaFon",
						Timestamp:  "Wed Sep 9 13:57:21 2020",
						Parts:      2,
						TotalParts: 2,
						Text:       "636-269 — your personal login code.\nDo not share this code with anyone.",
					},
					{
						ID:         "sim-3",
						Read:       true,
						From:       "+79658283425",
						Timestamp:  "Wed Sep 9 16:32:26 2020",
						Parts:      1,
						TotalParts: 1,
						Text:       "Our time to your time to yes to",
					},
					{
						ID:         "sim-4",
						Read:       true,
						From:       "+79658283425",
						Timestamp:  "Mon Sep 14 17:14:11 2020",
						Parts:      1,
						TotalParts: 1,
						Text:       "Ok",
					},
					{
						ID:         "sim-5",
						Read:       true,
						From:       "MegaFon",
						Timestamp:  "Wed Sep 16 10:24:46 2020",
						Parts:      7,
						TotalParts: 7,
						Text: "Listen to audiobooks on management, leadership,\n" +
							"personal efficiency and self-development 2 weeks free!\n" +
							"Just subscribe to the MegaFon AudioBooks and listen to them without advertising on any convenient device.\n" +
							"The cost after the trial period - 1 euro / day.\n" +
							"Payment from the phone account without card binding. Cancel subscriptions at any time: pay only for days\n" +
							"of usage. Learn more:\n" +
							"http://i.megafon.com/Q2XadzRp9xusLwS1",
					},
					{
						ID:         "sim-12",
						Read:       false,
						From:       "+79252384670",
						Timestamp:  "Fri Sep 18 19:02:27 2020",
						Parts:      3,
						TotalParts: 4,
						Text: "This subscriber left you 18.09.2020 at 18:35\n" +
							"voice message. You can listen to it for free by\n" +
							"number 0525. / Listen to podcasts and book parodies in\n" +
							"convenient application without advertising for 5 e/d.",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSMSBytes(tt.args.raw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSMSBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

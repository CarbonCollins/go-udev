package udev

import "testing"

type Testcases []Testcase

type Testcase struct {
	Object interface{}
	Valid  bool
}

func TestRules(testing *testing.T) {
	t := testingWrapper{testing}

	uevent := UEvent{
		Action: ADD,
		KObj:   "/devices/pci0000:00/0000:00:14.0/usb2/2-1/2-1:1.2/0003:04F2:0976.0008/hidraw/hidraw4",
		Env: map[string]string{
			"ACTION":    "add",
			"DEVPATH":   "/devices/pci0000:00/0000:00:14.0/usb2/2-1/2-1:1.2/0003:04F2:0976.0008/hidraw/hidraw4",
			"SUBSYSTEM": "hidraw",
			"MAJOR":     "247",
			"MINOR":     "4",
			"DEVNAME":   "hidraw4",
			"SEQNUM":    "2569",
		},
	}

	add := ADD.String()
	wrongAction := "can't match"

	rules := []RuleDefinition{
		RuleDefinition{
			Action: nil,
			Env: map[string]string{
				"DEVNAME": "hidraw\\d+",
			},
		},

		RuleDefinition{
			Action: &add,
			Env:    make(map[string]string, 0),
		},

		RuleDefinition{
			Action: nil,
			Env: map[string]string{
				"SUBSYSTEM": "can't match",
				"MAJOR":     "247",
			},
		},

		RuleDefinition{
			Action: &add,
			Env: map[string]string{
				"SUBSYSTEM": "hidraw",
				"MAJOR":     "\\d+",
			},
		},

		RuleDefinition{
			Action: &wrongAction,
			Env: map[string]string{
				"SUBSYSTEM": "hidraw",
				"MAJOR":     "\\d+",
			},
		},
	}

	testcases := []Testcase{
		Testcase{
			Object: &rules[0],
			Valid:  true,
		},
		Testcase{
			Object: &rules[1],
			Valid:  true,
		},
		Testcase{
			Object: &rules[2],
			Valid:  false,
		},
		Testcase{
			Object: &rules[3],
			Valid:  true,
		},
		Testcase{
			Object: &rules[4],
			Valid:  false,
		},
		Testcase{
			Object: &Or{[]RuleDefinition{rules[0], rules[4]}},
			Valid:  true,
		},
		Testcase{
			Object: &Or{[]RuleDefinition{rules[4], rules[0]}},
			Valid:  true,
		},
		Testcase{
			Object: &Or{[]RuleDefinition{rules[2], rules[4]}},
			Valid:  false,
		},
		Testcase{
			Object: &Or{[]RuleDefinition{rules[3], rules[1]}},
			Valid:  true,
		},
	}

	for k, tcase := range testcases {
		matcher := tcase.Object.(Matcher)

		err := matcher.Compile()
		t.FatalfIf(err != nil, "Testcase n°%d should compile without error, err: %v", k+1, err)

		ok := matcher.Evaluate(uevent)
		t.FatalfIf((ok != tcase.Valid) && tcase.Valid, "Testcase n°%d should evaluate event", k+1)
		t.FatalfIf((ok != tcase.Valid) && !tcase.Valid, "Testcase n°%d shouldn't evaluate event", k+1)
	}
}
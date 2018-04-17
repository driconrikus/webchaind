// Copyright 2015 The go-ethereum Authors
// This file is part of Webchain.
//
// Webchain is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Webchain is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Webchain. If not, see <http://www.gnu.org/licenses/>.

package whisper

import (
	"bytes"

	"testing"
)

var filterTopicsCreationTests = []struct {
	topics [][]string
	filter [][][4]byte
}{
	{ // Simple topic filter
		topics: [][]string{
			{"abc", "def", "ghi"},
			{"def"},
			{"ghi", "abc"},
		},
		filter: [][][4]byte{
			{{0x4e, 0x03, 0x65, 0x7a}, {0x34, 0x60, 0x7c, 0x9b}, {0x21, 0x41, 0x7d, 0xf9}},
			{{0x34, 0x60, 0x7c, 0x9b}},
			{{0x21, 0x41, 0x7d, 0xf9}, {0x4e, 0x03, 0x65, 0x7a}},
		},
	},
	{ // Wild-carded topic filter
		topics: [][]string{
			{"abc", "def", "ghi"},
			{},
			{""},
			{"def"},
		},
		filter: [][][4]byte{
			{{0x4e, 0x03, 0x65, 0x7a}, {0x34, 0x60, 0x7c, 0x9b}, {0x21, 0x41, 0x7d, 0xf9}},
			{},
			{},
			{{0x34, 0x60, 0x7c, 0x9b}},
		},
	},
}

var filterTopicsCreationFlatTests = []struct {
	topics []string
	filter [][][4]byte
}{
	{ // Simple topic list
		topics: []string{"abc", "def", "ghi"},
		filter: [][][4]byte{
			{{0x4e, 0x03, 0x65, 0x7a}},
			{{0x34, 0x60, 0x7c, 0x9b}},
			{{0x21, 0x41, 0x7d, 0xf9}},
		},
	},
	{ // Wild-carded topic list
		topics: []string{"abc", "", "ghi"},
		filter: [][][4]byte{
			{{0x4e, 0x03, 0x65, 0x7a}},
			{},
			{{0x21, 0x41, 0x7d, 0xf9}},
		},
	},
}

func TestFilterTopicsCreation(t *testing.T) {
	// Check full filter creation
	for i, tt := range filterTopicsCreationTests {
		// Check the textual creation
		filter := newFilterTopicsFromStrings(tt.topics...)
		if len(filter) != len(tt.topics) {
			t.Errorf("test %d: condition count mismatch: have %v, want %v", i, len(filter), len(tt.topics))
			continue
		}
		for j, condition := range filter {
			if len(condition) != len(tt.filter[j]) {
				t.Errorf("test %d, condition %d: size mismatch: have %v, want %v", i, j, len(condition), len(tt.filter[j]))
				continue
			}
			for k := 0; k < len(condition); k++ {
				if bytes.Compare(condition[k][:], tt.filter[j][k][:]) != 0 {
					t.Errorf("test %d, condition %d, segment %d: filter mismatch: have 0x%x, want 0x%x", i, j, k, condition[k], tt.filter[j][k])
				}
			}
		}
		// Check the binary creation
		binary := make([][][]byte, len(tt.topics))
		for j, condition := range tt.topics {
			binary[j] = make([][]byte, len(condition))
			for k, segment := range condition {
				binary[j][k] = []byte(segment)
			}
		}
		filter = NewFilterTopics(binary...)
		if len(filter) != len(tt.topics) {
			t.Errorf("test %d: condition count mismatch: have %v, want %v", i, len(filter), len(tt.topics))
			continue
		}
		for j, condition := range filter {
			if len(condition) != len(tt.filter[j]) {
				t.Errorf("test %d, condition %d: size mismatch: have %v, want %v", i, j, len(condition), len(tt.filter[j]))
				continue
			}
			for k := 0; k < len(condition); k++ {
				if bytes.Compare(condition[k][:], tt.filter[j][k][:]) != 0 {
					t.Errorf("test %d, condition %d, segment %d: filter mismatch: have 0x%x, want 0x%x", i, j, k, condition[k], tt.filter[j][k])
				}
			}
		}
	}
	// Check flat filter creation
	for i, tt := range filterTopicsCreationFlatTests {
		// Check the textual creation
		filter := newFilterTopicsFromStringsFlat(tt.topics...)
		if len(filter) != len(tt.topics) {
			t.Errorf("test %d: condition count mismatch: have %v, want %v", i, len(filter), len(tt.topics))
			continue
		}
		for j, condition := range filter {
			if len(condition) != len(tt.filter[j]) {
				t.Errorf("test %d, condition %d: size mismatch: have %v, want %v", i, j, len(condition), len(tt.filter[j]))
				continue
			}
			for k := 0; k < len(condition); k++ {
				if bytes.Compare(condition[k][:], tt.filter[j][k][:]) != 0 {
					t.Errorf("test %d, condition %d, segment %d: filter mismatch: have 0x%x, want 0x%x", i, j, k, condition[k], tt.filter[j][k])
				}
			}
		}
		// Check the binary creation
		binary := make([][]byte, len(tt.topics))
		for j, topic := range tt.topics {
			binary[j] = []byte(topic)
		}
		filter = newFilterTopicsFlat(binary...)
		if len(filter) != len(tt.topics) {
			t.Errorf("test %d: condition count mismatch: have %v, want %v", i, len(filter), len(tt.topics))
			continue
		}
		for j, condition := range filter {
			if len(condition) != len(tt.filter[j]) {
				t.Errorf("test %d, condition %d: size mismatch: have %v, want %v", i, j, len(condition), len(tt.filter[j]))
				continue
			}
			for k := 0; k < len(condition); k++ {
				if bytes.Compare(condition[k][:], tt.filter[j][k][:]) != 0 {
					t.Errorf("test %d, condition %d, segment %d: filter mismatch: have 0x%x, want 0x%x", i, j, k, condition[k], tt.filter[j][k])
				}
			}
		}
	}
}

var filterCompareTests = []struct {
	matcher filterer
	message filterer
	match   bool
}{
	{ // Wild-card filter matching anything
		matcher: filterer{to: "", from: "", matcher: newTopicMatcher()},
		message: filterer{to: "to", from: "from", matcher: newTopicMatcher(newFilterTopicsFromStringsFlat("topic")...)},
		match:   true,
	},
	{ // Filter matching the to field
		matcher: filterer{to: "to", from: "", matcher: newTopicMatcher()},
		message: filterer{to: "to", from: "from", matcher: newTopicMatcher(newFilterTopicsFromStringsFlat("topic")...)},
		match:   true,
	},
	{ // Filter rejecting the to field
		matcher: filterer{to: "to", from: "", matcher: newTopicMatcher()},
		message: filterer{to: "", from: "from", matcher: newTopicMatcher(newFilterTopicsFromStringsFlat("topic")...)},
		match:   false,
	},
	{ // Filter matching the from field
		matcher: filterer{to: "", from: "from", matcher: newTopicMatcher()},
		message: filterer{to: "to", from: "from", matcher: newTopicMatcher(newFilterTopicsFromStringsFlat("topic")...)},
		match:   true,
	},
	{ // Filter rejecting the from field
		matcher: filterer{to: "", from: "from", matcher: newTopicMatcher()},
		message: filterer{to: "to", from: "", matcher: newTopicMatcher(newFilterTopicsFromStringsFlat("topic")...)},
		match:   false,
	},
	{ // Filter matching the topic field
		matcher: filterer{to: "", from: "from", matcher: newTopicMatcher(newFilterTopicsFromStringsFlat("topic")...)},
		message: filterer{to: "to", from: "from", matcher: newTopicMatcher(newFilterTopicsFromStringsFlat("topic")...)},
		match:   true,
	},
	{ // Filter rejecting the topic field
		matcher: filterer{to: "", from: "", matcher: newTopicMatcher(newFilterTopicsFromStringsFlat("topic")...)},
		message: filterer{to: "to", from: "from", matcher: newTopicMatcher()},
		match:   false,
	},
}

func TestFilterCompare(t *testing.T) {
	for i, tt := range filterCompareTests {
		if match := tt.matcher.Compare(tt.message); match != tt.match {
			t.Errorf("test %d: match mismatch: have %v, want %v", i, match, tt.match)
		}
	}
}

// NewFilterTopicsFlat creates a 2D topic array used by whisper.Filter from flat
// binary data elements.
func newFilterTopicsFlat(data ...[]byte) [][]Topic {
	filter := make([][]Topic, len(data))
	for i, element := range data {
		// Only add non-wildcard topics
		filter[i] = make([]Topic, 0, 1)
		if len(element) > 0 {
			filter[i] = append(filter[i], NewTopic(element))
		}
	}
	return filter
}

// NewFilterTopicsFromStrings creates a 2D topic array used by whisper.Filter
// from textual data elements.
func newFilterTopicsFromStrings(data ...[]string) [][]Topic {
	filter := make([][]Topic, len(data))
	for i, condition := range data {
		// Handle the special case of condition == [""]
		if len(condition) == 1 && condition[0] == "" {
			filter[i] = []Topic{}
			continue
		}
		// Otherwise flatten normally
		filter[i] = newTopicsFromStrings(condition...)
	}
	return filter
}

// NewFilterTopicsFromStringsFlat creates a 2D topic array used by whisper.Filter from flat
// binary data elements.
func newFilterTopicsFromStringsFlat(data ...string) [][]Topic {
	filter := make([][]Topic, len(data))
	for i, element := range data {
		// Only add non-wildcard topics
		filter[i] = make([]Topic, 0, 1)
		if element != "" {
			filter[i] = append(filter[i], newTopicFromString(element))
		}
	}
	return filter
}

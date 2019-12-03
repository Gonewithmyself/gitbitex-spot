// Copyright 2019 GitBitEx.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conf

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

var hosts = []string{
	"127.0.0.1:9092",
	"127.0.0.1:9093",
}

func TestGetConfig(t *testing.T) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      hosts,
		Topic:        "golang",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 5 * time.Millisecond,
	})

	t.Log(w.WriteMessages(context.Background(), kafka.Message{
		Value: []byte("123"),
	}))
	t.Error()
	_ = w
}

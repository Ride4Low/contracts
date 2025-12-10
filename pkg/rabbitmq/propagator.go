package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// AMQPHeadersCarrier adapts amqp.Table to satisfy propagation.TextMapCarrier
type AMQPHeadersCarrier amqp.Table

// Get returns the value associated with the passed key.
func (c AMQPHeadersCarrier) Get(key string) string {
	val, ok := c[key]
	if !ok {
		return ""
	}
	strVal, ok := val.(string)
	if !ok {
		return ""
	}
	return strVal
}

// Set stores the key-value pair.
func (c AMQPHeadersCarrier) Set(key string, value string) {
	c[key] = value
}

// Keys returns the keys for which this carrier has a value.
func (c AMQPHeadersCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

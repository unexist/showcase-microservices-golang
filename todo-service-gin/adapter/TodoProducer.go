//
// @package Showcase-Microservices-Golang
//
// @file Todo producer
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package adapter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"

	"braces.dev/errtrace"

	"github.com/unexist/showcase-microservices-golang/domain"
)

type TodoProducer struct {
	conn sarama.SyncProducer
}

func NewTodoProducer() *TodoProducer {
	return &TodoProducer{}
}

func (producer *TodoProducer) Open(connectionString string) (err error) {
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer.conn, err = sarama.NewSyncProducer([]string{connectionString}, config)

	err = errtrace.Wrap(err)

	return
}

func (producer *TodoProducer) Publish(topic string, todo domain.Todo) {
	message, err := json.Marshal(todo)
	if nil != err {
		fmt.Println("Error: ", err.Error())

		panic(err)
	}

	/* Publish sync */
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, publishErr := producer.conn.SendMessage(msg)
	if publishErr != nil {
		fmt.Println("Error publish: ", publishErr.Error())
		panic(publishErr)
	}

	fmt.Println("Partition: ", partition)
	fmt.Println("Offset: ", offset)
	fmt.Println(todo)
}

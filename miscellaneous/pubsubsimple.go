package main

import (
	"fmt"
	"sync"
)

type Broker struct {
	Topic map[string][]Consumer
	mutex sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		Topic: make(map[string][]Consumer),
	}
}
func (b *Broker) Subscribe(topic string, c Consumer) {
	b.mutex.Lock()
	b.Topic[topic] = append(b.Topic[topic], c)
	b.mutex.Unlock()
}
func (b *Broker) Publish(topic, message string) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	consumers := b.Topic[topic]
	for _, consumer := range consumers {
		go func() {
			consumer.DataChan <- message
		}()
	}
}

type Consumer struct {
	Id       string
	DataChan chan string
}

func (c *Consumer) Consume(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(c.Id, "-is consuming")
	fmt.Println(<-c.DataChan)
}
func NewConsumer(id string) Consumer {
	return Consumer{
		Id:       id,
		DataChan: make(chan string, 1),
	}
}

func main() {
	broker := NewBroker()
	prajwalp := NewConsumer("prajwalp")
	sahils := NewConsumer("sahils")
	raj := NewConsumer("rajl")
	broker.Subscribe("npci", prajwalp)
	broker.Subscribe("npci", sahils)
	broker.Subscribe("tcs", raj)
	var wg sync.WaitGroup
	wg.Add(3)
	go prajwalp.Consume(&wg)
	go raj.Consume(&wg)
	go sahils.Consume(&wg)
	broker.Publish("npci", "npci is hiring golang developers")
	broker.Publish("tcs", "tcs is hiring react developers")
	wg.Wait()
}

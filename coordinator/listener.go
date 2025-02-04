// -----------------------------------------------------------------------------
// Coordinator package used for defining queue listener and event aggregator.
// -----------------------------------------------------------------------------
package coordinator

import (
	"bytes"
	"encoding/gob"
	"godistributed-rabbitmq/common"
	"godistributed-rabbitmq/common/dto"
	"log"

	"github.com/streadway/amqp"
)

// -----------------------------------------------------------------------------
// SensorListener - Struct that contains the logic for discovering
// sensor data queues, receiving messages and translating them into events.
// -----------------------------------------------------------------------------
type SensorListener struct {
	connection *amqp.Connection                // Connection to RabbitMQ.
	channel    *amqp.Channel                   // Channel created with the connection.
	sources    map[string]<-chan amqp.Delivery // Subscribed sensors with their delivery channels.
	aggregator *EventAggregator                // Used for funneling events.
}

// -----------------------------------------------------------------------------
// NewListener - Creates new sensor listener.
// -----------------------------------------------------------------------------
func NewListener(aggregator *EventAggregator) *SensorListener {
	listener := SensorListener{
		sources:    make(map[string]<-chan amqp.Delivery),
		aggregator: aggregator,
	}

	listener.connection, listener.channel = common.GetChannel(common.URL_GUEST)
	return &listener
}

// -----------------------------------------------------------------------------
// Start - Method used for starting sensor observation process.
// SensorListener will receive advertisement messages when the new sensor
// gets pluged into the system.
// -----------------------------------------------------------------------------
func (listener *SensorListener) Start() {
	// Passing "" will result with unique queue name creation by RabbitMQ.
	queue := common.GetQueue("", listener.channel, true)

	// Rebind queue from default exchange to the fanout.
	listener.channel.QueueBind(queue.Name, "", common.FANOUT_EXCHANGE, false, nil)

	// Receive sensor advertisement message when the new sensor is up and running.
	advertisements, _ := listener.channel.Consume(
		queue.Name, "", true, false, false, false, nil)

	// Send sensor discovery request.
	listener.DiscoveryRequest()
	log.Println("SensorListener: Listening for new sensors")

	for advertisement := range advertisements {
		log.Println("SensorListener: New sensor received")

		sensorName := string(advertisement.Body)
		listener.aggregator.Publish(common.SENSOR_DISCOVER_EVENT, sensorName)
		sensor, _ := listener.channel.Consume(sensorName, "", true, false, false, false, nil)

		if listener.sources[sensorName] == nil {
			listener.sources[sensorName] = sensor

			// Launch goroutine for observing incoming sensor readouts.
			go listener.observe(sensor)
		}
	}
}

// -----------------------------------------------------------------------------
// Stop - Method used for stoping sensor observation process.
// SensorListener will receive shutdown messages when the sensor
// gets pluged out of the system.
// -----------------------------------------------------------------------------
func (listener *SensorListener) Stop() {
	defer listener.channel.Close()
	defer listener.connection.Close()
}

// -----------------------------------------------------------------------------
// DiscoveryRequest - Method used for discovering already present sensors.
// -----------------------------------------------------------------------------
func (listener *SensorListener) DiscoveryRequest() {
	// Using fanout to send messages to every queue bound to this exchange.
	listener.channel.ExchangeDeclare(
		common.DISCOVERY_EXCHANGE, common.FANOUT, false, false, false, false, nil)

	log.Println("SensorListener: DiscoveryRequest sent")
	listener.channel.Publish(
		common.DISCOVERY_EXCHANGE, common.DISCOVERY_QUEUE, false, false, amqp.Publishing{})
}

// -----------------------------------------------------------------------------
// observe - Method used for observing incoming messages
// received from the subscribed sensor channel.
//
// sensor - incoming channel of amqp.Delivery containing sensor messages.
// -----------------------------------------------------------------------------
func (listener *SensorListener) observe(sensor <-chan amqp.Delivery) {
	for payload := range sensor {
		reader := bytes.NewReader(payload.Body)
		decoder := gob.NewDecoder(reader)

		readout := new(dto.Readout)
		decoder.Decode(readout)
		log.Printf("Received readout: %v\n", readout)

		// Event is prefixed sensor name
		event := common.NewEvent(common.MESSAGE_RECEIVED_EVENT, payload.RoutingKey)
		data := dto.NewEventData(readout.Name, readout.Value, readout.Timestamp)
		listener.aggregator.Publish(event, *data)
	}
}

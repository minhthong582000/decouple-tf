package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	amqp "github.com/rabbitmq/amqp091-go"

	"cdk.tf/go/stack/pkg/ec2"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func tfApply(tfId string) error {
	pathToStateFiles := "./cdktf.out/stacks/" + tfId
	cmd := exec.Command(
		"bash",
		"-c",
		"terraform init; terraform apply -auto-approve", 
	)
	cmd.Dir = pathToStateFiles
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	// Print the output
	fmt.Println(string(stdout))
	return nil
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var instanceName string = string(d.Body)

			log.Printf("Received instance name: %s", instanceName)
			
			// Best practice: get tfId from user input
			tfId := "tf-deployment"
			app := cdktf.NewApp(nil)
			// Create EC2 instances
			stack := ec2.NewEC2Stack(instanceName)
			stack.CreateStack(app, tfId)
			app.Synth()
			// Run terraform apply
			err := tfApply(tfId)
			if err != nil {
				fmt.Println(err.Error())
				// Send message back to queue
				d.Nack(false, true)
			}

			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			// Finished proccessing request
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
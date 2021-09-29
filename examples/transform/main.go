package main

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/pkg/uuid"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clientA, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30501),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}
	clientB, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30502),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}
	clientC, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30503),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}
	clientD, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30504),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}
	clientE, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30505),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}
	clientF, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30506),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}
	clientG, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30507),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		errCh := make(chan error)
		eventsCh, err := clientD.SubscribeToEvents(ctx, ">", "", errCh)
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case err := <-errCh:
				log.Fatal(err)
				return
			case event, more := <-eventsCh:
				if !more {
					log.Println("client d, done")
					return
				}
				_, err := clientD.NewQueueMessage().
					SetChannel("queue.e").
					SetMetadata(event.Metadata).
					SetBody(event.Body).Send(ctx)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("client D transform events on channel %s to queue messages", event.Channel)

			case <-ctx.Done():
				return
			}
		}

	}()
	go func() {
		for {
			results, err := clientE.ReceiveQueueMessages(ctx, &kubemq.ReceiveQueueMessagesRequest{
				RequestID:           "id",
				ClientID:            "client e",
				Channel:             "queue",
				MaxNumberOfMessages: 1,
				WaitTimeSeconds:     60,
				IsPeak:              false,
			})
			if err != nil {
				log.Fatal(err)
			}
			for _, msg := range results.Messages {
				log.Printf("client e queue message received, %s", string(msg.Body))
			}
		}
	}()
	go func() {
		for {
			results, err := clientF.ReceiveQueueMessages(ctx, &kubemq.ReceiveQueueMessagesRequest{
				RequestID:           "id",
				ClientID:            "client f",
				Channel:             "queue",
				MaxNumberOfMessages: 1,
				WaitTimeSeconds:     60,
				IsPeak:              false,
			})
			if err != nil {
				log.Fatal(err)
			}
			for _, msg := range results.Messages {
				log.Printf("client f queue message received, %s", string(msg.Body))
			}
		}
	}()
	go func() {
		for {
			results, err := clientG.ReceiveQueueMessages(ctx, &kubemq.ReceiveQueueMessagesRequest{
				RequestID:           "id",
				ClientID:            "client g",
				Channel:             "queue",
				MaxNumberOfMessages: 1,
				WaitTimeSeconds:     60,
				IsPeak:              false,
			})
			if err != nil {
				log.Fatal(err)
			}
			for _, msg := range results.Messages {
				log.Printf("client g queue message received, %s", string(msg.Body))
			}
		}
	}()
	// give some time to connect a receiver
	time.Sleep(1 * time.Second)
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)
	counter := 0
	for {
		counter++
		err := clientA.NewEvent().
			SetId("event").
			SetChannel("events.a").
			SetMetadata("").
			SetBody([]byte(fmt.Sprintf("client a send event %d", counter))).
			Send(ctx)
		if err != nil {
			log.Println(fmt.Sprintf("error sending event %d, error: %s", counter, err))
		}
		err = clientB.NewEvent().
			SetId("event").
			SetChannel("events.b").
			SetMetadata("").
			SetBody([]byte(fmt.Sprintf("client b send event %d", counter))).
			Send(ctx)
		if err != nil {
			log.Println(fmt.Sprintf("error sending event %d, error: %s", counter, err))
		}
		err = clientC.NewEvent().
			SetId("event").
			SetChannel("events.c").
			SetMetadata("").
			SetBody([]byte(fmt.Sprintf("client c send event %d", counter))).
			Send(ctx)
		if err != nil {
			log.Println(fmt.Sprintf("error sending event %d, error: %s", counter, err))
		}
		select {
		case <-gracefulShutdown:
			break
		default:
			time.Sleep(time.Second)
		}
	}
}

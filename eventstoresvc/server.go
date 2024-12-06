package main

import (
	"context"
	"event-driven-distributed-sys/cockroachdb/eventstorerepository"
	"flag"
	"fmt"
	"log"
	"net"

	"eventstore"
	"natsutil"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// publishEvent publishes an event to JetStream stream
func publishEvent(component *natsutil.NATSComponent, event *eventstore.Event) {
	// creates JetStreamContext to publish the event to JetStream stream
	JetStreamContext, _ := component.JetStreamContext()
	subject := event.EventType
	eventMsg, _ := []byte(event.EventData)
	// publish the event to JetStream stream
	JetStreamContext.Publish(subject, eventMsg)
	log.Println("Event published to JetStream stream on subject: ", subject)
}

// server implements eventstore.EventStoreServer interface
type server struct {
		eventstore.UnimplementedEventStoreServer
		repository eventstore.Repository
		nats 	 *natsutil.NATSComponent
}

// CreateEvent creates a new event to the event store
func (s *server) CreateEvent(ctx context.Context, eventRequest *eventstore.CreateEventRequest) (*eventstore.CreateEventResponse, error) {
		err := s.repository.CreateEvent(ctx, eventRequest.Event)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal error")
		}
		log.Println("Event created successfully")
		// publish the event to JetStream stream
		go publishEvent(s.nats, eventRequest.Event)
		return &eventstore.CreateEventResponse{IsSuccess: true, Error: ""}, nil
}

// GetEvents gets all events for the given aggregate and event
func (s *server) GetEvents(ctx context.Context, filter *eventstore.GetEventsRequest) (*eventstore.GetEventsResponse, error) {
		events, err := s.repository.GetEvents(ctx, filter)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal error")
		}
		return &eventstore.GetEventsResponse{Events: events}, nil
}

// GetEventsStream gets stream of events for the given event - not implemented
func (s *server) GetEventsStream(filter *eventstore.GetEventsRequest, stream eventstore.EventStore_GetEventsStreamServer) error {
		return nil
}

func getServer() *server {
		eventstoreDB, _ := eventstore.NewEventStoreDB()
		repository, _ := eventstorerepository.New(eventstoreDB.DB)
		natsComponent := natsutil.NewNATSComponent("eventstore-service")
		natsComponent.ConnectToServer(nats.DefaultURL)
		server := &server{
			repository: repository,
			nats: natsComponent,
		}
		return server
}


func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	server := getServer()
	eventstore.RegisterEventStoreServer(grpcServer, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

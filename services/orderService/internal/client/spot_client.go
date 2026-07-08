package client

import (
	"time"

	spotv1 "grpc-exchange/gen/spot"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewSpotClient(address string) (spotv1.SpotInstrumentServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		 grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, nil, err
	}

	return spotv1.NewSpotInstrumentServiceClient(conn), conn, nil
}
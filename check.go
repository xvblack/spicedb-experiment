package main

import (
	"context"
	"io"
	"log"
	"strconv"
	"time"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type secureMetadataCreds map[string]string

func (c secureMetadataCreds) RequireTransportSecurity() bool { return false }
func (c secureMetadataCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return c, nil
}

func main() {
	client, err := authzed.NewClient(
		"localhost:50051",
		grpc.WithPerRPCCredentials(secureMetadataCreds{"authorization": "Bearer " + "somerandomkeyhere"}),
		// grpcutil.WithBearerToken(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithInsecure(),
		// grpcutil.WithSystemCerts(grpcutil.SkipVerifyCA),
	)
	ctx := context.Background()

	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	start := time.Now()
	for k := 0; ; k++ {
		i := k % 3125
		this_round := time.Now()
		resp, err := client.LookupResources(
			ctx,
			&pb.LookupResourcesRequest{
				ResourceObjectType: "workday/profile",
				Permission:         "read",
				Subject: &v1.SubjectReference{
					Object: &v1.ObjectReference{
						ObjectType: "workday/user",
						ObjectId:   strconv.Itoa(i * 10),
					},
					// OptionalRelation: "",
				},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		count := 0
		for {
			_, err := resp.Recv()
			if err == io.EOF {
				break
			}
			count++
		}
		log.Printf("%d relations looked up in %dus", count, time.Now().Sub(this_round).Microseconds())
		// client.CheckPermission(
		// 	ctx,
		// 	&pb.CheckPermissionRequest{
		// 		// Consistency: &v1.Consistency{},
		// 		Resource: &v1.ObjectReference{
		// 			ObjectType: "workday/profile",
		// 			ObjectId:   strconv.Itoa(i),
		// 		},
		// 		Permission: "read",
		// 		Subject: &v1.SubjectReference{
		// 			Object: &v1.ObjectReference{
		// 				ObjectType: "workday/user",
		// 				ObjectId:   strconv.Itoa(i),
		// 			},
		// 			// OptionalRelation: "",
		// 		},
		// 	},
		// )
	}
	end := time.Now()
	log.Printf("Spent %dus time", end.Sub(start).Microseconds()/100)
}

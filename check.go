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
)

func run_check(client *authzed.Client) {

	ctx := context.Background()
	start := time.Now()
	for k := 0; k < 100000; k++ {
		this_round := time.Now()
		i := k % 3125

		// i := 0
		resp, err := client.LookupResources(
			ctx,
			&pb.LookupResourcesRequest{
				// Consistency: Consistency_FullyConsistent{},
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
		// 			ObjectType: "workday/system_role",
		// 			ObjectId:   "SINGLETON",
		// 		},
		// 		Permission: "any_profile_read",
		// 		Subject: &v1.SubjectReference{
		// 			Object: &v1.ObjectReference{
		// 				ObjectType: "workday/user",
		// 				ObjectId:   strconv.Itoa(i),
		// 			},
		// 			// OptionalRelation: "",
		// 		},
		// 	},
		// )
		// log.Printf("system role checked in %dus", time.Now().Sub(this_round).Microseconds())

	}
	end := time.Now()
	log.Printf("Spent %dus time", end.Sub(start).Microseconds()/100)
}

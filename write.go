package main

import (
	"context"
	"log"
	"strconv"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
)

func run_write(client *authzed.Client) {
	ctx := context.Background()

	relations := make([]*v1.RelationshipUpdate, 0)
	// branchSize := 5
	leafSize := 10
	for i := 0; i < 3125; i++ {
		org := &v1.ObjectReference{
			ObjectType: "workday/org",
			ObjectId:   strconv.Itoa(i),
		}
		admin := &v1.ObjectReference{
			ObjectType: "workday/user",
			ObjectId:   strconv.Itoa(i * leafSize),
		}
		relations = append(relations, &pb.RelationshipUpdate{
			Operation: pb.RelationshipUpdate_OPERATION_TOUCH,
			Relationship: &v1.Relationship{
				Resource: org,
				Relation: "admin",
				Subject: &v1.SubjectReference{
					Object: admin,
				},
			},
		})
		if i != 0 {
			relations = append(relations, &pb.RelationshipUpdate{
				Operation: pb.RelationshipUpdate_OPERATION_TOUCH,
				Relationship: &v1.Relationship{
					Resource: org,
					Relation: "parent",
					Subject: &v1.SubjectReference{
						Object: &v1.ObjectReference{
							ObjectType: "workday/org",
							ObjectId:   strconv.Itoa(i / 5),
						},
					},
				},
			})
		}
		for leaf := 0; leaf < leafSize; leaf++ {
			relations = append(relations, &pb.RelationshipUpdate{
				Operation: pb.RelationshipUpdate_OPERATION_TOUCH,
				Relationship: &v1.Relationship{
					Resource: &v1.ObjectReference{
						ObjectType: "workday/profile",
						ObjectId:   strconv.Itoa(i*leafSize + leaf),
					},
					Relation: "belong_to",
					Subject: &v1.SubjectReference{
						Object: org,
					},
				},
			})
			relations = append(relations, &pb.RelationshipUpdate{
				Operation: pb.RelationshipUpdate_OPERATION_TOUCH,
				Relationship: &v1.Relationship{
					Resource: &v1.ObjectReference{
						ObjectType: "workday/system_role",
						ObjectId:   "SINGLETON",
					},
					Relation: "profile",
					Subject: &v1.SubjectReference{
						Object: &v1.ObjectReference{
							ObjectType: "workday/profile",
							ObjectId:   strconv.Itoa(i*leafSize + leaf),
						},
					},
				},
			})
		}
	}
	log.Printf("Total number of relations %d", len(relations))

	for i := 0; i < len(relations); i += 5000 {
		end := i + 5000
		if end > len(relations) {
			end = len(relations)
		}
		resp, err := client.WriteRelationships(ctx, &pb.WriteRelationshipsRequest{
			Updates: relations[i:end],
		})
		if err != nil {
			log.Fatalf("failed to check permission: %s", err)
		}
		log.Print(resp.String())

	}

	// resp, err = client.CheckPermission(ctx, &pb.CheckPermissionRequest{
	// 	Resource:   firstPost,
	// 	Permission: "write",
	// 	Subject:    emilia,
	// })
	// if err != nil {
	// 	log.Fatalf("failed to check permission: %s", err)
	// }
	// // resp.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION

	// resp, err = client.CheckPermission(ctx, &pb.CheckPermissionRequest{
	// 	Resource:   firstPost,
	// 	Permission: "read",
	// 	Subject:    beatrice,
	// })
	// if err != nil {
	// 	log.Fatalf("failed to check permission: %s", err)
	// }
	// // resp.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION

	// resp, err = client.CheckPermission(ctx, &pb.CheckPermissionRequest{
	// 	Resource:   firstPost,
	// 	Permission: "write",
	// 	Subject:    beatrice,
	// })
	// if err != nil {
	// 	log.Fatalf("failed to check permission: %s", err)
	// }
	// resp.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_NO_PERMISSION
}

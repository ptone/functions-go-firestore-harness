package main

import (
	"context"
	"fmt"
	"log"
	"p"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var projectID = "ptone-serverless"

func main() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// write your query here to match your trigger path
	q := client.Collection("cities").
		Where("population", ">", 10).
		OrderBy("population", firestore.Desc).
		Limit(10)
	qsnapIter := q.Snapshots(ctx)
	// Listen forever for changes to the query's results.
	lastValues := map[string]p.FirestoreValue{}
	for {
		qsnap, err := qsnapIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range qsnap.Changes {
			// You might check the type of change and only call the function with the right subset of events
			newValue := p.FirestoreValue{
				CreateTime: c.Doc.CreateTime,
				UpdateTime: c.Doc.UpdateTime,
				Name:       fmt.Sprintf("projects/%s/databases/(default)/documents/%s/%s", projectID, "cities", c.Doc.Ref.ID),
				Fields:     c.Doc.Data(),
			}
			e := p.FirestoreEvent{
				Value: newValue,
				// have not sorted out how to mock UpdateMask
			}

			if val, ok := lastValues[c.Doc.Ref.ID]; ok {
				e.OldValue = val
			}
			lastValues[c.Doc.Ref.ID] = newValue

			// call the cloud function with the change
			p.DocChange(context.Background(), e)
		}
	}
}

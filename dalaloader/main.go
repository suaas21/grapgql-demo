package main

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"log"
)

func main() {
	batchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		var results []*dataloader.Result

		// do your stuff
		return results
	}

	loader := dataloader.NewBatchedLoader(batchFn)

	thunk := loader.Load(context.TODO(), dataloader.StringKey("key1"))

	result, err := thunk()
	if err != nil {

	}

	log.Printf("value: %#v", result)
}

package example

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/grpc/status"
)

func Save(data string, esc ElementServiceClient) {
	element := &Element{}
	if err := json.Unmarshal([]byte(data), element); err != nil {
		fmt.Println("Unable to decode the give data", data, "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	element, err := esc.PersistElement(ctx, element)
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			panic("error is not grpc")
		}
		fmt.Println("Unable to save the given data grpc code:", status.Code(), "message:", status.Message())
		os.Exit(1)
	}
	fmt.Println("Element Created:")
	element.print()
}

func List(filter string, esc ElementServiceClient) {
	elef := &ElementFilter{}
	if filter != "" {
		if err := json.Unmarshal([]byte(filter), elef); err != nil {
			fmt.Println("Unable to decode the give data", elef, "error", err)
			os.Exit(1)
		}
	}

	ctx := context.Background()
	eles, err := esc.ListElement(ctx, elef)
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			panic("error is not grpc")
		}
		fmt.Println("Unable to list the element grpc code:", status.Code(), "message:", status.Message())
		os.Exit(1)
	}
	if eles == nil || len(eles.Elements) == 0 {
		fmt.Println("No element found")
		os.Exit(0)
	}
	fmt.Println("Elements:")
	eles.print()
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/customvision/training"
)

var (
    training_key string = "1aa008d180f7491b8fc10498e08f6eab"
    prediction_key string = "1aa008d180f7491b8fc10498e08f6eab"
    prediction_resource_id = "/subscriptions/65f4b813-c139-4b12-9523-5533498a0ff7/resourceGroups/resource_group0/providers/Microsoft.CognitiveServices/accounts/mhl-congitive-services"
    endpoint string = "https://mhl-congitive-services.cognitiveservices.azure.com/"
    project_name string = "Go Sample Project"
    iteration_publish_name = "classifyModel"
    sampleDataDirectory = "../backend/images"
)

func main() {
    fmt.Println("Creating project...")

    ctx := context.Background()

    trainer := training.New(training_key, endpoint)

    _, err := trainer.CreateProject(ctx, project_name, "sample project", nil, string(training.Multiclass), nil)
    if (err != nil) {
        log.Fatal(err)
    }
}
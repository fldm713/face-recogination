package customvision

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/customvision/prediction"
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/customvision/training"
	"github.com/gofrs/uuid"
)

var (
    training_key string = "1aa008d180f7491b8fc10498e08f6eab"
    prediction_key string = "1aa008d180f7491b8fc10498e08f6eab"
    prediction_resource_id = "/subscriptions/65f4b813-c139-4b12-9523-5533498a0ff7/resourceGroups/resource_group0/providers/Microsoft.CognitiveServices/accounts/mhl-congitive-services"
    endpoint string = "https://mhl-congitive-services.cognitiveservices.azure.com/"
    project_name string = "FaceRecognition"
    iteration_publish_name = "classifyModel"
    sampleDataDirectory = "./imagesFace"

    ctx context.Context = nil
    trainer training.BaseClient
    projectId uuid.UUID
    tagMap map[string]uuid.UUID
    predictionEndpoint string
    predictor prediction.BaseClient
)

func Init(projectID string, tags map[string]string) {
    ctx  = context.Background()
    trainer = training.New(training_key, endpoint)
    projectId, _ = uuid.FromString(projectID)
    tagMap = make(map[string]uuid.UUID)
    for k, v := range tags {
        tagMap[k], _ = uuid.FromString(v)
    }
    predictionEndpoint = "https://mhl-congitive-services.cognitiveservices.azure.com/customvision/v3.0/Prediction/" + projectID + "/classify/iterations/classifyModel/image"
    predictor = prediction.New(endpoint)
}

func CreateProject(name string, description string) uuid.UUID {
    project, err := trainer.CreateProject(ctx, name, description, nil, string(training.Multiclass), nil)
    if (err != nil) {
        log.Fatal(err)
    }

	return *project.ID
}

func CreateTag(projectID uuid.UUID, name string, description string) (training.Tag, error) {
    tag, err :=  trainer.CreateTag(ctx, projectID, name, description, string(training.Regular))
    return tag, err
}

func UpLoadImage(image []byte, tagName string) (err error) {
    tagId, ok := tagMap[tagName]
    if !ok {
        return errors.New("tag name not found in tag map")
    }

    imageData := ioutil.NopCloser(bytes.NewReader(image))
    _, err = trainer.CreateImagesFromData(ctx, projectId, imageData, []uuid.UUID{ tagId })
    
    return err
}

func UpLoadImages(projectID uuid.UUID, hemlockTag training.Tag, cherryTag training.Tag) {
    fmt.Println("Adding images...")
    cherryImages, err := ioutil.ReadDir(path.Join(sampleDataDirectory, "Cherry"))
    if err != nil {
        fmt.Println("Error finding Sample images")
    }

    hemLockImages, err := ioutil.ReadDir(path.Join(sampleDataDirectory, "Hemlock"))
    if err != nil {
        fmt.Println("Error finding Sample images")
    }

    for _, file := range hemLockImages {
        imageFile, _ := ioutil.ReadFile(path.Join(sampleDataDirectory, "Hemlock", file.Name()))
        imageData := ioutil.NopCloser(bytes.NewReader(imageFile))

        trainer.CreateImagesFromData(ctx, projectID, imageData, []uuid.UUID{ *hemlockTag.ID })
    }

    for _, file := range cherryImages {
        imageFile, _ := ioutil.ReadFile(path.Join(sampleDataDirectory, "Cherry", file.Name()))
        imageData := ioutil.NopCloser(bytes.NewReader(imageFile))
        trainer.CreateImagesFromData(ctx, projectID, imageData, []uuid.UUID{ *cherryTag.ID })
    }
}

func Train() (err error) {
    fmt.Println("Training...")
    var hour int32 = 1
    var forceTrain bool = false
    iteration, _ := trainer.TrainProject(ctx, projectId, "Regular", &hour, &forceTrain, "")
    for {
        if *iteration.Status != "Training" {
            break
        }
        fmt.Println("Training status: " + *iteration.Status)
        time.Sleep(1 * time.Second)
        iteration, _ = trainer.GetIteration(ctx, projectId, *iteration.ID)
    }
    fmt.Println("Training status: " + *iteration.Status)

    _, err = trainer.PublishIteration(ctx, projectId, *iteration.ID, iteration_publish_name, prediction_resource_id)

    return err
}

func Predict2() {
    fmt.Println("Predicting...")
    predictor := prediction.New(endpoint)

    testImageData, _ := ioutil.ReadFile(path.Join(sampleDataDirectory, "Test", "test_image.jpg"))
    // iterationId, _ := uuid.FromString("0be97a4b-6a10-407c-99cc-f9e40efdd847")
    // req, err := predictor.PredictImagePreparer(ctx, projectId, ioutil.NopCloser(bytes.NewReader(testImageData)), &iterationId, "")
    
    urlPath := "https://mhl-congitive-services.cognitiveservices.azure.com/customvision/v3.0/Prediction/" + projectId.String() + "/classify/iterations/classifyModel/image"
    req, err := http.NewRequest("POST", urlPath, bytes.NewReader(testImageData))
    req.Header.Set("Content-Type", "application/octet-stream")
    req.Header.Set("Prediction-Key", "1aa008d180f7491b8fc10498e08f6eab")
    if err != nil {
        panic(err)
    }

    resp, err := predictor.ClassifyImageSender(req)

    if err != nil {
        panic(err)
    }
    fmt.Println("response body")
    fmt.Printf("response body:%v\n", resp.Body)

    b, err := io.ReadAll(resp.Body)

    if err != nil {
        panic(err)
    }

    fmt.Printf("Predict response body: %v\n", string(b))

    // fmt.Println("Predicting...")
    // predictor := prediction.New(endpoint)

    // testImageData, _ := ioutil.ReadFile(path.Join(sampleDataDirectory, "Test", "test_image.jpg"))
    // results, _ := predictor.ClassifyImage(ctx, projectId, iteration_publish_name, ioutil.NopCloser(bytes.NewReader(testImageData)), "")

    // for _, prediction := range *results.Predictions    {
    //     fmt.Printf("\t%s: %.2f%%", *prediction.TagName, *prediction.Probability * 100)
    //     fmt.Println("")
    // }

}

func Predict(fileName string) (result string, err error) {
    fmt.Println("Predicting...")
    
    testImageData, _ := ioutil.ReadFile(fileName)
    req, err := http.NewRequest("POST", predictionEndpoint, bytes.NewReader(testImageData))
    if err != nil {
        fmt.Printf("=======%v\n", err)
        panic(err)
    }
    req.Header.Set("Content-Type", "application/octet-stream")
    req.Header.Set("Prediction-Key", "1aa008d180f7491b8fc10498e08f6eab")

    resp, err := predictor.ClassifyImageSender(req)

    if err != nil {
        fmt.Printf("=======%v\n", err)
        panic(err)
    }

    b, err := io.ReadAll(resp.Body)

    return string(b), err
    // return "", err
}



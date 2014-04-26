package slimgfast

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
)

type Job struct {
	ImageSource  ImageSource
	ImageRequest ImageRequest
	Result       chan []byte
	Error        chan error
}

type WorkerGroup struct {
	Transformers []Transformer
	NumWorkers   int
	jobs         chan Job
}

func (wg *WorkerGroup) Start() {
	wg.jobs = make(chan Job, wg.NumWorkers)
	for i := 0; i < wg.NumWorkers; i++ {
		go work(wg)
	}
}

func (wg *WorkerGroup) Close() {
	close(wg.jobs)
}

func (wg *WorkerGroup) Resize(imageSource *ImageSource, imageRequest *ImageRequest) ([]byte, error) {
	job := Job{
		ImageSource:  *imageSource,
		ImageRequest: *imageRequest,
		Result:       make(chan []byte),
		Error:        make(chan error),
	}
	defer close(job.Result)
	defer close(job.Error)

	wg.jobs <- job

	var resizedBytes []byte
	var err error
	select {
	case resizedBytes = <-job.Result:
	case err = <-job.Error:
	}
	return resizedBytes, err
}

// work consumes the job queue and sends results back on the job's result
// channel (and errors back on the job's error channel)
func work(wg *WorkerGroup) {
	for job := range wg.jobs {
		data, err := job.ImageSource.GetImageData(&job.ImageRequest)
		if err != nil {
			job.Error <- err
			return
		}
		resizedData, err := resizeImg(&job.ImageRequest, wg.Transformers, data)
		if err == nil {
			job.Result <- resizedData
		} else {
			job.Error <- err
		}
	}
}

func resizeImg(req *ImageRequest, transformers []Transformer, data []byte) ([]byte, error) {
	// Middle variable is format name that was used
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Println("Error decoding image", err)
		return nil, err
	}
	for _, transformer := range transformers {
		img, err = transformer.Transform(req, img)
		if err != nil {
			return nil, err
		}
	}
	var buf bytes.Buffer
	if err = jpeg.Encode(&buf, img, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

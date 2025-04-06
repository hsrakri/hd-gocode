package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"gocv.io/x/gocv"
)

type Block struct {
	X      int
	Y      int
	Width  int
	Height int
	Color  color.RGBA
}

type EyeTracker struct {
	webcam       *gocv.VideoCapture
	window       *gocv.Window
	faceCascade  *gocv.CascadeClassifier
	eyeCascade   *gocv.CascadeClassifier
	block        Block
	sensitivity  float64
	screenWidth  int
	screenHeight int
}

func NewEyeTracker() (*EyeTracker, error) {
	// Initialize webcam
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return nil, fmt.Errorf("error opening webcam: %v", err)
	}

	// Load cascade classifiers
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("error getting executable path: %v", err)
	}
	basePath := filepath.Dir(execPath)

	faceCascade := gocv.NewCascadeClassifier()
	if !faceCascade.Load(filepath.Join(basePath, "assets/haarcascade_frontalface_default.xml")) {
		return nil, fmt.Errorf("error loading face cascade classifier")
	}

	eyeCascade := gocv.NewCascadeClassifier()
	if !eyeCascade.Load(filepath.Join(basePath, "assets/haarcascade_eye.xml")) {
		return nil, fmt.Errorf("error loading eye cascade classifier")
	}

	// Create window
	window := gocv.NewWindow("Eye Tracker")

	// Initialize block
	block := Block{
		X:      320,
		Y:      240,
		Width:  50,
		Height: 50,
		Color:  color.RGBA{0, 0, 255, 255}, // Blue color
	}

	return &EyeTracker{
		webcam:       webcam,
		window:       window,
		faceCascade:  &faceCascade,
		eyeCascade:   &eyeCascade,
		block:        block,
		sensitivity:  2.0,
		screenWidth:  640,
		screenHeight: 480,
	}, nil
}

func (et *EyeTracker) Close() {
	et.webcam.Close()
	et.window.Close()
	et.faceCascade.Close()
	et.eyeCascade.Close()
}

func (et *EyeTracker) detectEyes(img gocv.Mat) ([]image.Point, error) {
	gray := gocv.NewMat()
	defer gray.Close()

	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	// Detect faces first
	faces := et.faceCascade.DetectMultiScale(gray)
	var eyeCenters []image.Point

	for _, face := range faces {
		faceROI := gray.Region(face)

		// Detect eyes within the face region
		eyes := et.eyeCascade.DetectMultiScale(faceROI)

		for _, eye := range eyes {
			// Calculate eye center relative to the whole image
			centerX := face.Min.X + eye.Min.X + eye.Dx()/2
			centerY := face.Min.Y + eye.Min.Y + eye.Dy()/2
			eyeCenters = append(eyeCenters, image.Point{X: centerX, Y: centerY})

			// Draw rectangle around detected eye
			gocv.Rectangle(&img,
				image.Rect(
					face.Min.X+eye.Min.X,
					face.Min.Y+eye.Min.Y,
					face.Min.X+eye.Max.X,
					face.Min.Y+eye.Max.Y,
				),
				color.RGBA{0, 255, 0, 0}, 2)
		}
	}

	return eyeCenters, nil
}

func (et *EyeTracker) updateBlockPosition(eyeCenters []image.Point) {
	if len(eyeCenters) == 0 {
		return
	}

	// Calculate average eye position
	var avgX, avgY int
	for _, center := range eyeCenters {
		avgX += center.X
		avgY += center.Y
	}
	avgX /= len(eyeCenters)
	avgY /= len(eyeCenters)

	// Map eye position to block movement
	targetX := int(float64(avgX) * et.sensitivity)
	targetY := int(float64(avgY) * et.sensitivity)

	// Apply bounds
	if targetX < 0 {
		targetX = 0
	} else if targetX > et.screenWidth-et.block.Width {
		targetX = et.screenWidth - et.block.Width
	}

	if targetY < 0 {
		targetY = 0
	} else if targetY > et.screenHeight-et.block.Height {
		targetY = et.screenHeight - et.block.Height
	}

	// Update block position with smooth movement
	et.block.X = targetX
	et.block.Y = targetY
}

func (et *EyeTracker) drawBlock(img gocv.Mat) {
	gocv.Rectangle(&img,
		image.Rect(et.block.X, et.block.Y, et.block.X+et.block.Width, et.block.Y+et.block.Height),
		et.block.Color,
		-1) // Filled rectangle
}

func (et *EyeTracker) Run() error {
	frame := gocv.NewMat()
	defer frame.Close()

	for {
		if ok := et.webcam.Read(&frame); !ok {
			return fmt.Errorf("error reading from webcam")
		}
		if frame.Empty() {
			continue
		}

		// Detect eyes
		eyeCenters, err := et.detectEyes(frame)
		if err != nil {
			return fmt.Errorf("error detecting eyes: %v", err)
		}

		// Update block position based on eye movement
		et.updateBlockPosition(eyeCenters)

		// Draw the block
		et.drawBlock(frame)

		// Show the frame
		et.window.IMShow(frame)

		// Check for key press
		key := et.window.WaitKey(1)
		if key == 27 { // ESC key
			break
		} else if key == 82 || key == 114 { // 'R' or 'r' key
			// Reset block position
			et.block.X = et.screenWidth/2 - et.block.Width/2
			et.block.Y = et.screenHeight/2 - et.block.Height/2
		} else if key == 43 { // '+' key
			et.sensitivity += 0.1
		} else if key == 45 { // '-' key
			et.sensitivity -= 0.1
			if et.sensitivity < 0.1 {
				et.sensitivity = 0.1
			}
		}
	}

	return nil
}

func main() {
	// Create eye tracker
	eyeTracker, err := NewEyeTracker()
	if err != nil {
		log.Fatalf("Error creating eye tracker: %v", err)
	}
	defer eyeTracker.Close()

	// Run the eye tracker
	if err := eyeTracker.Run(); err != nil {
		log.Fatalf("Error running eye tracker: %v", err)
	}
}

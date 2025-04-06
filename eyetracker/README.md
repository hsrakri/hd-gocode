# Eye Tracker

A Go program that uses computer vision to track eye movements and control a block on screen.

## Features

- Real-time eye tracking using webcam
- Block movement controlled by eye position
- Visual feedback of eye detection
- Adjustable sensitivity and tracking parameters

## Requirements

### System Requirements
- Go 1.16 or later
- OpenCV 4.x
- Webcam

### Installation Dependencies

For macOS:
```bash
# Method 1: Using Homebrew (for stable macOS versions)
brew install opencv
brew install pkg-config

# Method 2: Manual Installation (for newer macOS versions)
# 1. Download OpenCV from https://opencv.org/releases/
# 2. Extract and build from source:
mkdir build && cd build
cmake -DBUILD_SHARED_LIBS=OFF ..
make -j8
sudo make install

# Method 3: Using MacPorts
sudo port install opencv4
sudo port install pkgconfig
```

For Linux:
```bash
# Install OpenCV dependencies
sudo apt-get update
sudo apt-get install -y libopencv-dev

# Install pkg-config
sudo apt-get install -y pkg-config
```

### Installing GoCV
```bash
go get -u -d gocv.io/x/gocv
```

## Known Issues

### macOS Installation
- On newer versions of macOS (15+), Homebrew installation might fail
- In such cases, use either the manual installation method or MacPorts
- Make sure to set the following environment variables after installation:
  ```bash
  export PKG_CONFIG_PATH="/usr/local/lib/pkgconfig:/opt/local/lib/pkgconfig"
  export CGO_CPPFLAGS="-I/usr/local/include -I/opt/local/include"
  export CGO_LDFLAGS="-L/usr/local/lib -L/opt/local/lib"
  ```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/hd-gocode.git
cd hd-gocode/eyetracker
```

2. Install dependencies (see Installation Dependencies section above)

3. Initialize Go module:
```bash
go mod init eyetracker
go mod tidy
```

4. Run the program:
```bash
go run main.go
```

## Usage

1. Launch the program
2. Position yourself in front of the webcam
3. The program will detect your eyes and display the tracking window
4. A blue block will move based on your eye position
5. Press 'q' to quit the program

## Controls

- ESC: Exit the program
- R: Reset block position
- +/-: Adjust tracking sensitivity

## Troubleshooting

### Common Issues

1. **OpenCV Not Found**
   - Verify OpenCV installation
   - Check PKG_CONFIG_PATH environment variable
   - Ensure pkg-config is installed

2. **Webcam Access**
   - Grant camera permissions to the terminal/IDE
   - Verify webcam is not in use by another application
   - Check webcam index (default is 0)

3. **Build Errors**
   - Ensure all environment variables are set correctly
   - Try rebuilding with `go clean -cache && go build`

## How it Works

1. The program captures video from your webcam
2. Uses Haar Cascade Classifier to detect faces and eyes
3. Tracks the center point of detected eyes
4. Maps eye position to block movement coordinates
5. Updates the display in real-time

## Contributing

Contributions are welcome! Please feel free to submit pull requests with improvements or bug fixes.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
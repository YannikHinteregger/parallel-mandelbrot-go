# Parallel Mandelbrot using Go

This project features the computation and rendering of the Mandelbrot set to show off the concurrency capabilities of Go. This code is based on the work of [this repository](https://github.com/GiselaMD/parallel-mandelbrot-go).

To run the project use `go run .`

To render the mandelbrot set this projects uses the pixel library which does not support windows.
Consider using Ubuntu and ensure that these packages are installed: `apt install gcc libgl1-mesa-dev xorg-dev`

If issues occur refer to the (Github repository of Pixel](https://github.com/faiface/pixel)

![parallel-mandelbrot-go](https://user-images.githubusercontent.com/34191327/122488349-afc5bf00-cfb3-11eb-9767-e6724b05078c.gif)

![mandelbrot fractal](https://user-images.githubusercontent.com/34191327/122481879-ea285f80-cfa5-11eb-985a-7b3e7c73a88d.png)
![mandelbrot fractal blue](https://user-images.githubusercontent.com/34191327/122482067-412e3480-cfa6-11eb-8505-f76ce7d2e3a2.png)


package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hschendel/stl"
)

const (
	width  = 800
	height = 800
)

var (
	transformLoc int32
)

type mesh struct {
	vao      uint32
	vbo      uint32
	points   []float32
	tricount uint32
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	vertexShaderSource, err := os.ReadFile("vertexShader.vert")
	if err != nil {
		panic(err)
	}

	fragmentShaderSource, err := os.ReadFile("fragmentShader.frag")
	if err != nil {
		panic(err)
	}

	program := initOpenGL(vertexShaderSource, fragmentShaderSource)

	stlread, err := stl.ReadFile("assets/3dbenchy-1.stl")
	if err != nil {
		panic(err)
	}

	stlread.ScaleLinearDowntoSizeBox(stl.Vec3{1.5, 1.5, 1.5})
	stlread.MoveToPositive()
	stlread.Translate(stl.Vec3{-0.5, -0.5, -0.5})
	tricount := uint32(len(stlread.Triangles))
	benchyPoints := make([]float32, tricount*9)

	for i := range benchyPoints {
		benchyPoints[i] = stlread.Triangles[i/9].Vertices[i%9/3][i%3]
	}

	angle := float32(0)

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	rotate := mgl32.Ident4()
	rotateUniform := gl.GetUniformLocation(program, gl.Str("rotate\x00"))
	gl.UniformMatrix4fv(rotateUniform, 1, false, &rotate[0])

	vao, vbo := makeVAO(benchyPoints)

	benchyMesh := mesh{
		vao:      vao,
		vbo:      vbo,
		points:   benchyPoints,
		tricount: tricount,
	}

	for i := 0; i < int(benchyMesh.tricount)*3; i++ {
		rotated := mgl32.QuatRotate(-1*math.Pi/2, mgl32.Vec3{1, 0, 0}).Rotate(mgl32.Vec3{benchyMesh.points[i*3], benchyMesh.points[i*3+1], benchyMesh.points[i*3+2]})
		benchyMesh.points[i*3] = rotated[0]
		benchyMesh.points[i*3+1] = rotated[1]
		benchyMesh.points[i*3+2] = rotated[2]
	}

	for !window.ShouldClose() {

		rotate = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		gl.UniformMatrix4fv(rotateUniform, 1, false, &rotate[0])
		gl.BindBuffer(gl.ARRAY_BUFFER, benchyMesh.vbo)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(benchyMesh.points)*4, gl.Ptr(benchyMesh.points))
		draw(&benchyMesh, window, program)
		angle += 0.01
	}
}

func initGlfw() *glfw.Window {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Floating, glfw.False)

	window, err := glfw.CreateWindow(width, height, "OpenGL Renderer", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	return window
}

func initOpenGL(vertexShaderSource []byte, fragmentShaderSource []byte) uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(string(vertexShaderSource), gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(string(fragmentShaderSource), gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	//glfw.SwapInterval(0)
	gl.Viewport(0, 0, int32(width), int32(height))

	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func draw(mesh *mesh, window *glfw.Window, program uint32) {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(mesh.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(mesh.tricount)*3)

	glfw.PollEvents()
	window.SwapBuffers()
}

func makeVAO(points []float32) (uint32, uint32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, int32(12), 0)
	gl.EnableVertexAttribArray(0)

	return vao, vbo
}

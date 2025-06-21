package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func initOpenGL() uint32 {
	vertexShaderSourceByte, err := os.ReadFile("vertexShader.vert")
	if err != nil {
		panic(err)
	}

	fragmentShaderSourceByte, err := os.ReadFile("fragmentShader.frag")
	if err != nil {
		panic(err)
	}

	vertexShaderSource := strings.TrimSpace(string(vertexShaderSourceByte)) + "\x00"
	fragmentShaderSource := strings.TrimSpace(string(fragmentShaderSourceByte)) + "\x00"
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
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

func draw(scene scene, window *glfw.Window, program uint32) {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	if len(scene.objects) > 0 {
		for _, mesh := range scene.objects {
			gl.BindVertexArray(mesh.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(mesh.Vertices)))
		}
	}
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

package main

import (
        "log"
        "runtime"

        "github.com/go-gl/gl/v4.6-core/gl"
        "github.com/go-gl/glfw/v3.3/glfw"
        "github.com/hschendel/stl"
)

const (
        width  = 800
        height = 600
)

var ()

func main() {
        runtime.LockOSThread()

        window := initGlfw()
        defer glfw.Terminate()

        program := initOpenGL()

        stlread, err := stl.ReadFile("assets/3dbenchy-1.stl")
        if err != nil {
                panic(err)
        }

        mesh := [len(stlread.Triangles)*9]

        for i, t := range stlread.Triangles {
                t.Vertices
        }

        for !window.ShouldClose() {
                draw
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

        window, err := glfw.CreateWindow(width, height, "OpenGL Renderer", nil, nil)
        if err != nil {
                panic(err)
        }

        window.MakeContextCurrent()

        return window
}

func initOpenGL() uint32 {
        if err := gl.Init(); err != nil {
                panic(err)
        }
        version := gl.GoStr(gl.GetString(gl.VERSION))
        log.Println("OpenGL version", version)

        prog := gl.CreateProgram()
        gl.LinkProgram(prog)
        //glfw.SwapInterval(0)
        return prog
}

func draw(triangles []*stl.Triangle, window *glfw.Window, program uint32) {
        gl.ClearColor(0.0, 0.0, 0.0, 1.0)
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        gl.UseProgram(program)

        for x := range triangles {
                for _, c := range cells[x] {
                        c.draw()
                }
        }

        glfw.PollEvents()
        window.SwapBuffers()
}

func (triangle *stl.Triangle) draw() {
        gl.BindVertexArray(c.vao)
        gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(len(squareConst)/6))
}

func makeVAO(points []float32) (uint32, uint32) {
        var vbo uint32
        gl.GenBuffers(1, &vbo)
        gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
        gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

        var vao uint32
        gl.GenVertexArrays(1, &vao)
        gl.BindVertexArray(vao)
        gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, int32(24), 0)
        gl.EnableVertexAttribArray(0)
        gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, int32(24), 12)
        gl.EnableVertexAttribArray(1)

        return vao, vbo
}

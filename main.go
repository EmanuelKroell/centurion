package main

import (
	"math"
	"runtime"

	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/math32"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width    = 800
	height   = 800
	invINDEX = math.MaxUint32
)

var (
	transformLoc int32
)

type object struct {
	vao      uint32
	vbo      uint32
	position mgl32.Vec3
	rotation mgl32.Vec3
	scale    mgl32.Vec3
	Vertices []float32
	Uvs      []float32
	Normals  []float32
}

type scene struct {
	objects []object
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	modelImport, err := obj.Decode("assets/Bamboo_House/Bambo_House.obj", "assets/Bamboo_House.mtl")
	if err != nil {
		panic(err)
	}

	objects := make([]object, 0, len(modelImport.Objects))

	for _, o := range modelImport.Objects {
		obj := makeObject(modelImport, &o)
		objects = append(objects, obj)
	}

	scene := scene{
		objects: objects,
	}

	angle := 0.0

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{15, 5, 5}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	for !window.ShouldClose() {
		angle += 0.005
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		draw(scene, window, program)
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
	glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)

	window, err := glfw.CreateWindow(width, height, "OpenGL Renderer", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	return window
}

func makeObject(model *obj.Decoder, o *obj.Object) object {
	Vertices := math32.NewArrayF32(0, model.Vertices.Len())
	Uvs := math32.NewArrayF32(0, model.Uvs.Len())
	Normals := math32.NewArrayF32(0, model.Normals.Len())

	copyVertex := func(face *obj.Face, idx int) {
		var vec2 math32.Vector2
		var vec3 math32.Vector3

		model.Vertices.GetVector3(3*face.Vertices[idx], &vec3)
		Vertices.AppendVector3(&vec3)

		if face.Normals[idx] != invINDEX {
			model.Normals.GetVector3(3*face.Normals[idx], &vec3)
			Normals.AppendVector3(&vec3)
		}

		if face.Uvs[idx] != invINDEX {
			model.Vertices.GetVector2(2*face.Uvs[idx], &vec2)
			Uvs.AppendVector2(&vec2)

		}
	}

	for _, f := range o.Faces {
		for idx := 1; idx < len(f.Vertices)-1; idx++ {
			copyVertex(&f, 0)
			copyVertex(&f, idx)
			copyVertex(&f, idx+1)
		}

	}

	vao, vbo := makeVAO(Vertices)

	return object{
		vao:      vao,
		vbo:      vbo,
		Vertices: Vertices,
		Uvs:      Uvs,
		Normals:  Normals,
	}
}

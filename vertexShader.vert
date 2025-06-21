#version 460 core

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

layout(location = 0) in vec3 inPosition;
out float color;

void main() {
    gl_Position = projection * camera * model * vec4(inPosition, 1.0);
    color = sqrt(inPosition.x*inPosition.x + inPosition.y*inPosition.y + inPosition.z*inPosition.z)/20 + 0.2;
}
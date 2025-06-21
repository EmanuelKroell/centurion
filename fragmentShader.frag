#version 460
in float color;
out vec4 frag_colour;

void main() {
    frag_colour = vec4(color+0.2, color+0.1, color+0.2, 1.0);
}
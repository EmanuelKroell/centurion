#version 460
in float color;
out vec4 frag_colour;

void main() {
    frag_colour = vec4(color, color, color, 1.0);
}

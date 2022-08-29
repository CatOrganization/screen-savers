#version 120

// Input vertex attributes (from vertex shader)
varying vec2 fragTexCoord;
varying vec4 fragColor;

uniform sampler2D texture0;

void main()
{
    vec4 mask_color = texture2D(texture0, fragTexCoord);
    float val = mask_color.r;

    gl_FragColor = vec4(val, val, val, 1-max(0.1, val));
}
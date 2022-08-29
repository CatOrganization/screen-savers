#version 120

// Input vertex attributes (from vertex shader)
varying vec2 fragTexCoord;
varying vec4 fragColor;

uniform float scale;

uniform sampler2D texture0;
uniform sampler2D texture1;

void main()
{
    //vec2 screen_pos = vec2(gl_FragCoord.x / scale, (screen_height - gl_FragCoord.y) / scale);
//    vec2 screen_pos = gl_FragCoord.xy;
    vec4 mask_color = texture2D(texture0, fragTexCoord);

    float val = mask_color.r;
    gl_FragColor = vec4(val, val, val, 1-max(0.1, val));
}
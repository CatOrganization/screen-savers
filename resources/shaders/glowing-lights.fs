#version 120

#define MAX_LIGHTS 128

uniform float screen_height;
uniform float scale;

uniform float num_lights; // TODO: why don't uniform ints work??
uniform float max_light_distance;
uniform vec2 light_positions[MAX_LIGHTS];
uniform float light_intensities[MAX_LIGHTS];


void main()
{
    vec2 screen_pos = vec2(gl_FragCoord.x / scale, (screen_height - gl_FragCoord.y) / scale);

    float brightness = 0;


    for (int n = 0; n < num_lights; n++)
    {
        // Light is off, move on
        if (light_intensities[n] == 0)
        {
            continue;
        }

        float dx = screen_pos.x - light_positions[n].x;
        float dy = screen_pos.y - light_positions[n].y;
        float dist_to_light_sq = (dx*dx) + (dy*dy);

        float actual_light_max_distance = max_light_distance * light_intensities[n];

        // We're too far away from this light source, move on to the next
        if (dist_to_light_sq > (actual_light_max_distance*actual_light_max_distance))
        {
            continue;
        }

        float new_brightness = 1 - (dist_to_light_sq / (actual_light_max_distance*actual_light_max_distance));
        if (new_brightness > brightness)
        {
            brightness = new_brightness;
        }
    }

    gl_FragColor = vec4(0, 0, 0, 1-brightness);

}
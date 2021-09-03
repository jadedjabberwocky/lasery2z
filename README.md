# lasery2z
Tool to modify Snapmaker laser GCode for a rotary module to include z movements.

## Background
Snapmaker's Luban tool to generate 4D laser projects only uses two axes:

* y (backwards / forwards) and
* b (rotation of rotary module)

The X axis is unused - the laser is always positioned in the centre of the
job job.

The Z axis is unused - the tool assumes a cylindrical shape to be lasered.

It is this Z axis that we wish to exploit. If we still assume a circular
cross-section of the job, but relax the requirement for parallel sides,
then by adjusting the z axis based on the y axis, we can follow contours
on the job, for example a wine glass.

## Inputs
The inputs to the tool are:
* GCode file for cylindrical object
* Mechanism for mapping Y values to Z values.
  For now this will take the form of a black and white image.


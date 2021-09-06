package main

type Config struct {
	InputFilename     string `kong:"name='if',short='i',help='Input filename. Use - for stdin'"`
	OutputFilename    string `kong:"name='of',short='o',help='Output filename. Use - for stdout'"`
	MapOutputFilename string `kong:"name='mof',help='Optional filename to output map to'"`

	Help struct{} `kong:"cmd,default"`

	ImageMap struct {
		ImageInputFilename string  `kong:"name='iif',help='Filename of image to use for mapping',required"`
		ImageWidth         float64 `kong:"name='iw',help='Width of image in same units used in gcode',required"`
		ImageHeight        float64 `kong:"name='ih',help='Optional height of image in same units used in gcode'"`
	} `kong:"cmd,aliases='im',help='Use an image to map Y to Z'"`
}

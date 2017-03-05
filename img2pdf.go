package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func absPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return abs
}
func noExtFilename(path string) string {
	name := filepath.Base(path)
	noExt := name[:len(name)-len(filepath.Ext(name))]
	return noExt
}

func main() {
	flag.Usage = func() {
		fmt.Println(`
Usage: conv2pdf [OPTION]... SOURCE...
Description: Image to PDF converter (depend on Inkscape)
Options:`[1:])
		flag.PrintDefaults()
	}

	dir := flag.String("d", "", "output to the specified directory")
	force := flag.Bool("f", false, "do not prompt before overwriting")

	flag.Parse()
	if *dir != "" {
		*dir = absPath(*dir)
	}

	reader := bufio.NewReader(os.Stdin)
	for i, val := range flag.Args() {
		input := absPath(val)
		noExt := noExtFilename(input)
		outdir := filepath.Dir(input)
		if *dir != "" {
			outdir = *dir
		}

		output := filepath.Join(outdir, noExt+".pdf")
		fmt.Printf("[%[1]v] input : %[2]v\n[%[1]v] output: %[3]v\n", i, input, output)
		if exist(output) && !*force {
			fmt.Printf("overwrite %v? (y/n [n]) ", output)
			if line, err := reader.ReadString('\n'); err != nil {
				panic(err)
			} else if line[0] != 'y' {
				fmt.Println("not overwrite")
				continue
			}
		}
		out, _ := exec.Command("inkscape", input, "--export-pdf", output).CombinedOutput()
		println(string(out))
	}
}

package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

type CtrlSeqWriter struct {
	writer       io.Writer
	startCtrlSeq string
	endCtrlSeq   string
}

func main() {
	// fmt.Println("starting zsh")

	stdoutWriter := &CtrlSeqWriter{
		writer:       os.Stdout,
		startCtrlSeq: "<out>",
		endCtrlSeq:   "</out>",
	}

	stderrWriter := &CtrlSeqWriter{
		writer:       os.Stderr,
		startCtrlSeq: "<err> ",
		endCtrlSeq:   "</err>",
	}

	// -i forces interactive mode
	// -l forces a login shell
	// cmd := exec.Command("/bin/bash", "-i", "-l")
	cmd := exec.Command("/bin/zsh", "-i")
	cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal("fatal error starting ", err)
	}

	go func() {
		// This blocks until the pipe closes.
		_, err := io.Copy(stdoutWriter, stdout)
		if err != nil {
			log.Fatal("copy error", err)
		}
	}()

	go func() {
		// This blocks until the pipe closes.
		_, err := io.Copy(stderrWriter, stderr)
		if err != nil {
			log.Fatal("copy error", err)
		}
	}()

	// buf := make([]byte, 1000)
	// fmt.Println("Trying to read stdout pipe")
	// i := 0
	// for {
	// 	_, err := stdout.Read(buf)
	// 	if err != nil {
	// 		log.Fatal("Fatal pipe read", err)
	// 	}
	// 	// fmt.Fprint(os.Stdout, string(buf))
	// 	os.Stdout.Write(buf)
	// 	fmt.Fprint(os.Stdout, i)
	// 	// fmt.Println("Read ", n, " bytes to stdout")
	// 	// fmt.Println("Read ", n, " bytes to stdout", string(buf))
	// 	i++
	// }

	if err := cmd.Wait(); err != nil {
		log.Fatal("fatal error waiting ", err)
	}

	// cmd.Stdin = os.Stdin
	// cmd.Stdout = stdoutWriter
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func (w *CtrlSeqWriter) Write(p []byte) (n int, err error) {
	_, err = io.WriteString(w.writer, w.startCtrlSeq)
	if err != nil {
		return
	}
	bodyLen, err := w.writer.Write(p)
	if err != nil {
		return
	}
	_, err = io.WriteString(w.writer, w.endCtrlSeq)
	if err != nil {
		return
	}
	// Note that you get an error if you return the total length.
	n = bodyLen
	return
}

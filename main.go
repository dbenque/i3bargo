package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type I3Element struct {
	Name     string `json:"name"`
	Instance string `json:"instance,omitempty"`
	Color    string `json:"color,omitempty"`
	Markup   string `json:"markup"`
	FullText string `json:"full_text"`
}
type I3Line []I3Element

var kubeCurrentContext string
var kubeCurrentNamespace string

func kubedata() {
	for {
		kubeCurrentContext = ""
		kubeCurrentNamespace = ""
		if currentContext, err := exec.Command("kubectl", "config", "current-context").Output(); err == nil {
			kubeCurrentContext = string(currentContext)
			kubeCurrentContext = strings.TrimSuffix(kubeCurrentContext, "\n")
			cc, _ := exec.Command("kubectl", "config", "view", "-ojsonpath={..current-context}").Output()
			ccstr := strings.TrimSuffix(string(cc), "\n")
			currentNS, _ := exec.Command("kubectl", "config", "view", "-ojsonpath={.Contexts[?(@.Name==\""+ccstr+"\")]..namespace}").Output()
			kubeCurrentNamespace = string(currentNS)
			kubeCurrentNamespace = strings.TrimSuffix(kubeCurrentNamespace, "\n")

		} else {
			fmt.Printf("error kube: %v", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	cmd := exec.Command("i3status")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	go kubedata()

	go func() {
		r := bufio.NewReader(stdout)
		versionline, _, _ := r.ReadLine()
		fmt.Println(string(versionline))
		bra, _, _ := r.ReadLine() // opening brackets
		fmt.Println(string(bra))
		for {
			line := []byte{}
			i3line := I3Line{}
			for {
				b, _ := r.ReadByte()
				line = append(line, b)
				if b == ']' {
					if err := json.Unmarshal(line, &i3line); err == nil {
						break
					} else {
						log.Fatalf("Error: %v\n%s\n", err, string(line))
					}
				}
			}
			processLine(i3line)
			rc, _ := r.ReadByte()
			if rc != '\n' {
				log.Fatalf("Missing rc between records")
			}
			comma, _ := r.ReadByte()
			if comma != ',' {
				log.Fatalf("Missing comma between records")
			}
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func processLine(line I3Line) {
	l := completeLine(line)
	buf, _ := json.Marshal(l)
	fmt.Print(string(buf) + "\n,")
}

func completeLine(line I3Line) I3Line {
	newLine := I3Line{}

	if kubeCurrentContext != "" || kubeCurrentNamespace != "" {

		newLine = append(newLine, I3Element{
			Name:     "kube",
			Markup:   "none",
			Color:    "#3030FF",
			FullText: "☸ " + kubeCurrentContext + "/" + kubeCurrentNamespace + " ☸",
		})
	}
	newLine = append(newLine, line...)
	return newLine
}

package main

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/itchyny/volume-go"
)

const (
	defaultTrueValue = "spaudio_yes"
	internalName     = "Built-in Output"
)

type output struct {
	SPAudioDataType []struct {
		Items []struct {
			Name    string `json:"_name"`
			Default string `json:"coreaudio_default_audio_output_device"`
		} `json:"_items"`
	}
}

func main() {
	var lastDevice string

	for {
		current, err := getCurrentAudioDevice()
		if err != nil {
			log.Fatal(err)
		}

		if lastDevice != current {
			lastDevice = current
			if current == internalName {
				if err := volume.Mute(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func getCurrentAudioDevice() (string, error) {
	cmd := exec.Command("system_profiler", "SPAudioDataType", "-json")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var output output
	json.Unmarshal(out, &output)

	for _, dataType := range output.SPAudioDataType {
		for _, item := range dataType.Items {
			if item.Default == defaultTrueValue {
				return strings.TrimSpace(item.Name), nil
			}
		}
	}

	return "", errors.New("no audio device found")
}

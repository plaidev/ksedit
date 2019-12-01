package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/urfave/cli"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// DecodedSecret ...
type DecodedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// fix map[string][]bytes -> map[string]string to decode base64 secret data.
	Data       map[string]string `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`
	StringData map[string]string `json:"stringData,omitempty" protobuf:"bytes,4,rep,name=stringData"`
	Type       apiv1.SecretType  `json:"type,omitempty" protobuf:"bytes,3,opt,name=type,casttype=SecretType"`
}

func encodeSecret(secret *DecodedSecret) {
	for key, data := range secret.Data {
		decoded := base64.StdEncoding.EncodeToString([]byte(data))
		secret.Data[key] = decoded
	}
}

func decodeSecret(secret *DecodedSecret) {
	for key, data := range secret.Data {
		decoded, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			panic(err)
		}
		secret.Data[key] = string(decoded)
	}
}

func editSecretWithEditor(secret *DecodedSecret, editor string) {
	tempFile, err := ioutil.TempFile("", "secret")
	if err != nil {
		panic(err)
	}
	bytes, err := yaml.Marshal(secret)
	if err != nil {
		panic(err)
	}
	tempFile.Write(bytes)
	tempFile.Close()

	c := exec.Command(editor, tempFile.Name())
	tty, err := os.Open("/dev/tty")
	if err != nil {
		panic(err)
	}
	defer tty.Close()
	c.Stdin = tty
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		panic(err)
	}

	edited, err := ioutil.ReadFile(tempFile.Name())

	*secret = *&DecodedSecret{}
	err = yaml.Unmarshal(edited, secret)
	if err != nil {
		panic(err)
	}
}

func readSecretYml(filepath string, secret *DecodedSecret) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, secret)
	if err != nil {
		panic(err)
	}
}

func readSecretYmlFromStdin(secret *DecodedSecret) {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, secret)
	if err != nil {
		panic(err)
	}
}

func main() {
	app := cli.NewApp()

	app.Name = "ksedit"
	app.Usage = "kubernetest secret resource edit"
	app.UsageText = "ksedit [global options] filepath"
	app.Description = ""
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "write, w",
			Usage: "write secret",
		},
		cli.BoolFlag{
			Name:  "encode, e",
			Usage: "encode secret",
		},
		cli.BoolFlag{
			Name:  "decode, d",
			Usage: "decode secret",
		},
		cli.StringFlag{
			Name:  "editor",
			Usage: "editor",
			Value: "vim",
		},
	}

	app.Action = func(context *cli.Context) error {
		filepath := context.Args().Get(0)
		writeOpt := context.GlobalBool("write")
		encodeOpt := context.GlobalBool("encode")
		decodeOpt := context.GlobalBool("decode")
		editor := context.GlobalString("editor")

		secret := &DecodedSecret{}

		if filepath != "" {
			readSecretYml(filepath, secret)
		} else {
			readSecretYmlFromStdin(secret)
		}

		if encodeOpt == true {
			encodeSecret(secret)
		} else if decodeOpt == true {
			decodeSecret(secret)
		} else {
			decodeSecret(secret)
			editSecretWithEditor(secret, editor)
			encodeSecret(secret)
		}

		bytes, err := yaml.Marshal(secret)
		if err != nil {
			panic(err)
		}

		if writeOpt == true && filepath != "" {
			info, err := os.Stat(filepath)
			if err != nil {
				panic(err)
			}
			ioutil.WriteFile(filepath, bytes, info.Mode())
			return nil
		}

		fmt.Println(string(bytes))
		return nil
	}

	app.Run(os.Args)
}

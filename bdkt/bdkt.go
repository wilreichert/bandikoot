package main

import (
  "bufio"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "text/template"

  "gopkg.in/yaml.v2"
)



var projects = []string{"keystone", "glance"}

func load_config(file string) (map[interface{}]interface{}, error) {
  data, err := ioutil.ReadFile(file)
  if err != nil {
    return nil, err
  }
  m := make(map[interface{}]interface{})
  err = yaml.Unmarshal([]byte(data), &m)
  if err != nil {
    return nil, err
  }
  fmt.Printf("config:\n%v\n\n", m)
  return m, nil
}

func render_dockerfile(dir string, c map[interface{}]interface{}) error {
  pattern := filepath.Join(dir, "Dockerfile.tmpl")
  t := template.Must(template.ParseGlob(pattern))

  f, err := os.Create(filepath.Join(dir, "Dockerfile"))
  if err != nil {
    return err
  }
  w := bufio.NewWriter(f)
  err = t.Execute(w, c)
  if err != nil {
    return err
  }
  w.Flush()
  return nil
}

// main
func main() {
	config_file := flag.String("config", "config.yaml", "location of configuration yaml file")
  flag.Parse()

	config, err := load_config(*config_file)
  if err != nil {
    log.Fatalf("error: %v", err)
    os.Exit(1)
  }

  for _, project := range projects {
    err = render_dockerfile(project, config)
    if err != nil {
      log.Fatalf("error: %v", err)
      os.Exit(1)
    }
  }
}

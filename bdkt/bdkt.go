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

// data to populate dockerfile templates
type S struct {
  Values map[string]interface{}
}

var sep = "--------------------------------"
var services = []string{"keystone", "glance"}
var debug *bool

func check(e error) {
  if e != nil {
    log.Fatalf("error: %v", e)
    panic(e)
  }
}

func global_config(file string) (map[interface{}]interface{}, error) {
  data, err := ioutil.ReadFile(file)
  if err != nil {
    return nil, err
  }
  m := make(map[interface{}]interface{})
  err = yaml.Unmarshal([]byte(data), &m)
  if err != nil {
    return nil, err
  }

  if *debug {
    fmt.Printf("\nDEBUG Global Config\n%s\n%v\n\n", sep, m)
  }
  return m, nil
}

func service_config(name string, c map[interface{}]interface{}) S {
  c1 := c["openstack"].(map[interface{}]interface{})
  c2 := c[name].(map[interface{}]interface{})

  s := S{
     Values: make(map[string]interface{}),
  }

  s.Values["name"] = name
  s.Values["packages"] = []string{}

  if c["packages"] != nil {
    for _, x := range c["packages"].([]interface{}) {
      s.Values["packages"] = append(s.Values["packages"].([]string), x.(string))
    }
  }
  if c1["packages"] != nil {
    for _, x := range c1["packages"].([]interface{}) {
      s.Values["packages"] = append(s.Values["packages"].([]string), x.(string))
    }
  }
  if c2["packages"] != nil {
    for _, x := range c2["packages"].([]interface{}) {
      s.Values["packages"] = append(s.Values["packages"].([]string), x.(string))
    }
  }
  s.Values["build_packages"] = []string{}
  if c["build_packages"] != nil {
    for _, x := range c["build_packages"].([]interface{}) {
      s.Values["build_packages"] = append(s.Values["build_packages"].([]string), x.(string))
    }
  }

  if c2["base_image"] != nil {
    s.Values["base_image"] = c2["base_image"]
  } else {
    s.Values["base_image"] = c1["base_image"]
  }
  if c2["install"] != nil {
    s.Values["install"] = c2["install"]
  } else {
    s.Values["install"] = c1["install"]
  }
  if c2["location"] != nil {
    s.Values["location"] = c2["location"]
  } else {
    s.Values["location"] = string(c1["location"].(string) + "/" + name)
  }
  if c2["branch"] != nil {
    s.Values["branch"]= c2["branch"]
  } else if c1["branch"] != nil {
    s.Values["branch"] = c1["branch"]
  }

  if *debug {
    fmt.Printf("\nDEBUG Service Config\n%s\n%v\n\n", sep, s)
  }
  return s
}

func render_dockerfile(s S) error {
  pattern := filepath.Join(s.Values["name"].(string), "Dockerfile.tmpl")
  t := template.Must(template.ParseGlob(pattern))

  f, err := os.Create(filepath.Join(s.Values["name"].(string), "Dockerfile"))
  if err != nil {
    return err
  }
  w := bufio.NewWriter(f)
  err = t.Execute(w, s)
  if err != nil {
    return err
  }
  w.Flush()

  if *debug {
    dat, err := ioutil.ReadFile(filepath.Join(s.Values["name"].(string), "Dockerfile"))
    check(err)
    fmt.Printf("\nDEBUG %s Dockerfile\n%s\n%v\n\n", s.Values["name"], sep, string(dat))
  }
  return nil
}

func main() {
	config_file := flag.String("config", "config.yaml", "location of configuration yaml file")
  service_name := flag.String("service", "all", "service to build")
  debug = flag.Bool("debug", false, "add debug info")
  flag.Parse()

  config, err := global_config(*config_file)
  check(err)

  if *service_name == "all" {
    for _, p := range services {
      svc := service_config(p, config)
      err = render_dockerfile(svc)
      check(err)
    }
  } else {
    svc := service_config(*service_name, config)
    err = render_dockerfile(svc)
    check(err)
  }
}

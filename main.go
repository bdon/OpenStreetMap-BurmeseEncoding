package main

import (
  "fmt"
  "os"
  "github.com/paulmach/osm"
  "github.com/paulmach/osm/osmpbf"
  "github.com/google/myanmar-tools/clients/go"
  "context"
  "runtime"
  "time"
)

func main() {
  start := time.Now()
  file, err := os.Open("./myanmar.osm.pbf")
  if err != nil {
    panic(err)
  }
  defer file.Close()

  zgDetector := myanmartools.NewZawgyiDetector()

  scanner := osmpbf.New(context.Background(), file, runtime.GOMAXPROCS(-1))
  defer scanner.Close()

  for scanner.Scan() {
    switch o := scanner.Object().(type) {
    case *osm.Node:
      for k, v := range o.Tags.Map() {
          score := zgDetector.GetZawgyiProbability(v)
          if score > 0.5 {
            fmt.Printf("  Tag: %s = %s\n", k, v)
            // fmt.Println(score)
          }
      }
    case *osm.Way:
      for k, v := range o.Tags.Map() {
          score := zgDetector.GetZawgyiProbability(v)
          if score > 0.5 {
            fmt.Printf("  Tag: %s = %s\n", k, v)
            //fmt.Println(score)
          }
      }
    case *osm.Relation:
      for k, v := range o.Tags.Map() {
          score := zgDetector.GetZawgyiProbability(v)
          if score > 0.5 {
            fmt.Printf("  Tag: %s = %s\n", k, v)
            //fmt.Println(score)
          }
      }
    }
  }

  elapsed := time.Since(start)
  fmt.Printf("took %s", elapsed)

  if err := scanner.Err(); err != nil {
    panic(err)
  }
}
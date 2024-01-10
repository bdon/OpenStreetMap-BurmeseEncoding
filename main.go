package main

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/google/myanmar-tools/clients/go"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"os"
	"runtime"
	"strconv"
	"time"
	"unicode/utf8"
)

// interface Result struct {
//   Datestamp string
//   Threshold float32
//   LikelyZawgyiCount int64
//   HasBurmeseCount int64
// }

func hasBurmeseCodepoint(s string) bool {
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r >= '\u1000' && r <= '\u109F' {
			return true
		}
		i += size
	}
	return false
}

func main() {
	start := time.Now()

	// open the OSM PBF data
	file, err := os.Open("./myanmar.osm.pbf")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// open the output writer
	output, err := os.Create("output.csv.gz")
	if err != nil {
		panic(err)
	}
	defer output.Close()

	gw := gzip.NewWriter(output)
	defer gw.Close()
	csvw := csv.NewWriter(gw)

	zgDetector := myanmartools.NewZawgyiDetector()

	hasBurmeseCount := 0
	likelyZawgyiCount := 0

	rows := make([][]string, 0)

	// osm_type osm_id key text suggested_text

	processTags := func(osmType string, osmId int64, tags map[string]string) {
		for key, v := range tags {
			if hasBurmeseCodepoint(v) {
				score := zgDetector.GetZawgyiProbability(v)
				hasBurmeseCount += 1
				if score > 0.8 {
					likelyZawgyiCount += 1
					rows = append(rows, []string{fmt.Sprintf("%.2f", score), osmType, strconv.FormatInt(osmId, 10), key, v, Zg2uni(v)})
				}
			}
		}
	}

	scanner := osmpbf.New(context.Background(), file, runtime.GOMAXPROCS(-1))
	defer scanner.Close()

	for scanner.Scan() {
		switch o := scanner.Object().(type) {
		case *osm.Node:
			processTags("node", o.ObjectID().Ref(), o.Tags.Map())
		case *osm.Way:
			processTags("way", o.ObjectID().Ref(), o.Tags.Map())
		case *osm.Relation:
			processTags("relation", o.ObjectID().Ref(), o.Tags.Map())
		}
	}

	elapsed := time.Since(start)
	fmt.Println("took %s", elapsed)
	fmt.Println(likelyZawgyiCount, hasBurmeseCount, float32(likelyZawgyiCount)/float32(hasBurmeseCount)*100)

	for _, row := range rows {
		if err := csvw.Write(row); err != nil {
			panic(err)
		}
	}

	csvw.Flush()
	if err := csvw.Error(); err != nil {
		panic(err)
	}

	// number of strings with burmese text: ABCD
	// number of those strings

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

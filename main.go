package main

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/google/myanmar-tools/clients/go"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"os"
	"runtime"
	"sort"
	"strconv"
	"time"
	"unicode/utf8"
)

type Results struct {
	OsmReplicationTimestamp string
	OsmReplicationSeqnum    uint64
	Threshold               float64
	LikelyZawgyiCount       int
	HasBurmeseCount         int
	ElapsedSeconds          float64
}

type Row struct {
	Score     float64
	OsmType   int
	OsmId     int64
	Key       string
	Value     string
	Suggested string
}

func osmTypeString(i int) string {
	if i == 0 {
		return "node"
	} else if i == 1 {
		return "way"
	} else {
		return "relation"
	}
}

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
	threshold := 0.8
	start := time.Now()

	// open the OSM PBF data
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// open the output writer
	output, err := os.Create(os.Args[2])
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

	rows := make([]Row, 0)

	// osm_type osm_id key text suggested_text

	processTags := func(osmType int, osmId int64, tags map[string]string) {
		for key, v := range tags {
			if hasBurmeseCodepoint(v) {
				score := zgDetector.GetZawgyiProbability(v)
				hasBurmeseCount += 1
				if score > threshold {
					likelyZawgyiCount += 1
					rows = append(rows, Row{
						Score:     score,
						OsmType:   osmType,
						OsmId:     osmId,
						Key:       key,
						Value:     v,
						Suggested: Zg2uni(v),
					})
				}
			}
		}
	}

	scanner := osmpbf.New(context.Background(), file, runtime.GOMAXPROCS(-1))

	defer scanner.Close()

	for scanner.Scan() {
		switch o := scanner.Object().(type) {
		case *osm.Node:
			processTags(0, o.ObjectID().Ref(), o.Tags.Map())
		case *osm.Way:
			processTags(1, o.ObjectID().Ref(), o.Tags.Map())
		case *osm.Relation:
			processTags(2, o.ObjectID().Ref(), o.Tags.Map())
		}
	}

	// sort by score > osm type > osm id > key
	sort.Slice(rows, func(i int, j int) bool {
		if rows[i].Score != rows[j].Score {
			return rows[i].Score > rows[j].Score
		}
		if rows[i].OsmType != rows[j].OsmType {
			return rows[i].OsmType < rows[j].OsmType
		}
		if rows[i].OsmId != rows[j].OsmId {
			return rows[i].OsmId < rows[j].OsmId
		}
		return rows[i].Key < rows[j].Key
	})

	for _, row := range rows {
		row_str := []string{fmt.Sprintf("%.2f", row.Score), osmTypeString(row.OsmType), strconv.FormatInt(row.OsmId, 10), row.Key, row.Value, row.Suggested}
		if err := csvw.Write(row_str); err != nil {
			panic(err)
		}
	}

	csvw.Flush()
	if err := csvw.Error(); err != nil {
		panic(err)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	header, _ := scanner.Header()
	results := Results{
		Threshold:               threshold,
		ElapsedSeconds:          time.Since(start).Seconds(),
		OsmReplicationTimestamp: header.ReplicationTimestamp.Format(time.RFC3339),
		OsmReplicationSeqnum:    header.ReplicationSeqNum,
		LikelyZawgyiCount:       likelyZawgyiCount,
		HasBurmeseCount:         hasBurmeseCount,
	}
	jsonData, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(jsonData))
}

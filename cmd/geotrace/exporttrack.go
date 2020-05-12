package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

type gpxContainer struct {
	XMLName xml.Name
	Track   gpxTrack `xml:"trk"`
}

type gpxTrackPoint struct {
	Lat  float64 `xml:"lat,attr"`
	Lon  float64 `xml:"lon,attr"`
	Time string  `xml:"time"`
}

type gpxTrackSegment struct {
	Points []gpxTrackPoint `xml:"trkpt"`
}

type gpxTrack struct {
	Name     string            `xml:"name"`
	Segments []gpxTrackSegment `xml:"trkseg"`
}

func generateExportTrackCmd() *Command {
	var sqlitePath string
	var startTime string
	var endTime string
	cmd := &cobra.Command{
		Use: "export-track",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			if sqlitePath == "" {
				return fmt.Errorf("specify a database file")
			}
			if startTime == "" {
				return fmt.Errorf("specify a start time")
			}
			if endTime == "" {
				return fmt.Errorf("specify a end time")
			}
			start, err := time.Parse(time.RFC3339, startTime)
			if err != nil {
				return err
			}
			end, err := time.Parse(time.RFC3339, endTime)
			if err != nil {
				return err
			}
			start = start.In(time.UTC)
			end = end.In(time.UTC)
			db, err := sql.Open("sqlite3", sqlitePath+"?mode=ro")
			if err != nil {
				return err
			}
			defer db.Close()
			result, err := db.QueryContext(ctx, "SELECT time, lat, lon FROM traces WHERE time >= ? AND time <= ?", start.Format(time.RFC3339), end.Format(time.RFC3339))
			if err != nil {
				return err
			}
			defer result.Close()
			track := gpxTrack{}
			segment := gpxTrackSegment{
				Points: make([]gpxTrackPoint, 0, 10),
			}
			for result.Next() {
				var time string
				var lat float64
				var lon float64
				if err := result.Scan(&time, &lat, &lon); err != nil {
					return err
				}
				segment.Points = append(segment.Points, gpxTrackPoint{
					Time: time,
					Lat:  lat,
					Lon:  lon,
				})
			}
			track.Segments = append(track.Segments, segment)
			doc := gpxContainer{Track: track}
			doc.XMLName.Local = "gpx"
			doc.XMLName.Space = "http://www.topografix.com/GPX/1/1"
			xml.NewEncoder(os.Stdout).Encode(doc)
			return nil
		},
	}
	cmd.Flags().StringVar(&sqlitePath, "sqlite-store", "", "Path to sqlite file")
	cmd.Flags().StringVar(&startTime, "start-time", "", "")
	cmd.Flags().StringVar(&endTime, "end-time", "", "")
	return &Command{cmd}
}

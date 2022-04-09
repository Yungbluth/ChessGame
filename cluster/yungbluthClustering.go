/*
yungbluthClustering.go
CSc 372
Project part 2
Due: November 23rd 2020 at the beginning of class
This program reads a path to a text file from user input and prints out the clusters when stable
The text file must have the format: k on the first line, n on the second line, and n x,y coordinate pairs on following lines
where k is the number of clusters and n is the number of coordinate pairs.

Link to Google Doc: https://docs.google.com/document/d/1apVQn3P4VnPpFiEYju66s0_Hh4uxHYkPL8PGI9-Y-Cw/edit?usp=sharing
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/* This function calculates and returns the distance between two points */
func distance(x1 int, y1 int, x2 float64, y2 float64) float64 {
	return ((float64(x1) - x2) * (float64(x1) - x2)) + ((float64(y1) - y2) * (float64(y1) - y2))
}

/* This function checks if two 2d slices are equivalent, and returns a boolean */
func equalCheck(a [][]int, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	input := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the path to the text file: ")
	text, _ := input.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	file, err := os.Open(text)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var k, n, i int
	var plot [][]int
	i = 0
	for scanner.Scan() {
		i++
		if i == 1 {
			/* set k */
			k, err = strconv.Atoi(scanner.Text())
		}
		if i == 2 {
			/* set n and initialize the empty array */
			n, err = strconv.Atoi(scanner.Text())
			plot = make([][]int, n)
			for j := range plot {
				plot[j] = make([]int, n)
			}
			if k > n {
				fmt.Println("k > n, next time don't be a bad boy")
				os.Exit(1)
			}
		}
		if i > 2 {
			/* fill the array and do some splitting and conversion from string to int */
			splitter := strings.Split(scanner.Text(), " ")
			strArr := make([]int, len(splitter))
			for j := range strArr {
				strArr[j], _ = strconv.Atoi(splitter[j])
			}
			plot[(i - 3)] = strArr
		}
	}
	/* setup the clusters to be the first k elements */
	clusters := make([][]float64, k)
	for i := range clusters {
		clusters[i] = append(clusters[i], float64(plot[i][0]))
		clusters[i] = append(clusters[i], float64(plot[i][1]))
	}
	iterations := 0
	notStable := true
	centroidPoints := make([][]int, k)
	oldCentroid := make([][]int, k)
	for notStable == true {
		iterations++
		oldCentroid = centroidPoints
		centroidPoints = make([][]int, k)
		for points := range plot {
			minDistance := distance(plot[points][0], plot[points][1], clusters[0][0], clusters[0][1])
			j := 0
			for i := range clusters {
				curDistance := distance(plot[points][0], plot[points][1], clusters[i][0], clusters[i][1])
				if curDistance < minDistance {
					minDistance = curDistance
					j = i
				}
			}
			centroidPoints[j] = append(centroidPoints[j], points)
		}
		if equalCheck(centroidPoints, oldCentroid) == true {
			/* Centroids haven't changed, we are done */
			notStable = false
		} else {
			/* recompute centroid locations, put them in clusters*/
			for i := range centroidPoints {
				length := float64(len(centroidPoints[i]))
				var sumX, sumY float64
				for j := range centroidPoints[i] {
					/* each point in the cluster is centroidPoints[i][j] */
					sumX += float64(plot[centroidPoints[i][j]][0])
					sumY += float64(plot[centroidPoints[i][j]][1])
				}
				xPos := sumX / length
				yPos := sumY / length
				clusters[i][0] = xPos
				clusters[i][1] = yPos
			}
			/* 	centroidPoints is the i value of each point in each cluster (ex. [[0] [1 2 3 4]])
			plot is the list of each point in order (ex. [[3 4] [3 3] [1 2] [4 2] [3 1]])
			clusters is the current location of the clusters (ex. [[3 4] [3 3]])*/
		}
	}
	fmt.Println("The final centroid locations are:\n")
	for i := range clusters {
		/* the Sprintf and TrimRights are just for formatting to get rid of trailing 0 and . where applicable to make things pretty
		   In the formatting I also round to 3 decimal points and round */
		xVal := fmt.Sprintf("%.3f", clusters[i][0])
		xVal = strings.TrimRight(strings.TrimRight(xVal, "0"), ".")
		yVal := fmt.Sprintf("%.3f", clusters[i][1])
		yVal = strings.TrimRight(strings.TrimRight(yVal, "0"), ".")
		fmt.Printf("u(%v) = (%v,%v)\n", i, xVal, yVal)
	}
	fmt.Printf("\n%v iterations were required.", iterations)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

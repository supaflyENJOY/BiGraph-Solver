// biGraph-Solver project main.go
package main

import "fmt"

type SolutionData struct {
	sourceData         [][]int
	solutions          [][]int
	temporarySolution  []int
	usedVertices       []bool
	resultCount        int
	maxSecondVertice   int
	secondVerticesData map[int]int
}

func countSecondVerticesRelations(solution *SolutionData) {
	solution.secondVerticesData = make(map[int]int)
	for _, v := range solution.sourceData {
		for _, v2 := range v {
			solution.secondVerticesData[v2]++
			if v2 > solution.maxSecondVertice {
				solution.maxSecondVertice = v2
			}
		}
	}
}

func iterateOver(v int, solution *SolutionData) {
	fmt.Println(v, len(solution.sourceData))
	if v >= len(solution.sourceData) {
		//fmt.Println(len(solution.temporarySolution), solution.resultCount)
		tCount := len(solution.temporarySolution)
		if tCount >= solution.resultCount {
			if tCount > solution.resultCount {
				solution.solutions = nil
				solution.resultCount = tCount
			}
			ts := make([]int, len(solution.temporarySolution))
			copy(ts, solution.temporarySolution)
			solution.solutions = append(solution.solutions, ts)
		}
		return
	}
	if len(solution.sourceData[v]) == 0 {
		iterateOver(v+1, solution)
	}

	for _, vertice := range solution.sourceData[v] {
		if solution.usedVertices[vertice] {
			continue
		}
		solution.usedVertices[vertice] = true
		tempLen := len(solution.temporarySolution)
		solution.temporarySolution = append(solution.temporarySolution, v, vertice)
		iterateOver(v+1, solution)
		solution.temporarySolution = solution.temporarySolution[:tempLen]
		solution.usedVertices[vertice] = false
	}

	iterateOver(v+1, solution)
}

func GetSolutions(data *[][]int) *[][]int {
	solution := new(SolutionData)
	solution.sourceData = *data
	solution.solutions = nil
	solution.temporarySolution = make([]int, 0, len(*data)*2)
	countSecondVerticesRelations(solution)
	solution.usedVertices = make([]bool, solution.maxSecondVertice+1)
	iterateOver(0, solution)
	return &solution.solutions
}

/*func main() {
	var n, m, k int
	//fmt.Print("Write n, m: ")
	fmt.Scan(&n, &m)
	var data = make([][]int, n)
	for i := 0; i < n; i++ {
		//fmt.Printf("Number of relations for vertice %v: ", i+1)
		fmt.Scan(&k)
		data[i] = make([]int, k)
		//fmt.Printf("Write relations %v: ", i+1)
		for j := 0; j < k; j++ {
			fmt.Scan(&data[i][j])
		}
	}
	//fmt.Println(data)
	solutions := GetSolutions(&data)
	fmt.Println(solutions)

}*/

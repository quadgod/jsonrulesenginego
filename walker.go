package pathresolver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type pathWalkerNode struct {
	stringValue string
	indexValue  int
}

type pathPathWalker struct {
	currentIndex int
	lastIndex    int
	nodes        []pathWalkerNode
}

type indexPair struct {
	start int
	end   int
}

func (w *pathPathWalker) IsCurrentNodeIndex() bool {
	if w.currentIndex == -1 {
		return false
	}

	return w.nodes[w.currentIndex].stringValue == ""
}

func (w *pathPathWalker) IndexValue() int {
	if w.currentIndex == -1 {
		return -1
	}

	return w.nodes[w.currentIndex].indexValue
}

func (w *pathPathWalker) StringValue() string {
	if w.currentIndex == -1 {
		return ""
	}

	return w.nodes[w.currentIndex].stringValue
}

func (w *pathPathWalker) MoveToNextNode() bool {
	if w.currentIndex < w.lastIndex {
		w.currentIndex++
		return true
	}

	return false
}

func newPathWalker(path string) (*pathPathWalker, error) {
	pathWithoutSpaces := strings.TrimSpace(path)

	if pathWithoutSpaces == "" {
		return nil, errors.New("path can't be empty")
	}

	parts := strings.Split(pathWithoutSpaces, ".")
	partsLen := len(parts)

	pw := &pathPathWalker{
		currentIndex: -1,
		lastIndex:    -1,
		nodes:        make([]pathWalkerNode, 0),
	}

	for i := 0; i < partsLen; i++ {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			return nil, fmt.Errorf("invalid path. path \"%s\"", path)
		}

		pairs, err := extractIndexPairs(&part)
		if err != nil {
			return nil, fmt.Errorf("invalid path. path \"%s\". invalid part \"%s\"", path, part)
		}

		isContainsIndexPairs := len(pairs) > 0

		if part[0] != '[' {
			if isContainsIndexPairs {
				name := part[0:pairs[0].start]
				pw.nodes = append(pw.nodes, pathWalkerNode{
					indexValue:  -1,
					stringValue: name,
				})
			} else {
				pw.nodes = append(pw.nodes, pathWalkerNode{
					indexValue:  -1,
					stringValue: part,
				})
			}
		}

		if isContainsIndexPairs {
			for _, idxPair := range pairs {
				pairStrValue := part[idxPair.start+1 : idxPair.end]
				index, err := strconv.Atoi(pairStrValue)
				if err != nil || index < 0 {
					return nil, fmt.Errorf("invalid path. path \"%s\". invalid part \"%s\"", path, part)
				}

				pw.nodes = append(pw.nodes, pathWalkerNode{
					indexValue:  index,
					stringValue: "",
				})
			}
		}
	}

	nodesCount := len(pw.nodes)

	if nodesCount > 0 {
		pw.currentIndex = 0
		pw.lastIndex = nodesCount - 1
	}

	return pw, nil
}

func extractIndexPairs(pathPart *string) ([]indexPair, error) {
	openIndexes := make([]int, 0)
	closeIndexes := make([]int, 0)

	for idx, prune := range *pathPart {
		if prune == '[' {
			openIndexes = append(openIndexes, idx)
		} else if prune == ']' {
			closeIndexes = append(closeIndexes, idx)
		}
	}

	openIndexesLen := len(openIndexes)
	closeIndexesLen := len(closeIndexes)

	if openIndexesLen != closeIndexesLen {
		return nil, fmt.Errorf("invalid path part \"%s\"", *pathPart)
	}

	if closeIndexesLen > 0 && (*pathPart)[len(*pathPart)-1] != ']' {
		return nil, fmt.Errorf("invalid path part \"%s\"", *pathPart)
	}

	checkIndex := 0
	result := make([]indexPair, 0)

	for {
		if checkIndex == openIndexesLen {
			break
		}

		openIndex := openIndexes[checkIndex]
		closeIndex := closeIndexes[checkIndex]

		if closeIndex <= openIndex || closeIndex-openIndex <= 0 {
			return nil, fmt.Errorf("invalid path part \"%s\"", *pathPart)
		}

		result = append(result, indexPair{
			start: openIndex,
			end:   closeIndex,
		})

		checkIndex++
	}

	return result, nil
}

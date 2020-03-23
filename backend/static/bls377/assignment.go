/*
Copyright © 2020 ConsenSys

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by gnark/internal/templates/generator DO NOT EDIT

package backend

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/consensys/gurvy/bls377/fr"

	"github.com/consensys/gnark/backend"
)

// Assignment is used to specify inputs to the Prove and Verify functions
type Assignment struct {
	Value    fr.Element
	IsPublic bool // default == false (assignemnt is private)
}

// Assignments is used to specify inputs to the Prove and Verify functions
type Assignments map[string]Assignment

// NewAssignment returns an empty Assigments object
func NewAssignment() Assignments {
	return make(Assignments)
}

// Assign assign a value to a Secret/Public input identified by its name
func (a Assignments) Assign(visibility backend.Visibility, name string, v interface{}) {
	if _, ok := a[name]; ok {
		panic(name + " already assigned")
	}
	switch visibility {
	case backend.Secret:
		a[name] = Assignment{Value: fr.FromInterface(v)}
	case backend.Public:
		a[name] = Assignment{
			Value:    fr.FromInterface(v),
			IsPublic: true,
		}
	default:
		panic("supported visibility attributes are SECRET and PUBLIC")
	}
}

// Read parse r1cs.Assigments from given file
// file line structure: secret/public, assignmentName, assignmentValue
// note this is a cs/ subpackage because we need to instantiate internal/fr.Element
func (assigment Assignments) Read(filePath string) error {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if len(line) != 3 {
			return backend.ErrInvalidInputFormat
		}
		visibility := strings.ToLower(strings.TrimSpace(line[0]))
		name := strings.TrimSpace(line[1])
		value := strings.TrimSpace(line[2])

		assigment.Assign(backend.Visibility(visibility), name, value)
	}
	return nil
}

// Write serialize given assigment to disk
// file line structure: secret/public, assignmentName, assignmentValue
func (assignment Assignments) Write(path string) error {
	csvFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	for k, v := range assignment {
		r := v.Value
		record := []string{string(backend.Secret), k, r.String()}
		if v.IsPublic {
			record[0] = string(backend.Public)
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

// DiscardSecrets returns a copy of self without Secret Assigment
func (assignments Assignments) DiscardSecrets() Assignments {
	toReturn := NewAssignment()
	for k, v := range assignments {
		if v.IsPublic {
			toReturn[k] = v
		}
	}
	return toReturn
}

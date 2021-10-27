package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Student struct {
	ID string `json:"name"`
	Subject string `json:"subject"`
	Note    float64 `json:"note"`
}

var students map[string]map[string]float64

func Get() ([]byte, error) {
	jsonData, err := json.MarshalIndent(students, "", "    ") 
	if err != nil {
		return jsonData, nil
	}
	return jsonData, err
}

func GetID(id string) ([]byte, error) {
	jsonData := []byte(`{}`)
	student, ok := students[id]
	if !ok {
		return jsonData, nil
	}
	jsonData, err := json.MarshalIndent(student, "", "    ")
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}

func Add(student Student) ([] byte) {
	jsonData := []byte(`{code: "OK"}`)
	subject := make(map[string]float64)
	subject[student.Subject] = student.Note
	_, ok := students[student.ID]
	if ok {
		students[student.ID][student.Subject] = student.Note
	} else {
		students[student.ID] = subject
	}
	return jsonData
}

func Delete(id string) ([] byte) {
	_, ok := students[id]
	if !ok {
		return []byte(`{code: "Undefined"}`)
	}
	delete(students, id)
	return []byte(`{code: "OK"}`)
}

func Update(id string, student Student) ([] byte) {
	_, ok := students[id]
	if !ok {
		return []byte(`{code: "Undefined"}`)
	}
	students[id][student.Subject] = student.Note
	return []byte(`{code: "OK"}`)
}

func student(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
		case "GET": // devuelve a todos los alumnos
			json, err := Get() // obtenemos los alumnos codificados en formato JSON
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return	
			}
			res.Header().Set(
				"Content-Type",
				"application/json",
			)
			res.Write(json) // enviamos la respuesta
		case "POST": // agregar alumno, materia y calificación
			var student Student
			err := json.NewDecoder(req.Body).Decode(&student) // leemos los datos mandados por POST
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			json := Add(student) // añadimos el nuevo estudiante
			res.Header().Set(
				"Content-Type",
				"application/json",
			)
			res.Write(json) // enviamos la respuesta
	}
}

func studentId(res http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/student/")
	fmt.Println(req.Method)
	switch req.Method {
		case "GET": // devuelve alumno por id
			json, err := GetID(id)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return	
			}
			res.Header().Set(
				"Content-Type",
				"application/json",
			)
			res.Write(json)
		case "PUT": // modificar alumno por id
			var student Student
			err := json.NewDecoder(req.Body).Decode(&student)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			json := Update(id, student)
			res.Header().Set(
				"Content-Type",
				"application/json",
			)
			res.Write(json)
		case "DELETE": // eliminar alumno por id
			json := Delete(id)
			res.Header().Set(
				"Content-Type",
				"application/json",
			)
			res.Write(json)
	}
}

func main() {
	students = make(map[string]map[string]float64)
	http.HandleFunc("/student", student)
	http.HandleFunc("/student/", studentId)
	fmt.Println("Servidor corriendo en el puerto 9000...")
	http.ListenAndServe(":9000", nil)
}
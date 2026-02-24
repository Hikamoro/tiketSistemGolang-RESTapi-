package main

import (
	"encoding/json"
	//"fmt"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var id int

func getId(id *int) int {
	*id += 1
	return *id
}

type Issue struct {
	Subject  string
	Text     string
	Priority string
}

type Respon struct {
	Err  string
	Code int
}

var Issues = make(map[int]Issue)

// func cheackMethod(m Method) {
// 	if r.Method != http.MethodDelete {
// 		res := Respon{
// 			Err:  "use method delete",
// 			Code: 999,
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(res)
// 		return
// 	}
// }

// Issues = append(Issues, struct{Id:123,Subject:"dadada",Text:"big text", Priority:"high"})
func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, "use get")
		return
	}
	if len(r.URL.Query()) == 0 {
		// fmt.Fprintf(w, strconv.Itoa(getId(&id)))
		// fmt.Fprintf(w, strconv.Itoa(getId(&id)))
		if len(Issues) == 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Respon{Err: "0 запросов", Code: 200})
			//fmt.Fprintln(w, "0 запросов")
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Issues)
			//fmt.Fprintln(w, Issues[i])
		}
	} else {
		wId := r.URL.Query().Get("id")
		num, err := strconv.Atoi(wId)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Respon{Err: "используйте id типа int!!!", Code: 999})
			// fmt.Fprintf(w, "используйте id типа int!!!")
			// log.Println("Ошибка:", err)
			return
		} else if num < 0 {
			//fmt.Fprintf(w, Respon{err:"id не может быть отрицательным!",code:999})
			//fmt.Fprintln(w, Respon{err: "id не может быть отрицательным!", code: 999})
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Respon{Err: "id не может быть отрицательным!", Code: 999})
			return
		}
		if num < len(Issues) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Issues[num])
			//fmt.Fprintln(w, Issues[num])
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Respon{Err: "такого тикета нет", Code: 999})
			//fmt.Fprintln(w, Respon{err: "такого тикета нет", code: 999})
		}

	}
}
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res := Respon{
			Err:  "use method post",
			Code: 999,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}
	var is Issue
	//fmt.Fprintf(w, "post")
	err := json.NewDecoder(r.Body).Decode(&is)
	if err != nil {
		res := Respon{
			Err:  "Ошибка парсинга JSON",
			Code: 999,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		//http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}
	Issues[getId(&id)] = is

	res := Respon{
		Err:  "good",
		Code: 200,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		res := Respon{
			Err:  "use method delete",
			Code: 999,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}
	wId := r.URL.Query().Get("id")
	num, err := strconv.Atoi(wId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Respon{Err: "используйте id типа int!!!", Code: 999})
		return
	}
	_, key := Issues[num]
	if key != false {
		delete(Issues, num)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Respon{Err: "done", Code: 200})
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Respon{Err: "такого id нет", Code: 999})
	}
}
func PutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		res := Respon{
			Err:  "use method Put",
			Code: 999,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}
	wId := r.URL.Query().Get("id")
	num, err := strconv.Atoi(wId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Respon{Err: "используйте id типа int!!!", Code: 999})
		return
	}
	var is Issue
	//fmt.Fprintf(w, "post")
	err = json.NewDecoder(r.Body).Decode(&is)
	if err != nil {
		res := Respon{
			Err:  "Ошибка парсинга JSON",
			Code: 999,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		//http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}
	_, key := Issues[num]
	if key != false {
		Issues[num] = is
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Respon{Err: "done", Code: 200})
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Respon{Err: "такого id нет", Code: 999})
	}
}

func main() {
	Issues[0] = Issue{
		Subject:  "dadada",
		Text:     "big text",
		Priority: "high",
	}
	http.HandleFunc("/get", GetHandler)
	http.HandleFunc("/post", PostHandler)
	http.HandleFunc("/put", PutHandler)
	http.HandleFunc("/delete", DeleteHandler)
	log.Println("Запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}

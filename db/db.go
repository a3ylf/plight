package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// TODO: TOTAL TIME WORKED

type Plight struct {
	mux      sync.RWMutex
	filename string
}

func StartDB(filename ...string) *Plight {
	base := "data.json"
	if len(filename) > 0 {
		base = filename[0]
	}
	return &Plight{
		mux:      sync.RWMutex{},
		filename: base,
	}

}

func (p *Plight) EnsureDB() error {
	name := "data.json"
	if p.filename != "" {
		name = p.filename

	}

	_, err := os.Open(name)

	if err != nil {
		err = createDB(name)
	}

	return err
}

func createDB(filename string) error {
	data := &Data{
		Sessions: make(map[string]Timers), // Garantir que o mapa 'z' esteja inicializado
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, b, 0644)
	return err
}

type Period struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Timers struct {
	Total   string              `json:"total"`
	Periods map[string][]Period `json:"periods"`
}

type Data struct {
	Sessions map[string]Timers `json:"sessions"`
}

func (p *Plight) ReadDB() (*Data, error) {
	b, err := os.ReadFile(p.filename)

	if err != nil {
		return &Data{}, err
	}

	data := &Data{}

	err = json.Unmarshal(b, data)
	if err != nil {
		return &Data{}, err
	}
	return data, nil
}

// //Format to save on file
// // a := time.Now().Format(time.DateTime)
// // fmt.Println(a)
// // b := time.Now().Format(time.DateTime)
// // fmt.Println(b)
// //
// //Retrieve from file
// fmt.Println()
// c , _ := time.Parse(time.DateTime,a)
//    d , _ := time.Parse(time.DateTime,b)
//    e := d.Sub(c).String()
// fmt.Println(e)
//    x, _ := time.ParseDuration(e)
//    fmt.Println(x)

func (p *Plight) WriteDB(session string) error {

	data, err := p.ReadDB()

	if err != nil {
		return err
	}
	timenow := fmt.Sprint(time.Now().Date())

	if data.Sessions == nil {
		data.Sessions = make(map[string]Timers)
	}
	if _, e := data.Sessions[session]; !e {
		data.Sessions[session] = Timers{
			Periods: make(map[string][]Period),

			// xd
			Total: time.Now().Sub(time.Now()).String(),
		}
	}
	last := len(data.Sessions[session].Periods[timenow]) - 1

	if last == -1 {

		data.Sessions[session].Periods[timenow] = append(data.Sessions[session].Periods[timenow],
			Period{From: time.Now().Format(time.TimeOnly)})

	} else if data.Sessions[session].Periods[timenow][last].To == "" {
		now := time.Now().Format(time.TimeOnly)
		data.Sessions[session].Periods[timenow][last].To = now
		dur, err := time.ParseDuration(data.Sessions[session].Total)
        from, err:= time.Parse(time.TimeOnly,data.Sessions[session].Periods[timenow][last].From)
        to, err:= time.Parse(time.TimeOnly,data.Sessions[session].Periods[timenow][last].To)
        if err != nil {
            return err
        }
		newTotal := to.
			Add(dur).
			Sub(from).
			String()

        s := data.Sessions[session]
        s.Total = newTotal
		data.Sessions[session] = s

	} else {
		data.Sessions[session].Periods[timenow] = append(data.Sessions[session].Periods[timenow],
			Period{From: time.Now().Format(time.TimeOnly)})
	}
	d, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}

	p.mux.Lock()
	err = os.WriteFile(p.filename, d, 0644)
	p.mux.Unlock()

	if err != nil {
		return err
	}
	return nil
}

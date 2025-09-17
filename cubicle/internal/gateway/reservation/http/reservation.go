package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"bibliotecaup.com/cubicle/internal/gateway"
	discovery "bibliotecaup.com/pkg/registry"
	"bibliotecaup.com/reservation/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
}

func (g *Gateway) PutReservation(ctx context.Context, recordID model.RecordID, recordType model.RecordType, reservation *model.Reservation) error {
	addrs, err := g.registry.ServiceAddress(ctx, "reservation")
	if err != nil {
		return err
	}

	// Construir la URL del servicio de reservaciones
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/reservation"
	log.Printf("Calling reservation service, request: POST %s", url)

	// Crear la solicitud HTTP POST
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	// Agregar los datos como parámetros de consulta
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", string(recordType))
	values.Add("userId", string(reservation.UserID))
	values.Add("start", reservation.Start.Format(time.RFC3339))
	values.Add("end", reservation.End.Format(time.RFC3339))
	values.Add("status", string(reservation.Status))
	req.URL.RawQuery = values.Encode()

	// Enviar la solicitud
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}

func (g *Gateway) GetAvailability(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (*model.Availability, error) {
	addrs, err := g.registry.ServiceAddress(ctx, "reservation")
	if err != nil {
		return nil, err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/reservation"
	log.Printf("%s", "Calling reservation service, request: GET "+url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("recordId", string(recordID))
	values.Add("recordType", fmt.Sprintf("%v", recordType))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}
	var availability model.Availability
	if err := json.NewDecoder(resp.Body).Decode(&availability); err != nil {
		return nil, err
	}
	return &availability, nil
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// import (
//     "context"
//     "encoding/json"
//     "fmt"
//     "net/http"
//     "time"
// )

// type Availability struct {
//     AvailableNow  bool      `json:"availableNow"`
//     NextAvailable time.Time `json:"nextAvailable"`
// }

// type Gateway struct {
//     url    string
//     client *http.Client
// }

// func New(url string) *Gateway {
//     return &Gateway{
//         url:    url,
//         client: &http.Client{},
//     }
// }

// func (g *Gateway) GetAggregated(ctx context.Context, id string) (*Availability, error) {
//     u := fmt.Sprintf("%s/?recordId=%s&recordType=cubicle", g.url, id)
//     req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
//     if err != nil {
//         return nil, err
//     }

//     resp, err := g.client.Do(req)
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()

//     var availability Availability
//     if err := json.NewDecoder(resp.Body).Decode(&availability); err != nil {
//         return nil, err
//     }

//     return &availability, nil
// }
